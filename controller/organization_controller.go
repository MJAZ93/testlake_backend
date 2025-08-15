package controller

import (
	"errors"
	"math"
	"net/http"
	"strconv"
	"time"

	"testlake/dao"
	"testlake/inout"
	"testlake/inout/organization"
	"testlake/model"
	"testlake/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrganizationController struct{}

// CreateOrganization creates a new organization
func (controller OrganizationController) CreateOrganization(context *gin.Context) {
	userID, err := utils.ExtractUserID(context)
	if err != nil {
		utils.ReportUnauthorized(context, "Authentication required")
		return
	}

	var req organization.CreateOrganizationRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		utils.ReportBadRequest(context, "Invalid request data: "+err.Error())
		return
	}

	// Check if slug already exists
	orgDao := dao.NewOrganizationDao()
	exists, err := orgDao.SlugExists(req.Slug)
	if err != nil {
		utils.ReportInternalServerError(context, "Database error")
		return
	}
	if exists {
		utils.ReportBadRequest(context, "Organization slug already exists")
		return
	}

	// Create organization
	org := &model.Organization{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		LogoURL:     req.LogoURL,
		PlanType:    req.PlanType,
		CreatedBy:   userID,
		Status:      model.OrganizationStatusActive,
	}

	// Set default values if not provided
	if req.MaxUsers != nil {
		org.MaxUsers = *req.MaxUsers
	} else {
		org.MaxUsers = 10 // Default value
	}

	if req.MaxProjects != nil {
		org.MaxProjects = *req.MaxProjects
	} else {
		org.MaxProjects = 5 // Default value
	}

	// Set default plan type if not provided
	if org.PlanType == "" {
		org.PlanType = model.PlanTypeStarter
	}

	if err := orgDao.Create(org); err != nil {
		utils.ReportInternalServerError(context, "Failed to create organization")
		return
	}

	response := organization.OrganizationOut{
		BaseResponse: inout.BaseResponse{
			ErrorCode:        0,
			ErrorDescription: "Success",
		},
		Data: organization.FromModel(org),
	}

	context.JSON(http.StatusCreated, response)
}

// GetOrganizations returns paginated list of organizations for current user
func (controller OrganizationController) GetOrganizations(context *gin.Context) {
	userID, err := utils.ExtractUserID(context)
	if err != nil {
		utils.ReportUnauthorized(context, "Authentication required")
		return
	}

	pageStr := context.DefaultQuery("page", "0")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 0 {
		page = 0
	}

	orgDao := dao.NewOrganizationDao()
	orgs, total, err := orgDao.GetByCreatedBy(userID, page)
	if err != nil {
		utils.ReportInternalServerError(context, "Database error")
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(orgDao.Limit)))

	response := organization.OrganizationListOut{
		BaseResponse: inout.BaseResponse{
			ErrorCode:        0,
			ErrorDescription: "Success",
		},
		List: organization.FromModelList(orgs),
		Meta: inout.PaginationMeta{
			Page:       page,
			Limit:      orgDao.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}

	context.JSON(http.StatusOK, response)
}

// GetPendingInvites returns pending invitations for an organization
func (controller OrganizationController) GetPendingInvites(context *gin.Context) {
	userID, err := utils.ExtractUserID(context)
	if err != nil {
		utils.ReportUnauthorized(context, "Authentication required")
		return
	}

	idParam := context.Param("id")
	orgID, err := uuid.Parse(idParam)
	if err != nil {
		utils.ReportBadRequest(context, "Invalid organization ID")
		return
	}

	orgDao := dao.NewOrganizationDao()
	org, err := orgDao.GetByID(orgID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ReportNotFound(context, "Organization not found")
		} else {
			utils.ReportInternalServerError(context, "Database error")
		}
		return
	}

	// Check if user has access to view invitations (creator or admin)
	memberDao := dao.NewOrganizationMemberDao()
	if org.CreatedBy != userID {
		role, err := memberDao.GetUserRole(orgID, userID)
		if err != nil {
			utils.ReportForbidden(context, "Access denied")
			return
		}
		if role != model.OrganizationMemberRoleAdmin {
			utils.ReportForbidden(context, "Only admins can view pending invitations")
			return
		}
	}

	// Get pending invitations
	invitationDao := dao.NewOrganizationInvitationDao()
	invitations, err := invitationDao.GetPendingInvitations(orgID)
	if err != nil {
		utils.ReportInternalServerError(context, "Database error")
		return
	}

	// Convert to response format
	invitesData := make([]organization.PendingInvite, len(invitations))
	for i, invitation := range invitations {
		invitesData[i] = organization.FromInvitationModel(&invitation)
	}

	response := organization.PendingInvitesOut{
		BaseResponse: inout.BaseResponse{
			ErrorCode:        0,
			ErrorDescription: "Success",
		},
		Data: invitesData,
	}

	context.JSON(http.StatusOK, response)
}

