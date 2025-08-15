package subscription

import (
	"testlake/inout"
	"testlake/model"
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID                   uuid.UUID                `json:"id"`
	OrganizationID       uuid.UUID                `json:"organization_id"`
	PlanID               uuid.UUID                `json:"plan_id"`
	PayPalSubscriptionID string                   `json:"paypal_subscription_id"`
	Status               model.SubscriptionStatus `json:"status"`
	BillingCycle         model.BillingCycle       `json:"billing_cycle"`
	CurrentPeriodStart   time.Time                `json:"current_period_start"`
	CurrentPeriodEnd     time.Time                `json:"current_period_end"`
	TrialEnd             *time.Time               `json:"trial_end"`
	CancelAtPeriodEnd    bool                     `json:"cancel_at_period_end"`
	CancelledAt          *time.Time               `json:"cancelled_at"`
	CreatedBy            uuid.UUID                `json:"created_by"`
	CreatedAt            time.Time                `json:"created_at"`
	UpdatedAt            time.Time                `json:"updated_at"`
}

type SubscriptionOut struct {
	inout.BaseResponse
	Data Subscription `json:"data"`
}

type UsageMetrics struct {
	UsersCount        int `json:"users_count"`
	ProjectsCount     int `json:"projects_count"`
	EnvironmentsCount int `json:"environments_count"`
	SchemasCount      int `json:"schemas_count"`
	TestRecordsCount  int `json:"test_records_count"`
	APIRequestsCount  int `json:"api_requests_count"`
}

type PlanLimits struct {
	MaxUsers                int `json:"max_users"`
	MaxProjects             int `json:"max_projects"`
	MaxEnvironments         int `json:"max_environments"`
	MaxSchemas              int `json:"max_schemas"`
	MaxTestRecordsPerSchema int `json:"max_test_records_per_schema"`
}

type UsageData struct {
	CurrentUsage UsageMetrics `json:"current_usage"`
	PlanLimits   PlanLimits   `json:"plan_limits"`
	PeriodStart  time.Time    `json:"period_start"`
	PeriodEnd    time.Time    `json:"period_end"`
}

type SubscriptionUsageOut struct {
	inout.BaseResponse
	Data UsageData `json:"data"`
}

func FromSubscriptionModel(s *model.Subscription) Subscription {
	return Subscription{
		ID:                   s.ID,
		OrganizationID:       s.OrganizationID,
		PlanID:               s.PlanID,
		PayPalSubscriptionID: s.PayPalSubscriptionID,
		Status:               s.Status,
		BillingCycle:         s.BillingCycle,
		CurrentPeriodStart:   s.CurrentPeriodStart,
		CurrentPeriodEnd:     s.CurrentPeriodEnd,
		TrialEnd:             s.TrialEnd,
		CancelAtPeriodEnd:    s.CancelAtPeriodEnd,
		CancelledAt:          s.CancelledAt,
		CreatedBy:            s.CreatedBy,
		CreatedAt:            s.CreatedAt,
		UpdatedAt:            s.UpdatedAt,
	}
}
