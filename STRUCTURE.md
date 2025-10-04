# TestLake Backend - Project Structure & Development Guide

This document provides a comprehensive guide to the TestLake backend architecture, including detailed templates, configurations, and development patterns for building scalable Go REST APIs.

## Technology Stack

- **Language**: Go 1.24+
- **Web Framework**: Gin (github.com/gin-gonic/gin)
- **Database**: PostgreSQL with GORM ORM
- **Authentication**: JWT with refresh tokens
- **Documentation**: Swagger/OpenAPI with swaggo
- **Testing**: testify framework with comprehensive test coverage
- **Email**: GoMail with HTML templates
- **Environment**: godotenv for configuration management
- **Containerization**: Docker and Docker Compose
- **Payments**: PayPal API integration

## Architecture Pattern

The project follows a **layered MVC architecture** with additional separation of concerns:

```
┌─────────────────┐
│   Controllers   │ ← Business logic, HTTP handlers, request/response processing
├─────────────────┤
│    Services     │ ← API route declarations and Swagger documentation
├─────────────────┤
│      DAOs       │ ← Data Access Objects (database layer)
├─────────────────┤
│     Models      │ ← Database entity definitions
└─────────────────┘
```

## Complete Directory Structure

```
testlake/
├── main.go               # Application entry point - loads env, connects DB, starts server
├── go.mod                # Go module definition with dependencies
├── go.sum                # Go module checksums for reproducible builds
├── Dockerfile            # Multi-stage Docker build configuration
├── docker-compose.yml    # PostgreSQL + API service orchestration
├── server.log            # Application runtime logs
├── .env                  # Environment variables (not in version control)
├── README.md             # Project overview and setup instructions
├── CLAUDE.md             # AI assistant project context and instructions
├── SPEC.md               # Business requirements and API specifications
└── STRUCTURE.md          # This file - architecture and development guide
│
├── app/                  # 🚀 Application Bootstrap & Routing Layer
│   ├── app.go           # Main application setup, middleware configuration, server start
│   ├── routes.go        # Route registration (public/private), service orchestration
│   └── swagger.go       # Swagger UI configuration and documentation setup
│
├── controller/           # 🎯 Business Logic & HTTP Request Handlers
│   ├── auth_controller.go          # Authentication: signup, signin, password reset, email verification
│   ├── user_controller.go          # User management: profile, dashboard, notifications, invites
│   ├── organization_controller.go  # Organization: CRUD, members, invitations, roles
│   ├── billing_controller.go       # Billing: overview, history, invoices, payments
│   ├── subscription_controller.go  # Subscriptions: CRUD, plan changes, usage tracking
│   ├── payment_method_controller.go # Payment methods: CRUD, default selection, PayPal integration
│   └── plan_controller.go          # Plans: listing, comparison, feature details
│
├── service/              # 🛣️  Route Definitions & Swagger Documentation
│   ├── auth_service.go          # Auth route declarations with Swagger annotations
│   ├── user_service.go          # User management route declarations
│   ├── organization_service.go  # Organization route declarations
│   ├── billing_service.go       # Billing and invoice route declarations
│   ├── subscription_service.go  # Subscription route declarations
│   ├── payment_method_service.go # Payment method route declarations
│   └── plan_service.go          # Plan route declarations
│
├── dao/                  # 🗄️  Data Access Layer & Database Operations
│   ├── main.go                     # Database connection, migration, configuration
│   ├── user_dao.go                 # User CRUD operations and queries
│   ├── organization_dao.go         # Organization CRUD and complex queries
│   ├── organization_member_dao.go  # Membership management, role assignments
│   ├── organization_usage_dao.go   # Usage tracking and billing calculations
│   ├── email_verification_dao.go   # Email verification token management
│   ├── billing_event_dao.go        # Billing event logging and retrieval
│   ├── invoice_dao.go              # Invoice generation, payment tracking
│   ├── payment_dao.go              # Payment processing, transaction history
│   ├── payment_method_dao.go       # Payment method storage, PayPal integration
│   ├── plan_dao.go                 # Plan definitions, features, pricing
│   └── subscription_dao.go         # Subscription lifecycle, plan changes
│
├── model/                # 📊 Database Entity Definitions & Relationships
│   ├── user.go                     # User accounts, authentication, profiles
│   ├── organization.go             # Organizations, settings, metadata
│   ├── organization_member.go      # Membership relationships, roles, permissions
│   ├── organization_usage.go       # Usage tracking, quotas, billing metrics
│   ├── email_verification_token.go # Email verification workflow
│   ├── project.go                  # TestLake projects (core business entity)
│   ├── project_access.go           # Project access control, permissions
│   ├── team.go                     # Team structures within organizations
│   ├── environment.go              # Development environments (dev, staging, prod)
│   ├── feature.go                  # Feature definitions, test data generation
│   ├── data_schema.go              # Schema definitions for test data
│   ├── test_data.go                # Generated test data records
│   ├── billing_event.go            # Billing events, usage tracking
│   ├── invoice.go                  # Invoice generation, payment status
│   ├── payment.go                  # Payment transactions, reconciliation
│   ├── payment_method.go           # Stored payment methods, PayPal accounts
│   ├── plan.go                     # Subscription plans, features, limits
│   └── subscription.go             # Active subscriptions, billing cycles
│
├── inout/                # 📥📤 Input/Output DTOs & API Contracts
│   ├── base_response.go    # Common response structures, error handling
│   ├── auth/              # Authentication request/response DTOs
│   │   ├── auth_in.go     # Login, signup, password reset requests
│   │   └── auth_out.go    # Authentication responses, token data
│   ├── user/              # User management DTOs
│   │   ├── user_in.go     # Profile updates, preference changes
│   │   └── user_out.go    # User profile data, dashboard information
│   ├── organization/      # Organization management DTOs
│   │   ├── organization_in.go  # Organization creation, updates, invitations
│   │   └── organization_out.go # Organization data, member lists, roles
│   ├── billing/           # Billing and financial DTOs
│   │   ├── billing_in.go  # Payment processing, invoice requests
│   │   └── billing_out.go # Billing history, invoice data, usage reports
│   ├── payment/           # Payment method DTOs
│   │   ├── payment_in.go  # Payment method creation, updates
│   │   └── payment_out.go # Payment method details, transaction history
│   ├── subscription/      # Subscription management DTOs
│   │   ├── subscription_in.go  # Plan changes, cancellation requests
│   │   └── subscription_out.go # Subscription status, usage data
│   ├── plan/              # Plan comparison DTOs
│   │   └── plan_out.go    # Plan features, pricing, comparison data
│   └── usage/             # Usage tracking DTOs
│       └── usage_out.go   # Usage metrics, quota information
│
├── middleware/           # 🛡️  HTTP Middleware Components
│   ├── cors.go          # Cross-Origin Resource Sharing configuration
│   └── jwt.go           # JWT authentication, token validation, user context
│
├── utils/               # 🔧 Utility Functions & Helper Services
│   ├── email.go         # Email sending, template rendering, SMTP configuration
│   ├── error_handler.go # Centralized error handling, HTTP status codes
│   ├── helpers.go       # General utility functions, string manipulation, validation
│   ├── jwt.go           # JWT token generation, validation, claims extraction
│   └── password.go      # Password hashing, validation, security utilities
│
├── templates/           # 📧 HTML Email Templates
│   ├── email_confirmation.html      # Initial email verification template
│   ├── email_confirmation_resend.html # Resend verification template
│   ├── email_verification_error.html  # Verification error template
│   └── email_verified_success.html    # Successful verification template
│
├── docs/                # 📖 Auto-Generated API Documentation
│   ├── docs.go          # Generated Swagger documentation code
│   ├── swagger.json     # OpenAPI 3.0 specification (JSON format)
│   └── swagger.yaml     # OpenAPI 3.0 specification (YAML format)
│
├── proj_docs/           # 📚 Project Documentation & Specifications
│   ├── GO_REST_API_DOCUMENTATION.md # Complete backend architecture guide
│   ├── testlake_spec.md             # TestLake platform business requirements
│   ├── testlake_example.md          # Real-world implementation examples
│   └── testlake_payments.md         # Payment system and billing documentation
│
├── routes/              # 🛤️  Additional Route Grouping (Future Use)
│
└── *_test/             # 🧪 Comprehensive Testing Suite
    ├── controller_test/ # HTTP endpoint integration tests
    │   ├── auth_controller_test.go      # Authentication flow testing
    │   ├── user_controller_test.go      # User management testing
    │   └── logs/                        # Test execution logs
    │       └── email_errors.log
    ├── service_test/    # Service layer integration tests
    │   ├── auth_service_test.go         # Authentication service testing
    │   ├── user_service_test.go         # User service testing
    │   └── user_management_service_test.go # Complete user management testing
    ├── dao_test/        # Database layer unit tests
    │   ├── user_dao_test.go             # User database operations testing
    │   └── organization_dao_test.go     # Organization database operations testing
    ├── model_test/      # Model validation and relationship tests
    │   └── email_verification_token_test.go # Email verification model testing
    └── utils_test/      # Utility function unit tests
        ├── email_utils_test.go          # Email utility testing
        └── logs/                        # Utility test logs
            └── email_errors.log
```

