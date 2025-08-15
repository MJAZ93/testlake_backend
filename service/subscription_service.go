package service

import (
	"testlake/controller"

	"github.com/gin-gonic/gin"
)

type SubscriptionService struct {
	Route      string
	Controller controller.SubscriptionController
}

// GetSubscription godoc
// @Summary Get organization subscription
// @Description Get current subscription for an organization
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token" format(Bearer {token})
// @Param id path string true "Organization ID"
// @Success 200 {object} subscription.SubscriptionOut
// @Failure 400 {object} inout.BaseResponse
// @Failure 401 {object} inout.BaseResponse
// @Failure 403 {object} inout.BaseResponse
// @Failure 404 {object} inout.BaseResponse
// @Router /api/v1/organizations/{id}/subscription [GET]
func (s SubscriptionService) GetSubscription(r *gin.RouterGroup, route string) {
	r.GET("/"+s.Route+"/"+route, s.Controller.GetSubscription)
}

// CreateSubscription godoc
// @Summary Create subscription
// @Description Create a new subscription for an organization
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token" format(Bearer {token})
// @Param id path string true "Organization ID"
// @Param subscription body subscription.CreateSubscriptionRequest true "Subscription data"
// @Success 201 {object} subscription.SubscriptionOut
// @Failure 400 {object} inout.BaseResponse
// @Failure 401 {object} inout.BaseResponse
// @Failure 403 {object} inout.BaseResponse
// @Router /api/v1/organizations/{id}/subscription/create [POST]
func (s SubscriptionService) CreateSubscription(r *gin.RouterGroup, route string) {
	r.POST("/"+s.Route+"/"+route, s.Controller.CreateSubscription)
}

// ChangePlan godoc
// @Summary Change subscription plan
// @Description Change the plan for an existing subscription
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token" format(Bearer {token})
// @Param id path string true "Organization ID"
// @Param plan body subscription.ChangePlanRequest true "Plan change data"
// @Success 200 {object} subscription.SubscriptionOut
// @Failure 400 {object} inout.BaseResponse
// @Failure 401 {object} inout.BaseResponse
// @Failure 403 {object} inout.BaseResponse
// @Failure 404 {object} inout.BaseResponse
// @Router /api/v1/organizations/{id}/subscription/change-plan [PUT]
func (s SubscriptionService) ChangePlan(r *gin.RouterGroup, route string) {
	r.PUT("/"+s.Route+"/"+route, s.Controller.ChangePlan)
}

// CancelSubscription godoc
// @Summary Cancel subscription
// @Description Cancel an organization's subscription
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token" format(Bearer {token})
// @Param id path string true "Organization ID"
// @Success 200 {object} inout.BaseResponse
// @Failure 400 {object} inout.BaseResponse
// @Failure 401 {object} inout.BaseResponse
// @Failure 403 {object} inout.BaseResponse
// @Failure 404 {object} inout.BaseResponse
// @Router /api/v1/organizations/{id}/subscription/cancel [POST]
func (s SubscriptionService) CancelSubscription(r *gin.RouterGroup, route string) {
	r.POST("/"+s.Route+"/"+route, s.Controller.CancelSubscription)
}

// ReactivateSubscription godoc
// @Summary Reactivate subscription
// @Description Reactivate a cancelled subscription
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token" format(Bearer {token})
// @Param id path string true "Organization ID"
// @Success 200 {object} subscription.SubscriptionOut
// @Failure 400 {object} inout.BaseResponse
// @Failure 401 {object} inout.BaseResponse
// @Failure 403 {object} inout.BaseResponse
// @Failure 404 {object} inout.BaseResponse
// @Router /api/v1/organizations/{id}/subscription/reactivate [POST]
func (s SubscriptionService) ReactivateSubscription(r *gin.RouterGroup, route string) {
	r.POST("/"+s.Route+"/"+route, s.Controller.ReactivateSubscription)
}

// GetSubscriptionUsage godoc
// @Summary Get subscription usage
// @Description Get current usage metrics for an organization's subscription
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token" format(Bearer {token})
// @Param id path string true "Organization ID"
// @Success 200 {object} subscription.SubscriptionUsageOut
// @Failure 400 {object} inout.BaseResponse
// @Failure 401 {object} inout.BaseResponse
// @Failure 403 {object} inout.BaseResponse
// @Router /api/v1/organizations/{id}/subscription/usage [GET]
func (s SubscriptionService) GetSubscriptionUsage(r *gin.RouterGroup, route string) {
	r.GET("/"+s.Route+"/"+route, s.Controller.GetSubscriptionUsage)
}
