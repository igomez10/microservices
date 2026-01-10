package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/exaring/otelpgx"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/igomez10/microservices/socialapp/internal/middlewares/cache"
	"github.com/igomez10/microservices/socialapp/internal/server"
	"github.com/igomez10/microservices/socialapp/pkg/dbpgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jessevdk/go-flags"
	_ "github.com/lib/pq"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// Configuration constants
const (
	DefaultDBMinConns          = 3
	DefaultDBConnTimeout       = 5 * time.Second
	DefaultDBPingTimeout       = 5 * time.Second
	DefaultHTTPTimeout         = 15 * time.Second
	DefaultRetryMax            = 10
	DefaultRedisPoolSize       = 10
	DefaultRedisMinIdleConns   = 10
	DefaultShutdownTimeout     = 30 * time.Second
	DefaultMaxIdleConns        = 100
	DefaultMaxIdleConnsPerHost = 100
	DefaultIdleConnTimeout     = 120 * time.Second
	DefaultTLSHandshakeTimeout = 10 * time.Second
	DefaultResponseTimeout     = 10 * time.Second
	DefaultExpectContinue      = 10 * time.Second
	DefaultNumDBPools          = 5
)

type appOptions struct {
	AppName               string        `short:"n" long:"name" description:"name of the app" default:"socialapp"`
	AppPort               int           `short:"p" long:"port" description:"main port for application" default:"8080" env:"PORT"`
	ProxyHost             string        `short:"x" long:"proxy" description:"proxy url, \"http://localhost:9091\"" env:"HTTP_PROXY"`
	LogLevel              string        `short:"l" long:"logLevel" description:"log level info/error/warning" default:"info" choice:"info" choice:"error" choice:"debug" choice:"warning" env:"LOG_LEVEL"`
	LogHost               string        `long:"logHost" description:"log host url" required:"true" env:"LOGSTASH_HOST"`
	PropertiesSubdomain   string        `long:"propertiesSubdomain" description:"Properties subdomain" required:"true" env:"PROPERTIES_SUBDOMAIN"`
	DefaultTimeout        time.Duration `long:"defaultTimeout" description:"Default timeout for requests" default:"10s"`
	AgentURL              string        `long:"agentURL" description:"Agent URL" required:"true" env:"OTEL_EXPORTER_OTLP_ENDPOINT"`
	DatabaseURL           string        `long:"databaseURL" description:"Database URL" required:"true" env:"DATABASE_URL"`
	PuttyknifeDomain      string        `long:"puttyknifeDomain" description:"Puttyknife domain: puttyknife.{...}.com" required:"true" env:"PUTTYKNIFE_DOMAIN"`
	PuttyknifeURL         string        `long:"puttyknife-url" description:"Puttyknife url" required:"true" env:"PUTTYKNIFE_URL"`
	UrlShortenerSubdomain string        `long:"urlShortenerSubdomain" description:"URL shortener subdomain" env:"URLSHORTENER_SUBDOMAIN"`
	URLShortenerURL       string        `long:"urlShortenerURL" description:"URL shortener URL" required:"true" env:"URLSHORTENER_URL"`
	SocialappSubdomain    string        `long:"socialappSubdomain" description:"Socialapp subdomain" required:"true" env:"SOCIALAPP_SUBDOMAIN"`
	JwtSecret             string        `long:"jwtSecret" description:"jwt secret" required:"true" env:"JWT_SECRET"`
	RedisURL              string        `long:"redisURL" description:"redis url" required:"true" env:"REDIS_URL"`
	KibanaSubdomain       string        `long:"kibanaSubdomain" description:"Kibana subdomain" required:"true" env:"KIBANA_SUBDOMAIN"`
	KibanaURL             string        `long:"kibanaURL" description:"Kibana URL" required:"true" env:"KIBANA_URL"`
	LocalSubdomain        string        `long:"localSubdomain" description:"Local subdomain" required:"true" env:"LOCAL_SUBDOMAIN" default:"google.com"`
}

func main() {
	// Create cancellable context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config, shutdown := getConfig(ctx)
	defer func() {
		slog.Info("Running shutdown functions")
		for _, fn := range shutdown {
			if err := fn(); err != nil {
				slog.Error("failed to shutdown", "error", err)
			}
		}
	}()

	// Set up signal handling for graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Run server in a goroutine
	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- run(ctx, config)
	}()

	// Wait for shutdown signal or server error
	select {
	case sig := <-signalChan:
		slog.Info("Received shutdown signal", "signal", sig.String())
		cancel() // Cancel context to trigger graceful shutdown
	case err := <-serverErrors:
		if err != nil {
			slog.Error("Server error", "error", err)
		}
	}

	slog.Info("Main shutdown complete")
}

