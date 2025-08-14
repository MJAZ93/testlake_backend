package auth

import "testlake/model"

type SignUpRequest struct {
	Email        string                `json:"email" binding:"required,email"`
	Username     string                `json:"username" binding:"required,min=3,max=100"`
	FirstName    *string               `json:"first_name"`
	LastName     *string               `json:"last_name"`
	Password     string                `json:"password" binding:"required,min=6"`
	AuthProvider model.AuthProvider    `json:"auth_provider" binding:"required"`
}

type SignInRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}