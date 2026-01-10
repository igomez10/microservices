package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/igomez10/microservices/urlshortener/generated/server"
	"github.com/igomez10/microservices/urlshortener/pkg/contexthelper"
	"github.com/igomez10/microservices/urlshortener/pkg/controllers/url"
	"github.com/igomez10/microservices/urlshortener/pkg/db"
	"github.com/igomez10/microservices/urlshortener/pkg/tracerhelper"
	flags "github.com/jessevdk/go-flags"
	_ "github.com/newrelic/go-agent/v3/integrations/nrpq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/slok/go-http-metrics/metrics/prometheus"
	metricsMiddleware "github.com/slok/go-http-metrics/middleware"
	"github.com/slok/go-http-metrics/middleware/std"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

var opts struct {
	Port       int    `long:"port" env:"PORT" default:"8080" description:"HTTP port"`
	HTTPAddr   string `long:"http-addr" env:"HTTP_ADDR" default:"" description:"HTTP address"`
	DBURL      string `long:"db-url" env:"DB_URL" default:"postgres://postgres:password@localhost:5432/urlshortener?sslmode=disable" description:"Database URL"`
	logLevel   string `long:"log-level" env:"LOG_LEVEL" default:"info" description:"Log level"`
	MetaServer struct {
		Addr string `long:"addr" env:"ADDR" default:"localhost" description:"Meta service address"`
		Port int    `long:"port" env:"PORT" default:"8081" description:"Meta service port"`
	} `group:"Meta service" namespace:"meta" env-namespace:"META"`
	AgentURL string `long:"agent-url" env:"OTEL_EXPORTER_OTLP_ENDPOINT" default:"" description:"Agent URL" required:"true"`
	AppName  string `long:"app-name" env:"APP_NAME" default:"urlshortener" description:"Application name"`
}

func main() {
	mainCtx := context.Background()
	// Parse flags
	if _, err := flags.Parse(&opts); err != nil {
		if err.(*flags.Error).Type != flags.ErrHelp {
			slog.Error("failed to parse flags", "err", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	instanceID := uuid.NewString()
	level, err := parseLogLevel(opts.logLevel)
	if err != nil {
		slog.Error("failed to parse log level", "level", opts.logLevel, "err", err)
		os.Exit(1)
	}
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	})).With(
		"app", opts.AppName,
		"instance", instanceID,
	)
	slog.SetDefault(logger)

	// Connect to database
	dbConn, err := sql.Open("nrpostgres", opts.DBURL)
	if err != nil {
		slog.Error("failed to open database", "err", err)
		os.Exit(1)
	}

	defer func() {
		if err := dbConn.Close(); err != nil {
			slog.Error("failed to close database connection", "err", err)
		}
	}()

	if dbConn == nil {
		slog.Error("db is nil")
		os.Exit(1)
	}

	pingCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := dbConn.PingContext(pingCtx); err != nil {
		slog.Error("failed to ping database, shutting down", "err", err)
		os.Exit(1)
	}

	// Create queries instance for db calls
	queries := db.New()

	// Create URL API service and controller
	URLAPIService := &url.URLApiService{
		DB:     queries,
		DBConn: dbConn,
	}
	URLAPIController := server.NewURLAPIController(URLAPIService)

	// health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Start meta service
	go func() {
		meta := NewMetaRouter()
		slog.Info("starting meta service", "addr", fmt.Sprintf("%s:%d", opts.MetaServer.Addr, opts.MetaServer.Port))
		if err := http.ListenAndServe(fmt.Sprintf("%s:%d", opts.MetaServer.Addr, opts.MetaServer.Port), meta); err != nil {
			slog.Error("failed to start meta service", "err", err)
			os.Exit(1)
		}
	}()

	exporter, err := otlptracegrpc.New(mainCtx, otlptracegrpc.WithInsecure(), otlptracegrpc.WithEndpointURL(opts.AgentURL))
	if err != nil {
		slog.Error("failed to create otlp exporter for tracing", "endpoint", opts.AgentURL, "err", err)
		os.Exit(1)
	}

	// Create a new tracer provider with a batch span processor and the otlp exporter.
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("urlshortener"),
			attribute.KeyValue{Key: attribute.Key("instance_id"), Value: attribute.StringValue(instanceID)},
		// Add more attributes as needed
		)),
	)

	// Register the tracer provider as the global provider.
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Start HTTP server
	middlewares := []func(http.Handler) http.Handler{
		middleware.RequestID,
		middleware.RealIP,
		middleware.Recoverer,
		ObservabilityMiddleware(),
	}
	urlRouter := NewRouter(middlewares, []server.Router{URLAPIController})

	addr := fmt.Sprintf("%s:%d", opts.HTTPAddr, opts.Port)
	slog.Info("starting HTTP server", "addr", addr)
	if err := http.ListenAndServe(addr, urlRouter); err != nil {
		slog.Error("failed to start HTTP server", "err", err)
		os.Exit(1)
	}
}

