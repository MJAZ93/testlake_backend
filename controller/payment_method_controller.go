package controller

import (
	"errors"
	"net/http"
	"testlake/dao"
	"testlake/inout"
	"testlake/inout/payment"
	"testlake/model"
	"testlake/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentMethodController struct{}

// GetPaymentMethods returns all payment methods for an organization
func (controller PaymentMethodController) GetPaymentMethods(context *gin.Context) {
	orgIDParam := context.Param("id")
	organizationID, err := uuid.Parse(orgIDParam)
	if err != nil {
		utils.ReportBadRequest(context, "Invalid organization ID")
		return
	}

	// Verify user has access to organization
	userID, err := utils.ExtractUserID(context)
	if err != nil {
		utils.ReportUnauthorized(context, "Authentication required")
		return
	}

	if !controller.verifyOrganizationAccess(userID, organizationID) {
		utils.ReportForbidden(context, "Access denied")
		return
	}

	paymentMethodDao := dao.NewPaymentMethodDao()
	paymentMethods, err := paymentMethodDao.GetByOrganizationID(organizationID)
	if err != nil {
		utils.ReportInternalServerError(context, "Database error")
		return
	}

	response := payment.PaymentMethodListOut{
		BaseResponse: inout.BaseResponse{
			ErrorCode:        0,
			ErrorDescription: "Success",
		},
		Data: payment.FromPaymentMethodModelList(paymentMethods),
	}

	context.JSON(http.StatusOK, response)
}

// CreatePaymentMethod creates a new payment method
func (controller PaymentMethodController) CreatePaymentMethod(context *gin.Context) {
	orgIDParam := context.Param("id")
	organizationID, err := uuid.Parse(orgIDParam)
	if err != nil {
		utils.ReportBadRequest(context, "Invalid organization ID")
		return
	}

	userID, err := utils.ExtractUserID(context)
	if err != nil {
		utils.ReportUnauthorized(context, "Authentication required")
		return
	}

	if !controller.verifyOrganizationAccess(userID, organizationID) {
		utils.ReportForbidden(context, "Access denied")
		return
	}

	var request payment.CreatePaymentMethodRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		utils.ReportBadRequest(context, "Invalid request body")
		return
	}

	paymentMethod := &model.PaymentMethod{
		OrganizationID:    organizationID,
		PayPalEmail:       &request.PayPalEmail,
		PayPalPayerID:     request.PayPalPayerID,
		PaymentMethodType: request.PaymentMethodType,
		IsDefault:         request.IsDefault,
		IsActive:          true,
		CreatedBy:         userID,
	}

	paymentMethodDao := dao.NewPaymentMethodDao()

	// If this is set as default, make sure no other payment method is default
	if request.IsDefault {
		err = paymentMethodDao.SetDefault(organizationID, paymentMethod.ID)
		if err != nil {
			utils.ReportInternalServerError(context, "Failed to update default payment method")
			return
		}
	}

	err = paymentMethodDao.Create(paymentMethod)
	if err != nil {
		utils.ReportInternalServerError(context, "Failed to create payment method")
		return
	}

	response := payment.PaymentMethodOut{
		BaseResponse: inout.BaseResponse{
			ErrorCode:        0,
			ErrorDescription: "Success",
		},
		Data: payment.FromPaymentMethodModel(paymentMethod),
	}

	context.JSON(http.StatusCreated, response)
}

// UpdatePaymentMethod updates an existing payment method
func (controller PaymentMethodController) UpdatePaymentMethod(context *gin.Context) {
	orgIDParam := context.Param("id")
	organizationID, err := uuid.Parse(orgIDParam)
	if err != nil {
		utils.ReportBadRequest(context, "Invalid organization ID")
		return
	}

	idParam := context.Param("pmId")
	paymentMethodID, err := uuid.Parse(idParam)
	if err != nil {
		utils.ReportBadRequest(context, "Invalid payment method ID")
		return
	}

	userID, err := utils.ExtractUserID(context)
	if err != nil {
		utils.ReportUnauthorized(context, "Authentication required")
		return
	}

	if !controller.verifyOrganizationAccess(userID, organizationID) {
		utils.ReportForbidden(context, "Access denied")
		return
	}

	var request payment.UpdatePaymentMethodRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		utils.ReportBadRequest(context, "Invalid request body")
		return
	}

	paymentMethodDao := dao.NewPaymentMethodDao()
	paymentMethod, err := paymentMethodDao.GetByID(paymentMethodID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ReportNotFound(context, "Payment method not found")
		} else {
			utils.ReportInternalServerError(context, "Database error")
		}
		return
	}

	// Verify the payment method belongs to the organization
	if paymentMethod.OrganizationID != organizationID {
		utils.ReportForbidden(context, "Access denied")
		return
	}

	// Update fields
	if request.PayPalEmail != nil {
		paymentMethod.PayPalEmail = request.PayPalEmail
	}
	if request.PayPalPayerID != nil {
		paymentMethod.PayPalPayerID = request.PayPalPayerID
	}
	if request.IsDefault != nil && *request.IsDefault {
		err = paymentMethodDao.SetDefault(organizationID, paymentMethodID)
		if err != nil {
			utils.ReportInternalServerError(context, "Failed to update default payment method")
			return
		}
	}

	err = paymentMethodDao.Update(paymentMethod)
	if err != nil {
		utils.ReportInternalServerError(context, "Failed to update payment method")
		return
	}

	response := payment.PaymentMethodOut{
		BaseResponse: inout.BaseResponse{
			ErrorCode:        0,
			ErrorDescription: "Success",
		},
		Data: payment.FromPaymentMethodModel(paymentMethod),
	}

	context.JSON(http.StatusOK, response)
}

