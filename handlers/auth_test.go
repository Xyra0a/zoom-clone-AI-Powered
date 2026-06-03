package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"zoom-clone/backend/models"
)

func TestRegister(t *testing.T) {
	// Clear users for test
	users = make(map[string]*models.User)

	payload := models.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
		Username: "testuser",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(body))
	w := httptest.NewRecorder()

	Register(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	var response models.AuthResponse
	json.NewDecoder(w.Body).Decode(&response)

	if response.User.Email != payload.Email {
		t.Errorf("Expected email %s, got %s", payload.Email, response.User.Email)
	}

	if response.Token == "" {
		t.Error("Expected token to be generated")
	}
}

func TestLogin(t *testing.T) {
	// Clear and setup users for test
	users = make(map[string]*models.User)

	// Register a user first
	payload := models.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
		Username: "testuser",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(body))
	w := httptest.NewRecorder()
	Register(w, req)

	// Now test login
	loginPayload := models.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	loginBody, _ := json.Marshal(loginPayload)
	loginReq := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(loginBody))
	loginW := httptest.NewRecorder()

	Login(loginW, loginReq)

	if loginW.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, loginW.Code)
	}

	var response models.AuthResponse
	json.NewDecoder(loginW.Body).Decode(&response)

	if response.Token == "" {
		t.Error("Expected token to be generated")
	}
}

func TestLoginWithInvalidCredentials(t *testing.T) {
	users = make(map[string]*models.User)

	loginPayload := models.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}

	body, _ := json.Marshal(loginPayload)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(body))
	w := httptest.NewRecorder()

	Login(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}
}
