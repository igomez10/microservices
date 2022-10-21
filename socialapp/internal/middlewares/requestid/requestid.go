package requestid

import (
	"net/http"
	"socialapp/internal/contexthelper"

	"github.com/google/uuid"
)

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()
		w.Header().Set("X-Request-ID", requestID)
		r.Header.Set("X-Request-ID", requestID)
		r = contexthelper.SetRequestIDInContext(r, uuid.NewString())

		// ---------
		//  HANDLE REQUEST
		next.ServeHTTP(w, r)
		// HANDLE RESPONSE
		// ---------
	})
}
