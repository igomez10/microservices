package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"socialapp/internal/authorizationparser"
	"socialapp/internal/middlewares/authorization"
	"socialapp/internal/middlewares/gandalf"
	"socialapp/pkg/controller/authentication"
	"socialapp/pkg/controller/comment"
	"socialapp/pkg/controller/role"
	"socialapp/pkg/controller/scope"
	"socialapp/pkg/controller/user"
	"socialapp/pkg/db"
	"socialapp/socialappapi/openapi"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
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

type Middleware func(http.Handler) http.Handler

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

	// Role service
	RoleAPIService := &role.RoleApiService{DB: queries, DBConn: dbConn}
	RoleAPIController := openapi.NewRoleApiController(RoleAPIService)

	// Scope service
	ScopeAPIService := &scope.ScopeApiService{DB: queries, DBConn: dbConn}
	ScopeAPIController := openapi.NewScopeApiController(ScopeAPIService)

	routers := []openapi.Router{
		CommentApiController,
		UserApiController,
		AuthApiController,
		RoleAPIController,
		ScopeAPIController,
	}

	authenticationMiddleware := gandalf.Middleware{DB: queries, DBConn: dbConn}
	middlewares := []Middleware{
		middleware.Logger,
		middleware.Heartbeat("/health"),
		authenticationMiddleware.Authenticate,
		cors.AllowAll().Handler,
		middleware.RequestID,
		middleware.RealIP,
		middleware.Recoverer,
		middleware.Timeout(60 * time.Second),
	}

	// open file
	openAPIPath := "openapi.yaml"
	openapiFile, err := os.Open(openAPIPath)
	if err != nil {
		log.Fatal().Err(err).Str("path", openAPIPath).Msg("failed to open openapi file")
	}

	// read file
	content, err := ioutil.ReadAll(openapiFile)
	if err != nil {
		log.Fatal().Err(err).Str("path", openAPIPath).Msg("failed to read openapi file")
	}

	// parse file
	doc, err := openapi3.NewLoader().LoadFromData(content)
	if err != nil {
		log.Fatal().Err(err)
	}

	authorizationParse := authorizationparser.FromOpenAPIToEndpointScopes(doc)
	router := NewRouter(middlewares, routers, authorizationParse)

	// Expose the registered metrics via HTTP.
	router.Handle("/metrics", promhttp.Handler())

	// Expose the api spec via HTTP.
	router.HandleFunc("/apispec", func(w http.ResponseWriter, r *http.Request) {
		log.Info().Msg("Health check")
		// send open api file
		// open api file
		file := "openapi.yaml"
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

func NewRouter(middlewares []Middleware, routers []openapi.Router, authorizationParse authorizationparser.EndpointAuthorizations) chi.Router {
	router := chi.NewRouter()

	for i := range middlewares {
		router.Use(middlewares[i])
	}

	// Custom misc middleware
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			customW := NewCustomResponseWriter(w)
			defer func() {
				// log failed auth requests
				logEvent := log.Warn().
					Str("Path", r.URL.Path).
					Str("X-Request-ID", r.Header.Get("X-Request-ID")).
					Str("ngrok", r.Header.Get("ngrok")).
					Str("Method", r.Method).
					Str("Path", r.URL.Path).
					Str("RemoteAddr", r.RemoteAddr).
					Str("UserAgent", r.UserAgent()).
					Str("Referer", r.Referer())

				if customW.statusCode == http.StatusUnauthorized {
					logEvent.Str("AuthHeader", r.Header.Get("Authorization")).
						Msgf("Failed request: %d", customW.statusCode)
				}

				logEvent.Str("StatusCode", fmt.Sprintf("%d", customW.statusCode)).
					Msgf("Finished Request")
			}()

			requestID := uuid.NewString()
			customW.Header().Set("X-Request-ID", requestID)
			r.Header.Set("X-Request-ID", requestID)
			r = r.WithContext(context.WithValue(r.Context(), "X-Request-ID", requestID))
			next.ServeHTTP(customW, r)
		})
	})

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

				// authorization
				requiredScopesForEndpoint := authorizationParse[route.Pattern][route.Method]
				mapRequiredScopes := map[string]bool{}
				for _, scope := range requiredScopesForEndpoint {
					mapRequiredScopes[scope] = true
				}
				authorizationRuler := authorization.Middleware{
					RequiredScopes: mapRequiredScopes,
				}

				r.Use(authorizationRuler.Authorize)
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
