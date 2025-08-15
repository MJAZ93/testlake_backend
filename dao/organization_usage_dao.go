package dao

import (
	"testlake/model"
	"time"

	"github.com/google/uuid"
)

type OrganizationUsageDao struct {
	Limit int
}

func NewOrganizationUsageDao() *OrganizationUsageDao {
	return &OrganizationUsageDao{Limit: 50}
}

func (dao *OrganizationUsageDao) Create(usage *model.OrganizationUsage) error {
	return Database.Create(usage).Error
}

func (dao *OrganizationUsageDao) GetByID(id uuid.UUID) (*model.OrganizationUsage, error) {
	var usage model.OrganizationUsage
	err := Database.First(&usage, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &usage, nil
}

func (dao *OrganizationUsageDao) GetCurrentUsage(organizationID uuid.UUID) (*model.OrganizationUsage, error) {
	var usage model.OrganizationUsage
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	err := Database.Where("organization_id = ? AND period_start = ? AND period_end = ?",
		organizationID, startOfMonth, endOfMonth).
		First(&usage).Error
	if err != nil {
		return nil, err
	}
	return &usage, nil
}

func (dao *OrganizationUsageDao) GetByOrganizationAndPeriod(organizationID uuid.UUID, periodStart, periodEnd time.Time) (*model.OrganizationUsage, error) {
	var usage model.OrganizationUsage
	err := Database.Where("organization_id = ? AND period_start = ? AND period_end = ?",
		organizationID, periodStart, periodEnd).
		First(&usage).Error
	if err != nil {
		return nil, err
	}
	return &usage, nil
}

func (dao *OrganizationUsageDao) GetUsageHistory(organizationID uuid.UUID, page int) ([]model.OrganizationUsage, int64, error) {
	var usages []model.OrganizationUsage
	var total int64

	err := Database.Model(&model.OrganizationUsage{}).
		Where("organization_id = ?", organizationID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := page * dao.Limit
	err = Database.Where("organization_id = ?", organizationID).
		Order("period_start DESC").
		Offset(offset).Limit(dao.Limit).
		Find(&usages).Error
	if err != nil {
		return nil, 0, err
	}

	return usages, total, nil
}

func (dao *OrganizationUsageDao) Update(usage *model.OrganizationUsage) error {
	return Database.Save(usage).Error
}

func (dao *OrganizationUsageDao) UpsertCurrentUsage(organizationID uuid.UUID, updates map[string]interface{}) error {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	var usage model.OrganizationUsage
	err := Database.Where("organization_id = ? AND period_start = ? AND period_end = ?",
		organizationID, startOfMonth, endOfMonth).
		First(&usage).Error

	if err != nil {
		// Create new usage record
		usage = model.OrganizationUsage{
			OrganizationID: organizationID,
			PeriodStart:    startOfMonth,
			PeriodEnd:      endOfMonth,
			RecordedAt:     now,
		}

		// Apply updates
		for key, value := range updates {
			switch key {
			case "users_count":
				if count, ok := value.(int); ok {
					usage.UsersCount = count
				}
			case "projects_count":
				if count, ok := value.(int); ok {
					usage.ProjectsCount = count
				}
			case "environments_count":
				if count, ok := value.(int); ok {
					usage.EnvironmentsCount = count
				}
			case "schemas_count":
				if count, ok := value.(int); ok {
					usage.SchemasCount = count
				}
			case "test_records_count":
				if count, ok := value.(int); ok {
					usage.TestRecordsCount = count
				}
			case "api_requests_count":
				if count, ok := value.(int); ok {
					usage.APIRequestsCount = count
				}
			}
		}

		return Database.Create(&usage).Error
	}

	// Update existing record
	updates["recorded_at"] = now
	return Database.Model(&usage).Updates(updates).Error
}

func (dao *OrganizationUsageDao) IncrementUsage(organizationID uuid.UUID, field string, increment int) error {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	// First try to update existing record
	result := Database.Model(&model.OrganizationUsage{}).
		Where("organization_id = ? AND period_start = ? AND period_end = ?",
			organizationID, startOfMonth, endOfMonth).
		UpdateColumn(field, Database.Raw(field+" + ?", increment))

	if result.RowsAffected == 0 {
		// Create new record if none exists
		usage := model.OrganizationUsage{
			OrganizationID: organizationID,
			PeriodStart:    startOfMonth,
			PeriodEnd:      endOfMonth,
			RecordedAt:     now,
		}

		switch field {
		case "users_count":
			usage.UsersCount = increment
		case "projects_count":
			usage.ProjectsCount = increment
		case "environments_count":
			usage.EnvironmentsCount = increment
		case "schemas_count":
			usage.SchemasCount = increment
		case "test_records_count":
			usage.TestRecordsCount = increment
		case "api_requests_count":
			usage.APIRequestsCount = increment
		}

		return Database.Create(&usage).Error
	}

	return result.Error
}

func (dao *OrganizationUsageDao) Delete(id uuid.UUID) error {
	return Database.Delete(&model.OrganizationUsage{}, "id = ?", id).Error
}
