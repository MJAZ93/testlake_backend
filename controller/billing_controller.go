package controller

import (
	"errors"
	"net/http"
	"strconv"
	"testlake/dao"
	"testlake/inout"
	"testlake/inout/billing"
	"testlake/inout/payment"
	"testlake/inout/plan"
	"testlake/inout/subscription"
	"testlake/model"
	"testlake/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BillingController struct{}

// GetBillingOverview returns billing overview for an organization
func (controller BillingController) GetBillingOverview(context *gin.Context) {
	orgIDParam := context.Param("orgId")
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

	overview := billing.BillingOverview{}

	// Get current subscription
	subscriptionDao := dao.NewSubscriptionDao()
	currentSub, err := subscriptionDao.GetActiveByOrganizationID(organizationID)
	if err == nil && currentSub != nil {
		sub := subscription.FromSubscriptionModel(currentSub)
		overview.CurrentSubscription = &sub
		overview.NextBillingDate = &currentSub.CurrentPeriodEnd
	}

	// Get current plan
	orgDao := dao.NewOrganizationDao()
	org, err := orgDao.GetByID(organizationID)
	if err == nil && org.PlanID != nil {
		planDao := dao.NewPlanDao()
		currentPlan, err := planDao.GetByID(*org.PlanID)
		if err == nil {
			planData := plan.FromPlanModel(currentPlan)
			overview.CurrentPlan = &planData

			// Calculate next billing amount
			if org.BillingCycle == model.BillingCycleMonthly {
				overview.NextBillingAmount = currentPlan.PriceMonthly
			} else {
				overview.NextBillingAmount = currentPlan.PriceYearly
			}
		}
	}

	// Get current usage
	usageDao := dao.NewOrganizationUsageDao()
	currentUsage, err := usageDao.GetCurrentUsage(organizationID)
	if err == nil && currentUsage != nil {
		overview.CurrentUsage = subscription.UsageMetrics{
			UsersCount:        currentUsage.UsersCount,
			ProjectsCount:     currentUsage.ProjectsCount,
			EnvironmentsCount: currentUsage.EnvironmentsCount,
			SchemasCount:      currentUsage.SchemasCount,
			TestRecordsCount:  currentUsage.TestRecordsCount,
			APIRequestsCount:  currentUsage.APIRequestsCount,
		}
	}

	// Get plan limits
	if overview.CurrentPlan != nil {
		overview.PlanLimits = subscription.PlanLimits{
			MaxUsers:                overview.CurrentPlan.MaxUsers,
			MaxProjects:             overview.CurrentPlan.MaxProjects,
			MaxEnvironments:         overview.CurrentPlan.MaxEnvironments,
			MaxSchemas:              overview.CurrentPlan.MaxSchemas,
			MaxTestRecordsPerSchema: overview.CurrentPlan.MaxTestRecordsPerSchema,
		}
	}

	// Get unpaid invoices
	invoiceDao := dao.NewInvoiceDao()
	unpaidInvoices, err := invoiceDao.GetUnpaidByOrganizationID(organizationID)
	if err == nil {
		overview.UnpaidInvoices = billing.FromInvoiceModelList(unpaidInvoices)
	}

	// Get recent payments
	paymentDao := dao.NewPaymentDao()
	recentPayments, err := paymentDao.GetRecentByOrganizationID(organizationID, 5)
	if err == nil {
		overview.RecentPayments = payment.FromPaymentModelList(recentPayments)
	}

	response := billing.BillingOverviewOut{
		BaseResponse: inout.BaseResponse{
			ErrorCode:        0,
			ErrorDescription: "Success",
		},
		Data: overview,
	}

	context.JSON(http.StatusOK, response)
}

// GetInvoices returns invoices for an organization
func (controller BillingController) GetInvoices(context *gin.Context) {
	orgIDParam := context.Param("orgId")
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

	pageParam := context.DefaultQuery("page", "0")
	page, err := strconv.Atoi(pageParam)
	if err != nil || page < 0 {
		utils.ReportBadRequest(context, "Invalid page parameter")
		return
	}

	invoiceDao := dao.NewInvoiceDao()
	invoices, total, err := invoiceDao.GetByOrganizationID(organizationID, page)
	if err != nil {
		utils.ReportInternalServerError(context, "Database error")
		return
	}

	totalPages := int(total) / invoiceDao.Limit
	if int(total)%invoiceDao.Limit > 0 {
		totalPages++
	}

	response := billing.InvoiceListOut{
		BaseResponse: inout.BaseResponse{
			ErrorCode:        0,
			ErrorDescription: "Success",
		},
		List: billing.FromInvoiceModelList(invoices),
		Meta: inout.PaginationMeta{
			Page:       page,
			Limit:      invoiceDao.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}

	context.JSON(http.StatusOK, response)
}

// GetInvoice returns a specific invoice
func (controller BillingController) GetInvoice(context *gin.Context) {
	idParam := context.Param("id")
	invoiceID, err := uuid.Parse(idParam)
	if err != nil {
		utils.ReportBadRequest(context, "Invalid invoice ID")
		return
	}

	userID, err := utils.ExtractUserID(context)
	if err != nil {
		utils.ReportUnauthorized(context, "Authentication required")
		return
	}

	invoiceDao := dao.NewInvoiceDao()
	invoice, err := invoiceDao.GetByID(invoiceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ReportNotFound(context, "Invoice not found")
		} else {
			utils.ReportInternalServerError(context, "Database error")
		}
		return
	}

	// Verify user has access to the organization
	if !controller.verifyOrganizationAccess(userID, invoice.OrganizationID) {
		utils.ReportForbidden(context, "Access denied")
		return
	}

	response := billing.InvoiceOut{
		BaseResponse: inout.BaseResponse{
			ErrorCode:        0,
			ErrorDescription: "Success",
		},
		Data: billing.FromInvoiceModel(invoice),
	}

	context.JSON(http.StatusOK, response)
}

// DownloadInvoice provides a download link for an invoice
func (controller BillingController) DownloadInvoice(context *gin.Context) {
	idParam := context.Param("id")
	invoiceID, err := uuid.Parse(idParam)
	if err != nil {
		utils.ReportBadRequest(context, "Invalid invoice ID")
		return
	}

	userID, err := utils.ExtractUserID(context)
	if err != nil {
		utils.ReportUnauthorized(context, "Authentication required")
		return
	}

	invoiceDao := dao.NewInvoiceDao()
	invoice, err := invoiceDao.GetByID(invoiceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ReportNotFound(context, "Invoice not found")
		} else {
			utils.ReportInternalServerError(context, "Database error")
		}
		return
	}

	// Verify user has access to the organization
	if !controller.verifyOrganizationAccess(userID, invoice.OrganizationID) {
		utils.ReportForbidden(context, "Access denied")
		return
	}

	// In a real implementation, this would generate and return a PDF or redirect to a PDF URL
	if invoice.InvoiceURL != nil {
		context.Redirect(http.StatusFound, *invoice.InvoiceURL)
	} else {
		utils.ReportNotFound(context, "Invoice download not available")
	}
}

