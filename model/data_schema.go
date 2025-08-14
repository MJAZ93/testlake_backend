package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DataSchemaStatus string

const (
	DataSchemaStatusActive   DataSchemaStatus = "active"
	DataSchemaStatusArchived DataSchemaStatus = "archived"
)

type FieldType string

const (
	FieldTypeString    FieldType = "string"
	FieldTypeNumber    FieldType = "number"
	FieldTypeDate      FieldType = "date"
	FieldTypeBoolean   FieldType = "boolean"
	FieldTypeOptions   FieldType = "options"
	FieldTypeReference FieldType = "reference"
)

type DataSchema struct {
	ID               uuid.UUID        `gorm:"type:uuid;primaryKey" json:"id"`
	Name             string           `gorm:"type:varchar(200);not null" json:"name"`
	Description      *string          `gorm:"type:text" json:"description"`
	ProjectID        uuid.UUID        `gorm:"type:uuid;not null" json:"project_id"`
	IsReusable       bool             `gorm:"default:true" json:"is_reusable"`
	SchemaDefinition string           `gorm:"type:jsonb;not null" json:"schema_definition"`
	CreatedBy        uuid.UUID        `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
	Status           DataSchemaStatus `gorm:"type:varchar(20);default:active" json:"status"`
	DeletedAt        gorm.DeletedAt   `gorm:"index" json:"-"`

	// Relationships
	Project Project `gorm:"foreignKey:ProjectID;references:ID" json:"-"`
	Creator User    `gorm:"foreignKey:CreatedBy;references:ID" json:"-"`
}

func (ds *DataSchema) BeforeCreate(tx *gorm.DB) (err error) {
	if ds.ID == uuid.Nil {
		ds.ID = uuid.New()
	}
	return
}

type FeatureSchema struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	FeatureID uuid.UUID      `gorm:"type:uuid;not null" json:"feature_id"`
	SchemaID  uuid.UUID      `gorm:"type:uuid;not null" json:"schema_id"`
	IsPrimary bool           `gorm:"default:false" json:"is_primary"`
	CreatedBy uuid.UUID      `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Feature   Feature    `gorm:"foreignKey:FeatureID;references:ID" json:"-"`
	Schema    DataSchema `gorm:"foreignKey:SchemaID;references:ID" json:"-"`
	Creator   User       `gorm:"foreignKey:CreatedBy;references:ID" json:"-"`
}

func (fs *FeatureSchema) BeforeCreate(tx *gorm.DB) (err error) {
	if fs.ID == uuid.Nil {
		fs.ID = uuid.New()
	}
	return
}

type SchemaField struct {
	ID                  uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	SchemaID            uuid.UUID      `gorm:"type:uuid;not null" json:"schema_id"`
	FieldName           string         `gorm:"type:varchar(100);not null" json:"field_name"`
	FieldType           FieldType      `gorm:"type:varchar(20);not null" json:"field_type"`
	IsRequired          bool           `gorm:"default:false" json:"is_required"`
	ValidationRegex     *string        `gorm:"type:varchar(500)" json:"validation_regex"`
	MinValue            *string        `gorm:"type:varchar(100)" json:"min_value"`
	MaxValue            *string        `gorm:"type:varchar(100)" json:"max_value"`
	Options             *string        `gorm:"type:jsonb" json:"options"`
	ReferenceSchemaID   *uuid.UUID     `gorm:"type:uuid" json:"reference_schema_id"`
	DisplayOrder        int            `gorm:"default:0" json:"display_order"`
	CreatedAt           time.Time      `json:"created_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Schema          DataSchema  `gorm:"foreignKey:SchemaID;references:ID" json:"-"`
	ReferenceSchema *DataSchema `gorm:"foreignKey:ReferenceSchemaID;references:ID" json:"-"`
}

func (sf *SchemaField) BeforeCreate(tx *gorm.DB) (err error) {
	if sf.ID == uuid.Nil {
		sf.ID = uuid.New()
	}
	return
}