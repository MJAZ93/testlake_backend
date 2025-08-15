# Go REST API - Generic Backend Documentation

## Table of Contents
- [Overview](#overview)
- [Architecture](#architecture)
- [Quick Start](#quick-start)
- [API Structure](#api-structure)
- [Authentication](#authentication)
- [Database Integration](#database-integration)
- [Error Handling Utilities](#error-handling-utilities)
- [Microservices Architecture](#microservices-architecture)
- [Testing](#testing)
- [Deployment](#deployment)
- [Examples](#examples)
- [Best Practices](#best-practices)

## Overview

This Go-based REST API follows a layered architecture pattern designed for building scalable web services. It uses modern Go frameworks and tools to provide a robust foundation for enterprise applications.

### Technology Stack
- **Language**: Go 1.24+
- **Web Framework**: Gin (github.com/gin-gonic/gin)
- **Documentation**: Swagger (github.com/swaggo/gin-swagger)
- **Database ORM**: GORM
- **Authentication**: JWT with refresh tokens
- **Environment**: godotenv

## Architecture

The system follows an MVC pattern with additional layers for better separation of concerns:

```
├── app/              # Application initialization and routing
├── service/          # Service layer with business logic
├── controller/       # HTTP handlers and request/response logic
├── dao/              # Data Access Objects (database layer)
├── model/            # Database entity definitions
├── inout/            # Request/Response DTOs
├── middleware/       # Gin middlewares
├── utils/            # Utility functions
├── service_test/     # Service layer tests
├── controller_test/  # Controller tests
└── dao_test/         # Database tests
```

## Quick Start

### Prerequisites
- Go 1.24 or higher
- PostgreSQL database
- Environment variables configured

### Installation
```bash
# Clone the repository
git clone <repository-url>
cd generic_backend

# Install dependencies
go mod tidy

# Set up environment variables
cp .env.example .env
# Edit .env with your configuration

# Run the application
go run main.go
```

### Environment Configuration (.env)
```env
# Server Configuration
IP=localhost
PORT=8000
SCHEME=http

# Database Configuration
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=your_database
DB_PORT=5432

# JWT Configuration
TOKEN_TTL=2000
JWT_PRIVATE_KEY=your_secret_key

# Logging
LOG_PATH=/path/to/logs
```

## API Structure

### Base Application Setup
```go
// app/app.go
// @BasePath /api/
func ServeApplication() {
    router := gin.Default()
    
    // Add Swagger documentation
    Swagger(router)
    
    // Apply default middleware
    router.Use(middleware.DefaultAuthMiddleware())
    
    // Define route groups
    baseRoute := router.Group("/api/v1")
    
    // Public routes (no authentication required)
    publicRoutes := baseRoute.Group("/public")
    PublicRoutes(publicRoutes)
    
    // Private routes (JWT authentication required)
    privateRoutes := baseRoute.Group("/private")
    privateRoutes.Use(middleware.JWTAuthMiddleware())
    PrivateRoutes(privateRoutes)
    
    // Start server
    ip := os.Getenv("IP")
    port := os.Getenv("PORT")
    router.Run(ip + ":" + port)
}
```

### Route Registration
```go
// app/routes.go
func PublicRoutes(r *gin.RouterGroup) {
    // Example service registration
    clientService := service.ClientService{
        Route: "client", 
        Controller: controller.ClientController{},
    }
    clientService.GetDetails(r, "details")
}

func PrivateRoutes(r *gin.RouterGroup) {
    // Protected routes
    userService := service.UserService{
        Route: "user", 
        Controller: controller.UserController{},
    }
    userService.List(r, "list")
}
```

## Authentication

### JWT Middleware
```go
// middleware/jwt.go
func JWTAuthMiddleware() gin.HandlerFunc {
    return func(context *gin.Context) {
        err := util.ValidateJWT(context)
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

### CORS Middleware
```go
// middleware/cors.go
func DefaultAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Content-Type", "application/json")
        c.Header("Access-Control-Allow-Origin", "*")
        
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
        
        if c.Request.Method == http.MethodOptions {
            c.AbortWithStatus(http.StatusOK)
            return
        }
        
        c.Next()
    }
}
```

## Database Integration

### Database Connection
```go
// dao/main.go
var Database *gorm.DB

func Connect() {
    host := os.Getenv("DB_HOST")
    username := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASSWORD")
    databaseName := os.Getenv("DB_NAME")
    port := os.Getenv("DB_PORT")

    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", 
        host, username, password, databaseName, port)
    
    Database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    
    if err != nil {
        panic(err)
    }
}
```

### Model Definition
```go
// model/user.go
type User struct {
    ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    Name      string    `gorm:"not null" json:"name"`
    Email     string    `gorm:"unique;not null" json:"email"`
    Password  string    `gorm:"not null" json:"-"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### DAO Implementation
```go
// dao/user_dao.go
type UserDao struct {
    Limit int
}

func (dao *UserDao) Create(user *model.User) error {
    return Database.Create(&user).Error
}

func (dao *UserDao) GetByID(id uint) (error, model.User) {
    var user model.User
    err := Database.First(&user, id).Error
    return err, user
}

func (dao *UserDao) GetAll(page int) (error, []model.User) {
    var users []model.User
    query := Database.Offset(page * dao.Limit).Limit(dao.Limit)
    err := query.Find(&users).Error
    return err, users
}
```

## Service Layer

### Service Structure
```go
// service/user_service.go
type UserService struct {
    Route      string
    Controller controller.UserController
}

// GetUser godoc
// @Summary Get user by ID
// @Tags User
// @Description Get user details by ID
// @Accept json
// @Produce json
// @Success 200 {object} user.UserOut
// @Param id path int true "User ID"
// @Router /private/user/details/{id} [GET]
func (s UserService) GetUser(r *gin.RouterGroup, route string) {
    r.GET("/"+s.Route+"/"+route+"/:id", s.Controller.GetUser)
}

// ListUsers godoc
// @Summary List all users
// @Tags User
// @Description Get paginated list of users
// @Accept json
// @Produce json
// @Success 200 {object} user.UserListOut
// @Param page query int false "Page number"
// @Router /private/user/list [GET]
func (s UserService) ListUsers(r *gin.RouterGroup, route string) {
    r.GET("/"+s.Route+"/"+route, s.Controller.ListUsers)
}
```

## Controller Layer

### Controller Implementation
```go
// controller/user_controller.go
type UserController struct{}

func (controller UserController) GetUser(context *gin.Context) {
    idParam := context.Param("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{
            "error_code": 400,
            "error_description": "Invalid user ID",
        })
        return
    }
    
    dao := dao2.UserDao{Limit: 1}
    err, user := dao.GetByID(uint(id))
    if err != nil {
        context.JSON(http.StatusNotFound, gin.H{
            "error_code": 404,
            "error_description": "User not found",
        })
        return
    }
    
    response := user.UserOut{
        BaseResponse: inout.BaseResponse{
            ErrorCode: 0,
            ErrorDescription: "Success",
        },
        Data: user,
    }
    
    context.JSON(http.StatusOK, response)
}
```

## Input/Output DTOs

### Base Response Structure
```go
// inout/base_response.go
type BaseResponse struct {
    ErrorCode        int    `json:"error_code"`
    ErrorDescription string `json:"error_description"`
}
```

### Specific Response DTOs
```go
// inout/user/user_out.go
type UserOut struct {
    inout.BaseResponse
    Data User `json:"data"`
}

type UserListOut struct {
    inout.BaseResponse
    List []User `json:"list"`
    Meta Pagination `json:"meta"`
}

type User struct {
    ID        uint      `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type Pagination struct {
    Page       int `json:"page"`
    Limit      int `json:"limit"`
    Total      int `json:"total"`
    TotalPages int `json:"total_pages"`
}
```

## Error Handling Utilities

The application provides a centralized error handling utility in `utils/error_handler.go` that offers a generic and consistent approach to handling HTTP errors, including 404 Not Found errors and other common HTTP status codes.

### ErrorHandler Structure

```go
// utils/error_handler.go
type StandardError struct {
    ErrorCode        int    `json:"error_code"`
    ErrorDescription string `json:"error_description"`
}

type ErrorHandler struct{}

func NewErrorHandler() *ErrorHandler {
    return &ErrorHandler{}
}
```

### Generic 404 Error Handling

The error handler provides a flexible generic approach for handling 404 errors with customizable messages:

```go
// NotFound handles 404 Not Found errors with optional custom message
func (eh *ErrorHandler) NotFound(c *gin.Context, message ...string) {
    defaultMessage := "RESOURCE_NOT_FOUND"
    if len(message) > 0 {
        defaultMessage = message[0]
    }
    
    c.JSON(http.StatusNotFound, StandardError{
        ErrorCode:        404,
        ErrorDescription: defaultMessage,
    })
}
```

**Key Features of Generic 404 Handling:**
- **Variadic Parameters**: Accepts zero or more custom messages using `message ...string`
- **Default Fallback**: Uses "RESOURCE_NOT_FOUND" when no custom message is provided
- **Consistent Response**: Always returns standardized JSON structure with error code and description
- **Type Safety**: Leverages Go's type system for reliable error responses

### Complete Error Handler Methods

The error handler supports all common HTTP error status codes:

```go
// 400 Bad Request
func (eh *ErrorHandler) BadRequest(c *gin.Context, message ...string)

// 401 Unauthorized
func (eh *ErrorHandler) Unauthorized(c *gin.Context, message ...string)

// 403 Forbidden
func (eh *ErrorHandler) Forbidden(c *gin.Context, message ...string)

// 404 Not Found
func (eh *ErrorHandler) NotFound(c *gin.Context, message ...string)

// 500 Internal Server Error
func (eh *ErrorHandler) InternalServerError(c *gin.Context, message ...string)

// Validation errors with field-specific information
func (eh *ErrorHandler) ValidationError(c *gin.Context, field string, message string)

// Custom errors with specific status and error codes
func (eh *ErrorHandler) CustomError(c *gin.Context, statusCode int, errorCode int, message string)
```

### Convenience Functions

Global convenience functions are provided for easier usage throughout the application:

```go
// Global convenience functions
func ReportNotFound(c *gin.Context, message ...string)
func ReportBadRequest(c *gin.Context, message ...string)
func ReportUnauthorized(c *gin.Context, message ...string)
func ReportForbidden(c *gin.Context, message ...string)
func ReportInternalServerError(c *gin.Context, message ...string)
func ReportValidationError(c *gin.Context, field string, message string)
func ReportCustomError(c *gin.Context, statusCode int, errorCode int, message string)
```

### Usage Examples

#### Basic 404 Error Handling
```go
// controller/user_controller.go
func (controller UserController) GetUser(context *gin.Context) {
    idParam := context.Param("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        utils.ReportBadRequest(context, "Invalid user ID format")
        return
    }
    
    dao := dao2.UserDao{Limit: 1}
    err, user := dao.GetByID(uint(id))
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            // Generic 404 handling with custom message
            utils.ReportNotFound(context, "User not found")
        } else {
            utils.ReportInternalServerError(context, "Database error occurred")
        }
        return
    }
    
    context.JSON(http.StatusOK, gin.H{
        "error_code": 0,
        "error_description": "Success",
        "data": user,
    })
}
```

#### Multiple Error Scenarios
```go
// controller/product_controller.go
func (controller ProductController) GetProduct(context *gin.Context) {
    idParam := context.Param("id")
    if idParam == "" {
        utils.ReportBadRequest(context, "Product ID is required")
        return
    }
    
    id, err := strconv.Atoi(idParam)
    if err != nil {
        utils.ReportBadRequest(context, "Invalid product ID format")
        return
    }
    
    if id <= 0 {
        utils.ReportBadRequest(context, "Product ID must be positive")
        return
    }
    
    dao := dao2.ProductDao{}
    err, product := dao.GetByID(uint(id))
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            // Using generic 404 handling with specific context
            utils.ReportNotFound(context, "Product with ID "+idParam+" not found")
        } else {
            utils.ReportInternalServerError(context)
        }
        return
    }
    
    context.JSON(http.StatusOK, gin.H{
        "error_code": 0,
        "error_description": "Success",
        "data": product,
    })
}
```

#### Route Handler for Undefined Endpoints
```go
// app/app.go
func ServeApplication() {
    router := gin.Default()
    
    // ... other route configurations ...
    
    // Handle 404 for undefined routes
    router.NoRoute(utils.HandleNoRoute())
    
    // Start server
    ip := os.Getenv("IP")
    port := os.Getenv("PORT")
    router.Run(ip + ":" + port)
}
```

### Error Response Structure

All errors follow a consistent JSON response structure:

```json
{
    "error_code": 404,
    "error_description": "User not found"
}
```

For validation errors, additional fields are included:
```json
{
    "error_code": 400,
    "error_description": "VALIDATION_ERROR",
    "field": "email",
    "message": "Email format is invalid"
}
```

### Benefits of Generic Error Handling

1. **Consistency**: All errors follow the same response format across the API
2. **Flexibility**: Custom messages can be provided while maintaining structure
3. **Type Safety**: Go's type system ensures reliable error handling
4. **Maintainability**: Centralized error handling reduces code duplication
5. **Debugging**: Standardized error codes make debugging easier
6. **API Documentation**: Consistent error responses simplify API documentation
7. **Client Integration**: Predictable error format improves client-side error handling

## Microservices Architecture

This section covers how to decompose the monolithic Go REST API into a microservices architecture where each service is containerized using Docker and can be deployed independently.

### Overview

The microservices architecture transforms the layered monolithic application into independent, loosely coupled services that communicate over well-defined APIs. Each service:

- Runs in its own Docker container
- Has its own database (if needed)
- Can be developed, tested, and deployed independently
- Communicates via HTTP REST APIs or message queues
- Implements a single business capability

### Microservices Decomposition Strategy

#### Service Boundaries

Based on the existing architecture, we can decompose into the following services:

```
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│   API Gateway   │  │  Auth Service   │  │  User Service   │
│   (Port 8000)   │  │   (Port 8001)   │  │   (Port 8002)   │
└─────────────────┘  └─────────────────┘  └─────────────────┘
         │                      │                      │
         ▼                      ▼                      ▼
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│ Product Service │  │ Order Service   │  │ Notification    │
│   (Port 8003)   │  │   (Port 8004)   │  │   Service       │
└─────────────────┘  └─────────────────┘  │   (Port 8005)   │
                                           └─────────────────┘
```

#### Service Responsibilities

- **API Gateway**: Request routing, load balancing, rate limiting, authentication
- **Auth Service**: JWT token management, user authentication, authorization
- **User Service**: User management, profiles, user-related operations
- **Product Service**: Product catalog, inventory management
- **Order Service**: Order processing, order history
- **Notification Service**: Email, SMS, push notifications

### Docker Containerization

#### Individual Service Dockerfile

Each service gets its own optimized Dockerfile:

```dockerfile
# services/user-service/Dockerfile
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY ../go.mod go.sum ./
RUN go mod download

# Copy source code
COPY .. .

# Build the service
RUN CGO_ENABLED=0 GOOS=linux go build -o user-service ./cmd/user-service/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy binary and config
COPY --from=builder /app/user-service .
COPY --from=builder /app/.env .

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8002/health || exit 1

EXPOSE 8002
CMD ["./user-service"]
```

#### Multi-Service Docker Compose

Complete Docker Compose configuration for all microservices:

```yaml
# docker-compose.yml
version: '3.8'

services:
  # Infrastructure Services
  postgres:
    image: postgres:13
    environment:
      POSTGRES_DB: microservices_db
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - microservices-network

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    networks:
      - microservices-network

  # API Gateway
  api-gateway:
    build:
      context: ./services/api-gateway
      dockerfile: Dockerfile
    ports:
      - "8000:8000"
    environment:
      - PORT=8000
      - AUTH_SERVICE_URL=http://auth-service:8001
      - USER_SERVICE_URL=http://user-service:8002
      - PRODUCT_SERVICE_URL=http://product-service:8003
      - ORDER_SERVICE_URL=http://order-service:8004
      - NOTIFICATION_SERVICE_URL=http://notification-service:8005
    depends_on:
      - auth-service
      - user-service
      - product-service
      - order-service
      - notification-service
    networks:
      - microservices-network

  # Auth Service
  auth-service:
    build:
      context: ./services/auth-service
      dockerfile: Dockerfile
    ports:
      - "8001:8001"
    environment:
      - PORT=8001
      - DB_HOST=postgres
      - DB_USER=postgres
      - DB_PASSWORD=password
      - DB_NAME=microservices_db
      - DB_PORT=5432
      - JWT_PRIVATE_KEY=${JWT_PRIVATE_KEY}
      - REDIS_HOST=redis
      - REDIS_PORT=6379
    depends_on:
      - postgres
      - redis
    networks:
      - microservices-network
    healthcheck:
      test: ["CMD", "wget", "--spider", "-q", "http://localhost:8001/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  # User Service
  user-service:
    build:
      context: ./services/user-service
      dockerfile: Dockerfile
    ports:
      - "8002:8002"
    environment:
      - PORT=8002
      - DB_HOST=postgres
      - DB_USER=postgres
      - DB_PASSWORD=password
      - DB_NAME=microservices_db
      - DB_PORT=5432
      - AUTH_SERVICE_URL=http://auth-service:8001
    depends_on:
      - postgres
      - auth-service
    networks:
      - microservices-network
    healthcheck:
      test: ["CMD", "wget", "--spider", "-q", "http://localhost:8002/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  # Product Service
  product-service:
    build:
      context: ./services/product-service
      dockerfile: Dockerfile
    ports:
      - "8003:8003"
    environment:
      - PORT=8003
      - DB_HOST=postgres
      - DB_USER=postgres
      - DB_PASSWORD=password
      - DB_NAME=microservices_db
      - DB_PORT=5432
    depends_on:
      - postgres
    networks:
      - microservices-network
    healthcheck:
      test: ["CMD", "wget", "--spider", "-q", "http://localhost:8003/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  # Order Service
  order-service:
    build:
      context: ./services/order-service
      dockerfile: Dockerfile
    ports:
      - "8004:8004"
    environment:
      - PORT=8004
      - DB_HOST=postgres
      - DB_USER=postgres
      - DB_PASSWORD=password
      - DB_NAME=microservices_db
      - DB_PORT=5432
      - USER_SERVICE_URL=http://user-service:8002
      - PRODUCT_SERVICE_URL=http://product-service:8003
      - NOTIFICATION_SERVICE_URL=http://notification-service:8005
    depends_on:
      - postgres
      - user-service
      - product-service
      - notification-service
    networks:
      - microservices-network
    healthcheck:
      test: ["CMD", "wget", "--spider", "-q", "http://localhost:8004/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  # Notification Service
  notification-service:
    build:
      context: ./services/notification-service
      dockerfile: Dockerfile
    ports:
      - "8005:8005"
    environment:
      - PORT=8005
      - REDIS_HOST=redis
      - REDIS_PORT=6379
      - EMAIL_SERVICE_URL=${EMAIL_SERVICE_URL}
      - SMS_SERVICE_URL=${SMS_SERVICE_URL}
    depends_on:
      - redis
    networks:
      - microservices-network
    healthcheck:
      test: ["CMD", "wget", "--spider", "-q", "http://localhost:8005/health"]
      interval: 30s
      timeout: 10s
      retries: 3

volumes:
  postgres_data:

networks:
  microservices-network:
    driver: bridge
```

### API Gateway Implementation

The API Gateway acts as a single entry point and handles routing, authentication, and cross-cutting concerns:

```go
// services/api-gateway/main.go
package main

import (
    "net/http"
    "net/http/httputil"
    "net/url"
    "os"
    "strings"

    "github.com/gin-gonic/gin"
)

type Gateway struct {
    authService         *httputil.ReverseProxy
    userService         *httputil.ReverseProxy
    productService      *httputil.ReverseProxy
    orderService        *httputil.ReverseProxy
    notificationService *httputil.ReverseProxy
}

func NewGateway() *Gateway {
    return &Gateway{
        authService:         createReverseProxy(os.Getenv("AUTH_SERVICE_URL")),
        userService:         createReverseProxy(os.Getenv("USER_SERVICE_URL")),
        productService:      createReverseProxy(os.Getenv("PRODUCT_SERVICE_URL")),
        orderService:        createReverseProxy(os.Getenv("ORDER_SERVICE_URL")),
        notificationService: createReverseProxy(os.Getenv("NOTIFICATION_SERVICE_URL")),
    }
}

func createReverseProxy(target string) *httputil.ReverseProxy {
    url, _ := url.Parse(target)
    return httputil.NewSingleHostReverseProxy(url)
}

func (g *Gateway) setupRoutes() *gin.Engine {
    router := gin.Default()

    // CORS middleware
    router.Use(func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(http.StatusOK)
            return
        }
        
        c.Next()
    })

    // Health check
    router.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "api-gateway"})
    })

    // API versioning
    v1 := router.Group("/api/v1")

    // Auth routes (no authentication required)
    auth := v1.Group("/auth")
    auth.Any("/*path", g.proxyToService(g.authService, "/api/v1/auth"))

    // Protected routes (require authentication)
    protected := v1.Group("")
    protected.Use(g.authMiddleware())

    // User service routes
    protected.Any("/users/*path", g.proxyToService(g.userService, "/api/v1/users"))

    // Product service routes
    protected.Any("/products/*path", g.proxyToService(g.productService, "/api/v1/products"))

    // Order service routes
    protected.Any("/orders/*path", g.proxyToService(g.orderService, "/api/v1/orders"))

    // Notification service routes
    protected.Any("/notifications/*path", g.proxyToService(g.notificationService, "/api/v1/notifications"))

    return router
}

func (g *Gateway) proxyToService(proxy *httputil.ReverseProxy, stripPrefix string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Strip the prefix from the request path
        c.Request.URL.Path = strings.TrimPrefix(c.Request.URL.Path, stripPrefix)
        
        // Add service-specific headers
        c.Request.Header.Set("X-Gateway", "true")
        c.Request.Header.Set("X-Original-Path", c.Request.URL.Path)
        
        proxy.ServeHTTP(c.Writer, c.Request)
    }
}

func (g *Gateway) authMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Extract JWT token
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error_code": 401,
                "error_description": "Authorization header required",
            })
            c.Abort()
            return
        }

        // Validate token with auth service
        if !g.validateTokenWithAuthService(token) {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error_code": 401,
                "error_description": "Invalid token",
            })
            c.Abort()
            return
        }

        c.Next()
    }
}

func (g *Gateway) validateTokenWithAuthService(token string) bool {
    // Implementation to validate token with auth service
    // This would make a call to the auth service to validate the token
    return true // Simplified for example
}

func main() {
    gateway := NewGateway()
    router := gateway.setupRoutes()
    
    port := os.Getenv("PORT")
    if port == "" {
        port = "8000"
    }
    
    router.Run(":" + port)
}
```

### Service Communication Patterns

#### HTTP Client for Inter-Service Communication

```go
// pkg/httpclient/client.go
package httpclient

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type ServiceClient struct {
    BaseURL    string
    HTTPClient *http.Client
    AuthToken  string
}

func NewServiceClient(baseURL string) *ServiceClient {
    return &ServiceClient{
        BaseURL: baseURL,
        HTTPClient: &http.Client{
            Timeout: 10 * time.Second,
        },
    }
}

func (c *ServiceClient) SetAuthToken(token string) {
    c.AuthToken = token
}

func (c *ServiceClient) Get(endpoint string, result interface{}) error {
    req, err := http.NewRequest("GET", c.BaseURL+endpoint, nil)
    if err != nil {
        return err
    }

    if c.AuthToken != "" {
        req.Header.Set("Authorization", c.AuthToken)
    }
    req.Header.Set("Content-Type", "application/json")

    resp, err := c.HTTPClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
    }

    return json.NewDecoder(resp.Body).Decode(result)
}

func (c *ServiceClient) Post(endpoint string, payload interface{}, result interface{}) error {
    jsonData, err := json.Marshal(payload)
    if err != nil {
        return err
    }

    req, err := http.NewRequest("POST", c.BaseURL+endpoint, bytes.NewBuffer(jsonData))
    if err != nil {
        return err
    }

    if c.AuthToken != "" {
        req.Header.Set("Authorization", c.AuthToken)
    }
    req.Header.Set("Content-Type", "application/json")

    resp, err := c.HTTPClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
        return fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
    }

    if result != nil {
        return json.NewDecoder(resp.Body).Decode(result)
    }

    return nil
}
```

### Service Discovery with Consul

#### Consul Integration

```go
// pkg/discovery/consul.go
package discovery

import (
    "fmt"
    "strconv"

    "github.com/hashicorp/consul/api"
)

type ConsulRegistry struct {
    client *api.Client
}

func NewConsulRegistry(consulURL string) (*ConsulRegistry, error) {
    config := api.DefaultConfig()
    config.Address = consulURL
    
    client, err := api.NewClient(config)
    if err != nil {
        return nil, err
    }
    
    return &ConsulRegistry{client: client}, nil
}

func (r *ConsulRegistry) RegisterService(name, host string, port int, tags []string) error {
    service := &api.AgentServiceRegistration{
        ID:      fmt.Sprintf("%s-%s-%d", name, host, port),
        Name:    name,
        Tags:    tags,
        Port:    port,
        Address: host,
        Check: &api.AgentServiceCheck{
            HTTP:                           fmt.Sprintf("http://%s:%d/health", host, port),
            Timeout:                        "10s",
            Interval:                       "30s",
            DeregisterCriticalServiceAfter: "60s",
        },
    }

    return r.client.Agent().ServiceRegister(service)
}

func (r *ConsulRegistry) DiscoverService(serviceName string) (string, error) {
    services, _, err := r.client.Health().Service(serviceName, "", true, nil)
    if err != nil {
        return "", err
    }

    if len(services) == 0 {
        return "", fmt.Errorf("service %s not found", serviceName)
    }

    service := services[0].Service
    return fmt.Sprintf("http://%s:%d", service.Address, service.Port), nil
}
```

### Individual Service Structure

#### User Service Example

```go
// services/user-service/main.go
package main

import (
    "os"

    "github.com/gin-gonic/gin"
    "your-app/services/user-service/controller"
    "your-app/services/user-service/dao"
    "your-app/services/user-service/service"
    "your-app/pkg/discovery"
    "your-app/pkg/httpclient"
)

func main() {
    // Connect to database
    dao.Connect()

    // Initialize service registry
    consul, err := discovery.NewConsulRegistry("localhost:8500")
    if err != nil {
        panic("Failed to connect to Consul: " + err.Error())
    }

    // Register service
    err = consul.RegisterService("user-service", "localhost", 8002, []string{"api", "user"})
    if err != nil {
        panic("Failed to register service: " + err.Error())
    }

    // Initialize HTTP clients for other services
    authClient := httpclient.NewServiceClient("http://auth-service:8001")

    // Initialize controllers with dependencies
    userController := controller.NewUserController(authClient)

    // Setup routes
    router := gin.Default()

    // Health check
    router.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "healthy", "service": "user-service"})
    })

    // API routes
    api := router.Group("/api/v1")
    userService := service.UserService{
        Route:      "users",
        Controller: userController,
    }
    
    userService.RegisterRoutes(api)

    // Start server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8002"
    }
    
    router.Run(":" + port)
}
```

### Monitoring and Logging

#### Centralized Logging with ELK Stack

```yaml
# docker-compose.monitoring.yml
version: '3.8'

services:
  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:7.14.0
    environment:
      - discovery.type=single-node
      - "ES_JAVA_OPTS=-Xms512m -Xmx512m"
    ports:
      - "9200:9200"
    networks:
      - microservices-network

  logstash:
    image: docker.elastic.co/logstash/logstash:7.14.0
    volumes:
      - ./logstash/config/logstash.yml:/usr/share/logstash/config/logstash.yml:ro
      - ./logstash/pipeline:/usr/share/logstash/pipeline:ro
    ports:
      - "5044:5044"
      - "5000:5000/tcp"
      - "5000:5000/udp"
      - "9600:9600"
    environment:
      LS_JAVA_OPTS: "-Xmx256m -Xms256m"
    networks:
      - microservices-network
    depends_on:
      - elasticsearch

  kibana:
    image: docker.elastic.co/kibana/kibana:7.14.0
    ports:
      - "5601:5601"
    environment:
      ELASTICSEARCH_URL: http://elasticsearch:9200
      ELASTICSEARCH_HOSTS: '["http://elasticsearch:9200"]'
    networks:
      - microservices-network
    depends_on:
      - elasticsearch

  prometheus:
    image: prom/prometheus:latest
    ports:
      - "9090:9090"
    volumes:
      - ./monitoring/prometheus.yml:/etc/prometheus/prometheus.yml
    networks:
      - microservices-network

  grafana:
    image: grafana/grafana:latest
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
    volumes:
      - grafana-storage:/var/lib/grafana
    networks:
      - microservices-network
    depends_on:
      - prometheus

volumes:
  grafana-storage:

networks:
  microservices-network:
    external: true
```

#### Structured Logging in Services

```go
// pkg/logger/logger.go
package logger

import (
    "os"

    "github.com/sirupsen/logrus"
)

type ServiceLogger struct {
    *logrus.Logger
    ServiceName string
}

func NewServiceLogger(serviceName string) *ServiceLogger {
    logger := logrus.New()
    
    // Set JSON formatter for structured logging
    logger.SetFormatter(&logrus.JSONFormatter{
        FieldMap: logrus.FieldMap{
            logrus.FieldKeyTime:  "@timestamp",
            logrus.FieldKeyLevel: "level",
            logrus.FieldKeyMsg:   "message",
        },
    })
    
    // Set output to stdout for container logging
    logger.SetOutput(os.Stdout)
    
    // Set log level from environment
    level, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
    if err != nil {
        level = logrus.InfoLevel
    }
    logger.SetLevel(level)
    
    return &ServiceLogger{
        Logger:      logger,
        ServiceName: serviceName,
    }
}

func (sl *ServiceLogger) WithRequestID(requestID string) *logrus.Entry {
    return sl.WithFields(logrus.Fields{
        "service":    sl.ServiceName,
        "request_id": requestID,
    })
}

func (sl *ServiceLogger) WithUserID(userID string) *logrus.Entry {
    return sl.WithFields(logrus.Fields{
        "service": sl.ServiceName,
        "user_id": userID,
    })
}
```

### Deployment Commands

#### Development Environment

```bash
# Start all services
docker-compose up --build -d

# Start with monitoring
docker-compose -f docker-compose.yml -f docker-compose.monitoring.yml up --build -d

# Scale specific services
docker-compose up --scale user-service=3 --scale product-service=2

# View logs for specific service
docker-compose logs -f user-service

# Stop all services
docker-compose down

# Clean up volumes
docker-compose down -v
```

#### Production Deployment with Kubernetes

```yaml
# k8s/user-service-deployment.yml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: user-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: user-service
  template:
    metadata:
      labels:
        app: user-service
    spec:
      containers:
      - name: user-service
        image: your-registry/user-service:latest
        ports:
        - containerPort: 8002
        env:
        - name: PORT
          value: "8002"
        - name: DB_HOST
          value: "postgres-service"
        - name: DB_USER
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: username
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: password
        livenessProbe:
          httpGet:
            path: /health
            port: 8002
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8002
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: user-service
spec:
  selector:
    app: user-service
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8002
  type: ClusterIP
```

### Best Practices for Microservices

#### 1. Database Per Service
- Each service owns its data
- No direct database access between services
- Use API calls for data exchange
- Implement eventual consistency patterns

#### 2. Service Communication
- Use async messaging for non-critical operations
- Implement circuit breakers for resilience
- Add request/response timeouts
- Use service mesh for complex communication patterns

#### 3. Configuration Management
- Use environment variables for service configuration
- Implement configuration hot-reloading
- Use secrets management for sensitive data
- Maintain separate configs for different environments

#### 4. Testing Strategy
- Unit tests for individual services
- Integration tests between services
- Contract testing for API compatibility
- End-to-end tests for critical user journeys

#### 5. Monitoring and Observability
- Implement distributed tracing
- Use centralized logging
- Monitor service health and performance
- Set up alerting for critical metrics

This microservices architecture provides scalability, maintainability, and deployment flexibility while maintaining the existing API structure and functionality.

## Testing

### Service Tests
```go
// service_test/user_test.go
func TestUserService(t *testing.T) {
    // Load environment variables
    err := godotenv.Load("../.env")
    if err != nil {
        t.Fatal("Error loading .env file")
    }
    
    // Test user creation
    err, user := CreateTestUser()
    if err != nil {
        t.Fatalf("Failed to create user: %v", err)
    }
    
    fmt.Println("User created successfully:", user)
}

func CreateTestUser() (error, user.UserOut) {
    ip := os.Getenv("IP")
    port := os.Getenv("PORT")
    scheme := os.Getenv("SCHEME")
    
    url := scheme + "://" + ip + ":" + port + "/api/v1/private/user/create"
    
    userData := map[string]interface{}{
        "name":  "Test User",
        "email": "test@example.com",
    }
    
    jsonData, _ := json.Marshal(userData)
    
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return err, user.UserOut{}
    }
    
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+getTestToken())
    
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err, user.UserOut{}
    }
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return err, user.UserOut{}
    }
    
    var result user.UserOut
    if err := json.Unmarshal(body, &result); err != nil {
        return err, user.UserOut{}
    }
    
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("Expected status 200, got %v", resp.StatusCode), user.UserOut{}
    }
    
    return nil, result
}
```

### Controller Tests
```go
// controller_test/user_controller_test.go
func TestUserController(t *testing.T) {
    // Initialize test database
    dao.Connect()
    
    // Create test router
    router := gin.Default()
    controller := controller.UserController{}
    
    router.GET("/user/:id", controller.GetUser)
    
    // Test valid user ID
    req, _ := http.NewRequest("GET", "/user/1", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    assert.Equal(t, 200, w.Code)
    
    // Test invalid user ID
    req, _ = http.NewRequest("GET", "/user/invalid", nil)
    w = httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    assert.Equal(t, 400, w.Code)
}
```

## Swagger Documentation

### Main Setup
```go
// app/routes.go
func Swagger(r *gin.Engine) {
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
```

### Swagger Annotations
```go
// @Summary Create new user
// @Description Create a new user account
// @Tags User
// @Accept json
// @Produce json
// @Param user body user.CreateUserRequest true "User data"
// @Success 201 {object} user.UserOut
// @Failure 400 {object} inout.BaseResponse
// @Router /private/user/create [POST]
func (s UserService) CreateUser(r *gin.RouterGroup, route string) {
    r.POST("/"+s.Route+"/"+route, s.Controller.CreateUser)
}
```

## Deployment

### Docker Configuration
```dockerfile
# Dockerfile
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/.env .

EXPOSE 8000
CMD ["./main"]
```

### Docker Compose
```yaml
# docker-compose.yml
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
      - DB_USER=postgres
      - DB_PASSWORD=password
      - DB_NAME=myapp
      - DB_PORT=5432

  postgres:
    image: postgres:13
    environment:
      - POSTGRES_DB=myapp
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

## Examples

### Complete CRUD Example

#### 1. Model
```go
// model/product.go
type Product struct {
    ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    Name        string    `gorm:"not null" json:"name"`
    Description string    `json:"description"`
    Price       float64   `gorm:"not null" json:"price"`
    Stock       int       `gorm:"default:0" json:"stock"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

#### 2. DAO
```go
// dao/product_dao.go
type ProductDao struct {
    Limit int
}

func (dao *ProductDao) Create(product *model.Product) error {
    return Database.Create(&product).Error
}

func (dao *ProductDao) GetByID(id uint) (error, model.Product) {
    var product model.Product
    err := Database.First(&product, id).Error
    return err, product
}

func (dao *ProductDao) Update(product *model.Product) error {
    return Database.Save(&product).Error
}

func (dao *ProductDao) Delete(id uint) error {
    return Database.Delete(&model.Product{}, id).Error
}

func (dao *ProductDao) GetAll(page int) (error, []model.Product) {
    var products []model.Product
    query := Database.Offset(page * dao.Limit).Limit(dao.Limit)
    err := query.Find(&products).Error
    return err, products
}
```

#### 3. Controller
```go
// controller/product_controller.go
type ProductController struct{}

func (controller ProductController) Create(context *gin.Context) {
    var product model.Product
    if err := context.ShouldBindJSON(&product); err != nil {
        context.JSON(http.StatusBadRequest, gin.H{
            "error_code": 400,
            "error_description": err.Error(),
        })
        return
    }
    
    dao := dao2.ProductDao{Limit: 1}
    if err := dao.Create(&product); err != nil {
        context.JSON(http.StatusInternalServerError, gin.H{
            "error_code": 500,
            "error_description": "Failed to create product",
        })
        return
    }
    
    context.JSON(http.StatusCreated, gin.H{
        "error_code": 0,
        "error_description": "Success",
        "data": product,
    })
}
```

#### 4. Service
```go
// service/product_service.go
type ProductService struct {
    Route      string
    Controller controller.ProductController
}

// CreateProduct godoc
// @Summary Create product
// @Tags Product
// @Accept json
// @Produce json
// @Param product body model.Product true "Product data"
// @Success 201 {object} product.ProductOut
// @Router /private/product/create [POST]
func (s ProductService) Create(r *gin.RouterGroup, route string) {
    r.POST("/"+s.Route+"/"+route, s.Controller.Create)
}
```

### Error Handling Example
```go
// util/error_handler.go
func HandleError(context *gin.Context, err error, statusCode int) {
    response := inout.BaseResponse{
        ErrorCode:        statusCode,
        ErrorDescription: err.Error(),
    }
    context.JSON(statusCode, response)
}

// Usage in controller
func (controller ProductController) GetProduct(context *gin.Context) {
    id, err := strconv.Atoi(context.Param("id"))
    if err != nil {
        util.HandleError(context, errors.New("invalid product ID"), http.StatusBadRequest)
        return
    }
    
    dao := dao2.ProductDao{}
    err, product := dao.GetByID(uint(id))
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            util.HandleError(context, errors.New("product not found"), http.StatusNotFound)
        } else {
            util.HandleError(context, err, http.StatusInternalServerError)
        }
        return
    }
    
    context.JSON(http.StatusOK, gin.H{
        "error_code": 0,
        "error_description": "Success",
        "data": product,
    })
}
```

## Best Practices

### 1. Project Structure
- Keep related functionality grouped together
- Use consistent naming conventions
- Separate concerns across layers
- Keep models, DTOs, and business logic separate

### 2. Error Handling
```go
// Always use consistent error responses
type ErrorResponse struct {
    ErrorCode        int    `json:"error_code"`
    ErrorDescription string `json:"error_description"`
    Details          interface{} `json:"details,omitempty"`
}

// Use custom error types for better handling
type ValidationError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
}
```

### 3. Security
- Always validate input data
- Use parameterized queries (GORM handles this)
- Never expose sensitive information in responses
- Implement proper authentication and authorization
- Use HTTPS in production
- Validate JWT tokens properly

### 4. Database
- Use transactions for complex operations
- Implement proper indexing
- Use connection pooling
- Handle database migrations properly

```go
// Example transaction usage
func (dao *ProductDao) CreateWithStock(product *model.Product, stock *model.Stock) error {
    tx := Database.Begin()
    
    if err := tx.Create(&product).Error; err != nil {
        tx.Rollback()
        return err
    }
    
    stock.ProductID = product.ID
    if err := tx.Create(&stock).Error; err != nil {
        tx.Rollback()
        return err
    }
    
    return tx.Commit().Error
}
```

### 5. Testing
- Write unit tests for all layers
- Use test databases
- Mock external dependencies
- Test error scenarios
- Implement integration tests

### 6. Performance
- Use pagination for large datasets
- Implement caching where appropriate
- Use database indexes
- Profile your application regularly
- Monitor memory usage

### 7. Configuration
- Use environment variables for configuration
- Never commit secrets to version control
- Use different configurations for different environments
- Validate configuration at startup

### 8. Logging
```go
// Use structured logging
import "github.com/sirupsen/logrus"

func setupLogging() {
    logrus.SetFormatter(&logrus.JSONFormatter{})
    logrus.SetLevel(logrus.InfoLevel)
}

// Log important events
logrus.WithFields(logrus.Fields{
    "user_id": userID,
    "action":  "user_created",
}).Info("New user registered")
```

### 9. API Versioning
- Always version your APIs
- Maintain backward compatibility
- Use semantic versioning
- Document breaking changes

### 10. Documentation
- Keep API documentation up to date
- Use Swagger annotations
- Document complex business logic
- Provide examples for all endpoints

## Common Patterns

### Repository Pattern
```go
type UserRepository interface {
    Create(user *model.User) error
    GetByID(id uint) (*model.User, error)
    Update(user *model.User) error
    Delete(id uint) error
}

type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{db: db}
}
```

### Service Pattern
```go
type UserService interface {
    CreateUser(userData CreateUserRequest) (*User, error)
    GetUser(id uint) (*User, error)
}

type userService struct {
    repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
    return &userService{repo: repo}
}
```

This documentation provides a comprehensive guide for building scalable REST APIs with Go, following the established patterns in your codebase while incorporating modern best practices and real-world examples.