## Key Design Principles

### 1. **Separation of Concerns**
- **Controllers**: Contain business logic, handle HTTP requests/responses, validation, serialization
- **Services**: Define API routes and contain Swagger documentation for endpoints
- **DAOs**: Handle database operations, queries, transactions
- **Models**: Define data structures and relationships

### 2. **Dependency Flow**
- Controllers depend on DAOs directly for data operations
- Services depend on Controllers for handling business logic
- Clean interfaces between layers

### 3. **Input/Output DTOs**
- Separate request/response structures from database models
- Located in `inout/` directory, organized by domain
- Provides API versioning flexibility

### 4. **Error Handling**
- Centralized error handling utilities in `utils/error_handler.go`
- Consistent error response format via `base_response.go`
- Proper HTTP status codes

### 5. **Authentication & Security**
- JWT-based authentication with refresh tokens
- Middleware for route protection
- Password hashing utilities
- Email verification system

### 6. **Testing Strategy**
- Comprehensive test coverage across all layers
- Separate test directories for organization
- Unit tests for services and utilities
- Integration tests for DAOs and controllers

### 7. **Documentation**
- Auto-generated Swagger/OpenAPI documentation
- Comprehensive project documentation in `proj_docs/`
- Code-level documentation following Go conventions

## Development Workflow

1. **Model First**: Define database entities in `model/`
2. **DAO Layer**: Implement database operations in `dao/`
3. **DTOs**: Define request/response structures in `inout/`
4. **Controller Layer**: Implement business logic and HTTP handlers in `controller/`
5. **Service Layer**: Define API routes and add Swagger documentation in `service/`
6. **Routes**: Register services in `app/routes.go`
7. **Testing**: Write tests for each layer
8. **Documentation**: Swagger annotations are in service layer

## Configuration & Environment Setup

### Environment Variables (.env)

```bash
# Database Configuration
DB_HOST=localhost
DB_USER=testlake_user
DB_PASSWORD=secure_password
DB_NAME=testlake_db
DB_PORT=5432

# Server Configuration
IP=0.0.0.0
PORT=8000
SCHEME=http

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-min-32-chars
JWT_REFRESH_SECRET=your-refresh-secret-key-min-32-chars
JWT_EXPIRES_IN=24h
JWT_REFRESH_EXPIRES_IN=168h

# Email Configuration (SMTP)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
FROM_EMAIL=noreply@testlake.com
FROM_NAME=TestLake Platform

# PayPal Configuration
PAYPAL_CLIENT_ID=your-paypal-client-id
PAYPAL_CLIENT_SECRET=your-paypal-client-secret
PAYPAL_MODE=sandbox # or live for production

# Application URLs
FRONTEND_URL=http://localhost:3000
API_BASE_URL=http://localhost:8000/api/v1

# File Upload Configuration
MAX_UPLOAD_SIZE=10MB
UPLOAD_PATH=./uploads

# Rate Limiting
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=1h

# Logging
LOG_LEVEL=info
LOG_FILE=server.log
```

### Database Configuration (dao/main.go)

The database configuration handles PostgreSQL connection, GORM setup, and automatic migrations:

```go
// dao/main.go - Database Configuration Template
package dao

import (
	"fmt"
	"log"
	"os"
	"testlake/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Database *gorm.DB

func Connect() {
	// Load environment variables
	host := os.Getenv("DB_HOST")
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	databaseName := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	// Build PostgreSQL connection string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		host, username, password, databaseName, port)

	// Connect to database with GORM
	var err error
	Database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run automatic migrations for all models
	err = Database.AutoMigrate(
		// Core TestLake Models
		&model.User{},
		&model.Organization{},
		&model.OrganizationMember{},
		&model.OrganizationInvitation{},
		&model.Project{},
		&model.Team{},
		&model.TeamMember{},
		&model.ProjectAccess{},
		&model.Environment{},
		&model.Feature{},
		&model.FeatureEnvironmentStatus{},
		&model.FeatureErrorLog{},
		&model.ErrorImage{},
		&model.DataSchema{},
		&model.FeatureSchema{},
		&model.SchemaField{},
		&model.TestData{},
		&model.TestDataRequest{},
		
		// Authentication & Verification
		&model.EmailVerificationToken{},
		
		// Payment & Billing Models
		&model.PaymentMethod{},
		&model.Subscription{},
		&model.Invoice{},
	)
	
	if err != nil {
		log.Fatal("Failed to run database migrations:", err)
	}

	log.Println("✅ Database connected and migrated successfully")
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return Database
}
```

