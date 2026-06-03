# API Documentation

Complete API documentation for the Go JWT Authentication System.

## Base URL

```
http://localhost:8080
```

## Authentication

All protected endpoints require a valid JWT token in the Authorization header:

```
Authorization: Bearer <token>
```

## Endpoints

### 1. Register User

Create a new user account.

**Endpoint:** `POST /api/auth/register`

**Headers:**
```
Content-Type: application/json
```

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securePassword123",
  "username": "johndoe"
}
```

**Response (201 Created):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "1",
    "email": "user@example.com",
    "username": "johndoe",
    "created_at": "2024-06-02T15:04:05Z",
    "updated_at": "2024-06-02T15:04:05Z"
  },
  "expires_at": "2024-06-03T15:04:05Z"
}
```

**Error Responses:**

- `400 Bad Request` - Invalid request body
```json
{
  "error": "Invalid request body"
}
```

- `409 Conflict` - Email or username already exists
```json
{
  "error": "User with this email already exists"
}
```

- `500 Internal Server Error` - Server error
```json
{
  "error": "Failed to create user"
}
```

**Example Request:**
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "StrongPassword123",
    "username": "johndoe"
  }'
```

---

### 2. Login User

Authenticate with email and password.

**Endpoint:** `POST /api/auth/login`

**Headers:**
```
Content-Type: application/json
```

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securePassword123"
}
```

**Response (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "1",
    "email": "user@example.com",
    "username": "johndoe",
    "created_at": "2024-06-02T15:04:05Z",
    "updated_at": "2024-06-02T15:04:05Z"
  },
  "expires_at": "2024-06-03T15:04:05Z"
}
```

**Error Responses:**

- `400 Bad Request` - Invalid request body
- `401 Unauthorized` - Invalid credentials
```json
{
  "error": "Invalid credentials"
}
```

**Example Request:**
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "StrongPassword123"
  }'
```

---

### 3. Refresh Token

Generate a new token using an existing token.

**Endpoint:** `POST /api/auth/refresh`

**Headers:**
```
Content-Type: application/json
```

**Request Body:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_at": "2024-06-03T15:04:05Z"
}
```

**Error Responses:**

- `400 Bad Request` - Invalid request body
- `401 Unauthorized` - Invalid or expired token
```json
{
  "error": "Invalid or expired token: token has expired"
}
```

**Example Request:**
```bash
curl -X POST http://localhost:8080/api/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }'
```

---

### 4. Get User Profile

Retrieve the current authenticated user's profile.

**Endpoint:** `GET /api/user/profile`

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Response (200 OK):**
```json
{
  "id": "1",
  "email": "user@example.com",
  "username": "johndoe",
  "created_at": "2024-06-02T15:04:05Z",
  "updated_at": "2024-06-02T15:04:05Z"
}
```

**Error Responses:**

- `401 Unauthorized` - Missing or invalid token
```json
{
  "error": "Authorization header is required"
}
```

- `404 Not Found` - User not found
```json
{
  "error": "User not found"
}
```

**Example Request:**
```bash
curl -X GET http://localhost:8080/api/user/profile \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

---

### 5. Health Check

Check if the server is running.

**Endpoint:** `GET /health`

**Response (200 OK):**
```json
{
  "status": "ok"
}
```

**Example Request:**
```bash
curl -X GET http://localhost:8080/health
```

---

## JWT Token Structure

The JWT token contains the following payload:

```json
{
  "user_id": "1",
  "email": "user@example.com",
  "username": "johndoe",
  "exp": 1717345445,
  "iat": 1717259045
}
```

**Claims:**
- `user_id` - Unique user identifier
- `email` - User's email address
- `username` - User's username
- `exp` - Token expiration time (Unix timestamp)
- `iat` - Token issued at time (Unix timestamp)

---

## Error Handling

### HTTP Status Codes

| Status | Meaning |
|--------|---------|
| 200 | OK - Request successful |
| 201 | Created - Resource created successfully |
| 400 | Bad Request - Invalid request body or parameters |
| 401 | Unauthorized - Missing or invalid authentication |
| 404 | Not Found - Resource not found |
| 409 | Conflict - Resource already exists |
| 500 | Internal Server Error - Server error |

