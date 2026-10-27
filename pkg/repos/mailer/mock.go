package mailer

import (
	"fmt"
	"sync"
)

// SentMessage records a delivered message metadata in mock mode
type SentMessage struct {
	From    string
	To      string
	Subject string
	Body    string
	HTML    string
	Text    string
}

// MockMailClient provides a thread-safe test and development mock email client
type MockMailClient struct {
	mu           sync.Mutex
	SentMessages []SentMessage
	ShouldFail   bool
	FailureError error
}

// NewMockMailClient creates a new MockMailClient
func NewMockMailClient() *MockMailClient {
	return &MockMailClient{
		SentMessages: make([]SentMessage, 0),
	}
}

// Send records the email in memory or returns a simulated error
func (m *MockMailClient) Send(email *mail) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ShouldFail {
		if m.FailureError != nil {
			return m.FailureError
		}
		return fmt.Errorf("mock mail delivery failure")
	}

	m.SentMessages = append(m.SentMessages, SentMessage{
		From:    email.from,
		To:      email.to,
		Subject: email.subject,
		Body:    email.body,
		HTML:    email.html,
		Text:    email.text,
	})
	return nil
}


// GetSentCount returns the total count of emails delivered to the mock
func (m *MockMailClient) GetSentCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.SentMessages)
}

// GetLastSent returns the most recent message or nil
func (m *MockMailClient) GetLastSent() *SentMessage {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.SentMessages) == 0 {
		return nil
	}
	msg := m.SentMessages[len(m.SentMessages)-1]
	return &msg
}

// Reset clears the sent message log and error states
func (m *MockMailClient) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.SentMessages = make([]SentMessage, 0)
	m.ShouldFail = false
	m.FailureError = nil
}
