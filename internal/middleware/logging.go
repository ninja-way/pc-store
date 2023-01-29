package middleware

import (
	"log"
	"net/http"
)

// statusRecorder helps track and record the status code of the response
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

// Logging all requests and response status code to them
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			rec := statusRecorder{w, 200}
			next.ServeHTTP(&rec, r)
			log.Printf("[%s] Status: %d Endpoint: %s", r.Method, rec.status, r.URL)
		})
}
