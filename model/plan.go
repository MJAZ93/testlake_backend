package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Plan struct {
	ID                      uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name                    string    `gorm:"type:varchar(100);not null" json:"name"`
	Slug                    string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"slug"`
	Description             *string   `gorm:"type:text" json:"description"`
	PriceMonthly            float64   `gorm:"type:decimal(10,2);not null" json:"price_monthly"`
	PriceYearly             float64   `gorm:"type:decimal(10,2);not null" json:"price_yearly"`
	MaxUsers                int       `gorm:"not null" json:"max_users"`
	MaxProjects             int       `gorm:"not null" json:"max_projects"`
	MaxEnvironments         int       `gorm:"not null" json:"max_environments"`
	MaxSchemas              int       `gorm:"not null" json:"max_schemas"`
	MaxTestRecordsPerSchema int       `gorm:"not null" json:"max_test_records_per_schema"`
	Features                string    `gorm:"type:jsonb;not null" json:"features"` // JSON array of enabled features
	PayPalMonthlyPlanID     *string   `gorm:"type:varchar(100)" json:"paypal_monthly_plan_id"`
	PayPalYearlyPlanID      *string   `gorm:"type:varchar(100)" json:"paypal_yearly_plan_id"`
	IsActive                bool      `gorm:"default:true" json:"is_active"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

func (p *Plan) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return
}
