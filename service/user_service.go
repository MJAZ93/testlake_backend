package service

import (
	"testlake/controller"

	"github.com/gin-gonic/gin"
)

type UserService struct {
	Route      string
	Controller controller.UserController
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get current user's profile information
// @Tags User Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} user.UserOut
// @Failure 401 {object} inout.BaseResponse
// @Router /api/v1/users/profile [GET]
func (s UserService) GetProfile(r *gin.RouterGroup, route string) {
	r.GET("/"+s.Route+"/"+route, s.Controller.GetProfile)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update current user's profile information
// @Tags User Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user body user.UpdateUserRequest true "User profile data"
// @Success 200 {object} user.UserOut
// @Failure 400 {object} inout.BaseResponse
// @Failure 401 {object} inout.BaseResponse
// @Router /api/v1/users/profile [PUT]
func (s UserService) UpdateProfile(r *gin.RouterGroup, route string) {
	r.PUT("/"+s.Route+"/"+route, s.Controller.UpdateProfile)
}

// DeleteAccount godoc
// @Summary Delete user account
// @Description Delete current user's account permanently
// @Tags User Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} inout.BaseResponse
// @Failure 401 {object} inout.BaseResponse
// @Router /api/v1/users/account [DELETE]
func (s UserService) DeleteAccount(r *gin.RouterGroup, route string) {
	r.DELETE("/"+s.Route+"/"+route, s.Controller.DeleteAccount)
}

// GetDashboard godoc
// @Summary Get user dashboard
// @Description Get user dashboard with overview of projects and activities
// @Tags User Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} user.DashboardOut
// @Failure 401 {object} inout.BaseResponse
// @Router /api/v1/users/dashboard [GET]
func (s UserService) GetDashboard(r *gin.RouterGroup, route string) {
	r.GET("/"+s.Route+"/"+route, s.Controller.GetDashboard)
}

// GetNotifications godoc
// @Summary Get user notifications
// @Description Get list of user notifications
// @Tags User Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} user.NotificationsOut
// @Failure 401 {object} inout.BaseResponse
// @Router /api/v1/users/notifications [GET]
func (s UserService) GetNotifications(r *gin.RouterGroup, route string) {
	r.GET("/"+s.Route+"/"+route, s.Controller.GetNotifications)
}

// MarkNotificationRead godoc
// @Summary Mark notification as read
// @Description Mark a specific notification as read
// @Tags User Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Notification ID"
// @Success 200 {object} inout.BaseResponse
// @Failure 400 {object} inout.BaseResponse
// @Failure 401 {object} inout.BaseResponse
// @Router /api/v1/users/notifications/{id}/read [PUT]
func (s UserService) MarkNotificationRead(r *gin.RouterGroup, route string) {
	r.PUT("/"+s.Route+"/"+route+"/:id/read", s.Controller.MarkNotificationRead)
}
