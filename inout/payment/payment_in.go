package payment

import (
	"testlake/model"
)

type CreatePaymentMethodRequest struct {
	PayPalEmail       string                  `json:"paypal_email" binding:"required"`
	PayPalPayerID     *string                 `json:"paypal_payer_id"`
	PaymentMethodType model.PaymentMethodType `json:"payment_method_type"`
	IsDefault         bool                    `json:"is_default"`
}

type UpdatePaymentMethodRequest struct {
	PayPalEmail   *string `json:"paypal_email"`
	PayPalPayerID *string `json:"paypal_payer_id"`
	IsDefault     *bool   `json:"is_default"`
}
