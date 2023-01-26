package proxyrouter

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog/log"
)

type ProxyRouter struct {
	Router chi.Router
	Host   string
	Target *url.URL
	Proxy  *httputil.ReverseProxy
}

var prometheusProxyRequests = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "proxy_requests_total",
	Help: "The total number of proxy requests",
}, []string{"host"})
var prometheusProxyResponseTime = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "proxy_response_time_milliseconds",
	Help: "The response time of proxy requests",
}, []string{"host"})

func NewProxyRouter(target *url.URL, middlewares []func(http.Handler) http.Handler) *ProxyRouter {
	router := chi.NewRouter()

	for i := range middlewares {
		router.Use(middlewares[i])
	}

	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: target.Scheme,
		Host:   target.Host,
	})

	router.HandleFunc("/*", func(w http.ResponseWriter, req *http.Request) {
		// metrics for proxy
		startTime := time.Now()
		log.Info().Msgf("Proxying request to %s", req.Host)
		prometheusProxyRequests.WithLabelValues(req.Host).Inc()
		req.Host = target.Host
		req.URL.Host = target.Host

		// remove auth header
		req.Header.Del("Authorization")
		proxy.ServeHTTP(w, req)

		latency := float64(time.Since(startTime).Milliseconds())
		prometheusProxyResponseTime.
			WithLabelValues(req.Host).
			Observe(latency)
		return
	})

	return &ProxyRouter{
		Router: router,
		Target: target,
		Proxy:  proxy,
	}
}
