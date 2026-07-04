// Package server provides the HTTP server setup for socialapp
package server

import (
	"context"
	"encoding/base64"
	"hash/fnv"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
	"github.com/igomez10/microservices/socialapp/pkg/snowflake"
	"github.com/igomez10/microservices/socialapp/socialappapi/openapi"
	urlClient "github.com/igomez10/microservices/urlshortener/generated/clients/go/client"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/cors"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/trace"
)

// Configuration constants
const (
	DefaultMiddlewareTimeout = 20 * time.Second
	CompressionLevel         = 5
)

// Config holds all the configuration for the server
type Config struct {
	InstanceID            string
	AppName               string
	AppPort               int
	ProxyURL              string
	LogLevel              slog.Level
	LogDestinations       []io.Writer
	LoggerProvider        *sdklog.LoggerProvider
	DBPool                *pgxpool.Pool
	Queries               *dbpgx.Queries
	Cache                 *cache.Cache
	PropertiesSubdomain   *url.URL
	DefaultTimeout        time.Duration
	URLShortenerSubdomain string
	URLShortenerURL       *url.URL
	SocialappSubdomain    string
	JWTSecret             string
	PuttyknifeSubDomain   string
	PuttyknifeURL         *url.URL
	KibanaSubdomain       string
	KibanaURL             string
	LocalSubdomain        string
	OpenAPIContent        []byte
}

// getHost safely returns the host from a URL or empty string if nil
func getHost(u *url.URL) string {
	if u == nil {
		return ""
	}
	return u.Host
}

// createServiceGenerator creates a snowflake generator with a unique node ID
// derived from hostname and service name combination
func createServiceGenerator(hostname, serviceName string) (*snowflake.Generator, error) {
	h := fnv.New32a()
	h.Write([]byte(hostname + "-" + serviceName))
	nodeID := int64(h.Sum32() % 1024)
	return snowflake.NewGenerator(nodeID)
}

