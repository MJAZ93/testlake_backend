package subscription

import (
	"testlake/model"

	"github.com/google/uuid"
)

type CreateSubscriptionRequest struct {
	PlanID       uuid.UUID          `json:"plan_id" binding:"required"`
	BillingCycle model.BillingCycle `json:"billing_cycle" binding:"required"`
}

type ChangePlanRequest struct {
	NewPlanID    uuid.UUID          `json:"new_plan_id" binding:"required"`
	BillingCycle model.BillingCycle `json:"billing_cycle" binding:"required"`
}
