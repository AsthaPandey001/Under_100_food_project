package authrepo

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashToken(t *testing.T) {
	token := "4a8e2b9c7d1f0e3a5b6c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f90"
	expectedSum := sha256.Sum256([]byte(token))
	expected := hex.EncodeToString(expectedSum[:])

	actual := HashToken(token)
	assert.Equal(t, expected, actual)

	// Ensure hashing is deterministic
	assert.Equal(t, actual, HashToken(token))

	// Ensure different token yields different hash
	assert.NotEqual(t, actual, HashToken("different_token"))
}

func TestMaskEmail(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Standard email", "aarav@example.com", "a***v@e*****e.com"},
		{"Long username", "bengaluru_foodie@gmail.com", "b**************e@g***l.com"},
		{"Short 2 char name", "me@domain.org", "m*@d****n.org"},
		{"Single char name", "a@test.in", "a*@t**t.in"},
		{"Subdomain email", "support@kitchen.under100.in", "s*****t@k*****n.under100.in"},
		{"Invalid email without @", "invalidemail", "***@***.***"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := MaskEmail(tt.input)
			assert.Equal(t, tt.expected, actual)
			// Crucial security invariant: never return plaintext input
			if len(tt.input) > 4 && tt.input != "***@***.***" {
				assert.NotEqual(t, tt.input, actual)
			}
		})
	}
}

func TestAuthErrors(t *testing.T) {
	// Verify error messages match expected security descriptions
	assert.Equal(t, "email verification link has expired", ErrTokenExpired.Error())
	assert.Equal(t, "email verification link is invalid", ErrTokenInvalid.Error())
	assert.Equal(t, "email verification link has already been used", ErrTokenAlreadyUsed.Error())
	assert.Equal(t, "please wait 60 seconds before requesting another verification email", ErrResendCooldown.Error())
	assert.Equal(t, "maximum verification attempts exceeded. Please request a new OTP", ErrOTPMaxAttemptsExceeded.Error())
	assert.Equal(t, "please wait 60 seconds before requesting another OTP", ErrOTPCooldown.Error())
}
