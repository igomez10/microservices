package beacon

import (
	"net/http"
	"time"

	"github.com/igomez10/microservices/socialapp/internal/contexthelper"
	"github.com/igomez10/microservices/socialapp/internal/responseWriter"
	"github.com/rs/zerolog"
)

type Beacon struct {
	Logger zerolog.Logger
}

// Middleware returns a handler that logs all the request/response information for
// the centralized logging platform
func (b *Beacon) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		customW := responseWriter.NewCustomResponseWriter(w)
		r = r.WithContext(contexthelper.SetLoggerInContext(r.Context(), b.Logger))
		startTime := time.Now()
		// ---------
		//  HANDLE REQUEST
		next.ServeHTTP(customW, r)
		// HANDLE RESPONSE
		// ---------

		requestPattern := contexthelper.GetRequestPatternInContext(r.Context())
		requestID := contexthelper.GetRequestIDInContext(r.Context())
		latency := time.Since(startTime).Milliseconds()
		var logEvent *zerolog.Event
		if customW.StatusCode == http.StatusUnauthorized {
			username, _ := contexthelper.GetUsernameInContext(r.Context())
			logEvent = b.Logger.WithLevel(zerolog.ErrorLevel).
				Str("authorization", r.Header.Get("Authorization")).
				Str("username", username).
				Str("error", "Unauthorized")
		} else {
			logEvent = b.Logger.WithLevel(zerolog.InfoLevel)
		}

		logEvent.Str("path", r.URL.Path).
			Str("query", r.URL.RawQuery).
			Str("pattern", requestPattern).
			Str("X-Request-ID", requestID).
			Str("method", r.Method).
			Str("remote_addr", r.RemoteAddr).
			Str("user_agent", r.UserAgent()).
			Str("referer", r.Referer()).
			Str("request_host", r.Host).
			Int("status_code", customW.StatusCode).
			Int64("latency_ms", latency).
			Msgf("finished request")
	})
}
