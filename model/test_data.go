package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TestDataStatus string

const (
	TestDataStatusActive  TestDataStatus = "active"
	TestDataStatusUsed    TestDataStatus = "used"
	TestDataStatusInvalid TestDataStatus = "invalid"
)

type TestDataRequestStatus string

const (
	TestDataRequestStatusPending   TestDataRequestStatus = "pending"
	TestDataRequestStatusFulfilled TestDataRequestStatus = "fulfilled"
	TestDataRequestStatusRejected  TestDataRequestStatus = "rejected"
)

type TestData struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	SchemaID      uuid.UUID      `gorm:"type:uuid;not null" json:"schema_id"`
	EnvironmentID uuid.UUID      `gorm:"type:uuid;not null" json:"environment_id"`
	DataValues    string         `gorm:"type:jsonb;not null" json:"data_values"`
	IsUsed        bool           `gorm:"default:false" json:"is_used"`
	UsedAt        *time.Time     `json:"used_at"`
	UsedBy        *uuid.UUID     `gorm:"type:uuid" json:"used_by"`
	FeatureID     *uuid.UUID     `gorm:"type:uuid" json:"feature_id"`
	CreatedBy     uuid.UUID      `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	Status        TestDataStatus `gorm:"type:varchar(20);default:active" json:"status"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Schema      DataSchema  `gorm:"foreignKey:SchemaID;references:ID" json:"-"`
	Environment Environment `gorm:"foreignKey:EnvironmentID;references:ID" json:"-"`
	UsedByUser  *User       `gorm:"foreignKey:UsedBy;references:ID" json:"-"`
	Feature     *Feature    `gorm:"foreignKey:FeatureID;references:ID" json:"-"`
	Creator     User        `gorm:"foreignKey:CreatedBy;references:ID" json:"-"`
}

func (td *TestData) BeforeCreate(tx *gorm.DB) (err error) {
	if td.ID == uuid.Nil {
		td.ID = uuid.New()
	}
	return
}

type TestDataRequest struct {
	ID              uuid.UUID             `gorm:"type:uuid;primaryKey" json:"id"`
	FeatureID       uuid.UUID             `gorm:"type:uuid;not null" json:"feature_id"`
	EnvironmentID   uuid.UUID             `gorm:"type:uuid;not null" json:"environment_id"`
	SchemaID        uuid.UUID             `gorm:"type:uuid;not null" json:"schema_id"`
	RequestedBy     uuid.UUID             `gorm:"type:uuid;not null" json:"requested_by"`
	ProvidedDataID  *uuid.UUID            `gorm:"type:uuid" json:"provided_data_id"`
	RequestNotes    *string               `gorm:"type:text" json:"request_notes"`
	ResponseNotes   *string               `gorm:"type:text" json:"response_notes"`
	RequestedAt     time.Time             `gorm:"default:now()" json:"requested_at"`
	FulfilledAt     *time.Time            `json:"fulfilled_at"`
	Status          TestDataRequestStatus `gorm:"type:varchar(20);default:pending" json:"status"`
	DeletedAt       gorm.DeletedAt        `gorm:"index" json:"-"`

	// Relationships
	Feature      Feature     `gorm:"foreignKey:FeatureID;references:ID" json:"-"`
	Environment  Environment `gorm:"foreignKey:EnvironmentID;references:ID" json:"-"`
	Schema       DataSchema  `gorm:"foreignKey:SchemaID;references:ID" json:"-"`
	Requester    User        `gorm:"foreignKey:RequestedBy;references:ID" json:"-"`
	ProvidedData *TestData   `gorm:"foreignKey:ProvidedDataID;references:ID" json:"-"`
}

func (tdr *TestDataRequest) BeforeCreate(tx *gorm.DB) (err error) {
	if tdr.ID == uuid.Nil {
		tdr.ID = uuid.New()
	}
	return
}