type Pattern struct {
	Pattern string
}

type customHTTPResponseWriter struct {
	http.ResponseWriter
	statusCode    int
	headerWritten bool
}

func (w *customHTTPResponseWriter) WriteHeader(statusCode int) {
	if w.headerWritten {
		return
	}
	w.statusCode = statusCode
	w.headerWritten = true
	w.ResponseWriter.WriteHeader(statusCode)
}

func ObservabilityMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, span := tracerhelper.GetTracer().Start(r.Context(), "ObservabilityMiddleware")
			defer span.End()

			span.SetAttributes(attribute.String("x-request-id", middleware.GetReqID(r.Context())))

			logger := contexthelper.GetLoggerInContext(ctx)
			logger = logger.With(
				"trace_id", span.SpanContext().TraceID().String(),
				"x-request-id", middleware.GetReqID(r.Context()),
			)
			ctx = contexthelper.SetLoggerInContext(ctx, logger)
			r = r.WithContext(ctx)

			customW := &customHTTPResponseWriter{
				ResponseWriter: w,
			}
			next.ServeHTTP(customW, r)
			logger.Info("finished request",
				"method", r.Method,
				"status_code", customW.statusCode,
				"path", r.URL.Path,
				"remote_addr", r.RemoteAddr,
				"user_agent", r.UserAgent(),
			)
		})
	}
}

func NewMetaRouter() chi.Router {
	mainRouter := chi.NewRouter()

	mainRouter.Mount("/debug", middleware.Profiler())
	// Expose health the registered metrics via HTTP, no logging for those requests
	mainRouter.Group(func(r chi.Router) {
		// HEALTH
		r.MethodFunc("GET", "/health", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("OK"))
		})

		// METRICS
		r.Handle("/metrics", promhttp.Handler())

		// OPENAPI
		// Expose the api spec via HTTP.
		r.HandleFunc("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
			// send open api file
			// open api file
			file := "openapi.yaml"
			content, err := os.ReadFile(file)
			if err != nil {
				slog.Error("error reading file", "err", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Write(content)
		})
	})

	return mainRouter
}

// Adds the request pattern into the context, this is only required because chi.mux
// does not provide a way to get the pattern from the request. This middleware will update the
// string pointer saved in the context as "pattern"
func (p *Pattern) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracerhelper.GetTracer().Start(r.Context(), "PatternMiddleware")
		defer span.End()
		span.SetAttributes(attribute.String("x-pattern", p.Pattern))

		logger := contexthelper.GetLoggerInContext(ctx)
		logger = logger.With("x-pattern", p.Pattern)
		ctx = contexthelper.SetLoggerInContext(ctx, logger)
		ctx = contexthelper.SetRequestPatternInContext(ctx, p.Pattern)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// NewRouter creates a new router for any number of api routers
func NewRouter(middlewares []func(http.Handler) http.Handler, routers []server.Router) chi.Router {
	mainRouter := chi.NewRouter()

	// Main router group for api logic
	mainRouter.Group(func(r chi.Router) {
		mdlw := metricsMiddleware.New(metricsMiddleware.Config{
			Recorder: prometheus.NewRecorder(prometheus.Config{}),
			Service:  opts.AppName,
		})

		for _, api := range routers {
			for _, route := range api.Routes() {

				r.Group(func(r chi.Router) {
					// Add open telemetry traces before any middleware can short-circuit.
					resourceName := fmt.Sprintf("%s_%s", route.Method, route.Pattern)
					r.Use(otelhttp.NewMiddleware(
						resourceName,
						otelhttp.WithTracerProvider(otel.GetTracerProvider()),
						otelhttp.WithPropagators(otel.GetTextMapPropagator()),
					))

					// use a  custom middleware to record the metrics on the route pattern.
					r.Use(std.HandlerProvider(route.Pattern, mdlw))

					pattern := Pattern{Pattern: route.Pattern}
					r.Use(pattern.Middleware)

					for i := range middlewares {
						r.Use(middlewares[i])
					}

					handler := route.HandlerFunc
					r.Method(route.Method, route.Pattern, handler)
				})
			}
		}
	})

	return mainRouter
}

func parseLogLevel(input string) (slog.Level, error) {
	switch strings.ToLower(input) {
	case "debug":
		return slog.LevelDebug, nil
	case "info", "":
		return slog.LevelInfo, nil
	case "warn", "warning":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return slog.LevelInfo, fmt.Errorf("invalid log level: %s", input)
	}
}
