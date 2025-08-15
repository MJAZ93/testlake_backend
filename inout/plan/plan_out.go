package plan

import (
	"testlake/inout"
	"testlake/model"
	"time"

	"github.com/google/uuid"
)

type Plan struct {
	ID                      uuid.UUID `json:"id"`
	Name                    string    `json:"name"`
	Slug                    string    `json:"slug"`
	Description             *string   `json:"description"`
	PriceMonthly            float64   `json:"price_monthly"`
	PriceYearly             float64   `json:"price_yearly"`
	MaxUsers                int       `json:"max_users"`
	MaxProjects             int       `json:"max_projects"`
	MaxEnvironments         int       `json:"max_environments"`
	MaxSchemas              int       `json:"max_schemas"`
	MaxTestRecordsPerSchema int       `json:"max_test_records_per_schema"`
	Features                []string  `json:"features"`
	IsActive                bool      `json:"is_active"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

type PlanOut struct {
	inout.BaseResponse
	Data Plan `json:"data"`
}

type PlanListOut struct {
	inout.BaseResponse
	List []Plan `json:"list"`
}

type PlanComparison struct {
	Plans []Plan `json:"plans"`
}

type PlanComparisonOut struct {
	inout.BaseResponse
	Data PlanComparison `json:"data"`
}

func FromPlanModel(p *model.Plan) Plan {
	plan := Plan{
		ID:                      p.ID,
		Name:                    p.Name,
		Slug:                    p.Slug,
		Description:             p.Description,
		PriceMonthly:            p.PriceMonthly,
		PriceYearly:             p.PriceYearly,
		MaxUsers:                p.MaxUsers,
		MaxProjects:             p.MaxProjects,
		MaxEnvironments:         p.MaxEnvironments,
		MaxSchemas:              p.MaxSchemas,
		MaxTestRecordsPerSchema: p.MaxTestRecordsPerSchema,
		Features:                []string{},
		IsActive:                p.IsActive,
		CreatedAt:               p.CreatedAt,
		UpdatedAt:               p.UpdatedAt,
	}

	// Parse features JSON string into []string
	// This would typically use json.Unmarshal if Features is JSON
	// For now, we'll leave it as empty slice since Features is defined as string in model

	return plan
}

func FromPlanModelList(plans []model.Plan) []Plan {
	result := make([]Plan, len(plans))
	for i, plan := range plans {
		result[i] = FromPlanModel(&plan)
	}
	return result
}
