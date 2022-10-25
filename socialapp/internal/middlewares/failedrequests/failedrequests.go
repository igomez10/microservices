package failedrequests

import (
	"fmt"
	"net/http"
	"socialapp/internal/contexthelper"
	"socialapp/internal/responseWriter"
	"time"

	"github.com/rs/zerolog"
)

// LogFailedRequestsMiddleware logs failed requests, specially useful for detecting attacks
// with the same token or against the same user
func FailedRequestsMiddleware(next http.Handler) http.Handler {
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
				Str("Authorization", r.Header.Get("Authorization")).
				Str("Username", username).
				Str("Error", "Unauthorized")
		} else {
			logEvent = log.WithLevel(zerolog.InfoLevel)
		}

		logEvent.Str("Path", r.URL.Path).
			Str("Method", r.Method).
			Str("RemoteAddr", r.RemoteAddr).
			Str("UserAgent", r.UserAgent()).
			Str("Referer", r.Referer()).
			Str("Host", r.Host).
			Str("Code", fmt.Sprintf("%d", customW.StatusCode)).
			Int64("Latency_ms", latency).
			Msgf("Finished Request")
	})
}
