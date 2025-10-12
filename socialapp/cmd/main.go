package main

import (
	//  With pprof to enable profiling
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/http"
	_ "net/http/pprof"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/exaring/otelpgx"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/hashicorp/go-retryablehttp"
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
	"github.com/igomez10/microservices/socialapp/pkg/dbpgx"
	"github.com/igomez10/microservices/socialapp/socialappapi/openapi"
	urlClient "github.com/igomez10/microservices/urlshortener/generated/clients/go/client"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jessevdk/go-flags"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

type Configuration struct {
	appName               string
	appPort               int
	proxyURL              string
	logLevel              zerolog.Level
	logDestination        io.Writer
	connections           *ForcedConnectionPool
	queries               *dbpgx.Queries
	cache                 *cache.Cache
	propertiesSubdomain   *url.URL
	defaultTimeout        time.Duration
	urlService            *url.URL
	urlShortenerSubdomain string
	socialappSubdomain    string
	instanceID            string
	jwtSecret             string
	puttyknifeDomain      string
	puttyknifeURL         *url.URL
	KibanaSubdomain       string
	KibanaURL             string
	urlShortenerURL       string
	localSubdomain        string
}

func main() {
	var opts struct {
		AppName               string        `short:"n" long:"name" description:"name of the app" default:"socialapp"`
		AppPort               int           `short:"p" long:"port" description:"main port for application" default:"8080"`
		ProxyHost             string        `short:"x" long:"proxy" description:"proxy url, \"http://localhost:9091\"" env:"HTTP_PROXY"`
		LogLevel              string        `short:"l" long:"logLevel" description:"log level info/error/warning" default:"info"`
		LogHost               string        `long:"logHost" description:"log host url" env:"LOGSTASH_HOST"`
		PropertiesSubdomain   string        `long:"propertiesSubdomain" description:"Properties subdomain" env:"PROPERTIES_SUBDOMAIN"`
		DefaultTimeout        time.Duration `long:"defaultTimeout" description:"Default timeout for requests" default:"10s"`
		AgentURL              string        `long:"agentURL" description:"Agent URL" default:"http://localhost:4317"`
		DatabaseURL           string        `long:"databaseURL" description:"Database URL" env:"DATABASE_URL"`
		PuttyknifeDomain      string        `long:"puttyknifeDomain" description:"Puttyknife domain: puttyknife.{...}.com" env:"PUTTYKNIFE_DOMAIN"`
		PuttyknifeURL         string        `long:"puttyknife-url" description:"Puttyknife url" env:"PUTTYKNIFE_URL"`
		UrlShortenerSubdomain string        `long:"urlShortenerSubdomain" description:"URL shortener subdomain" env:"URLSHORTENER_SUBDOMAIN"`
		URLShortenerURL       string        `long:"urlShortenerURL" description:"URL shortener URL" env:"URLSHORTENER_URL"`
		SocialappSubdomain    string        `long:"socialappSubdomain" description:"Socialapp subdomain" env:"SOCIALAPP_SUBDOMAIN"`
		JwtSecret             string        `long:"jwtSecret" description:"jwt secret" env:"JWT_SECRET"`
		RedisURL              string        `long:"redisURL" description:"redis url" env:"REDIS_URL"`
		KibanaSubdomain       string        `long:"kibanaSubdomain" description:"Kibana subdomain" env:"KIBANA_SUBDOMAIN"`
		KibanaURL             string        `long:"kibanaURL" description:"Kibana URL" env:"KIBANA_URL"`
		LocalSubdomain        string        `long:"localSubdomain" description:"Local subdomain" env:"LOCAL_SUBDOMAIN" default:"google.com"`
	}

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

	// setup retryable http client
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

	//set max connections age to 10 seconds
	retryClient.HTTPClient.Transport = &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		// IdleConnTmeout is the maximum amount of time an idle (keep-alive) connection will remain idle before closing itself.
		IdleConnTimeout:       120 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 10 * time.Second,
	}

	retryClient.RetryMax = 10
	retryClient.HTTPClient.Timeout = 15 * time.Second
	retryClient.Backoff = retryablehttp.LinearJitterBackoff

	retryClient.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		// Retry on network errors or 5xx status codes
		if err != nil {
			log.Warn().
				Err(err).
				Stringer("url", resp.Request.URL).
				Str("method", resp.Request.Method).
				Str("status", resp.Status).
				Int("status_code", resp.StatusCode).
				Msg("http retry")
			return true, err
		}

		if resp.StatusCode >= 500 {
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

	// parse log level
	parsedLogLevel, err := zerolog.ParseLevel(opts.LogLevel)
	if err != nil {
		log.Fatal().Err(err).Msgf("Invalid log level, %s", opts.LogLevel)
	}

	var logDestination io.Writer = os.Stdout
	// Validate logHost is a url
	if opts.LogHost != "" {
		u, err := url.Parse(opts.LogHost)
		if err != nil {
			log.Fatal().Err(err).Msgf("failed to parse log host url")
		}

		conn, err := net.Dial(u.Scheme, u.Host)
		if err != nil {
			log.Fatal().Err(err).Msgf("failed to establish connection with log host")
		}

		logDestination = conn
	}

	// parse properties subdomain
	var propertiesSubdomainURL *url.URL
	if len(opts.PropertiesSubdomain) != 0 {
		u, err := url.Parse(opts.PropertiesSubdomain)
		if err != nil {
			log.Fatal().Err(err).Msgf("failed to parse properties subdomain url %s", opts.PropertiesSubdomain)
		}
		propertiesSubdomainURL = u
	}

	// parse url service host
	var urlService *url.URL
	if len(opts.URLShortenerURL) != 0 {
		u, err := url.Parse(opts.URLShortenerURL)
		if err != nil {
			log.Fatal().Err(err).Msgf("failed to parse url service host url %s", opts.UrlShortenerSubdomain)
		}
		urlService = u
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
	// setup tracing
	http.DefaultClient = &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	ctx := context.Background()

	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpointURL(opts.AgentURL),
	)
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to create otlp exporter for tracing %q", opts.AgentURL)
	}

	res, err := resource.New(ctx,
		resource.WithProcessRuntimeDescription(),
		resource.WithFromEnv(),      // Reads OTEL_RESOURCE_ATTRIBUTES and OTEL_SERVICE_NAME
		resource.WithTelemetrySDK(), // SDK info
		resource.WithContainer(),    // Container info, including container.id
		resource.WithContainerID(),
		resource.WithAttributes(attribute.KeyValue{
			Key:   attribute.Key("instance_id"),
			Value: attribute.StringValue(instanceID),
		}),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create resource")
	}

	// Create a new tracer provider with a batch span processor and the otlp exporter.
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)

	// The OTLP exporter sends data to the OpenTelemetry collector.
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

	// Register the tracer provider as the global provider.
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(tp)
	otel.SetMeterProvider(meterProvider)

	// Connect to database
	// force creation of 8 connections, one per service
	connections := CreateDBPools(ctx, opts.DatabaseURL, 5, fmt.Sprintf("%s-%s", opts.AppName, instanceID))
	defer connections.Close()
	redisOpts, err := redis.ParseURL(opts.RedisURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse redis url")
	}

	redisOpts.PoolSize = 10
	redisOpts.MinIdleConns = 10

	cache := cache.NewCache(cache.CacheConfig{
		RedisOpts: redisOpts,
	})
	c := Configuration{
		appPort:               opts.AppPort,
		logLevel:              parsedLogLevel,
		logDestination:        logDestination,
		appName:               opts.AppName,
		connections:           connections,
		queries:               queries,
		cache:                 cache,
		propertiesSubdomain:   propertiesSubdomainURL,
		defaultTimeout:        opts.DefaultTimeout,
		urlService:            urlService,
		puttyknifeDomain:      opts.PuttyknifeDomain,
		puttyknifeURL:         puttyknifeURL,
		urlShortenerSubdomain: opts.UrlShortenerSubdomain,
		socialappSubdomain:    opts.SocialappSubdomain,
		instanceID:            instanceID,
		jwtSecret:             opts.JwtSecret,
		urlShortenerURL:       opts.UrlShortenerSubdomain,
		KibanaSubdomain:       opts.KibanaSubdomain,
		KibanaURL:             opts.KibanaURL,
	}

	// listen on sigkill or sigint to gracefully shutdown the server
	signalChan := make(chan os.Signal, 1)
	run(ctx, c)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
	fmt.Println("Shutting down")

}
func run(_ context.Context, config Configuration) {
	// Setup logger
	zerolog.SetGlobalLevel(config.logLevel)
	log.Logger = zerolog.New(config.logDestination)

	log.Logger = log.With().
		Str("app", config.appName).
		Str("instance", config.instanceID).
		Timestamp().
		Caller().
		Logger()

	// EventRecorder for event sourcing
	eventRecorder := eventRecorder.EventRecorder{
		DB: config.queries,
	}

	// Comment service
	CommentApiService := &comment.CommentService{
		DB:     config.queries,
		DBConn: config.connections.GetPool(),
	}
	CommentApiController := openapi.NewCommentAPIController(CommentApiService)

	// User service
	UserApiService := &user.UserApiService{
		DB:            config.queries,
		DBConn:        config.connections.GetPool(),
		EventRecorder: eventRecorder,
	}
	UserApiController := openapi.NewUserAPIController(UserApiService)

	// Auth service
	AuthApiService := &authentication.AuthenticationService{
		JWTSecret: config.jwtSecret,
	}
	AuthApiController := openapi.NewAuthenticationAPIController(AuthApiService)

	// Role service
	RoleAPIService := &role.RoleApiService{
		DB:     config.queries,
		DBConn: config.connections.GetPool(),
	}
	RoleAPIController := openapi.NewRoleAPIController(RoleAPIService)

	// Scope service
	ScopeAPIService := &scope.ScopeApiService{
		DB:     config.queries,
		DBConn: config.connections.GetPool(),
	}
	ScopeAPIController := openapi.NewScopeAPIController(ScopeAPIService)

	// URL service
	uc := urlClient.NewConfiguration()
	uc.Host = config.urlService.Host
	uc.Scheme = config.urlService.Scheme
	uc.HTTPClient = http.DefaultClient
	uc.UserAgent = config.appName
	urlServiceClient := urlClient.NewAPIClient(uc)

	URLAPIService := &socialappurl.URLApiService{
		Client: urlServiceClient,
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
		"/openapi.yaml": {
			"GET": true,
		},
	}
	socialappAuthenticationMiddleware := gandalf.Middleware{
		DB:               config.queries,
		DBConn:           config.connections.GetPool(),
		AllowlistedPaths: socialappAllowlistedPaths,
		JWTSecret:        config.jwtSecret,
		AllowBasicAuth:   false,
		AuthEndpoint:     "/v1/oauth/token",
	}

	beacon := beacon.Beacon{}

	// open apispec file
	openAPIPath := "openapi.yaml"
	openapiFile, err := os.Open(openAPIPath)
	if err != nil {
		log.Fatal().Err(err).Str("path", openAPIPath).Msg("failed to open openapi file")
	}

	// kibanaTargetURL, err := url.Parse(os.Getenv("KIBANA_URL"))
	kibanaTargetURL, err := url.Parse(config.KibanaURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse target url")
	}

	// 1. Kibana router (proxy)
	kibanaAuthMiddleware := gandalf.Middleware{
		DB:               config.queries,
		DBConn:           config.connections.GetPool(),
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
		middleware.Timeout(20 * time.Second),
		kibanaAuthMiddleware.Authenticate,
		authorizationRuler.Authorize,
		middleware.RealIP,
	}
	kibanaSubdomain := config.KibanaSubdomain
	authKibanaRouter := proxyrouter.NewProxyRouter(kibanaTargetURL, kibanaRouterMiddlewares)

	// 2. SocialApp router
	// read api spec
	content, err := io.ReadAll(io.Reader(openapiFile))
	if err != nil {
		log.Fatal().Err(err).Str("path", openAPIPath).Msg("failed to read openapi file")
	}
	openapiFile.Close()

	// parse api spec
	doc, err := openapi3.NewLoader().LoadFromData(content)
	if err != nil {
		log.Fatal().Err(err)
	}

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
		middleware.Timeout(config.defaultTimeout),
		socialappAuthenticationMiddleware.Authenticate,
		middleware.RealIP,
		// config.cache.Middleware,
	}

	socialappRouter := socialapprouter.NewSocialAppRouter(socialappMiddlewares, routers, authorizationParse, nil)

	propertiesMiddleware := []func(http.Handler) http.Handler{
		cors.AllowAll().Handler,
		middleware.Heartbeat("/health"),
		requestid.Middleware,
		beacon.Middleware,
		middleware.Recoverer,
		middleware.Timeout(20 * time.Second),
		middleware.RealIP,
	}
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse properties target url")
	}
	propertiesProxy := proxyrouter.NewProxyRouter(kibanaTargetURL, propertiesMiddleware)
	puttyknifeAllowlistedPaths := map[string]map[string]bool{
		"/openapi.yaml": {
			"GET": true,
		},
	}
	PuttyKnifeAuthenticationMiddleware := gandalf.Middleware{
		DB:               config.queries,
		DBConn:           config.connections.GetPool(),
		AllowlistedPaths: puttyknifeAllowlistedPaths,
		JWTSecret:        config.jwtSecret,
		AllowBasicAuth:   false,
		AuthEndpoint:     "/v1/oauth/token",
	}
	puttyknifeAuthorizationRuler := authorization.Middleware{
		RequiredScopes: map[string]bool{"puttyknife.properties:read": true},
	}
	puttyKnifeMiddlewares := []func(http.Handler) http.Handler{
		cors.AllowAll().Handler,
		middleware.Heartbeat("/health"),
		requestid.Middleware,
		beacon.Middleware,
		PuttyKnifeAuthenticationMiddleware.Authenticate,
		puttyknifeAuthorizationRuler.Authorize,
		middleware.Recoverer,
		middleware.Timeout(20 * time.Second),
		middleware.RealIP,
	}
	puttyknifeServerProxy := proxyrouter.NewProxyRouter(config.puttyknifeURL, puttyKnifeMiddlewares)

	// URL shortener
	urlShortenerURL, err := url.Parse(config.urlShortenerURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse urlshortener target url")
	}
	urlshortenerMiddleware := []func(http.Handler) http.Handler{
		cors.AllowAll().Handler,
		middleware.Heartbeat("/health"),
		requestid.Middleware,
		beacon.Middleware,
		middleware.Recoverer,
		middleware.Timeout(20 * time.Second),
		middleware.RealIP,
	}
	urlshortenerProxy := proxyrouter.NewProxyRouter(urlShortenerURL, urlshortenerMiddleware)

	// LOCAL
	localSubdomain := config.localSubdomain

	// 3. Main router for routing to different routers based on subdomain
	mainRouter := chi.NewRouter()
	mainRouter.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		log.Info().Str("host", r.Host).
			Str("path", r.URL.Path).
			Msg("Request host")

		switch r.Host {
		case kibanaSubdomain:
			// check for auth cookie
			if cookie, err := r.Cookie("kibanaauthtoken"); err != nil {
				usr, pwd, ok := r.BasicAuth()
				if !ok {
					w.Header().Set("WWW-Authenticate", `Basic realm="`+"Please enter your username and password for this site"+`"`)
					w.WriteHeader(http.StatusUnauthorized)
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
		case config.propertiesSubdomain.Host:
			propertiesProxy.Router.ServeHTTP(w, r)
		case config.socialappSubdomain, config.socialappSubdomain + ":80", config.socialappSubdomain + ":8085", config.socialappSubdomain + ":443":
			socialappRouter.Router.ServeHTTP(w, r)
		case config.urlShortenerSubdomain:
			urlshortenerProxy.Router.ServeHTTP(w, r)
		case localSubdomain:
			socialappRouter.Router.ServeHTTP(w, r)
		case config.puttyknifeDomain:
			puttyknifeServerProxy.Router.ServeHTTP(w, r)
		default:
			socialappRouter.Router.ServeHTTP(w, r)
		}
	})

	log.Info().Msgf("Listening on port %d", config.appPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.appPort), mainRouter); err != nil {
		log.Fatal().Err(err).Msgf("Shutting down")
	}
}

// CreateDBPools creates a pool of connections to the database, in go's implementation of sql, the sql.DB is a connection pool
// but we want to manually control the minimum number of connections to the database
func CreateDBPools(ctx context.Context, databaseURL string, numPools int, applicationName string) *ForcedConnectionPool {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse database url")
	}

	config.MinConns = 3
	config.ConnConfig.ConnectTimeout = 5 * time.Second
	config.ConnConfig.RuntimeParams = map[string]string{
		"application_name": applicationName,
	}

	config.ConnConfig.Tracer = otelpgx.NewTracer(otelpgx.WithTracerProvider(otel.GetTracerProvider()))
	pools := make([]*pgxpool.Pool, 0, numPools)
	for i := 0; i < numPools; i++ {
		dbConn, err := pgxpool.NewWithConfig(ctx, config)
		if err != nil {
			log.Fatal().Err(err)
		}

		if dbConn == nil {
			log.Fatal().Msg("db is nil")
		}

		pingCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := dbConn.Ping(pingCtx); err != nil {
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
	connections       []*pgxpool.Pool
	numPools          int
	currentRoundRobin int
}

func (f *ForcedConnectionPool) GetPool() *pgxpool.Pool {
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
