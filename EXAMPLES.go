package main

// This file demonstrates example usage patterns for the authentication system.
// It's not meant to be compiled, just for reference.

// ============================================================================
// EXAMPLE 1: Adding More Protected Routes
// ============================================================================

// In your handlers package:

/*
package handlers

import (
	"encoding/json"
	"net/http"
	"zoom-clone/backend/middleware"
	"zoom-clone/backend/utils"
)

// UpdateSettings handles updating user settings
func UpdateSettings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get authenticated user ID from headers
	userID := utils.GetUserIDFromRequest(r)
	if userID == "" {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Parse request
	var req struct {
		Notifications bool   `json:"notifications"`
		Theme        string `json:"theme"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update settings logic...
	// You would typically save this to a database

	utils.SuccessResponse(w, http.StatusOK, "Settings updated", map[string]interface{}{
		"user_id":       userID,
		"notifications": req.Notifications,
		"theme":         req.Theme,
	})
}

// GetSettings handles retrieving user settings
func GetSettings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := utils.GetUserIDFromRequest(r)
	if userID == "" {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Retrieve settings from database...
	settings := map[string]interface{}{
		"user_id":       userID,
		"notifications": true,
		"theme":         "dark",
	}

	utils.SuccessResponse(w, http.StatusOK, "Settings retrieved", settings)
}
*/

// ============================================================================
// EXAMPLE 2: Registering Protected Routes in main.go
// ============================================================================

/*
func main() {
	// ... existing code ...

	// Protected routes with AuthMiddleware
	mux.Handle("/api/user/settings", 
		middleware.AuthMiddleware(http.HandlerFunc(handlers.GetSettings)))
	
	mux.Handle("/api/user/settings/update", 
		middleware.AuthMiddleware(http.HandlerFunc(handlers.UpdateSettings)))

	mux.Handle("/api/user/delete", 
		middleware.AuthMiddleware(http.HandlerFunc(handlers.DeleteAccount)))

	// Optional auth routes
	mux.Handle("/api/posts", 
		middleware.OptionalAuthMiddleware(http.HandlerFunc(handlers.GetPosts)))

	// ... start server ...
}
*/

// ============================================================================
// EXAMPLE 3: Using Auth Utils in Handlers
// ============================================================================

/*
package handlers

import (
	"encoding/json"
	"net/http"
	"zoom-clone/backend/utils"
)

func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	// Get all user info from headers set by auth middleware
	userID := utils.GetUserIDFromRequest(r)
	email := utils.GetUserEmailFromRequest(r)
	username := utils.GetUserUsernameFromRequest(r)

	userInfo := map[string]string{
		"user_id":  userID,
		"email":    email,
		"username": username,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userInfo)
}
*/

// ============================================================================
// EXAMPLE 4: Error Handling Middleware
// ============================================================================

/*
package middleware

import (
	"log"
	"net/http"
	"zoom-clone/backend/utils"
)

// ErrorHandlingMiddleware logs errors and recovers from panics
func ErrorHandlingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				utils.ErrorResponseJSON(w, http.StatusInternalServerError, 
					"Internal server error", nil)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// LoggingMiddleware logs all requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}
*/

// ============================================================================
// EXAMPLE 5: CORS Middleware
// ============================================================================

/*
package middleware

import "net/http"

// CORSMiddleware adds CORS headers to responses
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Usage in main.go:
// mux.Use(middleware.CORSMiddleware)
*/

// ============================================================================
// EXAMPLE 6: Role-Based Access Control (RBAC)
// ============================================================================

/*
package middleware

import (
	"net/http"
	"zoom-clone/backend/utils"
)

// AdminMiddleware ensures user is an admin
func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := utils.GetUserIDFromRequest(r)
		
		// Check if user is admin (from database)
		isAdmin := checkUserIsAdmin(userID)
		if !isAdmin {
			http.Error(w, "Forbidden - Admin access required", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func checkUserIsAdmin(userID string) bool {
	// Query database to check admin status
	// This is a placeholder
	return false
}

// Usage:
// mux.Handle("/api/admin/users", 
//     middleware.AuthMiddleware(middleware.AdminMiddleware(
//         http.HandlerFunc(handlers.ListAllUsers))))
*/

