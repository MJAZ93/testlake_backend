# CLAUDE.md - Project Reference Guide

This document provides quick references to the key documentation files for this project. Instead of duplicating content, refer to these files for specific information:

## Documentation Files (in ./proj_docs)

### 1. `GO_REST_API_DOCUMENTATION.md`
**Purpose**: Complete Go REST API backend documentation  
**Contains**:
- Go backend architecture (MVC pattern, directory structure)
- Database integration with GORM and PostgreSQL
- Authentication system (JWT middleware)
- Error handling utilities (generic 404 handling, convenience functions)
- Microservices architecture with Docker
- Service/Controller/DAO patterns
- Testing strategies
- Deployment configurations
- Best practices and examples

**Key for**: Backend development, API structure, error handling patterns, microservices setup

### 2. `testlake_spec.md`
**Purpose**: Complete TestLake platform specification  
**Contains**:
- Core business concept and features
- Complete database schema (users, organizations, projects, environments, features, schemas, test data)
- REST API endpoints (authentication, user management, project management, etc.)
- Application screens and user flows
- Business logic and feature relationships
- Simplified onboarding flow ("Zero to Test Data in 5 Minutes")
- Real-world e-commerce example

**Key for**: Understanding business requirements, database design, user workflows, feature specifications

### 3. `testlake_example.md`
**Purpose**: Detailed real-world implementation example  
**Contains**:
- Complete e-commerce platform example
- Step-by-step implementation (project setup → schema creation → feature assignment → test data generation)
- SQL examples with actual data
- API request/response examples
- Team collaboration scenarios
- Scaling considerations

**Key for**: Implementation patterns, real-world usage scenarios, understanding feature relationships

### 4. `testlake_payments.md`
**Purpose**: Billing and payment system documentation  
**Contains**:
- Payment-related database schema updates (plans, subscriptions, invoices, payments)
- PayPal API integration details
- Billing REST API endpoints
- Subscription lifecycle management
- Usage tracking and plan limits enforcement
- Revenue recognition and compliance

**Key for**: Payment processing, subscription management, billing features, PayPal integration

## Quick Context

**This is a two-part project:**

1. **Go REST API Backend**: Generic, reusable backend architecture using Go + Gin + GORM + PostgreSQL
   - Refer to: `GO_REST_API_DOCUMENTATION.md`

2. **TestLake Platform**: Test data management SaaS application built on the Go backend
   - Refer to: `testlake_spec.md` for core features
   - Refer to: `testlake_example.md` for implementation examples  
   - Refer to: `testlake_payments.md` for billing features

## Development Workflow

1. **For backend/API questions**: Check `GO_REST_API_DOCUMENTATION.md`
2. **For business logic/features**: Check `testlake_spec.md`
3. **For implementation examples**: Check `testlake_example.md`
4. **For payment/billing**: Check `testlake_payments.md`

## Key Technologies
- **Backend**: Go 1.24+, Gin, GORM, PostgreSQL, JWT, Docker
- **Payments**: PayPal API integration
- **Architecture**: Microservices-ready with API Gateway pattern
- **Deployment**: Docker Compose, Kubernetes-ready

This approach keeps documentation up-to-date in source files while providing quick navigation for development needs.