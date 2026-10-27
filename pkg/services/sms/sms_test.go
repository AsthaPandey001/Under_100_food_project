package sms

import (
	"context"
	"strings"
	"sync"
	"testing"

	"github.com/mikestefanello/pagoda/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeE164(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"10 digits plain Indian", "9876543210", "+919876543210"},
		{"10 digits formatted with spaces and dashes", "98765-43210", "+919876543210"},
		{"Leading zero 11 digits", "09876543210", "+919876543210"},
		{"Prefix 91 without plus 12 digits", "919876543210", "+919876543210"},
		{"Already valid E.164 India", "+919876543210", "+919876543210"},
		{"US number with plus", "+14155552671", "+14155552671"},
		{"UK number with plus", "+447911123456", "+447911123456"},
		{"Empty string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := NormalizeE164(tt.input)
			assert.Equal(t, tt.expected, res)
		})
	}
}

func TestCleanIndianPhone(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"+919876543210", "9876543210"},
		{"919876543210", "9876543210"},
		{"09876543210", "9876543210"},
		{"9876543210", "9876543210"},
		{"98765 43210", "9876543210"},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, CleanIndianPhone(tt.input))
	}
}

func TestSimulatedProvider_DevMode(t *testing.T) {
	provider := NewSimulatedProvider(string(config.EnvDevelop))
	require.True(t, provider.IsSimulated())
	assert.Equal(t, "SimulatedDevOTP", provider.Name())

	ctx := context.Background()
	phone := "+919876543210"
	otp := "741258"

	err := provider.SendOTP(ctx, phone, otp, "customer")
	require.NoError(t, err)

	lastOTP := provider.GetLastSentOTP(phone)
	assert.Equal(t, otp, lastOTP)
}

func TestSimulatedProvider_ProductionRefusal(t *testing.T) {
	// Strict requirement: Fail immediately / refuse simulated OTP mode in production
	prodProvider := NewSimulatedProvider(string(config.EnvProduction))
	ctx := context.Background()

	err := prodProvider.SendOTP(ctx, "+919876543210", "123456", "customer")
	require.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "production"), "expected production error message")
}

func TestSimulatedProvider_ThreadSafety(t *testing.T) {
	provider := NewSimulatedProvider("dev")
	ctx := context.Background()

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			phone := "+919876543210"
			_ = provider.SendOTP(ctx, phone, "112233", "customer")
			_ = provider.GetLastSentOTP(phone)
		}(i)
	}
	wg.Wait()
}

func TestUnconfiguredProductionProvider(t *testing.T) {
	u := &UnconfiguredProductionProvider{}
	assert.False(t, u.IsSimulated())
	assert.Equal(t, "UnconfiguredProductionProvider", u.Name())

	err := u.SendOTP(context.Background(), "+919876543210", "123456", "customer")
	require.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "unconfigured in production"))
}
