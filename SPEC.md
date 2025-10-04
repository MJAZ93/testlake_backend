# Go REST API Architecture Specification

## Overview

This document serves as a comprehensive architecture specification for building REST API applications using Go. The architecture follows a clean, modular design pattern that promotes separation of concerns, testability, and maintainability. This structure has been successfully implemented in production environments and can be used as a template for new projects.

## Architecture Pattern

The architecture follows a **layered architecture** pattern with clear separation of concerns:

```
┌─────────────────┐
│   Client/UI     │ HTTP Requests
│   (Frontend)    │
└─────────────────┘
         │
         ▼
┌─────────────────┐
│   Middleware    │ Authentication, CORS, Logging
│ (Auth, CORS)    │
└─────────────────┘
         │
         ▼
┌─────────────────┐
│    Service      │ Route Registration & Swagger Docs
│  (Route Layer)  │
└─────────────────┘
         │
         ▼
┌─────────────────┐
│   Controller    │ HTTP Request Handling & Business Logic
│ (HTTP Handlers) │
└─────────────────┘
         │
         ▼
┌─────────────────┐
│      DAO        │ Database Operations & Queries
│ (Data Access)   │
└─────────────────┘
         │
         ▼
┌─────────────────┐
│   Database      │ Data Persistence
│  (PostgreSQL)   │
└─────────────────┘

Supporting Components:
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│     Models      │    │   Input/Output  │    │     Utils       │
│ (GORM Entities) │    │ (Request/Response)│    │   (Helpers)     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## Project Structure

```
project_root/
├── app/                    # Application setup and routing
│   ├── app.go             # Main application server setup
│   ├── routes.go          # Route definitions (public/private)
│   └── swagger.go         # Swagger documentation setup
├── controller/            # HTTP request handlers (MVC Controllers)
│   ├── auth_controller.go
│   ├── user_controller.go
│   └── ...
├── dao/                   # Data Access Objects
│   ├── main.go           # Database connection and migrations
│   ├── user_dao.go
│   └── ...
├── docs/                  # Auto-generated Swagger documentation
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── inout/                 # Input/Output data structures
│   ├── base_response.go
│   ├── auth/
│   ├── user/
│   └── ...
├── middleware/            # HTTP middleware
│   ├── cors.go
│   ├── jwt.go
│   └── ...
├── model/                 # Database models (GORM entities)
│   ├── user.go
│   ├── organization.go
│   └── ...
├── service/               # Service layer (route registration)
│   ├── auth_service.go
│   ├── user_service.go
│   └── ...
├── templates/             # HTML templates
│   └── email_confirmation.html
├── utils/                 # Utility functions
│   ├── error_handler.go
│   ├── jwt.go
│   ├── password.go
│   └── ...
├── main.go               # Application entry point
├── go.mod                # Go module dependencies
└── docker-compose.yml    # Docker setup
```

## Core Dependencies

### Required Dependencies (go.mod)

```go
module your-project-name

go 1.24

require (
    github.com/gin-gonic/gin v1.10.0           // Web framework
    github.com/golang-jwt/jwt/v5 v5.2.1        // JWT authentication
    github.com/google/uuid v1.6.0              // UUID generation
    github.com/joho/godotenv v1.5.1            // Environment variables
    github.com/stretchr/testify v1.9.0         // Testing framework
    github.com/swaggo/files v1.0.1             // Swagger files
    github.com/swaggo/gin-swagger v1.6.0       // Swagger integration
    github.com/swaggo/swag v1.16.3             // Swagger generator
    golang.org/x/crypto v0.28.0               // Cryptography
    gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df // Email
    gorm.io/driver/postgres v1.5.6            // PostgreSQL driver
    gorm.io/gorm v1.30.0                      // ORM
)
```

## Package Details

### 1. `app/` - Application Bootstrap

**Purpose**: Application initialization, server setup, and route configuration.

**Key Files**:
- `app.go`: Main server setup with middleware configuration
- `routes.go`: Route definitions separated into public and private routes
- `swagger.go`: Swagger documentation configuration

**Example (`app/app.go`)**:
```go
package app

import (
    "os"
    "testlake/middleware"
    "testlake/utils"
    "github.com/gin-gonic/gin"
)

