package billing

import (
	"testlake/inout"
	"testlake/inout/payment"
	"testlake/inout/plan"
	"testlake/inout/subscription"
	"testlake/model"
	"time"

	"github.com/google/uuid"
)

type InvoiceLineItem struct {
	ID          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	Quantity    int       `json:"quantity"`
	UnitPrice   float64   `json:"unit_price"`
	TotalPrice  float64   `json:"total_price"`
}

type Invoice struct {
	ID                 uuid.UUID           `json:"id"`
	OrganizationID     uuid.UUID           `json:"organization_id"`
	SubscriptionID     *uuid.UUID          `json:"subscription_id"`
	PayPalInvoiceID    *string             `json:"paypal_invoice_id"`
	InvoiceNumber      string              `json:"invoice_number"`
	Amount             float64             `json:"amount"`
	TaxAmount          float64             `json:"tax_amount"`
	TotalAmount        float64             `json:"total_amount"`
	Currency           string              `json:"currency"`
	Status             model.InvoiceStatus `json:"status"`
	BillingPeriodStart *time.Time          `json:"billing_period_start"`
	BillingPeriodEnd   *time.Time          `json:"billing_period_end"`
	DueDate            *time.Time          `json:"due_date"`
	PaidAt             *time.Time          `json:"paid_at"`
	InvoiceURL         *string             `json:"invoice_url"`
	CreatedAt          time.Time           `json:"created_at"`
	UpdatedAt          time.Time           `json:"updated_at"`
	LineItems          []InvoiceLineItem   `json:"line_items,omitempty"`
}

type InvoiceOut struct {
	inout.BaseResponse
	Data Invoice `json:"data"`
}

type InvoiceListOut struct {
	inout.BaseResponse
	List []Invoice            `json:"list"`
	Meta inout.PaginationMeta `json:"meta"`
}

type BillingOverview struct {
	CurrentSubscription *subscription.Subscription `json:"current_subscription"`
	CurrentPlan         *plan.Plan                 `json:"current_plan"`
	NextBillingDate     *time.Time                 `json:"next_billing_date"`
	NextBillingAmount   float64                    `json:"next_billing_amount"`
	CurrentUsage        subscription.UsageMetrics  `json:"current_usage"`
	PlanLimits          subscription.PlanLimits    `json:"plan_limits"`
	UnpaidInvoices      []Invoice                  `json:"unpaid_invoices"`
	RecentPayments      []payment.Payment          `json:"recent_payments"`
}

type BillingOverviewOut struct {
	inout.BaseResponse
	Data BillingOverview `json:"data"`
}

type BillingHistoryItem struct {
	ID          uuid.UUID `json:"id"`
	Type        string    `json:"type"` // "invoice" or "payment"
	Amount      float64   `json:"amount"`
	Currency    string    `json:"currency"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
}

type BillingHistory struct {
	Items []BillingHistoryItem `json:"items"`
}

type BillingHistoryOut struct {
	inout.BaseResponse
	List []BillingHistoryItem `json:"list"`
	Meta inout.PaginationMeta `json:"meta"`
}

type PaymentOut struct {
	inout.BaseResponse
	Data payment.Payment `json:"data"`
}

func FromInvoiceLineItemModel(item *model.InvoiceLineItem) InvoiceLineItem {
	return InvoiceLineItem{
		ID:          item.ID,
		Description: item.Description,
		Quantity:    item.Quantity,
		UnitPrice:   item.UnitPrice,
		TotalPrice:  item.TotalPrice,
	}
}

func FromInvoiceLineItemModelList(items []model.InvoiceLineItem) []InvoiceLineItem {
	result := make([]InvoiceLineItem, len(items))
	for i, item := range items {
		result[i] = FromInvoiceLineItemModel(&item)
	}
	return result
}

func FromInvoiceModel(invoice *model.Invoice) Invoice {
	result := Invoice{
		ID:                 invoice.ID,
		OrganizationID:     invoice.OrganizationID,
		SubscriptionID:     invoice.SubscriptionID,
		PayPalInvoiceID:    invoice.PayPalInvoiceID,
		InvoiceNumber:      invoice.InvoiceNumber,
		Amount:             invoice.Amount,
		TaxAmount:          invoice.TaxAmount,
		TotalAmount:        invoice.TotalAmount,
		Currency:           invoice.Currency,
		Status:             invoice.Status,
		BillingPeriodStart: invoice.BillingPeriodStart,
		BillingPeriodEnd:   invoice.BillingPeriodEnd,
		DueDate:            invoice.DueDate,
		PaidAt:             invoice.PaidAt,
		InvoiceURL:         invoice.InvoiceURL,
		CreatedAt:          invoice.CreatedAt,
		UpdatedAt:          invoice.UpdatedAt,
	}

	if invoice.LineItems != nil {
		result.LineItems = FromInvoiceLineItemModelList(invoice.LineItems)
	}

	return result
}

func FromInvoiceModelList(invoices []model.Invoice) []Invoice {
	result := make([]Invoice, len(invoices))
	for i, invoice := range invoices {
		result[i] = FromInvoiceModel(&invoice)
	}
	return result
}
