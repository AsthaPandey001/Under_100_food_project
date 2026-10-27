package emailservice

import (
	"fmt"
	"html"
	"strings"
)

// BaseEmailTemplate wraps email body in responsive Food Under ₹100 branded HTML container
func wrapBaseHTML(title, preheader, contentHTML string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <meta name="x-apple-disable-message-reformatting">
  <title>%s</title>
  <style>
    body {
      margin: 0;
      padding: 0;
      background-color: #f3f4f6;
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
      color: #1f2937;
      -webkit-font-smoothing: antialiased;
    }
    .wrapper {
      width: 100%%;
      table-layout: fixed;
      background-color: #f3f4f6;
      padding: 40px 0;
    }
    .main-card {
      background-color: #ffffff;
      margin: 0 auto;
      max-width: 600px;
      border-radius: 12px;
      overflow: hidden;
      box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05), 0 2px 4px -1px rgba(0, 0, 0, 0.03);
      border: 1px solid #e5e7eb;
    }
    .header {
      background: linear-gradient(135deg, #e23744 0%%, #ff5200 100%%);
      padding: 32px 24px;
      text-align: center;
    }
    .header h1 {
      margin: 0;
      color: #ffffff;
      font-size: 26px;
      font-weight: 800;
      letter-spacing: -0.5px;
    }
    .header p {
      margin: 6px 0 0;
      color: #fee2e2;
      font-size: 13px;
      font-weight: 500;
      text-transform: uppercase;
      letter-spacing: 1px;
    }
    .content {
      padding: 36px 32px;
      line-height: 1.6;
      font-size: 15px;
      color: #374151;
    }
    .content h2 {
      font-size: 20px;
      font-weight: 700;
      color: #111827;
      margin-top: 0;
      margin-bottom: 16px;
    }
    .btn-container {
      margin: 28px 0;
      text-align: center;
    }
    .btn {
      display: inline-block;
      background: linear-gradient(135deg, #e23744 0%%, #ff5200 100%%);
      color: #ffffff !important;
      font-weight: 700;
      font-size: 15px;
      padding: 14px 28px;
      text-decoration: none;
      border-radius: 8px;
      box-shadow: 0 2px 4px rgba(226, 55, 68, 0.25);
    }
    .raw-url {
      background-color: #f9fafb;
      border: 1px solid #e5e7eb;
      border-radius: 6px;
      padding: 12px;
      word-break: break-all;
      font-family: monospace;
      font-size: 12px;
      color: #4b5563;
      margin: 16px 0;
    }
    .footer {
      background-color: #f9fafb;
      padding: 24px 32px;
      text-align: center;
      border-top: 1px solid #e5e7eb;
      font-size: 12px;
      color: #6b7280;
      line-height: 1.5;
    }
    .footer a {
      color: #e23744;
      text-decoration: none;
    }
    @media only screen and (max-width: 600px) {
      .main-card {
        border-radius: 0 !important;
        border: none !important;
      }
      .content {
        padding: 24px 18px !important;
      }
    }
  </style>
</head>
<body>
  <!-- Preheader hidden text for email client inbox snippet -->
  <div style="display:none;font-size:1px;color:#333333;line-height:1px;max-height:0px;max-width:0px;opacity:0;overflow:hidden;">
    %s
  </div>
  <table class="wrapper" role="presentation" border="0" cellpadding="0" cellspacing="0">
    <tr>
      <td align="center">
        <table class="main-card" role="presentation" border="0" cellpadding="0" cellspacing="0" width="100%%">
          <!-- Header -->
          <tr>
            <td class="header">
              <h1>Food Under ₹100</h1>
              <p>Bengaluru's Best Eats Under ₹100</p>
            </td>
          </tr>
          <!-- Body Content -->
          <tr>
            <td class="content">
              %s
            </td>
          </tr>
          <!-- Footer -->
          <tr>
            <td class="footer">
              <p style="margin: 0 0 8px 0; font-weight: 600; color: #374151;">Food Under ₹100 — Discover Bengaluru's budget culinary gems.</p>
              <p style="margin: 0 0 8px 0;">This is an automated transactional notification. If you did not request this, no action is required.</p>
              <p style="margin: 0;">&copy; 2026 Food Under ₹100 &bull; Bengaluru, Karnataka, India</p>
            </td>
          </tr>
        </table>
      </td>
    </tr>
  </table>
</body>
</html>`, title, preheader, contentHTML)
}

// BuildWelcomeEmail renders registration confirmation / welcome email
func BuildWelcomeEmail(name, appBaseURL string) (htmlContent, textContent, subject string) {
	safeName := html.EscapeString(strings.TrimSpace(name))
	if safeName == "" {
		safeName = "Foodie"
	}

	subject = "Welcome to Food Under ₹100!"
	preheader := "Your registration has been received successfully. Start discovering delicious dishes under ₹100 in Bengaluru!"

	contentHTML := fmt.Sprintf(`
<h2>Welcome to Food Under ₹100!</h2>
<p>Hello <strong>%s</strong>,</p>
<p>Your registration request has been received successfully.</p>
<p>We're thrilled to have you join our community! Food Under ₹100 was built specifically to help you discover Bengaluru's legendary heritage darshinis, iconic breakfast outlets, and hidden street delicacies—all priced strictly at or below <strong>₹100</strong>.</p>
<div class="btn-container">
  <a href="%s" class="btn" target="_blank">Explore Dishes Under ₹100 &rarr;</a>
</div>
<p style="font-size: 13px; color: #6b7280; margin-top: 24px; padding-top: 16px; border-top: 1px solid #f3f4f6;">
  <strong>Notice:</strong> This email confirms your registration request only. If you signed up with email verification required, please also click the verification link sent separately to activate your account.
</p>
<p>Happy eating,<br><strong>The Food Under ₹100 Team</strong></p>
`, safeName, html.EscapeString(appBaseURL))

	htmlContent = wrapBaseHTML(subject, preheader, contentHTML)

	textContent = fmt.Sprintf(`Welcome to Food Under ₹100!

Hello %s,

Welcome to Food Under ₹100!

Your registration request has been received successfully.

We're excited to help you discover delicious and affordable food options across Bangalore.

Start exploring dishes under ₹100: %s

Thank you for joining us!

Regards,
Food Under ₹100 Team
`, name, appBaseURL)

	return htmlContent, textContent, subject
}

// BuildVerificationEmail renders the email ownership verification link email
func BuildVerificationEmail(name, verifyURL string) (htmlContent, textContent, subject string) {
	safeName := html.EscapeString(strings.TrimSpace(name))
	if safeName == "" {
		safeName = "Valued Customer"
	}

	subject = "Verify Your Email Address — Food Under ₹100"
	preheader := "Please click the button to verify your email address and activate your account."

	contentHTML := fmt.Sprintf(`
<h2>Verify Your Email Address</h2>
<p>Hello <strong>%s</strong>,</p>
<p>Thank you for signing up with Food Under ₹100. To activate your account and start discovering Bengaluru's best budget food, please verify your email address by clicking the button below:</p>
<div class="btn-container">
  <a href="%s" class="btn" target="_blank">Verify Email Address &rarr;</a>
</div>
<p>If the button doesn't work, copy and paste the following link into your browser:</p>
<div class="raw-url">%s</div>
<p style="font-size: 13px; color: #6b7280; margin-top: 20px;">
  &bull; This verification link will expire in <strong>30 minutes</strong>.<br>
  &bull; This link is single-use and will be invalidated once verified.<br>
  &bull; If you did not create an account with Food Under ₹100, you can safely ignore this email.
</p>
<p>Regards,<br><strong>The Food Under ₹100 Team</strong></p>
`, safeName, html.EscapeString(verifyURL), html.EscapeString(verifyURL))

	htmlContent = wrapBaseHTML(subject, preheader, contentHTML)

	textContent = fmt.Sprintf(`Verify Your Email Address — Food Under ₹100

Hello %s,

Thank you for signing up with Food Under ₹100!

Please verify your email address to activate your account by clicking the link below:
%s

This verification link will expire in 30 minutes and can only be used once.

If you did not create an account, you can safely ignore this message.

Regards,
Food Under ₹100 Team
`, name, verifyURL)

	return htmlContent, textContent, subject
}

// BuildPasswordResetEmail renders the secure password reset instructions email
func BuildPasswordResetEmail(name, resetURL string) (htmlContent, textContent, subject string) {
	safeName := html.EscapeString(strings.TrimSpace(name))
	if safeName == "" {
		safeName = "Customer"
	}

	subject = "Reset Your Password — Food Under ₹100"
	preheader := "Instructions to reset your Food Under ₹100 account password."

	contentHTML := fmt.Sprintf(`
<h2>Password Reset Request</h2>
<p>Hello <strong>%s</strong>,</p>
<p>We received a request to reset the password for your Food Under ₹100 account. Click the button below to choose a new, secure password:</p>
<div class="btn-container">
  <a href="%s" class="btn" target="_blank">Reset Password &rarr;</a>
</div>
<p>If you prefer, copy and paste the link below directly into your web browser:</p>
<div class="raw-url">%s</div>
<p style="font-size: 13px; color: #6b7280; margin-top: 20px;">
  &bull; This password reset link is valid for <strong>30 minutes</strong> only.<br>
  &bull; This link can only be used once.<br>
  &bull; If you did not request a password reset, your account is safe and you can ignore this email.
</p>
<p>Regards,<br><strong>The Food Under ₹100 Security Team</strong></p>
`, safeName, html.EscapeString(resetURL), html.EscapeString(resetURL))

	htmlContent = wrapBaseHTML(subject, preheader, contentHTML)

	textContent = fmt.Sprintf(`Reset Your Password — Food Under ₹100

Hello %s,

We received a request to reset the password for your Food Under ₹100 account.

Click the link below to set a new password:
%s

This link is valid for 30 minutes and can only be used once.

If you did not request a password reset, no further action is needed and your existing password remains unchanged.

Regards,
Food Under ₹100 Security Team
`, name, resetURL)

	return htmlContent, textContent, subject
}

// BuildPasswordChangedNotification renders account security alert when password is changed
func BuildPasswordChangedNotification(name, appBaseURL string) (htmlContent, textContent, subject string) {
	safeName := html.EscapeString(strings.TrimSpace(name))
	if safeName == "" {
		safeName = "Customer"
	}

	subject = "Security Alert: Password Changed — Food Under ₹100"
	preheader := "Your Food Under ₹100 account password was recently changed."

	contentHTML := fmt.Sprintf(`
<h2>Your Password Has Been Changed</h2>
<p>Hello <strong>%s</strong>,</p>
<p>This is a confirmation that the password for your Food Under ₹100 account was recently changed.</p>
<p>All previously active sessions have been signed out for your security.</p>
<div class="btn-container">
  <a href="%s/login" class="btn" target="_blank">Sign In with New Password &rarr;</a>
</div>
<div style="background-color: #fee2e2; border-left: 4px solid #e23744; padding: 12px 16px; border-radius: 4px; margin: 24px 0;">
  <p style="margin: 0; font-size: 13px; color: #991b1b;">
    <strong>Did you not make this change?</strong><br>
    If you did not change your password, someone else may have gained access to your account. Please use the <a href="%s/forgot-password" style="color: #991b1b; text-decoration: underline; font-weight: bold;">Forgot Password</a> page immediately to secure your account.
  </p>
</div>
<p>Regards,<br><strong>The Food Under ₹100 Security Team</strong></p>
`, safeName, html.EscapeString(appBaseURL), html.EscapeString(appBaseURL))

	htmlContent = wrapBaseHTML(subject, preheader, contentHTML)

	textContent = fmt.Sprintf(`Security Alert: Password Changed — Food Under ₹100

Hello %s,

This is a confirmation that the password for your Food Under ₹100 account was recently changed.

All existing sessions on other devices have been invalidated.

If you made this change, you can safely disregard this email.

If you did NOT make this change, please reset your password immediately at:
%s/forgot-password

Regards,
Food Under ₹100 Security Team
`, name, appBaseURL)

	return htmlContent, textContent, subject
}
