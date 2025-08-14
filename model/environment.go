package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EnvironmentStatus string

const (
	EnvironmentStatusActive   EnvironmentStatus = "active"
	EnvironmentStatusArchived EnvironmentStatus = "archived"
)

type Environment struct {
	ID          uuid.UUID         `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string            `gorm:"type:varchar(100);not null" json:"name"`
	Slug        string            `gorm:"type:varchar(100);not null" json:"slug"`
	Description *string           `gorm:"type:text" json:"description"`
	Color       string            `gorm:"type:varchar(7);default:#3B82F6" json:"color"`
	ProjectID   uuid.UUID         `gorm:"type:uuid;not null" json:"project_id"`
	IsDefault   bool              `gorm:"default:false" json:"is_default"`
	CreatedBy   uuid.UUID         `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Status      EnvironmentStatus `gorm:"type:varchar(20);default:active" json:"status"`
	DeletedAt   gorm.DeletedAt    `gorm:"index" json:"-"`

	// Relationships
	Project Project `gorm:"foreignKey:ProjectID;references:ID" json:"-"`
	Creator User    `gorm:"foreignKey:CreatedBy;references:ID" json:"-"`
}

func (e *Environment) BeforeCreate(tx *gorm.DB) (err error) {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return
}