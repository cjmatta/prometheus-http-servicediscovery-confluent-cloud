package middleware

import (
	"net/http"
	"strings"
)

// AuthMiddleware creates a middleware that validates the Authorization header
func AuthMiddleware(apiKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip auth for health endpoint
			if r.URL.Path == "/health" {
				next.ServeHTTP(w, r)
				return
			}

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Unauthorized: Missing Authorization header", http.StatusUnauthorized)
				return
			}

			// Check if the header starts with "Bearer "
			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Unauthorized: Invalid Authorization format", http.StatusUnauthorized)
				return
			}

			// Extract token value
			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token == "" {
				http.Error(w, "Unauthorized: Empty token", http.StatusUnauthorized)
				return
			}

			// Validate token against expected API key
			if token != apiKey {
				http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
				return
			}

			// Token is valid, proceed
			next.ServeHTTP(w, r)
		})
	}
}