package mailer

import (
	"context"
	"errors"
	"testing"

	"github.com/mikestefanello/pagoda/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockMailClient(t *testing.T) {
	cfg := &config.Config{
		Mail: config.MailConfig{
			FromAddress: "noreply@under100.example.com",
		},
	}
	mock := NewMockMailClient()
	client, err := NewMailClient(cfg, mock)
	require.NoError(t, err)

	ctx := context.Background()

	// Initially empty
	assert.Equal(t, 0, mock.GetSentCount())
	assert.Nil(t, mock.GetLastSent())

	// Send an email
	err = client.Compose().
		To("user@example.com").
		Subject("Password Reset").
		Body("Click here to reset: https://example.com/reset").
		Send(ctx)

	assert.NoError(t, err)
	assert.Equal(t, 1, mock.GetSentCount())

	last := mock.GetLastSent()
	require.NotNil(t, last)
	assert.Equal(t, "user@example.com", last.To)
	assert.Equal(t, "noreply@under100.example.com", last.From)
	assert.Equal(t, "Password Reset", last.Subject)
	assert.Contains(t, last.Body, "https://example.com/reset")

	// Test error simulation
	mock.ShouldFail = true
	mock.FailureError = errors.New("provider temporary outage")

	err = client.Compose().
		To("user2@example.com").
		Subject("Alert").
		Body("Body").
		Send(ctx)

	assert.Error(t, err)
	assert.Equal(t, "provider temporary outage", err.Error())
	assert.Equal(t, 1, mock.GetSentCount())

	// Reset mock
	mock.Reset()
	assert.Equal(t, 0, mock.GetSentCount())
	assert.Nil(t, mock.GetLastSent())
}
