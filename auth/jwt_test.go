package auth

import (
	"testing"

	"zoom-clone/backend/models"
)

func TestHashAndVerifyPassword(t *testing.T) {
	password := "securePassword123"

	// Hash password
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// Verify with correct password
	err = VerifyPassword(hashedPassword, password)
	if err != nil {
		t.Error("Password verification failed with correct password")
	}

	// Verify with incorrect password
	err = VerifyPassword(hashedPassword, "wrongPassword")
	if err == nil {
		t.Error("Password verification should fail with incorrect password")
	}
}

func TestTokenGeneration(t *testing.T) {
	// Initialize JWT
	t.Setenv("JWT_SECRET", "test_secret_key")
	t.Setenv("JWT_EXPIRATION", "24")
	Init()

	user := &models.User{
		ID:       "user_123",
		Email:    "test@example.com",
		Username: "testuser",
	}

	token, expiresAt, err := GenerateToken(user)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if token == "" {
		t.Error("Expected token to be generated")
	}

	if expiresAt.IsZero() {
		t.Error("Expected expiration time to be set")
	}
}

func TestTokenValidation(t *testing.T) {
	// Initialize JWT
	t.Setenv("JWT_SECRET", "test_secret_key")
	t.Setenv("JWT_EXPIRATION", "24")
	Init()

	user := &models.User{
		ID:       "user_123",
		Email:    "test@example.com",
		Username: "testuser",
	}

	token, _, err := GenerateToken(user)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Validate token
	claims, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if claims.UserID != user.ID {
		t.Errorf("Expected user ID %s, got %s", user.ID, claims.UserID)
	}

	if claims.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, claims.Email)
	}

	if claims.Username != user.Username {
		t.Errorf("Expected username %s, got %s", user.Username, claims.Username)
	}
}

func TestInvalidToken(t *testing.T) {
	// Initialize JWT
	t.Setenv("JWT_SECRET", "test_secret_key")
	Init()

	_, err := ValidateToken("invalid_token")
	if err == nil {
		t.Error("Expected error for invalid token")
	}
}
