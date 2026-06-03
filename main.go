package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"zoom-clone/backend/auth"
	"zoom-clone/backend/config"
	"zoom-clone/backend/handlers"
	"zoom-clone/backend/middleware"
)

func main() {
	// Load environment variables
	_ = godotenv.Load()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize JWT configuration
	if err := auth.Init(); err != nil {
		log.Fatalf("Failed to initialize auth: %v", err)
	}

	// Setup routes
	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("/api/auth/register", handlers.Register)
	mux.HandleFunc("/api/auth/login", handlers.Login)
	mux.HandleFunc("/api/auth/refresh", handlers.RefreshToken)

	// Protected routes
	mux.Handle("/api/user/profile", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetProfile)))

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status":"ok"}`)
	})

	address := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server starting on %s (Environment: %s)", address, cfg.Environment)

	// Start server
	if err := http.ListenAndServe(address, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
