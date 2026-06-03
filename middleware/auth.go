package middleware

import (
	"net/http"
	"strings"

	"zoom-clone/backend/auth"
)

// AuthMiddleware is a middleware that validates JWT tokens
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		token := parts[1]

		// Validate token
		claims, err := auth.ValidateToken(token)
		if err != nil {
			http.Error(w, "Invalid or expired token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Store claims in request context for use in handlers
		r.Header.Set("X-User-ID", claims.UserID)
		r.Header.Set("X-User-Email", claims.Email)
		r.Header.Set("X-User-Username", claims.Username)

		next.ServeHTTP(w, r)
	})
}

// OptionalAuthMiddleware validates JWT tokens but doesn't reject if missing
func OptionalAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				token := parts[1]
				claims, err := auth.ValidateToken(token)
				if err == nil {
					r.Header.Set("X-User-ID", claims.UserID)
					r.Header.Set("X-User-Email", claims.Email)
					r.Header.Set("X-User-Username", claims.Username)
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}
