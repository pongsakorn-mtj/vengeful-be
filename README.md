# Vengeful Backend API

## Overview
This is a secure REST API service built with Go, featuring user registration and management with MongoDB integration.

## Features
- User registration endpoint
- User listing with pagination
- Rate limiting
- Secure input validation
- Logging with logrus
- CORS protection
- MongoDB integration

## API Documentation

### Register User
Register a new user in the system.

```
POST /api/register
```

#### Request Body
```json
{
    "firstName": "John",
    "lastName": "Doe",
    "phoneNo": "0812345678",
    "email": "john.doe@example.com",
    "isAcceptTnc": true,
    "isAcceptPrivacyPolicy": true
}
```

#### Validation Rules
| Field | Type | Rules |
|-------|------|-------|
| firstName | string | Required, 2-50 characters |
| lastName | string | Required, 2-50 characters |
| phoneNo | string | Required, 10-15 characters |
| email | string | Required, valid email format |
| isAcceptTnc | boolean | Required, must be true |
| isAcceptPrivacyPolicy | boolean | Required, must be true |

#### Responses

##### Success Response (200 OK)
```json
{
    "status": "success",
    "message": "Registration successful"
}
```

##### Error Responses

###### Invalid Request Format (400 Bad Request)
```json
{
    "status": "error",
    "message": "Invalid request format"
}
```

###### Terms Not Accepted (400 Bad Request)
```json
{
    "status": "error",
    "message": "Must accept terms and conditions and privacy policy"
}
```

###### Email Already Exists (409 Conflict)
```json
{
    "status": "error",
    "message": "Email already registered"
}
```

###### Server Error (500 Internal Server Error)
```json
{
    "status": "error",
    "message": "Internal server error"
}
```

### Get All Users
Retrieve a paginated list of all users.

```
GET /api/users
```

#### Query Parameters
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| page | integer | 1 | Page number |
| limit | integer | 10 | Items per page (max: 100) |

#### Success Response (200 OK)
```json
{
    "status": "success",
    "message": "Users retrieved successfully",
    "data": {
        "users": [
            {
                "firstName": "John",
                "lastName": "Doe",
                "phoneNo": "0812345678",
                "email": "john.doe@example.com",
                "isAcceptTnc": true,
                "isAcceptPrivacyPolicy": true,
                "createdAt": "2024-01-31T10:00:00Z"
            }
        ],
        "pagination": {
            "currentPage": 1,
            "totalPages": 5,
            "totalRecords": 50,
            "limit": 10
        }
    }
}
```

#### Error Response (500 Internal Server Error)
```json
{
    "status": "error",
    "message": "Failed to get users"
}
```

## Security Features
- Rate limiting per IP (60 requests per minute)
- Input validation
- CORS protection
- Request logging
- Secure headers
- MongoDB authentication

## Run Instructions

### Development
```bash
# Run directly with Go
go run cmd/main.go

# Or build and run
go build -o vengeful-be cmd/main.go
./vengeful-be
```

### Production
```bash
# Build optimized binary
go build -tags netgo -ldflags '-s -w' -o vengeful-be cmd/main.go

# Run with production environment
GIN_MODE=release ./vengeful-be

# Run in background with nohup
nohup GIN_MODE=release ./vengeful-be > app.log 2>&1 &

# Run with specific port
PORT=3000 GIN_MODE=release ./vengeful-be

# Run with PM2 (if installed)
pm2 start ./vengeful-be --name "vengeful-api"
```

### Environment Variables
```bash
# Required variables
export MONGO_ROOT_USERNAME=your_username
export MONGO_ROOT_PASSWORD=your_password
export MONGODB_DATABASE=your_database
export MONGODB_COLLECTION_NAME=users
export MONGODB_HOST=your_mongodb_host
export PORT=8080

# Optional variables
export GIN_MODE=release     # Set Gin to release mode
export LOG_LEVEL=info      # Set log level (debug, info, warn, error)
```

### Using Docker (if containerized)
```bash
# Build Docker image
docker build -t vengeful-be .

# Run container
docker run -d \
  -p 8080:8080 \
  -e MONGO_ROOT_USERNAME=your_username \
  -e MONGO_ROOT_PASSWORD=your_password \
  -e MONGODB_DATABASE=your_database \
  -e MONGODB_HOST=your_mongodb_host \
  --name vengeful-api \
  vengeful-be
```

### Health Check
```bash
# Check if service is running
curl http://localhost:8080/api/health

# Monitor logs
tail -f app.log
```

## Build Instructions

### Development Build
```bash
go build -o vengeful-be cmd/main.go
```

### Production Build
```bash
# Optimized build with reduced binary size
go build -tags netgo -ldflags '-s -w' -o vengeful-be cmd/main.go

# Build for specific platform (example: Linux AMD64)
GOOS=linux GOARCH=amd64 go build -tags netgo -ldflags '-s -w' -o vengeful-be cmd/main.go

# Build for multiple platforms
# Linux
GOOS=linux GOARCH=amd64 go build -tags netgo -ldflags '-s -w' -o vengeful-be-linux-amd64 cmd/main.go
# macOS
GOOS=darwin GOARCH=amd64 go build -tags netgo -ldflags '-s -w' -o vengeful-be-darwin-amd64 cmd/main.go
# Windows
GOOS=windows GOARCH=amd64 go build -tags netgo -ldflags '-s -w' -o vengeful-be-windows-amd64.exe cmd/main.go
```

### Running the Built Binary

#### Linux/macOS
```bash
# Basic run
./vengeful-be

# Run with environment variables
MONGO_ROOT_USERNAME=user MONGO_ROOT_PASSWORD=pass ./vengeful-be

# Run in production mode
GIN_MODE=release ./vengeful-be

# Run with custom port
PORT=3000 ./vengeful-be

# Run in background
nohup ./vengeful-be > app.log 2>&1 &

# Run with all options
GIN_MODE=release PORT=3000 MONGO_ROOT_USERNAME=user MONGO_ROOT_PASSWORD=pass ./vengeful-be
```

#### Windows
```bash
# Basic run
vengeful-be.exe

# Run with environment variables (PowerShell)
$env:MONGO_ROOT_USERNAME="user"; $env:MONGO_ROOT_PASSWORD="pass"; .\vengeful-be.exe

# Run in production mode
$env:GIN_MODE="release"; .\vengeful-be.exe

# Run with custom port
$env:PORT="3000"; .\vengeful-be.exe
```

Build flags explanation:
- `-tags netgo`: Forces the use of Go's built-in DNS resolver
- `-ldflags '-s -w'`: Reduces binary size by removing debug information and symbol tables
- `-o vengeful-be`: Specifies the output binary name

## Setup
1. Install Go 1.21 or later
2. Clone the repository
3. Set up environment variables in `.env` file:
   ```
   MONGO_ROOT_USERNAME=your_username
   MONGO_ROOT_PASSWORD=your_password
   MONGODB_DATABASE=your_database
   MONGODB_COLLECTION_NAME=users
   MONGODB_HOST=your_mongodb_host
   PORT=8080
   ```
4. Run `go mod download` to install dependencies
5. Run `go run cmd/main.go` to start the server

## Development
- The API uses Gin framework for routing
- MongoDB for data storage
- Logrus for structured logging
- Rate limiting middleware for request throttling
- Custom validation for input data 