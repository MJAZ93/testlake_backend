package service

import (
	"testlake/controller"

	"github.com/gin-gonic/gin"
)

type BillingService struct {
	Route      string
	Controller controller.BillingController
}

// GetBillingOverview godoc
// @Summary Get billing overview
// @Description Get billing overview for an organization
// @Tags Billing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token" format(Bearer {token})
// @Param id path string true "Organization ID"
// @Success 200 {object} billing.BillingOverviewOut
// @Failure 400 {object} inout.BaseResponse
// @Failure 401 {object} inout.BaseResponse
// @Failure 403 {object} inout.BaseResponse
// @Router /api/v1/organizations/{id}/billing/overview [GET]
func (s BillingService) GetBillingOverview(r *gin.RouterGroup, route string) {
	r.GET("/"+s.Route+"/"+route, s.Controller.GetBillingOverview)
}

// GetInvoices godoc
// @Summary Get invoices
// @Description Get invoices for an organization
// @Tags Billing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token" format(Bearer {token})
// @Param id path string true "Organization ID"
// @Param page query int false "Page number"
// @Success 200 {object} billing.InvoiceListOut
// @Failure 400 {object} inout.BaseResponse
// @Failure 401 {object} inout.BaseResponse
// @Failure 403 {object} inout.BaseResponse
// @Router /api/v1/organizations/{id}/invoices [GET]
func (s BillingService) GetInvoices(r *gin.RouterGroup, route string) {
	r.GET("/"+s.Route+"/"+route, s.Controller.GetInvoices)
}

// GetInvoice godoc
// @Summary Get invoice
// @Description Get a specific invoice by ID
// @Tags Billing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token" format(Bearer {token})
// @Param id path string true "Invoice ID"
// @Success 200 {object} billing.InvoiceOut
// @Failure 400 {object} inout.BaseResponse
// @Failure 401 {object} inout.BaseResponse
// @Failure 403 {object} inout.BaseResponse
// @Failure 404 {object} inout.BaseResponse
// @Router /api/v1/invoices/{id} [GET]
func (s BillingService) GetInvoice(r *gin.RouterGroup, route string) {
	r.GET("/"+s.Route+"/"+route+"/:id", s.Controller.GetInvoice)
}

// DownloadInvoice godoc
// @Summary Download invoice
// @Description Download a specific invoice
// @Tags Billing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token" format(Bearer {token})
// @Param id path string true "Invoice ID"
// @Success 302 "Redirect to invoice PDF"
// @Failure 400 {object} inout.BaseResponse
// @Failure 401 {object} inout.BaseResponse
// @Failure 403 {object} inout.BaseResponse
// @Failure 404 {object} inout.BaseResponse
// @Router /api/v1/invoices/{id}/download [GET]
func (s BillingService) DownloadInvoice(r *gin.RouterGroup, route string) {
	r.GET("/"+s.Route+"/"+route+"/:id/download", s.Controller.DownloadInvoice)
}

// PayInvoice godoc
// @Summary Pay invoice
// @Description Process payment for an invoice
// @Tags Billing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token" format(Bearer {token})
// @Param id path string true "Invoice ID"
// @Param payment body billing.PayInvoiceRequest true "Payment data"
// @Success 200 {object} billing.PaymentOut
// @Failure 400 {object} inout.BaseResponse
// @Failure 401 {object} inout.BaseResponse
// @Failure 403 {object} inout.BaseResponse
// @Failure 404 {object} inout.BaseResponse
// @Router /api/v1/invoices/{id}/pay [POST]
func (s BillingService) PayInvoice(r *gin.RouterGroup, route string) {
	r.POST("/"+s.Route+"/"+route+"/:id/pay", s.Controller.PayInvoice)
}

// GetBillingHistory godoc
// @Summary Get billing history
// @Description Get billing history for an organization
// @Tags Billing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token" format(Bearer {token})
// @Param id path string true "Organization ID"
// @Param page query int false "Page number"
// @Success 200 {object} billing.BillingHistoryOut
// @Failure 400 {object} inout.BaseResponse
// @Failure 401 {object} inout.BaseResponse
// @Failure 403 {object} inout.BaseResponse
// @Router /api/v1/organizations/{id}/billing/history [GET]
func (s BillingService) GetBillingHistory(r *gin.RouterGroup, route string) {
	r.GET("/"+s.Route+"/"+route, s.Controller.GetBillingHistory)
}
