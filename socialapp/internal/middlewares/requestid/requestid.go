package requestid

import (
	"net/http"
	"os"
	"socialapp/internal/contexthelper"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()
		w.Header().Set("X-Request-ID", requestID)
		r.Header.Set("X-Request-ID", requestID)
		r = contexthelper.SetRequestIDInContext(r, requestID)

		// setup logger in context
		log := zerolog.
			New(os.Stdout).
			With().
			Str("X-Request-ID", requestID).
			Logger()
		r = contexthelper.SetLoggerInContext(r, log)

		// ---------
		//  HANDLE REQUEST
		next.ServeHTTP(w, r)
		// HANDLE RESPONSE
		// ---------
	})
}
