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
}, []string{"host", "path"})
var prometheusProxyResponseTime = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "proxy_response_time_milliseconds",
	Help: "The response time of proxy requests",
}, []string{"host", "path"})

func NewProxyRouter(target *url.URL, middlewares []func(http.Handler) http.Handler) *ProxyRouter {
	router := chi.NewRouter()

	for i := range middlewares {
		router.Use(middlewares[i])
	}

	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: target.Scheme,
		Host:   target.Host,
	})

	// Serve static files from the static/public folder
	fs := http.FileServer(http.Dir("static/public"))
	fs = http.StripPrefix("/static/public/", fs)

	// Expose the static public folder
	router.HandleFunc("/static/public/*", func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		prometheusProxyRequests.
			WithLabelValues(r.Host, r.URL.Path).
			Inc()

		fs.ServeHTTP(w, r)

		latency := float64(time.Since(startTime).Milliseconds())
		prometheusProxyResponseTime.
			WithLabelValues(r.Host, r.URL.Path).
			Observe(latency)
	})

	router.HandleFunc("/*", func(w http.ResponseWriter, req *http.Request) {
		// metrics for proxy
		startTime := time.Now()
		log.Info().Msgf("Proxying request to %s", req.Host)
		prometheusProxyRequests.
			WithLabelValues(req.Host, req.URL.Path).
			Inc()
		req.Host = target.Host
		req.URL.Host = target.Host

		// remove auth header
		req.Header.Del("Authorization")
		proxy.ServeHTTP(w, req)

		latency := float64(time.Since(startTime).Milliseconds())
		prometheusProxyResponseTime.
			WithLabelValues(req.Host, req.URL.Path).
			Observe(latency)
		return
	})

	return &ProxyRouter{
		Router: router,
		Target: target,
		Proxy:  proxy,
	}
}
