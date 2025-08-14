package controller

import (
	"errors"
	"net/http"

	"testlake/dao"
	"testlake/inout"
	"testlake/inout/auth"
	"testlake/model"
	"testlake/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthController struct{}

// SignUp creates a new user account
func (controller AuthController) SignUp(context *gin.Context) {
	var request auth.SignUpRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		utils.ReportBadRequest(context, "Invalid request data")
		return
	}

	userDao := dao.NewUserDao()

	// Check if email already exists
	emailExists, err := userDao.EmailExists(request.Email)
	if err != nil {
		utils.ReportInternalServerError(context, "Database error")
		return
	}
	if emailExists {
		utils.ReportBadRequest(context, "Email already exists")
		return
	}

	// Check if username already exists
	usernameExists, err := userDao.UsernameExists(request.Username)
	if err != nil {
		utils.ReportInternalServerError(context, "Database error")
		return
	}
	if usernameExists {
		utils.ReportBadRequest(context, "Username already exists")
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		utils.ReportInternalServerError(context, "Failed to process password")
		return
	}

	// Create user
	newUser := &model.User{
		Email:           request.Email,
		Username:        request.Username,
		FirstName:       request.FirstName,
		LastName:        request.LastName,
		AuthProvider:    request.AuthProvider,
		PasswordHash:    &hashedPassword,
		Status:          model.UserStatusActive,
		IsEmailVerified: false,
	}

	if err := userDao.Create(newUser); err != nil {
		utils.ReportInternalServerError(context, "Failed to create user")
		return
	}

	// Send email confirmation
	if err := utils.SendEmailConfirmation(newUser.Email, newUser.Username, newUser.ID); err != nil {
		// Log error but don't fail registration
		// In production, you might want to use a proper logger
		// log.Printf("Failed to send confirmation email: %v", err)
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(newUser.ID, newUser.Email, newUser.Username)
	if err != nil {
		utils.ReportInternalServerError(context, "Failed to generate token")
		return
	}

	response := auth.SignUpOut{
		BaseResponse: inout.BaseResponse{
			ErrorCode:        0,
			ErrorDescription: "Success",
		},
		Data: auth.AuthData{
			Token: token,
			User:  auth.UserFromModel(newUser),
		},
	}

	context.JSON(http.StatusCreated, response)
}

// SignIn authenticates user and returns JWT token
func (controller AuthController) SignIn(context *gin.Context) {
	var request auth.SignInRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		utils.ReportBadRequest(context, "Invalid request data")
		return
	}

	userDao := dao.NewUserDao()
	foundUser, err := userDao.GetByEmail(request.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ReportUnauthorized(context, "Invalid credentials")
		} else {
			utils.ReportInternalServerError(context, "Database error")
		}
		return
	}

	// Check password
	if foundUser.PasswordHash == nil || !utils.CheckPasswordHash(request.Password, *foundUser.PasswordHash) {
		utils.ReportUnauthorized(context, "Invalid credentials")
		return
	}

	// Check user status
	if foundUser.Status != model.UserStatusActive {
		utils.ReportForbidden(context, "Account is not active")
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(foundUser.ID, foundUser.Email, foundUser.Username)
	if err != nil {
		utils.ReportInternalServerError(context, "Failed to generate token")
		return
	}

	// Update last login
	userDao.UpdateLastLogin(foundUser.ID)

	response := auth.SignInOut{
		BaseResponse: inout.BaseResponse{
			ErrorCode:        0,
			ErrorDescription: "Success",
		},
		Data: auth.AuthData{
			Token: token,
			User:  auth.UserFromModel(foundUser),
		},
	}

	context.JSON(http.StatusOK, response)
}

// SignOut invalidates the current JWT token
func (controller AuthController) SignOut(context *gin.Context) {
	// For JWT-based auth, client-side token removal is sufficient
	// In a more sophisticated setup, we could maintain a token blacklist
	response := inout.BaseResponse{
		ErrorCode:        0,
		ErrorDescription: "Signed out successfully",
	}

	context.JSON(http.StatusOK, response)
}

