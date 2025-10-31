package middleware

import (
	"log"
	"net/http"
	"time"
)

// Instrumentation logs request metadata and latency for performance budgeting.
func Instrumentation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("http_request method=%s path=%s duration=%s", r.Method, r.URL.Path, time.Since(start))
	})
}
