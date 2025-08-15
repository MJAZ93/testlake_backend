package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentMethodEnum string

const (
	PaymentMethodEnumPayPal PaymentMethodEnum = "paypal"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusCancelled PaymentStatus = "cancelled"
	PaymentStatusRefunded  PaymentStatus = "refunded"
)

type Payment struct {
	ID              uuid.UUID         `gorm:"type:uuid;primaryKey" json:"id"`
	OrganizationID  uuid.UUID         `gorm:"type:uuid;not null" json:"organization_id"`
	InvoiceID       *uuid.UUID        `gorm:"type:uuid" json:"invoice_id"`
	SubscriptionID  *uuid.UUID        `gorm:"type:uuid" json:"subscription_id"`
	PayPalPaymentID *string           `gorm:"type:varchar(100);uniqueIndex" json:"paypal_payment_id"`
	PayPalPayerID   *string           `gorm:"type:varchar(100)" json:"paypal_payer_id"`
	Amount          float64           `gorm:"type:decimal(10,2);not null" json:"amount"`
	Currency        string            `gorm:"type:varchar(3);default:USD" json:"currency"`
	PaymentMethod   PaymentMethodEnum `gorm:"type:varchar(20);default:paypal" json:"payment_method"`
	Status          PaymentStatus     `gorm:"type:varchar(20);default:pending" json:"status"`
	FailureReason   *string           `gorm:"type:text" json:"failure_reason"`
	ProcessedAt     *time.Time        `json:"processed_at"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`

	// Relationships
	Organization Organization  `gorm:"foreignKey:OrganizationID;references:ID" json:"-"`
	Invoice      *Invoice      `gorm:"foreignKey:InvoiceID;references:ID" json:"-"`
	Subscription *Subscription `gorm:"foreignKey:SubscriptionID;references:ID" json:"-"`
}

func (p *Payment) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return
}
