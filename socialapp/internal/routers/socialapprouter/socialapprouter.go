package socialapprouter

import (
	"net/http"
	"os"
	"socialapp/internal/authorizationparser"
	"socialapp/internal/middlewares/authorization"
	"socialapp/internal/middlewares/pattern"
	"socialapp/socialappapi/openapi"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"github.com/slok/go-http-metrics/metrics/prometheus"
	metricsMiddleware "github.com/slok/go-http-metrics/middleware"
	"github.com/slok/go-http-metrics/middleware/std"
)

type SocialAppRouter struct {
	Router chi.Router
}

type Middleware func(http.Handler) http.Handler

func NewSocialAppRouter(middlewares []func(http.Handler) http.Handler, routers []openapi.Router, authorizationParse authorizationparser.EndpointAuthorizations, newrelicApp *newrelic.Application) SocialAppRouter {
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

					pattern := pattern.Pattern{Pattern: route.Pattern}
					r.Use(pattern.Middleware)

					for i := range middlewares {
						r.Use(middlewares[i])
					}

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
					_, wrappedHandler := newrelic.WrapHandle(newrelicApp, route.Pattern, handler)
					r.Method(route.Method, route.Pattern, wrappedHandler)
				})
			}
		}
	})
	s := SocialAppRouter{
		Router: mainRouter,
	}
	return s
}
