package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/google/uuid"
	"gopkg.in/gomail.v2"
)

func SendEmailConfirmation(email, username string, userID uuid.UUID) error {
	emailService := NewEmailService()
	return emailService.SendEmailConfirmation(email, username, userID)
}

func ResendEmailConfirmation(email, username string, userID uuid.UUID) error {
	emailService := NewEmailService()
	return emailService.ResendEmailConfirmation(email, username, userID)
}

type EmailService struct {
	dialer *gomail.Dialer
	from   string
	name   string
}

type EmailTemplateData struct {
	Username string
	Token    string
	BaseURL  string
}

func NewEmailService() *EmailService {
	host := os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	fromAddress := os.Getenv("EMAIL_FROM_ADDRESS")
	fromName := os.Getenv("EMAIL_FROM_NAME")

	port, err := strconv.Atoi(portStr)
	if err != nil {
		port = 587
	}

	dialer := gomail.NewDialer(host, port, username, password)

	return &EmailService{
		dialer: dialer,
		from:   fromAddress,
		name:   fromName,
	}
}

func (e *EmailService) SendEmailConfirmation(email, username string, userID uuid.UUID) error {
	token := uuid.New().String()
	subject := "Welcome to TestLake - Please Confirm Your Email"

	data := EmailTemplateData{
		Username: username,
		Token:    token,
		BaseURL:  e.getBaseURL(),
	}

	body, err := e.loadTemplate("email_confirmation.html", data)
	if err != nil {
		e.logError("Failed to load email confirmation template", err, email)
		return fmt.Errorf("failed to load email template: %w", err)
	}

	return e.sendEmail(email, subject, body)
}

func (e *EmailService) ResendEmailConfirmation(email, username string, userID uuid.UUID) error {
	token := uuid.New().String()
	subject := "TestLake - Email Confirmation Resent"

	data := EmailTemplateData{
		Username: username,
		Token:    token,
		BaseURL:  e.getBaseURL(),
	}

	body, err := e.loadTemplate("email_confirmation_resend.html", data)
	if err != nil {
		e.logError("Failed to load email confirmation resend template", err, email)
		return fmt.Errorf("failed to load email template: %w", err)
	}

	return e.sendEmail(email, subject, body)
}

func (e *EmailService) sendEmail(to, subject, body string) error {
	message := gomail.NewMessage()
	message.SetHeader("From", message.FormatAddress(e.from, e.name))
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", body)

	if err := e.dialer.DialAndSend(message); err != nil {
		e.logError("Failed to send email", err, to)
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (e *EmailService) loadTemplate(templateName string, data EmailTemplateData) (string, error) {
	templatePath := filepath.Join("templates", templateName)

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to parse template %s: %w", templateName, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", templateName, err)
	}

	return buf.String(), nil
}

func (e *EmailService) logError(message string, err error, email string) {
	// Ensure logs directory exists
	if err := os.MkdirAll("logs", 0755); err != nil {
		log.Printf("Failed to create logs directory: %v", err)
		return
	}

	// Open or create log file
	logFile, fileErr := os.OpenFile("logs/email_errors.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if fileErr != nil {
		log.Printf("Failed to open email error log file: %v", fileErr)
		return
	}
	defer logFile.Close()

	// Create logger for the file
	logger := log.New(logFile, "", log.LstdFlags)

	// Log the error
	logger.Printf("[ERROR] %s | Email: %s | Error: %v", message, email, err)

	// Also log to console for immediate visibility
	log.Printf("[EMAIL ERROR] %s | Email: %s | Error: %v", message, email, err)
}

func (e *EmailService) getBaseURL() string {
	scheme := os.Getenv("SCHEME")
	if scheme == "" {
		scheme = "http"
	}

	ip := os.Getenv("IP")
	if ip == "" {
		ip = "localhost"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	return fmt.Sprintf("%s://%s:%s", scheme, ip, port)
}