### Application Bootstrap (main.go)

```go
// main.go - Application Entry Point Template
package main

import (
	"log"
	"testlake/app"
	"testlake/dao"
	_ "testlake/docs" // Import for Swagger docs generation
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("❌ Error loading .env file:", err)
	}

	// Connect to database and run migrations
	dao.Connect()

	// Start the web server
	app.ServeApplication()
}
```

### Application Setup (app/app.go)

```go
// app/app.go - Application Configuration Template
package app

import (
	"os"
	"testlake/middleware"
	"testlake/utils"
	"github.com/gin-gonic/gin"
)

// @title TestLake API
// @version 1.0
// @description Test Data Management Platform API
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8000
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func ServeApplication() {
	// Initialize Gin router with default middleware (logger, recovery)
	router := gin.Default()

	// Apply global middleware
	router.Use(middleware.DefaultAuthMiddleware())
	router.Use(middleware.CORSMiddleware())

	// Setup Swagger documentation
	Swagger(router)

	// Setup API route groups
	baseRoute := router.Group("/api/v1")

	// Public routes (no authentication required)
	publicRoutes := baseRoute.Group("")
	PublicRoutes(publicRoutes)

	// Private routes (JWT authentication required)
	privateRoutes := baseRoute.Group("")
	privateRoutes.Use(middleware.JWTAuthMiddleware())
	PrivateRoutes(privateRoutes)

	// Handle 404 routes
	router.NoRoute(utils.HandleNoRoute())

	// Start server
	ip := os.Getenv("IP")
	port := os.Getenv("PORT")
	
	log.Printf("🚀 Server starting on %s:%s", ip, port)
	log.Printf("📖 Swagger docs available at http://%s:%s/swagger/index.html", ip, port)
	
	router.Run(ip + ":" + port)
}
```

## Features

- **Authentication**: User registration, login, email verification
- **Organizations**: Multi-tenant organization management
- **Billing**: Subscription management, payment processing
- **User Management**: Profile management, password reset
- **Email System**: HTML templates, verification workflows
- **API Documentation**: Auto-generated Swagger docs
- **Error Handling**: Centralized, consistent error responses
- **Testing**: Comprehensive test suite
- **Docker Support**: Ready for containerization

## Code Examples

### Model Definition Example
```go
// model/user.go
package model

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID     `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Email        string        `gorm:"uniqueIndex;not null" json:"email"`
	Username     string        `gorm:"uniqueIndex;not null" json:"username"`
	FirstName    *string       `json:"first_name,omitempty"`
	LastName     *string       `json:"last_name,omitempty"`
	PasswordHash *string       `json:"-"`
	AuthProvider AuthProvider  `gorm:"not null" json:"auth_provider"`
	Status       UserStatus    `gorm:"default:'active'" json:"status"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type AuthProvider string
const (
	AuthProviderEmail  AuthProvider = "email"
	AuthProviderGoogle AuthProvider = "google"
)

type UserStatus string
const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusBanned   UserStatus = "banned"
)
```

### DAO Implementation Example
```go
// dao/user_dao.go
package dao

import (
	"testlake/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao() *UserDao {
	return &UserDao{db: GetDB()}
}

func (dao *UserDao) Create(user *model.User) error {
	return dao.db.Create(user).Error
}

func (dao *UserDao) GetByID(id uuid.UUID) (*model.User, error) {
	var user model.User
	err := dao.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (dao *UserDao) GetByEmail(email string) (*model.User, error) {
	var user model.User
	err := dao.db.First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (dao *UserDao) EmailExists(email string) (bool, error) {
	var count int64
	err := dao.db.Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

func (dao *UserDao) Update(user *model.User) error {
	return dao.db.Save(user).Error
}

func (dao *UserDao) Delete(id uuid.UUID) error {
	return dao.db.Delete(&model.User{}, "id = ?", id).Error
}
```

### Input/Output DTOs Example
```go
// inout/user/user_in.go
package user

import "testlake/model"

type UpdateUserRequest struct {
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Username  string  `json:"username" binding:"required"`
}

// inout/user/user_out.go
package user

import (
	"testlake/inout"
	"testlake/model"
	"time"
	"github.com/google/uuid"
)

type UserOut struct {
	inout.BaseResponse
	Data UserData `json:"data"`
}

type UserData struct {
	ID           uuid.UUID           `json:"id"`
	Email        string              `json:"email"`
	Username     string              `json:"username"`
	FirstName    *string             `json:"first_name,omitempty"`
	LastName     *string             `json:"last_name,omitempty"`
	AuthProvider model.AuthProvider  `json:"auth_provider"`
	Status       model.UserStatus    `json:"status"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}

