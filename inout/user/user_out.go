package user

import (
	"time"

	"testlake/inout"
	"testlake/model"

	"github.com/google/uuid"
)

type User struct {
	ID              uuid.UUID        `json:"id"`
	Email           string           `json:"email"`
	Username        string           `json:"username"`
	FirstName       *string          `json:"first_name"`
	LastName        *string          `json:"last_name"`
	AvatarURL       *string          `json:"avatar_url"`
	AuthProvider    model.AuthProvider `json:"auth_provider"`
	IsEmailVerified bool             `json:"is_email_verified"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
	LastLoginAt     *time.Time       `json:"last_login_at"`
	Status          model.UserStatus `json:"status"`
}

type UserOut struct {
	inout.BaseResponse
	Data User `json:"data"`
}

type UserListOut struct {
	inout.BaseResponse
	List []User                `json:"list"`
	Meta inout.PaginationMeta  `json:"meta"`
}

type LoginOut struct {
	inout.BaseResponse
	Data struct {
		User  User   `json:"user"`
		Token string `json:"token"`
	} `json:"data"`
}

func FromModel(user *model.User) User {
	return User{
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

func FromModelList(users []model.User) []User {
	result := make([]User, len(users))
	for i, user := range users {
		result[i] = FromModel(&user)
	}
	return result
}