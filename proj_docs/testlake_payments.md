# TestLake - Payments & Billing System

## Payment-Related Database Schema Updates

### Plans Table
```sql
CREATE TABLE plans (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    price_monthly DECIMAL(10,2) NOT NULL,
    price_yearly DECIMAL(10,2) NOT NULL,
    max_users INTEGER NOT NULL,
    max_projects INTEGER NOT NULL,
    max_environments INTEGER NOT NULL,
    max_schemas INTEGER NOT NULL,
    max_test_records_per_schema INTEGER NOT NULL,
    features JSONB NOT NULL, -- JSON array of enabled features
    paypal_monthly_plan_id VARCHAR(100), -- PayPal subscription plan ID
    paypal_yearly_plan_id VARCHAR(100), -- PayPal subscription plan ID
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Insert default plans
INSERT INTO plans (id, name, slug, description, price_monthly, price_yearly, max_users, max_projects, max_environments, max_schemas, max_test_records_per_schema, features) VALUES
('00000000-0000-0000-0000-000000000001', 'Free', 'free', 'Perfect for personal projects', 0.00, 0.00, 1, 2, 2, 5, 100, '["basic_analytics", "email_support"]'),
('00000000-0000-0000-0000-000000000002', 'Starter', 'starter', 'Great for small teams', 29.00, 290.00, 5, 10, 5, 25, 1000, '["basic_analytics", "email_support", "team_collaboration", "csv_import"]'),
('00000000-0000-0000-0000-000000000003', 'Professional', 'professional', 'For growing development teams', 79.00, 790.00, 25, 50, 15, 100, 10000, '["advanced_analytics", "priority_support", "team_collaboration", "csv_import", "api_access", "custom_validation"]'),
('00000000-0000-0000-0000-000000000004', 'Enterprise', 'enterprise', 'For large organizations', 199.00, 1990.00, 100, 200, 50, 500, 100000, '["advanced_analytics", "priority_support", "team_collaboration", "csv_import", "api_access", "custom_validation", "sso", "audit_logs", "custom_integrations"]');
```

### Update Organizations Table
```sql
-- Add payment-related fields to organizations table
ALTER TABLE organizations 
ADD COLUMN plan_id UUID REFERENCES plans(id),
ADD COLUMN billing_cycle ENUM('monthly', 'yearly') DEFAULT 'monthly',
ADD COLUMN subscription_status ENUM('active', 'past_due', 'cancelled', 'suspended', 'trialing') DEFAULT 'active',
ADD COLUMN trial_ends_at TIMESTAMP,
ADD COLUMN next_billing_date TIMESTAMP,
ADD COLUMN paypal_subscription_id VARCHAR(100),
ADD COLUMN billing_email VARCHAR(255);

-- Update existing organizations to reference free plan
UPDATE organizations SET plan_id = '00000000-0000-0000-0000-000000000001' WHERE plan_id IS NULL;
```

