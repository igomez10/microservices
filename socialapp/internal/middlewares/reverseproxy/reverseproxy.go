package reverseproxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type ReverseProxy struct {
	KibanaURL url.URL
}

func (r *ReverseProxy) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		hostHeader := req.Host
		if hostHeader == "kibana.gomezignacio.com" {

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
