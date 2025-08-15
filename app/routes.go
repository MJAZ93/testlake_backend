package app

import (
	"testlake/controller"
	"testlake/service"

	"github.com/gin-gonic/gin"
)

func PublicRoutes(r *gin.RouterGroup) {
	// Create public sub-group
	// Authentication endpoints (public)
	authService := service.AuthService{
		Route:      "auth",
		Controller: controller.AuthController{},
	}

	authService.SignUp(r, "signup")
	authService.SignIn(r, "signin")
	authService.SignOut(r, "signout")
	authService.ForgotPassword(r, "forgot-password")
	authService.ResetPassword(r, "reset-password")
	authService.VerifyEmail(r, "verify-email")
	authService.ResendEmailConfirmation(r, "resend-email-confirmation")

	// Plan endpoints (public)
	planService := service.PlanService{
		Route:      "plans",
		Controller: controller.PlanController{},
	}

	planService.GetAllPlans(r, "")
	planService.GetPlan(r, "")
	planService.ComparePlans(r, "compare")
}

func PrivateRoutes(r *gin.RouterGroup) {
	// Create private sub-group

	// Authentication endpoints (require JWT)
	authService := service.AuthService{
		Route:      "auth",
		Controller: controller.AuthController{},
	}

	authService.RefreshToken(r, "refresh")

	// User Management endpoints
	userService := service.UserService{
		Route:      "users",
		Controller: controller.UserController{},
	}

	userService.GetProfile(r, "profile")
	userService.UpdateProfile(r, "profile")
	userService.DeleteAccount(r, "account")
	userService.GetDashboard(r, "dashboard")
	userService.GetNotifications(r, "notifications")
	userService.MarkNotificationRead(r, "notifications")
	userService.GetPendingInvites(r, "invites")
	userService.AcceptInvite(r, "invites")
	userService.DenyInvite(r, "invites")

	// Organization Management endpoints
	organizationService := service.OrganizationService{
		Route:      "organizations",
		Controller: controller.OrganizationController{},
	}

	organizationService.CreateOrganization(r)
	organizationService.GetOrganizations(r)
	organizationService.GetOrganization(r)
	organizationService.UpdateOrganization(r)
	organizationService.DeleteOrganization(r)
	organizationService.GetOrganizationMembers(r)
	organizationService.InviteMember(r)
	organizationService.GetPendingInvites(r)
	organizationService.RemoveMember(r)
	organizationService.UpdateMemberRole(r)

	// Payment Method endpoints
	paymentMethodService := service.PaymentMethodService{
		Route:      "organizations/:id/payment-methods",
		Controller: controller.PaymentMethodController{},
	}

	paymentMethodService.GetPaymentMethods(r, "")
	paymentMethodService.CreatePaymentMethod(r, "")
	paymentMethodService.UpdatePaymentMethod(r, "")
	paymentMethodService.DeletePaymentMethod(r, "")
	paymentMethodService.SetDefaultPaymentMethod(r, "")

	// Subscription endpoints
	subscriptionService := service.SubscriptionService{
		Route:      "organizations/:id/subscription",
		Controller: controller.SubscriptionController{},
	}

	subscriptionService.GetSubscription(r, "")
	subscriptionService.CreateSubscription(r, "create")
	subscriptionService.ChangePlan(r, "change-plan")
	subscriptionService.CancelSubscription(r, "cancel")
	subscriptionService.ReactivateSubscription(r, "reactivate")
	subscriptionService.GetSubscriptionUsage(r, "usage")

	// Billing endpoints
	billingService := service.BillingService{
		Route:      "organizations/:id/billing",
		Controller: controller.BillingController{},
	}

	billingService.GetBillingOverview(r, "overview")
	billingService.GetBillingHistory(r, "history")

	// Invoice endpoints
	invoiceService := service.BillingService{
		Route:      "organizations/:id/invoices",
		Controller: controller.BillingController{},
	}

	invoiceService.GetInvoices(r, "")

	// Global invoice endpoints (not organization-specific routes)
	globalInvoiceService := service.BillingService{
		Route:      "invoices",
		Controller: controller.BillingController{},
	}

	globalInvoiceService.GetInvoice(r, "")
	globalInvoiceService.DownloadInvoice(r, "")
	globalInvoiceService.PayInvoice(r, "")
}
