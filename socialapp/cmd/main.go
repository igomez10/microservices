package main

import (
	//  With pprof to enable profiling
	// _ "net/http/pprof"

	"context"
	"database/sql"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/igomez10/microservices/socialapp/internal/authorizationparser"
	"github.com/igomez10/microservices/socialapp/internal/eventRecorder"
	"github.com/igomez10/microservices/socialapp/internal/middlewares/authorization"
	"github.com/igomez10/microservices/socialapp/internal/middlewares/beacon"
	"github.com/igomez10/microservices/socialapp/internal/middlewares/cache"
	"github.com/igomez10/microservices/socialapp/internal/middlewares/gandalf"
	"github.com/igomez10/microservices/socialapp/internal/middlewares/requestid"
	"github.com/igomez10/microservices/socialapp/internal/routers/proxyrouter"
	"github.com/igomez10/microservices/socialapp/internal/routers/socialapprouter"
	"github.com/igomez10/microservices/socialapp/pkg/controller/authentication"
	"github.com/igomez10/microservices/socialapp/pkg/controller/comment"
	"github.com/igomez10/microservices/socialapp/pkg/controller/role"
	"github.com/igomez10/microservices/socialapp/pkg/controller/scope"
	socialappurl "github.com/igomez10/microservices/socialapp/pkg/controller/url"
	"github.com/igomez10/microservices/socialapp/pkg/controller/user"
	"github.com/igomez10/microservices/socialapp/pkg/db"
	"github.com/igomez10/microservices/socialapp/socialappapi/openapi"
	_ "github.com/lib/pq"
	_ "github.com/newrelic/go-agent/v3/integrations/nrpq"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Configuration struct {
	appName             string
	appPort             int
	proxyURL            string
	logLevel            zerolog.Level
	logDestination      io.Writer
	dbConnections       *ForcedConnectionPool
	queries             *db.Queries
	cache               *cache.Cache
	propertiesSubdomain *url.URL
	newRelicApp         *newrelic.Application
}

func main() {
	appPort := flag.Int("port", 8080, "main port for application")
	proxyURL := flag.String("proxy", "", "proxy url, \"http://localhost:9091\"")
	logHost := flag.String("logHost", os.Getenv("LOGSTASH_HOST"), "log host url \"tcp://localhost:5000\"")
	logLevel := flag.String("logLevel", "info", "log level info/error/warning")
	appName := flag.String("appName", "socialapp", "name of the app for logs")
	propertiesSubdomain := flag.String("propertiesSubdomain", os.Getenv("PROPERTIES_SUBDOMAIN"), "Properties subdomain")
	newRelicLicense := flag.String("newRelicLicense", os.Getenv("NEW_RELIC_LICENSE"), "New relic license API Key")

	flag.Parse()

	// Set proxy
	if *proxyURL != "" {
		if u, err := url.Parse(*proxyURL); err != nil {
			http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(u)}
		} else {
			log.Err(err).Msgf("Failed to parse proxy URL")
		}
	}

	// parse log level
	parsedLogLevel, err := zerolog.ParseLevel(*logLevel)
	if err != nil {
		log.Fatal().Err(err).Msgf("Invalid log level, %s", *logLevel)
	}

	var logDestination io.Writer = os.Stdout
	// Validate logHost is a url
	if *logHost != "" && len(*logHost) != 0 {
		u, err := url.Parse(*logHost)
		if err != nil {
			log.Fatal().Err(err).Msgf("failed to parse log host url")
		}

		conn, err := net.Dial(u.Scheme, u.Host)
		if err != nil {
			log.Fatal().Err(err).Msgf("failed to establish connection with log host")
		}

		logDestination = conn
	}

	// Connect to database
	// force creation of 8 connections, one per service
	connections := CreateDBPools(os.Getenv("DATABASE_URL"), 1)

	queries := db.New()

	redisOpts, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse redis url")
	}

	redisOpts.PoolSize = 10
	redisOpts.MinIdleConns = 10

	cache := cache.NewCache(cache.CacheConfig{
		RedisOpts: redisOpts,
	})

	// parse properties subdomain
	var propertiesSubdomainURL *url.URL
	if len(*propertiesSubdomain) != 0 && *propertiesSubdomain != "" {
		u, err := url.Parse(*propertiesSubdomain)
		if err != nil {
			log.Fatal().Err(err).Msgf("failed to parse properties subdomain url %s", *propertiesSubdomain)
		}
		propertiesSubdomainURL = u
	}

	var newrelicApp *newrelic.Application
	if *newRelicLicense != "" {
		app, err := newrelic.NewApplication(
			newrelic.ConfigAppName("socialapp"),
			newrelic.ConfigLicense(*newRelicLicense),
			newrelic.ConfigAppLogForwardingEnabled(false),
			newrelic.ConfigAppLogEnabled(false),
			newrelic.ConfigDistributedTracerEnabled(false),
		)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to create new relic application")
		} else {
			newrelicApp = app
		}
	}

	c := Configuration{
		appPort:             *appPort,
		proxyURL:            *proxyURL,
		logLevel:            parsedLogLevel,
		logDestination:      logDestination,
		appName:             *appName,
		dbConnections:       connections,
		queries:             queries,
		cache:               cache,
		propertiesSubdomain: propertiesSubdomainURL,
		newRelicApp:         newrelicApp,
	}

	defer connections.Close()
	run(c)
}