func run(ctx context.Context, config server.Config) error {
	// Create the router using the server package
	router, err := server.NewRouter(ctx, config)
	if err != nil {
		return fmt.Errorf("failed to create router: %w", err)
	}

	// Create HTTP server with graceful shutdown support
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.AppPort),
		Handler: router,
	}

	// Start graceful shutdown listener
	go func() {
		<-ctx.Done()
		slog.Info("Context cancelled, starting graceful shutdown")

		// Give server time to finish existing requests
		shutdownCtx, cancel := context.WithTimeout(context.Background(), DefaultShutdownTimeout)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			slog.Error("Server shutdown error", "error", err)
		}
	}()

	slog.Info("Server listening on port", "port", config.AppPort)
	err = httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server error: %w", err)
	}

	slog.Info("Server stopped")
	return nil
}

// CreateDBPools creates a pool of connections to the database
func CreateDBPools(ctx context.Context, databaseURL string, numPools int, applicationName string) *ForcedConnectionPool {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		slog.Error("failed to parse database url", "error", err)
		os.Exit(1)
	}

	config.MinConns = DefaultDBMinConns
	config.ConnConfig.ConnectTimeout = DefaultDBConnTimeout
	config.ConnConfig.RuntimeParams = map[string]string{
		"application_name": applicationName,
	}

	config.ConnConfig.Tracer = otelpgx.NewTracer(otelpgx.WithTracerProvider(otel.GetTracerProvider()))
	pools := make([]*pgxpool.Pool, 0, numPools)
	for i := 0; i < numPools; i++ {
		// Use context.Background() for pool creation so the pool isn't tied to request context lifecycle
		dbConn, err := pgxpool.NewWithConfig(context.Background(), config)
		if err != nil {
			slog.Error("failed to create db pool", "error", err)
			os.Exit(1)
		}

		if dbConn == nil {
			slog.Error("db is nil")
			os.Exit(1)
		}

		pingCtx, cancel := context.WithTimeout(context.Background(), DefaultDBPingTimeout)
		defer cancel()
		if err := dbConn.Ping(pingCtx); err != nil {
			slog.Error("failed to ping database, shutting down", "error", err)
			os.Exit(1)
		}

		pools = append(pools, dbConn)
	}

	f := &ForcedConnectionPool{
		numPools:    numPools,
		connections: pools,
	}

	return f
}

// ForcedConnectionPool is a wrapper around native Go sql.DB, this allows us to force the minimum number of connections
type ForcedConnectionPool struct {
	connections       []*pgxpool.Pool
	numPools          int
	currentRoundRobin atomic.Uint32
}

func (f *ForcedConnectionPool) GetPool() *pgxpool.Pool {
	// round robin with thread-safe atomic operations
	idx := f.currentRoundRobin.Add(1) - 1
	return f.connections[idx%uint32(f.numPools)]
}

func (f *ForcedConnectionPool) Close() {
	for _, conn := range f.connections {
		conn.Close()
	}
}