// GetOrganization returns organization by ID
func (controller OrganizationController) GetOrganization(context *gin.Context) {
	userID, err := utils.ExtractUserID(context)
	if err != nil {
		utils.ReportUnauthorized(context, "Authentication required")
		return
	}

	idParam := context.Param("id")
	orgID, err := uuid.Parse(idParam)
	if err != nil {
		utils.ReportBadRequest(context, "Invalid organization ID")
		return
	}

	orgDao := dao.NewOrganizationDao()
	org, err := orgDao.GetByID(orgID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ReportNotFound(context, "Organization not found")
		} else {
			utils.ReportInternalServerError(context, "Database error")
		}
		return
	}

	// Check if user has access to this organization (creator for now)
	if org.CreatedBy != userID {
		utils.ReportForbidden(context, "Access denied")
		return
	}

	response := organization.OrganizationOut{
		BaseResponse: inout.BaseResponse{
			ErrorCode:        0,
			ErrorDescription: "Success",
		},
		Data: organization.FromModel(org),
	}

	context.JSON(http.StatusOK, response)
}

// UpdateOrganization updates an organization
func (controller OrganizationController) UpdateOrganization(context *gin.Context) {
	userID, err := utils.ExtractUserID(context)
	if err != nil {
		utils.ReportUnauthorized(context, "Authentication required")
		return
	}

	idParam := context.Param("id")
	orgID, err := uuid.Parse(idParam)
	if err != nil {
		utils.ReportBadRequest(context, "Invalid organization ID")
		return
	}

	var req organization.UpdateOrganizationRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		utils.ReportBadRequest(context, "Invalid request data: "+err.Error())
		return
	}

	orgDao := dao.NewOrganizationDao()
	org, err := orgDao.GetByID(orgID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ReportNotFound(context, "Organization not found")
		} else {
			utils.ReportInternalServerError(context, "Database error")
		}
		return
	}

	// Check if user has access to update this organization (creator for now)
	if org.CreatedBy != userID {
		utils.ReportForbidden(context, "Access denied")
		return
	}

	// Update fields if provided
	if req.Name != nil {
		org.Name = *req.Name
	}
	if req.Description != nil {
		org.Description = req.Description
	}
	if req.LogoURL != nil {
		org.LogoURL = req.LogoURL
	}
	if req.PlanType != nil {
		org.PlanType = *req.PlanType
	}
	if req.MaxUsers != nil {
		org.MaxUsers = *req.MaxUsers
	}
	if req.MaxProjects != nil {
		org.MaxProjects = *req.MaxProjects
	}

	if err := orgDao.Update(org); err != nil {
		utils.ReportInternalServerError(context, "Failed to update organization")
		return
	}

	response := organization.OrganizationOut{
		BaseResponse: inout.BaseResponse{
			ErrorCode:        0,
			ErrorDescription: "Success",
		},
		Data: organization.FromModel(org),
	}

	context.JSON(http.StatusOK, response)
}

// DeleteOrganization soft deletes an organization
func (controller OrganizationController) DeleteOrganization(context *gin.Context) {
	userID, err := utils.ExtractUserID(context)
	if err != nil {
		utils.ReportUnauthorized(context, "Authentication required")
		return
	}

	idParam := context.Param("id")
	orgID, err := uuid.Parse(idParam)
	if err != nil {
		utils.ReportBadRequest(context, "Invalid organization ID")
		return
	}

	orgDao := dao.NewOrganizationDao()
	org, err := orgDao.GetByID(orgID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ReportNotFound(context, "Organization not found")
		} else {
			utils.ReportInternalServerError(context, "Database error")
		}
		return
	}

	// Check if user has access to delete this organization (creator for now)
	if org.CreatedBy != userID {
		utils.ReportForbidden(context, "Access denied")
		return
	}

	if err := orgDao.Delete(orgID); err != nil {
		utils.ReportInternalServerError(context, "Failed to delete organization")
		return
	}

	response := inout.BaseResponse{
		ErrorCode:        0,
		ErrorDescription: "Organization deleted successfully",
	}

	context.JSON(http.StatusOK, response)
}

