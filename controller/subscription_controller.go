package controller

import (
	"errors"
	"net/http"
	"testlake/dao"
	"testlake/inout"
	"testlake/inout/subscription"
	"testlake/model"
	"testlake/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionController struct{}

// GetSubscription returns the current subscription for an organization
func (controller SubscriptionController) GetSubscription(context *gin.Context) {
	orgIDParam := context.Param("orgId")
	organizationID, err := uuid.Parse(orgIDParam)
	if err != nil {
		utils.ReportBadRequest(context, "Invalid organization ID")
		return
	}

	userID, err := utils.ExtractUserID(context)
	if err != nil {
		utils.ReportUnauthorized(context, "Authentication required")
		return
	}

	if !controller.verifyOrganizationAccess(userID, organizationID) {
		utils.ReportForbidden(context, "Access denied")
		return
	}

	subscriptionDao := dao.NewSubscriptionDao()
	sub, err := subscriptionDao.GetByOrganizationID(organizationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ReportNotFound(context, "Subscription not found")
		} else {
			utils.ReportInternalServerError(context, "Database error")
		}
		return
	}

	response := subscription.SubscriptionOut{
		BaseResponse: inout.BaseResponse{
			ErrorCode:        0,
			ErrorDescription: "Success",
		},
		Data: subscription.FromSubscriptionModel(sub),
	}

	context.JSON(http.StatusOK, response)
}

// CreateSubscription creates a new subscription
func (controller SubscriptionController) CreateSubscription(context *gin.Context) {
	orgIDParam := context.Param("id")
	organizationID, err := uuid.Parse(orgIDParam)
	if err != nil {
		utils.ReportBadRequest(context, "Invalid organization ID")
		return
	}

	userID, err := utils.ExtractUserID(context)
	if err != nil {
		utils.ReportUnauthorized(context, "Authentication required")
		return
	}

	if !controller.verifyOrganizationAccess(userID, organizationID) {
		utils.ReportForbidden(context, "Access denied")
		return
	}

	var request subscription.CreateSubscriptionRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		utils.ReportBadRequest(context, "Invalid request body")
		return
	}

	// Validate plan exists
	planDao := dao.NewPlanDao()
	plan, err := planDao.GetByID(request.PlanID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ReportNotFound(context, "Plan not found")
		} else {
			utils.ReportInternalServerError(context, "Database error")
		}
		return
	}

	// Check if organization already has an active subscription
	subscriptionDao := dao.NewSubscriptionDao()
	existingSub, err := subscriptionDao.GetActiveByOrganizationID(organizationID)
	if err == nil && existingSub != nil {
		utils.ReportBadRequest(context, "Organization already has an active subscription")
		return
	}

	// Create subscription - in a real implementation, this would involve PayPal API calls
	now := time.Now()
	var periodEnd time.Time
	if request.BillingCycle == model.BillingCycleMonthly {
		periodEnd = now.AddDate(0, 1, 0)
	} else {
		periodEnd = now.AddDate(1, 0, 0)
	}

	newSubscription := &model.Subscription{
		OrganizationID:       organizationID,
		PlanID:               request.PlanID,
		PayPalSubscriptionID: "TEMP_" + uuid.New().String(), // This would be replaced with actual PayPal subscription ID
		Status:               model.SubscriptionStatusPending,
		BillingCycle:         request.BillingCycle,
		CurrentPeriodStart:   now,
		CurrentPeriodEnd:     periodEnd,
		CreatedBy:            userID,
	}

	err = subscriptionDao.Create(newSubscription)
	if err != nil {
		utils.ReportInternalServerError(context, "Failed to create subscription")
		return
	}

	// Update organization plan information
	orgDao := dao.NewOrganizationDao()
	org, err := orgDao.GetByID(organizationID)
	if err != nil {
		utils.ReportInternalServerError(context, "Failed to get organization")
		return
	}

	org.PlanID = &plan.ID
	org.BillingCycle = request.BillingCycle
	org.SubscriptionStatus = model.OrganizationSubscriptionStatusActive
	org.NextBillingDate = &periodEnd
	err = orgDao.Update(org)
	if err != nil {
		utils.ReportInternalServerError(context, "Failed to update organization")
		return
	}

	response := subscription.SubscriptionOut{
		BaseResponse: inout.BaseResponse{
			ErrorCode:        0,
			ErrorDescription: "Subscription created successfully",
		},
		Data: subscription.FromSubscriptionModel(newSubscription),
	}

	context.JSON(http.StatusCreated, response)
}

