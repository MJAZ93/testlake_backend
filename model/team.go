package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TeamMemberRole string

const (
	TeamMemberRoleMember TeamMemberRole = "member"
	TeamMemberRoleAdmin  TeamMemberRole = "admin"
)

type Team struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Name           string         `gorm:"type:varchar(200);not null" json:"name"`
	Description    *string        `gorm:"type:text" json:"description"`
	OrganizationID uuid.UUID      `gorm:"type:uuid;not null" json:"organization_id"`
	CreatedBy      uuid.UUID      `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Organization Organization `gorm:"foreignKey:OrganizationID;references:ID" json:"-"`
	Creator      User         `gorm:"foreignKey:CreatedBy;references:ID" json:"-"`
}

func (t *Team) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return
}

type TeamMember struct {
	ID      uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	TeamID  uuid.UUID      `gorm:"type:uuid;not null" json:"team_id"`
	UserID  uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	Role    TeamMemberRole `gorm:"type:varchar(20);default:member" json:"role"`
	AddedBy uuid.UUID      `gorm:"type:uuid;not null" json:"added_by"`
	AddedAt time.Time      `gorm:"default:now()" json:"added_at"`

	// Relationships
	Team    Team `gorm:"foreignKey:TeamID;references:ID" json:"-"`
	User    User `gorm:"foreignKey:UserID;references:ID" json:"-"`
	AddedByUser User `gorm:"foreignKey:AddedBy;references:ID" json:"-"`
}

func (tm *TeamMember) BeforeCreate(tx *gorm.DB) (err error) {
	if tm.ID == uuid.Nil {
		tm.ID = uuid.New()
	}
	return
}