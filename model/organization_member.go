package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrganizationMemberRole string

const (
	OrganizationMemberRoleMember OrganizationMemberRole = "member"
	OrganizationMemberRoleAdmin  OrganizationMemberRole = "admin"
	OrganizationMemberRoleOwner  OrganizationMemberRole = "owner"
)

type OrganizationMember struct {
	ID             uuid.UUID              `gorm:"type:uuid;primaryKey" json:"id"`
	OrganizationID uuid.UUID              `gorm:"type:uuid;not null;uniqueIndex:idx_org_user" json:"organization_id"`
	UserID         uuid.UUID              `gorm:"type:uuid;not null;uniqueIndex:idx_org_user" json:"user_id"`
	Role           OrganizationMemberRole `gorm:"type:varchar(20);default:member" json:"role"`
	InvitedBy      uuid.UUID              `gorm:"type:uuid;not null" json:"invited_by"`
	InvitedAt      time.Time              `gorm:"default:now()" json:"invited_at"`
	JoinedAt       *time.Time             `json:"joined_at"`
	Status         string                 `gorm:"type:varchar(20);default:invited" json:"status"` // invited, joined, left

	// Relationships
	Organization  Organization `gorm:"foreignKey:OrganizationID;references:ID" json:"-"`
	User          User         `gorm:"foreignKey:UserID;references:ID" json:"-"`
	InvitedByUser User         `gorm:"foreignKey:InvitedBy;references:ID" json:"-"`
}

func (om *OrganizationMember) BeforeCreate(tx *gorm.DB) (err error) {
	if om.ID == uuid.Nil {
		om.ID = uuid.New()
	}
	return
}

type OrganizationInvitation struct {
	ID             uuid.UUID              `gorm:"type:uuid;primaryKey" json:"id"`
	OrganizationID uuid.UUID              `gorm:"type:uuid;not null" json:"organization_id"`
	Email          string                 `gorm:"type:varchar(255);not null" json:"email"`
	Role           OrganizationMemberRole `gorm:"type:varchar(20);default:member" json:"role"`
	Token          string                 `gorm:"type:varchar(255);uniqueIndex;not null" json:"token"`
	InvitedBy      uuid.UUID              `gorm:"type:uuid;not null" json:"invited_by"`
	InvitedAt      time.Time              `gorm:"default:now()" json:"invited_at"`
	ExpiresAt      time.Time              `json:"expires_at"`
	UsedAt         *time.Time             `json:"used_at"`
	Status         string                 `gorm:"type:varchar(20);default:pending" json:"status"` // pending, accepted, expired, cancelled

	// Relationships
	Organization  Organization `gorm:"foreignKey:OrganizationID;references:ID" json:"-"`
	InvitedByUser User         `gorm:"foreignKey:InvitedBy;references:ID" json:"-"`
}

func (oi *OrganizationInvitation) BeforeCreate(tx *gorm.DB) (err error) {
	if oi.ID == uuid.Nil {
		oi.ID = uuid.New()
	}
	return
}
