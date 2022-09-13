package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"socialapp/pkg/controller/authentication"
	"socialapp/pkg/controller/comment"
	"socialapp/pkg/controller/user"
	"socialapp/pkg/db"
	"socialapp/socialappapi/openapi"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/slok/go-http-metrics/metrics/prometheus"
	metricsMiddleware "github.com/slok/go-http-metrics/middleware"

	"github.com/rs/zerolog/log"
	"github.com/slok/go-http-metrics/middleware/std"
)

var (
	appPort *int = flag.Int("port", 8080, "main port for application")
)

func main() {
	flag.Parse()
	log.Info().Msgf("Starting PORT: %d", *appPort)

	dbConn, err := sql.Open("mysql", os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatal().Err(err)
	}
	defer dbConn.Close()

	if dbConn == nil {
		log.Fatal().Msg("db is nil")
	}
	defer dbConn.Close()

	queries := db.New()

	// Comment service
	CommentApiService := &comment.CommentService{DB: queries, DBConn: dbConn}
	CommentApiController := openapi.NewCommentApiController(CommentApiService)

	// User service
	UserApiService := &user.UserApiService{DB: queries, DBConn: dbConn}
	UserApiController := openapi.NewUserApiController(UserApiService)

	// Auth service
	AuthApiService := &authentication.AuthenticationService{DB: queries, DBConn: dbConn}
	AuthApiController := openapi.NewAuthenticationApiController(AuthApiService)

	router := NewRouter(CommentApiController, UserApiController, AuthApiController)

	// Expose the registered metrics via HTTP.
	router.Handle("/metrics", promhttp.Handler())

	// Expose the api spec via HTTP.
	router.HandleFunc("/apispec", func(w http.ResponseWriter, r *http.Request) {
		log.Info().Msg("H . ealth check")
		// send open api file
		// open api file
		file := "./openapi.yaml"
		content, err := os.ReadFile(file)
		if err != nil {
			log.Error().Err(err).Msg("Error reading file")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// w.Header().Set("Content-Type", "application/json")
		w.Write(content)
	})

	if err := http.ListenAndServe(fmt.Sprintf(":%d", *appPort), router); err != nil {
		log.Fatal().Err(err).Msgf("Shutting down")
	}

}

func NewRouter(routers ...openapi.Router) chi.Router {
	router := chi.NewRouter()

	// cors middleware
	router.Use(cors.AllowAll().Handler)

	// health check
	router.Use(middleware.Heartbeat("/health"))

	// Custom misc middleware
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			customW := NewCustomResponseWriter(w)
			defer func() {
				// log failed auth requests
				if customW.statusCode == http.StatusUnauthorized {
					log.Warn().
						Str("Path", r.URL.Path).
						Str("X-Request-ID", r.Header.Get("X-Request-ID")).
						Str("AuthHeader", r.Header.Get("Authorization")).
						Str("Secret", r.Header.Get("Secret")).
						Str("ngrok", r.Header.Get("ngrok")).
						Str("Method", r.Method).
						Str("Path", r.URL.Path).
						Msg("Unauthorized request")

				}
				log.Info().Msgf("%s %s %d", r.Method, r.RequestURI, customW.statusCode)
			}()

			requestID := uuid.NewString()
			customW.Header().Set("X-Request-ID", requestID)
			r.Header.Set("X-Request-ID", requestID)
			r = r.WithContext(context.WithValue(r.Context(), "X-Request-ID", requestID))
			next.ServeHTTP(customW, r)

		})
	})

	// auth middleware
	router.Use(middleware.BasicAuth("realm", map[string]string{"admin": "admin"}))
	// request id middleware
	router.Use(middleware.RequestID)
	// realip middleware
	router.Use(middleware.RealIP)
	// usual logger middleware
	router.Use(middleware.Logger)
	// recover middleware for recovering from panics
	router.Use(middleware.Recoverer)
	// timeout middleware
	router.Use(middleware.Timeout(60 * time.Second))

	mdlw := metricsMiddleware.New(metricsMiddleware.Config{
		Recorder: prometheus.NewRecorder(prometheus.Config{}),
	})

	for _, api := range routers {
		for _, route := range api.Routes() {
			var handler http.Handler
			handler = route.HandlerFunc

			router.Group(func(r chi.Router) {
				// use a middleware to record the metrics on the route pattern.
				r.Use(std.HandlerProvider(route.Pattern, mdlw))

				r.Method(route.Method, route.Pattern, handler)
			})
		}
	}

	return router
}

// custom response writer for capturing status code in the response
type customResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewCustomResponseWriter(w http.ResponseWriter) *customResponseWriter {
	return &customResponseWriter{w, http.StatusOK}
}

func (lrw *customResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