func FromModel(user *model.User) UserData {
	return UserData{
		ID:           user.ID,
		Email:        user.Email,
		Username:     user.Username,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		AuthProvider: user.AuthProvider,
		Status:       user.Status,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}
```

### Controller Implementation Example
```go
// controller/user_controller.go
package controller

import (
	"errors"
	"net/http"
	"testlake/dao"
	"testlake/inout"
	"testlake/inout/user"
	"testlake/model"
	"testlake/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct{}

func (controller UserController) GetProfile(context *gin.Context) {
	// Extract user ID from JWT token
	userID, err := utils.ExtractUserID(context)
	if err != nil {
		utils.ReportUnauthorized(context, "Authentication required")
		return
	}

	// Get user from database
	userDao := dao.NewUserDao()
	foundUser, err := userDao.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ReportNotFound(context, "User not found")
		} else {
			utils.ReportInternalServerError(context, "Database error")
		}
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

func (controller UserController) UpdateProfile(context *gin.Context) {
	// Extract user ID from JWT token
	userID, err := utils.ExtractUserID(context)
	if err != nil {
		utils.ReportUnauthorized(context, "Authentication required")
		return
	}

	// Parse request body
	var request user.UpdateUserRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		utils.ReportBadRequest(context, "Invalid request data")
		return
	}

	// Get existing user
	userDao := dao.NewUserDao()
	foundUser, err := userDao.GetByID(userID)
	if err != nil {
		utils.ReportNotFound(context, "User not found")
		return
	}

	// Update user fields
	foundUser.Username = request.Username
	foundUser.FirstName = request.FirstName
	foundUser.LastName = request.LastName

	// Save to database
	if err := userDao.Update(foundUser); err != nil {
		utils.ReportInternalServerError(context, "Failed to update user")
		return
	}

	// Return updated user
	response := user.UserOut{
		BaseResponse: inout.BaseResponse{
			ErrorCode:        0,
			ErrorDescription: "Profile updated successfully",
		},
		Data: user.FromModel(foundUser),
	}

	context.JSON(http.StatusOK, response)
}
```

### Service Declaration Example
```go
// service/user_service.go
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
// @Param Authorization header string true "Bearer token" format(Bearer {token})
// @Success 200 {object} user.UserOut
// @Failure 401 {object} inout.BaseResponse
// @Router /api/v1/users/profile [GET]
func (s UserService) GetProfile(r *gin.RouterGroup, route string) {
	r.GET("/"+s.Route+"/"+route, s.Controller.GetProfile)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update current user's profile information
// @Tags User Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token" format(Bearer {token})
// @Param user body user.UpdateUserRequest true "User profile data"
// @Success 200 {object} user.UserOut
// @Failure 400 {object} inout.BaseResponse
// @Failure 401 {object} inout.BaseResponse
// @Router /api/v1/users/profile [PUT]
func (s UserService) UpdateProfile(r *gin.RouterGroup, route string) {
	r.PUT("/"+s.Route+"/"+route, s.Controller.UpdateProfile)
}
```

### Middleware Example
```go
// middleware/jwt.go
package middleware

import (
	"net/http"
	"strings"
	"testlake/utils"

	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error_code":        401,
				"error_description": "Authorization header required",
			})
			c.Abort()
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error_code":        401,
				"error_description": "Invalid token",
			})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Next()
	}
}
```

## Testing Documentation

### Testing Strategy Overview

The project implements a **three-tier testing approach**:

1. **Unit Tests**: Test individual components in isolation
2. **Integration Tests**: Test API endpoints with database interactions
3. **Service Tests**: Test complete request-response cycles against a running server

### Test Structure

```
*_test/               # Test directories organized by layer
├── dao_test/         # Database layer tests
├── controller_test/  # HTTP endpoint tests
├── service_test/     # Integration tests with live server
├── model_test/       # Model validation tests
└── utils_test/       # Utility function tests
```

### Testing Technologies

- **Framework**: `testify/assert` for assertions
- **HTTP Testing**: `httptest` for mocking HTTP requests
- **Database**: Uses same database as development with test data cleanup
- **Environment**: `godotenv` for loading test environment variables

### DAO Testing Example

```go
// dao_test/user_dao_test.go
package dao_test

import (
	"testing"
	"time"
	"fmt"
	"os"

	"testlake/dao"
	"testlake/model"
	
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	err := godotenv.Load("../.env")
	if err != nil {
		panic("Error loading .env file")
	}
	
	dao.Connect()
	code := m.Run()
	os.Exit(code)
}

func TestUserDao_Create(t *testing.T) {
	userDao := dao.NewUserDao()
	
	timestamp := time.Now().UnixNano()
	user := &model.User{
		Email:        fmt.Sprintf("test%d@example.com", timestamp),
		Username:     fmt.Sprintf("testuser%d", timestamp),
		AuthProvider: model.AuthProviderEmail,
		Status:       model.UserStatusActive,
	}
	
	err := userDao.Create(user)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, user.ID)
	
	// Cleanup
	defer userDao.Delete(user.ID)
}

func TestUserDao_GetByEmail(t *testing.T) {
	userDao := dao.NewUserDao()
	
	// Create test user
	timestamp := time.Now().UnixNano()
	user := &model.User{
		Email:        fmt.Sprintf("test%d@example.com", timestamp),
		Username:     fmt.Sprintf("testuser%d", timestamp),
		AuthProvider: model.AuthProviderEmail,
		Status:       model.UserStatusActive,
	}
	
	err := userDao.Create(user)
	assert.NoError(t, err)
	
	// Test retrieval
	foundUser, err := userDao.GetByEmail(user.Email)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, foundUser.ID)
	assert.Equal(t, user.Username, foundUser.Username)
	
	// Cleanup
	defer userDao.Delete(user.ID)
}
```

### Controller Testing Example

```go
// controller_test/user_controller_test.go
package controller_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"testlake/app"
	"testlake/dao"
	"testlake/inout/auth"
	"testlake/model"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	router := gin.New()
	baseRoute := router.Group("/api/v1")
	publicRoutes := baseRoute.Group("/public")
	app.PublicRoutes(publicRoutes)
	return router
}

func TestCreateUser(t *testing.T) {
	router := setupRouter()
	
	uniqueID := time.Now().UnixNano()
	requestData := auth.SignUpRequest{
		Email:        fmt.Sprintf("test%d@example.com", uniqueID),
		Username:     fmt.Sprintf("testuser%d", uniqueID),
		Password:     "password123",
		AuthProvider: model.AuthProviderEmail,
	}
	
	jsonData, _ := json.Marshal(requestData)
	
	req, _ := http.NewRequest("POST", "/api/v1/public/auth/signup", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusCreated, w.Code)
	
	var response auth.SignUpOut
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 0, response.ErrorCode)
	assert.Equal(t, requestData.Email, response.Data.User.Email)
	
	// Cleanup
	userDao := dao.NewUserDao()
	defer userDao.Delete(response.Data.User.ID)
}

func TestLogin(t *testing.T) {
	router := setupRouter()
	userDao := dao.NewUserDao()
	
	// Create test user with hashed password
	uniqueID := time.Now().UnixNano()
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	
	testUser := &model.User{
		Email:        fmt.Sprintf("login%d@example.com", uniqueID),
		Username:     fmt.Sprintf("loginuser%d", uniqueID),
		AuthProvider: model.AuthProviderEmail,
		PasswordHash: string(hashedPassword),
		Status:       model.UserStatusActive,
	}
	
	err := userDao.Create(testUser)
	assert.NoError(t, err)
	defer userDao.Delete(testUser.ID)
	
	// Test login
	loginRequest := auth.SignInRequest{
		Email:    testUser.Email,
		Password: password,
	}
	
	jsonData, _ := json.Marshal(loginRequest)
	req, _ := http.NewRequest("POST", "/api/v1/public/auth/signin", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response auth.SignInOut
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 0, response.ErrorCode)
	assert.NotEmpty(t, response.Data.Token)
}
```

### Service Integration Testing Example

```go
// service_test/user_service_test.go
package service_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"testlake/inout/auth"
	"testlake/model"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func getBaseURL() string {
	ip := os.Getenv("IP")
	port := os.Getenv("PORT")
	scheme := os.Getenv("SCHEME")
	return scheme + "://" + ip + ":" + port
}

