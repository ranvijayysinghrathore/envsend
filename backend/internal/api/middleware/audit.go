package middleware

import (
	"fmt"
	"net/http"
	"time"
)

// AuditLogger logs all requests for audit purposes (excluding sensitive data).
func AuditLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create response writer wrapper to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Process request
		next.ServeHTTP(wrapped, r)

		// Log request (skip if sensitive)
		isSensitive := r.Header.Get("X-Sensitive-Request") == "true"
		if !isSensitive {
			duration := time.Since(start)
			fmt.Printf("[AUDIT] %s %s %d %v %s\n",
				r.Method,
				r.URL.Path,
				wrapped.statusCode,
				duration,
				getClientIP(r),
			)
		} else {
			// Log minimal info for sensitive requests
			duration := time.Since(start)
			fmt.Printf("[AUDIT] %s %s %d %v [SENSITIVE]\n",
				r.Method,
				r.URL.Path,
				wrapped.statusCode,
				duration,
			)
		}
	})
}

// responseWriter wraps http.ResponseWriter to capture status code.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
