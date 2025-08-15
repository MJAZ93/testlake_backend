package dao

import (
	"testlake/model"

	"github.com/google/uuid"
)

type BillingEventDao struct {
	Limit int
}

func NewBillingEventDao() *BillingEventDao {
	return &BillingEventDao{Limit: 50}
}

func (dao *BillingEventDao) Create(event *model.BillingEvent) error {
	return Database.Create(event).Error
}

func (dao *BillingEventDao) GetByID(id uuid.UUID) (*model.BillingEvent, error) {
	var event model.BillingEvent
	err := Database.First(&event, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (dao *BillingEventDao) GetByOrganizationID(organizationID uuid.UUID, page int) ([]model.BillingEvent, int64, error) {
	var events []model.BillingEvent
	var total int64

	err := Database.Model(&model.BillingEvent{}).
		Where("organization_id = ?", organizationID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := page * dao.Limit
	err = Database.Where("organization_id = ?", organizationID).
		Order("created_at DESC").
		Offset(offset).Limit(dao.Limit).
		Find(&events).Error
	if err != nil {
		return nil, 0, err
	}

	return events, total, nil
}

func (dao *BillingEventDao) GetByPayPalEventID(paypalEventID string) (*model.BillingEvent, error) {
	var event model.BillingEvent
	err := Database.First(&event, "paypal_event_id = ?", paypalEventID).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (dao *BillingEventDao) GetByEventType(organizationID uuid.UUID, eventType model.BillingEventType, page int) ([]model.BillingEvent, int64, error) {
	var events []model.BillingEvent
	var total int64

	err := Database.Model(&model.BillingEvent{}).
		Where("organization_id = ? AND event_type = ?", organizationID, eventType).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := page * dao.Limit
	err = Database.Where("organization_id = ? AND event_type = ?", organizationID, eventType).
		Order("created_at DESC").
		Offset(offset).Limit(dao.Limit).
		Find(&events).Error
	if err != nil {
		return nil, 0, err
	}

	return events, total, nil
}

func (dao *BillingEventDao) Update(event *model.BillingEvent) error {
	return Database.Save(event).Error
}

func (dao *BillingEventDao) Delete(id uuid.UUID) error {
	return Database.Delete(&model.BillingEvent{}, "id = ?", id).Error
}

func (dao *BillingEventDao) GetAll(page int) ([]model.BillingEvent, int64, error) {
	var events []model.BillingEvent
	var total int64

	err := Database.Model(&model.BillingEvent{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := page * dao.Limit
	err = Database.Preload("Organization").
		Order("created_at DESC").
		Offset(offset).Limit(dao.Limit).
		Find(&events).Error
	if err != nil {
		return nil, 0, err
	}

	return events, total, nil
}

func (dao *BillingEventDao) EventExists(paypalEventID string) (bool, error) {
	var count int64
	err := Database.Model(&model.BillingEvent{}).
		Where("paypal_event_id = ?", paypalEventID).
		Count(&count).Error
	return count > 0, err
}