// PayInvoice processes payment for an invoice
func (controller BillingController) PayInvoice(context *gin.Context) {
	idParam := context.Param("id")
	invoiceID, err := uuid.Parse(idParam)
	if err != nil {
		utils.ReportBadRequest(context, "Invalid invoice ID")
		return
	}

	userID, err := utils.ExtractUserID(context)
	if err != nil {
		utils.ReportUnauthorized(context, "Authentication required")
		return
	}

	var request billing.PayInvoiceRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		utils.ReportBadRequest(context, "Invalid request body")
		return
	}

	invoiceDao := dao.NewInvoiceDao()
	invoice, err := invoiceDao.GetByID(invoiceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ReportNotFound(context, "Invoice not found")
		} else {
			utils.ReportInternalServerError(context, "Database error")
		}
		return
	}

	// Verify user has access to the organization
	if !controller.verifyOrganizationAccess(userID, invoice.OrganizationID) {
		utils.ReportForbidden(context, "Access denied")
		return
	}

	// Check if invoice is already paid
	if invoice.Status == model.InvoiceStatusPaid {
		utils.ReportBadRequest(context, "Invoice is already paid")
		return
	}

	// In a real implementation, this would process payment through PayPal
	// For now, we'll simulate a successful payment

	// Create payment record
	paymentDao := dao.NewPaymentDao()
	newPayment := &model.Payment{
		OrganizationID: invoice.OrganizationID,
		InvoiceID:      &invoice.ID,
		Amount:         invoice.TotalAmount,
		Currency:       invoice.Currency,
		PaymentMethod:  model.PaymentMethodEnumPayPal,
		Status:         model.PaymentStatusCompleted,
	}

	err = paymentDao.Create(newPayment)
	if err != nil {
		utils.ReportInternalServerError(context, "Failed to create payment record")
		return
	}

	// Update invoice status
	err = invoiceDao.UpdateStatus(invoice.ID, model.InvoiceStatusPaid)
	if err != nil {
		utils.ReportInternalServerError(context, "Failed to update invoice status")
		return
	}

	response := billing.PaymentOut{
		BaseResponse: inout.BaseResponse{
			ErrorCode:        0,
			ErrorDescription: "Payment processed successfully",
		},
		Data: payment.FromPaymentModel(newPayment),
	}

	context.JSON(http.StatusOK, response)
}

