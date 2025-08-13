package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthProvider string

const (
	AuthProviderEmail AuthProvider = "email"
	AuthProviderGmail AuthProvider = "gmail"
	AuthProviderApple AuthProvider = "apple"
)

type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusInactive  UserStatus = "inactive"
)

type User struct {
	ID                uuid.UUID     `gorm:"type:uuid;primaryKey" json:"id"`
	Email             string        `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Username          string        `gorm:"type:varchar(100);uniqueIndex;not null" json:"username"`
	FirstName         *string       `gorm:"type:varchar(100)" json:"first_name"`
	LastName          *string       `gorm:"type:varchar(100)" json:"last_name"`
	AvatarURL         *string       `gorm:"type:varchar(500)" json:"avatar_url"`
	AuthProvider      AuthProvider  `gorm:"type:varchar(20);not null" json:"auth_provider"`
	AuthProviderID    *string       `gorm:"type:varchar(255)" json:"auth_provider_id"`
	PasswordHash      *string       `gorm:"type:varchar(255)" json:"-"`
	IsEmailVerified   bool          `gorm:"default:false" json:"is_email_verified"`
	CreatedAt         time.Time     `json:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at"`
	LastLoginAt       *time.Time    `json:"last_login_at"`
	Status            UserStatus    `gorm:"type:varchar(20);default:active" json:"status"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}