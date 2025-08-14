package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectStatus string

const (
	ProjectStatusActive   ProjectStatus = "active"
	ProjectStatusArchived ProjectStatus = "archived"
)

type Project struct {
	ID             uuid.UUID     `gorm:"type:uuid;primaryKey" json:"id"`
	Name           string        `gorm:"type:varchar(200);not null" json:"name"`
	Description    *string       `gorm:"type:text" json:"description"`
	OrganizationID *uuid.UUID    `gorm:"type:uuid" json:"organization_id"`
	UserID         *uuid.UUID    `gorm:"type:uuid" json:"user_id"`
	IsPersonal     bool          `gorm:"default:false" json:"is_personal"`
	CreatedBy      uuid.UUID     `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
	Status         ProjectStatus `gorm:"type:varchar(20);default:active" json:"status"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Organization *Organization `gorm:"foreignKey:OrganizationID;references:ID" json:"-"`
	User         *User         `gorm:"foreignKey:UserID;references:ID" json:"-"`
	Creator      User          `gorm:"foreignKey:CreatedBy;references:ID" json:"-"`
}

func (p *Project) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return
}