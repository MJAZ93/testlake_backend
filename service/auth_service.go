package service

import (
	"testlake/controller"

	"github.com/gin-gonic/gin"
)

type AuthService struct {
	Route      string
	Controller controller.AuthController
}

// SignUp godoc
// @Summary User registration
// @Description Create a new user account with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body auth.SignUpRequest true "Registration data"
// @Success 201 {object} auth.SignUpOut
// @Failure 400 {object} inout.BaseResponse
// @Router /api/v1/auth/signup [POST]
func (s AuthService) SignUp(r *gin.RouterGroup, route string) {
	r.POST("/"+s.Route+"/"+route, s.Controller.SignUp)
}

// SignIn godoc
// @Summary User login
// @Description Authenticate user with email and password, returns JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body auth.SignInRequest true "Login credentials"
// @Success 200 {object} auth.SignInOut
// @Failure 401 {object} inout.BaseResponse
// @Router /api/v1/auth/signin [POST]
func (s AuthService) SignIn(r *gin.RouterGroup, route string) {
	r.POST("/"+s.Route+"/"+route, s.Controller.SignIn)
}

// SignOut godoc
// @Summary User logout
// @Description Sign out current user (invalidate JWT token)
// @Tags Authentication
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token" format(Bearer {token})
// @Success 200 {object} inout.BaseResponse
// @Router /api/v1/auth/signout [POST]
func (s AuthService) SignOut(r *gin.RouterGroup, route string) {
	r.POST("/"+s.Route+"/"+route, s.Controller.SignOut)
}

// RefreshToken godoc
// @Summary Refresh JWT token
// @Description Generate a new JWT token from valid existing token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token" format(Bearer {token})
// @Security BearerAuth
// @Success 200 {object} auth.RefreshTokenOut
// @Failure 401 {object} inout.BaseResponse
// @Router /api/v1/auth/refresh [POST]
func (s AuthService) RefreshToken(r *gin.RouterGroup, route string) {
	r.POST("/"+s.Route+"/"+route, s.Controller.RefreshToken)
}

// ForgotPassword godoc
// @Summary Request password reset
// @Description Send password reset email to user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body auth.ForgotPasswordRequest true "Email for password reset"
// @Success 200 {object} inout.BaseResponse
// @Router /api/v1/auth/forgot-password [POST]
func (s AuthService) ForgotPassword(r *gin.RouterGroup, route string) {
	r.POST("/"+s.Route+"/"+route, s.Controller.ForgotPassword)
}

// ResetPassword godoc
// @Summary Reset user password
// @Description Reset password with valid reset token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token" format(Bearer {token})
// @Param request body auth.ResetPasswordRequest true "Password reset data"
// @Success 200 {object} inout.BaseResponse
// @Router /api/v1/auth/reset-password [POST]
func (s AuthService) ResetPassword(r *gin.RouterGroup, route string) {
	r.POST("/"+s.Route+"/"+route, s.Controller.ResetPassword)
}

// VerifyEmail godoc
// @Summary Verify email address
// @Description Verify user email with verification token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param token path string true "Email verification token"
// @Success 200 {object} inout.BaseResponse
// @Router /api/v1/auth/verify-email/{token} [GET]
func (s AuthService) VerifyEmail(r *gin.RouterGroup, route string) {
	r.GET("/"+s.Route+"/"+route+"/:token", s.Controller.VerifyEmail)
}

// ResendEmailConfirmation godoc
// @Summary Resend email confirmation
// @Description Resend email confirmation to user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body auth.ResendEmailConfirmationRequest true "Email to resend confirmation"
// @Success 200 {object} inout.BaseResponse
// @Router /api/v1/auth/resend-email-confirmation [POST]
func (s AuthService) ResendEmailConfirmation(r *gin.RouterGroup, route string) {
	r.POST("/"+s.Route+"/"+route, s.Controller.ResendEmailConfirmation)
}
