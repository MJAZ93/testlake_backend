package payment

import (
	"testlake/inout"
	"testlake/model"
	"time"

	"github.com/google/uuid"
)

type PaymentMethod struct {
	ID                uuid.UUID               `json:"id"`
	OrganizationID    uuid.UUID               `json:"organization_id"`
	PayPalPayerID     *string                 `json:"paypal_payer_id"`
	PayPalEmail       *string                 `json:"paypal_email"`
	PaymentMethodType model.PaymentMethodType `json:"payment_method_type"`
	IsDefault         bool                    `json:"is_default"`
	IsActive          bool                    `json:"is_active"`
	CreatedBy         uuid.UUID               `json:"created_by"`
	CreatedAt         time.Time               `json:"created_at"`
	UpdatedAt         time.Time               `json:"updated_at"`
}

type PaymentMethodOut struct {
	inout.BaseResponse
	Data PaymentMethod `json:"data"`
}

type PaymentMethodListOut struct {
	inout.BaseResponse
	Data []PaymentMethod `json:"data"`
}

type Payment struct {
	ID              uuid.UUID               `json:"id"`
	OrganizationID  uuid.UUID               `json:"organization_id"`
	InvoiceID       *uuid.UUID              `json:"invoice_id"`
	SubscriptionID  *uuid.UUID              `json:"subscription_id"`
	PayPalPaymentID *string                 `json:"paypal_payment_id"`
	PayPalPayerID   *string                 `json:"paypal_payer_id"`
	Amount          float64                 `json:"amount"`
	Currency        string                  `json:"currency"`
	PaymentMethod   model.PaymentMethodEnum `json:"payment_method"`
	Status          model.PaymentStatus     `json:"status"`
	FailureReason   *string                 `json:"failure_reason"`
	ProcessedAt     *time.Time              `json:"processed_at"`
	CreatedAt       time.Time               `json:"created_at"`
	UpdatedAt       time.Time               `json:"updated_at"`
}

type PaymentOut struct {
	inout.BaseResponse
	Data Payment `json:"data"`
}

type PaymentListOut struct {
	inout.BaseResponse
	List []Payment            `json:"list"`
	Meta inout.PaginationMeta `json:"meta"`
}

func FromPaymentMethodModel(pm *model.PaymentMethod) PaymentMethod {
	return PaymentMethod{
		ID:                pm.ID,
		OrganizationID:    pm.OrganizationID,
		PayPalPayerID:     pm.PayPalPayerID,
		PayPalEmail:       pm.PayPalEmail,
		PaymentMethodType: pm.PaymentMethodType,
		IsDefault:         pm.IsDefault,
		IsActive:          pm.IsActive,
		CreatedBy:         pm.CreatedBy,
		CreatedAt:         pm.CreatedAt,
		UpdatedAt:         pm.UpdatedAt,
	}
}

func FromPaymentMethodModelList(pms []model.PaymentMethod) []PaymentMethod {
	result := make([]PaymentMethod, len(pms))
	for i, pm := range pms {
		result[i] = FromPaymentMethodModel(&pm)
	}
	return result
}

func FromPaymentModel(payment *model.Payment) Payment {
	return Payment{
		ID:              payment.ID,
		OrganizationID:  payment.OrganizationID,
		InvoiceID:       payment.InvoiceID,
		SubscriptionID:  payment.SubscriptionID,
		PayPalPaymentID: payment.PayPalPaymentID,
		PayPalPayerID:   payment.PayPalPayerID,
		Amount:          payment.Amount,
		Currency:        payment.Currency,
		PaymentMethod:   payment.PaymentMethod,
		Status:          payment.Status,
		FailureReason:   payment.FailureReason,
		ProcessedAt:     payment.ProcessedAt,
		CreatedAt:       payment.CreatedAt,
		UpdatedAt:       payment.UpdatedAt,
	}
}

func FromPaymentModelList(payments []model.Payment) []Payment {
	result := make([]Payment, len(payments))
	for i, payment := range payments {
		result[i] = FromPaymentModel(&payment)
	}
	return result
}
