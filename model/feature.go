package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Feature struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(200);not null" json:"name"`
	Description *string        `gorm:"type:text" json:"description"`
	ProjectID   uuid.UUID      `gorm:"type:uuid;not null" json:"project_id"`
	CreatedBy   uuid.UUID      `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Project Project `gorm:"foreignKey:ProjectID;references:ID" json:"-"`
	Creator User    `gorm:"foreignKey:CreatedBy;references:ID" json:"-"`
}

func (f *Feature) BeforeCreate(tx *gorm.DB) (err error) {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return
}

type FeatureEnvironmentStatus struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	FeatureID     uuid.UUID      `gorm:"type:uuid;not null" json:"feature_id"`
	EnvironmentID uuid.UUID      `gorm:"type:uuid;not null" json:"environment_id"`
	IsWorking     bool           `gorm:"default:true" json:"is_working"`
	ErrorMessage  *string        `gorm:"type:text" json:"error_message"`
	LastTestedAt  *time.Time     `json:"last_tested_at"`
	LastTestedBy  *uuid.UUID     `gorm:"type:uuid" json:"last_tested_by"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Feature       Feature     `gorm:"foreignKey:FeatureID;references:ID" json:"-"`
	Environment   Environment `gorm:"foreignKey:EnvironmentID;references:ID" json:"-"`
	LastTestedByUser *User    `gorm:"foreignKey:LastTestedBy;references:ID" json:"-"`
}

func (fes *FeatureEnvironmentStatus) BeforeCreate(tx *gorm.DB) (err error) {
	if fes.ID == uuid.Nil {
		fes.ID = uuid.New()
	}
	return
}

type FeatureErrorLog struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	FeatureID     uuid.UUID      `gorm:"type:uuid;not null" json:"feature_id"`
	EnvironmentID uuid.UUID      `gorm:"type:uuid;not null" json:"environment_id"`
	ErrorMessage  string         `gorm:"type:text;not null" json:"error_message"`
	ErrorDetails  *string        `gorm:"type:jsonb" json:"error_details"`
	ReportedBy    uuid.UUID      `gorm:"type:uuid;not null" json:"reported_by"`
	ReportedAt    time.Time      `gorm:"default:now()" json:"reported_at"`
	ResolvedAt    *time.Time     `json:"resolved_at"`
	ResolvedBy    *uuid.UUID     `gorm:"type:uuid" json:"resolved_by"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Feature       Feature     `gorm:"foreignKey:FeatureID;references:ID" json:"-"`
	Environment   Environment `gorm:"foreignKey:EnvironmentID;references:ID" json:"-"`
	Reporter      User        `gorm:"foreignKey:ReportedBy;references:ID" json:"-"`
	Resolver      *User       `gorm:"foreignKey:ResolvedBy;references:ID" json:"-"`
}

func (fel *FeatureErrorLog) BeforeCreate(tx *gorm.DB) (err error) {
	if fel.ID == uuid.Nil {
		fel.ID = uuid.New()
	}
	return
}

type ErrorImage struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	ErrorLogID  uuid.UUID      `gorm:"type:uuid;not null" json:"error_log_id"`
	ImageURL    string         `gorm:"type:varchar(500);not null" json:"image_url"`
	ImageName   *string        `gorm:"type:varchar(200)" json:"image_name"`
	UploadedAt  time.Time      `gorm:"default:now()" json:"uploaded_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	ErrorLog FeatureErrorLog `gorm:"foreignKey:ErrorLogID;references:ID" json:"-"`
}

func (ei *ErrorImage) BeforeCreate(tx *gorm.DB) (err error) {
	if ei.ID == uuid.Nil {
		ei.ID = uuid.New()
	}
	return
}