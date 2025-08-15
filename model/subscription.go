package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BillingCycle string

const (
	BillingCycleMonthly BillingCycle = "monthly"
	BillingCycleYearly  BillingCycle = "yearly"
)

type SubscriptionStatus string

const (
	SubscriptionStatusActive    SubscriptionStatus = "active"
	SubscriptionStatusCancelled SubscriptionStatus = "cancelled"
	SubscriptionStatusSuspended SubscriptionStatus = "suspended"
	SubscriptionStatusExpired   SubscriptionStatus = "expired"
	SubscriptionStatusPending   SubscriptionStatus = "pending"
)

type Subscription struct {
	ID                   uuid.UUID          `gorm:"type:uuid;primaryKey" json:"id"`
	OrganizationID       uuid.UUID          `gorm:"type:uuid;not null" json:"organization_id"`
	PlanID               uuid.UUID          `gorm:"type:uuid;not null" json:"plan_id"`
	PayPalSubscriptionID string             `gorm:"type:varchar(100);uniqueIndex;not null" json:"paypal_subscription_id"`
	Status               SubscriptionStatus `gorm:"type:varchar(20);default:pending" json:"status"`
	BillingCycle         BillingCycle       `gorm:"type:varchar(20);not null" json:"billing_cycle"`
	CurrentPeriodStart   time.Time          `gorm:"not null" json:"current_period_start"`
	CurrentPeriodEnd     time.Time          `gorm:"not null" json:"current_period_end"`
	TrialEnd             *time.Time         `json:"trial_end"`
	CancelAtPeriodEnd    bool               `gorm:"default:false" json:"cancel_at_period_end"`
	CancelledAt          *time.Time         `json:"cancelled_at"`
	CreatedBy            uuid.UUID          `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt            time.Time          `json:"created_at"`
	UpdatedAt            time.Time          `json:"updated_at"`

	// Relationships
	Organization Organization `gorm:"foreignKey:OrganizationID;references:ID" json:"-"`
	Plan         Plan         `gorm:"foreignKey:PlanID;references:ID" json:"-"`
	Creator      User         `gorm:"foreignKey:CreatedBy;references:ID" json:"-"`
}

func (s *Subscription) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return
}
