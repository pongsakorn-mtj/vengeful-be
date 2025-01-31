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