# TestLake Go REST API

A Go-based REST API for the TestLake test data management platform, implementing user management functionality.

## Architecture

This project follows the layered architecture pattern from the Go REST API documentation:

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

## Features

### User Management
- User registration with email, Google, and Apple authentication
- JWT-based authentication system
- User profile management
- Account status management
- Password hashing with bcrypt

### API Endpoints

#### Public Endpoints (No Authentication Required)
- `POST /api/v1/public/user/create` - Create new user
- `POST /api/v1/public/user/login` - User login

#### Private Endpoints (JWT Authentication Required)
- `GET /api/v1/private/user/details/{id}` - Get user by ID
- `GET /api/v1/private/user/list` - List users (paginated)
- `PUT /api/v1/private/user/update/{id}` - Update user profile
- `DELETE /api/v1/private/user/delete/{id}` - Delete user account
- `PUT /api/v1/private/user/status/{id}` - Update user status

## Quick Start

### Prerequisites
- Go 1.24 or higher
- PostgreSQL database
- Environment variables configured

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd testlake
```

2. Install dependencies:
```bash
go mod tidy
```

3. Set up environment variables:
```bash
cp .env.example .env
# Edit .env with your configuration
```

4. Run the application:
```bash
go run main.go
```

### Using Docker

1. Start with Docker Compose:
```bash
docker-compose up --build
```

This will start both the API server and PostgreSQL database.

### Environment Configuration

The `.env` file should contain:

```env
# Server Configuration
IP=localhost
PORT=8000
SCHEME=http

# Database Configuration
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=testlake
DB_PORT=5432

# JWT Configuration
TOKEN_TTL=2000
JWT_PRIVATE_KEY=your_secret_key

# Logging
LOG_PATH=./logs
```

## API Documentation

Once the server is running, you can access the Swagger documentation at:
```
http://localhost:8000/swagger/index.html
```

## Testing

### Run Unit Tests
```bash
# Run DAO tests
go test ./dao_test/...

# Run Controller tests
go test ./controller_test/...

# Run Service integration tests (requires server to be running)
go test ./service_test/...

# Run all tests
go test ./...
```

### Test Coverage
```bash
go test -cover ./...
```

## Database Schema

The user model includes the following fields:
- `id` - UUID primary key
- `email` - Unique email address
- `username` - Unique username
- `first_name` - Optional first name
- `last_name` - Optional last name
- `avatar_url` - Optional profile picture URL
- `auth_provider` - Authentication method (email, gmail, apple)
- `auth_provider_id` - External provider ID
- `password_hash` - Bcrypt hashed password
- `is_email_verified` - Email verification status
- `created_at` - Creation timestamp
- `updated_at` - Last update timestamp
- `last_login_at` - Last login timestamp
- `status` - User status (active, suspended, inactive)

## Example Usage

### Create a New User
```bash
curl -X POST http://localhost:8000/api/v1/public/user/create \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "username": "newuser",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe",
    "auth_provider": "email"
  }'
```

### Login User
```bash
curl -X POST http://localhost:8000/api/v1/public/user/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

### Get User Profile (Authenticated)
```bash
curl -X GET http://localhost:8000/api/v1/private/user/details/{user-id} \
  -H "Authorization: Bearer {jwt-token}"
```

## Development

### Project Structure Principles
- **Separation of Concerns**: Each layer has a specific responsibility
- **Dependency Injection**: Controllers depend on DAOs, not vice versa
- **Error Handling**: Centralized error handling with consistent response format
- **Security**: JWT authentication, password hashing, input validation

### Adding New Features
1. Create model in `model/`
2. Implement DAO in `dao/`
3. Create DTOs in `inout/`
4. Implement controller in `controller/`
5. Create service in `service/`
6. Add routes in `app/routes.go`
7. Write tests in respective test directories

## Contributing

1. Follow the existing architecture patterns
2. Write tests for all new functionality
3. Use consistent error handling
4. Document API endpoints with Swagger annotations
5. Validate all inputs
6. Never expose sensitive information in responses

## License

This project is licensed under the MIT License.