### Error Response Format

```json
{
  "error": "Error message describing what went wrong"
}
```

---

## Request/Response Examples

### JavaScript/Fetch

```javascript
// Register
async function register() {
  const response = await fetch('http://localhost:8080/api/auth/register', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      email: 'user@example.com',
      password: 'password123',
      username: 'johndoe'
    })
  });
  
  const data = await response.json();
  return data.token;
}

// Login
async function login() {
  const response = await fetch('http://localhost:8080/api/auth/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      email: 'user@example.com',
      password: 'password123'
    })
  });
  
  const data = await response.json();
  localStorage.setItem('token', data.token);
  return data;
}

// Access protected route
async function getProfile(token) {
  const response = await fetch('http://localhost:8080/api/user/profile', {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
  
  return await response.json();
}

// Refresh token
async function refreshToken(token) {
  const response = await fetch('http://localhost:8080/api/auth/refresh', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      token: token
    })
  });
  
  const data = await response.json();
  localStorage.setItem('token', data.token);
  return data.token;
}
```

### Python

```python
import requests
import json

BASE_URL = 'http://localhost:8080'

# Register
def register():
    response = requests.post(
        f'{BASE_URL}/api/auth/register',
        json={
            'email': 'user@example.com',
            'password': 'password123',
            'username': 'johndoe'
        }
    )
    return response.json()

# Login
def login():
    response = requests.post(
        f'{BASE_URL}/api/auth/login',
        json={
            'email': 'user@example.com',
            'password': 'password123'
        }
    )
    return response.json()

# Get profile
def get_profile(token):
    headers = {
        'Authorization': f'Bearer {token}'
    }
    response = requests.get(
        f'{BASE_URL}/api/user/profile',
        headers=headers
    )
    return response.json()

# Refresh token
def refresh_token(token):
    response = requests.post(
        f'{BASE_URL}/api/auth/refresh',
        json={'token': token}
    )
    return response.json()
```

### cURL

```bash
# Register
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "username": "johndoe"
  }'

# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'

# Get Profile (replace TOKEN with actual token)
curl -X GET http://localhost:8080/api/user/profile \
  -H "Authorization: Bearer TOKEN"

# Refresh Token
curl -X POST http://localhost:8080/api/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"token": "TOKEN"}'

# Health Check
curl http://localhost:8080/health
```

---

## Rate Limiting

Currently, the API does not implement rate limiting. For production, consider implementing:

1. **Per-endpoint rate limits** - Limit requests to specific endpoints
2. **Per-user rate limits** - Limit based on user ID or IP address
3. **Token bucket algorithm** - Classic rate limiting approach
4. **Sliding window** - More sophisticated rate limiting

Example middleware implementation:

```go
import "golang.org/x/time/rate"

var limiter = rate.NewLimiter(rate.Limit(10), 1) // 10 requests per second

func rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
```

---

## Best Practices

1. **Token Storage** - Store tokens securely (HttpOnly cookies or localStorage)
2. **HTTPS** - Always use HTTPS in production
3. **Token Expiration** - Set appropriate expiration times
4. **Refresh Tokens** - Implement refresh token rotation
5. **Error Messages** - Don't expose sensitive information in error messages
6. **Rate Limiting** - Implement rate limiting for authentication endpoints
7. **Logging** - Log authentication events for security monitoring
8. **CORS** - Configure CORS properly for your frontend
9. **Input Validation** - Validate all user inputs
10. **Password Policy** - Enforce strong password requirements

---

## Troubleshooting

### Token Expired

**Problem:** Getting "token has expired" error

**Solution:** Refresh the token using the `/api/auth/refresh` endpoint

### Invalid Credentials

**Problem:** Login fails with "Invalid credentials"

**Solution:** 
1. Verify email and password are correct
2. Make sure user is registered first
3. Check for typos

### Missing Authorization Header

**Problem:** Getting "Authorization header is required"

**Solution:** Include the Authorization header with Bearer token:
```
Authorization: Bearer <your_token>
```

### CORS Errors

**Problem:** Frontend can't connect to backend

**Solution:** Add CORS middleware to your application or configure frontend proxy

---

For more information, see [README.md](README.md)
