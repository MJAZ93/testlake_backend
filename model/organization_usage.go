package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrganizationUsage struct {
	ID                uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	OrganizationID    uuid.UUID `gorm:"type:uuid;not null" json:"organization_id"`
	PeriodStart       time.Time `gorm:"not null" json:"period_start"`
	PeriodEnd         time.Time `gorm:"not null" json:"period_end"`
	UsersCount        int       `gorm:"default:0" json:"users_count"`
	ProjectsCount     int       `gorm:"default:0" json:"projects_count"`
	EnvironmentsCount int       `gorm:"default:0" json:"environments_count"`
	SchemasCount      int       `gorm:"default:0" json:"schemas_count"`
	TestRecordsCount  int       `gorm:"default:0" json:"test_records_count"`
	APIRequestsCount  int       `gorm:"default:0" json:"api_requests_count"`
	RecordedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"recorded_at"`

	// Unique constraint on organization and period
	_ struct{} `gorm:"uniqueIndex:idx_organization_usage_period,column:organization_id,period_start,period_end"`

	// Relationships
	Organization Organization `gorm:"foreignKey:OrganizationID;references:ID" json:"-"`
}

func (ou *OrganizationUsage) BeforeCreate(tx *gorm.DB) (err error) {
	if ou.ID == uuid.Nil {
		ou.ID = uuid.New()
	}
	return
}
