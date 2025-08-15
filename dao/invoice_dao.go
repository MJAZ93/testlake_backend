package dao

import (
	"fmt"
	"testlake/model"
	"time"

	"github.com/google/uuid"
)

type InvoiceDao struct {
	Limit int
}

func NewInvoiceDao() *InvoiceDao {
	return &InvoiceDao{Limit: 50}
}

func (dao *InvoiceDao) Create(invoice *model.Invoice) error {
	return Database.Create(invoice).Error
}

func (dao *InvoiceDao) CreateWithLineItems(invoice *model.Invoice, lineItems []model.InvoiceLineItem) error {
	tx := Database.Begin()

	if err := tx.Create(invoice).Error; err != nil {
		tx.Rollback()
		return err
	}

	for i := range lineItems {
		lineItems[i].InvoiceID = invoice.ID
		if err := tx.Create(&lineItems[i]).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (dao *InvoiceDao) GetByID(id uuid.UUID) (*model.Invoice, error) {
	var invoice model.Invoice
	err := Database.Preload("LineItems").
		Preload("Organization").
		Preload("Subscription").
		First(&invoice, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &invoice, nil
}

func (dao *InvoiceDao) GetByInvoiceNumber(invoiceNumber string) (*model.Invoice, error) {
	var invoice model.Invoice
	err := Database.Preload("LineItems").
		Preload("Organization").
		Preload("Subscription").
		First(&invoice, "invoice_number = ?", invoiceNumber).Error
	if err != nil {
		return nil, err
	}
	return &invoice, nil
}

func (dao *InvoiceDao) GetByOrganizationID(organizationID uuid.UUID, page int) ([]model.Invoice, int64, error) {
	var invoices []model.Invoice
	var total int64

	err := Database.Model(&model.Invoice{}).
		Where("organization_id = ?", organizationID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := page * dao.Limit
	err = Database.Preload("LineItems").
		Where("organization_id = ?", organizationID).
		Order("created_at DESC").
		Offset(offset).Limit(dao.Limit).
		Find(&invoices).Error
	if err != nil {
		return nil, 0, err
	}

	return invoices, total, nil
}

func (dao *InvoiceDao) GetUnpaidByOrganizationID(organizationID uuid.UUID) ([]model.Invoice, error) {
	var invoices []model.Invoice
	err := Database.Preload("LineItems").
		Where("organization_id = ? AND status IN ?", organizationID, []model.InvoiceStatus{
			model.InvoiceStatusDraft,
			model.InvoiceStatusSent,
		}).
		Order("created_at DESC").
		Find(&invoices).Error
	if err != nil {
		return nil, err
	}
	return invoices, nil
}

func (dao *InvoiceDao) Update(invoice *model.Invoice) error {
	return Database.Save(invoice).Error
}

func (dao *InvoiceDao) UpdateStatus(id uuid.UUID, status model.InvoiceStatus) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if status == model.InvoiceStatusPaid {
		updates["paid_at"] = "NOW()"
	}
	return Database.Model(&model.Invoice{}).Where("id = ?", id).Updates(updates).Error
}

func (dao *InvoiceDao) UpdatePayPalInvoiceID(id uuid.UUID, paypalInvoiceID string) error {
	return Database.Model(&model.Invoice{}).
		Where("id = ?", id).
		Update("paypal_invoice_id", paypalInvoiceID).Error
}

func (dao *InvoiceDao) Delete(id uuid.UUID) error {
	tx := Database.Begin()

	// Delete line items first
	if err := tx.Delete(&model.InvoiceLineItem{}, "invoice_id = ?", id).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete invoice
	if err := tx.Delete(&model.Invoice{}, "id = ?", id).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (dao *InvoiceDao) GetAll(page int) ([]model.Invoice, int64, error) {
	var invoices []model.Invoice
	var total int64

	err := Database.Model(&model.Invoice{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := page * dao.Limit
	err = Database.Preload("LineItems").
		Preload("Organization").
		Preload("Subscription").
		Order("created_at DESC").
		Offset(offset).Limit(dao.Limit).
		Find(&invoices).Error
	if err != nil {
		return nil, 0, err
	}

	return invoices, total, nil
}

func (dao *InvoiceDao) GenerateInvoiceNumber() (string, error) {
	var count int64
	err := Database.Model(&model.Invoice{}).Count(&count).Error
	if err != nil {
		return "", err
	}

	// Generate invoice number in format INV-YYYY-NNNNNN
	return fmt.Sprintf("INV-%d-%06d", time.Now().Year(), count+1), nil
}
