package requestid

import (
	"context"
	"net/http"
	"runtime/pprof"
	"socialapp/internal/contexthelper"

	"github.com/google/uuid"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Request-ID", requestID)
		r.Header.Set("X-Request-ID", requestID)

		// Add request id as attribute to the logger
		r = contexthelper.SetRequestIDInContext(r, requestID)
		log := contexthelper.GetLoggerInContext(r.Context())
		log = log.With().
			Str("X-Request-ID", requestID).
			Logger()

		r = contexthelper.SetLoggerInContext(r, log)

		// ---------
		//  HANDLE REQUEST
		labels := pprof.Labels("path", r.URL.Path)
		pprof.Do(r.Context(), labels, func(ctx context.Context) {
			// Do some work...
			next.ServeHTTP(w, r)
		})

		// HANDLE RESPONSE
		// ---------
	})
}
