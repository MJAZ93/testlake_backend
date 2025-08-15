package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BillingEventType string

const (
	BillingEventTypeSubscriptionCreated   BillingEventType = "subscription_created"
	BillingEventTypeSubscriptionUpdated   BillingEventType = "subscription_updated"
	BillingEventTypeSubscriptionCancelled BillingEventType = "subscription_cancelled"
	BillingEventTypePaymentSucceeded      BillingEventType = "payment_succeeded"
	BillingEventTypePaymentFailed         BillingEventType = "payment_failed"
	BillingEventTypeInvoiceCreated        BillingEventType = "invoice_created"
	BillingEventTypePlanChanged           BillingEventType = "plan_changed"
)

type BillingEvent struct {
	ID             uuid.UUID        `gorm:"type:uuid;primaryKey" json:"id"`
	OrganizationID uuid.UUID        `gorm:"type:uuid;not null" json:"organization_id"`
	EventType      BillingEventType `gorm:"type:varchar(50);not null" json:"event_type"`
	EventData      *string          `gorm:"type:jsonb" json:"event_data"`
	PayPalEventID  *string          `gorm:"type:varchar(100)" json:"paypal_event_id"`
	ProcessedAt    time.Time        `gorm:"default:CURRENT_TIMESTAMP" json:"processed_at"`
	CreatedAt      time.Time        `json:"created_at"`

	// Relationships
	Organization Organization `gorm:"foreignKey:OrganizationID;references:ID" json:"-"`
}

func (be *BillingEvent) BeforeCreate(tx *gorm.DB) (err error) {
	if be.ID == uuid.Nil {
		be.ID = uuid.New()
	}
	return
}
