package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PlanType string

const (
	PlanTypeFree         PlanType = "free"
	PlanTypeStarter      PlanType = "starter"
	PlanTypeProfessional PlanType = "professional"
	PlanTypeEnterprise   PlanType = "enterprise"
)

type OrganizationStatus string

const (
	OrganizationStatusActive    OrganizationStatus = "active"
	OrganizationStatusSuspended OrganizationStatus = "suspended"
	OrganizationStatusCancelled OrganizationStatus = "cancelled"
)

type OrganizationSubscriptionStatus string

const (
	OrganizationSubscriptionStatusActive    OrganizationSubscriptionStatus = "active"
	OrganizationSubscriptionStatusPastDue   OrganizationSubscriptionStatus = "past_due"
	OrganizationSubscriptionStatusCancelled OrganizationSubscriptionStatus = "cancelled"
	OrganizationSubscriptionStatusSuspended OrganizationSubscriptionStatus = "suspended"
	OrganizationSubscriptionStatusTrialing  OrganizationSubscriptionStatus = "trialing"
)

type Organization struct {
	ID          uuid.UUID          `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string             `gorm:"type:varchar(200);not null" json:"name"`
	Slug        string             `gorm:"type:varchar(100);uniqueIndex;not null" json:"slug"`
	Description *string            `gorm:"type:text" json:"description"`
	LogoURL     *string            `gorm:"type:varchar(500)" json:"logo_url"`
	PlanType    PlanType           `gorm:"type:varchar(20);default:starter" json:"plan_type"`
	MaxUsers    int                `gorm:"default:10" json:"max_users"`
	MaxProjects int                `gorm:"default:5" json:"max_projects"`
	CreatedBy   uuid.UUID          `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	Status      OrganizationStatus `gorm:"type:varchar(20);default:active" json:"status"`
	DeletedAt   gorm.DeletedAt     `gorm:"index" json:"-"`

	// Payment-related fields
	PlanID               *uuid.UUID                     `gorm:"type:uuid" json:"plan_id"`
	BillingCycle         BillingCycle                   `gorm:"type:varchar(20);default:monthly" json:"billing_cycle"`
	SubscriptionStatus   OrganizationSubscriptionStatus `gorm:"type:varchar(20);default:active" json:"subscription_status"`
	TrialEndsAt          *time.Time                     `json:"trial_ends_at"`
	NextBillingDate      *time.Time                     `json:"next_billing_date"`
	PayPalSubscriptionID *string                        `gorm:"type:varchar(100)" json:"paypal_subscription_id"`
	BillingEmail         *string                        `gorm:"type:varchar(255)" json:"billing_email"`

	// Relationships
	Creator User  `gorm:"foreignKey:CreatedBy;references:ID" json:"-"`
	Plan    *Plan `gorm:"foreignKey:PlanID;references:ID" json:"-"`
}

func (o *Organization) BeforeCreate(tx *gorm.DB) (err error) {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return
}
