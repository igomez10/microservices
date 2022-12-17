package reverseproxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type ReverseProxy struct {
	KibanaURL url.URL
}

var prometheusProxyRequests = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "proxy_requests_total",
	Help: "The total number of proxy requests",
}, []string{"host"})
var prometheusProxyResponseTime = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "proxy_response_time_milliseconds",
	Help: "The response time of proxy requests",
}, []string{"host"})

func (r *ReverseProxy) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		hostHeader := req.Host
		if hostHeader == "kibana.gomezignacio.com" {
			// metrics for proxy
			prometheusProxyRequests.WithLabelValues(hostHeader).Inc()
			startTime := time.Now()
			proxy := httputil.NewSingleHostReverseProxy(&url.URL{
				Scheme: r.KibanaURL.Scheme,
				Host:   r.KibanaURL.Host,
			})
			proxy.ServeHTTP(w, req)
			prometheusProxyResponseTime.WithLabelValues(hostHeader).Observe(float64(time.Since(startTime).Milliseconds()))

			return
		}
		// ---------
		//  HANDLE REQUEST
		next.ServeHTTP(w, req)
		// HANDLE RESPONSE
		// ---------
	})
}
