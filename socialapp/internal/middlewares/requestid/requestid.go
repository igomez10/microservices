package requestid

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"socialapp/internal/contexthelper"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func RequestIDMiddleware(next http.Handler) http.Handler {

	conn, err := net.Dial("udp", os.Getenv("LOGSTASH_HOST"))
	if err != nil {
		log.Warn().Msg(fmt.Sprintf("error dialing logsdtash host udp: %v", err))
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Request-ID", requestID)
		r.Header.Set("X-Request-ID", requestID)
		r = contexthelper.SetRequestIDInContext(r, requestID)

		// setup logger in context
		var log zerolog.Logger
		if conn == nil {
			log = zerolog.New(os.Stdout)
		} else {
			log = zerolog.New(conn)
		}

		log = log.With().
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
