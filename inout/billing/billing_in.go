package billing

type PayInvoiceRequest struct {
	PaymentMethodID string `json:"payment_method_id" binding:"required"`
}
