package reverseproxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type ReverseProxy struct {
	KibanaURL url.URL
}

var prometheusProxyRequests = promauto.NewCounter(prometheus.CounterOpts{
	Name: "proxy_requests",
	Help: "Number of requests to the proxy",
})

func (r *ReverseProxy) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		hostHeader := req.Host
		if hostHeader == "kibana.gomezignacio.com" {
			// metrics for proxy
			prometheusProxyRequests.Inc()

			proxy := httputil.NewSingleHostReverseProxy(&url.URL{
				Scheme: r.KibanaURL.Scheme,
				Host:   r.KibanaURL.Host,
			})

			proxy.ServeHTTP(w, req)
			return

		}
		// ---------
		//  HANDLE REQUEST
		next.ServeHTTP(w, req)
		// HANDLE RESPONSE
		// ---------
	})
}