// NewRouter creates and configures the main chi router
func NewRouter(ctx context.Context, config Config) (chi.Router, error) {
	// Setup logger
	multi := io.MultiWriter(config.LogDestinations...)
	handler := slog.NewJSONHandler(multi, &slog.HandlerOptions{
		AddSource: true,
		Level:     config.LogLevel,
	})
	logger := slog.New(handler).With(
		"app", config.AppName,
		"instance", config.InstanceID,
	)
	slog.SetDefault(logger)

	// Get hostname for snowflake node ID generation
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	// Create separate snowflake generators for each service
	userSnowflakeGen, err := createServiceGenerator(hostname, "user")
	if err != nil {
		return nil, err
	}
	commentSnowflakeGen, err := createServiceGenerator(hostname, "comment")
	if err != nil {
		return nil, err
	}
	roleSnowflakeGen, err := createServiceGenerator(hostname, "role")
	if err != nil {
		return nil, err
	}
	scopeSnowflakeGen, err := createServiceGenerator(hostname, "scope")
	if err != nil {
		return nil, err
	}
	eventSnowflakeGen, err := createServiceGenerator(hostname, "event")
	if err != nil {
		return nil, err
	}

	slog.Info("Initialized snowflake generators for all services", "hostname", hostname)

	// EventRecorder for event sourcing
	eventRec := eventRecorder.EventRecorder{
		DB:                 config.Queries,
		SnowflakeGenerator: eventSnowflakeGen,
	}

	// Comment service
	commentApiService := &comment.CommentService{
		DB:                 config.Queries,
		DBConn:             config.DBPool,
		EventRecorder:      eventRec,
		SnowflakeGenerator: commentSnowflakeGen,
	}
	commentApiController := openapi.NewCommentAPIController(commentApiService)

	// User service
	userApiService := &user.UserApiService{
		DB:                 config.Queries,
		DBConn:             config.DBPool,
		EventRecorder:      eventRec,
		SnowflakeGenerator: userSnowflakeGen,
	}
	userApiController := openapi.NewUserAPIController(userApiService)

	// Auth service
	authApiService := &authentication.AuthenticationService{
		JWTSecret: config.JWTSecret,
	}
	authApiController := openapi.NewAuthenticationAPIController(authApiService)

	// Role service
	roleAPIService := &role.RoleApiService{
		DB:                 config.Queries,
		DBConn:             config.DBPool,
		EventRecorder:      eventRec,
		SnowflakeGenerator: roleSnowflakeGen,
	}
	roleAPIController := openapi.NewRoleAPIController(roleAPIService)

	// Scope service
	scopeAPIService := &scope.ScopeApiService{
		DB:                 config.Queries,
		DBConn:             config.DBPool,
		EventRecorder:      eventRec,
		SnowflakeGenerator: scopeSnowflakeGen,
	}
	scopeAPIController := openapi.NewScopeAPIController(scopeAPIService)

	// URL service
	uc := urlClient.NewConfiguration()
	if config.URLShortenerURL != nil {
		uc.Host = config.URLShortenerURL.Host
		uc.Scheme = config.URLShortenerURL.Scheme
	}
	urlHTTPClient := &http.Client{
		Transport: otelhttp.NewTransport(
			http.DefaultClient.Transport,
			otelhttp.WithTracerProvider(otel.GetTracerProvider()),
			otelhttp.WithPropagators(otel.GetTextMapPropagator()),
			otelhttp.WithSpanOptions(trace.WithAttributes(
				attribute.String("peer.service", "urlshortener"),
			)),
		),
		Timeout: http.DefaultClient.Timeout,
	}
	uc.HTTPClient = urlHTTPClient
	uc.UserAgent = config.AppName
	urlServiceClient := urlClient.NewAPIClient(uc)

	urlAPIService := &socialappurl.URLApiService{
		Client: urlServiceClient,
	}
	urlAPIController := openapi.NewURLAPIController(urlAPIService)

	routers := []openapi.Router{
		commentApiController,
		userApiController,
		authApiController,
		roleAPIController,
		scopeAPIController,
		urlAPIController,
	}

	socialappAllowlistedPaths := map[string]map[string]bool{
		"/openapi.yaml": {
			"GET": true,
		},
	}
	socialappAuthenticationMiddleware := gandalf.Middleware{
		DB:               config.Queries,
		DBConn:           config.DBPool,
		AllowlistedPaths: socialappAllowlistedPaths,
		JWTSecret:        config.JWTSecret,
		AllowBasicAuth:   false,
		AuthEndpoint:     "/v1/oauth/token",
	}

	beaconMW := beacon.Beacon{}

	// Parse API spec for authorization
	var authorizationParse authorizationparser.EndpointAuthorizations
	if len(config.OpenAPIContent) > 0 {
		doc, err := openapi3.NewLoader().LoadFromData(config.OpenAPIContent)
		if err != nil {
			return nil, err
		}
		authorizationParse = authorizationparser.FromOpenAPIToEndpointScopes(doc)
	} else {
		authorizationParse = authorizationparser.EndpointAuthorizations{}
	}

	// Compress responses
	compressor := middleware.NewCompressor(CompressionLevel, "application/json", "application/x-yaml", "gzip", "application/json; charset=UTF-8")

	socialappMiddlewares := []func(http.Handler) http.Handler{
		compressor.Handler,
		cors.AllowAll().Handler,
		middleware.Heartbeat("/health"),
		requestid.Middleware,
		beaconMW.Middleware,
		middleware.Recoverer,
		middleware.Timeout(config.DefaultTimeout),
		socialappAuthenticationMiddleware.Authenticate,
		middleware.RealIP,
	}

	socialappRouter := socialapprouter.NewSocialAppRouter(socialappMiddlewares, routers, authorizationParse)

	// Setup proxy routers (for Kibana, properties, etc.)
	kibanaTargetURL, _ := url.Parse(config.KibanaURL)
	if kibanaTargetURL == nil {
		kibanaTargetURL = &url.URL{Scheme: "http", Host: "localhost:5601"}
	}

	kibanaAuthMiddleware := gandalf.Middleware{
		DB:               config.Queries,
		DBConn:           config.DBPool,
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
		beaconMW.Middleware,
		middleware.Recoverer,
		middleware.Timeout(DefaultMiddlewareTimeout),
		kibanaAuthMiddleware.Authenticate,
		authorizationRuler.Authorize,
		middleware.RealIP,
	}
	authKibanaRouter := proxyrouter.NewProxyRouter(kibanaTargetURL, kibanaRouterMiddlewares)

	propertiesMiddleware := []func(http.Handler) http.Handler{
		cors.AllowAll().Handler,
		middleware.Heartbeat("/health"),
		requestid.Middleware,
		beaconMW.Middleware,
		middleware.Recoverer,
		middleware.Timeout(DefaultMiddlewareTimeout),
		middleware.RealIP,
	}
	propertiesProxy := proxyrouter.NewProxyRouter(kibanaTargetURL, propertiesMiddleware)

	urlshortenerMiddleware := []func(http.Handler) http.Handler{
		cors.AllowAll().Handler,
		middleware.Heartbeat("/health"),
		requestid.Middleware,
		beaconMW.Middleware,
		middleware.Recoverer,
		middleware.Timeout(DefaultMiddlewareTimeout),
		middleware.RealIP,
	}
	urlshortenerProxy := proxyrouter.NewProxyRouter(config.URLShortenerURL, urlshortenerMiddleware)

	puttyknifeAllowlistedPaths := map[string]map[string]bool{
		"/openapi.yaml": {
			"GET": true,
		},
	}
	puttyknifeAuthMiddleware := gandalf.Middleware{
		DB:               config.Queries,
		DBConn:           config.DBPool,
		AllowlistedPaths: puttyknifeAllowlistedPaths,
		JWTSecret:        config.JWTSecret,
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
		beaconMW.Middleware,
		puttyknifeAuthMiddleware.Authenticate,
		puttyknifeAuthorizationRuler.Authorize,
		middleware.Recoverer,
		middleware.Timeout(DefaultMiddlewareTimeout),
		middleware.RealIP,
	}
	var puttyknifeServerProxy *proxyrouter.ProxyRouter
	if config.PuttyknifeURL != nil {
		puttyknifeServerProxy = proxyrouter.NewProxyRouter(config.PuttyknifeURL, puttyKnifeMiddlewares)
	}

	// Main router for routing to different routers based on subdomain
	mainRouter := chi.NewRouter()
	mainRouter.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Request host", "host", r.Host, "path", r.URL.Path)

		switch r.Host {
		case config.KibanaSubdomain:
			// Check for auth cookie
			if cookie, err := r.Cookie("kibanaauthtoken"); err != nil {
				usr, pwd, ok := r.BasicAuth()
				if !ok {
					w.Header().Set("WWW-Authenticate", `Basic realm="`+"Please enter your username and password for this site"+`"`)
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Unauthorized.\n"))
					return
				} else {
					// Add form value to body
					r.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(usr+":"+pwd)))
				}
				slog.Warn("failed to get auth cookie", "error", err)
			} else {
				r.Header.Add("Authorization", "Bearer "+cookie.Value)
			}

			r.Form = url.Values{}
			r.Form.Set("scope", "kibana:read")
			authKibanaRouter.Router.ServeHTTP(w, r)
		case getHost(config.PropertiesSubdomain):
			propertiesProxy.Router.ServeHTTP(w, r)
		case config.SocialappSubdomain, config.SocialappSubdomain + ":80", config.SocialappSubdomain + ":8085", config.SocialappSubdomain + ":443":
			socialappRouter.Router.ServeHTTP(w, r)
		case config.URLShortenerSubdomain:
			urlshortenerProxy.Router.ServeHTTP(w, r)
		case config.LocalSubdomain:
			socialappRouter.Router.ServeHTTP(w, r)
		case config.PuttyknifeSubDomain:
			if puttyknifeServerProxy != nil && puttyknifeServerProxy.Router != nil {
				puttyknifeServerProxy.Router.ServeHTTP(w, r)
			}
		default:
			// Default to socialapp router
			socialappRouter.Router.ServeHTTP(w, r)
		}
	})

	return mainRouter, nil
}