func TestCreateUserIntegration(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatal("Error loading .env file")
	}
	
	baseURL := getBaseURL()
	url := baseURL + "/api/v1/public/auth/signup"
	
	userData := auth.SignUpRequest{
		Email:        fmt.Sprintf("integration%d@example.com", time.Now().UnixNano()),
		Username:     fmt.Sprintf("integrationuser%d", time.Now().UnixNano()),
		Password:     "password123",
		AuthProvider: model.AuthProviderEmail,
	}
	
	jsonData, _ := json.Marshal(userData)
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Skip("Server not running - skipping integration test")
		return
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	
	var result auth.SignUpOut
	err = json.Unmarshal(body, &result)
	assert.NoError(t, err)
	
	if resp.StatusCode == http.StatusCreated {
		assert.Equal(t, 0, result.ErrorCode)
		assert.Equal(t, userData.Email, result.Data.User.Email)
	} else {
		t.Logf("Expected status 201, got %v. Response: %s", resp.StatusCode, string(body))
	}
}
```

### Test Execution Commands

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test package
go test ./dao_test/
go test ./controller_test/
go test ./service_test/

# Run tests with verbose output
go test -v ./...

# Run tests with race detection
go test -race ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Testing Best Practices

1. **Isolation**: Each test should be independent and not rely on other tests
2. **Cleanup**: Always clean up test data using `defer` statements
3. **Unique Data**: Use timestamps or UUIDs to ensure unique test data
4. **Environment**: Load test environment variables in `TestMain`
5. **Assertions**: Use descriptive assertions that clearly indicate what failed
6. **Mocking**: Mock external dependencies when testing business logic
7. **Coverage**: Aim for high test coverage, especially for critical business logic

### Test Data Management

- **Unique Identifiers**: Use `time.Now().UnixNano()` for unique test data
- **Cleanup Strategy**: Always use `defer` to clean up created test data
- **Database State**: Tests should not depend on existing database state
- **Isolation**: Each test creates and cleans up its own data

This comprehensive testing strategy ensures reliability, maintainability, and confidence in code changes across all layers of the application.

## Complete Development Templates

### 1. Complete Model Template with Advanced Features

```go
// model/example_entity.go - Complete Model Template
package model

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ExampleEntity represents a business entity with full GORM features
type ExampleEntity struct {
	// Primary Key with UUID
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	
	// Required fields with validation
	Name        string `gorm:"not null;size:255" json:"name" validate:"required,min=1,max=255"`
	Email       string `gorm:"uniqueIndex;not null;size:255" json:"email" validate:"required,email"`
	Status      EntityStatus `gorm:"default:'active'" json:"status"`
	
	// Optional fields
	Description *string `gorm:"type:text" json:"description,omitempty"`
	Metadata    JSON    `gorm:"type:jsonb" json:"metadata,omitempty"`
	
	// Relationships
	OrganizationID uuid.UUID    `gorm:"type:uuid;not null" json:"organization_id"`
	Organization   Organization `gorm:"foreignKey:OrganizationID;constraint:OnDelete:CASCADE" json:"organization,omitempty"`
	
	// One-to-Many relationship
	RelatedItems []RelatedItem `gorm:"foreignKey:ExampleEntityID" json:"related_items,omitempty"`
	
	// Many-to-Many relationship
	Tags []Tag `gorm:"many2many:example_entity_tags" json:"tags,omitempty"`
	
	// Timestamps (GORM automatically manages these)
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete
}

// Custom JSON type for PostgreSQL JSONB
type JSON map[string]interface{}

// Entity status enumeration
type EntityStatus string

const (
	EntityStatusActive   EntityStatus = "active"
	EntityStatusInactive EntityStatus = "inactive"
	EntityStatusArchived EntityStatus = "archived"
)

// GORM Hooks - Called automatically by GORM
func (e *ExampleEntity) BeforeCreate(tx *gorm.DB) error {
	// Custom logic before creating entity
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}

func (e *ExampleEntity) AfterCreate(tx *gorm.DB) error {
	// Custom logic after creating entity (e.g., logging, notifications)
	return nil
}

// Custom methods
func (e *ExampleEntity) IsActive() bool {
	return e.Status == EntityStatusActive
}

func (e *ExampleEntity) GetDisplayName() string {
	if e.Name != "" {
		return e.Name
	}
	return e.Email
}

// Table name override (optional)
func (ExampleEntity) TableName() string {
	return "custom_table_name"
}
```

### 2. Complete DAO Template with Advanced Queries

```go
// dao/example_entity_dao.go - Complete DAO Template
package dao

import (
	"fmt"
	"testlake/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExampleEntityDao struct {
	db *gorm.DB
}

func NewExampleEntityDao() *ExampleEntityDao {
	return &ExampleEntityDao{db: GetDB()}
}

// Basic CRUD Operations
func (dao *ExampleEntityDao) Create(entity *model.ExampleEntity) error {
	return dao.db.Create(entity).Error
}

func (dao *ExampleEntityDao) GetByID(id uuid.UUID) (*model.ExampleEntity, error) {
	var entity model.ExampleEntity
	err := dao.db.Preload("Organization").Preload("RelatedItems").Preload("Tags").
		First(&entity, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (dao *ExampleEntityDao) Update(entity *model.ExampleEntity) error {
	return dao.db.Save(entity).Error
}

func (dao *ExampleEntityDao) Delete(id uuid.UUID) error {
	return dao.db.Delete(&model.ExampleEntity{}, "id = ?", id).Error
}

// Advanced Query Methods
func (dao *ExampleEntityDao) GetByOrganization(orgID uuid.UUID, page, pageSize int) ([]model.ExampleEntity, int64, error) {
	var entities []model.ExampleEntity
	var total int64
	
	// Count total records
	dao.db.Model(&model.ExampleEntity{}).Where("organization_id = ?", orgID).Count(&total)
	
	// Get paginated results
	offset := (page - 1) * pageSize
	err := dao.db.Where("organization_id = ?", orgID).
		Preload("Organization").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&entities).Error
		
	return entities, total, err
}

func (dao *ExampleEntityDao) SearchByName(name string, orgID uuid.UUID) ([]model.ExampleEntity, error) {
	var entities []model.ExampleEntity
	searchPattern := fmt.Sprintf("%%%s%%", name)
	
	err := dao.db.Where("name ILIKE ? AND organization_id = ?", searchPattern, orgID).
		Preload("Organization").
		Find(&entities).Error
		
	return entities, err
}

func (dao *ExampleEntityDao) GetActiveEntities(orgID uuid.UUID) ([]model.ExampleEntity, error) {
	var entities []model.ExampleEntity
	err := dao.db.Where("organization_id = ? AND status = ?", orgID, model.EntityStatusActive).
		Find(&entities).Error
	return entities, err
}

// Complex queries with joins
func (dao *ExampleEntityDao) GetEntitiesWithStats(orgID uuid.UUID) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	
	err := dao.db.Table("example_entities").
		Select(`
			example_entities.id,
			example_entities.name,
			example_entities.status,
			COUNT(related_items.id) as related_items_count,
			MAX(related_items.created_at) as last_related_item
		`).
		Joins("LEFT JOIN related_items ON related_items.example_entity_id = example_entities.id").
		Where("example_entities.organization_id = ?", orgID).
		Group("example_entities.id, example_entities.name, example_entities.status").
		Scan(&results).Error
		
	return results, err
}

