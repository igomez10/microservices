package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/igomez10/microservices/urlshortener/generated/server"
	"github.com/igomez10/microservices/urlshortener/pkg/contexthelper"
	"github.com/igomez10/microservices/urlshortener/pkg/controllers/url"
	"github.com/igomez10/microservices/urlshortener/pkg/db"
	flags "github.com/jessevdk/go-flags"
	_ "github.com/newrelic/go-agent/v3/integrations/nrpq"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/slok/go-http-metrics/metrics/prometheus"
	metricsMiddleware "github.com/slok/go-http-metrics/middleware"
	"github.com/slok/go-http-metrics/middleware/std"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

var opts struct {
	Port       int    `long:"port" env:"PORT" default:"8080" description:"HTTP port"`
	HTTPAddr   string `long:"http-addr" env:"HTTP_ADDR" defatult:"" description:"HTTP address"`
	DBURL      string `long:"db-url" env:"DB_URL" default:"postgres://postgres:password@localhost:5432/urlshortener?sslmode=disable" description:"Database URL"`
	logLevel   string `long:"log-level" env:"LOG_LEVEL" default:"info" description:"Log level"`
	MetaServer struct {
		Addr string `long:"addr" env:"ADDR" default:"localhost" description:"Meta service address"`
		Port int    `long:"port" env:"PORT" default:"8081" description:"Meta service port"`
	} `group:"Meta service" namespace:"meta" env-namespace:"META"`
	NewRelicLicense string `long:"newrelic-license" env:"NEWRELIC_LICENSE" default:"" description:"New relic license"`
	AgentURL        string `long:"agent-url" env:"AGENT_URL" default:"" description:"Agent URL"`
	AppName         string `long:"app-name" env:"APP_NAME" default:"urlshortener" description:"Application name"`
}

func main() {
	mainCtx := context.Background()
	// Parse flags
	if _, err := flags.Parse(&opts); err != nil {
		if err.(*flags.Error).Type != flags.ErrHelp {
			log.Fatal().Err(err).Msg("failed to parse flags")
		}
		os.Exit(0)
	}

	instanceID := uuid.NewString()
	log.Logger = log.With().
		Str("app", opts.AppName).
		Str("instance", instanceID).
		Timestamp().
		Caller().
		Logger()

	// log config print stack trace
	log.Logger = log.With().Caller().Logger()
	// Set log level zerolog
	if _, err := zerolog.ParseLevel(opts.logLevel); err != nil {
		log.Fatal().Err(err).Msg("failed to parse log level")
	}

	// Connect to database
	dbConn, err := sql.Open("nrpostgres", opts.DBURL)
	if err != nil {
		log.Fatal().Err(err)
	}

	defer func() {
		if err := dbConn.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close database connection")
		}
	}()

	if dbConn == nil {
		log.Fatal().Msg("db is nil")
	}

	pingCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := dbConn.PingContext(pingCtx); err != nil {
		log.Fatal().Err(err).Msg("failed to ping database, shutting down")
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
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Start meta service
	go func() {
		meta := NewMetaRouter()
		log.Info().Str("addr", fmt.Sprintf("%s:%d", opts.MetaServer.Addr, opts.MetaServer.Port)).Msg("starting meta service")
		if err := http.ListenAndServe(fmt.Sprintf("%s:%d", opts.MetaServer.Addr, opts.MetaServer.Port), meta); err != nil {
			log.Fatal().Err(err).Msg("failed to start meta service")
		}
	}()

	// start new relic
	var newrelicApp *newrelic.Application
	if opts.NewRelicLicense != "" {
		na, err := newrelic.NewApplication(
			newrelic.ConfigAppName("urlshortener"),
			newrelic.ConfigLicense(opts.NewRelicLicense),
			newrelic.ConfigAppLogForwardingEnabled(false),
			newrelic.ConfigAppLogEnabled(false),
			newrelic.ConfigDistributedTracerEnabled(true),
		)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to create new relic application")
		}
		newrelicApp = na
	} else {
		log.Warn().Msg("new relic license not provided, new relic will not be enabled")
	}

	exporter, err := otlptracegrpc.New(mainCtx, otlptracegrpc.WithInsecure(), otlptracegrpc.WithEndpointURL(opts.AgentURL))
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to create otlp exporter for tracing %q", opts.AgentURL)
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

	// Start HTTP server
	middlewares := []func(http.Handler) http.Handler{
		middleware.RequestID,
		middleware.RealIP,
		middleware.Recoverer,
	}
	urlRouter := NewRouter(middlewares, []server.Router{URLAPIController}, newrelicApp)

	addr := fmt.Sprintf("%s:%d", opts.HTTPAddr, opts.Port)
	log.Info().Str("addr", addr).Msg("starting HTTP server")
	if err := http.ListenAndServe(addr, urlRouter); err != nil {
		log.Fatal().Err(err).Msg("failed to start HTTP server")
	}
}

type Pattern struct {
	Pattern string
}

func NewMetaRouter() chi.Router {
	mainRouter := chi.NewRouter()

	mainRouter.Mount("/debug", middleware.Profiler())
	// Expose health the registered metrics via HTTP, no logging for those requests
	mainRouter.Group(func(r chi.Router) {
		// HEALTH
		r.MethodFunc("GET", "/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})

		// METRICS
		r.Handle("/metrics", promhttp.Handler())

		// OPENAPI
		// Expose the api spec via HTTP.
		r.HandleFunc("/apispec", func(w http.ResponseWriter, r *http.Request) {
			// send open api file
			// open api file
			file := "openapi.yaml"
			content, err := os.ReadFile(file)
			if err != nil {
				log.Error().Err(err).Msg("Error reading file")
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
		ctx := contexthelper.SetRequestPatternInContext(r.Context(), p.Pattern)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// NewRouter creates a new router for any number of api routers
func NewRouter(middlewares []func(http.Handler) http.Handler, routers []server.Router, newrelicApp *newrelic.Application) chi.Router {
	mainRouter := chi.NewRouter()

	// Main router group for api logic
	mainRouter.Group(func(r chi.Router) {

		mdlw := metricsMiddleware.New(metricsMiddleware.Config{
			Recorder: prometheus.NewRecorder(prometheus.Config{}),
		})

		for _, api := range routers {
			for _, route := range api.Routes() {
				var handler http.Handler
				handler = route.HandlerFunc

				r.Group(func(r chi.Router) {
					// use a  custom middleware to record the metrics on the route pattern.
					r.Use(std.HandlerProvider(route.Pattern, mdlw))

					pattern := Pattern{Pattern: route.Pattern}
					r.Use(pattern.Middleware)

					for i := range middlewares {
						r.Use(middlewares[i])
					}

					if newrelicApp != nil {
						_, handler = newrelic.WrapHandle(newrelicApp, route.Pattern, handler)
					}

					resourceName := fmt.Sprintf("%s_%s", route.Method, route.Pattern)
					otelHandler := otelhttp.NewHandler(http.Handler(handler), resourceName)
					r.Method(route.Method, route.Pattern, otelHandler)
				})
			}
		}
	})

	return mainRouter
}
