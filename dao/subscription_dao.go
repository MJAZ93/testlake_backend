package dao

import (
	"testlake/model"

	"github.com/google/uuid"
)

type SubscriptionDao struct {
	Limit int
}

func NewSubscriptionDao() *SubscriptionDao {
	return &SubscriptionDao{Limit: 50}
}

func (dao *SubscriptionDao) Create(subscription *model.Subscription) error {
	return Database.Create(subscription).Error
}

func (dao *SubscriptionDao) GetByID(id uuid.UUID) (*model.Subscription, error) {
	var subscription model.Subscription
	err := Database.Preload("Organization").Preload("Plan").First(&subscription, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

func (dao *SubscriptionDao) GetByOrganizationID(organizationID uuid.UUID) (*model.Subscription, error) {
	var subscription model.Subscription
	err := Database.Preload("Organization").Preload("Plan").
		Where("organization_id = ?", organizationID).
		Order("created_at DESC").
		First(&subscription).Error
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

func (dao *SubscriptionDao) GetByPayPalSubscriptionID(paypalSubID string) (*model.Subscription, error) {
	var subscription model.Subscription
	err := Database.Preload("Organization").Preload("Plan").
		First(&subscription, "paypal_subscription_id = ?", paypalSubID).Error
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

func (dao *SubscriptionDao) GetActiveByOrganizationID(organizationID uuid.UUID) (*model.Subscription, error) {
	var subscription model.Subscription
	err := Database.Preload("Organization").Preload("Plan").
		Where("organization_id = ? AND status = ?", organizationID, model.SubscriptionStatusActive).
		Order("created_at DESC").
		First(&subscription).Error
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

func (dao *SubscriptionDao) Update(subscription *model.Subscription) error {
	return Database.Save(subscription).Error
}

func (dao *SubscriptionDao) UpdateStatus(id uuid.UUID, status model.SubscriptionStatus) error {
	return Database.Model(&model.Subscription{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (dao *SubscriptionDao) Cancel(id uuid.UUID, cancelAtPeriodEnd bool) error {
	updates := map[string]interface{}{
		"cancel_at_period_end": cancelAtPeriodEnd,
	}
	if !cancelAtPeriodEnd {
		updates["status"] = model.SubscriptionStatusCancelled
		updates["cancelled_at"] = "NOW()"
	}
	return Database.Model(&model.Subscription{}).Where("id = ?", id).Updates(updates).Error
}

func (dao *SubscriptionDao) GetAll(page int) ([]model.Subscription, int64, error) {
	var subscriptions []model.Subscription
	var total int64

	err := Database.Model(&model.Subscription{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := page * dao.Limit
	err = Database.Preload("Organization").Preload("Plan").
		Offset(offset).Limit(dao.Limit).
		Find(&subscriptions).Error
	if err != nil {
		return nil, 0, err
	}

	return subscriptions, total, nil
}

func (dao *SubscriptionDao) Delete(id uuid.UUID) error {
	return Database.Delete(&model.Subscription{}, "id = ?", id).Error
}
