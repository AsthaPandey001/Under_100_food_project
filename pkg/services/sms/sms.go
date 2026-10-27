package sms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/mikestefanello/pagoda/config"
)

// Provider defines the clean interface for OTP delivery across development and production
type Provider interface {
	Name() string
	SendOTP(ctx context.Context, phoneE164, otp, portal string) error
	IsSimulated() bool
}

// SimulatedProvider provides a ₹0 zero-cost development and testing mock OTP dispatcher
type SimulatedProvider struct {
	mu          sync.RWMutex
	sentOTPs    map[string]string // phoneE164 -> last OTP
	environment string
}

// NewSimulatedProvider creates a thread-safe development simulated OTP provider
func NewSimulatedProvider(env string) *SimulatedProvider {
	return &SimulatedProvider{
		sentOTPs:    make(map[string]string),
		environment: env,
	}
}

func (s *SimulatedProvider) Name() string {
	return "SimulatedDevOTP"
}

func (s *SimulatedProvider) IsSimulated() bool {
	return true
}

func (s *SimulatedProvider) SendOTP(ctx context.Context, phoneE164, otp, portal string) error {
	// Security rule: fail immediately if simulated provider is somehow invoked in production
	if s.environment == string(config.EnvProduction) {
		return fmt.Errorf("security violation: simulated OTP provider cannot be used in production environment")
	}

	s.mu.Lock()
	s.sentOTPs[phoneE164] = otp
	s.mu.Unlock()

	// Clearly state in logs that this is a development mock and no SMS was charged
	log.Printf("[DEV SIMULATED OTP] [Portal: %s] Target: %s | Verification Code: %s (Simulated - no real SMS sent, zero cost)", portal, phoneE164, otp)
	return nil
}

// GetLastSentOTP returns the latest sent OTP for testing assertions
func (s *SimulatedProvider) GetLastSentOTP(phoneE164 string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.sentOTPs[phoneE164]
}

// Clear resets recorded simulated OTPs
func (s *SimulatedProvider) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sentOTPs = make(map[string]string)
}

// Fast2SMSProvider implements real SMS delivery via Fast2SMS Indian OTP gateway
type Fast2SMSProvider struct {
	apiKey string
	client *http.Client
}

func NewFast2SMSProvider(apiKey string) *Fast2SMSProvider {
	return &Fast2SMSProvider{
		apiKey: apiKey,
		client: &http.Client{Timeout: 6 * time.Second},
	}
}

func (f *Fast2SMSProvider) Name() string {
	return "Fast2SMS"
}

func (f *Fast2SMSProvider) IsSimulated() bool {
	return false
}

func (f *Fast2SMSProvider) SendOTP(ctx context.Context, phoneE164, otp, portal string) error {
	clean := CleanIndianPhone(phoneE164)
	payload, err := json.Marshal(map[string]interface{}{
		"route":            "otp",
		"variables_values": otp,
		"numbers":          clean,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://www.fast2sms.com/dev/bulkV2", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	req.Header.Set("authorization", f.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := f.client.Do(req)
	if err != nil {
		return fmt.Errorf("Fast2SMS request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Fast2SMS gateway returned HTTP %d: %s", resp.StatusCode, string(body))
	}
	return nil
}

// TwilioProvider implements carrier SMS dispatch via Twilio API
type TwilioProvider struct {
	accountSID string
	authToken  string
	fromNumber string
	client     *http.Client
}

func NewTwilioProvider(sid, token, from string) *TwilioProvider {
	return &TwilioProvider{
		accountSID: sid,
		authToken:  token,
		fromNumber: from,
		client:     &http.Client{Timeout: 6 * time.Second},
	}
}

func (t *TwilioProvider) Name() string {
	return "Twilio"
}

func (t *TwilioProvider) IsSimulated() bool {
	return false
}

func (t *TwilioProvider) SendOTP(ctx context.Context, phoneE164, otp, portal string) error {
	targetE164 := NormalizeE164(phoneE164)
	msgBody := fmt.Sprintf("Your Under100 verification code is %s. Valid for 5 minutes.", otp)
	form := url.Values{}
	form.Set("To", targetE164)
	form.Set("From", t.fromNumber)
	form.Set("Body", msgBody)

	reqURL := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", t.accountSID)
	req, err := http.NewRequestWithContext(ctx, "POST", reqURL, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.SetBasicAuth(t.accountSID, t.authToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := t.client.Do(req)
	if err != nil {
		return fmt.Errorf("Twilio request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Twilio returned HTTP %d: %s", resp.StatusCode, string(body))
	}
	return nil
}

// UnconfiguredProductionProvider is returned when production lacks SMS gateway credentials
type UnconfiguredProductionProvider struct{}

func (u *UnconfiguredProductionProvider) Name() string {
	return "UnconfiguredProductionProvider"
}

func (u *UnconfiguredProductionProvider) IsSimulated() bool {
	return false
}

func (u *UnconfiguredProductionProvider) SendOTP(ctx context.Context, phoneE164, otp, portal string) error {
	return fmt.Errorf("telecom SMS delivery is unconfigured in production: set FAST2SMS_API_KEY, TWILIO_ACCOUNT_SID, or SMS_GATEWAY_URL")
}

// NewProvider resolves the active SMS provider based on configuration and environment
func NewProvider(cfg *config.Config) Provider {
	fast2smsKey := strings.TrimSpace(os.Getenv("FAST2SMS_API_KEY"))
	if fast2smsKey != "" {
		return NewFast2SMSProvider(fast2smsKey)
	}

	twilioSID := strings.TrimSpace(os.Getenv("TWILIO_ACCOUNT_SID"))
	twilioToken := strings.TrimSpace(os.Getenv("TWILIO_AUTH_TOKEN"))
	twilioFrom := strings.TrimSpace(os.Getenv("TWILIO_FROM_NUMBER"))
	if twilioSID != "" && twilioToken != "" && twilioFrom != "" {
		return NewTwilioProvider(twilioSID, twilioToken, twilioFrom)
	}

	// If in production and no gateway is configured, do not pretend simulated OTP works
	if cfg.App.Environment == config.EnvProduction {
		return &UnconfiguredProductionProvider{}
	}

	// Dev / test default: Zero-cost simulated provider
	return NewSimulatedProvider(string(cfg.App.Environment))
}

// NormalizeE164 ensures numbers are formatted as E.164 (defaulting to +91 India)
func NormalizeE164(input string) string {
	clean := strings.TrimSpace(input)
	var digits strings.Builder
	for _, r := range clean {
		if r >= '0' && r <= '9' {
			digits.WriteRune(r)
		}
	}
	num := digits.String()
	if num == "" {
		return ""
	}
	if strings.HasPrefix(num, "91") && len(num) == 12 {
		return "+" + num
	}
	if strings.HasPrefix(num, "0") && len(num) == 11 {
		return "+91" + num[1:]
	}
	if len(num) == 10 {
		return "+91" + num
	}
	if strings.HasPrefix(clean, "+") {
		return "+" + num
	}
	return "+" + num
}

// CleanIndianPhone extracts 10 digits for Indian gateways
func CleanIndianPhone(input string) string {
	e164 := NormalizeE164(input)
	digits := strings.TrimPrefix(e164, "+")
	if strings.HasPrefix(digits, "91") && len(digits) == 12 {
		return digits[2:]
	}
	if len(digits) > 10 {
		return digits[len(digits)-10:]
	}
	return digits
}