// GetOrganizationMembers returns organization members
func (controller OrganizationController) GetOrganizationMembers(context *gin.Context) {
	userID, err := utils.ExtractUserID(context)
	if err != nil {
		utils.ReportUnauthorized(context, "Authentication required")
		return
	}

	idParam := context.Param("id")
	orgID, err := uuid.Parse(idParam)
	if err != nil {
		utils.ReportBadRequest(context, "Invalid organization ID")
		return
	}

	orgDao := dao.NewOrganizationDao()
	org, err := orgDao.GetByID(orgID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ReportNotFound(context, "Organization not found")
		} else {
			utils.ReportInternalServerError(context, "Database error")
		}
		return
	}

	// Check if user has access to view members (creator or member)
	memberDao := dao.NewOrganizationMemberDao()
	if org.CreatedBy != userID {
		isMember, err := memberDao.IsUserMember(orgID, userID)
		if err != nil {
			utils.ReportInternalServerError(context, "Database error")
			return
		}
		if !isMember {
			utils.ReportForbidden(context, "Access denied")
			return
		}
	}

	// Get organization members
	orgMembers, err := orgDao.GetMembers(orgID)
	if err != nil {
		utils.ReportInternalServerError(context, "Database error")
		return
	}

	// Convert to response format
	members := make([]organization.Member, len(orgMembers))
	for i, member := range orgMembers {
		members[i] = organization.FromUserModel(&member.User, string(member.Role), member.InvitedAt)
		if member.JoinedAt != nil {
			members[i].JoinedAt = *member.JoinedAt
		}
	}

	response := organization.MembersOut{
		BaseResponse: inout.BaseResponse{
			ErrorCode:        0,
			ErrorDescription: "Success",
		},
		Data: members,
	}

	context.JSON(http.StatusOK, response)
}

// InviteMember invites a user to the organization
func (controller OrganizationController) InviteMember(context *gin.Context) {
	userID, err := utils.ExtractUserID(context)
	if err != nil {
		utils.ReportUnauthorized(context, "Authentication required")
		return
	}

	idParam := context.Param("id")
	orgID, err := uuid.Parse(idParam)
	if err != nil {
		utils.ReportBadRequest(context, "Invalid organization ID")
		return
	}

	var req organization.InviteMemberRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		utils.ReportBadRequest(context, "Invalid request data: "+err.Error())
		return
	}

	orgDao := dao.NewOrganizationDao()
	org, err := orgDao.GetByID(orgID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ReportNotFound(context, "Organization not found")
		} else {
			utils.ReportInternalServerError(context, "Database error")
		}
		return
	}

	// Check if user has access to invite members (creator or admin)
	memberDao := dao.NewOrganizationMemberDao()
	if org.CreatedBy != userID {
		role, err := memberDao.GetUserRole(orgID, userID)
		if err != nil {
			utils.ReportForbidden(context, "Access denied")
			return
		}
		if role != model.OrganizationMemberRoleAdmin {
			utils.ReportForbidden(context, "Only admins can invite members")
			return
		}
	}

	// Check if user is already a member
	userDao := dao.NewUserDao()
	invitedUser, err := userDao.GetByEmail(req.Email)
	if err == nil {
		// User exists, check if already a member
		isMember, err := memberDao.IsUserMember(orgID, invitedUser.ID)
		if err != nil {
			utils.ReportInternalServerError(context, "Database error")
			return
		}
		if isMember {
			utils.ReportBadRequest(context, "User is already a member of this organization")
			return
		}
	}

	// Check if there's already a pending invitation
	invitationDao := dao.NewOrganizationInvitationDao()
	existingInvitation, err := invitationDao.GetPendingInvitationByEmail(orgID, req.Email)
	if err == nil && existingInvitation != nil {
		utils.ReportBadRequest(context, "Invitation already sent to this email")
		return
	}

	// Create invitation token
	token, err := utils.GenerateSecureToken(32)
	if err != nil {
		utils.ReportInternalServerError(context, "Failed to generate invitation token")
		return
	}

	// Parse role
	var inviteRole model.OrganizationMemberRole
	switch req.Role {
	case "member":
		inviteRole = model.OrganizationMemberRoleMember
	case "admin":
		inviteRole = model.OrganizationMemberRoleAdmin
	default:
		inviteRole = model.OrganizationMemberRoleMember
	}

	// Create invitation record
	invitation := &model.OrganizationInvitation{
		OrganizationID: orgID,
		Email:          req.Email,
		Role:           inviteRole,
		Token:          token,
		InvitedBy:      userID,
		ExpiresAt:      time.Now().Add(7 * 24 * time.Hour), // Expires in 7 days
		Status:         "pending",
	}

	if err := invitationDao.CreateInvitation(invitation); err != nil {
		utils.ReportInternalServerError(context, "Failed to create invitation")
		return
	}

	// TODO: Send invitation email here
	// This would include sending an email with the invitation link

	result := organization.InviteResult{
		Email:   req.Email,
		Status:  "invited",
		Message: "Invitation sent successfully",
	}

	response := organization.InviteOut{
		BaseResponse: inout.BaseResponse{
			ErrorCode:        0,
			ErrorDescription: "Success",
		},
		Data: result,
	}

	context.JSON(http.StatusOK, response)
}

