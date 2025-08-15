package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InvoiceStatus string

const (
	InvoiceStatusDraft     InvoiceStatus = "draft"
	InvoiceStatusSent      InvoiceStatus = "sent"
	InvoiceStatusPaid      InvoiceStatus = "paid"
	InvoiceStatusCancelled InvoiceStatus = "cancelled"
	InvoiceStatusRefunded  InvoiceStatus = "refunded"
)

type Invoice struct {
	ID                 uuid.UUID     `gorm:"type:uuid;primaryKey" json:"id"`
	OrganizationID     uuid.UUID     `gorm:"type:uuid;not null" json:"organization_id"`
	SubscriptionID     *uuid.UUID    `gorm:"type:uuid" json:"subscription_id"`
	PayPalInvoiceID    *string       `gorm:"type:varchar(100)" json:"paypal_invoice_id"`
	InvoiceNumber      string        `gorm:"type:varchar(50);uniqueIndex;not null" json:"invoice_number"`
	Amount             float64       `gorm:"type:decimal(10,2);not null" json:"amount"`
	TaxAmount          float64       `gorm:"type:decimal(10,2);default:0" json:"tax_amount"`
	TotalAmount        float64       `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	Currency           string        `gorm:"type:varchar(3);default:USD" json:"currency"`
	Status             InvoiceStatus `gorm:"type:varchar(20);default:draft" json:"status"`
	BillingPeriodStart *time.Time    `json:"billing_period_start"`
	BillingPeriodEnd   *time.Time    `json:"billing_period_end"`
	DueDate            *time.Time    `json:"due_date"`
	PaidAt             *time.Time    `json:"paid_at"`
	InvoiceURL         *string       `gorm:"type:varchar(500)" json:"invoice_url"`
	CreatedAt          time.Time     `json:"created_at"`
	UpdatedAt          time.Time     `json:"updated_at"`

	// Relationships
	Organization Organization      `gorm:"foreignKey:OrganizationID;references:ID" json:"-"`
	Subscription *Subscription     `gorm:"foreignKey:SubscriptionID;references:ID" json:"-"`
	LineItems    []InvoiceLineItem `gorm:"foreignKey:InvoiceID;references:ID" json:"line_items,omitempty"`
}

type InvoiceLineItem struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	InvoiceID   uuid.UUID `gorm:"type:uuid;not null" json:"invoice_id"`
	Description string    `gorm:"type:varchar(255);not null" json:"description"`
	Quantity    int       `gorm:"default:1" json:"quantity"`
	UnitPrice   float64   `gorm:"type:decimal(10,2);not null" json:"unit_price"`
	TotalPrice  float64   `gorm:"type:decimal(10,2);not null" json:"total_price"`
	CreatedAt   time.Time `json:"created_at"`

	// Relationships
	Invoice Invoice `gorm:"foreignKey:InvoiceID;references:ID" json:"-"`
}

func (i *Invoice) BeforeCreate(tx *gorm.DB) (err error) {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return
}

func (ili *InvoiceLineItem) BeforeCreate(tx *gorm.DB) (err error) {
	if ili.ID == uuid.Nil {
		ili.ID = uuid.New()
	}
	return
}
