# Database Integration Guide

This guide explains how to integrate a PostgreSQL database with the authentication system.

## Overview

Currently, the authentication system uses in-memory storage for users. For production use, you should replace this with a real database like PostgreSQL, MongoDB, or MySQL.

This guide covers PostgreSQL integration using the `github.com/lib/pq` driver.

## Setup

### 1. Add Database Driver Dependency

```bash
go get github.com/lib/pq
```

Update `go.mod`:
```
require (
	...
	github.com/lib/pq v1.10.9
)
```

### 2. Create Database Schema

Connect to PostgreSQL and create the users table:

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    username VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_email ON users(email);
CREATE INDEX idx_username ON users(username);
```

### 3. Create Database Package

Create `db/database.go`:

```go
package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"zoom-clone/backend/models"
)

var DB *sql.DB

// Connect initializes database connection
func Connect(databaseURL string) error {
	var err error
	DB, err = sql.Open("postgres", databaseURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)

	return nil
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// CreateUser inserts a new user into the database
func CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (email, password, username, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	err := DB.QueryRow(query, user.Email, user.Password, user.Username, user.CreatedAt, user.UpdatedAt).
		Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetUserByEmail retrieves a user by email
func GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, password, username, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	err := DB.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.Password, &user.Username, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func GetUserByID(id string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, password, username, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	err := DB.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.Password, &user.Username, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetUserByUsername retrieves a user by username
func GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, password, username, created_at, updated_at
		FROM users
		WHERE username = $1
	`

	err := DB.QueryRow(query, username).Scan(
		&user.ID, &user.Email, &user.Password, &user.Username, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// UpdateUser updates user information
func UpdateUser(user *models.User) error {
	query := `
		UPDATE users
		SET email = $1, password = $2, username = $3, updated_at = $4
		WHERE id = $5
	`

	result, err := DB.Exec(query, user.Email, user.Password, user.Username, user.UpdatedAt, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// DeleteUser deletes a user
func DeleteUser(id string) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
```

### 4. Update Handlers to Use Database

Modify `handlers/auth.go` to use database instead of in-memory storage:

```go
package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"zoom-clone/backend/auth"
	"zoom-clone/backend/db"
	"zoom-clone/backend/models"
)

// Register handles user registration
func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if user already exists
	existingUser, err := db.GetUserByEmail(req.Email)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if existingUser != nil {
		http.Error(w, "User with this email already exists", http.StatusConflict)
		return
	}

	// Check username
	existingUser, err = db.GetUserByUsername(req.Username)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if existingUser != nil {
		http.Error(w, "Username already taken", http.StatusConflict)
		return
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Create user
	user := &models.User{
		Email:     req.Email,
		Password:  hashedPassword,
		Username:  req.Username,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save to database
	if err := db.CreateUser(user); err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Generate token
	token, expiresAt, err := auth.GenerateToken(user)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Return response without password
	user.Password = ""
	response := models.AuthResponse{
		Token:     token,
		User:      user,
		ExpiresAt: expiresAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Login handles user login
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Find user by email
	user, err := db.GetUserByEmail(req.Email)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Verify password
	if err := auth.VerifyPassword(user.Password, req.Password); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate token
	token, expiresAt, err := auth.GenerateToken(user)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Return response without password
	user.Password = ""
	response := models.AuthResponse{
		Token:     token,
		User:      user,
		ExpiresAt: expiresAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ... rest of the handlers remain the same
```

### 5. Update main.go

Add database initialization:

```go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"zoom-clone/backend/auth"
	"zoom-clone/backend/config"
	"zoom-clone/backend/db"
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

	// Initialize database
	if err := db.Connect(cfg.DatabaseURL); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize JWT
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

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status":"ok"}`)
	})

	address := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server starting on %s (Environment: %s)", address, cfg.Environment)

	if err := http.ListenAndServe(address, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
```

## Using Docker Compose

The `docker-compose.yml` file includes a PostgreSQL service. Update your `.env`:

```
DATABASE_URL=postgresql://user:password@postgres:5432/zoom_clone
DB_USER=user
DB_PASSWORD=password
DB_NAME=zoom_clone
```

Then run:

```bash
docker-compose up
```

## Migration Tools

For managing database schema changes, consider using:

1. **Flyway** - SQL-based migrations
2. **golang-migrate** - Go-based migrations
3. **GORM** - ORM with auto-migration support

Example with golang-migrate:

```bash
go get -u github.com/golang-migrate/migrate/cmd/migrate
migrate create -ext sql -dir db/migrations -seq create_users_table
```

## Performance Optimization

### Connection Pooling

```go
DB.SetMaxOpenConns(25)    // Maximum connections
DB.SetMaxIdleConns(5)     // Idle connections to keep
DB.SetConnMaxLifetime(5 * time.Minute)  // Connection lifetime
```

### Indexing

Add indexes for frequently queried fields:

```sql
CREATE INDEX idx_email ON users(email);
CREATE INDEX idx_username ON users(username);
```

### Query Optimization

Use prepared statements for repeated queries to benefit from query caching.

## Error Handling

Always handle database errors appropriately:

```go
if err == sql.ErrNoRows {
	// User not found
} else if err != nil {
	// Other database error
}
```

## Security Best Practices

1. Use parameterized queries (already done in examples above)
2. Never expose database errors to clients
3. Hash passwords before storing
4. Use prepared statements
5. Implement connection timeouts
6. Keep database credentials in environment variables
7. Use HTTPS for all connections

## Next Steps

1. Implement password reset functionality
2. Add email verification
3. Implement refresh token rotation
4. Add audit logging for authentication events
5. Implement rate limiting
6. Add user roles and permissions
