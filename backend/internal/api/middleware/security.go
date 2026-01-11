package middleware

import (
	"net/http"
)

// SecurityHeaders adds security headers to all responses.
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Prevent clickjacking
		w.Header().Set("X-Frame-Options", "DENY")

		// Prevent MIME sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// XSS protection
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		// Strict transport security (HTTPS only)
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		// Content security policy
		w.Header().Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'")

		// Referrer policy
		w.Header().Set("Referrer-Policy", "no-referrer")

		// Permissions policy
		w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		next.ServeHTTP(w, r)
	})
}

// NoSecretLogging prevents logging of secret-related request bodies.
func NoSecretLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Mark request as sensitive (for logging middleware)
		r.Header.Set("X-Sensitive-Request", "true")
		next.ServeHTTP(w, r)
	})
}
