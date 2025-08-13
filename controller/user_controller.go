package controller

import (
	"errors"
	"net/http"

	"testlake/dao"
	"testlake/inout"
	"testlake/inout/user"
	"testlake/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserController struct{}

// GetProfile returns current user's profile
func (controller UserController) GetProfile(context *gin.Context) {
	userID, err := utils.ExtractUserID(context)
	if err != nil {
		utils.ReportUnauthorized(context, "Authentication required")
		return
	}

	userDao := dao.NewUserDao()
	foundUser, err := userDao.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ReportNotFound(context, "User not found")
		} else {
			utils.ReportInternalServerError(context, "Database error")
		}
		return
	}

	response := user.UserOut{
		BaseResponse: inout.BaseResponse{
			ErrorCode:        0,
			ErrorDescription: "Success",
		},
		Data: user.FromModel(foundUser),
	}

	context.JSON(http.StatusOK, response)
}

// UpdateProfile updates current user's profile
func (controller UserController) UpdateProfile(context *gin.Context) {
	userID, err := utils.ExtractUserID(context)
	if err != nil {
		utils.ReportUnauthorized(context, "Authentication required")
		return
	}

	var request user.UpdateUserRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		utils.ReportBadRequest(context, "Invalid request data")
		return
	}

	userDao := dao.NewUserDao()
	existingUser, err := userDao.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ReportNotFound(context, "User not found")
		} else {
			utils.ReportInternalServerError(context, "Database error")
		}
		return
	}

	existingUser.FirstName = request.FirstName
	existingUser.LastName = request.LastName
	existingUser.AvatarURL = request.AvatarURL

	if err := userDao.Update(existingUser); err != nil {
		utils.ReportInternalServerError(context, "Failed to update user")
		return
	}

	response := user.UserOut{
		BaseResponse: inout.BaseResponse{
			ErrorCode:        0,
			ErrorDescription: "Success",
		},
		Data: user.FromModel(existingUser),
	}

	context.JSON(http.StatusOK, response)
}

// DeleteAccount deletes current user's account
func (controller UserController) DeleteAccount(context *gin.Context) {
	userID, err := utils.ExtractUserID(context)
	if err != nil {
		utils.ReportUnauthorized(context, "Authentication required")
		return
	}

	userDao := dao.NewUserDao()
	_, err = userDao.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ReportNotFound(context, "User not found")
		} else {
			utils.ReportInternalServerError(context, "Database error")
		}
		return
	}

	if err := userDao.Delete(userID); err != nil {
		utils.ReportInternalServerError(context, "Failed to delete user")
		return
	}

	response := inout.BaseResponse{
		ErrorCode:        0,
		ErrorDescription: "Account deleted successfully",
	}

	context.JSON(http.StatusOK, response)
}

// GetDashboard returns user dashboard data
func (controller UserController) GetDashboard(context *gin.Context) {
	userID, err := utils.ExtractUserID(context)
	if err != nil {
		utils.ReportUnauthorized(context, "Authentication required")
		return
	}

	userDao := dao.NewUserDao()
	foundUser, err := userDao.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ReportNotFound(context, "User not found")
		} else {
			utils.ReportInternalServerError(context, "Database error")
		}
		return
	}

	// TODO: Implement dashboard data aggregation
	dashboardData := user.DashboardOut{
		BaseResponse: inout.BaseResponse{
			ErrorCode:        0,
			ErrorDescription: "Success",
		},
		Data: user.DashboardData{
			User:              user.FromModel(foundUser),
			PersonalProjects:  0,                     // TODO: Get actual count from projects
			OrganizationCount: 0,                     // TODO: Get actual count from organizations
			RecentActivity:    []user.ActivityItem{}, // TODO: Get recent activity
		},
	}

	context.JSON(http.StatusOK, dashboardData)
}

// GetNotifications returns user notifications
func (controller UserController) GetNotifications(context *gin.Context) {
	userID, err := utils.ExtractUserID(context)
	if err != nil {
		utils.ReportUnauthorized(context, "Authentication required")
		return
	}

	// TODO: Implement notifications retrieval
	// For now, use userID in a placeholder way to avoid compiler error
	_ = userID

	response := user.NotificationsOut{
		BaseResponse: inout.BaseResponse{
			ErrorCode:        0,
			ErrorDescription: "Success",
		},
		Data: []user.Notification{}, // TODO: Get actual notifications
	}

	context.JSON(http.StatusOK, response)
}

// MarkNotificationRead marks a notification as read
func (controller UserController) MarkNotificationRead(context *gin.Context) {
	notificationIDParam := context.Param("id")
	notificationID, err := uuid.Parse(notificationIDParam)
	if err != nil {
		utils.ReportBadRequest(context, "Invalid notification ID format")
		return
	}

	userID, err := utils.ExtractUserID(context)
	if err != nil {
		utils.ReportUnauthorized(context, "Authentication required")
		return
	}

	// TODO: Implement notification mark as read functionality
	// For now, just return success
	_ = userID
	_ = notificationID

	response := inout.BaseResponse{
		ErrorCode:        0,
		ErrorDescription: "Notification marked as read",
	}

	context.JSON(http.StatusOK, response)
}