// Transaction support
func (dao *ExampleEntityDao) CreateWithRelatedItems(entity *model.ExampleEntity, items []model.RelatedItem) error {
	return dao.db.Transaction(func(tx *gorm.DB) error {
		// Create main entity
		if err := tx.Create(entity).Error; err != nil {
			return err
		}
		
		// Create related items
		for _, item := range items {
			item.ExampleEntityID = entity.ID
			if err := tx.Create(&item).Error; err != nil {
				return err
			}
		}
		
		return nil
	})
}

// Bulk operations
func (dao *ExampleEntityDao) BulkUpdateStatus(ids []uuid.UUID, status model.EntityStatus) error {
	return dao.db.Model(&model.ExampleEntity{}).
		Where("id IN ?", ids).
		Update("status", status).Error
}

func (dao *ExampleEntityDao) ExistsWithName(name string, orgID uuid.UUID) (bool, error) {
	var count int64
	err := dao.db.Model(&model.ExampleEntity{}).
		Where("name = ? AND organization_id = ?", name, orgID).
		Count(&count).Error
	return count > 0, err
}

// Custom raw queries for complex operations
func (dao *ExampleEntityDao) GetEntityStatistics(orgID uuid.UUID) (map[string]interface{}, error) {
	var result map[string]interface{}
	
	err := dao.db.Raw(`
		SELECT 
			COUNT(*) as total_entities,
			COUNT(CASE WHEN status = 'active' THEN 1 END) as active_entities,
			COUNT(CASE WHEN status = 'inactive' THEN 1 END) as inactive_entities,
			AVG(EXTRACT(EPOCH FROM (NOW() - created_at))) as avg_age_seconds
		FROM example_entities 
		WHERE organization_id = ? AND deleted_at IS NULL
	`, orgID).Scan(&result).Error
	
	return result, err
}
```

### 3. Complete Controller Template with Error Handling

```go
// controller/example_entity_controller.go - Complete Controller Template
package controller

import (
	"errors"
	"net/http"
	"strconv"
	"testlake/dao"
	"testlake/inout"
	"testlake/inout/example_entity"
	"testlake/model"
	"testlake/utils"
	
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExampleEntityController struct{}

// CreateEntity handles entity creation with full validation
func (controller ExampleEntityController) CreateEntity(context *gin.Context) {
	// Extract user context
	userID, err := utils.ExtractUserID(context)
	if err != nil {
		utils.ReportUnauthorized(context, "Authentication required")
		return
	}

	orgID, err := utils.ExtractOrganizationID(context)
	if err != nil {
		utils.ReportBadRequest(context, "Organization ID required")
		return
	}

	// Parse and validate request
	var request example_entity.CreateEntityRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		utils.ReportBadRequest(context, "Invalid request data: "+err.Error())
		return
	}

	// Additional business validation
	if err := controller.validateCreateRequest(&request, orgID); err != nil {
		utils.ReportBadRequest(context, err.Error())
		return
	}

	// Create entity
	entity := &model.ExampleEntity{
		Name:           request.Name,
		Email:          request.Email,
		Description:    request.Description,
		OrganizationID: orgID,
		Status:         model.EntityStatusActive,
	}

	entityDao := dao.NewExampleEntityDao()
	if err := entityDao.Create(entity); err != nil {
		utils.ReportInternalServerError(context, "Failed to create entity")
		return
	}

	// Return success response
	response := example_entity.EntityOut{
		BaseResponse: inout.BaseResponse{
			ErrorCode:        0,
			ErrorDescription: "Entity created successfully",
		},
		Data: example_entity.FromModel(entity),
	}

	context.JSON(http.StatusCreated, response)
}

// GetEntity retrieves a single entity with related data
func (controller ExampleEntityController) GetEntity(context *gin.Context) {
	// Parse entity ID
	entityIDStr := context.Param("id")
	entityID, err := uuid.Parse(entityIDStr)
	if err != nil {
		utils.ReportBadRequest(context, "Invalid entity ID")
		return
	}

	// Check permissions
	if !controller.canAccessEntity(context, entityID) {
		utils.ReportForbidden(context, "Access denied")
		return
	}

	// Retrieve entity
	entityDao := dao.NewExampleEntityDao()
	entity, err := entityDao.GetByID(entityID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ReportNotFound(context, "Entity not found")
		} else {
			utils.ReportInternalServerError(context, "Database error")
		}
		return
	}

	// Return entity data
	response := example_entity.EntityOut{
		BaseResponse: inout.BaseResponse{
			ErrorCode:        0,
			ErrorDescription: "Success",
		},
		Data: example_entity.FromModel(entity),
	}

	context.JSON(http.StatusOK, response)
}

