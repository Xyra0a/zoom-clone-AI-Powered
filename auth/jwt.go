package auth

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"zoom-clone/backend/models"
)

var (
	ErrInvalidToken    = errors.New("invalid token")
	ErrExpiredToken    = errors.New("token has expired")
	ErrInvalidClaims   = errors.New("invalid token claims")
	secretKey          string
	tokenExpiration    int
)

// Init initializes JWT configuration from environment variables
func Init() error {
	secretKey = os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return errors.New("JWT_SECRET environment variable is required")
	}

	expirationStr := os.Getenv("JWT_EXPIRATION")
	if expirationStr == "" {
		tokenExpiration = 24 // default 24 hours
	} else {
		var err error
		tokenExpiration, err = strconv.Atoi(expirationStr)
		if err != nil {
			return fmt.Errorf("invalid JWT_EXPIRATION value: %w", err)
		}
	}

	return nil
}

// GenerateToken generates a JWT token for the given user
func GenerateToken(user *models.User) (string, time.Time, error) {
	expirationTime := time.Now().Add(time.Duration(tokenExpiration) * time.Hour)

	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"email":    user.Email,
		"username": user.Username,
		"exp":      expirationTime.Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to generate token: %w", err)
	}

	return tokenString, expirationTime, nil
}

// ValidateToken validates and parses a JWT token
func ValidateToken(tokenString string) (*models.Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidClaims
	}

	// Check expiration
	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, ErrInvalidClaims
	}

	if time.Now().Unix() > int64(exp) {
		return nil, ErrExpiredToken
	}

	return &models.Claims{
		UserID:   claims["user_id"].(string),
		Email:    claims["email"].(string),
		Username: claims["username"].(string),
	}, nil
}

// RefreshToken generates a new token using an existing valid token
func RefreshToken(tokenString string) (string, time.Time, error) {
	claims, err := ValidateToken(tokenString)
	if err != nil {
		return "", time.Time{}, err
	}

	user := &models.User{
		ID:       claims.UserID,
		Email:    claims.Email,
		Username: claims.Username,
	}

	return GenerateToken(user)
}
