package auth

import (
	"time"
	"testlake/inout"
	"testlake/model"
	
	"github.com/google/uuid"
)

type AuthUser struct {
	ID                uuid.UUID            `json:"id"`
	Email             string               `json:"email"`
	Username          string               `json:"username"`
	FirstName         *string              `json:"first_name"`
	LastName          *string              `json:"last_name"`
	AvatarURL         *string              `json:"avatar_url"`
	AuthProvider      model.AuthProvider   `json:"auth_provider"`
	IsEmailVerified   bool                 `json:"is_email_verified"`
	CreatedAt         time.Time            `json:"created_at"`
	UpdatedAt         time.Time            `json:"updated_at"`
	LastLoginAt       *time.Time           `json:"last_login_at"`
	Status            model.UserStatus     `json:"status"`
}

type AuthData struct {
	Token string   `json:"token"`
	User  AuthUser `json:"user"`
}

type TokenData struct {
	Token string `json:"token"`
}

type SignUpOut struct {
	inout.BaseResponse
	Data AuthData `json:"data"`
}

type SignInOut struct {
	inout.BaseResponse
	Data AuthData `json:"data"`
}

type RefreshTokenOut struct {
	inout.BaseResponse
	Data TokenData `json:"data"`
}

func UserFromModel(user *model.User) AuthUser {
	return AuthUser{
		ID:              user.ID,
		Email:           user.Email,
		Username:        user.Username,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		AvatarURL:       user.AvatarURL,
		AuthProvider:    user.AuthProvider,
		IsEmailVerified: user.IsEmailVerified,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
		LastLoginAt:     user.LastLoginAt,
		Status:          user.Status,
	}
}