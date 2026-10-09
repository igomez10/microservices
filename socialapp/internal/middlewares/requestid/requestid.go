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
		// trace_id is what Grafana's Loki derived field matches to link a log
		// line to its trace in Tempo; every request-scoped logger inherits it
		logger := contexthelper.GetLoggerInContext(ctx).With(
			"X-Request-ID", requestID,
			"trace_id", span.SpanContext().TraceID().String(),
			"span_id", span.SpanContext().SpanID().String(),
		)
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
