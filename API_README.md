# TIC HCM1 2025 - Login API

This project now includes a simple login API with hardcoded credentials for demonstration purposes.

## Running the API Server

To start the API server, run:

```bash
go run . api
```

Or with a specific port:

```bash
API_PORT=3000 go run . api
```

The server will start on port 8080 by default (or the port specified in the `API_PORT` environment variable).

## Available Endpoints

### POST /api/login
Authenticates a user with username and password.

**Request Body:**
```json
{
  "username": "admin",
  "password": "admin123"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Login successful",
  "user_id": 1
}
```

**Error Response (401):**
```json
{
  "success": false,
  "message": "Invalid username or password"
}
```

### GET /api/health
Health check endpoint.

**Response (200):**
```json
{
  "status": "healthy",
  "timestamp": "2025-01-27T10:30:00Z",
  "service": "TIC HCM1 2025 API"
}
```

## Hardcoded Users

The following users are available for testing:

| Username | Password | User ID |
|----------|----------|---------|
| admin    | admin123 | 1       |
| user1    | password1| 101     |
| demo     | demo123  | 102     |
| testuser | testpass | 103     |
| john_doe | john123  | 104     |

## Testing the API

### Using curl

```bash
# Successful login
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin123"}'

# Failed login
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "wrongpassword"}'

# Health check
curl -X GET http://localhost:8080/api/health
```

### Using the test script

Run the provided test script:

```bash
./test_api.sh
```

This will test various login scenarios and the health endpoint.

## Features

- **Simple Authentication**: No JWT tokens, just basic username/password validation
- **Database Integration**: Updates user's last login time in the database
- **CORS Support**: Allows cross-origin requests
- **Logging**: Logs successful and failed login attempts
- **Error Handling**: Proper HTTP status codes and error messages

## Security Notes

⚠️ **This is for demonstration purposes only!**

- Credentials are hardcoded in the source code
- No password hashing or encryption
- No session management
- No rate limiting
- No HTTPS enforcement

For production use, implement proper security measures including:
- Password hashing (bcrypt, Argon2)
- JWT or session-based authentication
- Rate limiting
- HTTPS
- Input validation and sanitization
- Database connection security

## Running the Original Demo

To run the original AI streak prediction demo instead of the API server:

```bash
go run .
```

This will execute the full demo with database initialization, AI model training, and user analysis. 