// ============================================================================
// EXAMPLE 7: Rate Limiting for Login Attempts
// ============================================================================

/*
package middleware

import (
	"net/http"
	"sync"
	"time"
	"golang.org/x/time/rate"
)

var (
	loginAttempts = make(map[string]*rate.Limiter)
	mu             sync.Mutex
)

// GetLimiter gets or creates a rate limiter for an email
func GetLimiter(email string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := loginAttempts[email]
	if !exists {
		limiter = rate.NewLimiter(rate.Limit(5), 1) // 5 attempts per minute
		loginAttempts[email] = limiter
	}

	return limiter
}

// RateLimitLogin middleware limits login attempts
func RateLimitLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract email from request
		var req struct {
			Email string `json:"email"`
		}
		// Parse email...

		limiter := GetLimiter(req.Email)
		if !limiter.Allow() {
			http.Error(w, "Too many login attempts, please try again later", 
				http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Usage:
// mux.Handle("/api/auth/login", middleware.RateLimitLogin(
//     http.HandlerFunc(handlers.Login)))
*/

// ============================================================================
// EXAMPLE 8: Integration with React Frontend
// ============================================================================

/*
// React Hook for Authentication

import { useState, useEffect } from 'react';

export const useAuth = () => {
	const [token, setToken] = useState(localStorage.getItem('token'));
	const [user, setUser] = useState(null);
	const [loading, setLoading] = useState(false);

	const register = async (email, password, username) => {
		setLoading(true);
		try {
			const response = await fetch('/api/auth/register', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ email, password, username })
			});
			const data = await response.json();
			
			localStorage.setItem('token', data.token);
			setToken(data.token);
			setUser(data.user);
			return data;
		} catch (error) {
			console.error('Registration failed:', error);
			throw error;
		} finally {
			setLoading(false);
		}
	};

	const login = async (email, password) => {
		setLoading(true);
		try {
			const response = await fetch('/api/auth/login', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ email, password })
			});
			const data = await response.json();
			
			localStorage.setItem('token', data.token);
			setToken(data.token);
			setUser(data.user);
			return data;
		} catch (error) {
			console.error('Login failed:', error);
			throw error;
		} finally {
			setLoading(false);
		}
	};

	const logout = () => {
		localStorage.removeItem('token');
		setToken(null);
		setUser(null);
	};

	const getProfile = async () => {
		if (!token) return;
		
		try {
			const response = await fetch('/api/user/profile', {
				headers: { 'Authorization': `Bearer ${token}` }
			});
			const user = await response.json();
			setUser(user);
			return user;
		} catch (error) {
			console.error('Failed to get profile:', error);
		}
	};

	return { token, user, loading, register, login, logout, getProfile };
};

// Usage in component:
// const { token, user, login, register } = useAuth();
*/

// ============================================================================
// EXAMPLE 9: Token Refresh Loop
// ============================================================================

/*
package middleware

import (
	"context"
	"time"
	"zoom-clone/backend/auth"
)

// StartTokenRefreshLoop periodically refreshes tokens
func StartTokenRefreshLoop(ctx context.Context, token string, refreshInterval time.Duration) {
	ticker := time.NewTicker(refreshInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			newToken, _, err := auth.RefreshToken(token)
			if err != nil {
				// Handle refresh error - possibly redirect to login
				continue
			}
			token = newToken
			// Store new token...
		}
	}
}

// Usage:
// ctx, cancel := context.WithCancel(context.Background())
// go StartTokenRefreshLoop(ctx, token, 12*time.Hour)
*/

// ============================================================================
// EXAMPLE 10: Custom Claims Extension
// ============================================================================

/*
package auth

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
)

// ExtendedClaims extends standard JWT claims with custom fields
type ExtendedClaims struct {
	UserID   string   `json:"user_id"`
	Email    string   `json:"email"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	OrgID    string   `json:"org_id"`
	jwt.RegisteredClaims
}

// GenerateTokenWithClaims generates token with custom claims
func GenerateTokenWithClaims(userID, email, username string, roles []string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	
	claims := &ExtendedClaims{
		UserID:   userID,
		Email:    email,
		Username: username,
		Roles:    roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}
*/

// This file demonstrates advanced patterns and integration examples.
// Refer to this for implementing additional features in your application.