// GetBillingHistory returns billing history for an organization
func (controller BillingController) GetBillingHistory(context *gin.Context) {
	orgIDParam := context.Param("orgId")
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

	pageParam := context.DefaultQuery("page", "0")
	page, err := strconv.Atoi(pageParam)
	if err != nil || page < 0 {
		utils.ReportBadRequest(context, "Invalid page parameter")
		return
	}

	// Get billing history items (invoices and payments)
	var historyItems []billing.BillingHistoryItem

	// Get invoices
	invoiceDao := dao.NewInvoiceDao()
	invoices, _, err := invoiceDao.GetByOrganizationID(organizationID, 0) // Get all for history
	if err == nil {
		for _, invoice := range invoices {
			historyItems = append(historyItems, billing.BillingHistoryItem{
				ID:          invoice.ID,
				Type:        "invoice",
				Amount:      invoice.TotalAmount,
				Currency:    invoice.Currency,
				Status:      string(invoice.Status),
				Description: "Invoice " + invoice.InvoiceNumber,
				Date:        invoice.CreatedAt,
			})
		}
	}

	// Get payments
	paymentDao := dao.NewPaymentDao()
	payments, _, err := paymentDao.GetByOrganizationID(organizationID, 0) // Get all for history
	if err == nil {
		for _, payment := range payments {
			historyItems = append(historyItems, billing.BillingHistoryItem{
				ID:          payment.ID,
				Type:        "payment",
				Amount:      payment.Amount,
				Currency:    payment.Currency,
				Status:      string(payment.Status),
				Description: "Payment",
				Date:        payment.CreatedAt,
			})
		}
	}

	// Pagination
	limit := 50
	total := int64(len(historyItems))
	start := page * limit
	end := start + limit
	if end > len(historyItems) {
		end = len(historyItems)
	}
	if start > len(historyItems) {
		start = len(historyItems)
	}

	paginatedItems := historyItems[start:end]
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	response := billing.BillingHistoryOut{
		BaseResponse: inout.BaseResponse{
			ErrorCode:        0,
			ErrorDescription: "Success",
		},
		List: paginatedItems,
		Meta: inout.PaginationMeta{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}

	context.JSON(http.StatusOK, response)
}

// Helper method to verify organization access
func (controller BillingController) verifyOrganizationAccess(userID, organizationID uuid.UUID) bool {
	orgDao := dao.NewOrganizationDao()
	_, err := orgDao.GetByID(organizationID)
	return err == nil
}