// GetEntities retrieves paginated list of entities
func (controller ExampleEntityController) GetEntities(context *gin.Context) {
	// Extract organization ID
	orgID, err := utils.ExtractOrganizationID(context)
	if err != nil {
		utils.ReportBadRequest(context, "Organization ID required")
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(context.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(context.DefaultQuery("page_size", "20"))
	
	// Validate pagination
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// Retrieve entities
	entityDao := dao.NewExampleEntityDao()
	entities, total, err := entityDao.GetByOrganization(orgID, page, pageSize)
	if err != nil {
		utils.ReportInternalServerError(context, "Failed to retrieve entities")
		return
	}

	// Build response with pagination metadata
	response := example_entity.EntitiesListOut{
		BaseResponse: inout.BaseResponse{
			ErrorCode:        0,
			ErrorDescription: "Success",
		},
		Data: example_entity.EntitiesListData{
			Entities: example_entity.FromModels(entities),
			Pagination: inout.PaginationData{
				Page:       page,
				PageSize:   pageSize,
				Total:      total,
				TotalPages: (total + int64(pageSize) - 1) / int64(pageSize),
			},
		},
	}

	context.JSON(http.StatusOK, response)
}

// UpdateEntity handles entity updates with optimistic locking
func (controller ExampleEntityController) UpdateEntity(context *gin.Context) {
	// Parse entity ID
	entityIDStr := context.Param("id")
	entityID, err := uuid.Parse(entityIDStr)
	if err != nil {
		utils.ReportBadRequest(context, "Invalid entity ID")
		return
	}

	// Check permissions
	if !controller.canAccessEntity(context, entityID) {
		utils.ReportForbidden(context, "Access denied")
		return
	}

	// Parse request
	var request example_entity.UpdateEntityRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		utils.ReportBadRequest(context, "Invalid request data")
		return
	}

	// Retrieve existing entity
	entityDao := dao.NewExampleEntityDao()
	entity, err := entityDao.GetByID(entityID)
	if err != nil {
		utils.ReportNotFound(context, "Entity not found")
		return
	}

	// Update fields
	if request.Name != nil {
		entity.Name = *request.Name
	}
	if request.Description != nil {
		entity.Description = request.Description
	}
	if request.Status != nil {
		entity.Status = *request.Status
	}

	// Save changes
	if err := entityDao.Update(entity); err != nil {
		utils.ReportInternalServerError(context, "Failed to update entity")
		return
	}

	// Return updated entity
	response := example_entity.EntityOut{
		BaseResponse: inout.BaseResponse{
			ErrorCode:        0,
			ErrorDescription: "Entity updated successfully",
		},
		Data: example_entity.FromModel(entity),
	}

	context.JSON(http.StatusOK, response)
}

// DeleteEntity handles soft deletion
func (controller ExampleEntityController) DeleteEntity(context *gin.Context) {
	// Parse entity ID
	entityIDStr := context.Param("id")
	entityID, err := uuid.Parse(entityIDStr)
	if err != nil {
		utils.ReportBadRequest(context, "Invalid entity ID")
		return
	}

	// Check permissions
	if !controller.canAccessEntity(context, entityID) {
		utils.ReportForbidden(context, "Access denied")
		return
	}

	// Delete entity
	entityDao := dao.NewExampleEntityDao()
	if err := entityDao.Delete(entityID); err != nil {
		utils.ReportInternalServerError(context, "Failed to delete entity")
		return
	}

	// Return success response
	response := inout.BaseResponse{
		ErrorCode:        0,
		ErrorDescription: "Entity deleted successfully",
	}

	context.JSON(http.StatusOK, response)
}

// Private helper methods
func (controller ExampleEntityController) validateCreateRequest(request *example_entity.CreateEntityRequest, orgID uuid.UUID) error {
	// Check if entity with same name exists
	entityDao := dao.NewExampleEntityDao()
	exists, err := entityDao.ExistsWithName(request.Name, orgID)
	if err != nil {
		return errors.New("validation failed")
	}
	if exists {
		return errors.New("entity with this name already exists")
	}

	return nil
}

func (controller ExampleEntityController) canAccessEntity(context *gin.Context, entityID uuid.UUID) bool {
	// Implement permission checking logic
	// This could check organization membership, role permissions, etc.
	userID, _ := utils.ExtractUserID(context)
	// ... permission logic here ...
	_ = userID // Use userID for permission check
	return true // Simplified for template
}
```

### 4. Complete Input/Output DTO Templates

```go
// inout/example_entity/entity_in.go - Input DTOs Template
package example_entity

import (
	"testlake/model"
	"github.com/google/uuid"
)

// CreateEntityRequest represents data needed to create a new entity
type CreateEntityRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=255"`
	Email       string  `json:"email" binding:"required,email"`
	Description *string `json:"description,omitempty"`
	Tags        []uuid.UUID `json:"tags,omitempty"`
}

// UpdateEntityRequest represents data that can be updated
type UpdateEntityRequest struct {
	Name        *string             `json:"name,omitempty" binding:"omitempty,min=1,max=255"`
	Description *string             `json:"description,omitempty"`
	Status      *model.EntityStatus `json:"status,omitempty" binding:"omitempty,oneof=active inactive archived"`
	Tags        []uuid.UUID         `json:"tags,omitempty"`
}

// SearchEntityRequest for advanced search functionality
type SearchEntityRequest struct {
	Query          string             `json:"query,omitempty"`
	Status         *model.EntityStatus `json:"status,omitempty"`
	Tags           []uuid.UUID         `json:"tags,omitempty"`
	CreatedAfter   *time.Time          `json:"created_after,omitempty"`
	CreatedBefore  *time.Time          `json:"created_before,omitempty"`
	SortBy         string              `json:"sort_by,omitempty" binding:"omitempty,oneof=name created_at updated_at"`
	SortDirection  string              `json:"sort_direction,omitempty" binding:"omitempty,oneof=asc desc"`
}

// BulkActionRequest for bulk operations
type BulkActionRequest struct {
	EntityIDs []uuid.UUID         `json:"entity_ids" binding:"required,min=1"`
	Action    string              `json:"action" binding:"required,oneof=delete activate deactivate archive"`
	Status    *model.EntityStatus `json:"status,omitempty"`
}
```

```go
// inout/example_entity/entity_out.go - Output DTOs Template
package example_entity

import (
	"time"
	"testlake/inout"
	"testlake/model"
	"github.com/google/uuid"
)

// EntityOut represents a single entity response
type EntityOut struct {
	inout.BaseResponse
	Data EntityData `json:"data"`
}

// EntitiesListOut represents a list of entities with pagination
type EntitiesListOut struct {
	inout.BaseResponse
	Data EntitiesListData `json:"data"`
}

// EntityData represents the main entity data structure
type EntityData struct {
	ID             uuid.UUID           `json:"id"`
	Name           string              `json:"name"`
	Email          string              `json:"email"`
	Description    *string             `json:"description,omitempty"`
	Status         model.EntityStatus  `json:"status"`
	OrganizationID uuid.UUID           `json:"organization_id"`
	Organization   *OrganizationData   `json:"organization,omitempty"`
	RelatedItems   []RelatedItemData   `json:"related_items,omitempty"`
	Tags           []TagData           `json:"tags,omitempty"`
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at"`
}

// EntitiesListData represents paginated list data
type EntitiesListData struct {
	Entities   []EntityData         `json:"entities"`
	Pagination inout.PaginationData `json:"pagination"`
	Summary    *EntitySummary       `json:"summary,omitempty"`
}

// Related data structures
type OrganizationData struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type RelatedItemData struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type TagData struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Color *string   `json:"color,omitempty"`
}

// EntitySummary provides summary statistics
type EntitySummary struct {
	TotalEntities    int64            `json:"total_entities"`
	ActiveEntities   int64            `json:"active_entities"`
	InactiveEntities int64            `json:"inactive_entities"`
	ArchivedEntities int64            `json:"archived_entities"`
	StatusBreakdown  map[string]int64 `json:"status_breakdown"`
}

// EntityStats for dashboard/analytics
type EntityStatsOut struct {
	inout.BaseResponse
	Data EntityStatsData `json:"data"`
}

type EntityStatsData struct {
	Summary      EntitySummary         `json:"summary"`
	RecentItems  []EntityData          `json:"recent_items"`
	TopTags      []TagUsageData        `json:"top_tags"`
	GrowthMetrics EntityGrowthMetrics  `json:"growth_metrics"`
}

