package user

import "testlake/model"

type CreateUserRequest struct {
	Email        string             `json:"email" binding:"required,email"`
	Username     string             `json:"username" binding:"required,min=3,max=50"`
	FirstName    *string            `json:"first_name"`
	LastName     *string            `json:"last_name"`
	Password     string             `json:"password" binding:"required,min=6"`
	AuthProvider model.AuthProvider `json:"auth_provider"`
}

type UpdateUserRequest struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	AvatarURL *string `json:"avatar_url"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}

type UpdateStatusRequest struct {
	Status model.UserStatus `json:"status" binding:"required"`
}