// DeletePaymentMethod deletes a payment method
func (controller PaymentMethodController) DeletePaymentMethod(context *gin.Context) {
	orgIDParam := context.Param("id")
	organizationID, err := uuid.Parse(orgIDParam)
	if err != nil {
		utils.ReportBadRequest(context, "Invalid organization ID")
		return
	}

	idParam := context.Param("pmId")
	paymentMethodID, err := uuid.Parse(idParam)
	if err != nil {
		utils.ReportBadRequest(context, "Invalid payment method ID")
		return
	}

	userID, err := utils.ExtractUserID(context)
	if err != nil {
		utils.ReportUnauthorized(context, "Authentication required")
		return
	}

	if !controller.verifyOrganizationAccess(userID, organizationID) {
		utils.ReportForbidden(context, "Access denied")
		return
	}

	paymentMethodDao := dao.NewPaymentMethodDao()
	paymentMethod, err := paymentMethodDao.GetByID(paymentMethodID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ReportNotFound(context, "Payment method not found")
		} else {
			utils.ReportInternalServerError(context, "Database error")
		}
		return
	}

	// Verify the payment method belongs to the organization
	if paymentMethod.OrganizationID != organizationID {
		utils.ReportForbidden(context, "Access denied")
		return
	}

	err = paymentMethodDao.SetActive(paymentMethodID, false)
	if err != nil {
		utils.ReportInternalServerError(context, "Failed to delete payment method")
		return
	}

	response := inout.BaseResponse{
		ErrorCode:        0,
		ErrorDescription: "Payment method deleted successfully",
	}

	context.JSON(http.StatusOK, response)
}

// SetDefaultPaymentMethod sets a payment method as default
func (controller PaymentMethodController) SetDefaultPaymentMethod(context *gin.Context) {
	orgIDParam := context.Param("id")
	organizationID, err := uuid.Parse(orgIDParam)
	if err != nil {
		utils.ReportBadRequest(context, "Invalid organization ID")
		return
	}

	idParam := context.Param("pmId")
	paymentMethodID, err := uuid.Parse(idParam)
	if err != nil {
		utils.ReportBadRequest(context, "Invalid payment method ID")
		return
	}

	userID, err := utils.ExtractUserID(context)
	if err != nil {
		utils.ReportUnauthorized(context, "Authentication required")
		return
	}

	if !controller.verifyOrganizationAccess(userID, organizationID) {
		utils.ReportForbidden(context, "Access denied")
		return
	}

	paymentMethodDao := dao.NewPaymentMethodDao()
	paymentMethod, err := paymentMethodDao.GetByID(paymentMethodID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ReportNotFound(context, "Payment method not found")
		} else {
			utils.ReportInternalServerError(context, "Database error")
		}
		return
	}

	// Verify the payment method belongs to the organization
	if paymentMethod.OrganizationID != organizationID {
		utils.ReportForbidden(context, "Access denied")
		return
	}

	err = paymentMethodDao.SetDefault(organizationID, paymentMethodID)
	if err != nil {
		utils.ReportInternalServerError(context, "Failed to set default payment method")
		return
	}

	response := inout.BaseResponse{
		ErrorCode:        0,
		ErrorDescription: "Payment method set as default successfully",
	}

	context.JSON(http.StatusOK, response)
}

// Helper method to verify organization access
func (controller PaymentMethodController) verifyOrganizationAccess(userID, organizationID uuid.UUID) bool {
	// This should check if the user has access to the organization
	// For now, implementing a simple check - in production, this would verify membership
	orgDao := dao.NewOrganizationDao()
	_, err := orgDao.GetByID(organizationID)
	return err == nil
}