// ChangePlan changes the subscription plan
func (controller SubscriptionController) ChangePlan(context *gin.Context) {
	orgIDParam := context.Param("id")
	organizationID, err := uuid.Parse(orgIDParam)
	if err != nil {
		utils.ReportBadRequest(context, "Invalid organization ID")
		return
	}

	userID, err := utils.ExtractUserID(context)
	if err != nil {
		utils.ReportUnauthorized(context, "Authentication required")
		return
	}

	if !controller.verifyOrganizationAccess(userID, organizationID) {
		utils.ReportForbidden(context, "Access denied")
		return
	}

	var request subscription.ChangePlanRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		utils.ReportBadRequest(context, "Invalid request body")
		return
	}

	// Validate new plan exists
	planDao := dao.NewPlanDao()
	newPlan, err := planDao.GetByID(request.NewPlanID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ReportNotFound(context, "Plan not found")
		} else {
			utils.ReportInternalServerError(context, "Database error")
		}
		return
	}

	// Get current subscription
	subscriptionDao := dao.NewSubscriptionDao()
	currentSub, err := subscriptionDao.GetActiveByOrganizationID(organizationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ReportNotFound(context, "No active subscription found")
		} else {
			utils.ReportInternalServerError(context, "Database error")
		}
		return
	}

	// Update subscription
	currentSub.PlanID = request.NewPlanID
	currentSub.BillingCycle = request.BillingCycle
	err = subscriptionDao.Update(currentSub)
	if err != nil {
		utils.ReportInternalServerError(context, "Failed to update subscription")
		return
	}

	// Update organization
	orgDao := dao.NewOrganizationDao()
	org, err := orgDao.GetByID(organizationID)
	if err != nil {
		utils.ReportInternalServerError(context, "Failed to get organization")
		return
	}

	org.PlanID = &newPlan.ID
	org.BillingCycle = request.BillingCycle
	err = orgDao.Update(org)
	if err != nil {
		utils.ReportInternalServerError(context, "Failed to update organization")
		return
	}

	response := subscription.SubscriptionOut{
		BaseResponse: inout.BaseResponse{
			ErrorCode:        0,
			ErrorDescription: "Plan changed successfully",
		},
		Data: subscription.FromSubscriptionModel(currentSub),
	}

	context.JSON(http.StatusOK, response)
}

// CancelSubscription cancels a subscription
func (controller SubscriptionController) CancelSubscription(context *gin.Context) {
	orgIDParam := context.Param("orgId")
	organizationID, err := uuid.Parse(orgIDParam)
	if err != nil {
		utils.ReportBadRequest(context, "Invalid organization ID")
		return
	}

	userID, err := utils.ExtractUserID(context)
	if err != nil {
		utils.ReportUnauthorized(context, "Authentication required")
		return
	}

	if !controller.verifyOrganizationAccess(userID, organizationID) {
		utils.ReportForbidden(context, "Access denied")
		return
	}

	// Get current subscription
	subscriptionDao := dao.NewSubscriptionDao()
	currentSub, err := subscriptionDao.GetActiveByOrganizationID(organizationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ReportNotFound(context, "No active subscription found")
		} else {
			utils.ReportInternalServerError(context, "Database error")
		}
		return
	}

	// Cancel subscription at period end
	err = subscriptionDao.Cancel(currentSub.ID, true)
	if err != nil {
		utils.ReportInternalServerError(context, "Failed to cancel subscription")
		return
	}

	response := inout.BaseResponse{
		ErrorCode:        0,
		ErrorDescription: "Subscription will be cancelled at the end of the current billing period",
	}

	context.JSON(http.StatusOK, response)
}

