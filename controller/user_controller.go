package controller

import (
	"errors"
	"net/http"
	"testlake/model"
	"time"

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

// DenyInvite declines an organization invitation
func (controller UserController) DenyInvite(context *gin.Context) {
	userID, err := utils.ExtractUserID(context)
	if err != nil {
		utils.ReportUnauthorized(context, "Authentication required")
		return
	}

	token := context.Param("token")
	if token == "" {
		utils.ReportBadRequest(context, "Invitation token is required")
		return
	}

	// Get invitation by token
	invitationDao := dao.NewOrganizationInvitationDao()
	invitation, err := invitationDao.GetInvitationByToken(token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ReportNotFound(context, "Invitation not found or expired")
		} else {
			utils.ReportInternalServerError(context, "Database error")
		}
		return
	}

	// Check if invitation has expired
	if time.Now().After(invitation.ExpiresAt) {
		utils.ReportBadRequest(context, "Invitation has expired")
		return
	}

	// Get user to verify email matches
	userDao := dao.NewUserDao()
	u, err := userDao.GetByID(userID)
	if err != nil {
		utils.ReportInternalServerError(context, "Database error")
		return
	}

	if u.Email != invitation.Email {
		utils.ReportBadRequest(context, "Invitation email does not match your account email")
		return
	}

	// Mark invitation as cancelled (denied)
	if err := invitationDao.CancelInvitation(invitation.ID); err != nil {
		utils.ReportInternalServerError(context, "Failed to deny invitation")
		return
	}

	response := user.DenyInviteOut{
		BaseResponse: inout.BaseResponse{
			ErrorCode:        0,
			ErrorDescription: "Invitation declined successfully",
		},
		Data: user.DenyInviteResult{
			OrganizationID:   invitation.OrganizationID,
			OrganizationName: invitation.Organization.Name,
			Status:           "declined",
			Message:          "You have successfully declined the invitation",
		},
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

// GetPendingInvites returns pending invitations for current user
func (controller UserController) GetPendingInvites(context *gin.Context) {
	userID, err := utils.ExtractUserID(context)
	if err != nil {
		utils.ReportUnauthorized(context, "Authentication required")
		return
	}

	// Get current user to retrieve email
	userDao := dao.NewUserDao()
	currentUser, err := userDao.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ReportNotFound(context, "User not found")
		} else {
			utils.ReportInternalServerError(context, "Database error")
		}
		return
	}

	// Get pending invitations by email
	invitationDao := dao.NewOrganizationInvitationDao()
	invitations, err := invitationDao.GetPendingInvitationsByEmail(currentUser.Email)
	if err != nil {
		utils.ReportInternalServerError(context, "Database error")
		return
	}

	// Convert to response format
	pendingInvites := make([]user.PendingInvite, len(invitations))
	for i, invitation := range invitations {
		pendingInvites[i] = user.PendingInvite{
			ID:               invitation.ID,
			OrganizationID:   invitation.OrganizationID,
			OrganizationName: invitation.Organization.Name,
			Role:             string(invitation.Role),
			InvitedAt:        invitation.InvitedAt,
			ExpiresAt:        invitation.ExpiresAt,
			Token:            invitation.Token,
		}
	}

	response := user.PendingInvitesOut{
		BaseResponse: inout.BaseResponse{
			ErrorCode:        0,
			ErrorDescription: "Success",
		},
		Data: pendingInvites,
	}

	context.JSON(http.StatusOK, response)
}

// AcceptInvite accepts an organization invitation
func (controller UserController) AcceptInvite(context *gin.Context) {
	userID, err := utils.ExtractUserID(context)
	if err != nil {
		utils.ReportUnauthorized(context, "Authentication required")
		return
	}

	token := context.Param("token")
	if token == "" {
		utils.ReportBadRequest(context, "Invitation token is required")
		return
	}

	// Get invitation by token
	invitationDao := dao.NewOrganizationInvitationDao()
	invitation, err := invitationDao.GetInvitationByToken(token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ReportNotFound(context, "Invitation not found or expired")
		} else {
			utils.ReportInternalServerError(context, "Database error")
		}
		return
	}

	// Check if invitation has expired
	if time.Now().After(invitation.ExpiresAt) {
		utils.ReportBadRequest(context, "Invitation has expired")
		return
	}

	// Get user to verify email matches
	userDao := dao.NewUserDao()
	u, err := userDao.GetByID(userID)
	if err != nil {
		utils.ReportInternalServerError(context, "Database error")
		return
	}

	if u.Email != invitation.Email {
		utils.ReportBadRequest(context, "Invitation email does not match your account email")
		return
	}

	// Check if user is already a member
	memberDao := dao.NewOrganizationMemberDao()
	isMember, err := memberDao.IsUserMember(invitation.OrganizationID, userID)
	if err != nil {
		utils.ReportInternalServerError(context, "Database error")
		return
	}
	if isMember {
		utils.ReportBadRequest(context, "You are already a member of this organization")
		return
	}

	// Accept invitation and add user as member
	now := time.Now()
	member := &model.OrganizationMember{
		OrganizationID: invitation.OrganizationID,
		UserID:         userID,
		Role:           invitation.Role,
		InvitedBy:      invitation.InvitedBy,
		InvitedAt:      invitation.InvitedAt,
		JoinedAt:       &now,
		Status:         "joined",
	}

	if err := memberDao.AddMember(member); err != nil {
		utils.ReportInternalServerError(context, "Failed to add member")
		return
	}

	// Mark invitation as accepted
	if err := invitationDao.AcceptInvitation(token); err != nil {
		utils.ReportInternalServerError(context, "Failed to accept invitation")
		return
	}

	response := user.AcceptInviteOut{
		BaseResponse: inout.BaseResponse{
			ErrorCode:        0,
			ErrorDescription: "Invitation accepted successfully",
		},
		Data: user.AcceptInviteResult{
			OrganizationID:   invitation.OrganizationID,
			OrganizationName: invitation.Organization.Name,
			Role:             string(invitation.Role),
			Status:           "accepted",
			Message:          "Invitation accepted successfully",
		},
	}

	context.JSON(http.StatusOK, response)
}
