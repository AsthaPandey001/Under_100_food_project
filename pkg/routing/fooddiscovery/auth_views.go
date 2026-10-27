package fooddiscovery

import (
	"fmt"
)

// RenderLoginPage renders the unified authentication login page with Option A (Phone + OTP) and Option B (Email + Password)
func RenderLoginPage(flashMsg, errMsg string, returnTo ...string) string {
	retURL := "/"
	if len(returnTo) > 0 && returnTo[0] != "" {
		retURL = sanitizeReturnURL(returnTo[0])
	}

	var flashHTML string
	if flashMsg == "password_reset_success" {
		flashHTML = `<div class="auth-alert auth-alert-success">
			<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
			<span>Your password has been reset successfully! Please sign in with your new password.</span>
		</div>`
	} else if flashMsg == "verified_success" {
		flashHTML = `<div class="auth-alert auth-alert-success">
			<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
			<span>Your email address has been verified! You can now sign in.</span>
		</div>`
	} else if flashMsg == "auth_required" {
		flashHTML = `<div class="auth-alert auth-alert-warning">
			<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
			<span>Please sign in to access your account profile and orders.</span>
		</div>`
	} else if flashMsg == "admin_required" {
		flashHTML = `<div class="auth-alert auth-alert-warning">
			<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
			<span>Administrator credentials required to access the Admin Panel.</span>
		</div>`
	} else if flashMsg == "merchant_required" {
		flashHTML = `<div class="auth-alert auth-alert-warning">
			<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/><polyline points="9 22 9 12 15 12 15 22"/></svg>
			<span>Restaurant Partner credentials required to access the Merchant Portal.</span>
		</div>`
	} else if flashMsg != "" {
		flashHTML = fmt.Sprintf(`<div class="auth-alert auth-alert-info"><span>%s</span></div>`, escapeHTML(flashMsg))
	}

	var errHTML string
	if errMsg != "" {
		errHTML = fmt.Sprintf(`<div class="auth-alert auth-alert-error">
			<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="15" y1="9" x2="9" y2="15"/><line x1="9" y1="9" x2="15" y2="15"/></svg>
			<span>%s</span>
		</div>`, errMsg) // allows HTML links inside errors (like resend verification)
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Customer Sign In — Food Under ₹100</title>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@400;500;600;700;800&family=Outfit:wght@600;700;800;900&display=swap" rel="stylesheet">
  <style>
    * { box-sizing: border-box; margin: 0; padding: 0; }
    body {
      background: #f8fafc;
      color: #0f172a;
      font-family: 'Plus Jakarta Sans', sans-serif;
      min-height: 100vh;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      padding: 1.5rem;
    }
    .auth-container {
      width: 100%%;
      max-width: 460px;
    }
    .auth-card {
      background: #ffffff;
      border: 1px solid #e2e8f0;
      border-radius: 20px;
      padding: 2.25rem;
      box-shadow: 0 10px 30px rgba(0, 0, 0, 0.05);
    }
    .auth-header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      margin-bottom: 0.5rem;
    }
    .auth-brand {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      font-family: 'Outfit', sans-serif;
      font-size: 1.45rem;
      font-weight: 800;
      color: #0f172a;
      text-decoration: none;
    }
    .auth-brand span { color: #e23744; }
    .auth-city-badge {
      background: #fef08a;
      color: #854d0e;
      font-size: 0.72rem;
      font-weight: 800;
      padding: 0.25rem 0.6rem;
      border-radius: 6px;
      text-transform: uppercase;
      letter-spacing: 0.5px;
    }
    .auth-subtitle {
      color: #64748b;
      font-size: 0.88rem;
      margin-bottom: 1.5rem;
      line-height: 1.4;
    }
    .auth-alert {
      display: flex;
      align-items: flex-start;
      gap: 0.75rem;
      padding: 0.85rem 1rem;
      border-radius: 12px;
      margin-bottom: 1.25rem;
      font-size: 0.86rem;
      line-height: 1.45;
    }
    .auth-alert svg { flex-shrink: 0; margin-top: 1px; }
    .auth-alert-success { background: #dcfce7; border: 1px solid #bbf7d0; color: #166534; }
    .auth-alert-warning { background: #fef9c3; border: 1px solid #fde047; color: #854d0e; }
    .auth-alert-error { background: #fee2e2; border: 1px solid #fecaca; color: #991b1b; }
    .auth-alert-info { background: #e0f2fe; border: 1px solid #bae6fd; color: #0369a1; }
    .auth-alert a { color: inherit; font-weight: 700; text-decoration: underline; }

    /* Method switch buttons (Phone vs Email) */
    .method-tabs {
      display: flex;
      background: #f1f5f9;
      border-radius: 12px;
      padding: 4px;
      margin-bottom: 1.5rem;
      gap: 4px;
    }
    .method-tab {
      flex: 1;
      border: none;
      background: transparent;
      padding: 0.65rem 0.5rem;
      font-family: inherit;
      font-size: 0.84rem;
      font-weight: 700;
      color: #64748b;
      border-radius: 8px;
      cursor: pointer;
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 0.4rem;
      transition: all 0.15s ease;
    }
    .method-tab.active {
      background: #ffffff;
      color: #0f172a;
      box-shadow: 0 2px 6px rgba(0, 0, 0, 0.08);
    }

    .form-group { margin-bottom: 1.15rem; }
    .form-label {
      display: block;
      font-size: 0.82rem;
      font-weight: 700;
      color: #0f172a;
      margin-bottom: 0.4rem;
    }
    .phone-input-row {
      display: flex;
      gap: 0.5rem;
    }
    .country-code-select {
      background: #f8fafc;
      border: 1px solid #cbd5e1;
      border-radius: 10px;
      padding: 0.72rem 0.65rem;
      font-size: 0.9rem;
      font-weight: 600;
      color: #0f172a;
      cursor: pointer;
      outline: none;
    }
    .form-input {
      width: 100%%;
      background: #ffffff;
      border: 1px solid #cbd5e1;
      color: #0f172a;
      padding: 0.72rem 0.9rem;
      border-radius: 10px;
      font-size: 0.92rem;
      font-family: inherit;
      outline: none;
      transition: border-color 0.15s;
    }
    .form-input:focus { border-color: #e23744; }
    .input-password-wrapper {
      position: relative;
      display: flex;
      align-items: center;
    }
    .toggle-pass-btn {
      position: absolute;
      right: 10px;
      background: none;
      border: none;
      color: #64748b;
      cursor: pointer;
      padding: 4px;
      font-size: 0.8rem;
      font-weight: 600;
    }
    .toggle-pass-btn:hover { color: #0f172a; }

    .btn-submit {
      width: 100%%;
      background: #e23744;
      color: #ffffff;
      border: none;
      padding: 0.85rem;
      border-radius: 10px;
      font-size: 0.95rem;
      font-weight: 800;
      font-family: inherit;
      cursor: pointer;
      transition: background 0.15s, transform 0.05s;
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 0.5rem;
    }
    .btn-submit:hover { background: #cb202d; }
    .btn-submit:active { transform: scale(0.99); }
    .btn-submit:disabled { opacity: 0.65; cursor: not-allowed; }

    .btn-otp-action {
      background: #0f172a;
      color: #ffffff;
      border: none;
      padding: 0.72rem 1rem;
      border-radius: 10px;
      font-size: 0.84rem;
      font-weight: 700;
      font-family: inherit;
      cursor: pointer;
      white-space: nowrap;
      transition: background 0.15s;
    }
    .btn-otp-action:hover { background: #334155; }
    .btn-otp-action:disabled { opacity: 0.6; cursor: not-allowed; }

    .otp-code-input {
      letter-spacing: 8px;
      font-size: 1.35rem;
      font-weight: 800;
      text-align: center;
      background: #fafafa;
    }
    .auth-links-row {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 1.25rem;
      font-size: 0.82rem;
    }
    .auth-checkbox-label {
      display: flex;
      align-items: center;
      gap: 0.45rem;
      color: #475569;
      cursor: pointer;
      user-select: none;
    }
    .auth-link {
      color: #e23744;
      font-weight: 700;
      text-decoration: none;
    }
    .auth-link:hover { text-decoration: underline; }

    .auth-footer-nav {
      text-align: center;
      margin-top: 1.5rem;
      font-size: 0.88rem;
      color: #64748b;
    }
    .back-home {
      display: inline-block;
      margin-top: 1.25rem;
      color: #64748b;
      text-decoration: none;
      font-size: 0.84rem;
      text-align: center;
      width: 100%%;
    }
    .back-home:hover { color: #e23744; }

    .portal-quick-links {
      margin-top: 1.25rem;
      padding-top: 1.25rem;
      border-top: 1px dashed #e2e8f0;
      display: flex;
      justify-content: space-between;
      font-size: 0.78rem;
      color: #64748b;
    }
    .portal-quick-links a {
      color: #475569;
      text-decoration: none;
      font-weight: 600;
    }
    .portal-quick-links a:hover { color: #e23744; }
  </style>
</head>
<body>
  <div class="auth-container">
    <div class="auth-card">
      <div class="auth-header">
        <a href="/" class="auth-brand">Under<span>100</span></a>
        <span class="auth-city-badge">Bengaluru</span>
      </div>
      <p class="auth-subtitle">Sign in to your customer account to save favorites, discover dishes under ₹100, and place orders.</p>

      %s
      %s

      <!-- METHOD SELECTION TABS -->
      <div class="method-tabs">
        <button type="button" id="tab-phone" class="method-tab active" onclick="switchLoginMethod('phone')">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"/></svg>
          Phone &amp; OTP
        </button>
        <button type="button" id="tab-email" class="method-tab" onclick="switchLoginMethod('email')">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"/><polyline points="22,6 12,13 2,6"/></svg>
          Email &amp; Password
        </button>
      </div>

      <!-- OPTION A: PHONE NUMBER & OTP FORM -->
      <div id="method-phone-container">
        <div class="form-group">
          <label class="form-label" for="phoneInput">Mobile Phone Number</label>
          <div class="phone-input-row">
            <select id="countryCode" class="country-code-select" aria-label="Country Code">
              <option value="+91" selected>🇮🇳 +91 (IN)</option>
              <option value="+1">🇺🇸 +1 (US)</option>
              <option value="+44">🇬🇧 +44 (UK)</option>
            </select>
            <input type="tel" id="phoneInput" class="form-input" placeholder="98765 43210" maxlength="14" autocomplete="tel-national">
            <button type="button" id="btnSendOtp" class="btn-otp-action" onclick="requestPhoneOTP()">Send OTP</button>
          </div>
        </div>

        <div class="form-group" id="otpBoxGroup" style="display: none;">
          <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.4rem;">
            <label class="form-label" style="margin-bottom: 0;" for="otpInput">Enter 6-Digit OTP</label>
            <span id="otpCountdownText" style="font-size: 0.8rem; color: #64748b; font-weight: 600;">Valid for 5 mins</span>
          </div>
          <input type="text" id="otpInput" class="form-input otp-code-input" placeholder="••••••" maxlength="6" inputmode="numeric" autocomplete="one-time-code">
        </div>

        <div id="phoneOtpAlert" style="display: none; font-size: 0.82rem; margin-bottom: 1rem; font-weight: 600; line-height: 1.4;"></div>

        <button type="button" id="btnVerifyOtp" class="btn-submit" style="display: none;" onclick="verifyPhoneOTP()">
          Verify &amp; Continue
        </button>
      </div>

      <!-- OPTION B: EMAIL & PASSWORD FORM -->
      <div id="method-email-container" style="display: none;">
        <form id="emailLoginForm" method="POST" action="/login">
          <input type="hidden" name="return_to" value="%s">
          <div class="form-group">
            <label class="form-label" for="emailInput">Email Address</label>
            <input type="email" id="emailInput" name="email" class="form-input" placeholder="name@example.com" required autocomplete="email">
          </div>

          <div class="form-group">
            <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.4rem;">
              <label class="form-label" style="margin-bottom: 0;" for="passwordInput">Password</label>
              <a href="/forgot-password" class="auth-link">Forgot Password?</a>
            </div>
            <div class="input-password-wrapper">
              <input type="password" id="passwordInput" name="password" class="form-input" placeholder="••••••••" required autocomplete="current-password">
              <button type="button" class="toggle-pass-btn" onclick="togglePassVisibility('passwordInput', this)">Show</button>
            </div>
          </div>

          <div class="auth-links-row">
            <label class="auth-checkbox-label">
              <input type="checkbox" name="remember" value="true" checked>
              <span>Keep me signed in</span>
            </label>
            <a href="/recover-email" class="auth-link" style="color: #64748b; font-weight: 500;">Forgot Email?</a>
          </div>

          <button type="submit" class="btn-submit">
            Sign In with Email
          </button>
        </form>
      </div>

      <div class="auth-footer-nav">
        Don't have an account? <a href="/register" class="auth-link">Create Account</a>
      </div>

      <div class="portal-quick-links">
        <a href="/merchant/login">Restaurant Owner Portal &rarr;</a>
        <a href="/recover-email">Find My Account &rarr;</a>
      </div>
    </div>

    <a href="/" class="back-home">&larr; Return to Under100 Storefront</a>
  </div>

  <script>
    var resendTimer = null;
    var targetPhone = '';
    var returnURL = %q;

    function switchLoginMethod(method) {
      var tabPhone = document.getElementById('tab-phone');
      var tabEmail = document.getElementById('tab-email');
      var cPhone = document.getElementById('method-phone-container');
      var cEmail = document.getElementById('method-email-container');

      if (method === 'email') {
        tabPhone.classList.remove('active');
        tabEmail.classList.add('active');
        cPhone.style.display = 'none';
        cEmail.style.display = 'block';
      } else {
        tabEmail.classList.remove('active');
        tabPhone.classList.add('active');
        cEmail.style.display = 'none';
        cPhone.style.display = 'block';
      }
    }

    function togglePassVisibility(inputId, btn) {
      var input = document.getElementById(inputId);
      if (input.type === 'password') {
        input.type = 'text';
        btn.innerText = 'Hide';
      } else {
        input.type = 'password';
        btn.innerText = 'Show';
      }
    }

    function requestPhoneOTP() {
      var cc = document.getElementById('countryCode').value;
      var rawPhone = document.getElementById('phoneInput').value.trim();
      var alertBox = document.getElementById('phoneOtpAlert');
      var btnSend = document.getElementById('btnSendOtp');

      if (!rawPhone || rawPhone.replace(/\D/g, '').length < 10) {
        alertBox.style.display = 'block';
        alertBox.style.color = '#dc2626';
        alertBox.innerText = 'Please enter a valid 10-digit mobile number.';
        return;
      }

      var fullPhone = cc + rawPhone.replace(/\D/g, '');
      targetPhone = fullPhone;

      btnSend.disabled = true;
      btnSend.innerText = 'Sending...';
      alertBox.style.display = 'none';

      var params = new URLSearchParams();
      params.append('phone', fullPhone);
      params.append('portal', 'customer');

      fetch('/api/auth/send-otp', {
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
        body: params.toString()
      })
      .then(function(r) { return r.json(); })
      .then(function(data) {
        if (data.success) {
          document.getElementById('otpBoxGroup').style.display = 'block';
          document.getElementById('btnVerifyOtp').style.display = 'flex';
          var otpInp = document.getElementById('otpInput');
          otpInp.value = '';
          otpInp.focus();

          alertBox.style.display = 'block';
          alertBox.style.color = '#166534';
          alertBox.innerHTML = 'Verification OTP sent to <strong>' + (data.target || fullPhone) + '</strong>.';

          startCooldownTimer(60);
        } else {
          btnSend.disabled = false;
          btnSend.innerText = 'Send OTP';
          alertBox.style.display = 'block';
          alertBox.style.color = '#dc2626';
          alertBox.innerText = data.error || 'Failed to send OTP code.';
        }
      })
      .catch(function() {
        btnSend.disabled = false;
        btnSend.innerText = 'Send OTP';
        alertBox.style.display = 'block';
        alertBox.style.color = '#dc2626';
        alertBox.innerText = 'Network error while requesting OTP.';
      });
    }

    function startCooldownTimer(seconds) {
      var btnSend = document.getElementById('btnSendOtp');
      var countdownText = document.getElementById('otpCountdownText');
      if (resendTimer) clearInterval(resendTimer);

      var remaining = seconds;
      btnSend.disabled = true;
      btnSend.innerText = remaining + 's';

      resendTimer = setInterval(function() {
        remaining--;
        if (remaining <= 0) {
          clearInterval(resendTimer);
          btnSend.disabled = false;
          btnSend.innerText = 'Resend OTP';
          countdownText.innerText = 'Did not receive? Click Resend OTP';
        } else {
          btnSend.innerText = remaining + 's';
          countdownText.innerText = 'Resend available in ' + remaining + 's';
        }
      }, 1000);
    }

    function verifyPhoneOTP() {
      var otpCode = document.getElementById('otpInput').value.trim();
      var alertBox = document.getElementById('phoneOtpAlert');
      var btnVerify = document.getElementById('btnVerifyOtp');

      if (!otpCode || otpCode.length !== 6) {
        alertBox.style.display = 'block';
        alertBox.style.color = '#dc2626';
        alertBox.innerText = 'Please enter the complete 6-digit verification code.';
        return;
      }

      btnVerify.disabled = true;
      btnVerify.innerText = 'Verifying...';

      var params = new URLSearchParams();
      params.append('phone', targetPhone);
      params.append('otp', otpCode);
      params.append('portal', 'customer');
      params.append('return_to', returnURL);

      fetch('/api/auth/verify-otp', {
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
        body: params.toString()
      })
      .then(function(r) { return r.json(); })
      .then(function(data) {
        if (data.success) {
          alertBox.style.display = 'block';
          alertBox.style.color = '#166534';
          alertBox.innerText = 'Authentication successful! Redirecting...';
          window.location.href = data.redirect || returnURL || '/';
        } else {
          btnVerify.disabled = false;
          btnVerify.innerText = 'Verify & Continue';
          alertBox.style.display = 'block';
          alertBox.style.color = '#dc2626';
          alertBox.innerText = data.error || 'Invalid verification code.';
        }
      })
      .catch(function() {
        btnVerify.disabled = false;
        btnVerify.innerText = 'Verify & Continue';
        alertBox.style.display = 'block';
        alertBox.style.color = '#dc2626';
        alertBox.innerText = 'Network error during verification.';
      });
    }
  </script>
</body>
</html>`, flashHTML, errHTML, escapeHTML(retURL), retURL)
}

// RenderRegisterPage renders the dedicated user registration page
func RenderRegisterPage(flashMsg, errMsg string, returnTo ...string) string {
	retURL := "/"
	if len(returnTo) > 0 && returnTo[0] != "" {
		retURL = sanitizeReturnURL(returnTo[0])
	}

	var errHTML string
	if errMsg != "" {
		errHTML = fmt.Sprintf(`<div class="auth-alert auth-alert-error">
			<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="15" y1="9" x2="9" y2="15"/><line x1="9" y1="9" x2="15" y2="15"/></svg>
			<span>%s</span>
		</div>`, escapeHTML(errMsg))
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Create Your Account — Food Under ₹100</title>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@400;500;600;700;800&family=Outfit:wght@600;700;800;900&display=swap" rel="stylesheet">
  <style>
    * { box-sizing: border-box; margin: 0; padding: 0; }
    body {
      background: #f8fafc;
      color: #0f172a;
      font-family: 'Plus Jakarta Sans', sans-serif;
      min-height: 100vh;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      padding: 1.5rem;
    }
    .auth-container { width: 100%%; max-width: 480px; }
    .auth-card {
      background: #ffffff;
      border: 1px solid #e2e8f0;
      border-radius: 20px;
      padding: 2.25rem;
      box-shadow: 0 10px 30px rgba(0, 0, 0, 0.05);
    }
    .auth-header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      margin-bottom: 0.5rem;
    }
    .auth-brand {
      font-family: 'Outfit', sans-serif;
      font-size: 1.45rem;
      font-weight: 800;
      color: #0f172a;
      text-decoration: none;
    }
    .auth-brand span { color: #e23744; }
    .auth-subtitle {
      color: #64748b;
      font-size: 0.88rem;
      margin-bottom: 1.5rem;
      line-height: 1.4;
    }
    .auth-alert {
      display: flex;
      align-items: flex-start;
      gap: 0.75rem;
      padding: 0.85rem 1rem;
      border-radius: 12px;
      margin-bottom: 1.25rem;
      font-size: 0.86rem;
    }
    .auth-alert-error { background: #fee2e2; border: 1px solid #fecaca; color: #991b1b; }
    .form-group { margin-bottom: 1.1rem; }
    .form-label {
      display: block;
      font-size: 0.82rem;
      font-weight: 700;
      color: #0f172a;
      margin-bottom: 0.35rem;
    }
    .form-input {
      width: 100%%;
      background: #ffffff;
      border: 1px solid #cbd5e1;
      color: #0f172a;
      padding: 0.72rem 0.9rem;
      border-radius: 10px;
      font-size: 0.92rem;
      font-family: inherit;
      outline: none;
    }
    .form-input:focus { border-color: #e23744; }
    .phone-row { display: flex; gap: 0.5rem; }
    .country-code-select {
      background: #f8fafc;
      border: 1px solid #cbd5e1;
      border-radius: 10px;
      padding: 0.72rem 0.65rem;
      font-size: 0.9rem;
      font-weight: 600;
      color: #0f172a;
    }
    .input-password-wrapper { position: relative; display: flex; align-items: center; }
    .toggle-pass-btn {
      position: absolute;
      right: 10px;
      background: none;
      border: none;
      color: #64748b;
      cursor: pointer;
      font-size: 0.8rem;
      font-weight: 600;
    }
    .btn-submit {
      width: 100%%;
      background: #e23744;
      color: #ffffff;
      border: none;
      padding: 0.85rem;
      border-radius: 10px;
      font-size: 0.95rem;
      font-weight: 800;
      font-family: inherit;
      cursor: pointer;
      margin-top: 0.5rem;
      transition: background 0.15s;
    }
    .btn-submit:hover { background: #cb202d; }
    .auth-checkbox-label {
      display: flex;
      align-items: flex-start;
      gap: 0.55rem;
      font-size: 0.82rem;
      color: #475569;
      cursor: pointer;
      margin-bottom: 1.1rem;
      line-height: 1.4;
    }
    .auth-checkbox-label input { margin-top: 2px; }
    .auth-link { color: #e23744; font-weight: 700; text-decoration: none; }
    .auth-link:hover { text-decoration: underline; }
    .auth-footer-nav {
      text-align: center;
      margin-top: 1.25rem;
      font-size: 0.88rem;
      color: #64748b;
    }
    .back-home {
      display: inline-block;
      margin-top: 1.25rem;
      color: #64748b;
      text-decoration: none;
      font-size: 0.84rem;
      text-align: center;
      width: 100%%;
    }
  </style>
</head>
<body>
  <div class="auth-container">
    <div class="auth-card">
      <div class="auth-header">
        <a href="/" class="auth-brand">Under<span>100</span></a>
        <span style="background: #fef08a; color: #854d0e; font-size: 0.72rem; font-weight: 800; padding: 0.25rem 0.6rem; border-radius: 6px;">BENGALURU</span>
      </div>
      <p class="auth-subtitle">Create your customer account to explore top-rated budget dishes under ₹100 across Bangalore.</p>

      %s

      <form method="POST" action="/register" id="registerForm" onsubmit="return validateRegisterForm()">
        <input type="hidden" name="return_to" value="%s">

        <div class="form-group">
          <label class="form-label" for="regName">Full Name</label>
          <input type="text" id="regName" name="name" class="form-input" placeholder="Aarav Sharma" required autocomplete="name">
        </div>

        <div class="form-group">
          <label class="form-label" for="regEmail">Email Address</label>
          <input type="email" id="regEmail" name="email" class="form-input" placeholder="aarav@example.com" required autocomplete="email">
        </div>

        <div class="form-group">
          <label class="form-label" for="regPhone">Mobile Number (Optional)</label>
          <div class="phone-row">
            <select name="country_code" class="country-code-select">
              <option value="+91" selected>🇮🇳 +91</option>
              <option value="+1">🇺🇸 +1</option>
              <option value="+44">🇬🇧 +44</option>
            </select>
            <input type="tel" id="regPhone" name="phone" class="form-input" placeholder="98765 43210" maxlength="14" autocomplete="tel-national">
          </div>
        </div>

        <div class="form-group">
          <label class="form-label" for="regPass">Create Password</label>
          <div class="input-password-wrapper">
            <input type="password" id="regPass" name="password" class="form-input" placeholder="Min. 8 characters" minlength="8" required autocomplete="new-password">
            <button type="button" class="toggle-pass-btn" onclick="togglePass('regPass', this)">Show</button>
          </div>
        </div>

        <div class="form-group">
          <label class="form-label" for="regConfirmPass">Confirm Password</label>
          <div class="input-password-wrapper">
            <input type="password" id="regConfirmPass" name="confirm_password" class="form-input" placeholder="Re-enter password" minlength="8" required autocomplete="new-password">
            <button type="button" class="toggle-pass-btn" onclick="togglePass('regConfirmPass', this)">Show</button>
          </div>
          <div id="passMismatchText" style="display: none; color: #dc2626; font-size: 0.78rem; margin-top: 4px; font-weight: 600;">Passwords do not match.</div>
        </div>

        <label class="auth-checkbox-label">
          <input type="checkbox" name="terms" required checked>
          <span>I agree to the <a href="/terms" class="auth-link">Terms of Service</a> and <a href="/privacy-policy" class="auth-link">Privacy Policy</a>.</span>
        </label>

        <button type="submit" class="btn-submit">Create Customer Account</button>
      </form>

      <div class="auth-footer-nav">
        Already have an account? <a href="/login" class="auth-link">Sign In</a>
      </div>
    </div>

    <a href="/" class="back-home">&larr; Return to Under100 Storefront</a>
  </div>

  <script>
    function togglePass(id, btn) {
      var inp = document.getElementById(id);
      if (inp.type === 'password') {
        inp.type = 'text';
        btn.innerText = 'Hide';
      } else {
        inp.type = 'password';
        btn.innerText = 'Show';
      }
    }

    function validateRegisterForm() {
      var p1 = document.getElementById('regPass').value;
      var p2 = document.getElementById('regConfirmPass').value;
      var mismatch = document.getElementById('passMismatchText');

      if (p1 !== p2) {
        mismatch.style.display = 'block';
        return false;
      }
      mismatch.style.display = 'none';
      return true;
    }
  </script>
</body>
</html>`, errHTML, escapeHTML(retURL))
}

// RenderVerifyEmailPage renders email verification status, check your inbox card, and resend interface
func RenderVerifyEmailPage(state, email, token, errMsg string) string {
	var bodyContent string

	if state == "success" {
		bodyContent = `<div style="text-align: center; padding: 1rem 0;">
			<div style="width: 68px; height: 68px; background: #dcfce7; color: #16a34a; border-radius: 50%; display: flex; align-items: center; justify-content: center; margin: 0 auto 1.25rem;">
				<svg width="36" height="36" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="20 6 9 17 4 12"/></svg>
			</div>
			<h2 style="font-family: 'Outfit', sans-serif; font-size: 1.5rem; font-weight: 800; margin-bottom: 0.5rem; color: #0f172a;">Email Verified!</h2>
			<p style="color: #64748b; font-size: 0.9rem; line-height: 1.5; margin-bottom: 1.75rem;">Your email address has been verified successfully. Your account is active and ready to use.</p>
			<a href="/login?msg=verified_success" class="btn-submit" style="text-decoration: none; justify-content: center;">Continue to Sign In</a>
		</div>`
	} else if state == "expired" || state == "invalid" {
		title := "Verification Link Expired"
		desc := "This verification link has expired (links are valid for 30 minutes). Please enter your email below to receive a fresh verification link."
		if state == "invalid" {
			title = "Invalid Verification Link"
			desc = "This verification link is invalid or has already been used. Please request a new verification link below."
		}

		bodyContent = fmt.Sprintf(`<div style="text-align: center; padding: 0.5rem 0;">
			<div style="width: 68px; height: 68px; background: #fee2e2; color: #dc2626; border-radius: 50%%; display: flex; align-items: center; justify-content: center; margin: 0 auto 1.25rem;">
				<svg width="34" height="34" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
			</div>
			<h2 style="font-family: 'Outfit', sans-serif; font-size: 1.45rem; font-weight: 800; margin-bottom: 0.5rem; color: #0f172a;">%s</h2>
			<p style="color: #64748b; font-size: 0.88rem; line-height: 1.5; margin-bottom: 1.5rem;">%s</p>

			<div class="form-group" style="text-align: left;">
				<label class="form-label">Your Email Address</label>
				<input type="email" id="resendEmailInput" class="form-input" value="%s" placeholder="name@example.com">
			</div>

			<div id="resendAlert" style="display: none; font-size: 0.84rem; margin-bottom: 1rem; font-weight: 600;"></div>

			<button type="button" id="btnResendToken" class="btn-submit" onclick="triggerResendVerification()">
				Send New Verification Email
			</button>
		</div>`, title, desc, escapeHTML(email))
	} else {
		// Awaiting verification ("Check your inbox")
		bodyContent = fmt.Sprintf(`<div style="text-align: center; padding: 0.5rem 0;">
			<div style="width: 68px; height: 68px; background: #e0f2fe; color: #0284c7; border-radius: 50%%; display: flex; align-items: center; justify-content: center; margin: 0 auto 1.25rem;">
				<svg width="34" height="34" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"/><polyline points="22,6 12,13 2,6"/></svg>
			</div>
			<h2 style="font-family: 'Outfit', sans-serif; font-size: 1.45rem; font-weight: 800; margin-bottom: 0.5rem; color: #0f172a;">Verify Your Email Address</h2>
			<p style="color: #64748b; font-size: 0.9rem; line-height: 1.5; margin-bottom: 1.25rem;">
				We've sent a secure verification link to:<br>
				<strong style="color: #0f172a; word-break: break-all;">%s</strong>
			</p>
			<p style="color: #94a3b8; font-size: 0.8rem; margin-bottom: 1.5rem; line-height: 1.4;">
				Please click the link in your email to activate your account. If you don't see it within a couple minutes, please check your spam folder.
			</p>

			<input type="hidden" id="resendEmailInput" value="%s">
			<div id="resendAlert" style="display: none; font-size: 0.84rem; margin-bottom: 1rem; font-weight: 600;"></div>

			<button type="button" id="btnResendToken" class="btn-submit" style="background: #0f172a;" onclick="triggerResendVerification()">
				Resend Verification Email
			</button>
		</div>`, escapeHTML(email), escapeHTML(email))
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Email Verification — Food Under ₹100</title>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@400;500;600;700;800&family=Outfit:wght@600;700;800;900&display=swap" rel="stylesheet">
  <style>
    * { box-sizing: border-box; margin: 0; padding: 0; }
    body {
      background: #f8fafc;
      color: #0f172a;
      font-family: 'Plus Jakarta Sans', sans-serif;
      min-height: 100vh;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      padding: 1.5rem;
    }
    .auth-card {
      background: #ffffff;
      border: 1px solid #e2e8f0;
      border-radius: 20px;
      width: 460px;
      max-width: 100%%;
      padding: 2.25rem;
      box-shadow: 0 10px 30px rgba(0, 0, 0, 0.05);
    }
    .form-group { margin-bottom: 1rem; }
    .form-label { display: block; font-size: 0.82rem; font-weight: 700; margin-bottom: 0.35rem; }
    .form-input {
      width: 100%%;
      border: 1px solid #cbd5e1;
      padding: 0.72rem 0.9rem;
      border-radius: 10px;
      font-size: 0.92rem;
      outline: none;
    }
    .btn-submit {
      width: 100%%;
      background: #e23744;
      color: #ffffff;
      border: none;
      padding: 0.85rem;
      border-radius: 10px;
      font-size: 0.95rem;
      font-weight: 800;
      cursor: pointer;
      display: flex;
      align-items: center;
      justify-content: center;
      transition: background 0.15s;
    }
    .btn-submit:hover { opacity: 0.92; }
    .btn-submit:disabled { opacity: 0.6; cursor: not-allowed; }
    .back-link {
      display: inline-block;
      margin-top: 1.5rem;
      color: #64748b;
      text-decoration: none;
      font-size: 0.84rem;
      text-align: center;
      width: 100%%;
    }
    .back-link:hover { color: #e23744; }
  </style>
</head>
<body>
  <div class="auth-card">
    <div style="text-align: center; margin-bottom: 1rem;">
      <a href="/" style="font-family: 'Outfit', sans-serif; font-size: 1.45rem; font-weight: 800; color: #0f172a; text-decoration: none;">
        Under<span style="color: #e23744;">100</span>
      </a>
    </div>

    %s

    <a href="/login" class="back-link">&larr; Return to Sign In</a>
  </div>

  <script>
    function triggerResendVerification() {
      var emailVal = document.getElementById('resendEmailInput').value.trim();
      var alertBox = document.getElementById('resendAlert');
      var btn = document.getElementById('btnResendToken');

      if (!emailVal) {
        alertBox.style.display = 'block';
        alertBox.style.color = '#dc2626';
        alertBox.innerText = 'Please enter your registered email address.';
        return;
      }

      btn.disabled = true;
      btn.innerText = 'Sending...';

      var params = new URLSearchParams();
      params.append('email', emailVal);

      fetch('/api/auth/resend-verification', {
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
        body: params.toString()
      })
      .then(function(r) { return r.json(); })
      .then(function(data) {
        if (data.success) {
          alertBox.style.display = 'block';
          alertBox.style.color = '#166534';
          alertBox.innerText = data.message || 'Fresh verification link sent! Please check your inbox.';
          var count = 60;
          btn.innerText = 'Resend in ' + count + 's';
          var timer = setInterval(function() {
            count--;
            if (count <= 0) {
              clearInterval(timer);
              btn.disabled = false;
              btn.innerText = 'Resend Verification Email';
            } else {
              btn.innerText = 'Resend in ' + count + 's';
            }
          }, 1000);
        } else {
          btn.disabled = false;
          btn.innerText = 'Resend Verification Email';
          alertBox.style.display = 'block';
          alertBox.style.color = '#dc2626';
          alertBox.innerText = data.error || 'Failed to resend verification email.';
        }
      })
      .catch(function() {
        btn.disabled = false;
        btn.innerText = 'Resend Verification Email';
        alertBox.style.display = 'block';
        alertBox.style.color = '#dc2626';
        alertBox.innerText = 'Network error while requesting verification link.';
      });
    }
  </script>
</body>
</html>`, bodyContent)
}

// RenderForgotPasswordPage renders password recovery request form with anti-enumeration protection
func RenderForgotPasswordPage(flashMsg, errMsg string) string {
	var alertHTML string
	if flashMsg != "" {
		alertHTML = fmt.Sprintf(`<div style="background: #dcfce7; border: 1px solid #bbf7d0; color: #166534; padding: 0.85rem 1rem; border-radius: 12px; margin-bottom: 1.25rem; font-size: 0.86rem;">
			%s
		</div>`, escapeHTML(flashMsg))
	} else if errMsg != "" {
		alertHTML = fmt.Sprintf(`<div style="background: #fee2e2; border: 1px solid #fecaca; color: #991b1b; padding: 0.85rem 1rem; border-radius: 12px; margin-bottom: 1.25rem; font-size: 0.86rem;">
			%s
		</div>`, escapeHTML(errMsg))
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Reset Password — Food Under ₹100</title>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@400;500;600;700;800&family=Outfit:wght@600;700;800;900&display=swap" rel="stylesheet">
  <style>
    * { box-sizing: border-box; margin: 0; padding: 0; }
    body {
      background: #f8fafc;
      color: #0f172a;
      font-family: 'Plus Jakarta Sans', sans-serif;
      min-height: 100vh;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      padding: 1.5rem;
    }
    .auth-card {
      background: #ffffff;
      border: 1px solid #e2e8f0;
      border-radius: 20px;
      width: 450px;
      max-width: 100%%;
      padding: 2.25rem;
      box-shadow: 0 10px 30px rgba(0, 0, 0, 0.05);
    }
    .form-group { margin-bottom: 1.25rem; }
    .form-label { display: block; font-size: 0.82rem; font-weight: 700; margin-bottom: 0.4rem; }
    .form-input {
      width: 100%%;
      border: 1px solid #cbd5e1;
      padding: 0.72rem 0.9rem;
      border-radius: 10px;
      font-size: 0.92rem;
      outline: none;
    }
    .form-input:focus { border-color: #e23744; }
    .btn-submit {
      width: 100%%;
      background: #e23744;
      color: #ffffff;
      border: none;
      padding: 0.85rem;
      border-radius: 10px;
      font-size: 0.95rem;
      font-weight: 800;
      cursor: pointer;
      transition: background 0.15s;
    }
    .btn-submit:hover { background: #cb202d; }
    .auth-link { color: #e23744; font-weight: 700; text-decoration: none; }
    .auth-link:hover { text-decoration: underline; }
    .back-link {
      display: inline-block;
      margin-top: 1.5rem;
      color: #64748b;
      text-decoration: none;
      font-size: 0.84rem;
      text-align: center;
      width: 100%%;
    }
  </style>
</head>
<body>
  <div class="auth-card">
    <div style="display: flex; align-items: center; justify-content: space-between; margin-bottom: 0.5rem;">
      <a href="/" style="font-family: 'Outfit', sans-serif; font-size: 1.45rem; font-weight: 800; color: #0f172a; text-decoration: none;">
        Under<span style="color: #e23744;">100</span>
      </a>
      <span style="background: #fef08a; color: #854d0e; font-size: 0.72rem; font-weight: 800; padding: 0.25rem 0.6rem; border-radius: 6px;">BENGALURU</span>
    </div>

    <h2 style="font-family: 'Outfit', sans-serif; font-size: 1.3rem; font-weight: 800; margin-bottom: 0.35rem; color: #0f172a;">Forgot Your Password?</h2>
    <p style="color: #64748b; font-size: 0.88rem; line-height: 1.45; margin-bottom: 1.25rem;">
      Enter your registered email address below. If an account exists, we will send you a secure link to reset your password.
    </p>

    %s

    <form method="POST" action="/forgot-password">
      <div class="form-group">
        <label class="form-label" for="forgotEmail">Registered Email Address</label>
        <input type="email" id="forgotEmail" name="email" class="form-input" placeholder="name@example.com" required autocomplete="email">
      </div>
      <button type="submit" class="btn-submit">Send Password Reset Link</button>
    </form>

    <div style="margin-top: 1.25rem; padding-top: 1.25rem; border-top: 1px dashed #e2e8f0; text-align: center; font-size: 0.84rem; color: #64748b;">
      Forgot your registered email? <a href="/recover-email" class="auth-link">Recover with Phone OTP &rarr;</a>
    </div>

    <a href="/login" class="back-link">&larr; Back to Sign In</a>
  </div>
</body>
</html>`, alertHTML)
}

// RenderResetPasswordPage renders the password reset form with new password input
func RenderResetPasswordPage(token, userParam, passwordTokenParam, errMsg string) string {
	var errHTML string
	if errMsg != "" {
		errHTML = fmt.Sprintf(`<div style="background: #fee2e2; border: 1px solid #fecaca; color: #991b1b; padding: 0.85rem 1rem; border-radius: 12px; margin-bottom: 1.25rem; font-size: 0.86rem;">
			%s
		</div>`, escapeHTML(errMsg))
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Set New Password — Food Under ₹100</title>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@400;500;600;700;800&family=Outfit:wght@600;700;800;900&display=swap" rel="stylesheet">
  <style>
    * { box-sizing: border-box; margin: 0; padding: 0; }
    body {
      background: #f8fafc;
      color: #0f172a;
      font-family: 'Plus Jakarta Sans', sans-serif;
      min-height: 100vh;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      padding: 1.5rem;
    }
    .auth-card {
      background: #ffffff;
      border: 1px solid #e2e8f0;
      border-radius: 20px;
      width: 450px;
      max-width: 100%%;
      padding: 2.25rem;
      box-shadow: 0 10px 30px rgba(0, 0, 0, 0.05);
    }
    .form-group { margin-bottom: 1.25rem; }
    .form-label { display: block; font-size: 0.82rem; font-weight: 700; margin-bottom: 0.4rem; }
    .form-input {
      width: 100%%;
      border: 1px solid #cbd5e1;
      padding: 0.72rem 0.9rem;
      border-radius: 10px;
      font-size: 0.92rem;
      outline: none;
    }
    .form-input:focus { border-color: #e23744; }
    .input-password-wrapper { position: relative; display: flex; align-items: center; }
    .toggle-pass-btn {
      position: absolute;
      right: 10px;
      background: none;
      border: none;
      color: #64748b;
      cursor: pointer;
      font-size: 0.8rem;
      font-weight: 600;
    }
    .btn-submit {
      width: 100%%;
      background: #e23744;
      color: #ffffff;
      border: none;
      padding: 0.85rem;
      border-radius: 10px;
      font-size: 0.95rem;
      font-weight: 800;
      cursor: pointer;
      transition: background 0.15s;
    }
    .btn-submit:hover { background: #cb202d; }
    .back-link {
      display: inline-block;
      margin-top: 1.5rem;
      color: #64748b;
      text-decoration: none;
      font-size: 0.84rem;
      text-align: center;
      width: 100%%;
    }
  </style>
</head>
<body>
  <div class="auth-card">
    <div style="text-align: center; margin-bottom: 1rem;">
      <a href="/" style="font-family: 'Outfit', sans-serif; font-size: 1.45rem; font-weight: 800; color: #0f172a; text-decoration: none;">
        Under<span style="color: #e23744;">100</span>
      </a>
    </div>

    <h2 style="font-family: 'Outfit', sans-serif; font-size: 1.35rem; font-weight: 800; margin-bottom: 0.35rem; color: #0f172a;">Set New Password</h2>
    <p style="color: #64748b; font-size: 0.88rem; line-height: 1.45; margin-bottom: 1.25rem;">
      Please create a strong password for your account (minimum 8 characters).
    </p>

    %s

    <form method="POST" action="/reset-password" onsubmit="return validateResetPass()">
      <input type="hidden" name="token" value="%s">
      <input type="hidden" name="user" value="%s">
      <input type="hidden" name="password_token" value="%s">

      <div class="form-group">
        <label class="form-label" for="newPassword">New Password</label>
        <div class="input-password-wrapper">
          <input type="password" id="newPassword" name="password" class="form-input" placeholder="Min. 8 characters" minlength="8" required autocomplete="new-password">
          <button type="button" class="toggle-pass-btn" onclick="togglePass('newPassword', this)">Show</button>
        </div>
      </div>

      <div class="form-group">
        <label class="form-label" for="confirmPassword">Confirm New Password</label>
        <div class="input-password-wrapper">
          <input type="password" id="confirmPassword" name="confirm_password" class="form-input" placeholder="Re-enter new password" minlength="8" required autocomplete="new-password">
          <button type="button" class="toggle-pass-btn" onclick="togglePass('confirmPassword', this)">Show</button>
        </div>
        <div id="mismatchWarning" style="display: none; color: #dc2626; font-size: 0.78rem; margin-top: 4px; font-weight: 600;">Passwords do not match.</div>
      </div>

      <button type="submit" class="btn-submit">Reset Password &amp; Sign In</button>
    </form>

    <a href="/login" class="back-link">&larr; Back to Sign In</a>
  </div>

  <script>
    function togglePass(id, btn) {
      var inp = document.getElementById(id);
      if (inp.type === 'password') {
        inp.type = 'text';
        btn.innerText = 'Hide';
      } else {
        inp.type = 'password';
        btn.innerText = 'Show';
      }
    }

    function validateResetPass() {
      var p1 = document.getElementById('newPassword').value;
      var p2 = document.getElementById('confirmPassword').value;
      var w = document.getElementById('mismatchWarning');
      if (p1 !== p2) {
        w.style.display = 'block';
        return false;
      }
      w.style.display = 'none';
      return true;
    }
  </script>
</body>
</html>`, errHTML, escapeHTML(token), escapeHTML(userParam), escapeHTML(passwordTokenParam))
}

// RenderRecoverEmailPage renders phone-based account recovery to discover forgotten registered email
func RenderRecoverEmailPage(flashMsg, errMsg string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Account Recovery — Food Under ₹100</title>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@400;500;600;700;800&family=Outfit:wght@600;700;800;900&display=swap" rel="stylesheet">
  <style>
    * { box-sizing: border-box; margin: 0; padding: 0; }
    body {
      background: #f8fafc;
      color: #0f172a;
      font-family: 'Plus Jakarta Sans', sans-serif;
      min-height: 100vh;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      padding: 1.5rem;
    }
    .auth-card {
      background: #ffffff;
      border: 1px solid #e2e8f0;
      border-radius: 20px;
      width: 460px;
      max-width: 100%%;
      padding: 2.25rem;
      box-shadow: 0 10px 30px rgba(0, 0, 0, 0.05);
    }
    .form-group { margin-bottom: 1.15rem; }
    .form-label { display: block; font-size: 0.82rem; font-weight: 700; margin-bottom: 0.4rem; }
    .phone-row { display: flex; gap: 0.5rem; }
    .country-code-select {
      background: #f8fafc;
      border: 1px solid #cbd5e1;
      border-radius: 10px;
      padding: 0.72rem 0.65rem;
      font-size: 0.9rem;
      font-weight: 600;
      color: #0f172a;
    }
    .form-input {
      width: 100%%;
      border: 1px solid #cbd5e1;
      padding: 0.72rem 0.9rem;
      border-radius: 10px;
      font-size: 0.92rem;
      outline: none;
    }
    .form-input:focus { border-color: #e23744; }
    .btn-submit {
      width: 100%%;
      background: #e23744;
      color: #ffffff;
      border: none;
      padding: 0.85rem;
      border-radius: 10px;
      font-size: 0.95rem;
      font-weight: 800;
      cursor: pointer;
      transition: background 0.15s;
    }
    .btn-submit:hover { background: #cb202d; }
    .btn-submit:disabled { opacity: 0.6; cursor: not-allowed; }
    .btn-otp-action {
      background: #0f172a;
      color: #ffffff;
      border: none;
      padding: 0.72rem 1rem;
      border-radius: 10px;
      font-size: 0.84rem;
      font-weight: 700;
      cursor: pointer;
      white-space: nowrap;
    }
    .btn-otp-action:disabled { opacity: 0.6; cursor: not-allowed; }
    .otp-code-input {
      letter-spacing: 8px;
      font-size: 1.35rem;
      font-weight: 800;
      text-align: center;
      background: #fafafa;
    }
    .recovery-result-box {
      background: #f0fdf4;
      border: 1px solid #bbf7d0;
      border-radius: 12px;
      padding: 1.25rem;
      margin-top: 1.25rem;
      text-align: center;
    }
    .back-link {
      display: inline-block;
      margin-top: 1.5rem;
      color: #64748b;
      text-decoration: none;
      font-size: 0.84rem;
      text-align: center;
      width: 100%%;
    }
  </style>
</head>
<body>
  <div class="auth-card">
    <div style="display: flex; align-items: center; justify-content: space-between; margin-bottom: 0.5rem;">
      <a href="/" style="font-family: 'Outfit', sans-serif; font-size: 1.45rem; font-weight: 800; color: #0f172a; text-decoration: none;">
        Under<span style="color: #e23744;">100</span>
      </a>
      <span style="background: #fef08a; color: #854d0e; font-size: 0.72rem; font-weight: 800; padding: 0.25rem 0.6rem; border-radius: 6px;">BENGALURU</span>
    </div>

    <h2 style="font-family: 'Outfit', sans-serif; font-size: 1.35rem; font-weight: 800; margin-bottom: 0.35rem; color: #0f172a;">Account Recovery</h2>
    <p style="color: #64748b; font-size: 0.88rem; line-height: 1.45; margin-bottom: 1.25rem;">
      Forgot your registered email address? Verify your mobile phone number with OTP to look up your account.
    </p>

    <div id="recoveryInputSection">
      <div class="form-group">
        <label class="form-label" for="recPhone">Verified Mobile Number</label>
        <div class="phone-row">
          <select id="recCC" class="country-code-select">
            <option value="+91" selected>🇮🇳 +91</option>
            <option value="+1">🇺🇸 +1</option>
            <option value="+44">🇬🇧 +44</option>
          </select>
          <input type="tel" id="recPhone" class="form-input" placeholder="98765 43210" maxlength="14">
          <button type="button" id="btnRecSendOtp" class="btn-otp-action" onclick="sendRecoveryOTP()">Get OTP</button>
        </div>
      </div>

      <div class="form-group" id="recOtpGroup" style="display: none;">
        <label class="form-label" for="recOtp">6-Digit Verification Code</label>
        <input type="text" id="recOtp" class="form-input otp-code-input" placeholder="••••••" maxlength="6" inputmode="numeric">
      </div>

      <div id="recAlert" style="display: none; font-size: 0.84rem; margin-bottom: 1rem; font-weight: 600;"></div>

      <button type="button" id="btnRecVerify" class="btn-submit" style="display: none;" onclick="verifyRecoveryOTP()">
        Verify Phone &amp; Recover Email
      </button>
    </div>

    <!-- RESULT BOX -->
    <div id="recoveryResultBox" class="recovery-result-box" style="display: none;">
      <div style="color: #166534; font-weight: 800; font-size: 1.05rem; margin-bottom: 0.5rem;">Account Found!</div>
      <p style="color: #475569; font-size: 0.85rem; margin-bottom: 0.75rem;">Your registered email address is:</p>
      <div id="maskedEmailDisplay" style="font-size: 1.15rem; font-weight: 800; color: #0f172a; background: #ffffff; border: 1px dashed #86efac; padding: 0.6rem 1rem; border-radius: 8px; margin-bottom: 1.25rem;">
        a***v@e***e.com
      </div>
      <div style="display: flex; gap: 8px; flex-direction: column;">
        <a href="/login" class="btn-submit" style="text-decoration: none; text-align: center; display: block;">Sign In with this Account</a>
        <a id="btnResetForThisEmail" href="/forgot-password" style="color: #e23744; font-weight: 700; font-size: 0.85rem; text-decoration: none; padding: 0.5rem;">Reset Password &rarr;</a>
      </div>
    </div>

    <a href="/login" class="back-link">&larr; Return to Sign In</a>
  </div>

  <script>
    var recPhoneFull = '';

    function sendRecoveryOTP() {
      var cc = document.getElementById('recCC').value;
      var raw = document.getElementById('recPhone').value.trim();
      var alertBox = document.getElementById('recAlert');
      var btnSend = document.getElementById('btnRecSendOtp');

      if (!raw || raw.replace(/\D/g, '').length < 10) {
        alertBox.style.display = 'block';
        alertBox.style.color = '#dc2626';
        alertBox.innerText = 'Please enter a valid 10-digit mobile number.';
        return;
      }

      recPhoneFull = cc + raw.replace(/\D/g, '');
      btnSend.disabled = true;
      btnSend.innerText = 'Sending...';
      alertBox.style.display = 'none';

      var p = new URLSearchParams();
      p.append('phone', recPhoneFull);
      p.append('portal', 'customer');

      fetch('/api/auth/send-otp', {
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
        body: p.toString()
      })
      .then(function(r) { return r.json(); })
      .then(function(data) {
        if (data.success) {
          document.getElementById('recOtpGroup').style.display = 'block';
          document.getElementById('btnRecVerify').style.display = 'block';
          var otpInp = document.getElementById('recOtp');
          otpInp.value = '';
          otpInp.focus();

          alertBox.style.display = 'block';
          alertBox.style.color = '#166534';
          alertBox.innerText = 'Verification OTP sent to ' + (data.target || recPhoneFull);

          var count = 60;
          btnSend.innerText = count + 's';
          var timer = setInterval(function() {
            count--;
            if (count <= 0) {
              clearInterval(timer);
              btnSend.disabled = false;
              btnSend.innerText = 'Resend OTP';
            } else {
              btnSend.innerText = count + 's';
            }
          }, 1000);
        } else {
          btnSend.disabled = false;
          btnSend.innerText = 'Get OTP';
          alertBox.style.display = 'block';
          alertBox.style.color = '#dc2626';
          alertBox.innerText = data.error || 'Failed to send OTP.';
        }
      })
      .catch(function() {
        btnSend.disabled = false;
        btnSend.innerText = 'Get OTP';
        alertBox.style.display = 'block';
        alertBox.style.color = '#dc2626';
        alertBox.innerText = 'Network error while requesting OTP.';
      });
    }

    function verifyRecoveryOTP() {
      var otpVal = document.getElementById('recOtp').value.trim();
      var alertBox = document.getElementById('recAlert');
      var btnVerify = document.getElementById('btnRecVerify');

      if (!otpVal || otpVal.length !== 6) {
        alertBox.style.display = 'block';
        alertBox.style.color = '#dc2626';
        alertBox.innerText = 'Please enter the 6-digit OTP code.';
        return;
      }

      btnVerify.disabled = true;
      btnVerify.innerText = 'Verifying...';

      var p = new URLSearchParams();
      p.append('phone', recPhoneFull);
      p.append('otp', otpVal);

      fetch('/api/auth/recover-email', {
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
        body: p.toString()
      })
      .then(function(r) { return r.json(); })
      .then(function(data) {
        if (data.success) {
          document.getElementById('recoveryInputSection').style.display = 'none';
          document.getElementById('recoveryResultBox').style.display = 'block';
          document.getElementById('maskedEmailDisplay').innerText = data.masked_email;
        } else {
          btnVerify.disabled = false;
          btnVerify.innerText = 'Verify Phone & Recover Email';
          alertBox.style.display = 'block';
          alertBox.style.color = '#dc2626';
          alertBox.innerText = data.error || 'Invalid OTP code.';
        }
      })
      .catch(function() {
        btnVerify.disabled = false;
        btnVerify.innerText = 'Verify Phone & Recover Email';
        alertBox.style.display = 'block';
        alertBox.style.color = '#dc2626';
        alertBox.innerText = 'Network error during verification.';
      });
    }
  </script>
</body>
</html>`)
}