// RemoveMember removes a user from the organization
func (controller OrganizationController) RemoveMember(context *gin.Context) {
	userID, err := utils.ExtractUserID(context)
	if err != nil {
		utils.ReportUnauthorized(context, "Authentication required")
		return
	}

	idParam := context.Param("id")
	orgID, err := uuid.Parse(idParam)
	if err != nil {
		utils.ReportBadRequest(context, "Invalid organization ID")
		return
	}

	userIDParam := context.Param("userId")
	memberUserID, err := uuid.Parse(userIDParam)
	if err != nil {
		utils.ReportBadRequest(context, "Invalid user ID")
		return
	}

	orgDao := dao.NewOrganizationDao()
	org, err := orgDao.GetByID(orgID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ReportNotFound(context, "Organization not found")
		} else {
			utils.ReportInternalServerError(context, "Database error")
		}
		return
	}

	// Check if user has access to remove members (creator or admin)
	memberDao := dao.NewOrganizationMemberDao()
	if org.CreatedBy != userID {
		role, err := memberDao.GetUserRole(orgID, userID)
		if err != nil {
			utils.ReportForbidden(context, "Access denied")
			return
		}
		if role != model.OrganizationMemberRoleAdmin {
			utils.ReportForbidden(context, "Only admins can remove members")
			return
		}
	}

	// Cannot remove self
	if memberUserID == userID {
		utils.ReportBadRequest(context, "Cannot remove yourself from the organization")
		return
	}

	// Check if target user is a member
	isMember, err := memberDao.IsUserMember(orgID, memberUserID)
	if err != nil {
		utils.ReportInternalServerError(context, "Database error")
		return
	}
	if !isMember {
		utils.ReportNotFound(context, "User is not a member of this organization")
		return
	}

	// Remove member
	if err := memberDao.RemoveMember(orgID, memberUserID); err != nil {
		utils.ReportInternalServerError(context, "Failed to remove member")
		return
	}

	response := inout.BaseResponse{
		ErrorCode:        0,
		ErrorDescription: "Member removed successfully",
	}

	context.JSON(http.StatusOK, response)
}

// UpdateMemberRole updates a member's role in the organization
func (controller OrganizationController) UpdateMemberRole(context *gin.Context) {
	userID, err := utils.ExtractUserID(context)
	if err != nil {
		utils.ReportUnauthorized(context, "Authentication required")
		return
	}

	idParam := context.Param("id")
	orgID, err := uuid.Parse(idParam)
	if err != nil {
		utils.ReportBadRequest(context, "Invalid organization ID")
		return
	}

	userIDParam := context.Param("userId")
	memberUserID, err := uuid.Parse(userIDParam)
	if err != nil {
		utils.ReportBadRequest(context, "Invalid user ID")
		return
	}

	var req organization.UpdateMemberRoleRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		utils.ReportBadRequest(context, "Invalid request data: "+err.Error())
		return
	}

	orgDao := dao.NewOrganizationDao()
	org, err := orgDao.GetByID(orgID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ReportNotFound(context, "Organization not found")
		} else {
			utils.ReportInternalServerError(context, "Database error")
		}
		return
	}

	// Check if user has access to update member roles (creator or admin)
	memberDao := dao.NewOrganizationMemberDao()
	if org.CreatedBy != userID {
		role, err := memberDao.GetUserRole(orgID, userID)
		if err != nil {
			utils.ReportForbidden(context, "Access denied")
			return
		}
		if role != model.OrganizationMemberRoleAdmin {
			utils.ReportForbidden(context, "Only admins can update member roles")
			return
		}
	}

	// Cannot update own role
	if memberUserID == userID {
		utils.ReportBadRequest(context, "Cannot update your own role")
		return
	}

	// Check if target user is a member
	isMember, err := memberDao.IsUserMember(orgID, memberUserID)
	if err != nil {
		utils.ReportInternalServerError(context, "Database error")
		return
	}
	if !isMember {
		utils.ReportNotFound(context, "User is not a member of this organization")
		return
	}

	// Parse and validate new role
	var newRole model.OrganizationMemberRole
	switch req.Role {
	case "member":
		newRole = model.OrganizationMemberRoleMember
	case "admin":
		newRole = model.OrganizationMemberRoleAdmin
	default:
		utils.ReportBadRequest(context, "Invalid role specified")
		return
	}

	// Update member role
	if err := memberDao.UpdateMemberRole(orgID, memberUserID, newRole); err != nil {
		utils.ReportInternalServerError(context, "Failed to update member role")
		return
	}

	response := inout.BaseResponse{
		ErrorCode:        0,
		ErrorDescription: "Member role updated successfully",
	}

	context.JSON(http.StatusOK, response)
}
