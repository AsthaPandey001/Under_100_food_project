package emailservice

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/mail"
	"os"
	"strings"

	"github.com/mikestefanello/pagoda/pkg/repos/mailer"
)

// Service defines the centralized interface for transactional email operations
type Service interface {
	SendWelcomeEmail(ctx context.Context, toEmail, name string) error
	SendVerificationEmail(ctx context.Context, toEmail, name, verifyURL string) error
	SendPasswordResetEmail(ctx context.Context, toEmail, name, resetURL string) error
	SendPasswordChangedNotification(ctx context.Context, toEmail, name string) error
}

// EmailService implements Service with Resend and mock fallback support
type EmailService struct {
	client     *mailer.MailClient
	defaultFrom string
	appBaseURL  string
}

// NewEmailService creates a new centralized transactional email service
func NewEmailService(client *mailer.MailClient, appBaseURL string) *EmailService {
	from := strings.TrimSpace(os.Getenv("EMAIL_FROM"))
	if from == "" {
		from = strings.TrimSpace(os.Getenv("PAGODA_MAIL_FROMADDRESS"))
	}
	if from == "" {
		from = "Food Under ₹100 <onboarding@resend.dev>"
	}

	cleanBaseURL := strings.TrimRight(strings.TrimSpace(appBaseURL), "/")
	if cleanBaseURL == "" {
		cleanBaseURL = "http://localhost:8000"
	}

	return &EmailService{
		client:      client,
		defaultFrom: from,
		appBaseURL:  cleanBaseURL,
	}
}

// validateEmail ensures the recipient is a syntactically valid email address
func validateEmail(email string) error {
	clean := strings.TrimSpace(email)
	if clean == "" {
		return errors.New("recipient email address cannot be empty")
	}
	_, err := mail.ParseAddress(clean)
	if err != nil {
		return fmt.Errorf("invalid recipient email address %q: %w", clean, err)
	}
	return nil
}

// SendWelcomeEmail dispatches registration confirmation / welcome email
func (s *EmailService) SendWelcomeEmail(ctx context.Context, toEmail, name string) error {
	if err := validateEmail(toEmail); err != nil {
		return err
	}

	if s.client == nil {
		log.Printf("[EMAIL SERVICE MOCK] [Welcome] To: %s | Name: %s (Client nil)", toEmail, name)
		return nil
	}

	htmlBody, textBody, subject := BuildWelcomeEmail(name, s.appBaseURL)

	err := s.client.Compose().
		From(s.defaultFrom).
		To(toEmail).
		Subject(subject).
		HTML(htmlBody).
		Text(textBody).
		Send(ctx)

	if err != nil {
		log.Printf("[EMAIL SERVICE ERROR] Failed to send welcome email to %s: %v", toEmail, err)
		return fmt.Errorf("failed to dispatch welcome email: %w", err)
	}

	log.Printf("[EMAIL SERVICE] Welcome email dispatched to %s", toEmail)
	return nil
}

// SendVerificationEmail dispatches email address verification link
func (s *EmailService) SendVerificationEmail(ctx context.Context, toEmail, name, verifyURL string) error {
	if err := validateEmail(toEmail); err != nil {
		return err
	}

	if strings.TrimSpace(verifyURL) == "" {
		return errors.New("verification URL cannot be empty")
	}

	if s.client == nil {
		log.Printf("[EMAIL SERVICE MOCK] [Verify] To: %s | URL: %s", toEmail, verifyURL)
		return nil
	}

	htmlBody, textBody, subject := BuildVerificationEmail(name, verifyURL)

	err := s.client.Compose().
		From(s.defaultFrom).
		To(toEmail).
		Subject(subject).
		HTML(htmlBody).
		Text(textBody).
		Send(ctx)

	if err != nil {
		log.Printf("[EMAIL SERVICE ERROR] Failed to send verification email to %s: %v", toEmail, err)
		return fmt.Errorf("failed to dispatch verification email: %w", err)
	}

	log.Printf("[EMAIL SERVICE] Verification email dispatched to %s", toEmail)
	return nil
}

// SendPasswordResetEmail dispatches password recovery link
func (s *EmailService) SendPasswordResetEmail(ctx context.Context, toEmail, name, resetURL string) error {
	if err := validateEmail(toEmail); err != nil {
		return err
	}

	if strings.TrimSpace(resetURL) == "" {
		return errors.New("password reset URL cannot be empty")
	}

	if s.client == nil {
		log.Printf("[EMAIL SERVICE MOCK] [Reset] To: %s | URL: %s", toEmail, resetURL)
		return nil
	}

	htmlBody, textBody, subject := BuildPasswordResetEmail(name, resetURL)

	err := s.client.Compose().
		From(s.defaultFrom).
		To(toEmail).
		Subject(subject).
		HTML(htmlBody).
		Text(textBody).
		Send(ctx)

	if err != nil {
		log.Printf("[EMAIL SERVICE ERROR] Failed to send password reset email to %s: %v", toEmail, err)
		return fmt.Errorf("failed to dispatch password reset email: %w", err)
	}

	log.Printf("[EMAIL SERVICE] Password reset email dispatched to %s", toEmail)
	return nil
}

// SendPasswordChangedNotification dispatches account security alert after password change
func (s *EmailService) SendPasswordChangedNotification(ctx context.Context, toEmail, name string) error {
	if err := validateEmail(toEmail); err != nil {
		return err
	}

	if s.client == nil {
		log.Printf("[EMAIL SERVICE MOCK] [Security Alert] To: %s | Password changed notice", toEmail)
		return nil
	}

	htmlBody, textBody, subject := BuildPasswordChangedNotification(name, s.appBaseURL)

	err := s.client.Compose().
		From(s.defaultFrom).
		To(toEmail).
		Subject(subject).
		HTML(htmlBody).
		Text(textBody).
		Send(ctx)

	if err != nil {
		log.Printf("[EMAIL SERVICE ERROR] Failed to send password changed alert to %s: %v", toEmail, err)
		return fmt.Errorf("failed to dispatch password changed security alert: %w", err)
	}

	log.Printf("[EMAIL SERVICE] Password changed alert dispatched to %s", toEmail)
	return nil
}
