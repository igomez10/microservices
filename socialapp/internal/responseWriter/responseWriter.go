package responseWriter

import "net/http"

// custom response writer for capturing status code in the response
type customResponseWriter struct {
	http.ResponseWriter
	StatusCode int
	Body       []byte
}

func NewCustomResponseWriter(w http.ResponseWriter) *customResponseWriter {
	return &customResponseWriter{w, http.StatusOK, []byte{}}
}

func (lrw *customResponseWriter) WriteHeader(code int) {
	lrw.StatusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *customResponseWriter) Write(b []byte) (int, error) {
	lrw.Body = append(lrw.Body, b...)
	return lrw.ResponseWriter.Write(b)
}
