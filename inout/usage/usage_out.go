package usage

import (
	"testlake/inout"
	"testlake/inout/subscription"
	"testlake/model"
	"time"

	"github.com/google/uuid"
)

type UsageDataPoint struct {
	Date              time.Time `json:"date"`
	UsersCount        int       `json:"users_count"`
	ProjectsCount     int       `json:"projects_count"`
	EnvironmentsCount int       `json:"environments_count"`
	SchemasCount      int       `json:"schemas_count"`
	TestRecordsCount  int       `json:"test_records_count"`
	APIRequestsCount  int       `json:"api_requests_count"`
}

type CurrentUsage struct {
	OrganizationID uuid.UUID                 `json:"organization_id"`
	PeriodStart    time.Time                 `json:"period_start"`
	PeriodEnd      time.Time                 `json:"period_end"`
	Current        subscription.UsageMetrics `json:"current"`
	Limits         subscription.PlanLimits   `json:"limits"`
	UtilizationPct map[string]float64        `json:"utilization_pct"`
	LastUpdated    time.Time                 `json:"last_updated"`
}

type CurrentUsageOut struct {
	inout.BaseResponse
	Data CurrentUsage `json:"data"`
}

type UsageHistory struct {
	OrganizationID uuid.UUID        `json:"organization_id"`
	Period         string           `json:"period"`
	DataPoints     []UsageDataPoint `json:"data_points"`
}

type UsageHistoryOut struct {
	inout.BaseResponse
	Data UsageHistory `json:"data"`
}

type LimitCheck struct {
	Resource       string  `json:"resource"`
	Current        int     `json:"current"`
	Limit          int     `json:"limit"`
	UtilizationPct float64 `json:"utilization_pct"`
	Status         string  `json:"status"` // "ok", "warning", "critical", "exceeded"
}

type LimitsCheck struct {
	OrganizationID uuid.UUID    `json:"organization_id"`
	CheckedAt      time.Time    `json:"checked_at"`
	OverallStatus  string       `json:"overall_status"`
	Limits         []LimitCheck `json:"limits"`
}

type LimitsCheckOut struct {
	inout.BaseResponse
	Data LimitsCheck `json:"data"`
}

type BillingForecast struct {
	OrganizationID      uuid.UUID `json:"organization_id"`
	CurrentPeriodAmount float64   `json:"current_period_amount"`
	NextPeriodAmount    float64   `json:"next_period_amount"`
	ProjectedAmount     float64   `json:"projected_amount"`
	Currency            string    `json:"currency"`
	ForecastedAt        time.Time `json:"forecasted_at"`
	BillingCycle        string    `json:"billing_cycle"`
}

type BillingForecastOut struct {
	inout.BaseResponse
	Data BillingForecast `json:"data"`
}

func FromUsageModel(usage *model.OrganizationUsage, limits *model.Plan) CurrentUsage {
	utilizationPct := make(map[string]float64)

	if limits.MaxUsers > 0 {
		utilizationPct["users"] = float64(usage.UsersCount) / float64(limits.MaxUsers) * 100
	}
	if limits.MaxProjects > 0 {
		utilizationPct["projects"] = float64(usage.ProjectsCount) / float64(limits.MaxProjects) * 100
	}
	if limits.MaxEnvironments > 0 {
		utilizationPct["environments"] = float64(usage.EnvironmentsCount) / float64(limits.MaxEnvironments) * 100
	}
	if limits.MaxSchemas > 0 {
		utilizationPct["schemas"] = float64(usage.SchemasCount) / float64(limits.MaxSchemas) * 100
	}
	if limits.MaxTestRecordsPerSchema > 0 {
		utilizationPct["test_records"] = float64(usage.TestRecordsCount) / float64(limits.MaxTestRecordsPerSchema) * 100
	}

	return CurrentUsage{
		OrganizationID: usage.OrganizationID,
		PeriodStart:    usage.PeriodStart,
		PeriodEnd:      usage.PeriodEnd,
		Current: subscription.UsageMetrics{
			UsersCount:        usage.UsersCount,
			ProjectsCount:     usage.ProjectsCount,
			EnvironmentsCount: usage.EnvironmentsCount,
			SchemasCount:      usage.SchemasCount,
			TestRecordsCount:  usage.TestRecordsCount,
			APIRequestsCount:  usage.APIRequestsCount,
		},
		Limits: subscription.PlanLimits{
			MaxUsers:                limits.MaxUsers,
			MaxProjects:             limits.MaxProjects,
			MaxEnvironments:         limits.MaxEnvironments,
			MaxSchemas:              limits.MaxSchemas,
			MaxTestRecordsPerSchema: limits.MaxTestRecordsPerSchema,
		},
		UtilizationPct: utilizationPct,
		LastUpdated:    usage.RecordedAt,
	}
}

func FromUsageModelList(usages []model.OrganizationUsage) []UsageDataPoint {
	result := make([]UsageDataPoint, len(usages))
	for i, usage := range usages {
		result[i] = UsageDataPoint{
			Date:              usage.RecordedAt,
			UsersCount:        usage.UsersCount,
			ProjectsCount:     usage.ProjectsCount,
			EnvironmentsCount: usage.EnvironmentsCount,
			SchemasCount:      usage.SchemasCount,
			TestRecordsCount:  usage.TestRecordsCount,
			APIRequestsCount:  usage.APIRequestsCount,
		}
	}
	return result
}