func ServeApplication() {
    router := gin.Default()
    
    // Global middleware
    router.Use(middleware.DefaultAuthMiddleware())
    
    // Swagger setup
    Swagger(router)
    
    // API versioning
    baseRoute := router.Group("/api/v1")
    
    // Public routes (no authentication)
    publicRoutes := baseRoute.Group("")
    PublicRoutes(publicRoutes)
    
    // Private routes (JWT required)
    privateRoutes := baseRoute.Group("")
    privateRoutes.Use(middleware.JWTAuthMiddleware())
    PrivateRoutes(privateRoutes)
    
    // 404 handler
    router.NoRoute(utils.HandleNoRoute())
    
    // Start server
    ip := os.Getenv("IP")
    port := os.Getenv("PORT")
    router.Run(ip + ":" + port)
}
```

### 2. `controller/` - HTTP Request Handlers

**Purpose**: Handle HTTP requests, validate input, implement business logic, call DAO methods for data operations, and return responses. Controllers receive requests from the service layer and process them.

**Flow Position**: Client → Middleware → Service → **Controller** → DAO → Database

**Pattern**: Each controller focuses on a specific domain (User, Auth, etc.)

**Example (`controller/user_controller.go`)**:
```go
package controller

import (
    "net/http"
    "testlake/dao"
    "testlake/inout/user"
    "testlake/utils"
    "github.com/gin-gonic/gin"
)

type UserController struct{}

func (controller UserController) GetProfile(context *gin.Context) {
    // Extract user ID from JWT token
    userID, err := utils.ExtractUserID(context)
    if err != nil {
        utils.ReportUnauthorized(context, "Authentication required")
        return
    }

    // Query database
    userDao := dao.NewUserDao()
    foundUser, err := userDao.GetByID(userID)
    if err != nil {
        utils.ReportNotFound(context, "User not found")
        return
    }

    // Return response
    response := user.UserOut{
        BaseResponse: inout.BaseResponse{
            ErrorCode:        0,
            ErrorDescription: "Success",
        },
        Data: user.FromModel(foundUser),
    }

    context.JSON(http.StatusOK, response)
}
```

### 3. `dao/` - Data Access Objects

**Purpose**: Database operations and data persistence layer. DAOs receive calls from controllers and handle all database interactions using GORM.

**Flow Position**: Client → Middleware → Service → Controller → **DAO** → Database

**Key Features**:
- Database connection management
- GORM integration
- Automatic migrations
- CRUD operations

**Example (`dao/main.go`)**:
```go
package dao

import (
    "fmt"
    "log"
    "os"
    "testlake/model"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var Database *gorm.DB

func Connect() {
    // Database connection
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
        os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), 
        os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

    var err error
    Database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // Auto-migrate all models
    err = Database.AutoMigrate(&model.User{}, &model.Organization{})
    if err != nil {
        log.Fatal("Failed to run database migrations:", err)
    }
}
```

### 4. `docs/` - API Documentation

**Purpose**: Auto-generated Swagger/OpenAPI documentation.

**Generated by**: `swag init` command
**Files**: `docs.go`, `swagger.json`, `swagger.yaml`

### 5. `inout/` - Input/Output Structures

**Purpose**: Define request and response data structures separate from database models. This layer contains JSON tags and handles API serialization. Models convert themselves to these structures for external exposure.

**Key Principles**:
- Contains JSON tags for API serialization
- Separate from internal database models
- Input structures for request data validation
- Output structures for response formatting
- Conversion functions from models to inout structures

**Pattern**: Separate input (`*_in.go`) and output (`*_out.go`) structures for each domain.

**Example (`inout/base_response.go`)**:
```go
package inout

type BaseResponse struct {
    ErrorCode        int    `json:"error_code"`
    ErrorDescription string `json:"error_description"`
}

type PaginationMeta struct {
    Page       int   `json:"page"`
    Limit      int   `json:"limit"`
    Total      int64 `json:"total"`
    TotalPages int   `json:"total_pages"`
}
```

**Example (`inout/user/user_out.go`)**:
```go
package user

import (
    "time"
    "testlake/inout"
    "testlake/model"
    "github.com/google/uuid"
)

type UserData struct {
    ID              uuid.UUID `json:"id"`
    Email           string    `json:"email"`
    Username        string    `json:"username"`
    FirstName       *string   `json:"first_name"`
    LastName        *string   `json:"last_name"`
    IsEmailVerified bool      `json:"is_email_verified"`
    AuthProvider    string    `json:"auth_provider"`
    Status          string    `json:"status"`
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
}

type UserOut struct {
    inout.BaseResponse
    Data UserData `json:"data"`
}

// FromModel converts a model.User to UserData (no sensitive fields)
func FromModel(user *model.User) UserData {
    return UserData{
        ID:              user.ID,
        Email:           user.Email,
        Username:        user.Username,
        FirstName:       user.FirstName,
        LastName:        user.LastName,
        IsEmailVerified: user.IsEmailVerified,
        AuthProvider:    string(user.AuthProvider),
        Status:          string(user.Status),
        CreatedAt:       user.CreatedAt,
        UpdatedAt:       user.UpdatedAt,
        // Note: PasswordHash is intentionally excluded
    }
}
```

**Example (`inout/user/user_in.go`)**:
```go
package user

