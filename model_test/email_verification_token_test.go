package model_test

import (
	"testing"
	"testlake/model"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestEmailVerificationToken_BeforeCreate(t *testing.T) {
	// Create in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// Migrate the schema
	db.AutoMigrate(&model.EmailVerificationToken{})

	token := &model.EmailVerificationToken{
		UserID:    uuid.New(),
		Token:     "test-token",
		ExpiresAt: time.Now().Add(24 * time.Hour),
		IsUsed:    false,
	}

	// ID should be nil initially
	assert.Equal(t, uuid.Nil, token.ID)

	// Create the token - this should trigger BeforeCreate hook
	err = db.Create(token).Error
	assert.NoError(t, err)

	// ID should now be set
	assert.NotEqual(t, uuid.Nil, token.ID)
}

func TestEmailVerificationToken_BeforeCreateWithExistingID(t *testing.T) {
	// Create in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// Migrate the schema
	db.AutoMigrate(&model.EmailVerificationToken{})

	existingID := uuid.New()
	token := &model.EmailVerificationToken{
		ID:        existingID,
		UserID:    uuid.New(),
		Token:     "test-token",
		ExpiresAt: time.Now().Add(24 * time.Hour),
		IsUsed:    false,
	}

	// Create the token - BeforeCreate should not change existing ID
	err = db.Create(token).Error
	assert.NoError(t, err)

	// ID should remain the same
	assert.Equal(t, existingID, token.ID)
}

func TestEmailVerificationToken_IsExpired(t *testing.T) {
	// Test case 1: Token is not expired
	token := &model.EmailVerificationToken{
		ExpiresAt: time.Now().Add(1 * time.Hour), // Expires in 1 hour
	}
	assert.False(t, token.IsExpired())

	// Test case 2: Token is expired
	expiredToken := &model.EmailVerificationToken{
		ExpiresAt: time.Now().Add(-1 * time.Hour), // Expired 1 hour ago
	}
	assert.True(t, expiredToken.IsExpired())

	// Test case 3: Token expires exactly now (edge case)
	nowToken := &model.EmailVerificationToken{
		ExpiresAt: time.Now(),
	}
	// This might be true or false depending on timing, but should not panic
	_ = nowToken.IsExpired()
}

func TestEmailVerificationToken_IsValid(t *testing.T) {
	// Test case 1: Token is valid (not used and not expired)
	validToken := &model.EmailVerificationToken{
		ExpiresAt: time.Now().Add(1 * time.Hour),
		IsUsed:    false,
	}
	assert.True(t, validToken.IsValid())

	// Test case 2: Token is used but not expired
	usedToken := &model.EmailVerificationToken{
		ExpiresAt: time.Now().Add(1 * time.Hour),
		IsUsed:    true,
	}
	assert.False(t, usedToken.IsValid())

	// Test case 3: Token is not used but expired
	expiredToken := &model.EmailVerificationToken{
		ExpiresAt: time.Now().Add(-1 * time.Hour),
		IsUsed:    false,
	}
	assert.False(t, expiredToken.IsValid())

	// Test case 4: Token is both used and expired
	usedExpiredToken := &model.EmailVerificationToken{
		ExpiresAt: time.Now().Add(-1 * time.Hour),
		IsUsed:    true,
	}
	assert.False(t, usedExpiredToken.IsValid())
}

func TestEmailVerificationToken_StructFields(t *testing.T) {
	userID := uuid.New()
	tokenID := uuid.New()
	tokenStr := "test-token-123"
	expiresAt := time.Now().Add(24 * time.Hour)

	token := &model.EmailVerificationToken{
		ID:        tokenID,
		UserID:    userID,
		Token:     tokenStr,
		ExpiresAt: expiresAt,
		IsUsed:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Test all fields are set correctly
	assert.Equal(t, tokenID, token.ID)
	assert.Equal(t, userID, token.UserID)
	assert.Equal(t, tokenStr, token.Token)
	assert.Equal(t, expiresAt, token.ExpiresAt)
	assert.True(t, token.IsUsed)
	assert.NotZero(t, token.CreatedAt)
	assert.NotZero(t, token.UpdatedAt)
}

func TestEmailVerificationToken_DefaultValues(t *testing.T) {
	token := &model.EmailVerificationToken{}

	// Test default values
	assert.Equal(t, uuid.Nil, token.ID)
	assert.Equal(t, uuid.Nil, token.UserID)
	assert.Empty(t, token.Token)
	assert.True(t, token.ExpiresAt.IsZero())
	assert.False(t, token.IsUsed)
	assert.True(t, token.CreatedAt.IsZero())
	assert.True(t, token.UpdatedAt.IsZero())
}

func TestEmailVerificationToken_DatabaseIntegration(t *testing.T) {
	// Create in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// Migrate the schema
	db.AutoMigrate(&model.EmailVerificationToken{})

	userID := uuid.New()
	tokenStr := "integration-test-token"
	expiresAt := time.Now().Add(24 * time.Hour)

	// Create token
	token := &model.EmailVerificationToken{
		UserID:    userID,
		Token:     tokenStr,
		ExpiresAt: expiresAt,
		IsUsed:    false,
	}

	err = db.Create(token).Error
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, token.ID)

	// Retrieve token
	var retrievedToken model.EmailVerificationToken
	err = db.Where("token = ?", tokenStr).First(&retrievedToken).Error
	assert.NoError(t, err)

	assert.Equal(t, token.ID, retrievedToken.ID)
	assert.Equal(t, userID, retrievedToken.UserID)
	assert.Equal(t, tokenStr, retrievedToken.Token)
	assert.False(t, retrievedToken.IsUsed)
	assert.True(t, retrievedToken.IsValid())

	// Update token to used
	err = db.Model(&retrievedToken).Update("is_used", true).Error
	assert.NoError(t, err)

	// Verify update
	var updatedToken model.EmailVerificationToken
	err = db.Where("token = ?", tokenStr).First(&updatedToken).Error
	assert.NoError(t, err)
	assert.True(t, updatedToken.IsUsed)
	assert.False(t, updatedToken.IsValid()) // Should no longer be valid

	// Soft delete
	err = db.Delete(&updatedToken).Error
	assert.NoError(t, err)

	// Should not find deleted token in normal query
	var deletedToken model.EmailVerificationToken
	err = db.Where("token = ?", tokenStr).First(&deletedToken).Error
	assert.Error(t, err) // Should return "record not found"

	// But should find with Unscoped
	err = db.Unscoped().Where("token = ?", tokenStr).First(&deletedToken).Error
	assert.NoError(t, err)
	assert.NotNil(t, deletedToken.DeletedAt)
}
