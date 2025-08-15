package dao

import (
	"testlake/model"
	"time"

	"github.com/google/uuid"
)

type OrganizationMemberDao struct {
	Limit int
}

func NewOrganizationMemberDao() *OrganizationMemberDao {
	return &OrganizationMemberDao{Limit: 50}
}

// GetMembers returns all members of an organization
func (dao *OrganizationMemberDao) GetMembers(orgID uuid.UUID) ([]model.OrganizationMember, error) {
	var members []model.OrganizationMember
	err := Database.
		Preload("User").
		Where("organization_id = ? AND status = ?", orgID, "joined").
		Find(&members).Error
	return members, err
}

// GetMemberByUserID returns a specific member by organization and user ID
func (dao *OrganizationMemberDao) GetMemberByUserID(orgID, userID uuid.UUID) (*model.OrganizationMember, error) {
	var member model.OrganizationMember
	err := Database.
		Preload("User").
		Where("organization_id = ? AND user_id = ?", orgID, userID).
		First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

// AddMember adds a new member to an organization
func (dao *OrganizationMemberDao) AddMember(member *model.OrganizationMember) error {
	return Database.Create(member).Error
}

// UpdateMemberRole updates a member's role
func (dao *OrganizationMemberDao) UpdateMemberRole(orgID, userID uuid.UUID, role model.OrganizationMemberRole) error {
	return Database.Model(&model.OrganizationMember{}).
		Where("organization_id = ? AND user_id = ?", orgID, userID).
		Update("role", role).Error
}

// RemoveMember removes a member from an organization
func (dao *OrganizationMemberDao) RemoveMember(orgID, userID uuid.UUID) error {
	return Database.
		Where("organization_id = ? AND user_id = ?", orgID, userID).
		Delete(&model.OrganizationMember{}).Error
}

// IsUserMember checks if a user is a member of an organization
func (dao *OrganizationMemberDao) IsUserMember(orgID, userID uuid.UUID) (bool, error) {
	var count int64
	err := Database.Model(&model.OrganizationMember{}).
		Where("organization_id = ? AND user_id = ? AND status = ?", orgID, userID, "joined").
		Count(&count).Error
	return count > 0, err
}

// GetUserRole returns the user's role in an organization
func (dao *OrganizationMemberDao) GetUserRole(orgID, userID uuid.UUID) (model.OrganizationMemberRole, error) {
	var member model.OrganizationMember
	err := Database.
		Where("organization_id = ? AND user_id = ? AND status = ?", orgID, userID, "joined").
		First(&member).Error
	if err != nil {
		return "", err
	}
	return member.Role, nil
}

// OrganizationInvitationDao handles invitation operations
type OrganizationInvitationDao struct{}

func NewOrganizationInvitationDao() *OrganizationInvitationDao {
	return &OrganizationInvitationDao{}
}

// CreateInvitation creates a new organization invitation
func (dao *OrganizationInvitationDao) CreateInvitation(invitation *model.OrganizationInvitation) error {
	return Database.Create(invitation).Error
}

// GetInvitationByToken returns an invitation by token
func (dao *OrganizationInvitationDao) GetInvitationByToken(token string) (*model.OrganizationInvitation, error) {
	var invitation model.OrganizationInvitation
	err := Database.
		Preload("Organization").
		Where("token = ? AND status = ?", token, "pending").
		First(&invitation).Error
	if err != nil {
		return nil, err
	}
	return &invitation, nil
}

// GetPendingInvitationByEmail returns pending invitation by email and organization
func (dao *OrganizationInvitationDao) GetPendingInvitationByEmail(orgID uuid.UUID, email string) (*model.OrganizationInvitation, error) {
	var invitation model.OrganizationInvitation
	err := Database.
		Where("organization_id = ? AND email = ? AND status = ?", orgID, email, "pending").
		First(&invitation).Error
	if err != nil {
		return nil, err
	}
	return &invitation, nil
}

// AcceptInvitation marks an invitation as accepted
func (dao *OrganizationInvitationDao) AcceptInvitation(token string) error {
	now := time.Now()
	return Database.Model(&model.OrganizationInvitation{}).
		Where("token = ?", token).
		Updates(map[string]interface{}{
			"status":  "accepted",
			"used_at": &now,
		}).Error
}

// CancelInvitation cancels an invitation
func (dao *OrganizationInvitationDao) CancelInvitation(id uuid.UUID) error {
	return Database.Model(&model.OrganizationInvitation{}).
		Where("id = ?", id).
		Update("status", "cancelled").Error
}

// CleanupExpiredInvitations removes expired invitations
func (dao *OrganizationInvitationDao) CleanupExpiredInvitations() error {
	return Database.Model(&model.OrganizationInvitation{}).
		Where("expires_at < ? AND status = ?", time.Now(), "pending").
		Update("status", "expired").Error
}

// GetOrganizationInvitations returns all invitations for an organization
func (dao *OrganizationInvitationDao) GetOrganizationInvitations(orgID uuid.UUID) ([]model.OrganizationInvitation, error) {
	var invitations []model.OrganizationInvitation
	err := Database.
		Where("organization_id = ?", orgID).
		Find(&invitations).Error
	return invitations, err
}

// GetPendingInvitations returns only pending invitations for an organization
func (dao *OrganizationInvitationDao) GetPendingInvitations(orgID uuid.UUID) ([]model.OrganizationInvitation, error) {
	var invitations []model.OrganizationInvitation
	err := Database.
		Where("organization_id = ? AND status = ?", orgID, "pending").
		Find(&invitations).Error
	return invitations, err
}

// GetPendingInvitationsByEmail returns all pending invitations for a user by email
func (dao *OrganizationInvitationDao) GetPendingInvitationsByEmail(email string) ([]model.OrganizationInvitation, error) {
	var invitations []model.OrganizationInvitation
	err := Database.
		Preload("Organization").
		Where("email = ? AND status = ? AND expires_at > ?", email, "pending", time.Now()).
		Find(&invitations).Error
	return invitations, err
}
