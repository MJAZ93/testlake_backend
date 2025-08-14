package dao

import (
	"testlake/model"
	"time"

	"github.com/google/uuid"
)

type EmailVerificationDao struct{}

func NewEmailVerificationDao() *EmailVerificationDao {
	return &EmailVerificationDao{}
}

func (dao *EmailVerificationDao) Create(token *model.EmailVerificationToken) error {
	return Database.Create(token).Error
}

func (dao *EmailVerificationDao) GetByToken(tokenStr string) (*model.EmailVerificationToken, error) {
	var token model.EmailVerificationToken
	err := Database.Where("token = ? AND deleted_at IS NULL", tokenStr).
		Preload("User").
		First(&token).Error
	return &token, err
}

func (dao *EmailVerificationDao) MarkAsUsed(tokenStr string) error {
	return Database.Model(&model.EmailVerificationToken{}).
		Where("token = ? AND deleted_at IS NULL", tokenStr).
		Updates(map[string]interface{}{
			"is_used":    true,
			"updated_at": time.Now(),
		}).Error
}

func (dao *EmailVerificationDao) DeleteExpiredTokens() error {
	return Database.Where("expires_at < ?", time.Now()).
		Delete(&model.EmailVerificationToken{}).Error
}

func (dao *EmailVerificationDao) DeleteTokensForUser(userID uuid.UUID) error {
	return Database.Where("user_id = ?", userID).
		Delete(&model.EmailVerificationToken{}).Error
}

func (dao *EmailVerificationDao) GetActiveTokensForUser(userID uuid.UUID) ([]model.EmailVerificationToken, error) {
	var tokens []model.EmailVerificationToken
	err := Database.Where("user_id = ? AND is_used = ? AND expires_at > ? AND deleted_at IS NULL",
		userID, false, time.Now()).Find(&tokens).Error
	return tokens, err
}
