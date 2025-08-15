package service

import (
	"testlake/controller"

	"github.com/gin-gonic/gin"
)

type PaymentMethodService struct {
	Route      string
	Controller controller.PaymentMethodController
}

// GetPaymentMethods godoc
// @Summary Get payment methods
// @Description Get all payment methods for an organization
// @Tags Payment Methods
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token" format(Bearer {token})
// @Param id path string true "Organization ID"
// @Success 200 {object} payment.PaymentMethodListOut
// @Failure 400 {object} inout.BaseResponse
// @Failure 401 {object} inout.BaseResponse
// @Failure 403 {object} inout.BaseResponse
// @Router /api/v1/organizations/{id}/payment-methods [GET]
func (s PaymentMethodService) GetPaymentMethods(r *gin.RouterGroup, route string) {
	routePath := "/" + s.Route
	if route != "" {
		routePath += "/" + route
	}
	r.GET(routePath, s.Controller.GetPaymentMethods)
}

// CreatePaymentMethod godoc
// @Summary Create payment method
// @Description Create a new payment method for an organization
// @Tags Payment Methods
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token" format(Bearer {token})
// @Param id path string true "Organization ID"
// @Param paymentMethod body payment.CreatePaymentMethodRequest true "Payment method data"
// @Success 201 {object} payment.PaymentMethodOut
// @Failure 400 {object} inout.BaseResponse
// @Failure 401 {object} inout.BaseResponse
// @Failure 403 {object} inout.BaseResponse
// @Router /api/v1/organizations/{id}/payment-methods [POST]
func (s PaymentMethodService) CreatePaymentMethod(r *gin.RouterGroup, route string) {
	routePath := "/" + s.Route
	if route != "" {
		routePath += "/" + route
	}
	r.POST(routePath, s.Controller.CreatePaymentMethod)
}

// UpdatePaymentMethod godoc
// @Summary Update payment method
// @Description Update an existing payment method
// @Tags Payment Methods
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token" format(Bearer {token})
// @Param id path string true "Organization ID"
// @Param pmId path string true "Payment Method ID"
// @Param paymentMethod body payment.UpdatePaymentMethodRequest true "Payment method data"
// @Success 200 {object} payment.PaymentMethodOut
// @Failure 400 {object} inout.BaseResponse
// @Failure 401 {object} inout.BaseResponse
// @Failure 403 {object} inout.BaseResponse
// @Failure 404 {object} inout.BaseResponse
// @Router /api/v1/organizations/{id}/payment-methods/{pmId} [PUT]
func (s PaymentMethodService) UpdatePaymentMethod(r *gin.RouterGroup, route string) {
	routePath := "/" + s.Route
	if route != "" {
		routePath += "/" + route
	}
	routePath += "/:pmId"
	r.PUT(routePath, s.Controller.UpdatePaymentMethod)
}

// DeletePaymentMethod godoc
// @Summary Delete payment method
// @Description Delete a payment method
// @Tags Payment Methods
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token" format(Bearer {token})
// @Param id path string true "Organization ID"
// @Param pmId path string true "Payment Method ID"
// @Success 200 {object} inout.BaseResponse
// @Failure 400 {object} inout.BaseResponse
// @Failure 401 {object} inout.BaseResponse
// @Failure 403 {object} inout.BaseResponse
// @Failure 404 {object} inout.BaseResponse
// @Router /api/v1/organizations/{id}/payment-methods/{pmId} [DELETE]
func (s PaymentMethodService) DeletePaymentMethod(r *gin.RouterGroup, route string) {
	routePath := "/" + s.Route
	if route != "" {
		routePath += "/" + route
	}
	routePath += "/:pmId"
	r.DELETE(routePath, s.Controller.DeletePaymentMethod)
}

// SetDefaultPaymentMethod godoc
// @Summary Set default payment method
// @Description Set a payment method as the default for the organization
// @Tags Payment Methods
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token" format(Bearer {token})
// @Param id path string true "Organization ID"
// @Param pmId path string true "Payment Method ID"
// @Success 200 {object} inout.BaseResponse
// @Failure 400 {object} inout.BaseResponse
// @Failure 401 {object} inout.BaseResponse
// @Failure 403 {object} inout.BaseResponse
// @Failure 404 {object} inout.BaseResponse
// @Router /api/v1/organizations/{id}/payment-methods/{pmId}/set-default [PUT]
func (s PaymentMethodService) SetDefaultPaymentMethod(r *gin.RouterGroup, route string) {
	routePath := "/" + s.Route
	if route != "" {
		routePath += "/" + route
	}
	routePath += "/:pmId/set-default"
	r.PUT(routePath, s.Controller.SetDefaultPaymentMethod)
}
