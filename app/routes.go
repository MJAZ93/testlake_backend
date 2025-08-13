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
}