type UpdateUserRequest struct {
    FirstName *string `json:"first_name" validate:"omitempty,min=1,max=100"`
    LastName  *string `json:"last_name" validate:"omitempty,min=1,max=100"`
}

type CreateUserRequest struct {
    Email     string  `json:"email" validate:"required,email"`
    Username  string  `json:"username" validate:"required,min=3,max=50"`
    Password  string  `json:"password" validate:"required,min=8"`
    FirstName *string `json:"first_name" validate:"omitempty,min=1,max=100"`
    LastName  *string `json:"last_name" validate:"omitempty,min=1,max=100"`
}

// ToModel converts input data to a model.User for database operations
func (req *CreateUserRequest) ToModel() *model.User {
    return &model.User{
        Email:        req.Email,
        Username:     req.Username,
        FirstName:    req.FirstName,
        LastName:     req.LastName,
        AuthProvider: model.AuthProviderEmail,
        Status:       model.UserStatusActive,
        // PasswordHash will be set by the controller after hashing
    }
}
```

### 6. `middleware/` - HTTP Middleware

**Purpose**: Cross-cutting concerns like authentication, CORS, logging.

**Example (`middleware/jwt.go`)**:
```go
package middleware

import (
    "net/http"
    "testlake/inout"
    "testlake/utils"
    "github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
    return func(context *gin.Context) {
        err := utils.ValidateJWT(context)
        if err != nil {
            response := inout.BaseResponse{
                ErrorCode:        401,
                ErrorDescription: err.Error(),
            }
            context.JSON(http.StatusUnauthorized, response)
            context.Abort()
            return
        }
        context.Next()
    }
}
```

### 7. `model/` - Database Models

**Purpose**: GORM database models with relationships and constraints. Models are internal database entities and should NOT contain JSON tags. They convert themselves to inout structures for API responses.

**Key Principles**:
- No JSON tags in models (they're internal database entities)
- Models provide conversion methods to inout structures
- Only GORM tags for database mapping
- Sensitive fields (like PasswordHash) never exposed

**Example (`model/user.go`)**:
```go
package model