func run(config Configuration) {
	// Setup logger
	zerolog.SetGlobalLevel(config.logLevel)
	log.Logger = zerolog.New(config.logDestination)

	instanceID := uuid.NewString()
	log.Logger = log.With().
		Str("app", config.appName).
		Str("instance", instanceID).
		Timestamp().
		Logger()

	// EventRecorder for event sourcing
	eventRecorder := eventRecorder.EventRecorder{
		DB: config.queries,
	}

	// Comment service
	CommentApiService := &comment.CommentService{
		DB:     config.queries,
		DBConn: config.dbConnections.GetPool(),
	}
	CommentApiController := openapi.NewCommentAPIController(CommentApiService)

	// User service
	UserApiService := &user.UserApiService{
		DB:            config.queries,
		DBConn:        config.dbConnections.GetPool(),
		EventRecorder: eventRecorder,
	}
	UserApiController := openapi.NewUserAPIController(UserApiService)

	// Auth service
	AuthApiService := &authentication.AuthenticationService{
		DB:     config.queries,
		DBConn: config.dbConnections.GetPool(),
	}
	AuthApiController := openapi.NewAuthenticationAPIController(AuthApiService)

	// Role service
	RoleAPIService := &role.RoleApiService{
		DB:     config.queries,
		DBConn: config.dbConnections.GetPool(),
	}
	RoleAPIController := openapi.NewRoleAPIController(RoleAPIService)

	// Scope service
	ScopeAPIService := &scope.ScopeApiService{
		DB:     config.queries,
		DBConn: config.dbConnections.GetPool(),
	}
	ScopeAPIController := openapi.NewScopeAPIController(ScopeAPIService)

	URLAPIService := &socialappurl.URLApiService{
		DB:     config.queries,
		DBConn: config.dbConnections.GetPool(),
	}
	URLAPIController := openapi.NewURLAPIController(URLAPIService)

	routers := []openapi.Router{
		CommentApiController,
		UserApiController,
		AuthApiController,
		RoleAPIController,
		ScopeAPIController,
		URLAPIController,
	}

	socialappAllowlistedPaths := map[string]map[string]bool{
		"/metrics": {
			"GET": true,
		},
		"/apispec": {
			"GET": true,
		},
	}
	socialappAuthenticationMiddleware := gandalf.Middleware{
		DB:               config.queries,
		DBConn:           config.dbConnections.GetPool(),
		Cache:            config.cache,
		AllowlistedPaths: socialappAllowlistedPaths,
		AllowBasicAuth:   false,
		AuthEndpoint:     "/v1/oauth/token",
	}

	beacon := beacon.Beacon{Logger: log.Logger}

	// open apispec file
	openAPIPath := "openapi.yaml"
	openapiFile, err := os.Open(openAPIPath)
	if err != nil {
		log.Fatal().Err(err).Str("path", openAPIPath).Msg("failed to open openapi file")
	}

	kibanaTargetURL, err := url.Parse(os.Getenv("KIBANA_URL"))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse target url")
	}

	// 1. Kibana router (proxy)
	kibanaAuthMiddleware := gandalf.Middleware{
		DB:               config.queries,
		DBConn:           config.dbConnections.GetPool(),
		Cache:            config.cache,
		AllowlistedPaths: map[string]map[string]bool{},
		AllowBasicAuth:   true,
		AuthEndpoint:     "/v1/oauth/token",
	}
	authorizationRuler := authorization.Middleware{
		RequiredScopes: map[string]bool{"kibana:read": true},
	}
	kibanaRouterMiddlewares := []func(http.Handler) http.Handler{
		cors.AllowAll().Handler,
		middleware.Heartbeat("/health"),
		requestid.Middleware,
		beacon.Middleware,
		middleware.Recoverer,
		middleware.Timeout(60 * time.Second),
		kibanaAuthMiddleware.Authenticate,
		authorizationRuler.Authorize,
		middleware.RealIP,
	}
	kibanaSubdomain := os.Getenv("KIBANA_SUBDOMAIN")
	authKibanaRouter := proxyrouter.NewProxyRouter(kibanaTargetURL, kibanaRouterMiddlewares)

	// 2. SocialApp router
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
	socialappSubdomain := os.Getenv("SOCIALAPP_SUBDOMAIN")
	authorizationParse := authorizationparser.FromOpenAPIToEndpointScopes(doc)

	// compress responses with gzip to save bandwidth
	compressor := middleware.NewCompressor(5, "application/json", "application/x-yaml", "gzip", "application/json; charset=UTF-8")

	socialappMiddlewares := []func(http.Handler) http.Handler{
		compressor.Handler,
		cors.AllowAll().Handler,
		middleware.Heartbeat("/health"),
		requestid.Middleware,
		beacon.Middleware,
		middleware.Recoverer,
		middleware.Timeout(60 * time.Second),
		socialappAuthenticationMiddleware.Authenticate,
		middleware.RealIP,
		config.cache.Middleware,
	}

	socialappRouter := socialapprouter.NewSocialAppRouter(socialappMiddlewares, routers, authorizationParse, config.newRelicApp)

	propertiesMiddleware := []func(http.Handler) http.Handler{
		cors.AllowAll().Handler,
		middleware.Heartbeat("/health"),
		requestid.Middleware,
		beacon.Middleware,
		middleware.Recoverer,
		middleware.Timeout(60 * time.Second),
		middleware.RealIP,
	}
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse properties target url")
	}

	propertiesProxy := proxyrouter.NewProxyRouter(kibanaTargetURL, propertiesMiddleware)

	localSubdomain := os.Getenv("LOCAL_SUBDOMAIN")
	if localSubdomain == "" {
		// default to google.com
		localSubdomain = "google.com"
	}

	// 3. Main router for routing to different routers based on subdomain
	mainRouter := chi.NewRouter()
	mainRouter.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		switch r.Host {
		case kibanaSubdomain:
			// check for auth cookie
			if cookie, err := r.Cookie("kibanaauthtoken"); err != nil {
				usr, pwd, ok := r.BasicAuth()
				if !ok {
					w.Header().Set("WWW-Authenticate", `Basic realm="`+"Please enter your username and password for this site"+`"`)
					w.WriteHeader(401)
					w.Write([]byte("Unauthorized.\n"))
					return
				} else {
					// add form value to body
					r.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(usr+":"+pwd)))
				}
				log.Warn().Err(err).Msg("failed to get auth cookie")
			} else {
				r.Header.Add("Authorization", "Bearer "+cookie.Value)
			}

			r.Form = url.Values{}
			r.Form.Set("scope", "kibana:read")
			authKibanaRouter.Router.ServeHTTP(w, r)
		case config.propertiesSubdomain.Hostname():
			propertiesProxy.Router.ServeHTTP(w, r)
		case socialappSubdomain:
			socialappRouter.Router.ServeHTTP(w, r)
		case localSubdomain:
			socialappRouter.Router.ServeHTTP(w, r)
		default:
			w.Write([]byte("Host Not found"))
			w.WriteHeader(404)
		}
	})

	log.Info().Msgf("Listening on port %d", config.appPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.appPort), mainRouter); err != nil {
		log.Fatal().Err(err).Msgf("Shutting down")
	}
}

