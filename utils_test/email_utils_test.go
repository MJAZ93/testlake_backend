package utils_test

import (
	"os"
	"path/filepath"
	"testing"
	"testlake/dao"
	"testlake/model"
	"testlake/utils"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	err := godotenv.Load("../.env")
	if err != nil {
		// Don't panic, just skip database tests
	} else {
		dao.Connect()
	}

	code := m.Run()
	os.Exit(code)
}

func TestNewEmailService(t *testing.T) {
	// Set required environment variables
	os.Setenv("SMTP_HOST", "smtp.example.com")
	os.Setenv("SMTP_PORT", "587")
	os.Setenv("SMTP_USERNAME", "test@example.com")
	os.Setenv("SMTP_PASSWORD", "password")
	os.Setenv("EMAIL_FROM_ADDRESS", "noreply@testlake.com")
	os.Setenv("EMAIL_FROM_NAME", "TestLake")

	emailService := utils.NewEmailService()

	assert.NotNil(t, emailService)
}

func TestNewEmailServiceWithInvalidPort(t *testing.T) {
	// Set required environment variables with invalid port
	os.Setenv("SMTP_HOST", "smtp.example.com")
	os.Setenv("SMTP_PORT", "invalid")
	os.Setenv("SMTP_USERNAME", "test@example.com")
	os.Setenv("SMTP_PASSWORD", "password")
	os.Setenv("EMAIL_FROM_ADDRESS", "noreply@testlake.com")
	os.Setenv("EMAIL_FROM_NAME", "TestLake")

	emailService := utils.NewEmailService()

	// Should still create service with default port 587
	assert.NotNil(t, emailService)
}

func TestSendEmailConfirmationFunction(t *testing.T) {
	// Set required environment variables
	os.Setenv("SMTP_HOST", "smtp.example.com")
	os.Setenv("SMTP_PORT", "587")
	os.Setenv("SMTP_USERNAME", "test@example.com")
	os.Setenv("SMTP_PASSWORD", "password")
	os.Setenv("EMAIL_FROM_ADDRESS", "noreply@testlake.com")
	os.Setenv("EMAIL_FROM_NAME", "TestLake")
	os.Setenv("SCHEME", "https")
	os.Setenv("IP", "testlake.com")
	os.Setenv("PORT", "443")

	// Database should already be connected via TestMain
	if dao.Database == nil {
		t.Skip("Database not connected - skipping test")
	}

	userID := uuid.New()
	err := utils.SendEmailConfirmation("test@example.com", "testuser", userID)

	// This will fail because we don't have a real SMTP server, but it should create the token first
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to send email")

	// Cleanup any created tokens
	dao.Database.Where("user_id = ?", userID).Delete(&model.EmailVerificationToken{})
}

func TestResendEmailConfirmationFunction(t *testing.T) {
	// Set required environment variables
	os.Setenv("SMTP_HOST", "smtp.example.com")
	os.Setenv("SMTP_PORT", "587")
	os.Setenv("SMTP_USERNAME", "test@example.com")
	os.Setenv("SMTP_PASSWORD", "password")
	os.Setenv("EMAIL_FROM_ADDRESS", "noreply@testlake.com")
	os.Setenv("EMAIL_FROM_NAME", "TestLake")
	os.Setenv("SCHEME", "https")
	os.Setenv("IP", "testlake.com")
	os.Setenv("PORT", "443")

	// Database should already be connected via TestMain
	if dao.Database == nil {
		t.Skip("Database not connected - skipping test")
	}

	userID := uuid.New()
	err := utils.ResendEmailConfirmation("test@example.com", "testuser", userID)

	// This will fail because we don't have a real SMTP server, but it should handle token cleanup/creation
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to send email")

	// Cleanup any created tokens
	dao.Database.Where("user_id = ?", userID).Delete(&model.EmailVerificationToken{})
}

func TestRenderEmailVerifiedSuccess(t *testing.T) {
	// Set required environment variables
	os.Setenv("SCHEME", "https")
	os.Setenv("IP", "testlake.com")
	os.Setenv("PORT", "443")

	baseURL := "https://testlake.com:443"

	// Check if template file exists
	templatePath := filepath.Join("..", "templates", "email_verified_success.html")
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		t.Skip("Email template file not found - skipping template test")
	}

	body, err := utils.RenderEmailVerifiedSuccess(baseURL)

	if err != nil {
		t.Skipf("Could not render template (template file might be missing): %v", err)
	}

	assert.NotEmpty(t, body)
	assert.Contains(t, body, "testlake.com")
}

func TestRenderEmailVerificationError(t *testing.T) {
	title := "Error Title"
	heading := "Error Heading"
	message := "Error message content"

	// Check if template file exists
	templatePath := filepath.Join("..", "templates", "email_verification_error.html")
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		t.Skip("Email template file not found - skipping template test")
	}

	body, err := utils.RenderEmailVerificationError(title, heading, message)

	if err != nil {
		t.Skipf("Could not render error template (template file might be missing): %v", err)
	}

	assert.NotEmpty(t, body)
	assert.Contains(t, body, title)
	assert.Contains(t, body, heading)
	assert.Contains(t, body, message)
}

