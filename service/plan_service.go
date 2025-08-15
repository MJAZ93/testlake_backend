package service

import (
	"testlake/controller"

	"github.com/gin-gonic/gin"
)

type PlanService struct {
	Route      string
	Controller controller.PlanController
}

// GetAllPlans godoc
// @Summary Get all plans
// @Description Get all available subscription plans
// @Tags Plans
// @Accept json
// @Produce json
// @Success 200 {object} plan.PlanListOut
// @Failure 500 {object} inout.BaseResponse
// @Router /api/v1/plans [GET]
func (s PlanService) GetAllPlans(r *gin.RouterGroup, route string) {
	r.GET("/"+s.Route+"/"+route, s.Controller.GetAllPlans)
}

// GetPlan godoc
// @Summary Get plan by ID
// @Description Get specific plan information by ID
// @Tags Plans
// @Accept json
// @Produce json
// @Param id path string true "Plan ID"
// @Success 200 {object} plan.PlanOut
// @Failure 400 {object} inout.BaseResponse
// @Failure 404 {object} inout.BaseResponse
// @Router /api/v1/plans/{id} [GET]
func (s PlanService) GetPlan(r *gin.RouterGroup, route string) {
	r.GET("/"+s.Route+"/"+route+"/:id", s.Controller.GetPlan)
}

// ComparePlans godoc
// @Summary Compare plans
// @Description Get all plans for comparison
// @Tags Plans
// @Accept json
// @Produce json
// @Success 200 {object} plan.PlanComparisonOut
// @Failure 500 {object} inout.BaseResponse
// @Router /api/v1/plans/compare [GET]
func (s PlanService) ComparePlans(r *gin.RouterGroup, route string) {
	r.GET("/"+s.Route+"/"+route, s.Controller.ComparePlans)
}
