package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	_ "net/http/pprof"
	"net/url"
	"os"
	"os/signal"
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
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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
	"go.opentelemetry.io/otel/sdk/trace"
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

func main() {
	// Create cancellable context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config, shutdown := getConfig(ctx)
	defer func() {
		log.Info().Msg("Running shutdown functions")
		for _, fn := range shutdown {
			if err := fn(); err != nil {
				log.Error().Err(err).Msg("failed to shutdown")
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
		log.Info().Str("signal", sig.String()).Msg("Received shutdown signal")
		cancel() // Cancel context to trigger graceful shutdown
	case err := <-serverErrors:
		if err != nil {
			log.Error().Err(err).Msg("Server error")
		}
	}

	log.Info().Msg("Main shutdown complete")
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
		log.Info().Msg("Context cancelled, starting graceful shutdown")

		// Give server time to finish existing requests
		shutdownCtx, cancel := context.WithTimeout(context.Background(), DefaultShutdownTimeout)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Error().Err(err).Msg("Server shutdown error")
		}
	}()

	log.Info().Msgf("Server listening on port %d", config.AppPort)
	err = httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server error: %w", err)
	}

	log.Info().Msg("Server stopped")
	return nil
}

// CreateDBPools creates a pool of connections to the database
func CreateDBPools(ctx context.Context, databaseURL string, numPools int, applicationName string) *ForcedConnectionPool {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse database url")
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
			log.Fatal().Err(err)
		}

		if dbConn == nil {
			log.Fatal().Msg("db is nil")
		}

		pingCtx, cancel := context.WithTimeout(context.Background(), DefaultDBPingTimeout)
		defer cancel()
		if err := dbConn.Ping(pingCtx); err != nil {
			log.Fatal().Err(err).Msg("failed to ping database, shutting down")
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
	var opts struct {
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

	shutdown := []func() error{}
	_, err := flags.Parse(&opts)
	if err != nil {
		panic(err)
	}

	instanceID := uuid.NewString()
	log.Info().
		Str("logHost", opts.LogHost).
		Str("logLevel", opts.LogLevel).
		Str("appName", opts.AppName).
		Str("propertiesSubdomain", opts.PropertiesSubdomain).
		Str("urlServiceHost", opts.UrlShortenerSubdomain).
		Str("agentURL", opts.AgentURL).
		Str("urlShortenerSubdomain", opts.UrlShortenerSubdomain).
		Str("urlShortenerURL", opts.URLShortenerURL).
		Str("socialappSubdomain", opts.SocialappSubdomain).
		Str("instanceID", instanceID).
		Str("puttyknifeDomain", opts.PuttyknifeDomain).
		Str("puttyknifeURL", opts.PuttyknifeURL).
		Str("redisURL", opts.RedisURL).
		Str("kibanaSubdomain", opts.KibanaSubdomain).
		Msg("Starting SocialApp")

	// Setup retryable http client
	retryClient := retryablehttp.NewClient()
	retryClient.Logger = nil
	retryClient.RequestLogHook = func(_ retryablehttp.Logger, req *http.Request, attempt int) {
		if attempt >= 1 {
			log.Warn().
				Str("method", req.Method).
				Str("url", req.URL.String()).
				Int("attempt", attempt).
				Msgf("http retry")
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
			log.Warn().Err(err).Msg("http retry - network error")
			return true, err
		}
		if resp != nil && resp.StatusCode >= 500 {
			log.Warn().
				Stringer("url", resp.Request.URL).
				Str("method", resp.Request.Method).
				Str("status", resp.Status).
				Int("status_code", resp.StatusCode).
				Msg("http retry - 5xx status")
			return true, nil
		}
		return false, nil
	}
	http.DefaultClient = retryClient.StandardClient()

	// Set proxy
	if opts.ProxyHost != "" {
		if u, err := url.Parse(opts.ProxyHost); err != nil {
			log.Err(err).Msgf("Failed to parse proxy URL")
		} else {
			http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(u)}
		}
	}

	// Parse log level
	parsedLogLevel, err := zerolog.ParseLevel(opts.LogLevel)
	if err != nil {
		log.Fatal().Err(err).Msgf("Invalid log level, %s", opts.LogLevel)
	}

	// Setup log destinations
	var logDestinations []io.Writer = []io.Writer{os.Stdout}

	// Parse properties subdomain
	var propertiesSubdomainURL *url.URL
	if len(opts.PropertiesSubdomain) != 0 {
		u, err := url.Parse(opts.PropertiesSubdomain)
		if err != nil {
			log.Fatal().Err(err).Msgf("failed to parse properties subdomain url %s", opts.PropertiesSubdomain)
		}
		propertiesSubdomainURL = u
	}

	// Parse url service host
	urlService, err := url.Parse(opts.URLShortenerURL)
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to parse url service host url %s", opts.URLShortenerURL)
	}

	var puttyknifeURL *url.URL
	if len(opts.PuttyknifeURL) != 0 {
		u, err := url.Parse(opts.PuttyknifeURL)
		if err != nil {
			log.Fatal().Err(err).Msgf("failed to parse puttyknife url %s", opts.PuttyknifeURL)
		}
		puttyknifeURL = u
	}

	queries := dbpgx.New()

	// Setup tracing
	http.DefaultClient = &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpointURL(opts.AgentURL),
	)
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to create otlp exporter for tracing %q", opts.AgentURL)
	}

	res, err := resource.New(ctx,
		resource.WithProcessRuntimeDescription(),
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
		resource.WithContainer(),
		resource.WithContainerID(),
		resource.WithAttributes(attribute.KeyValue{
			Key:   attribute.Key("instance_id"),
			Value: attribute.StringValue(instanceID),
		}),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create resource")
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)

	exp, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpointURL(opts.AgentURL),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create otlp exporter for metrics")
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(exp)),
	)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(tp)
	otel.SetMeterProvider(meterProvider)

	shutdown = append(shutdown, func() error {
		log.Info().Msg("Shutting down tracer provider")
		return tp.Shutdown(ctx)
	})

	shutdown = append(shutdown, func() error {
		log.Info().Msg("Shutting down meter provider")
		return meterProvider.Shutdown(ctx)
	})

	// Setup OTLP logs exporter
	agentURL, err := url.Parse(opts.AgentURL)
	if err != nil {
		log.Fatal().Err(err).Str("agentURL", opts.AgentURL).Msg("failed to parse agent URL for log exporter")
	}
	logExporter, err := otlploghttp.New(ctx,
		otlploghttp.WithInsecure(),
		otlploghttp.WithEndpoint(agentURL.Host),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create otlp log exporter")
	}

	loggerProvider := sdklog.NewLoggerProvider(
		sdklog.WithResource(res),
		sdklog.WithProcessor(sdklog.NewBatchProcessor(logExporter)),
	)

	shutdown = append(shutdown, func() error {
		log.Info().Msg("Shutting down log provider")
		return loggerProvider.Shutdown(ctx)
	})

	// Connect to database
	connections := CreateDBPools(ctx, opts.DatabaseURL, DefaultNumDBPools, fmt.Sprintf("%s-%s", opts.AppName, instanceID))

	shutdown = append(shutdown, func() error {
		log.Info().Msg("Shutting down database connections")
		connections.Close()
		return nil
	})

	redisOpts, err := redis.ParseURL(opts.RedisURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse redis url")
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
		log.Fatal().Err(err).Str("path", openAPIPath).Msg("failed to read openapi file")
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
