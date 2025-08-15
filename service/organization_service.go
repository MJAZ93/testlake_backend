package service

import (
	"testlake/controller"

	"github.com/gin-gonic/gin"
)

type OrganizationService struct {
	Route      string
	Controller controller.OrganizationController
}

// CreateOrganization godoc
// @Summary Create organization
// @Description Create a new organization
// @Tags Organization Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token" format(Bearer {token})
// @Param organization body organization.CreateOrganizationRequest true "Organization data"
// @Success 201 {object} organization.OrganizationOut
// @Failure 400 {object} inout.BaseResponse
// @Failure 401 {object} inout.BaseResponse
// @Router /api/v1/organizations [POST]
func (s OrganizationService) CreateOrganization(r *gin.RouterGroup) {
	r.POST("/"+s.Route, s.Controller.CreateOrganization)
}

// GetOrganizations godoc
// @Summary Get organizations
// @Description Get paginated list of organizations for current user
// @Tags Organization Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token" format(Bearer {token})
// @Param page query int false "Page number (default 0)"
// @Success 200 {object} organization.OrganizationListOut
// @Failure 401 {object} inout.BaseResponse
// @Router /api/v1/organizations [GET]
func (s OrganizationService) GetOrganizations(r *gin.RouterGroup) {
	r.GET("/"+s.Route, s.Controller.GetOrganizations)
}

// GetOrganization godoc
// @Summary Get organization
// @Description Get organization details by ID
// @Tags Organization Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token" format(Bearer {token})
// @Param id path string true "Organization ID"
// @Success 200 {object} organization.OrganizationOut
// @Failure 400 {object} inout.BaseResponse
// @Failure 401 {object} inout.BaseResponse
// @Failure 403 {object} inout.BaseResponse
// @Failure 404 {object} inout.BaseResponse
// @Router /api/v1/organizations/{id} [GET]
func (s OrganizationService) GetOrganization(r *gin.RouterGroup) {
	r.GET("/"+s.Route+"/:id", s.Controller.GetOrganization)
}

// UpdateOrganization godoc
// @Summary Update organization
// @Description Update organization details
// @Tags Organization Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token" format(Bearer {token})
// @Param id path string true "Organization ID"
// @Param organization body organization.UpdateOrganizationRequest true "Organization update data"
// @Success 200 {object} organization.OrganizationOut
// @Failure 400 {object} inout.BaseResponse
// @Failure 401 {object} inout.BaseResponse
// @Failure 403 {object} inout.BaseResponse
// @Failure 404 {object} inout.BaseResponse
// @Router /api/v1/organizations/{id} [PUT]
func (s OrganizationService) UpdateOrganization(r *gin.RouterGroup) {
	r.PUT("/"+s.Route+"/:id", s.Controller.UpdateOrganization)
}

// DeleteOrganization godoc
// @Summary Delete organization
// @Description Soft delete an organization
// @Tags Organization Management
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
// @Router /api/v1/organizations/{id} [DELETE]
func (s OrganizationService) DeleteOrganization(r *gin.RouterGroup) {
	r.DELETE("/"+s.Route+"/:id", s.Controller.DeleteOrganization)
}

// GetOrganizationMembers godoc
// @Summary Get organization members
// @Description Get list of organization members
// @Tags Organization Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token" format(Bearer {token})
// @Param id path string true "Organization ID"
// @Success 200 {object} organization.MembersOut
// @Failure 400 {object} inout.BaseResponse
// @Failure 401 {object} inout.BaseResponse
// @Failure 403 {object} inout.BaseResponse
// @Failure 404 {object} inout.BaseResponse
// @Router /api/v1/organizations/{id}/members [GET]
func (s OrganizationService) GetOrganizationMembers(r *gin.RouterGroup) {
	r.GET("/"+s.Route+"/:id/members", s.Controller.GetOrganizationMembers)
}

// InviteMember godoc
// @Summary Invite member
// @Description Invite a user to join the organization
// @Tags Organization Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token" format(Bearer {token})
// @Param id path string true "Organization ID"
// @Param invite body organization.InviteMemberRequest true "Invitation data"
// @Success 200 {object} organization.InviteOut
// @Failure 400 {object} inout.BaseResponse
// @Failure 401 {object} inout.BaseResponse
// @Failure 403 {object} inout.BaseResponse
// @Failure 404 {object} inout.BaseResponse
// @Router /api/v1/organizations/{id}/invite [POST]
func (s OrganizationService) InviteMember(r *gin.RouterGroup) {
	r.POST("/"+s.Route+"/:id/invite", s.Controller.InviteMember)
}

// RemoveMember godoc
// @Summary Remove member
// @Description Remove a member from the organization
// @Tags Organization Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token" format(Bearer {token})
// @Param id path string true "Organization ID"
// @Param userId path string true "User ID to remove"
// @Success 200 {object} inout.BaseResponse
// @Failure 400 {object} inout.BaseResponse
// @Failure 401 {object} inout.BaseResponse
// @Failure 403 {object} inout.BaseResponse
// @Failure 404 {object} inout.BaseResponse
// @Router /api/v1/organizations/{id}/members/{userId} [DELETE]
func (s OrganizationService) RemoveMember(r *gin.RouterGroup) {
	r.DELETE("/"+s.Route+"/:id/members/:userId", s.Controller.RemoveMember)
}

// UpdateMemberRole godoc
// @Summary Update member role
// @Description Update a member's role in the organization
// @Tags Organization Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token" format(Bearer {token})
// @Param id path string true "Organization ID"
// @Param userId path string true "User ID"
// @Param role body organization.UpdateMemberRoleRequest true "Role update data"
// @Success 200 {object} inout.BaseResponse
// @Failure 400 {object} inout.BaseResponse
// @Failure 401 {object} inout.BaseResponse
// @Failure 403 {object} inout.BaseResponse
// @Failure 404 {object} inout.BaseResponse
// @Router /api/v1/organizations/{id}/members/{userId}/role [PUT]
func (s OrganizationService) UpdateMemberRole(r *gin.RouterGroup) {
	r.PUT("/"+s.Route+"/:id/members/:userId/role", s.Controller.UpdateMemberRole)
}

// GetPendingInvites godoc
// @Summary Get pending invites
// @Description Get list of pending invitations for the organization
// @Tags Organization Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token" format(Bearer {token})
// @Param id path string true "Organization ID"
// @Success 200 {object} organization.PendingInvitesOut
// @Failure 400 {object} inout.BaseResponse
// @Failure 401 {object} inout.BaseResponse
// @Failure 403 {object} inout.BaseResponse
// @Failure 404 {object} inout.BaseResponse
// @Router /api/v1/organizations/{id}/invites [GET]
func (s OrganizationService) GetPendingInvites(r *gin.RouterGroup) {
	r.GET("/"+s.Route+"/:id/invites", s.Controller.GetPendingInvites)
}