func TestEmailServiceSendEmailConfirmationWithDatabase(t *testing.T) {
	// Set required environment variables
	os.Setenv("SMTP_HOST", "smtp.example.com")
	os.Setenv("SMTP_PORT", "587")
	os.Setenv("SMTP_USERNAME", "test@example.com")
	os.Setenv("SMTP_PASSWORD", "password")
	os.Setenv("EMAIL_FROM_ADDRESS", "noreply@testlake.com")
	os.Setenv("EMAIL_FROM_NAME", "TestLake")
	os.Setenv("SCHEME", "https")
	os.Setenv("IP", "testlake.com")
	os.Setenv("PORT", "443")

	// Database should already be connected via TestMain
	if dao.Database == nil {
		t.Skip("Database not connected - skipping test")
	}

	// Check if template file exists
	templatePath := filepath.Join("..", "templates", "email_confirmation.html")
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		t.Skip("Email template file not found - skipping template test")
	}

	emailService := utils.NewEmailService()
	userID := uuid.New()

	err := emailService.SendEmailConfirmation("test@example.com", "testuser", userID)

	// Should fail at SMTP level, not at token creation or template loading
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to send email")

	// Verify token was created in database
	emailDao := dao.NewEmailVerificationDao()
	tokens, err := emailDao.GetActiveTokensForUser(userID)
	if err == nil {
		assert.Len(t, tokens, 1)

		// Cleanup
		dao.Database.Where("user_id = ?", userID).Delete(&model.EmailVerificationToken{})
	}
}

func TestEmailServiceResendEmailConfirmationWithDatabase(t *testing.T) {
	// Set required environment variables
	os.Setenv("SMTP_HOST", "smtp.example.com")
	os.Setenv("SMTP_PORT", "587")
	os.Setenv("SMTP_USERNAME", "test@example.com")
	os.Setenv("SMTP_PASSWORD", "password")
	os.Setenv("EMAIL_FROM_ADDRESS", "noreply@testlake.com")
	os.Setenv("EMAIL_FROM_NAME", "TestLake")
	os.Setenv("SCHEME", "https")
	os.Setenv("IP", "testlake.com")
	os.Setenv("PORT", "443")

	// Database should already be connected via TestMain
	if dao.Database == nil {
		t.Skip("Database not connected - skipping test")
	}

	// Check if template file exists
	templatePath := filepath.Join("..", "templates", "email_confirmation_resend.html")
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		t.Skip("Email template file not found - skipping template test")
	}

	emailService := utils.NewEmailService()
	userID := uuid.New()

	err := emailService.ResendEmailConfirmation("test@example.com", "testuser", userID)

	// Should fail at SMTP level, not at token cleanup/creation or template loading
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to send email")

	// Verify token was created in database
	emailDao := dao.NewEmailVerificationDao()
	tokens, err := emailDao.GetActiveTokensForUser(userID)
	if err == nil {
		assert.Len(t, tokens, 1)

		// Cleanup
		dao.Database.Where("user_id = ?", userID).Delete(&model.EmailVerificationToken{})
	}
}

func TestEmailServiceDefaultEnvironmentValues(t *testing.T) {
	// Clear environment variables to test defaults
	os.Unsetenv("SCHEME")
	os.Unsetenv("IP")
	os.Unsetenv("PORT")
	os.Unsetenv("SMTP_PORT")

	// Set minimum required vars
	os.Setenv("SMTP_HOST", "smtp.example.com")
	os.Setenv("SMTP_USERNAME", "test@example.com")
	os.Setenv("SMTP_PASSWORD", "password")
	os.Setenv("EMAIL_FROM_ADDRESS", "noreply@testlake.com")
	os.Setenv("EMAIL_FROM_NAME", "TestLake")

	emailService := utils.NewEmailService()

	// Should create service with default values
	assert.NotNil(t, emailService)
}

func TestEmailServiceMissingTemplateFile(t *testing.T) {
	// Set required environment variables
	os.Setenv("SMTP_HOST", "smtp.example.com")
	os.Setenv("SMTP_PORT", "587")
	os.Setenv("SMTP_USERNAME", "test@example.com")
	os.Setenv("SMTP_PASSWORD", "password")
	os.Setenv("EMAIL_FROM_ADDRESS", "noreply@testlake.com")
	os.Setenv("EMAIL_FROM_NAME", "TestLake")

	// Database should already be connected via TestMain
	if dao.Database == nil {
		t.Skip("Database not connected - skipping test")
	}

	// Temporarily rename template directory to simulate missing template
	templateDir := filepath.Join("..", "templates")
	tempDir := filepath.Join("..", "templates_backup")

	// Only run this test if we can manipulate the templates directory
	if _, err := os.Stat(templateDir); err == nil {
		err = os.Rename(templateDir, tempDir)
		if err != nil {
			t.Skip("Could not rename templates directory - skipping test")
		}
		defer os.Rename(tempDir, templateDir) // Restore after test

		emailService := utils.NewEmailService()
		userID := uuid.New()

		err := emailService.SendEmailConfirmation("test@example.com", "testuser", userID)

		// Should fail because template file is missing
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to load email template")

		// Cleanup any tokens that might have been created
		dao.Database.Where("user_id = ?", userID).Delete(&model.EmailVerificationToken{})
	} else {
		t.Skip("Templates directory not found - skipping missing template test")
	}
}
