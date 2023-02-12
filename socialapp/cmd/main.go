package main

import (
	//  With pprof to enable profiling
	// _ "net/http/pprof"
	"context"
	"database/sql"
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
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

var (
	appPort *int = flag.Int("port", 8080, "main port for application")
)

func main() {
	flag.Parse()

	useProxy := os.Getenv("USE_PROXY")
	if useProxy == "true" {
		proxyUrl, err := url.Parse("http://localhost:9091")
		if err != nil {
			panic(err)
		}
		http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	}

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

	instanceID := uuid.NewString()
	log.Logger = log.With().
		Str("app", "socialapp").
		Str("instance", instanceID).
		Timestamp().
		Logger()

	log.Info().Msgf("Starting PORT: %d", *appPort)

	// Connect to Kafka
	kafkaBrokers := os.Getenv("KAFKA_HOST")
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":            kafkaBrokers,
		"client.id":                    "go-producer-" + instanceID,
		"acks":                         "all",
		"retries":                      0,
		"linger.ms":                    1,
		"compression.type":             "snappy",
		"batch.num.messages":           1000,
		"queue.buffering.max.messages": 100000,
		"queue.buffering.max.ms":       1000,
		"message.send.max.retries":     3,
		"retry.backoff.ms":             5,
		"socket.keepalive.enable":      true,
		"socket.nagle.disable":         true,
		"socket.max.fails":             3,
		"broker.address.ttl":           1000,
		"broker.address.family":        "v4",
		"api.version.request":          true,
		"api.version.fallback.ms":      0,
		"security.protocol":            "plaintext",
		"ssl.key.location":             "",
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create producer")
	}

	// Connect to database
	// force creation of 8 connections, one per service
	connections := CreateDBPools(os.Getenv("DATABASE_URL"), 8)
	defer connections.Close()
	queries := db.New()

	// EventRecorder for event sourcing
	eventRecorder := eventRecorder.EventRecorder{
		Query: db.New(),
	}

	// Comment service
	CommentApiService := &comment.CommentService{
		DB:            queries,
		DBConn:        connections.GetPool(),
		KafkaProducer: p,
	}
	CommentApiController := openapi.NewCommentApiController(CommentApiService)

	// User service
	UserApiService := &user.UserApiService{
		DB:            queries,
		DBConn:        connections.GetPool(),
		KafkaProducer: p,
		EventRecorder: eventRecorder,
	}
	UserApiController := openapi.NewUserApiController(UserApiService)

	// Auth service
	AuthApiService := &authentication.AuthenticationService{
		DB:     queries,
		DBConn: connections.GetPool(),
	}
	AuthApiController := openapi.NewAuthenticationApiController(AuthApiService)

	// Role service
	RoleAPIService := &role.RoleApiService{
		DB:     queries,
		DBConn: connections.GetPool(),
	}
	RoleAPIController := openapi.NewRoleApiController(RoleAPIService)

	// Scope service
	ScopeAPIService := &scope.ScopeApiService{
		DB:     queries,
		DBConn: connections.GetPool(),
	}
	ScopeAPIController := openapi.NewScopeApiController(ScopeAPIService)

	URLAPIService := &socialappurl.URLApiService{
		DB:     queries,
		DBConn: connections.GetPool(),
	}
	URLAPIController := openapi.NewURLApiController(URLAPIService)

	routers := []openapi.Router{
		CommentApiController,
		UserApiController,
		AuthApiController,
		RoleAPIController,
		ScopeAPIController,
		URLAPIController,
	}

	redisOpts, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse redis url")
	}

	redisOpts.PoolSize = 10
	redisOpts.MinIdleConns = 10

	cache := cache.NewCache(cache.CacheConfig{
		RedisOpts: redisOpts,
	})

	socialappAllowlistedPaths := map[string]map[string]bool{
		"/metrics": {
			"GET": true,
		},
		"/apispec": {
			"GET": true,
		},
	}
	socialappAuthenticationMiddleware := gandalf.Middleware{
		DB:               queries,
		DBConn:           connections.GetPool(),
		Cache:            cache,
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

	kibanaTargetURL, err := url.Parse(os.Getenv("KIBANA_URL"))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse target url")
	}

	// 1. Kibana router (proxy)
	kibanaAuthMiddleware := gandalf.Middleware{
		DB:               queries,
		DBConn:           connections.GetPool(),
		Cache:            cache,
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
	socialappSubdomain := os.Getenv("SOCIALAPP_SUBDOMAIN")
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

	var newrelicApp *newrelic.Application
	newRelicLicense := os.Getenv("NEW_RELIC_LICENSE")
	if newRelicLicense != "" {
		app, err := newrelic.NewApplication(
			newrelic.ConfigAppName("socialapp"),
			newrelic.ConfigLicense(newRelicLicense),
			newrelic.ConfigAppLogForwardingEnabled(true),
			newrelic.ConfigAppLogEnabled(true),
			newrelic.ConfigDistributedTracerEnabled(true),
		)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to create new relic application")
		} else {
			newrelicApp = app
		}
	}

	socialappRouter := socialapprouter.NewSocialAppRouter(socialappMiddlewares, routers, authorizationParse, newrelicApp)

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
	propertiesSubdomain := os.Getenv("PROPERTIES_SUBDOMAIN")
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
		case propertiesSubdomain:
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

	log.Info().Msgf("Listening on port %d", *appPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *appPort), mainRouter); err != nil {
		log.Fatal().Err(err).Msgf("Shutting down")
	}

}

// CreateDBPools creates a pool of connections to the database, in go's implementation of sql, the sql.DB is a connection pool
// but we want to manually control the minimum number of connections to the database
func CreateDBPools(databaseURL string, numPools int) ForcedConnectionPool {
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

		dbConn.SetMaxOpenConns(10)
		dbConn.SetMaxIdleConns(10)

		pingCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := dbConn.PingContext(pingCtx); err != nil {
			log.Fatal().Err(err).Msg("failed to ping database, shutting down")
		}

		pools = append(pools, dbConn)
	}

	f := ForcedConnectionPool{
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
