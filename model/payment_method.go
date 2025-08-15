package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentMethodType string

const (
	PaymentMethodTypePayPal PaymentMethodType = "paypal"
)

type PaymentMethod struct {
	ID                uuid.UUID         `gorm:"type:uuid;primaryKey" json:"id"`
	OrganizationID    uuid.UUID         `gorm:"type:uuid;not null" json:"organization_id"`
	PayPalPayerID     *string           `gorm:"type:varchar(100)" json:"paypal_payer_id"`
	PayPalEmail       *string           `gorm:"type:varchar(255)" json:"paypal_email"`
	PaymentMethodType PaymentMethodType `gorm:"type:varchar(20);default:paypal" json:"payment_method_type"`
	IsDefault         bool              `gorm:"default:false" json:"is_default"`
	IsActive          bool              `gorm:"default:true" json:"is_active"`
	CreatedBy         uuid.UUID         `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`

	// Relationships
	Organization Organization `gorm:"foreignKey:OrganizationID;references:ID" json:"-"`
	Creator      User         `gorm:"foreignKey:CreatedBy;references:ID" json:"-"`
}

func (pm *PaymentMethod) BeforeCreate(tx *gorm.DB) (err error) {
	if pm.ID == uuid.Nil {
		pm.ID = uuid.New()
	}
	return
}
