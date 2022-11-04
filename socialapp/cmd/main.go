package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"socialapp/internal/authorizationparser"
	"socialapp/internal/middlewares/authorization"
	"socialapp/internal/middlewares/cache"
	"socialapp/internal/middlewares/failedrequests"
	"socialapp/internal/middlewares/gandalf"
	"socialapp/internal/middlewares/requestid"
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
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/slok/go-http-metrics/metrics/prometheus"
	metricsMiddleware "github.com/slok/go-http-metrics/middleware"
	"github.com/slok/go-http-metrics/middleware/std"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	appPort *int = flag.Int("port", 8080, "main port for application")
)

type Middleware func(http.Handler) http.Handler

func main() {
	flag.Parse()

	// Setup logger
	log.Logger = zerolog.New(os.Stdout).With().
		Str("app", "socialapp").
		Str("instance", uuid.NewString()).
		Timestamp().
		Logger()

	log.Info().Msgf("Starting PORT: %d", *appPort)

	dbConn, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatal().Err(err)
	}
	defer dbConn.Close()

	if dbConn == nil {
		log.Fatal().Msg("db is nil")
	}

	if err := dbConn.Ping(); err != nil {
		log.Fatal().Err(err).Msg("failed to ping database, shutting down")
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

	redisOpts, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse redis url")
	}

	cache := cache.NewCache(cache.CacheConfig{
		RedisOpts: redisOpts,
	})

	authenticationMiddleware := gandalf.Middleware{DB: queries, DBConn: dbConn, Cache: cache}

	middlewares := []Middleware{
		cors.AllowAll().Handler,
		middleware.Heartbeat("/health"),
		requestid.RequestIDMiddleware,
		failedrequests.FailedRequestsMiddleware,
		middleware.Recoverer,
		middleware.RequestID,
		middleware.Timeout(60 * time.Second),
		authenticationMiddleware.Authenticate,
		middleware.RealIP,
		cache.Middleware,
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
	log.Info().Msgf("Listening on port %d", *appPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *appPort), router); err != nil {
		log.Fatal().Err(err).Msgf("Shutting down")
	}

}

func NewRouter(middlewares []Middleware, routers []openapi.Router, authorizationParse authorizationparser.EndpointAuthorizations) chi.Router {
	mainRouter := chi.NewRouter()

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

	// Main router group, here is the main logic
	mainRouter.Group(func(r chi.Router) {
		for i := range middlewares {
			r.Use(middlewares[i])
		}

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
	})

	return mainRouter
}
