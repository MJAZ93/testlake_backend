package dao

import (
	"testlake/model"

	"github.com/google/uuid"
)

type PaymentDao struct {
	Limit int
}

func NewPaymentDao() *PaymentDao {
	return &PaymentDao{Limit: 50}
}

func (dao *PaymentDao) Create(payment *model.Payment) error {
	return Database.Create(payment).Error
}

func (dao *PaymentDao) GetByID(id uuid.UUID) (*model.Payment, error) {
	var payment model.Payment
	err := Database.Preload("Organization").
		Preload("Invoice").
		Preload("Subscription").
		First(&payment, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (dao *PaymentDao) GetByPayPalPaymentID(paypalPaymentID string) (*model.Payment, error) {
	var payment model.Payment
	err := Database.Preload("Organization").
		Preload("Invoice").
		Preload("Subscription").
		First(&payment, "paypal_payment_id = ?", paypalPaymentID).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (dao *PaymentDao) GetByOrganizationID(organizationID uuid.UUID, page int) ([]model.Payment, int64, error) {
	var payments []model.Payment
	var total int64

	err := Database.Model(&model.Payment{}).
		Where("organization_id = ?", organizationID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := page * dao.Limit
	err = Database.Preload("Invoice").
		Preload("Subscription").
		Where("organization_id = ?", organizationID).
		Order("created_at DESC").
		Offset(offset).Limit(dao.Limit).
		Find(&payments).Error
	if err != nil {
		return nil, 0, err
	}

	return payments, total, nil
}

func (dao *PaymentDao) GetRecentByOrganizationID(organizationID uuid.UUID, limit int) ([]model.Payment, error) {
	var payments []model.Payment
	err := Database.Preload("Invoice").
		Preload("Subscription").
		Where("organization_id = ?", organizationID).
		Order("created_at DESC").
		Limit(limit).
		Find(&payments).Error
	if err != nil {
		return nil, err
	}
	return payments, nil
}

func (dao *PaymentDao) GetByInvoiceID(invoiceID uuid.UUID) ([]model.Payment, error) {
	var payments []model.Payment
	err := Database.Where("invoice_id = ?", invoiceID).
		Order("created_at DESC").
		Find(&payments).Error
	if err != nil {
		return nil, err
	}
	return payments, nil
}

func (dao *PaymentDao) GetBySubscriptionID(subscriptionID uuid.UUID) ([]model.Payment, error) {
	var payments []model.Payment
	err := Database.Where("subscription_id = ?", subscriptionID).
		Order("created_at DESC").
		Find(&payments).Error
	if err != nil {
		return nil, err
	}
	return payments, nil
}

func (dao *PaymentDao) Update(payment *model.Payment) error {
	return Database.Save(payment).Error
}

func (dao *PaymentDao) UpdateStatus(id uuid.UUID, status model.PaymentStatus) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if status == model.PaymentStatusCompleted {
		updates["processed_at"] = "NOW()"
	}
	return Database.Model(&model.Payment{}).Where("id = ?", id).Updates(updates).Error
}

func (dao *PaymentDao) UpdateFailureReason(id uuid.UUID, reason string) error {
	return Database.Model(&model.Payment{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":         model.PaymentStatusFailed,
			"failure_reason": reason,
		}).Error
}

func (dao *PaymentDao) Delete(id uuid.UUID) error {
	return Database.Delete(&model.Payment{}, "id = ?", id).Error
}

func (dao *PaymentDao) GetAll(page int) ([]model.Payment, int64, error) {
	var payments []model.Payment
	var total int64

	err := Database.Model(&model.Payment{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := page * dao.Limit
	err = Database.Preload("Organization").
		Preload("Invoice").
		Preload("Subscription").
		Order("created_at DESC").
		Offset(offset).Limit(dao.Limit).
		Find(&payments).Error
	if err != nil {
		return nil, 0, err
	}

	return payments, total, nil
}