// ReactivateSubscription reactivates a cancelled subscription
func (controller SubscriptionController) ReactivateSubscription(context *gin.Context) {
	orgIDParam := context.Param("orgId")
	organizationID, err := uuid.Parse(orgIDParam)
	if err != nil {
		utils.ReportBadRequest(context, "Invalid organization ID")
		return
	}

	userID, err := utils.ExtractUserID(context)
	if err != nil {
		utils.ReportUnauthorized(context, "Authentication required")
		return
	}

	if !controller.verifyOrganizationAccess(userID, organizationID) {
		utils.ReportForbidden(context, "Access denied")
		return
	}

	// Get current subscription
	subscriptionDao := dao.NewSubscriptionDao()
	currentSub, err := subscriptionDao.GetByOrganizationID(organizationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ReportNotFound(context, "No subscription found")
		} else {
			utils.ReportInternalServerError(context, "Database error")
		}
		return
	}

	// Reactivate subscription
	currentSub.CancelAtPeriodEnd = false
	currentSub.Status = model.SubscriptionStatusActive
	err = subscriptionDao.Update(currentSub)
	if err != nil {
		utils.ReportInternalServerError(context, "Failed to reactivate subscription")
		return
	}

	response := subscription.SubscriptionOut{
		BaseResponse: inout.BaseResponse{
			ErrorCode:        0,
			ErrorDescription: "Subscription reactivated successfully",
		},
		Data: subscription.FromSubscriptionModel(currentSub),
	}

	context.JSON(http.StatusOK, response)
}

// GetSubscriptionUsage returns current usage metrics for the organization
func (controller SubscriptionController) GetSubscriptionUsage(context *gin.Context) {
	orgIDParam := context.Param("orgId")
	organizationID, err := uuid.Parse(orgIDParam)
	if err != nil {
		utils.ReportBadRequest(context, "Invalid organization ID")
		return
	}

	userID, err := utils.ExtractUserID(context)
	if err != nil {
		utils.ReportUnauthorized(context, "Authentication required")
		return
	}

	if !controller.verifyOrganizationAccess(userID, organizationID) {
		utils.ReportForbidden(context, "Access denied")
		return
	}

	// Get organization and plan limits
	orgDao := dao.NewOrganizationDao()
	org, err := orgDao.GetByID(organizationID)
	if err != nil {
		utils.ReportInternalServerError(context, "Failed to get organization")
		return
	}

	planLimits := subscription.PlanLimits{}
	if org.PlanID != nil {
		planDao := dao.NewPlanDao()
		plan, err := planDao.GetByID(*org.PlanID)
		if err == nil {
			planLimits = subscription.PlanLimits{
				MaxUsers:                plan.MaxUsers,
				MaxProjects:             plan.MaxProjects,
				MaxEnvironments:         plan.MaxEnvironments,
				MaxSchemas:              plan.MaxSchemas,
				MaxTestRecordsPerSchema: plan.MaxTestRecordsPerSchema,
			}
		}
	}

	// Get current usage
	usageDao := dao.NewOrganizationUsageDao()
	currentUsage, err := usageDao.GetCurrentUsage(organizationID)

	usageMetrics := subscription.UsageMetrics{}
	periodStart := time.Now()
	periodEnd := time.Now()

	if err == nil && currentUsage != nil {
		usageMetrics = subscription.UsageMetrics{
			UsersCount:        currentUsage.UsersCount,
			ProjectsCount:     currentUsage.ProjectsCount,
			EnvironmentsCount: currentUsage.EnvironmentsCount,
			SchemasCount:      currentUsage.SchemasCount,
			TestRecordsCount:  currentUsage.TestRecordsCount,
			APIRequestsCount:  currentUsage.APIRequestsCount,
		}
		periodStart = currentUsage.PeriodStart
		periodEnd = currentUsage.PeriodEnd
	}

	usageData := subscription.UsageData{
		CurrentUsage: usageMetrics,
		PlanLimits:   planLimits,
		PeriodStart:  periodStart,
		PeriodEnd:    periodEnd,
	}

	response := subscription.SubscriptionUsageOut{
		BaseResponse: inout.BaseResponse{
			ErrorCode:        0,
			ErrorDescription: "Success",
		},
		Data: usageData,
	}

	context.JSON(http.StatusOK, response)
}

// Helper method to verify organization access
func (controller SubscriptionController) verifyOrganizationAccess(userID, organizationID uuid.UUID) bool {
	orgDao := dao.NewOrganizationDao()
	_, err := orgDao.GetByID(organizationID)
	return err == nil
}
