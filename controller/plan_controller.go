package controller

import (
	"errors"
	"net/http"
	"strconv"
	"testlake/dao"
	"testlake/inout"
	"testlake/inout/plan"
	"testlake/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PlanController struct{}

// GetAllPlans returns all active plans
func (controller PlanController) GetAllPlans(context *gin.Context) {
	planDao := dao.NewPlanDao()
	plans, err := planDao.GetAll()
	if err != nil {
		utils.ReportInternalServerError(context, "Database error")
		return
	}

	response := plan.PlanListOut{
		BaseResponse: inout.BaseResponse{
			ErrorCode:        0,
			ErrorDescription: "Success",
		},
		List: plan.FromPlanModelList(plans),
	}

	context.JSON(http.StatusOK, response)
}

// GetPlan returns a specific plan by ID
func (controller PlanController) GetPlan(context *gin.Context) {
	idParam := context.Param("id")
	planID, err := uuid.Parse(idParam)
	if err != nil {
		utils.ReportBadRequest(context, "Invalid plan ID")
		return
	}

	planDao := dao.NewPlanDao()
	foundPlan, err := planDao.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ReportNotFound(context, "Plan not found")
		} else {
			utils.ReportInternalServerError(context, "Database error")
		}
		return
	}

	response := plan.PlanOut{
		BaseResponse: inout.BaseResponse{
			ErrorCode:        0,
			ErrorDescription: "Success",
		},
		Data: plan.FromPlanModel(foundPlan),
	}

	context.JSON(http.StatusOK, response)
}

// ComparePlans returns all active plans for comparison
func (controller PlanController) ComparePlans(context *gin.Context) {
	planDao := dao.NewPlanDao()
	plans, err := planDao.GetAll()
	if err != nil {
		utils.ReportInternalServerError(context, "Database error")
		return
	}

	comparison := plan.PlanComparison{
		Plans: plan.FromPlanModelList(plans),
	}

	response := plan.PlanComparisonOut{
		BaseResponse: inout.BaseResponse{
			ErrorCode:        0,
			ErrorDescription: "Success",
		},
		Data: comparison,
	}

	context.JSON(http.StatusOK, response)
}

// GetPlansWithPagination returns plans with pagination
func (controller PlanController) GetPlansWithPagination(context *gin.Context) {
	pageParam := context.DefaultQuery("page", "0")
	page, err := strconv.Atoi(pageParam)
	if err != nil || page < 0 {
		utils.ReportBadRequest(context, "Invalid page parameter")
		return
	}

	planDao := dao.NewPlanDao()
	plans, total, err := planDao.GetAllWithPagination(page)
	if err != nil {
		utils.ReportInternalServerError(context, "Database error")
		return
	}

	totalPages := int(total) / planDao.Limit
	if int(total)%planDao.Limit > 0 {
		totalPages++
	}

	response := struct {
		inout.BaseResponse
		List []plan.Plan          `json:"list"`
		Meta inout.PaginationMeta `json:"meta"`
	}{
		BaseResponse: inout.BaseResponse{
			ErrorCode:        0,
			ErrorDescription: "Success",
		},
		List: plan.FromPlanModelList(plans),
		Meta: inout.PaginationMeta{
			Page:       page,
			Limit:      planDao.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}

	context.JSON(http.StatusOK, response)
}
