package beacon

import (
	"net/http"
	"socialapp/internal/contexthelper"
	"socialapp/internal/responseWriter"
	"time"

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
		startTime := time.Now()
		// ---------
		//  HANDLE REQUEST
		next.ServeHTTP(customW, r)
		// HANDLE RESPONSE
		// ---------
		latency := time.Since(startTime).Milliseconds()
		log := contexthelper.GetLoggerInContext(r.Context())
		var logEvent *zerolog.Event
		if customW.StatusCode == http.StatusUnauthorized {
			username, _ := contexthelper.GetUsernameInContext(r.Context())
			logEvent = log.WithLevel(zerolog.ErrorLevel).
				Str("authorization", r.Header.Get("Authorization")).
				Str("username", username).
				Str("error", "Unauthorized")
		} else {
			logEvent = log.WithLevel(zerolog.InfoLevel)
		}

		logEvent.Str("path", r.URL.Path).
			Str("query", r.URL.RawQuery).
			Str("method", r.Method).
			Str("remote_addr", r.RemoteAddr).
			Str("user_agent", r.UserAgent()).
			Str("referer", r.Referer()).
			Str("host", r.Host).
			Int("status_code", customW.StatusCode).
			Str("pattern", contexthelper.GetRoutePatternInContext(r.Context())).
			Str("X-Request-ID", contexthelper.GetRequestIDInContext(r.Context())).
			Int64("latency_ms", latency).
			Msgf("finished request")
	})
}