func getConfig(ctx context.Context) (server.Config, []func() error) {
	var opts appOptions

	shutdown := []func() error{}
	_, err := flags.Parse(&opts)
	if err != nil {
		panic(err)
	}

	instanceID := uuid.NewString()

	http.DefaultClient = buildRetryHTTPClient()
	if opts.ProxyHost != "" {
		if u, err := url.Parse(opts.ProxyHost); err != nil {
			slog.Error("failed to parse proxy URL", "error", err)
		} else {
			http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(u)}
		}
	}

	// Parse log level
	parsedLogLevel, err := parseLogLevel(opts.LogLevel)
	if err != nil {
		slog.Error("Invalid log level", "log_level", opts.LogLevel, "error", err)
		os.Exit(1)
	}

	logger, logDestinations := configureLogger(parsedLogLevel, opts.AppName, instanceID)
	slog.SetDefault(logger)

	slog.Info(
		"Starting SocialApp",
		"logHost", opts.LogHost,
		"logLevel", opts.LogLevel,
		"appName", opts.AppName,
		"propertiesSubdomain", opts.PropertiesSubdomain,
		"urlServiceHost", opts.UrlShortenerSubdomain,
		"agentURL", opts.AgentURL,
		"urlShortenerSubdomain", opts.UrlShortenerSubdomain,
		"urlShortenerURL", opts.URLShortenerURL,
		"socialappSubdomain", opts.SocialappSubdomain,
		"instanceID", instanceID,
		"puttyknifeDomain", opts.PuttyknifeDomain,
		"puttyknifeURL", opts.PuttyknifeURL,
		"redisURL", opts.RedisURL,
		"kibanaSubdomain", opts.KibanaSubdomain,
	)

	// Parse properties subdomain
	var propertiesSubdomainURL *url.URL
	if len(opts.PropertiesSubdomain) != 0 {
		u, err := url.Parse(opts.PropertiesSubdomain)
		if err != nil {
			slog.Error("failed to parse properties subdomain url", "url", opts.PropertiesSubdomain, "error", err)
			os.Exit(1)
		}
		propertiesSubdomainURL = u
	}

	// Parse url service host
	urlService, err := url.Parse(opts.URLShortenerURL)
	if err != nil {
		slog.Error("failed to parse url service host url", "url", opts.URLShortenerURL, "error", err)
		os.Exit(1)
	}

	var puttyknifeURL *url.URL
	if len(opts.PuttyknifeURL) != 0 {
		u, err := url.Parse(opts.PuttyknifeURL)
		if err != nil {
			slog.Error("failed to parse puttyknife url", "url", opts.PuttyknifeURL, "error", err)
			os.Exit(1)
		}
		puttyknifeURL = u
	}

	queries := dbpgx.New()

	// Setup tracing
	_, _, loggerProvider, _, telemetryCleanup, err := setupTelemetry(ctx, &opts, instanceID)
	if err != nil {
		slog.Error("failed to configure telemetry", "error", err)
		os.Exit(1)
	}
	shutdown = append(shutdown, telemetryCleanup...)
	http.DefaultClient.Transport = tracingTransport(http.DefaultClient.Transport, opts.AppName)
	// Connect to database
	connections := CreateDBPools(ctx, opts.DatabaseURL, DefaultNumDBPools, fmt.Sprintf("%s-%s", opts.AppName, instanceID))

	shutdown = append(shutdown, func() error {
		slog.Info("Shutting down database connections")
		connections.Close()
		return nil
	})

	redisOpts, err := redis.ParseURL(opts.RedisURL)
	if err != nil {
		slog.Error("failed to parse redis url", "error", err)
		os.Exit(1)
	}

	redisOpts.PoolSize = DefaultRedisPoolSize
	redisOpts.MinIdleConns = DefaultRedisMinIdleConns

	cacheMiddleware := cache.NewCache(cache.CacheConfig{
		RedisOpts: redisOpts,
	})

	// Read OpenAPI spec file
	openAPIPath := "openapi.yaml"
	openAPIContent, err := os.ReadFile(openAPIPath)
	if err != nil {
		slog.Error("failed to read openapi file", "path", openAPIPath, "error", err)
		os.Exit(1)
	}

	config := server.Config{
		InstanceID:            instanceID,
		AppName:               opts.AppName,
		AppPort:               opts.AppPort,
		LogLevel:              parsedLogLevel,
		LogDestinations:       logDestinations,
		LoggerProvider:        loggerProvider,
		DBPool:                connections.GetPool(),
		Queries:               queries,
		Cache:                 cacheMiddleware,
		PropertiesSubdomain:   propertiesSubdomainURL,
		DefaultTimeout:        opts.DefaultTimeout,
		URLShortenerSubdomain: opts.UrlShortenerSubdomain,
		URLShortenerURL:       urlService,
		SocialappSubdomain:    opts.SocialappSubdomain,
		JWTSecret:             opts.JwtSecret,
		PuttyknifeSubDomain:   opts.PuttyknifeDomain,
		PuttyknifeURL:         puttyknifeURL,
		KibanaSubdomain:       opts.KibanaSubdomain,
		KibanaURL:             opts.KibanaURL,
		LocalSubdomain:        opts.LocalSubdomain,
		OpenAPIContent:        openAPIContent,
	}

	return config, shutdown
}

func parseLogLevel(level string) (slog.Level, error) {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn", "warning":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return slog.LevelInfo, fmt.Errorf("unsupported log level: %s", level)
	}
}

