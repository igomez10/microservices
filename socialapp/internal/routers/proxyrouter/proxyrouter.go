package proxyrouter

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"log/slog"
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
var prometheusProxyResponseTime = promauto.NewSummaryVec(prometheus.SummaryOpts{
	Name: "proxy_response_time_milliseconds",
	Help: "The response time of proxy requests",
	Objectives: map[float64]float64{
		0.5:  0.05,
		0.80: 0.01,
		0.85: 0.01,
		0.90: 0.01,
		0.95: 0.01,
		0.99: 0.001,
	},
}, []string{"host", "path"})

func NewProxyRouter(target *url.URL, middlewares []func(http.Handler) http.Handler) *ProxyRouter {
	router := chi.NewRouter()

	router.Use(otelhttp.NewMiddleware(
		"proxy",
		otelhttp.WithTracerProvider(otel.GetTracerProvider()),
		otelhttp.WithPropagators(otel.GetTextMapPropagator()),
	))

	for i := range middlewares {
		router.Use(middlewares[i])
	}

	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: target.Scheme,
		Host:   target.Host,
	})
	proxy.Transport = otelhttp.NewTransport(
		http.DefaultTransport,
		otelhttp.WithTracerProvider(otel.GetTracerProvider()),
		otelhttp.WithPropagators(otel.GetTextMapPropagator()),
	)

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

	proxyHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// metrics for proxy
		startTime := time.Now()
		prometheusProxyRequests.
			WithLabelValues(req.Host, req.URL.Path).
			Inc()
		req.Host = target.Host
		req.URL.Host = target.Host
		slog.Info("Proxying request", "host", req.Host)
		// remove auth header
		req.Header.Del("Authorization")
		proxy.ServeHTTP(w, req)

		latency := float64(time.Since(startTime).Milliseconds())
		prometheusProxyResponseTime.
			WithLabelValues(req.Host, req.URL.Path).
			Observe(latency)
		return
	})
	router.Handle("/*", proxyHandler)

	return &ProxyRouter{
		Router: router,
		Target: target,
		Proxy:  proxy,
	}
}
