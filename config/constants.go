package config

// Database configuration constants
const (
	DefaultJWTExpiration = 24 // hours
	DefaultPort          = "8080"
	DefaultEnvironment   = "development"
)

// Token constants
const (
	BearerScheme     = "Bearer"
	AuthHeaderName   = "Authorization"
	UserIDHeader     = "X-User-ID"
	UserEmailHeader  = "X-User-Email"
	UserUsernameHeader = "X-User-Username"
)

// Error messages
const (
	ErrMissingAuthHeader    = "Authorization header is required"
	ErrInvalidAuthHeader    = "Invalid authorization header format"
	ErrInvalidToken         = "Invalid or expired token"
	ErrUserNotFound         = "User not found"
	ErrUserAlreadyExists    = "User with this email already exists"
	ErrUsernameAlreadyTaken = "Username already taken"
	ErrInvalidCredentials   = "Invalid credentials"
	ErrFailedToHashPassword = "Failed to hash password"
	ErrFailedToGenerateToken = "Failed to generate token"
	ErrUserNotAuthenticated = "User not authenticated"
)

// HTTP status messages
const (
	MsgRegistrationSuccess = "User registered successfully"
	MsgLoginSuccess        = "Login successful"
	MsgTokenRefreshed      = "Token refreshed successfully"
	MsgProfileRetrieved    = "Profile retrieved successfully"
)