type TagUsageData struct {
	Tag   TagData `json:"tag"`
	Count int64   `json:"count"`
}

type EntityGrowthMetrics struct {
	ThisMonth     int64   `json:"this_month"`
	LastMonth     int64   `json:"last_month"`
	GrowthPercent float64 `json:"growth_percent"`
}

// Conversion functions
func FromModel(entity *model.ExampleEntity) EntityData {
	data := EntityData{
		ID:             entity.ID,
		Name:           entity.Name,
		Email:          entity.Email,
		Description:    entity.Description,
		Status:         entity.Status,
		OrganizationID: entity.OrganizationID,
		CreatedAt:      entity.CreatedAt,
		UpdatedAt:      entity.UpdatedAt,
	}

	// Convert organization if loaded
	if entity.Organization.ID != uuid.Nil {
		data.Organization = &OrganizationData{
			ID:   entity.Organization.ID,
			Name: entity.Organization.Name,
		}
	}

	// Convert related items if loaded
	if len(entity.RelatedItems) > 0 {
		data.RelatedItems = make([]RelatedItemData, len(entity.RelatedItems))
		for i, item := range entity.RelatedItems {
			data.RelatedItems[i] = RelatedItemData{
				ID:          item.ID,
				Name:        item.Name,
				Description: item.Description,
				CreatedAt:   item.CreatedAt,
			}
		}
	}

	// Convert tags if loaded
	if len(entity.Tags) > 0 {
		data.Tags = make([]TagData, len(entity.Tags))
		for i, tag := range entity.Tags {
			data.Tags[i] = TagData{
				ID:    tag.ID,
				Name:  tag.Name,
				Color: tag.Color,
			}
		}
	}

	return data
}

func FromModels(entities []model.ExampleEntity) []EntityData {
	result := make([]EntityData, len(entities))
	for i, entity := range entities {
		result[i] = FromModel(&entity)
	}
	return result
}
```

## Docker Configuration

### Dockerfile (Multi-stage Build)

```dockerfile
# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 -S testlake && \
    adduser -u 1001 -S testlake -G testlake

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/templates ./templates/

# Change ownership to non-root user
RUN chown -R testlake:testlake /app

# Switch to non-root user
USER testlake

# Expose port
EXPOSE 8000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8000/health || exit 1

# Run the application
CMD ["./main"]
```

### Docker Compose Configuration

```yaml
# docker-compose.yml - Complete Development Stack
version: '3.8'

services:
  # PostgreSQL Database
  postgres:
    image: postgres:15-alpine
    container_name: testlake-postgres
    environment:
      POSTGRES_DB: testlake_db
      POSTGRES_USER: testlake_user
      POSTGRES_PASSWORD: secure_password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./sql/init:/docker-entrypoint-initdb.d
    networks:
      - testlake-network
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U testlake_user -d testlake_db"]
      interval: 10s
      timeout: 5s
      retries: 5

  # Redis for caching and sessions
  redis:
    image: redis:7-alpine
    container_name: testlake-redis
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    networks:
      - testlake-network
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 3

  # TestLake API Service
  api:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: testlake-api
    environment:
      # Database
      DB_HOST: postgres
      DB_USER: testlake_user
      DB_PASSWORD: secure_password
      DB_NAME: testlake_db
      DB_PORT: 5432
      
      # Server
      IP: 0.0.0.0
      PORT: 8000
      SCHEME: http
      
      # JWT
      JWT_SECRET: your-super-secret-jwt-key-min-32-chars
      JWT_REFRESH_SECRET: your-refresh-secret-key-min-32-chars
      JWT_EXPIRES_IN: 24h
      JWT_REFRESH_EXPIRES_IN: 168h
      
      # Email (configure with your SMTP provider)
      SMTP_HOST: smtp.gmail.com
      SMTP_PORT: 587
      SMTP_USERNAME: your-email@gmail.com
      SMTP_PASSWORD: your-app-password
      FROM_EMAIL: noreply@testlake.com
      FROM_NAME: TestLake Platform
      
      # PayPal
      PAYPAL_CLIENT_ID: your-paypal-client-id
      PAYPAL_CLIENT_SECRET: your-paypal-client-secret
      PAYPAL_MODE: sandbox
      
      # URLs
      FRONTEND_URL: http://localhost:3000
      API_BASE_URL: http://localhost:8000/api/v1
      
      # Redis
      REDIS_HOST: redis
      REDIS_PORT: 6379
      
    ports:
      - "8000:8000"
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    networks:
      - testlake-network
    volumes:
      - ./uploads:/app/uploads
      - ./logs:/app/logs
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8000/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 30s

  # Nginx Reverse Proxy (optional)
  nginx:
    image: nginx:alpine
    container_name: testlake-nginx
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf
      - ./nginx/ssl:/etc/nginx/ssl
    depends_on:
      - api
    networks:
      - testlake-network

volumes:
  postgres_data:
    driver: local
  redis_data:
    driver: local

networks:
  testlake-network:
    driver: bridge
```

## Development Commands & Scripts

### Makefile for Development

```makefile
# Makefile - Development Commands
.PHONY: help build run test clean docker-up docker-down migrate swagger

# Default target
help:
	@echo "Available commands:"
	@echo "  build      - Build the application"
	@echo "  run        - Run the application"
	@echo "  test       - Run all tests"
	@echo "  test-cov   - Run tests with coverage"
	@echo "  clean      - Clean build artifacts"
	@echo "  docker-up  - Start Docker services"
	@echo "  docker-down - Stop Docker services"
	@echo "  migrate    - Run database migrations"
	@echo "  swagger    - Generate Swagger documentation"
	@echo "  lint       - Run linter"
	@echo "  format     - Format code"

# Build the application
build:
	go build -o bin/testlake .

# Run the application
run:
	go run .

# Run all tests
test:
	go test ./... -v

# Run tests with coverage
test-cov:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Docker commands
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-build:
	docker-compose build

# Database migrations
migrate:
	go run . migrate

# Generate Swagger documentation
swagger:
	swag init

# Linting
lint:
	golangci-lint run

# Format code
format:
	go fmt ./...
	gofumpt -w .

# Install development dependencies
deps:
	go mod download
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install mvdan.cc/gofumpt@latest

# Run development server with hot reload
dev:
	air

# Generate mocks for testing
mocks:
	mockgen -source=dao/interfaces.go -destination=mocks/dao_mocks.go

# Security scan
security:
	gosec ./...

# Performance benchmarks
bench:
	go test ./... -bench=.

# Check dependencies for updates
deps-check:
	go list -u -m all
```

This comprehensive structure documentation provides everything needed to understand, develop, and maintain the TestLake backend, including detailed templates, configuration examples, and development workflows.