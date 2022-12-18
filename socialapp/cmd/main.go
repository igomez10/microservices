package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	_ "net/http/pprof"
	"net/url"
	"os"
	"socialapp/internal/authorizationparser"
	"socialapp/internal/middlewares/beacon"
	"socialapp/internal/middlewares/cache"
	"socialapp/internal/middlewares/gandalf"
	"socialapp/internal/middlewares/requestid"
	"socialapp/internal/routers/proxyrouter"
	"socialapp/internal/routers/socialapprouter"
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
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	appPort *int = flag.Int("port", 8080, "main port for application")
)

func main() {
	flag.Parse()

	// Setup logger
	conn, err := net.Dial("udp", os.Getenv("LOGSTASH_HOST"))
	if err != nil {
		log.Warn().Err(err).Msg("Error connecting to logstash, default to stdout")
	}

	if conn == nil {
		log.Logger = zerolog.New(os.Stdout)
	} else {
		fmt.Printf("Writing logs to logstash: %q \n", conn.RemoteAddr())
		log.Logger = zerolog.New(conn)
	}

	log.Logger = log.With().
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

	// try to ping database with 5 seconds timeout
	pingCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := dbConn.PingContext(pingCtx); err != nil {
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

	socialappAllowlistedPaths := map[string]map[string]bool{
		"/users": {
			"POST": true,
		},
		"/metrics": {
			"GET": true,
		},
		"/apispec": {
			"GET": true,
		},
	}
	socialappAuthenticationMiddleware := gandalf.Middleware{
		DB:               queries,
		DBConn:           dbConn,
		Cache:            cache,
		AllowlistedPaths: socialappAllowlistedPaths,
		AllowBasicAuth:   false,
	}

	beacon := beacon.Beacon{Logger: log.Logger}

	// open apispec file
	openAPIPath := "openapi.yaml"
	openapiFile, err := os.Open(openAPIPath)
	if err != nil {
		log.Fatal().Err(err).Str("path", openAPIPath).Msg("failed to open openapi file")
	}

	// read api spec
	content, err := ioutil.ReadAll(openapiFile)
	if err != nil {
		log.Fatal().Err(err).Str("path", openAPIPath).Msg("failed to read openapi file")
	}
	openapiFile.Close()

	// parse api spec
	doc, err := openapi3.NewLoader().LoadFromData(content)
	if err != nil {
		log.Fatal().Err(err)
	}

	targetURL, err := url.Parse(os.Getenv("KIBANA_URL"))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse target url")
	}

	// 1. Kibana router (proxy)
	// kibanaAuthMiddleware := gandalf.Middleware{
	// 	DB:               queries,
	// 	DBConn:           dbConn,
	// 	Cache:            cache,
	// 	AllowlistedPaths: map[string]map[string]bool{},
	// 	AllowBasicAuth:   true,
	// }
	// authorizationRuler := authorization.Middleware{
	// 	RequiredScopes: map[string]bool{"kibana:read": true},
	// }
	kibanaRouterMiddlewares := []func(http.Handler) http.Handler{
		cors.AllowAll().Handler,
		middleware.Heartbeat("/health"),
		requestid.Middleware,
		beacon.Middleware,
		middleware.Recoverer,
		middleware.Timeout(60 * time.Second),
		// kibanaAuthMiddleware.Authenticate,
		// authorizationRuler.Authorize,
		middleware.RealIP,
	}
	kibanaRouter := proxyrouter.NewProxyRouter(os.Getenv("KIBANA_SUBDOMAIN"), targetURL, kibanaRouterMiddlewares)

	// 2. SocialApp router
	authorizationParse := authorizationparser.FromOpenAPIToEndpointScopes(doc)
	socialappMiddlewares := []func(http.Handler) http.Handler{
		cors.AllowAll().Handler,
		middleware.Heartbeat("/health"),
		requestid.Middleware,
		beacon.Middleware,
		middleware.Recoverer,
		middleware.Timeout(60 * time.Second),
		socialappAuthenticationMiddleware.Authenticate,
		middleware.RealIP,
		cache.Middleware,
	}
	socialappRouter := socialapprouter.NewSocialAppRouter(socialappMiddlewares, routers, authorizationParse)

	// 3. Main router for routing to different routers based on subdomain
	mainRouter := chi.NewRouter()
	mainRouter.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		switch r.Host {
		case os.Getenv("KIBANA_SUBDOMAIN"):
			kibanaRouter.Router.ServeHTTP(w, r)
		default:
			socialappRouter.Router.ServeHTTP(w, r)
		}
	})

	log.Info().Msgf("Listening on port %d", *appPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *appPort), mainRouter); err != nil {
		log.Fatal().Err(err).Msgf("Shutting down")
	}

}
