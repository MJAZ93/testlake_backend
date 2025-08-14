package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EmailVerificationToken struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	Token     string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"token"`
	ExpiresAt time.Time      `gorm:"not null" json:"expires_at"`
	IsUsed    bool           `gorm:"default:false" json:"is_used"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationship
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (e *EmailVerificationToken) BeforeCreate(tx *gorm.DB) (err error) {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return
}

func (e *EmailVerificationToken) IsExpired() bool {
	return time.Now().After(e.ExpiresAt)
}

func (e *EmailVerificationToken) IsValid() bool {
	return !e.IsUsed && !e.IsExpired()
}