func buildRetryHTTPClient() *http.Client {
	retryClient := retryablehttp.NewClient()
	retryClient.Logger = nil
	retryClient.RequestLogHook = func(_ retryablehttp.Logger, req *http.Request, attempt int) {
		if attempt >= 1 {
			slog.Warn("http retry", "method", req.Method, "url", req.URL.String(), "attempt", attempt)
		}
	}

	retryClient.HTTPClient.Transport = &http.Transport{
		MaxIdleConns:          DefaultMaxIdleConns,
		MaxIdleConnsPerHost:   DefaultMaxIdleConnsPerHost,
		IdleConnTimeout:       DefaultIdleConnTimeout,
		TLSHandshakeTimeout:   DefaultTLSHandshakeTimeout,
		ResponseHeaderTimeout: DefaultResponseTimeout,
		ExpectContinueTimeout: DefaultExpectContinue,
	}

	retryClient.RetryMax = DefaultRetryMax
	retryClient.HTTPClient.Timeout = DefaultHTTPTimeout
	retryClient.Backoff = retryablehttp.LinearJitterBackoff

	retryClient.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		if err != nil {
			slog.Warn("http retry - network error", "error", err)
			return true, err
		}
		if resp != nil && resp.StatusCode >= 500 {
			slog.Warn(
				"http retry - 5xx status",
				"url", resp.Request.URL.String(),
				"method", resp.Request.Method,
				"status", resp.Status,
				"status_code", resp.StatusCode,
			)
			return true, nil
		}
		return false, nil
	}

	return retryClient.StandardClient()
}

func configureLogger(level slog.Level, appName, instanceID string) (*slog.Logger, []io.Writer) {
	destinations := []io.Writer{os.Stdout}
	handler := slog.NewJSONHandler(io.MultiWriter(destinations...), &slog.HandlerOptions{
		AddSource: true,
		Level:     level,
	})
	logger := slog.New(handler).With("app", appName, "instance", instanceID)
	return logger, destinations
}

func setupTelemetry(ctx context.Context, opts *appOptions, instanceID string) (*sdktrace.TracerProvider, *metric.MeterProvider, *sdklog.LoggerProvider, *resource.Resource, []func() error, error) {
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpointURL(opts.AgentURL),
	)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("trace exporter: %w", err)
	}

	res, err := resource.New(ctx,
		resource.WithProcessRuntimeDescription(),
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
		resource.WithContainer(),
		resource.WithContainerID(),
		resource.WithAttributes(semconv.ServiceNameKey.String(opts.AppName)),
		resource.WithAttributes(attribute.KeyValue{
			Key:   attribute.Key("instance_id"),
			Value: attribute.StringValue(instanceID),
		}),
	)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("resource: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	exp, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpointURL(opts.AgentURL),
	)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("metric exporter: %w", err)
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(exp)),
	)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(tp)
	otel.SetMeterProvider(meterProvider)

	agentURL, err := url.Parse(opts.AgentURL)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("log exporter url: %w", err)
	}

	logExporter, err := otlploghttp.New(ctx,
		otlploghttp.WithInsecure(),
		otlploghttp.WithEndpoint(agentURL.Host),
	)
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("log exporter: %w", err)
	}

	loggerProvider := sdklog.NewLoggerProvider(
		sdklog.WithResource(res),
		sdklog.WithProcessor(sdklog.NewBatchProcessor(logExporter)),
	)

	cleanup := []func() error{
		func() error {
			slog.Info("Shutting down tracer provider")
			return tp.Shutdown(ctx)
		},
		func() error {
			slog.Info("Shutting down meter provider")
			return meterProvider.Shutdown(ctx)
		},
		func() error {
			slog.Info("Shutting down log provider")
			return loggerProvider.Shutdown(ctx)
		},
	}

	return tp, meterProvider, loggerProvider, res, cleanup, nil
}

func tracingTransport(base http.RoundTripper, serviceName string) http.RoundTripper {
	if base == nil {
		base = http.DefaultTransport
	}

	return otelhttp.NewTransport(
		base,
		otelhttp.WithTracerProvider(otel.GetTracerProvider()),
		otelhttp.WithPropagators(otel.GetTextMapPropagator()),
		otelhttp.WithSpanOptions(oteltrace.WithAttributes(
			attribute.String("peer.service", serviceName),
		)),
	)
}