// RefreshToken generates a new JWT token from a valid existing token
func (controller AuthController) RefreshToken(context *gin.Context) {
	userID, err := utils.ExtractUserID(context)
	if err != nil {
		utils.ReportUnauthorized(context, "Invalid token")
		return
	}

	userDao := dao.NewUserDao()
	foundUser, err := userDao.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ReportUnauthorized(context, "User not found")
		} else {
			utils.ReportInternalServerError(context, "Database error")
		}
		return
	}

	// Check user status
	if foundUser.Status != model.UserStatusActive {
		utils.ReportForbidden(context, "Account is not active")
		return
	}

	// Generate new JWT token
	newToken, err := utils.GenerateJWT(foundUser.ID, foundUser.Email, foundUser.Username)
	if err != nil {
		utils.ReportInternalServerError(context, "Failed to generate token")
		return
	}

	response := auth.RefreshTokenOut{
		BaseResponse: inout.BaseResponse{
			ErrorCode:        0,
			ErrorDescription: "Success",
		},
		Data: auth.TokenData{
			Token: newToken,
		},
	}

	context.JSON(http.StatusOK, response)
}

// ForgotPassword initiates password reset process
func (controller AuthController) ForgotPassword(context *gin.Context) {
	var request auth.ForgotPasswordRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		utils.ReportBadRequest(context, "Invalid request data")
		return
	}

	userDao := dao.NewUserDao()
	_, err := userDao.GetByEmail(request.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// For security, don't reveal if email exists
			response := inout.BaseResponse{
				ErrorCode:        0,
				ErrorDescription: "If the email exists, a password reset link will be sent",
			}
			context.JSON(http.StatusOK, response)
			return
		} else {
			utils.ReportInternalServerError(context, "Database error")
			return
		}
	}

	// TODO: Implement email sending logic for password reset
	// For now, just return success

	response := inout.BaseResponse{
		ErrorCode:        0,
		ErrorDescription: "If the email exists, a password reset link will be sent",
	}

	context.JSON(http.StatusOK, response)
}

// ResetPassword resets user password with a valid reset token
func (controller AuthController) ResetPassword(context *gin.Context) {
	var request auth.ResetPasswordRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		utils.ReportBadRequest(context, "Invalid request data")
		return
	}

	// TODO: Implement reset token validation
	// For now, this is a placeholder implementation

	response := inout.BaseResponse{
		ErrorCode:        0,
		ErrorDescription: "Password reset successfully",
	}

	context.JSON(http.StatusOK, response)
}

// VerifyEmail verifies user email with verification token
func (controller AuthController) VerifyEmail(context *gin.Context) {
	token := context.Param("token")
	if token == "" {
		utils.ReportBadRequest(context, "Verification token required")
		return
	}

	// TODO: Implement email verification token validation
	// For now, this is a placeholder implementation

	response := inout.BaseResponse{
		ErrorCode:        0,
		ErrorDescription: "Email verified successfully",
	}

	context.JSON(http.StatusOK, response)
}

// ResendEmailConfirmation resends email confirmation to user
func (controller AuthController) ResendEmailConfirmation(context *gin.Context) {
	var request auth.ResendEmailConfirmationRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		utils.ReportBadRequest(context, "Invalid request data")
		return
	}

	userDao := dao.NewUserDao()
	foundUser, err := userDao.GetByEmail(request.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// For security, don't reveal if email exists
			response := inout.BaseResponse{
				ErrorCode:        0,
				ErrorDescription: "If the email exists and is not verified, a confirmation email will be sent",
			}
			context.JSON(http.StatusOK, response)
			return
		} else {
			utils.ReportInternalServerError(context, "Database error")
			return
		}
	}

	// Check if email is already verified
	if foundUser.IsEmailVerified {
		utils.ReportBadRequest(context, "Email is already verified")
		return
	}

	// Send email confirmation
	if err := utils.ResendEmailConfirmation(foundUser.Email, foundUser.Username, foundUser.ID); err != nil {
		utils.ReportInternalServerError(context, "Failed to send confirmation email")
		return
	}

	response := inout.BaseResponse{
		ErrorCode:        0,
		ErrorDescription: "Confirmation email sent successfully",
	}

	context.JSON(http.StatusOK, response)
}
