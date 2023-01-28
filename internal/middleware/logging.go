package middleware

import (
	"log"
	"net/http"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			rec := statusRecorder{w, 200}
			next.ServeHTTP(&rec, r)
			log.Printf("[%s] Status: %d %s", r.Method, rec.status, r.URL)
		})
}
