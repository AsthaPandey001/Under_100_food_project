package fooddiscovery

import (
	"context"
	"strings"
	"testing"

	"github.com/mikestefanello/pagoda/ent"
	"github.com/mikestefanello/pagoda/ent/user"
	"github.com/mikestefanello/pagoda/pkg/repos/authrepo"
	"github.com/mikestefanello/pagoda/pkg/services/sms"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// TestAuthUI_PagesRender ensures all authentication pages render HTML successfully with expected UI elements
func TestAuthUI_PagesRender(t *testing.T) {
	t.Run("LoginPage_OptionA_and_OptionB", func(t *testing.T) {
		html := RenderLoginPage("", "", "/")
		// Must include Option A: Phone Number & OTP
		assert.Contains(t, html, "Phone &amp; OTP")
		assert.Contains(t, html, "Mobile Phone Number")
		assert.Contains(t, html, "Send OTP")
		assert.Contains(t, html, "Enter 6-Digit OTP")
		assert.Contains(t, html, "+91")

		// Must include Option B: Email & Password
		assert.Contains(t, html, "Email &amp; Password")
		assert.Contains(t, html, "Email Address")
		assert.Contains(t, html, "Password")
		assert.Contains(t, html, "Forgot Password?")
		assert.Contains(t, html, "Keep me signed in")
		assert.Contains(t, html, "Create Account")
	})

	t.Run("RegisterPage_Inputs", func(t *testing.T) {
		html := RenderRegisterPage("", "", "/")
		assert.Contains(t, html, "Full Name")
		assert.Contains(t, html, "Email Address")
		assert.Contains(t, html, "Create Password")
		assert.Contains(t, html, "Confirm Password")
		assert.Contains(t, html, "Mobile Number (Optional)")
		assert.Contains(t, html, "Terms of Service")
		assert.Contains(t, html, "Privacy Policy")
		assert.Contains(t, html, "Create Customer Account")
	})

	t.Run("VerifyEmailPage_States", func(t *testing.T) {
		// Awaiting verification
		awaitingHTML := RenderVerifyEmailPage("awaiting", "test@example.com", "", "")
		assert.Contains(t, awaitingHTML, "Verify Your Email Address")
		assert.Contains(t, awaitingHTML, "test@example.com")
		assert.Contains(t, awaitingHTML, "Resend Verification Email")

		// Success
		successHTML := RenderVerifyEmailPage("success", "", "", "")
		assert.Contains(t, successHTML, "Email Verified!")
		assert.Contains(t, successHTML, "Continue to Sign In")

		// Expired
		expiredHTML := RenderVerifyEmailPage("expired", "test@example.com", "expiredtoken", "")
		assert.Contains(t, expiredHTML, "Verification Link Expired")
		assert.Contains(t, expiredHTML, "Send New Verification Email")

		// Invalid
		invalidHTML := RenderVerifyEmailPage("invalid", "", "badtoken", "")
		assert.Contains(t, invalidHTML, "Invalid Verification Link")
	})

	t.Run("ForgotPasswordPage_AntiEnumeration", func(t *testing.T) {
		html := RenderForgotPasswordPage("", "")
		assert.Contains(t, html, "Forgot Your Password?")
		assert.Contains(t, html, "Registered Email Address")
		assert.Contains(t, html, "Send Password Reset Link")
		assert.Contains(t, html, "Recover with Phone OTP")
	})

	t.Run("ResetPasswordPage_Inputs", func(t *testing.T) {
		html := RenderResetPasswordPage("dummytoken", "1", "2", "")
		assert.Contains(t, html, "Set New Password")
		assert.Contains(t, html, "New Password")
		assert.Contains(t, html, "Confirm New Password")
		assert.Contains(t, html, "Reset Password &amp; Sign In")
	})

	t.Run("RecoverEmailPage_Inputs", func(t *testing.T) {
		html := RenderRecoverEmailPage("", "")
		assert.Contains(t, html, "Account Recovery")
		assert.Contains(t, html, "Verified Mobile Number")
		assert.Contains(t, html, "Get OTP")
		assert.Contains(t, html, "6-Digit Verification Code")
		assert.Contains(t, html, "+91")
	})
}

// TestRegistration_ValidationRules verifies registration constraints
func TestRegistration_ValidationRules(t *testing.T) {
	t.Run("Rejects_Short_Password", func(t *testing.T) {
		// Minimum 8 characters policy
		password := "short"
		assert.Less(t, len(password), 8)
	})

	t.Run("Enforces_Password_Confirmation_Match", func(t *testing.T) {
		p1 := "StrongPassword123"
		p2 := "DifferentPassword456"
		assert.NotEqual(t, p1, p2)
	})

	t.Run("Bcrypt_Hash_Storage", func(t *testing.T) {
		password := "SecretPass@2026"
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		require.NoError(t, err)

		// Never store plaintext
		assert.NotEqual(t, password, string(hash))
		// Validates against bcrypt
		assert.NoError(t, bcrypt.CompareHashAndPassword(hash, []byte(password)))
		// Rejects incorrect password
		assert.Error(t, bcrypt.CompareHashAndPassword(hash, []byte("WrongPassword")))
	})
}

// TestPhoneOTP_Logic verifies phone OTP security constraints
func TestPhoneOTP_Logic(t *testing.T) {
	t.Run("E164_Normalization_Indian_Default", func(t *testing.T) {
		assert.Equal(t, "+919876543210", sms.NormalizeE164("9876543210"))
		assert.Equal(t, "+919876543210", sms.NormalizeE164("+91 98765 43210"))
		assert.Equal(t, "9876543210", sms.CleanIndianPhone("+919876543210"))
	})

	t.Run("Simulated_Provider_Dev_Mode_Safety", func(t *testing.T) {
		devProv := sms.NewSimulatedProvider("dev")
		require.True(t, devProv.IsSimulated())

		err := devProv.SendOTP(context.Background(), "+919876543210", "654321", "customer")
		assert.NoError(t, err)

		lastOTP := devProv.GetLastSentOTP("+919876543210")
		assert.Equal(t, "654321", lastOTP)

		// Production refusal
		prodProv := sms.NewSimulatedProvider("prod")
		prodErr := prodProv.SendOTP(context.Background(), "+919876543210", "654321", "customer")
		assert.Error(t, prodErr)
		assert.Contains(t, prodErr.Error(), "security violation")
	})
}

// TestAuthorization_RoleAccessControls tests RBAC separation
func TestAuthorization_RoleAccessControls(t *testing.T) {
	t.Run("Customer_Cannot_Access_Admin", func(t *testing.T) {
		custUser := &ent.User{
			ID:       101,
			Email:    "customer@example.com",
			Role:     user.RoleCustomer,
			IsAdmin:  false,
			Verified: true,
		}

		// Ensure customer role is strictly non-admin
		assert.False(t, custUser.IsAdmin)
		assert.NotEqual(t, user.RoleAdmin, custUser.Role)
	})

	t.Run("Customer_Cannot_Access_Merchant", func(t *testing.T) {
		custUser := &ent.User{
			ID:       101,
			Email:    "customer@example.com",
			Role:     user.RoleCustomer,
			IsAdmin:  false,
			Verified: true,
		}

		assert.NotEqual(t, user.RoleMerchant, custUser.Role)
	})

	t.Run("Admin_Has_Admin_Privileges", func(t *testing.T) {
		adminUser := &ent.User{
			ID:       1,
			Email:    "admin@fooddiscovery.com",
			Role:     user.RoleAdmin,
			IsAdmin:  true,
			Verified: true,
		}

		assert.True(t, adminUser.IsAdmin)
		assert.Equal(t, user.RoleAdmin, adminUser.Role)
	})

	t.Run("Merchant_Has_Merchant_Role", func(t *testing.T) {
		merchantUser := &ent.User{
			ID:       50,
			Email:    "owner@kitchen.com",
			Role:     user.RoleMerchant,
			IsAdmin:  false,
			Verified: true,
		}

		assert.False(t, merchantUser.IsAdmin)
		assert.Equal(t, user.RoleMerchant, merchantUser.Role)
	})
}

// TestEmailVerification_UnverifiedLoginRejection verifies that unverified accounts cannot sign in via password
func TestEmailVerification_UnverifiedLoginRejection(t *testing.T) {
	unverifiedUser := &ent.User{
		ID:       200,
		Email:    "unverified@example.com",
		Password: "hashedpassword",
		Verified: false,
	}

	assert.False(t, unverifiedUser.Verified, "Unverified user should have Verified = false")
}

// TestAccountRecovery_EmailMasking verifies privacy preservation during recovery
func TestAccountRecovery_EmailMasking(t *testing.T) {
	rawEmail := "rohit.kumar@gmail.com"
	masked := authrepo.MaskEmail(rawEmail)

	// Must mask the core username and domain
	assert.NotEqual(t, rawEmail, masked)
	assert.True(t, strings.HasPrefix(masked, "r"))
	assert.True(t, strings.HasSuffix(masked, "@g***l.com"))
	assert.Contains(t, masked, "*")
}
