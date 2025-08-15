package organization

import "testlake/model"

type CreateOrganizationRequest struct {
	Name        string         `json:"name" binding:"required,min=2,max=200"`
	Slug        string         `json:"slug" binding:"required,min=2,max=100,alphanum"`
	Description *string        `json:"description"`
	LogoURL     *string        `json:"logo_url"`
	PlanType    model.PlanType `json:"plan_type"`
	MaxUsers    *int           `json:"max_users"`
	MaxProjects *int           `json:"max_projects"`
}

type UpdateOrganizationRequest struct {
	Name        *string         `json:"name" binding:"omitempty,min=2,max=200"`
	Description *string         `json:"description"`
	LogoURL     *string         `json:"logo_url"`
	PlanType    *model.PlanType `json:"plan_type"`
	MaxUsers    *int            `json:"max_users" binding:"omitempty,min=1"`
	MaxProjects *int            `json:"max_projects" binding:"omitempty,min=1"`
}

type InviteMemberRequest struct {
	Email string `json:"email" binding:"required,email"`
	Role  string `json:"role" binding:"required,oneof=member admin"`
}

type UpdateMemberRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=member admin"`
}
