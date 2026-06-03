# Go Authentication with JWT

A complete authentication system for the Zoom Clone backend using Go and JWT (JSON Web Tokens).

## Features

- User registration with email validation
- User login with password verification
- JWT token generation and validation
- Token refresh capability
- Bcrypt password hashing
- Authentication middleware
- Protected routes

## Project Structure

```
backend/
├── main.go                 # Application entry point
├── go.mod                  # Go module definition
├── .env.example            # Environment variables template
│
├── models/
│   └── user.go            # User and request/response models
│
├── auth/
│   ├── jwt.go             # JWT token generation and validation
│   └── password.go        # Password hashing and verification
│
├── middleware/
│   └── auth.go            # Authentication middleware
│
└── handlers/
    └── auth.go            # Authentication route handlers
```

## Setup

### 1. Install Dependencies

```bash
go mod download
```

### 2. Configure Environment Variables

Copy `.env.example` to `.env` and update the values:

```bash
cp .env.example .env
```

Edit `.env`:
```
JWT_SECRET=your_super_secret_key_change_this_in_production
JWT_EXPIRATION=24
PORT=8080
DATABASE_URL=postgresql://user:password@localhost:5432/zoom_clone
```

### 3. Run the Server

```bash
go run main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

### 1. Register User

**POST** `/api/auth/register`

Request body:
```json
{
  "email": "user@example.com",
  "password": "securepassword123",
  "username": "johndoe"
}
```

Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "user_20240602150405",
    "email": "user@example.com",
    "username": "johndoe",
    "created_at": "2024-06-02T15:04:05Z",
    "updated_at": "2024-06-02T15:04:05Z"
  },
  "expires_at": "2024-06-03T15:04:05Z"
}
```

### 2. Login User

**POST** `/api/auth/login`

Request body:
```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "user_20240602150405",
    "email": "user@example.com",
    "username": "johndoe",
    "created_at": "2024-06-02T15:04:05Z",
    "updated_at": "2024-06-02T15:04:05Z"
  },
  "expires_at": "2024-06-03T15:04:05Z"
}
```

### 3. Refresh Token

**POST** `/api/auth/refresh`

Request body:
```json
{
  "token": "existing_jwt_token"
}
```

Response:
```json
{
  "token": "new_jwt_token",
  "expires_at": "2024-06-03T15:04:05Z"
}
```

### 4. Get User Profile (Protected)

**GET** `/api/user/profile`

Headers:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

Response:
```json
{
  "id": "user_20240602150405",
  "email": "user@example.com",
  "username": "johndoe",
  "created_at": "2024-06-02T15:04:05Z",
  "updated_at": "2024-06-02T15:04:05Z"
}
```

## Usage Examples

### Using cURL

**Register:**
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "username": "johndoe"
  }'
```

**Login:**
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

**Access Protected Route:**
```bash
curl -X GET http://localhost:8080/api/user/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Using JavaScript/Fetch

```javascript
// Register
const registerResponse = await fetch('http://localhost:8080/api/auth/register', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    email: 'user@example.com',
    password: 'password123',
    username: 'johndoe'
  })
});
const { token } = await registerResponse.json();

// Access protected route
const profileResponse = await fetch('http://localhost:8080/api/user/profile', {
  headers: { 'Authorization': `Bearer ${token}` }
});
const profile = await profileResponse.json();
```

## Key Components

### JWT Token Structure

The JWT token contains the following claims:
- `user_id`: Unique user identifier
- `email`: User's email address
- `username`: User's username
- `exp`: Token expiration time (Unix timestamp)
- `iat`: Token issued at time (Unix timestamp)

### Password Security

- Passwords are hashed using bcrypt with default cost
- Passwords are never stored in plaintext
- Passwords are excluded from JSON responses

### Middleware

#### AuthMiddleware
- Required for protected routes
- Validates JWT token from Authorization header
- Returns 401 Unauthorized if token is invalid or missing
- Stores user claims in request headers for handlers to access

#### OptionalAuthMiddleware
- For routes that work with or without authentication
- Stores user claims if valid token is provided
- Does not reject requests without tokens

## Security Notes

1. **JWT Secret**: Change the `JWT_SECRET` environment variable in production to a strong, random value
2. **HTTPS**: Always use HTTPS in production (not just HTTP)
3. **Token Expiration**: Configure `JWT_EXPIRATION` based on your security requirements
4. **Database Integration**: Replace in-memory user storage with a real database (PostgreSQL, MongoDB, etc.)
5. **Rate Limiting**: Add rate limiting to prevent brute force attacks
6. **CORS**: Configure CORS properly for your frontend domain

## Extending the System

### Adding Database Support

Replace the in-memory `users` map in `handlers/auth.go` with database queries:

```go
// Example with database
user := &models.User{}
err := db.Where("email = ?", req.Email).First(user).Error
```

### Adding More Protected Routes

1. Create handler:
```go
func GetSettings(w http.ResponseWriter, r *http.Request) {
    userID := r.Header.Get("X-User-ID")
    // ... implementation
}
```

2. Register route:
```go
mux.Handle("/api/user/settings", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetSettings)))
```

### Adding Role-Based Access Control

Create a custom middleware that checks user roles:

```go
func AdminMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        userID := r.Header.Get("X-User-ID")
        // Check if user is admin
        if !isAdmin(userID) {
            http.Error(w, "Forbidden", http.StatusForbidden)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

## Dependencies

- `github.com/golang-jwt/jwt/v5` - JWT token generation and validation
- `golang.org/x/crypto` - Bcrypt password hashing
- `github.com/joho/godotenv` - Environment variable loading

## License

MIT
