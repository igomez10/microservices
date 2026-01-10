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

		requestID := uuid.New().String()
		w.Header().Set("X-Request-ID", requestID)
		r.Header.Set("X-Request-ID", requestID)

		ctx = contexthelper.SetRequestIDInContext(ctx, requestID)
		logger := contexthelper.GetLoggerInContext(ctx).With("X-Request-ID", requestID)
		ctx = contexthelper.SetLoggerInContext(ctx, logger)
		r = r.WithContext(ctx)

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
