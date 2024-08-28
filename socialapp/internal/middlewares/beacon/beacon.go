package beacon

import (
	"net/http"
	"time"

	"github.com/igomez10/microservices/socialapp/internal/contexthelper"
	"github.com/igomez10/microservices/socialapp/internal/responseWriter"
	"github.com/igomez10/microservices/socialapp/internal/tracerhelper"
	"github.com/rs/zerolog"
)

// Beacon is used to configure a centralized logging platform for the application.
type Beacon struct {
	Logger zerolog.Logger
}

// Middleware returns a handler that logs all the request/response information for
// the centralized logging platform
// Middleware function for handling HTTP requests and logging their responses.
func (b *Beacon) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracerhelper.GetTracer().Start(r.Context(), "middleware.beacon")
		defer span.End()

		r = r.WithContext(ctx)

		customW := responseWriter.NewCustomResponseWriter(w)
		r = r.WithContext(contexthelper.SetLoggerInContext(r.Context(), b.Logger))
		startTime := time.Now()

		// Handle the request
		next.ServeHTTP(customW, r)

		// Log the response
		statusCode := customW.StatusCode
		logEvent := createLogEvent(r, statusCode, startTime, &b.Logger)
		logEvent.Msgf("finished request")
	})
}

// Create a log event based on the request and response.
func createLogEvent(r *http.Request, statusCode int, startTime time.Time, logger *zerolog.Logger) *zerolog.Event {
	requestPattern := contexthelper.GetRequestPatternInContext(r.Context())
	requestID := contexthelper.GetRequestIDInContext(r.Context())
	latency := time.Since(startTime).Milliseconds()

	logevent := logger.WithLevel(zerolog.InfoLevel).
		Str("path", r.URL.Path).
		Str("query", r.URL.RawQuery).
		Str("pattern", requestPattern).
		Str("X-Request-ID", requestID).
		Str("method", r.Method).
		Str("remote_addr", r.RemoteAddr).
		Str("user_agent", r.UserAgent()).
		Str("referer", r.Referer()).
		Str("request_host", r.Host).
		Int("status_code", statusCode).
		Int64("latency_ms", latency)

	switch statusCode {
	case http.StatusUnauthorized:
		username, _ := contexthelper.GetUsernameInContext(r.Context())
		return logevent.
			Str("authorization", r.Header.Get("Authorization")).
			Str("username", username).
			Str("error", "Unauthorized")
	default:
		return logevent
	}
}