import (
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type AuthProvider string
type UserStatus string

const (
    AuthProviderEmail AuthProvider = "email"
    AuthProviderGmail AuthProvider = "gmail"
    UserStatusActive  UserStatus   = "active"
    UserStatusSuspended UserStatus = "suspended"
)

type User struct {
    ID              uuid.UUID      `gorm:"type:uuid;primaryKey"`
    Email           string         `gorm:"type:varchar(255);uniqueIndex;not null"`
    Username        string         `gorm:"type:varchar(100);uniqueIndex;not null"`
    FirstName       *string        `gorm:"type:varchar(100)"`
    LastName        *string        `gorm:"type:varchar(100)"`
    PasswordHash    string         `gorm:"type:varchar(255);not null"`
    IsEmailVerified bool           `gorm:"default:false"`
    AuthProvider    AuthProvider   `gorm:"type:varchar(50);default:'email'"`
    Status          UserStatus     `gorm:"type:varchar(50);default:'active'"`
    CreatedAt       time.Time
    UpdatedAt       time.Time
    DeletedAt       gorm.DeletedAt `gorm:"index"`
}

func (user *User) BeforeCreate(tx *gorm.DB) error {
    user.ID = uuid.New()
    return nil
}
```

### 8. `service/` - Service Layer (Route Registration)

**Purpose**: Route registration layer that connects HTTP routes to controllers. Handles Swagger documentation and HTTP method mapping. This is the entry point that receives requests from middleware and routes them to appropriate controllers.

**Flow Position**: Client → Middleware → **Service** → Controller → DAO → Database

**Example (`service/user_service.go`)**:
```go
package service

import (
    "testlake/controller"
    "github.com/gin-gonic/gin"
)

type UserService struct {
    Route      string
    Controller controller.UserController
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get current user's profile information
// @Tags User Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} user.UserOut
// @Failure 401 {object} inout.BaseResponse
// @Router /api/v1/users/profile [GET]
func (s UserService) GetProfile(r *gin.RouterGroup, route string) {
    r.GET("/"+s.Route+"/"+route, s.Controller.GetProfile)
}
```

### 9. `templates/` - HTML Templates

**Purpose**: HTML templates for emails, web pages, etc.

**Example**: Email verification templates, password reset forms.

### 10. `utils/` - Utility Functions

**Purpose**: Reusable utility functions for common operations.

**Key Files**:
- `error_handler.go`: Standardized error responses
- `jwt.go`: JWT token operations
- `password.go`: Password hashing and validation
- `email.go`: Email sending utilities

**Example (`utils/error_handler.go`)**:
```go
package utils

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func ReportNotFound(c *gin.Context, message string) {
    c.JSON(http.StatusNotFound, gin.H{
        "error_code":        404,
        "error_description": message,
    })
}

func ReportUnauthorized(c *gin.Context, message string) {
    c.JSON(http.StatusUnauthorized, gin.H{
        "error_code":        401,
        "error_description": message,
    })
}
```

### 11. `main.go` - Application Entry Point

**Purpose**: Application bootstrap and dependency injection.

```go
package main

import (
    "log"
    "testlake/app"
    "testlake/dao"
    _ "testlake/docs"
    "github.com/joho/godotenv"
)

func main() {
    // Load environment variables
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    // Connect to database
    dao.Connect()

    // Start web server
    app.ServeApplication()
}
```

## Environment Configuration

Create a `.env` file with the following variables:

```env
# Server Configuration
IP=0.0.0.0
PORT=8000

# Database Configuration
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=your_database
DB_PORT=5432

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key

# Email Configuration
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
```

## Request Flow Example

Here's how a typical request flows through the architecture:

```
1. Client sends HTTP request: GET /api/v1/users/profile
   ↓
2. Middleware processes request (CORS, Authentication)
   ↓
3. Service layer routes request to appropriate controller method
   ↓
4. Controller validates request, implements business logic
   ↓
5. Controller calls DAO method to fetch data
   ↓
6. DAO queries database using GORM
   ↓
7. Database returns data to DAO
   ↓
8. DAO returns model to Controller
   ↓
9. Controller transforms model to response format
   ↓
10. Controller returns JSON response to client
```

**Concrete Example**:
```go
// 1. Service registers route
func (s UserService) GetProfile(r *gin.RouterGroup, route string) {
    r.GET("/"+s.Route+"/"+route, s.Controller.GetProfile)  // Routes to controller
}

// 2. Controller handles request
func (controller UserController) GetProfile(context *gin.Context) {
    userID, _ := utils.ExtractUserID(context)              // Business logic
    
    userDao := dao.NewUserDao()                            // Create DAO
    foundUser, err := userDao.GetByID(userID)              // Call DAO method
    
    response := user.UserOut{                              // Transform response using inout
        BaseResponse: inout.BaseResponse{ErrorCode: 0, ErrorDescription: "Success"},
        Data: user.FromModel(foundUser),                   // Convert model to inout structure
    }
    context.JSON(http.StatusOK, response)                  // Return to client
}

// 3. DAO handles database operation
func (userDao *UserDao) GetByID(id uuid.UUID) (*model.User, error) {
    var user model.User
    err := userDao.db.First(&user, "id = ?", id).Error     // Database query
    return &user, err
}
```

## Development Workflow

### 1. Project Setup

```bash
# Initialize Go module
go mod init your-project-name

# Install dependencies
go mod tidy

# Generate Swagger docs
swag init

# Run the application
go run main.go
```

### 2. Adding New Features

1. **Create Model** (`model/new_feature.go`)
2. **Create DAO** (`dao/new_feature_dao.go`)
3. **Create Input/Output Structures** (`inout/new_feature/`)
4. **Create Controller** (`controller/new_feature_controller.go`)
5. **Create Service** (`service/new_feature_service.go`)
6. **Register Routes** (`app/routes.go`)
7. **Update Swagger** (`swag init`)

### 3. Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test ./controller_test -v
```

## Best Practices

### 1. Error Handling
- Use standardized error responses via `utils/error_handler.go`
- Always check and handle errors appropriately
- Use GORM error types for database-related errors

### 2. Security
- Always validate input data
- Use JWT for authentication
- Hash passwords with bcrypt
- Implement proper CORS middleware

### 3. Database
- Use GORM hooks for automatic field population (UUID, timestamps)
- Implement soft deletes where appropriate
- Use database migrations for schema changes

### 4. API Design
- Follow REST conventions
- Use appropriate HTTP status codes
- Implement proper pagination
- Version your APIs

### 5. Documentation
- Use Swagger annotations in service layer
- Document all public functions
- Maintain up-to-date README

## Docker Support

**Dockerfile**:
```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/.env .
CMD ["./main"]
```

**docker-compose.yml**:
```yaml
version: '3.8'
services:
  app:
    build: .
    ports:
      - "8000:8000"
    depends_on:
      - postgres
    environment:
      - DB_HOST=postgres
  
  postgres:
    image: postgres:13
    environment:
      POSTGRES_DB: your_database
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
    ports:
      - "5432:5432"
```

## Getting Started Checklist

- [ ] Clone or copy this architecture
- [ ] Update `go.mod` with your project name
- [ ] Configure `.env` file
- [ ] Update database models in `model/`
- [ ] Implement your business logic in controllers
- [ ] Register routes in `app/routes.go`
- [ ] Generate Swagger documentation with `swag init`
- [ ] Test your endpoints
- [ ] Deploy with Docker

This architecture provides a solid foundation for building scalable REST APIs in Go with clear separation of concerns and industry best practices.