### Payment Methods Table
```sql
CREATE TABLE payment_methods (
    id UUID PRIMARY KEY,
    organization_id UUID NOT NULL REFERENCES organizations(id),
    paypal_payer_id VARCHAR(100),
    paypal_email VARCHAR(255),
    payment_method_type ENUM('paypal') DEFAULT 'paypal',
    is_default BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

### Subscriptions Table
```sql
CREATE TABLE subscriptions (
    id UUID PRIMARY KEY,
    organization_id UUID NOT NULL REFERENCES organizations(id),
    plan_id UUID NOT NULL REFERENCES plans(id),
    paypal_subscription_id VARCHAR(100) UNIQUE NOT NULL,
    status ENUM('active', 'cancelled', 'suspended', 'expired', 'pending') DEFAULT 'pending',
    billing_cycle ENUM('monthly', 'yearly') NOT NULL,
    current_period_start TIMESTAMP NOT NULL,
    current_period_end TIMESTAMP NOT NULL,
    trial_end TIMESTAMP,
    cancel_at_period_end BOOLEAN DEFAULT FALSE,
    cancelled_at TIMESTAMP,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

### Invoices Table
```sql
CREATE TABLE invoices (
    id UUID PRIMARY KEY,
    organization_id UUID NOT NULL REFERENCES organizations(id),
    subscription_id UUID REFERENCES subscriptions(id),
    paypal_invoice_id VARCHAR(100),
    invoice_number VARCHAR(50) UNIQUE NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    tax_amount DECIMAL(10,2) DEFAULT 0,
    total_amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    status ENUM('draft', 'sent', 'paid', 'cancelled', 'refunded') DEFAULT 'draft',
    billing_period_start TIMESTAMP,
    billing_period_end TIMESTAMP,
    due_date TIMESTAMP,
    paid_at TIMESTAMP,
    invoice_url VARCHAR(500),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

### Invoice Line Items Table
```sql
CREATE TABLE invoice_line_items (
    id UUID PRIMARY KEY,
    invoice_id UUID NOT NULL REFERENCES invoices(id),
    description VARCHAR(255) NOT NULL,
    quantity INTEGER DEFAULT 1,
    unit_price DECIMAL(10,2) NOT NULL,
    total_price DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
```

### Payments Table
```sql
CREATE TABLE payments (
    id UUID PRIMARY KEY,
    organization_id UUID NOT NULL REFERENCES organizations(id),
    invoice_id UUID REFERENCES invoices(id),
    subscription_id UUID REFERENCES subscriptions(id),
    paypal_payment_id VARCHAR(100) UNIQUE,
    paypal_payer_id VARCHAR(100),
    amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    payment_method ENUM('paypal') DEFAULT 'paypal',
    status ENUM('pending', 'completed', 'failed', 'cancelled', 'refunded') DEFAULT 'pending',
    failure_reason TEXT,
    processed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

### Usage Tracking Table
```sql
CREATE TABLE organization_usage (
    id UUID PRIMARY KEY,
    organization_id UUID NOT NULL REFERENCES organizations(id),
    period_start TIMESTAMP NOT NULL,
    period_end TIMESTAMP NOT NULL,
    users_count INTEGER DEFAULT 0,
    projects_count INTEGER DEFAULT 0,
    environments_count INTEGER DEFAULT 0,
    schemas_count INTEGER DEFAULT 0,
    test_records_count INTEGER DEFAULT 0,
    api_requests_count INTEGER DEFAULT 0,
    recorded_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(organization_id, period_start, period_end)
);
```

### Billing Events Log Table
```sql
CREATE TABLE billing_events (
    id UUID PRIMARY KEY,
    organization_id UUID NOT NULL REFERENCES organizations(id),
    event_type ENUM('subscription_created', 'subscription_updated', 'subscription_cancelled', 'payment_succeeded', 'payment_failed', 'invoice_created', 'plan_changed') NOT NULL,
    event_data JSONB,
    paypal_event_id VARCHAR(100),
    processed_at TIMESTAMP DEFAULT NOW(),
    created_at TIMESTAMP DEFAULT NOW()
);
```

---

## REST API Endpoints - Billing & Payments

### Plan Management
```
GET    /api/v1/plans
GET    /api/v1/plans/{id}
GET    /api/v1/plans/compare
```

### Subscription Management
```
GET    /api/v1/organizations/{orgId}/subscription
POST   /api/v1/organizations/{orgId}/subscription/create
PUT    /api/v1/organizations/{orgId}/subscription/change-plan
POST   /api/v1/organizations/{orgId}/subscription/cancel
POST   /api/v1/organizations/{orgId}/subscription/reactivate
GET    /api/v1/organizations/{orgId}/subscription/usage
```

### Payment Methods
```
GET    /api/v1/organizations/{orgId}/payment-methods
POST   /api/v1/organizations/{orgId}/payment-methods
PUT    /api/v1/organizations/{orgId}/payment-methods/{id}
DELETE /api/v1/organizations/{orgId}/payment-methods/{id}
PUT    /api/v1/organizations/{orgId}/payment-methods/{id}/set-default
```

### Billing & Invoicing
```
GET    /api/v1/organizations/{orgId}/billing/overview
GET    /api/v1/organizations/{orgId}/invoices
GET    /api/v1/invoices/{id}
GET    /api/v1/invoices/{id}/download
POST   /api/v1/invoices/{id}/pay
GET    /api/v1/organizations/{orgId}/billing/history
```

### PayPal Integration
```
POST   /api/v1/payments/paypal/create-order
POST   /api/v1/payments/paypal/capture-order
POST   /api/v1/payments/paypal/create-subscription
GET    /api/v1/payments/paypal/subscription/{subscriptionId}
POST   /api/v1/payments/paypal/cancel-subscription/{subscriptionId}
POST   /api/v1/payments/paypal/webhooks
```

### Usage & Analytics
```
GET    /api/v1/organizations/{orgId}/usage/current
GET    /api/v1/organizations/{orgId}/usage/history
GET    /api/v1/organizations/{orgId}/limits/check
GET    /api/v1/organizations/{orgId}/billing/forecast
```

---

## PayPal API Integration

### Environment Configuration
```javascript
// PayPal SDK Configuration
const PAYPAL_CONFIG = {
  sandbox: {
    clientId: process.env.PAYPAL_SANDBOX_CLIENT_ID,
    clientSecret: process.env.PAYPAL_SANDBOX_CLIENT_SECRET,
    baseURL: 'https://api-m.sandbox.paypal.com'
  },
  production: {
    clientId: process.env.PAYPAL_LIVE_CLIENT_ID,
    clientSecret: process.env.PAYPAL_LIVE_CLIENT_SECRET,
    baseURL: 'https://api-m.paypal.com'
  }
};
```

### Subscription Creation Flow
```javascript
// 1. Create PayPal Product (One-time setup)
POST /v1/catalogs/products
{
  "name": "TestLake Professional Plan",
  "description": "Monthly subscription to TestLake Professional",
  "type": "SERVICE",
  "category": "SOFTWARE"
}

// 2. Create PayPal Billing Plan
POST /v1/billing/plans
{
  "product_id": "{product_id}",
  "name": "TestLake Professional Monthly",
  "description": "Monthly billing for TestLake Professional plan",
  "billing_cycles": [{
    "frequency": {
      "interval_unit": "MONTH",
      "interval_count": 1
    },
    "tenure_type": "REGULAR",
    "sequence": 1,
    "total_cycles": 0,
    "pricing_scheme": {
      "fixed_price": {
        "value": "79.00",
        "currency_code": "USD"
      }
    }
  }],
  "payment_preferences": {
    "auto_bill_outstanding": true,
    "setup_fee": {
      "value": "0",
      "currency_code": "USD"
    },
    "setup_fee_failure_action": "CONTINUE",
    "payment_failure_threshold": 3
  }
}

// 3. Create Subscription
POST /v1/billing/subscriptions
{
  "plan_id": "{plan_id}",
  "subscriber": {
    "name": {
      "given_name": "John",
      "surname": "Doe"
    },
    "email_address": "john@company.com"
  },
  "application_context": {
    "brand_name": "TestLake",
    "locale": "en-US",
    "user_action": "SUBSCRIBE_NOW",
    "payment_method": {
      "payer_selected": "PAYPAL",
      "payee_preferred": "IMMEDIATE_PAYMENT_REQUIRED"
    },
    "return_url": "https://app.testlake.com/billing/success",
    "cancel_url": "https://app.testlake.com/billing/cancel"
  }
}
```

### Webhook Handling
```javascript
// PayPal Webhook Events to Handle
const WEBHOOK_EVENTS = [
  'BILLING.SUBSCRIPTION.ACTIVATED',
  'BILLING.SUBSCRIPTION.CANCELLED',
  'BILLING.SUBSCRIPTION.SUSPENDED',
  'BILLING.SUBSCRIPTION.PAYMENT.FAILED',
  'PAYMENT.SALE.COMPLETED',
  'PAYMENT.SALE.DENIED',
  'INVOICING.INVOICE.PAID',
  'INVOICING.INVOICE.CANCELLED'
];

// Webhook endpoint
POST /api/v1/payments/paypal/webhooks
{
  "id": "WH-2WR32451HC0233532-67976317FL4543714",
  "create_time": "2023-01-01T00:00:00Z",
  "resource_type": "subscription",
  "event_type": "BILLING.SUBSCRIPTION.ACTIVATED",
  "summary": "A billing subscription was activated.",
  "resource": {
    "id": "I-BW452GLLEP1G",
    "status": "ACTIVE",
    "status_update_time": "2023-01-01T00:00:00Z"
  }
}
```

---

## Application Screens - Billing

### Plan Selection & Pricing
```
┌─ Choose Your Plan ─┐
│                                              │
│ ○ Free        ○ Starter      ● Professional │
│ $0/month      $29/month      $79/month       │
│                                              │
│ • 1 user      • 5 users      • 25 users     │
│ • 2 projects  • 10 projects  • 50 projects  │
│ • Basic       • CSV Import   • API Access   │
│                                              │
│ [Current Plan] [Get Started] [Upgrade Now]   │
│                                              │
│ 💡 Save 17% with yearly billing             │
│ ○ Monthly    ● Yearly                       │
└────────────────────────────────────────────┘
```

### PayPal Checkout Flow
```
┌─ Complete Your Subscription ─┐
│                                              │
│ Plan: Professional (Yearly)                 │
│ Price: $790.00/year (Save $158)             │
│                                              │
│ Billing Details:                            │
│ Company: [Acme Corp            ]             │
│ Email:   [billing@acme.com     ]             │
│                                              │
│ Payment Method:                              │
│ [🅿️ Pay with PayPal] [💳 Credit Card]       │
│                                              │
│ [⚡ Complete Subscription]                  │
│                                              │
│ 🔒 Secure checkout powered by PayPal        │
└────────────────────────────────────────────┘
```

### Billing Dashboard
```
┌─ Billing & Usage ─┐
│                                              │
│ Current Plan: Professional                   │
│ Status: Active                               │
│ Next billing: March 15, 2024 ($79.00)       │
│                                              │
│ Usage This Month:                            │
│ ████████░░ Users (18/25)                    │
│ ██████░░░░ Projects (34/50)                 │
│ ████░░░░░░ API Calls (4.2k/10k)            │
│                                              │
│ Recent Invoices:                             │
│ • Feb 2024 - $79.00 [Paid] [📥 Download]    │
│ • Jan 2024 - $79.00 [Paid] [📥 Download]    │
│                                              │
│ [Change Plan] [Update Payment] [View Usage] │
└────────────────────────────────────────────┘
```

### Payment Method Management
```
┌─ Payment Methods ─┐
│                                              │
│ Primary Payment Method:                      │
│ 🅿️ PayPal (john@company.com) [✓ Default]    │
│ Added: Jan 15, 2024                         │
│ [Edit] [Remove]                             │
│                                              │
│ [+ Add PayPal Account]                      │
│ [+ Add Credit Card]                         │
│                                              │
│ Backup Payment:                             │
│ ○ Use company credit card as backup         │
│ ● Suspend service if primary fails          │
│                                              │
│ [Save Changes]                              │
└────────────────────────────────────────────┘
```

### Usage Analytics
```
┌─ Usage Analytics ─┐
│                                              │
│ Current Period: Feb 1 - Feb 29, 2024        │
│                                              │
│ Resource Usage:                              │
│ ┌─────────────────────────────────────────┐ │
│ │ Users      ████████░░ 18/25 (72%)      │ │
│ │ Projects   ██████░░░░ 34/50 (68%)      │ │
│ │ Schemas    ████░░░░░░ 42/100 (42%)     │ │
│ │ Test Data  ███░░░░░░░ 3.2k/10k (32%)   │ │
│ │ API Calls  ██░░░░░░░░ 2.1k/5k (42%)    │ │
│ └─────────────────────────────────────────┘ │
│                                              │
│ Usage Trends (Last 6 months):               │
│ [📊 Interactive Chart]                      │
│                                              │
│ [Export Report] [Set Alerts] [Upgrade]      │
└────────────────────────────────────────────┘
```

### Invoice Details
```
┌─ Invoice #INV-2024-0234 ─┐
│                                              │
│ Invoice Date: February 15, 2024             │
│ Due Date: March 1, 2024                     │
│ Status: Paid                                │
│                                              │
│ Bill To:                                     │
│ Acme Corporation                             │
│ billing@acme.com                             │
│                                              │
│ Description              Qty    Amount       │
│ Professional Plan        1      $79.00       │
│ (Feb 15 - Mar 15, 2024)                     │
│                                              │
│ Subtotal:                       $79.00       │
│ Tax:                            $0.00        │
│ Total:                          $79.00       │
│                                              │
│ Payment Method: PayPal                       │
│ Paid: February 15, 2024                     │
│                                              │
│ [📥 Download PDF] [✉️ Email Invoice]         │
└────────────────────────────────────────────┘
```

---

## Business Logic - Billing Features

### Plan Limits Enforcement
```javascript
// Middleware to check plan limits
const checkPlanLimits = async (organizationId, resource) => {
  const org = await Organization.findById(organizationId).include('plan');
  const usage = await getCurrentUsage(organizationId);
  
  const limits = {
    users: org.plan.max_users,
    projects: org.plan.max_projects,
    environments: org.plan.max_environments,
    schemas: org.plan.max_schemas,
    test_records: org.plan.max_test_records_per_schema
  };
  
  if (usage[resource] >= limits[resource]) {
    throw new Error(`Plan limit exceeded for ${resource}. Upgrade to access more.`);
  }
};
```

### Subscription Lifecycle Management
1. **Trial Period**: 14-day free trial for paid plans
2. **Automatic Billing**: PayPal handles recurring payments
3. **Failed Payment Handling**: 
   - Retry payment 3 times over 7 days
   - Suspend service after failed retries
   - Send notification emails
4. **Plan Upgrades**: Immediate access, prorated billing
5. **Plan Downgrades**: Take effect at next billing cycle
6. **Cancellation**: Access continues until period end

### Usage Tracking & Alerts
- **Real-time Monitoring**: Track resource usage across organization
- **Usage Alerts**: Notify when approaching plan limits (80%, 90%, 95%)
- **Overage Protection**: Soft limits with upgrade prompts
- **Historical Analytics**: Track usage trends and growth patterns

### PayPal Integration Features
- **Subscription Management**: Create, modify, cancel subscriptions
- **Payment Processing**: Handle one-time and recurring payments
- **Webhook Processing**: Real-time status updates from PayPal
- **Refund Handling**: Process refunds through PayPal API
- **Tax Calculation**: Integrate with PayPal's tax services
- **Multi-currency Support**: Support international customers

### Billing Security & Compliance
- **PCI Compliance**: PayPal handles payment card data
- **Data Encryption**: Encrypt sensitive billing information
- **Audit Logging**: Track all billing-related actions
- **Fraud Prevention**: Leverage PayPal's fraud detection
- **GDPR Compliance**: Handle billing data according to privacy laws

### Revenue Recognition
- **Subscription Revenue**: Recognize monthly/yearly subscription fees
- **Usage-based Billing**: Track overage charges if applicable
- **Proration Logic**: Handle mid-cycle plan changes
- **Tax Handling**: Collect and remit taxes where required
- **Financial Reporting**: Generate revenue reports for accounting

This comprehensive payment system integrates seamlessly with TestLake's existing architecture while providing robust billing capabilities through PayPal's proven payment infrastructure.