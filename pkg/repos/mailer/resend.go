package mailer

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/resend/resend-go/v2"
)

type ResendMailClient struct {
	client *resend.Client
}

func NewResendMailClient(apiKey string) *ResendMailClient {
	client := resend.NewClient(apiKey)
	return &ResendMailClient{client: client}
}

func (r *ResendMailClient) Send(email *mail) error {
	params := &resend.SendEmailRequest{
		To:      []string{email.to},
		From:    email.from,
		Subject: email.subject,
	}

	if email.html != "" {
		params.Html = email.html
	} else if email.component != nil {
		params.Html = email.body
	}

	if email.text != "" {
		params.Text = email.text
	} else if params.Html == "" {
		params.Text = email.body
	}

	// 10-second request timeout for outbound transactional delivery
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sent, err := r.client.Emails.SendWithContext(ctx, params)
	if err != nil {
		return fmt.Errorf("resend email delivery rejected: %w", err)
	}

	log.Printf("[RESEND] Email delivered successfully. Resend ID: %s | To: %s | Subject: %s", sent.Id, email.to, email.subject)
	return nil
}

