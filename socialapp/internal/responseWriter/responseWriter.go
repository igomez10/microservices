package responseWriter

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

// custom response writer for capturing status code in the response
type customResponseWriter struct {
	http.ResponseWriter
	StatusCode  int
	Body        []byte
	wroteHeader bool
}

func NewCustomResponseWriter(w http.ResponseWriter) *customResponseWriter {
	return &customResponseWriter{
		ResponseWriter: w,
		StatusCode:     http.StatusOK,
		Body:           []byte{},
		wroteHeader:    false,
	}
}

func (lrw *customResponseWriter) WriteHeader(code int) {
	if lrw.wroteHeader {
		log.Error().Stack().Msg("WriteHeader called twice")
		return // Prevent duplicate WriteHeader calls
	}
	lrw.wroteHeader = true
	lrw.StatusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *customResponseWriter) Write(b []byte) (int, error) {
	if !lrw.wroteHeader {
		lrw.WriteHeader(http.StatusOK)
	}
	lrw.Body = append(lrw.Body, b...)
	return lrw.ResponseWriter.Write(b)
}
