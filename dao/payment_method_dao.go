package dao

import (
	"testlake/model"

	"github.com/google/uuid"
)

type PaymentMethodDao struct {
	Limit int
}

func NewPaymentMethodDao() *PaymentMethodDao {
	return &PaymentMethodDao{Limit: 50}
}

func (dao *PaymentMethodDao) Create(paymentMethod *model.PaymentMethod) error {
	return Database.Create(paymentMethod).Error
}

func (dao *PaymentMethodDao) GetByID(id uuid.UUID) (*model.PaymentMethod, error) {
	var paymentMethod model.PaymentMethod
	err := Database.First(&paymentMethod, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &paymentMethod, nil
}

func (dao *PaymentMethodDao) GetByOrganizationID(organizationID uuid.UUID) ([]model.PaymentMethod, error) {
	var paymentMethods []model.PaymentMethod
	err := Database.Where("organization_id = ? AND is_active = ?", organizationID, true).Find(&paymentMethods).Error
	if err != nil {
		return nil, err
	}
	return paymentMethods, nil
}

func (dao *PaymentMethodDao) GetDefaultByOrganizationID(organizationID uuid.UUID) (*model.PaymentMethod, error) {
	var paymentMethod model.PaymentMethod
	err := Database.Where("organization_id = ? AND is_default = ? AND is_active = ?", organizationID, true, true).First(&paymentMethod).Error
	if err != nil {
		return nil, err
	}
	return &paymentMethod, nil
}

func (dao *PaymentMethodDao) Update(paymentMethod *model.PaymentMethod) error {
	return Database.Save(paymentMethod).Error
}

func (dao *PaymentMethodDao) Delete(id uuid.UUID) error {
	return Database.Delete(&model.PaymentMethod{}, "id = ?", id).Error
}

func (dao *PaymentMethodDao) SetActive(id uuid.UUID, active bool) error {
	return Database.Model(&model.PaymentMethod{}).Where("id = ?", id).Update("is_active", active).Error
}

func (dao *PaymentMethodDao) SetDefault(organizationID, paymentMethodID uuid.UUID) error {
	// First, unset all default payment methods for this organization
	err := Database.Model(&model.PaymentMethod{}).
		Where("organization_id = ?", organizationID).
		Update("is_default", false).Error
	if err != nil {
		return err
	}

	// Then set the specified payment method as default
	return Database.Model(&model.PaymentMethod{}).
		Where("id = ? AND organization_id = ?", paymentMethodID, organizationID).
		Update("is_default", true).Error
}

func (dao *PaymentMethodDao) CountByOrganizationID(organizationID uuid.UUID) (int64, error) {
	var count int64
	err := Database.Model(&model.PaymentMethod{}).
		Where("organization_id = ? AND is_active = ?", organizationID, true).
		Count(&count).Error
	return count, err
}
