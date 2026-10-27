package emailservice

import (
	"context"
	"testing"

	"github.com/mikestefanello/pagoda/config"
	"github.com/mikestefanello/pagoda/pkg/repos/mailer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


func setupTestEmailService(t *testing.T) (*EmailService, *mailer.MockMailClient) {
	mockSender := mailer.NewMockMailClient()
	cfg := &config.Config{
		Mail: config.MailConfig{
			FromAddress: "auth@foodunder100.in",
		},
	}
	client, err := mailer.NewMailClient(cfg, mockSender)
	require.NoError(t, err)

	svc := NewEmailService(client, "http://localhost:8000")
	return svc, mockSender
}

func TestEmailService_Validation(t *testing.T) {
	svc, _ := setupTestEmailService(t)
	ctx := context.Background()

	t.Run("Empty_Email_Fails", func(t *testing.T) {
		err := svc.SendWelcomeEmail(ctx, "", "Arjun")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "recipient email address cannot be empty")
	})

	t.Run("Invalid_Email_Format_Fails", func(t *testing.T) {
		err := svc.SendWelcomeEmail(ctx, "invalid-email-address", "Arjun")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid recipient email address")
	})

	t.Run("Empty_Verification_URL_Fails", func(t *testing.T) {
		err := svc.SendVerificationEmail(ctx, "user@example.com", "Arjun", "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "verification URL cannot be empty")
	})

	t.Run("Empty_Reset_URL_Fails", func(t *testing.T) {
		err := svc.SendPasswordResetEmail(ctx, "user@example.com", "Arjun", "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "password reset URL cannot be empty")
	})
}

func TestEmailService_WelcomeEmail(t *testing.T) {
	svc, mockSender := setupTestEmailService(t)
	ctx := context.Background()

	err := svc.SendWelcomeEmail(ctx, "newcustomer@example.com", "Rohan Sharma")
	require.NoError(t, err)

	assert.Equal(t, 1, mockSender.GetSentCount())
	lastMsg := mockSender.GetLastSent()
	require.NotNil(t, lastMsg)

	assert.Equal(t, "newcustomer@example.com", lastMsg.To)
	assert.Equal(t, "Welcome to Food Under ₹100!", lastMsg.Subject)
	assert.Contains(t, lastMsg.HTML, "Rohan Sharma")
	assert.Contains(t, lastMsg.HTML, "Food Under ₹100")
	assert.Contains(t, lastMsg.HTML, "Explore Dishes Under ₹100")
	assert.Contains(t, lastMsg.Text, "Hello Rohan Sharma,")
	assert.Contains(t, lastMsg.Text, "Welcome to Food Under ₹100!")
}

func TestEmailService_VerificationEmail(t *testing.T) {
	svc, mockSender := setupTestEmailService(t)
	ctx := context.Background()

	verifyURL := "http://localhost:8000/verify-email?token=abcdef1234567890"
	err := svc.SendVerificationEmail(ctx, "verify.me@example.com", "Priya Nair", verifyURL)
	require.NoError(t, err)

	assert.Equal(t, 1, mockSender.GetSentCount())
	lastMsg := mockSender.GetLastSent()
	require.NotNil(t, lastMsg)

	assert.Equal(t, "verify.me@example.com", lastMsg.To)
	assert.Equal(t, "Verify Your Email Address — Food Under ₹100", lastMsg.Subject)
	assert.Contains(t, lastMsg.HTML, "Priya Nair")
	assert.Contains(t, lastMsg.HTML, verifyURL)
	assert.Contains(t, lastMsg.HTML, "30 minutes")
	assert.Contains(t, lastMsg.Text, verifyURL)
	assert.Contains(t, lastMsg.Text, "used once")
}

func TestEmailService_PasswordResetEmail(t *testing.T) {
	svc, mockSender := setupTestEmailService(t)
	ctx := context.Background()

	resetURL := "http://localhost:8000/reset-password?token=resettoken999&user=10"
	err := svc.SendPasswordResetEmail(ctx, "recover@example.com", "Vikram Patel", resetURL)
	require.NoError(t, err)

	assert.Equal(t, 1, mockSender.GetSentCount())
	lastMsg := mockSender.GetLastSent()
	require.NotNil(t, lastMsg)

	assert.Equal(t, "recover@example.com", lastMsg.To)
	assert.Equal(t, "Reset Your Password — Food Under ₹100", lastMsg.Subject)
	assert.Contains(t, lastMsg.HTML, "Vikram Patel")
	assert.Contains(t, lastMsg.HTML, "token=resettoken999")
	assert.Contains(t, lastMsg.HTML, "30 minutes")
	assert.Contains(t, lastMsg.Text, resetURL)
}


func TestEmailService_PasswordChangedNotification(t *testing.T) {
	svc, mockSender := setupTestEmailService(t)
	ctx := context.Background()

	err := svc.SendPasswordChangedNotification(ctx, "security@example.com", "Ananya Rao")
	require.NoError(t, err)

	assert.Equal(t, 1, mockSender.GetSentCount())
	lastMsg := mockSender.GetLastSent()
	require.NotNil(t, lastMsg)

	assert.Equal(t, "security@example.com", lastMsg.To)
	assert.Equal(t, "Security Alert: Password Changed — Food Under ₹100", lastMsg.Subject)
	assert.Contains(t, lastMsg.HTML, "Ananya Rao")
	assert.Contains(t, lastMsg.HTML, "Sign In with New Password")
	assert.Contains(t, lastMsg.HTML, "forgot-password")
	assert.Contains(t, lastMsg.Text, "invalidated")
}

func TestEmailService_HTMLInjectionProtection(t *testing.T) {
	// Test that malicious names are safely escaped in HTML output
	maliciousName := "<script>alert('xss')</script><b>Bold</b>"
	htmlBody, _, _ := BuildWelcomeEmail(maliciousName, "http://localhost:8000")

	assert.NotContains(t, htmlBody, "<script>")
	assert.Contains(t, htmlBody, "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;")
}

func TestEmailService_SenderFailurePropagation(t *testing.T) {
	svc, mockSender := setupTestEmailService(t)
	mockSender.ShouldFail = true
	ctx := context.Background()

	err := svc.SendWelcomeEmail(ctx, "fail@example.com", "Test User")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "mock mail delivery failure")
}
