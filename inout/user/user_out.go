package user

import (
	"time"

	"testlake/inout"
	"testlake/model"

	"github.com/google/uuid"
)

type User struct {
	ID              uuid.UUID          `json:"id"`
	Email           string             `json:"email"`
	Username        string             `json:"username"`
	FirstName       *string            `json:"first_name"`
	LastName        *string            `json:"last_name"`
	AvatarURL       *string            `json:"avatar_url"`
	AuthProvider    model.AuthProvider `json:"auth_provider"`
	IsEmailVerified bool               `json:"is_email_verified"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
	LastLoginAt     *time.Time         `json:"last_login_at"`
	Status          model.UserStatus   `json:"status"`
}

type UserOut struct {
	inout.BaseResponse
	Data User `json:"data"`
}

type UserListOut struct {
	inout.BaseResponse
	List []User               `json:"list"`
	Meta inout.PaginationMeta `json:"meta"`
}

type ActivityItem struct {
	ID          uuid.UUID `json:"id"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
}

type DashboardData struct {
	User              User           `json:"user"`
	PersonalProjects  int            `json:"personal_projects"`
	OrganizationCount int            `json:"organization_count"`
	RecentActivity    []ActivityItem `json:"recent_activity"`
}

type DashboardOut struct {
	inout.BaseResponse
	Data DashboardData `json:"data"`
}

type Notification struct {
	ID        uuid.UUID  `json:"id"`
	Title     string     `json:"title"`
	Message   string     `json:"message"`
	Type      string     `json:"type"`
	IsRead    bool       `json:"is_read"`
	CreatedAt time.Time  `json:"created_at"`
	ReadAt    *time.Time `json:"read_at"`
}

type NotificationsOut struct {
	inout.BaseResponse
	Data []Notification `json:"data"`
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
