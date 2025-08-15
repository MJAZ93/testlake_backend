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
// @Param Authorization header string true "Bearer token" format(Bearer {token})
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
// @Param Authorization header string true "Bearer token" format(Bearer {token})
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
// @Param Authorization header string true "Bearer token" format(Bearer {token})
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
// @Param Authorization header string true "Bearer token" format(Bearer {token})
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
// @Param Authorization header string true "Bearer token" format(Bearer {token})
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
// @Param Authorization header string true "Bearer token" format(Bearer {token})
// @Success 200 {object} inout.BaseResponse
// @Failure 400 {object} inout.BaseResponse
// @Failure 401 {object} inout.BaseResponse
// @Router /api/v1/users/notifications/{id}/read [PUT]
func (s UserService) MarkNotificationRead(r *gin.RouterGroup, route string) {
	r.PUT("/"+s.Route+"/"+route+"/:id/read", s.Controller.MarkNotificationRead)
}

// GetPendingInvites godoc
// @Summary Get pending invitations
// @Description Get all pending organization invitations for the current user
// @Tags User Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token" format(Bearer {token})
// @Success 200 {object} user.PendingInvitesOut
// @Failure 401 {object} inout.BaseResponse
// @Router /api/v1/users/invites [GET]
func (s UserService) GetPendingInvites(r *gin.RouterGroup, route string) {
	r.GET("/"+s.Route+"/"+route, s.Controller.GetPendingInvites)
}

// AcceptInvite godoc
// @Summary Accept organization invite
// @Description Accept an organization invitation using the invitation token
// @Tags User Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token" format(Bearer {token})
// @Param token path string true "Invitation token"
// @Success 200 {object} user.AcceptInviteOut
// @Failure 400 {object} inout.BaseResponse
// @Failure 401 {object} inout.BaseResponse
// @Failure 404 {object} inout.BaseResponse
// @Router /api/v1/users/invites/{token}/accept [POST]
func (s UserService) AcceptInvite(r *gin.RouterGroup, route string) {
	r.POST("/"+s.Route+"/"+route+"/:token/accept", s.Controller.AcceptInvite)
}

// DenyInvite godoc
// @Summary Deny organization invite
// @Description Decline an organization invitation using the invitation token
// @Tags User Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token" format(Bearer {token})
// @Param token path string true "Invitation token"
// @Success 200 {object} user.DenyInviteOut
// @Failure 400 {object} inout.BaseResponse
// @Failure 401 {object} inout.BaseResponse
// @Failure 404 {object} inout.BaseResponse
// @Router /api/v1/users/invites/{token}/deny [POST]
func (s UserService) DenyInvite(r *gin.RouterGroup, route string) {
	r.POST("/"+s.Route+"/"+route+"/:token/deny", s.Controller.DenyInvite)
}
