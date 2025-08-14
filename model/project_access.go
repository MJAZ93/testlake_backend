package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Permission string

const (
	PermissionRead  Permission = "read"
	PermissionWrite Permission = "write"
	PermissionAdmin Permission = "admin"
)

type ProjectAccess struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	ProjectID uuid.UUID      `gorm:"type:uuid;not null" json:"project_id"`
	TeamID    *uuid.UUID     `gorm:"type:uuid" json:"team_id"`
	UserID    *uuid.UUID     `gorm:"type:uuid" json:"user_id"`
	Permission Permission    `gorm:"type:varchar(20);default:read" json:"permission"`
	GrantedBy uuid.UUID      `gorm:"type:uuid;not null" json:"granted_by"`
	GrantedAt time.Time      `gorm:"default:now()" json:"granted_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Project     Project `gorm:"foreignKey:ProjectID;references:ID" json:"-"`
	Team        *Team   `gorm:"foreignKey:TeamID;references:ID" json:"-"`
	User        *User   `gorm:"foreignKey:UserID;references:ID" json:"-"`
	GrantedByUser User  `gorm:"foreignKey:GrantedBy;references:ID" json:"-"`
}

func (pa *ProjectAccess) BeforeCreate(tx *gorm.DB) (err error) {
	if pa.ID == uuid.Nil {
		pa.ID = uuid.New()
	}
	return
}