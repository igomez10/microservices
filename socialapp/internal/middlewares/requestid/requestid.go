package requestid

import (
	"context"
	"net/http"
	"runtime/pprof"

	"github.com/google/uuid"
	"github.com/igomez10/microservices/socialapp/internal/contexthelper"
	"github.com/igomez10/microservices/socialapp/internal/tracerhelper"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracerhelper.GetTracer().Start(r.Context(), "middleware.request_id")
		defer span.End()

		r = r.WithContext(ctx)

		requestID := uuid.New().String()
		w.Header().Set("X-Request-ID", requestID)
		r.Header.Set("X-Request-ID", requestID)

		// Add request id as attribute to the logger
		r = r.WithContext(contexthelper.SetRequestIDInContext(r.Context(), requestID))
		log := contexthelper.GetLoggerInContext(r.Context())
		log = log.With().
			Str("X-Request-ID", requestID).
			Logger()

		r = r.WithContext(contexthelper.SetLoggerInContext(r.Context(), log))

		// ---------
		//  HANDLE REQUEST

		// WITH PPROF PROFILING PYROSCOPE
		labels := pprof.Labels("path", r.URL.Path)
		pprof.Do(r.Context(), labels, func(ctx context.Context) {
			// Do some work...
			next.ServeHTTP(w, r)
		})

		// WITHOUT PPROF
		// next.ServeHTTP(w, r)

		// HANDLE RESPONSE
		// ---------
	})
}
