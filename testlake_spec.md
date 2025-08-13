# TestLake - Test Data Management Platform

## Overview

TestLake is a comprehensive test data management platform designed to eliminate the common developer frustration of "Can you please give me test data to test feature X?". The platform provides structured test data management, feature status tracking, and collaborative tools for development teams.

## Core Concept

TestLake centralizes test data management by providing:
- Structured test data schemas with validation rules
- Feature status tracking and monitoring
- Team collaboration and data sharing
- Automated test data generation
- Multi-platform access (Web, Mobile, API)

---

## Database Schema

### Users Table
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    avatar_url VARCHAR(500),
    auth_provider ENUM('email', 'gmail', 'apple') NOT NULL,
    auth_provider_id VARCHAR(255),
    password_hash VARCHAR(255),
    is_email_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    last_login_at TIMESTAMP,
    status ENUM('active', 'suspended', 'inactive') DEFAULT 'active'
);
```

### Organizations Table
```sql
CREATE TABLE organizations (
    id UUID PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    logo_url VARCHAR(500),
    plan_type ENUM('free', 'starter', 'professional', 'enterprise') DEFAULT 'starter',
    max_users INTEGER DEFAULT 10,
    max_projects INTEGER DEFAULT 5,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    status ENUM('active', 'suspended', 'cancelled') DEFAULT 'active'
);
```

### Projects Table
```sql
CREATE TABLE projects (
    id UUID PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    organization_id UUID REFERENCES organizations(id),
    user_id UUID REFERENCES users(id), -- For personal projects
    is_personal BOOLEAN DEFAULT FALSE,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    status ENUM('active', 'archived') DEFAULT 'active'
);
```

### Teams Table
```sql
CREATE TABLE teams (
    id UUID PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    organization_id UUID NOT NULL REFERENCES organizations(id),
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

### Team Members Table
```sql
CREATE TABLE team_members (
    id UUID PRIMARY KEY,
    team_id UUID NOT NULL REFERENCES teams(id),
    user_id UUID NOT NULL REFERENCES users(id),
    role ENUM('member', 'admin') DEFAULT 'member',
    added_by UUID NOT NULL REFERENCES users(id),
    added_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(team_id, user_id)
);
```

### Project Access Table
```sql
CREATE TABLE project_access (
    id UUID PRIMARY KEY,
    project_id UUID NOT NULL REFERENCES projects(id),
    team_id UUID REFERENCES teams(id),
    user_id UUID REFERENCES users(id),
    permission ENUM('read', 'write', 'admin') DEFAULT 'read',
    granted_by UUID NOT NULL REFERENCES users(id),
    granted_at TIMESTAMP DEFAULT NOW()
);
```

### Environments Table
```sql
CREATE TABLE environments (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL,
    description TEXT,
    color VARCHAR(7) DEFAULT '#3B82F6', -- Hex color for UI identification
    project_id UUID NOT NULL REFERENCES projects(id),
    is_default BOOLEAN DEFAULT FALSE,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    status ENUM('active', 'archived') DEFAULT 'active',
    UNIQUE(project_id, slug)
);
```

### Features Table
```sql
CREATE TABLE features (
    id UUID PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    project_id UUID NOT NULL REFERENCES projects(id),
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

### Feature Environment Status Table
```sql
CREATE TABLE feature_environment_status (
    id UUID PRIMARY KEY,
    feature_id UUID NOT NULL REFERENCES features(id),
    environment_id UUID NOT NULL REFERENCES environments(id),
    is_working BOOLEAN DEFAULT TRUE,
    error_message TEXT,
    last_tested_at TIMESTAMP,
    last_tested_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(feature_id, environment_id)
);
```

### Feature Error Logs Table
```sql
CREATE TABLE feature_error_logs (
    id UUID PRIMARY KEY,
    feature_id UUID NOT NULL REFERENCES features(id),
    environment_id UUID NOT NULL REFERENCES environments(id),
    error_message TEXT NOT NULL,
    error_details JSONB,
    reported_by UUID NOT NULL REFERENCES users(id),
    reported_at TIMESTAMP DEFAULT NOW(),
    resolved_at TIMESTAMP,
    resolved_by UUID REFERENCES users(id)
);
```

### Error Images Table
```sql
CREATE TABLE error_images (
    id UUID PRIMARY KEY,
    error_log_id UUID NOT NULL REFERENCES feature_error_logs(id),
    image_url VARCHAR(500) NOT NULL,
    image_name VARCHAR(200),
    uploaded_at TIMESTAMP DEFAULT NOW()
);
```

### Data Schemas Table
```sql
CREATE TABLE data_schemas (
    id UUID PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    project_id UUID NOT NULL REFERENCES projects(id),
    is_reusable BOOLEAN DEFAULT TRUE,
    schema_definition JSONB NOT NULL, -- JSON structure of the schema
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    status ENUM('active', 'archived') DEFAULT 'active'
);
```

### Feature Schemas Table
```sql
CREATE TABLE feature_schemas (
    id UUID PRIMARY KEY,
    feature_id UUID NOT NULL REFERENCES features(id),
    schema_id UUID NOT NULL REFERENCES data_schemas(id),
    is_primary BOOLEAN DEFAULT FALSE, -- Main schema for this feature
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(feature_id, schema_id)
);
```

### Schema Fields Table
```sql
CREATE TABLE schema_fields (
    id UUID PRIMARY KEY,
    schema_id UUID NOT NULL REFERENCES data_schemas(id),
    field_name VARCHAR(100) NOT NULL,
    field_type ENUM('string', 'number', 'date', 'boolean', 'options', 'reference') NOT NULL,
    is_required BOOLEAN DEFAULT FALSE,
    validation_regex VARCHAR(500),
    min_value VARCHAR(100),
    max_value VARCHAR(100),
    options JSONB, -- For options type
    reference_schema_id UUID REFERENCES data_schemas(id), -- For reference type
    display_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);
```

### Test Data Table
```sql
CREATE TABLE test_data (
    id UUID PRIMARY KEY,
    schema_id UUID NOT NULL REFERENCES data_schemas(id),
    environment_id UUID NOT NULL REFERENCES environments(id),
    data_values JSONB NOT NULL, -- JSON object with field values
    is_used BOOLEAN DEFAULT FALSE,
    used_at TIMESTAMP,
    used_by UUID REFERENCES users(id),
    feature_id UUID REFERENCES features(id), -- Which feature used this data
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    status ENUM('active', 'used', 'invalid') DEFAULT 'active'
);
```

### Test Data Requests Table
```sql
CREATE TABLE test_data_requests (
    id UUID PRIMARY KEY,
    feature_id UUID NOT NULL REFERENCES features(id),
    environment_id UUID NOT NULL REFERENCES environments(id),
    schema_id UUID NOT NULL REFERENCES data_schemas(id),
    requested_by UUID NOT NULL REFERENCES users(id),
    provided_data_id UUID REFERENCES test_data(id),
    request_notes TEXT,
    response_notes TEXT,
    requested_at TIMESTAMP DEFAULT NOW(),
    fulfilled_at TIMESTAMP,
    status ENUM('pending', 'fulfilled', 'rejected') DEFAULT 'pending'
);
```

---

## REST API Endpoints

### Authentication
```
POST   /api/v1/auth/signup
POST   /api/v1/auth/signin
POST   /api/v1/auth/signout
POST   /api/v1/auth/refresh
POST   /api/v1/auth/forgot-password
POST   /api/v1/auth/reset-password
GET    /api/v1/auth/verify-email/{token}
```

### User Management
```
GET    /api/v1/users/profile
PUT    /api/v1/users/profile
DELETE /api/v1/users/account
GET    /api/v1/users/dashboard
GET    /api/v1/users/notifications
PUT    /api/v1/users/notifications/{id}/read
```

### Organization Management
```
POST   /api/v1/organizations
GET    /api/v1/organizations
GET    /api/v1/organizations/{id}
PUT    /api/v1/organizations/{id}
DELETE /api/v1/organizations/{id}
GET    /api/v1/organizations/{id}/members
POST   /api/v1/organizations/{id}/invite
DELETE /api/v1/organizations/{id}/members/{userId}
PUT    /api/v1/organizations/{id}/members/{userId}/role
```

### Team Management
```
POST   /api/v1/organizations/{orgId}/teams
GET    /api/v1/organizations/{orgId}/teams
GET    /api/v1/teams/{id}
PUT    /api/v1/teams/{id}
DELETE /api/v1/teams/{id}
POST   /api/v1/teams/{id}/members
DELETE /api/v1/teams/{id}/members/{userId}
```

### Project Management
```
POST   /api/v1/projects
GET    /api/v1/projects
GET    /api/v1/projects/{id}
PUT    /api/v1/projects/{id}
DELETE /api/v1/projects/{id}
POST   /api/v1/projects/{id}/access
PUT    /api/v1/projects/{id}/access/{accessId}
DELETE /api/v1/projects/{id}/access/{accessId}
```

### Environment Management
```
POST   /api/v1/projects/{projectId}/environments
GET    /api/v1/projects/{projectId}/environments
GET    /api/v1/environments/{id}
PUT    /api/v1/environments/{id}
DELETE /api/v1/environments/{id}
PUT    /api/v1/environments/{id}/set-default
POST   /api/v1/environments/{id}/clone
```

### Feature Management
```
POST   /api/v1/projects/{projectId}/features
GET    /api/v1/projects/{projectId}/features
GET    /api/v1/features/{id}
PUT    /api/v1/features/{id}
DELETE /api/v1/features/{id}
POST   /api/v1/features/{id}/environments/{envId}/error-report
GET    /api/v1/features/{id}/environments/{envId}/error-logs
PUT    /api/v1/features/{id}/environments/{envId}/status
GET    /api/v1/features/{id}/environments
POST   /api/v1/features/{id}/environments/{envId}
DELETE /api/v1/features/{id}/environments/{envId}
```

### Data Schema Management
```
POST   /api/v1/projects/{projectId}/schemas
GET    /api/v1/projects/{projectId}/schemas
GET    /api/v1/schemas/{id}
PUT    /api/v1/schemas/{id}
DELETE /api/v1/schemas/{id}
POST   /api/v1/schemas/{id}/fields
PUT    /api/v1/schemas/{id}/fields/{fieldId}
DELETE /api/v1/schemas/{id}/fields/{fieldId}
GET    /api/v1/schemas/{id}/validate
GET    /api/v1/schemas/{id}/features
```

### Feature Schema Management
```
POST   /api/v1/features/{featureId}/schemas
GET    /api/v1/features/{featureId}/schemas
DELETE /api/v1/features/{featureId}/schemas/{schemaId}
PUT    /api/v1/features/{featureId}/schemas/{schemaId}/set-primary
GET    /api/v1/features/{featureId}/available-schemas
```

### Test Data Management
```
POST   /api/v1/schemas/{schemaId}/environments/{envId}/test-data
GET    /api/v1/schemas/{schemaId}/environments/{envId}/test-data
GET    /api/v1/test-data/{id}
PUT    /api/v1/test-data/{id}
DELETE /api/v1/test-data/{id}
POST   /api/v1/schemas/{schemaId}/environments/{envId}/import-csv
POST   /api/v1/schemas/{schemaId}/environments/{envId}/generate-sample
GET    /api/v1/features/{featureId}/environments/{envId}/request-data
GET    /api/v1/features/{featureId}/environments/{envId}/request-data/{schemaId}
POST   /api/v1/test-data/{id}/mark-used
POST   /api/v1/test-data/{id}/reactivate
POST   /api/v1/test-data/bulk-copy-between-environments
GET    /api/v1/features/{featureId}/test-data-summary
```

### Dashboard & Analytics
```
GET    /api/v1/dashboard/overview
GET    /api/v1/dashboard/features-status
GET    /api/v1/dashboard/features-status/{environmentId}
GET    /api/v1/dashboard/data-usage
GET    /api/v1/dashboard/environment-comparison
GET    /api/v1/analytics/project/{projectId}/summary
GET    /api/v1/analytics/project/{projectId}/environment/{envId}/summary
GET    /api/v1/analytics/schema/{schemaId}/usage
GET    /api/v1/analytics/environment/{envId}/health
```

---

## Application Screens

### Authentication Screens
- **Sign Up** - Email, Google, Apple registration
- **Sign In** - Multiple authentication options
- **Forgot Password** - Password recovery
- **Email Verification** - Account activation

### Dashboard Screens
- **Personal Dashboard** - Overview of personal projects and features
- **Organization Dashboard** - Company-wide metrics and status
- **Feature Status Board** - Real-time feature health monitoring
- **Data Usage Analytics** - Test data consumption patterns

### Organization Management
- **Organization Setup** - Company registration and configuration
- **Team Management** - Create and manage development teams
- **Member Invitation** - Invite and manage team members
- **Billing & Plans** - Subscription management

### Project Management
- **Project List** - View all accessible projects
- **Project Details** - Individual project overview
- **Access Control** - Manage project permissions
- **Project Settings** - Configuration and preferences

### Environment Management
- **Environment Setup** - Create and configure environments (dev, staging, prod, etc.)
- **Environment Dashboard** - Overview of all environments in a project
- **Environment Settings** - Configure environment-specific settings
- **Environment Cloning** - Duplicate environments with data
- **Environment Comparison** - Compare feature status across environments

### Feature Management
- **Feature List** - All features in a project with environment status and associated schemas
- **Feature Details** - Individual feature status across environments with schema relationships
- **Feature Schema Assignment** - Assign and manage schemas for each feature
- **Environment Status Board** - Matrix view of features vs environments
- **Error Reporting** - Report and track feature issues per environment
- **Error Gallery** - Visual error documentation with images per environment
- **Feature Environment Assignment** - Assign features to specific environments

### Data Schema Management
- **Schema Designer** - Visual schema creation tool with feature usage indicators
- **Field Configuration** - Define validation rules and types
- **Schema Preview** - Test and validate schema structure
- **Schema Templates** - Pre-built common schemas
- **Feature Usage View** - See which features use each schema
- **Schema Sharing** - Manage schema relationships across features

### Test Data Management
- **Environment Data Selector** - Choose environment before data operations
- **Data Entry Form** - Manual test data creation per environment
- **CSV Import** - Bulk data import with environment specification
- **Data Browser** - Search and filter test data by environment
- **Data Generator** - Auto-generate test data for specific environments
- **Usage Tracking** - Monitor data consumption per environment
- **Cross-Environment Data Copy** - Copy test data between environments
- **Environment Data Comparison** - Compare data sets across environments

### Request Management
- **Environment-Specific Data Request** - Request test data for feature in specific environment
- **Request Queue** - Manage incoming data requests by environment
- **Request History** - Track fulfillment history across environments
- **Cross-Environment Request Analytics** - Compare request patterns between environments

---

## Business Logic & Features

### User Authentication & Authorization
- **Multi-provider Auth**: Email, Google, Apple sign-in
- **Role-based Access**: Personal workspace vs. organization roles
- **Team Permissions**: Granular project and data access control

### Organization Management
- **Freemium Model**: Personal free workspace, paid organization features
- **Team Collaboration**: Multiple teams per organization
- **Cross-team Access**: Users can participate in multiple teams

### Environment Management
- **Environment Isolation**: Separate test data and feature status per environment
- **Environment Templates**: Pre-configured environment types (Development, Staging, Production, Testing)
- **Environment Cloning**: Duplicate environments with optional data copying
- **Default Environment**: Set primary environment for quick access
- **Environment-Specific Configuration**: Custom settings per environment

### Feature Status Tracking
- **Environment-Based Status**: Track feature health separately per environment
- **Cross-Environment Comparison**: Compare feature status across environments
- **Environment-Specific Error Logs**: Detailed error tracking per environment
- **Historical Tracking**: Feature status over time per environment
- **Environment Health Dashboard**: Overview of all features across environments

### Data Schema System
- **Project-Level Schemas**: Schemas belong to projects for maximum reusability
- **Feature-Schema Relationships**: Many-to-many relationship between features and schemas
- **Primary/Secondary Schema Classification**: Features can designate one primary schema and multiple secondary schemas
- **Type System**: String, Number, Date, Boolean, Options, References
- **Validation Engine**: Regex patterns and custom validation rules
- **Nested Schemas**: Reference other schemas within schemas
- **Reusability Control**: Configure if data can be reused or is single-use
- **Schema Sharing**: Multiple features can use the same schema without duplication

### Feature-Schema Workflow
- **Schema Assignment**: When creating features, assign relevant schemas (e.g., "user_registration" feature uses User + Address schemas)
- **Primary Schema Selection**: Designate the most important schema for each feature
- **Comprehensive Data Requests**: When requesting test data for a feature, get data from all associated schemas
- **Schema Reusability**: Common schemas like "User" can be shared across multiple features (login, registration, profile_update)
- **Contextual Data Management**: View and manage test data in the context of the features that use them

### Test Data Management
- **Environment-Scoped Data**: All test data belongs to specific environments
- **Smart Generation**: Auto-generate data based on schema rules per environment
- **Import/Export**: CSV bulk operations with environment targeting
- **Usage Tracking**: Monitor which data has been used per environment
- **Data Lifecycle**: Active → Used → Reactivated workflow per environment
- **Cross-Environment Operations**: Copy or migrate data between environments
- **Environment Data Isolation**: Ensure data separation between environments

### Request System
- **Environment-Specific Requests**: Request test data for features in specific environments
- **Multi-Schema Data Sets**: Features can request data from all associated schemas
- **Primary Schema Priority**: Primary schema data always included in requests
- **Secondary Schema Options**: Optional secondary schema data for comprehensive testing
- **Intelligent Distribution**: Randomly select unused data from correct environment and schemas
- **Fallback Generation**: Generate on-the-fly when no environment data available for any feature schema
- **Request Tracking**: Full audit trail of data requests per environment and schema
- **Environment Context**: All requests include environment and feature-schema information

---

## Simplified User Flow - Quick Start

### For New Users: "Zero to Test Data in 5 Minutes"

TestLake provides a streamlined onboarding flow that gets users from signup to working test data as quickly as possible, without requiring deep understanding of the full feature set.

#### **Step 1: Quick Project Setup**
- **One-Click Project Creation**: "Create My First Project" button on dashboard
- **Auto-Generated Environment**: System automatically creates a "Development" environment
- **Skip Advanced Settings**: No teams, permissions, or complex configuration required

#### **Step 2: Smart Feature Creation Wizard**
```
┌─ Create Feature Dialog ─┐
│ Feature Name: [add_user  ]                    │
│ Description: [Allow users to register]       │
│                                              │
│ ✅ I need test data for this feature         │
│ ✅ Auto-create basic schema                  │
│                                              │
│ [Create Feature & Setup Data] [Advanced...] │
└────────────────────────────────────────────┘
```

When user clicks "Create Feature & Setup Data":
1. **Creates the feature** in the default environment
2. **Suggests common schema** based on feature name (e.g., "add_user" → suggests User schema)
3. **Auto-creates basic schema** with intelligent field suggestions
4. **Links feature to schema** automatically

#### **Step 3: Intelligent Schema Generation**
```
┌─ Smart Schema Creator ─┐
│ For feature "add_user", we suggest:          │
│                                              │
│ Schema Name: User                            │
│ ┌─────────────────────────────────────────┐ │
│ │ username     [string] ✅                │ │
│ │ email        [string] ✅                │ │
│ │ password     [string] ✅                │ │
│ │ age          [number] ○                 │ │
│ │ phone        [string] ○                 │ │
│ └─────────────────────────────────────────┘ │
│                                              │
│ [✓] Auto-generate 10 sample records         │
│ [✓] Use realistic test data                  │
│                                              │
│ [Create Schema & Generate Data] [Customize]  │
└────────────────────────────────────────────┘
```

System automatically:
- **Analyzes feature name** to suggest relevant fields
- **Provides smart defaults** with validation rules
- **Pre-selects essential fields** (marked with ✅)
- **Offers to generate sample data** immediately

#### **Step 4: Instant Test Data Generation**
```
┌─ Generated Test Data Preview ─┐
│ ✅ Created 10 User records for Development   │
│                                              │
│ Sample Data:                                 │
│ • username: john_doe, email: john@test.com   │
│ • username: jane_smith, email: jane@test.com │
│ • username: bob_wilson, email: bob@test.com  │
│ ...                                          │
│                                              │
│ 🎉 Your feature "add_user" is ready!        │
│                                              │
│ [Get Test Data] [Add More Data] [Dashboard]  │
└────────────────────────────────────────────┘
```

#### **Step 5: One-Click Data Access**
```
┌─ Get Test Data ─┐
│ Feature: add_user                            │
│ Environment: Development                     │
│                                              │
│ 🎯 Here's your test data:                   │
│                                              │
│ {                                            │
│   "username": "alice_brown",                 │
│   "email": "alice@test.com",                 │
│   "password": "securePass123",               │
│   "age": 28,                                 │
│   "phone": "+1-555-0123"                    │
│ }                                            │
│                                              │
│ [📋 Copy JSON] [📧 Copy as Form] [New Data]  │
│ [✓ Mark as Used] [⚠️ Report Issue]          │
└────────────────────────────────────────────┘
```

### **Quick Action Buttons Throughout UI**

#### **Dashboard Quick Actions**
- **"+ New Feature"** → Opens smart feature creation wizard
- **"Get Test Data"** → Shows dropdown of all features for quick data access
- **"Add Sample Data"** → Bulk generate more test data for existing schemas

#### **Feature Card Quick Actions**
```
┌─ Feature: add_user ─┐
│ ✅ Working in Dev    │
│                      │
│ [Get Data] [+ Data]  │
│ [Report Issue]       │
└────────────────────┘
```

#### **Smart Suggestions Engine**
The system provides contextual help throughout:

- **Feature Name Analysis**: "add_car" → suggests Car schema with [brand, model, year, color]
- **Common Patterns**: "login" → suggests User schema with [username, password]
- **Field Type Detection**: "email" field → auto-adds email validation regex
- **Related Data**: Creating "Order" schema → suggests linking to existing "User" schema

### **Progressive Disclosure**
Users start simple but can access advanced features as needed:

1. **Beginner**: Use quick wizard, auto-generated schemas, one environment
2. **Intermediate**: Customize schemas, add multiple environments, manual data entry
3. **Advanced**: Teams, complex validation rules, CSV import, cross-environment data management

### **Mobile Quick Access**
The mobile app focuses on the most common use case:
```
┌─ TestLake Mobile ─┐
│ 🔍 Search features │
│                    │
│ add_user    [Get]  │
│ user_login  [Get]  │
│ checkout    [Get]  │
│                    │
│ [+ Quick Feature]  │
└──────────────────┘
```

**Tap "Get"** → Instantly receive test data
**Tap "+ Quick Feature"** → Voice-to-text feature creation

This simplified flow ensures that developers can get working test data within minutes of signing up, while still having access to all the powerful features of TestLake when they're ready to use them.

---