// CreateDBPools creates a pool of connections to the database, in go's implementation of sql, the sql.DB is a connection pool
// but we want to manually control the minimum number of connections to the database
func CreateDBPools(databaseURL string, numPools int) *ForcedConnectionPool {
	pools := make([]*sql.DB, 0, numPools)
	for i := 0; i < numPools; i++ {
		dbConn, err := sql.Open("nrpostgres", databaseURL)
		if err != nil {
			log.Fatal().Err(err)
		}
		// defer close connections will be executed calling ForcedConnectionPool.Close()

		if dbConn == nil {
			log.Fatal().Msg("db is nil")
		}

		// dbConn.SetMaxOpenConns(10)
		// dbConn.SetMaxIdleConns(10)

		pingCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := dbConn.PingContext(pingCtx); err != nil {
			log.Fatal().Err(err).Msg("failed to ping database, shutting down")
		}

		pools = append(pools, dbConn)
	}

	f := &ForcedConnectionPool{
		numPools:          numPools,
		connections:       pools,
		currentRoundRobin: 0,
	}

	return f
}

// ForcedConnectionPool is a wrapper around native Go sql.DB, this allows us to force the minium number of connections
type ForcedConnectionPool struct {
	connections       []*sql.DB
	numPools          int
	currentRoundRobin int
}

func (f *ForcedConnectionPool) GetPool() *sql.DB {
	// round robin
	pool := f.connections[f.currentRoundRobin]
	f.currentRoundRobin += 1
	f.currentRoundRobin = f.currentRoundRobin % f.numPools

	return pool
}

func (f *ForcedConnectionPool) Close() {
	for _, conn := range f.connections {
		conn.Close()
	}
}
