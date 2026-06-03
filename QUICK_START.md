# Quick Start Guide

Get up and running with the Go JWT authentication system in 5 minutes.

## 1. Install Go

Make sure you have Go 1.21 or later installed:
```bash
go version
```

## 2. Setup Environment Variables

```bash
cp .env.example .env
```

Edit `.env` and set your JWT secret (use a strong random string in production):
```
JWT_SECRET=my_super_secret_key_12345
JWT_EXPIRATION=24
PORT=8080
```

## 3. Download Dependencies

```bash
go mod download
```

Or using Make:
```bash
make install
```

## 4. Run the Server

```bash
go run main.go
```

Or using Make:
```bash
make run
```

You should see:
```
Server starting on :8080 (Environment: development)
```

## 5. Test the API

### Create a user account

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "username": "testuser"
  }'
```

Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "user_20240602150405",
    "email": "test@example.com",
    "username": "testuser",
    "created_at": "2024-06-02T15:04:05Z",
    "updated_at": "2024-06-02T15:04:05Z"
  },
  "expires_at": "2024-06-03T15:04:05Z"
}
```

Save the token from the response.

### Access protected route

Use the token to access protected endpoints:

```bash
curl -X GET http://localhost:8080/api/user/profile \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### Login

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

### Refresh Token

```bash
curl -X POST http://localhost:8080/api/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "token": "YOUR_EXISTING_TOKEN"
  }'
```

## Project Structure

```
backend/
├── main.go              # Entry point
├── config/              # Configuration management
├── auth/                # JWT and password utilities
├── models/              # Data structures
├── handlers/            # Route handlers
├── middleware/          # HTTP middleware
├── go.mod               # Module definition
├── .env.example         # Environment template
└── README.md            # Full documentation
```

## Common Commands

```bash
# Run server
make run

# Build binary
make build

# Run tests
make test

# Format code
make fmt

# Clean artifacts
make clean

# See all available commands
make help
```

## Next Steps

1. Read [README.md](README.md) for complete documentation
2. Configure a database (PostgreSQL/MongoDB) for user storage
3. Add role-based access control (RBAC)
4. Implement refresh token rotation
5. Add rate limiting for login attempts
6. Setup CORS for your frontend
7. Add email verification
8. Implement password reset functionality

## Troubleshooting

**Server won't start**
- Check `.env` file exists and `JWT_SECRET` is set
- Make sure port 8080 is not in use: `netstat -ano | findstr :8080`

**Invalid credentials error**
- Verify email and password are correct
- Make sure you registered the user first

**Token expired**
- Generate a new token using `/api/auth/refresh` endpoint
- Or login again with credentials

**CORS errors in frontend**
- Add CORS middleware to handle cross-origin requests
- Check browser console for specific CORS issues

For more help, see [README.md](README.md)
