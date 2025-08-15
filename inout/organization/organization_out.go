package organization

import (
	"time"

	"testlake/inout"
	"testlake/model"

	"github.com/google/uuid"
)

type Organization struct {
	ID          uuid.UUID                `json:"id"`
	Name        string                   `json:"name"`
	Slug        string                   `json:"slug"`
	Description *string                  `json:"description"`
	LogoURL     *string                  `json:"logo_url"`
	PlanType    model.PlanType           `json:"plan_type"`
	MaxUsers    int                      `json:"max_users"`
	MaxProjects int                      `json:"max_projects"`
	CreatedBy   uuid.UUID                `json:"created_by"`
	CreatedAt   time.Time                `json:"created_at"`
	UpdatedAt   time.Time                `json:"updated_at"`
	Status      model.OrganizationStatus `json:"status"`
}

type OrganizationOut struct {
	inout.BaseResponse
	Data Organization `json:"data"`
}

type OrganizationListOut struct {
	inout.BaseResponse
	List []Organization       `json:"list"`
	Meta inout.PaginationMeta `json:"meta"`
}

type Member struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
	FullName *string   `json:"full_name"`
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
}

type MembersOut struct {
	inout.BaseResponse
	Data []Member `json:"data"`
}

type InviteResult struct {
	Email   string `json:"email"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type InviteOut struct {
	inout.BaseResponse
	Data InviteResult `json:"data"`
}

type PendingInvite struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	InvitedAt time.Time `json:"invited_at"`
	ExpiresAt time.Time `json:"expires_at"`
	Token     string    `json:"token"`
}

type PendingInvitesOut struct {
	inout.BaseResponse
	Data []PendingInvite `json:"data"`
}

func FromModel(org *model.Organization) Organization {
	return Organization{
		ID:          org.ID,
		Name:        org.Name,
		Slug:        org.Slug,
		Description: org.Description,
		LogoURL:     org.LogoURL,
		PlanType:    org.PlanType,
		MaxUsers:    org.MaxUsers,
		MaxProjects: org.MaxProjects,
		CreatedBy:   org.CreatedBy,
		CreatedAt:   org.CreatedAt,
		UpdatedAt:   org.UpdatedAt,
		Status:      org.Status,
	}
}

func FromModelList(orgs []model.Organization) []Organization {
	result := make([]Organization, len(orgs))
	for i, org := range orgs {
		result[i] = FromModel(&org)
	}
	return result
}

func FromUserModel(user *model.User, role string, joinedAt time.Time) Member {
	var fullName *string
	if user.FirstName != nil && user.LastName != nil {
		full := *user.FirstName + " " + *user.LastName
		fullName = &full
	} else if user.FirstName != nil {
		fullName = user.FirstName
	} else if user.LastName != nil {
		fullName = user.LastName
	}

	return Member{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
		FullName: fullName,
		Role:     role,
		JoinedAt: joinedAt,
	}
}

func FromInvitationModel(invitation *model.OrganizationInvitation) PendingInvite {
	return PendingInvite{
		ID:        invitation.ID,
		Email:     invitation.Email,
		Role:      string(invitation.Role),
		InvitedAt: invitation.InvitedAt,
		ExpiresAt: invitation.ExpiresAt,
	}
}

func FromInvitationModelList(invitations []model.OrganizationInvitation) []PendingInvite {
	result := make([]PendingInvite, len(invitations))
	for i, invitation := range invitations {
		result[i] = FromInvitationModel(&invitation)
	}
	return result
}
