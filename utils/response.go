package utils

import (
	"encoding/json"
	"net/http"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// ErrorResponse represents an error API response
type ErrorResponse struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
}

// JSONResponse writes a JSON response
func JSONResponse(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

// SuccessResponse writes a success JSON response
func SuccessResponse(w http.ResponseWriter, status int, message string, data interface{}) error {
	response := Response{
		Success: true,
		Message: message,
		Data:    data,
	}
	return JSONResponse(w, status, response)
}

// ErrorResponseJSON writes an error JSON response
func ErrorResponseJSON(w http.ResponseWriter, status int, message string, err error) error {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	response := Response{
		Success: false,
		Message: message,
		Error:   errMsg,
	}
	return JSONResponse(w, status, response)
}

// GetUserIDFromRequest retrieves user ID from request headers
func GetUserIDFromRequest(r *http.Request) string {
	return r.Header.Get("X-User-ID")
}

// GetUserEmailFromRequest retrieves user email from request headers
func GetUserEmailFromRequest(r *http.Request) string {
	return r.Header.Get("X-User-Email")
}

// GetUserUsernameFromRequest retrieves username from request headers
func GetUserUsernameFromRequest(r *http.Request) string {
	return r.Header.Get("X-User-Username")
}
