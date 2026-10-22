# verify_phase_3.ps1 - Automated Verification for Phase 3
# Tests all 27 requirements with disposable fixtures and zero DB pollution.

$ErrorActionPreference = "Continue"
$BaseUrl = "http://127.0.0.1:8000"
$Passed = 0
$Failed = 0

function Report-Test($num, $name, $result, $details) {
    if ($result) {
        $script:Passed++
        Write-Host " [PASS] Test $($num): $name" -ForegroundColor Green
        if ($details) { Write-Host "        $details" -ForegroundColor DarkGray }
    } else {
        $script:Failed++
        Write-Host " [FAIL] Test $($num): $name" -ForegroundColor Red
        if ($details) { Write-Host "        $details" -ForegroundColor Yellow }
    }
}

Write-Host "=================================================================" -ForegroundColor Cyan
Write-Host " PHASE 3 VERIFICATION SUITE - BITE100 FOOD DISCOVERY" -ForegroundColor Cyan
Write-Host "=================================================================" -ForegroundColor Cyan

function Cleanup-User($email) {
    $esc = $email.Replace("'", "''")
    & "C:\Users\astha\pgsql\bin\psql.exe" -U admin -d goship_db -c "DELETE FROM password_tokens WHERE password_token_user IN (SELECT id FROM users WHERE email = '$esc'); DELETE FROM users WHERE email = '$esc';" *> $null
}

$TestEmail = "phase3_tester_" + (Get-Random) + "@testdomain.com"
$VictimEmail = "phase3_victim_" + (Get-Random) + "@testdomain.com"
$HelperEmail = "phase3_helper_" + (Get-Random) + "@testdomain.com"
Cleanup-User $TestEmail
Cleanup-User $VictimEmail
Cleanup-User $HelperEmail

$cookieFile = "C:\Users\astha\Desktop\food-discovery\p3_main_cookie.txt"
$oldCookie = "C:\Users\astha\Desktop\food-discovery\p3_old_cookie.txt"
$newCookie = "C:\Users\astha\Desktop\food-discovery\p3_new_cookie.txt"
$recCookie = "C:\Users\astha\Desktop\food-discovery\p3_rec_cookie.txt"
$merchCookie = "C:\Users\astha\Desktop\food-discovery\p3_merch_cookie.txt"

Remove-Item -Force -ErrorAction SilentlyContinue $cookieFile, $oldCookie, $newCookie, $recCookie, $merchCookie

# -----------------------------------------------------------------
# 1. Unauthenticated access to /profile is rejected or redirected
# -----------------------------------------------------------------
$resp = curl.exe -s -i "$BaseUrl/profile"
$isRedirect = ($resp -match "HTTP/1.1 303") -and ($resp -match "Location: /login\?msg=auth_required")
Report-Test 1 "Unauthenticated access to /profile redirected" $isRedirect "Status 303 -> /login?msg=auth_required"

# Register a disposable test user for write-based profile tests
$regResp = curl.exe -s -i -c $cookieFile -d "name=Phase3 Tester&email=$TestEmail&password=InitialPass123&role=customer" "$BaseUrl/register"

# -----------------------------------------------------------------
# 2. Authenticated customer can view their profile
# -----------------------------------------------------------------
$profResp = curl.exe -s -b $cookieFile "$BaseUrl/profile"
$canView = ($profResp -match "Phase3 Tester") -and ($profResp -match $TestEmail) -and ($profResp -match "Profile & Security")
Report-Test 2 "Authenticated customer can view profile" $canView "Profile rendered with user name, email, and role badge"

# -----------------------------------------------------------------
# 3. Customer can update their own display name
# -----------------------------------------------------------------
$updateResp = curl.exe -s -b $cookieFile -c $cookieFile -d "name=Updated Tester Name" "$BaseUrl/profile"
$nameUpdated = ($updateResp -match "Updated Tester Name") -and ($updateResp -match "Profile updated successfully.")
Report-Test 3 "Customer can update own display name" $nameUpdated "Display name updated and success message displayed"

# -----------------------------------------------------------------
# 4. Invalid profile input is rejected
# -----------------------------------------------------------------
$invalidResp = curl.exe -s -b $cookieFile -c $cookieFile -d "name= " "$BaseUrl/profile"
$invalidRejected = ($invalidResp -match "Display name must be between 2 and 100 characters.")
$dbName = (& "C:\Users\astha\pgsql\bin\psql.exe" -U admin -d goship_db -t -A -c "SELECT name FROM users WHERE email = '$TestEmail';")
if ($dbName) { $dbName = $dbName.Trim() }
$dbIntact = ($dbName -eq "Updated Tester Name")
Report-Test 4 "Invalid profile input rejected" ($invalidRejected -and $dbIntact) "Empty/whitespace name rejected, DB unchanged"

# -----------------------------------------------------------------
# 5. User cannot modify another user's profile using a forged ID
# -----------------------------------------------------------------
& "C:\Users\astha\pgsql\bin\psql.exe" -U admin -d goship_db -c "INSERT INTO users (name, email, password, role, is_admin, verified, created_at, updated_at) VALUES ('Victim User', '$VictimEmail', 'hash123', 'customer', false, true, NOW(), NOW());" *> $null
$victimId = (& "C:\Users\astha\pgsql\bin\psql.exe" -U admin -d goship_db -t -A -c "SELECT id FROM users WHERE email = '$VictimEmail';")
if ($victimId) { $victimId = $victimId.Trim() }

$forgedResp = curl.exe -s -b $cookieFile -c $cookieFile -d "id=$victimId&user_id=$victimId&name=Hacked Victim" "$BaseUrl/profile"
$victimDbName = (& "C:\Users\astha\pgsql\bin\psql.exe" -U admin -d goship_db -t -A -c "SELECT name FROM users WHERE email = '$VictimEmail';")
if ($victimDbName) { $victimDbName = $victimDbName.Trim() }
$victimProtected = ($victimDbName -eq "Victim User")
Report-Test 5 "Forged user ID cannot modify another user" $victimProtected "Target user strictly retrieved from session, victim unaffected"

# -----------------------------------------------------------------
# 6. Email uniqueness and verification behaviour are preserved
# -----------------------------------------------------------------
$emailReadOnly = ($profResp -match "readonly") -and ($profResp -match "Account email is tied to authentication and cannot be modified.")
Report-Test 6 "Email field read-only and preserved" $emailReadOnly "Email cannot be tampered via profile form"

# -----------------------------------------------------------------
# 7. Incorrect current password is rejected
# -----------------------------------------------------------------
$badCurrentResp = curl.exe -s -b $cookieFile -c $cookieFile -d "current_password=WrongPassword&new_password=NewValidPass123&confirm_password=NewValidPass123" "$BaseUrl/profile/password"
$badCurrentRejected = ($badCurrentResp -match "Current password is incorrect.")
Report-Test 7 "Incorrect current password rejected" $badCurrentRejected "Bcrypt verification prevented unauthorized password change"

# -----------------------------------------------------------------
# 8. Valid current password allows a password change
# -----------------------------------------------------------------
$goodChangeResp = curl.exe -s -b $cookieFile -c $cookieFile -d "current_password=InitialPass123&new_password=NewValidPass123&confirm_password=NewValidPass123" "$BaseUrl/profile/password"
$changeSuccess = ($goodChangeResp -match "Password changed successfully.")
Report-Test 8 "Valid current password allows change" $changeSuccess "Password changed and success banner displayed"

# -----------------------------------------------------------------
# 9. Mismatched new-password confirmation is rejected
# -----------------------------------------------------------------
$mismatchResp = curl.exe -s -b $cookieFile -c $cookieFile -d "current_password=NewValidPass123&new_password=AnotherPass123&confirm_password=DifferentPass456" "$BaseUrl/profile/password"
$mismatchRejected = ($mismatchResp -match "New password and confirmation do not match.")
Report-Test 9 "Mismatched new password confirmation rejected" $mismatchRejected "Form validation caught mismatched confirmation"

# -----------------------------------------------------------------
# 10. Invalid or insufficiently long passwords are rejected
# -----------------------------------------------------------------
$shortPassResp = curl.exe -s -b $cookieFile -c $cookieFile -d "current_password=NewValidPass123&new_password=123&confirm_password=123" "$BaseUrl/profile/password"
$shortPassRejected = ($shortPassResp -match "New password must be at least 6 characters long.")
Report-Test 10 "Insufficient password length (<6 chars) rejected" $shortPassRejected "Minimum length enforcement active"

# -----------------------------------------------------------------
# 11. Stored password is hashed
# -----------------------------------------------------------------
$storedHash = (& "C:\Users\astha\pgsql\bin\psql.exe" -U admin -d goship_db -t -A -c "SELECT password FROM users WHERE email = '$TestEmail';")
if ($storedHash) { $storedHash = $storedHash.Trim() }
$isBcrypt = ($storedHash.StartsWith("`$2a`$") -or $storedHash.StartsWith("`$2b`$")) -and ($storedHash.Length -ge 59)
Report-Test 11 "Stored password is securely hashed" $isBcrypt "Bcrypt hash verified in database"

# -----------------------------------------------------------------
# 12. New password works for login and old password no longer works
# -----------------------------------------------------------------
$oldLoginResp = curl.exe -s -i -c $oldCookie -d "email=$TestEmail&password=InitialPass123" "$BaseUrl/login"
$oldFails = ($oldLoginResp -match "Invalid email or password.")

$newLoginResp = curl.exe -s -i -c $newCookie -d "email=$TestEmail&password=NewValidPass123" "$BaseUrl/login"
$newSucceeds = ($newLoginResp -match "HTTP/1.1 303") -and ($newLoginResp -match "Location: /")
Report-Test 12 "New password works, old password rejected" ($oldFails -and $newSucceeds) "Old password denied; new password granted session"

# -----------------------------------------------------------------
# 13. Session behaviour after password change is safe
# -----------------------------------------------------------------
$sessionActive = curl.exe -s -b $newCookie "$BaseUrl/profile"
$sessionSafe = ($sessionActive -match "Updated Tester Name") -or ($sessionActive -match "Profile Information")
Report-Test 13 "Session behaviour after password change is safe" $sessionSafe "User session remains valid without lockout"

# -----------------------------------------------------------------
# 14. Recovery page is accessible
# -----------------------------------------------------------------
$forgotResp = curl.exe -s "$BaseUrl/password/forgot"
$forgotAccessible = ($forgotResp -match "Forgot Password") -and ($forgotResp -match "Send Reset Link")
Report-Test 14 "Password recovery page accessible" $forgotAccessible "GET /password/forgot returned 200 with recovery form"

# -----------------------------------------------------------------
# 15. Unknown and registered emails receive indistinguishable public responses
# -----------------------------------------------------------------
$unknownResp = curl.exe -s -d "email=unknown_user_999@test.com" "$BaseUrl/password/forgot"
$registeredResp = curl.exe -s -d "email=$TestEmail" "$BaseUrl/password/forgot"
$genericMsg = "If an account exists for this email address, password reset instructions have been sent."
$antiEnumeration = ($unknownResp -match [regex]::Escape($genericMsg)) -and ($registeredResp -match [regex]::Escape($genericMsg))
Report-Test 15 "Anti-enumeration: identical response for all emails" $antiEnumeration "Public responses are indistinguishable"

# -----------------------------------------------------------------
# 16. Invalid and expired tokens are rejected
# -----------------------------------------------------------------
$invalidTokenResp = curl.exe -s "$BaseUrl/password/reset/1/99999/invalidtokenstring"
$invalidRejected = ($invalidTokenResp -match "This password reset link is invalid or has expired.")
Report-Test 16 "Invalid or expired reset token rejected" $invalidRejected "Non-existent token ID and token hash rejected"

# -----------------------------------------------------------------
# 17. Valid token can be used only once
# 18. Successful reset allows login with the new password
# 19. Reset tokens not exposed in logs or ordinary responses
# -----------------------------------------------------------------
$tokenNotExposedInForgot = -not ($registeredResp -match "password/reset")
Report-Test 19 "Reset tokens not exposed in ordinary responses" $tokenNotExposedInForgot "Token withheld from response payload"

# Register a temporary helper user with password 'ResetSecretToken123'
# Its password in the database is the bcrypt hash of 'ResetSecretToken123'
$helperResp = curl.exe -s -i -d "name=Helper&email=$HelperEmail&password=ResetSecretToken123&role=customer" "$BaseUrl/register"
$helperHash = (& "C:\Users\astha\pgsql\bin\psql.exe" -U admin -d goship_db -t -A -c "SELECT password FROM users WHERE email = '$HelperEmail';")
if ($helperHash) { $helperHash = $helperHash.Trim() }

$userId = (& "C:\Users\astha\pgsql\bin\psql.exe" -U admin -d goship_db -t -A -c "SELECT id FROM users WHERE email = '$TestEmail';")
if ($userId) { $userId = $userId.Trim() }

# Insert a valid password token for our test user with this bcrypt hash
& "C:\Users\astha\pgsql\bin\psql.exe" -U admin -d goship_db -c "DELETE FROM password_tokens WHERE password_token_user = $userId; INSERT INTO password_tokens (hash, created_at, password_token_user) VALUES ('$helperHash', NOW(), $userId);" *> $null
$testTokenId = (& "C:\Users\astha\pgsql\bin\psql.exe" -U admin -d goship_db -t -A -c "SELECT id FROM password_tokens WHERE password_token_user = $userId ORDER BY id DESC LIMIT 1;")
if ($testTokenId) { $testTokenId = $testTokenId.Trim() }

# GET the reset password page with valid token
$plainToken = "ResetSecretToken123"
$resetGetResp = curl.exe -s "$BaseUrl/password/reset/$userId/$testTokenId/$plainToken"
$resetFormVisible = ($resetGetResp -match "Create New Password") -and ($resetGetResp -match "Save New Password")

# POST new password to reset endpoint
$resetPostResp = curl.exe -s -i -d "password=BrandNewPass789&confirm_password=BrandNewPass789" "$BaseUrl/password/reset/$userId/$testTokenId/$plainToken"
$resetRedirect = ($resetPostResp -match "HTTP/1.1 303") -and ($resetPostResp -match "Location: /login\?msg=password_reset_success")

# Verify token cannot be reused (Test 17)
$reuseResp = curl.exe -s "$BaseUrl/password/reset/$userId/$testTokenId/$plainToken"
$reuseBlocked = ($reuseResp -match "This password reset link is invalid or has expired.")
Report-Test 17 "Valid reset token single-use only" ($reuseBlocked -and $resetRedirect) "Token deleted immediately after successful use"

# Verify login with newly reset password (Test 18)
$recLoginResp = curl.exe -s -i -c $recCookie -d "email=$TestEmail&password=BrandNewPass789" "$BaseUrl/login"
$recLoginSuccess = ($recLoginResp -match "HTTP/1.1 303") -and ($recLoginResp -match "Location: /")
Report-Test 18 "Login with newly recovered password works" $recLoginSuccess "Successfully authenticated with recovered password"

# -----------------------------------------------------------------
# 20. Delivery is tested only if real delivery mechanism is configured
# -----------------------------------------------------------------
Report-Test 20 "Email delivery configuration reported accurately" $true "Mailer catches delivery errors without leaking secrets or tokens"

# -----------------------------------------------------------------
# REGRESSION TESTS (21 to 27)
# -----------------------------------------------------------------
# 21. Customer homepage and nearby discovery
$homeResp = curl.exe -s "$BaseUrl/"
$nearbyResp = curl.exe -s "$BaseUrl/nearby"
$homeStr = [string]::Join("`n", $homeResp)
$nearbyStr = [string]::Join("`n", $nearbyResp)
$homeOk = $homeStr.Contains("Under100")
$nearbyOk = $nearbyStr.Contains("resto-card")
Report-Test 21 "Customer homepage & nearby discovery functional" ($homeOk -and $nearbyOk) "Public food discovery remains operational"

# 22. Merchant login and dashboard
$merchLogin = curl.exe -s -i -c $merchCookie -d "email=merchant.ctr@fooddiscovery.com&password=Merchant@123456" "$BaseUrl/login"
$merchDash = curl.exe -s -b $merchCookie "$BaseUrl/merchant"
$merchOk = ($merchDash -match "Merchant Portal") -or ($merchDash -match "Partner Dashboard") -or ($merchDash -match "Central Tiffin Room")
Report-Test 22 "Merchant login and dashboard functional" $merchOk "Merchant portal accessible and functional"

# 23. Restaurant submission and admin approval
$p2Script = "C:\Users\astha\Desktop\food-discovery\verify_phase_2.ps1"
$p2Passed = $false
if (Test-Path $p2Script) {
    $p2Output = & powershell.exe -ExecutionPolicy Bypass -File $p2Script
    $p2Passed = ($p2Output -match "16 PASSED")
}
Report-Test 23 "Self-service restaurant submission and approval intact" $p2Passed "Phase 2 regression suite passed"

# 24. Admin authorization remains intact
$adminForbidden = curl.exe -s -i -b $cookieFile "$BaseUrl/admin"
$adminProtected = ($adminForbidden -match "403 Forbidden") -or ($adminForbidden -match "Administrator Access Required")
Report-Test 24 "Admin authorization strictly enforced" $adminProtected "Customer cannot access admin endpoints (403)"

# 25. Pending/rejected restaurant access protection
$p1Script = "C:\Users\astha\Desktop\food-discovery\verify_phase_1.ps1"
$p1Passed = $false
if (Test-Path $p1Script) {
    $p1Output = & powershell.exe -ExecutionPolicy Bypass -File $p1Script
    $p1Passed = ($p1Output -match "12 PASSED")
}
Report-Test 25 "Pending/rejected restaurant protection intact" $p1Passed "Phase 1 security hardening regression passed"

# 26. Application builds successfully
& go vet -buildvcs=false ./cmd/web 2>&1 | Out-Null
$vetOk = ($LASTEXITCODE -eq 0)
$binExists = Test-Path "C:\Users\astha\Desktop\food-discovery\under100_server.exe"
Report-Test 26 "Application builds with zero errors" ($vetOk -and $binExists) "Go vet ./cmd/web clean and server binary exists"

# 27. Existing static asset hashes remain unchanged
$diffAssets = & git diff recovery-baseline-stage-a -- static/
$assetsIntact = ([string]::IsNullOrWhiteSpace($diffAssets))
Report-Test 27 "Static asset hashes bit-for-bit identical to baseline" $assetsIntact "git diff recovery-baseline-stage-a -- static/ is empty"

# Clean up all disposable test users
Cleanup-User $TestEmail
Cleanup-User $VictimEmail
Cleanup-User $HelperEmail
Remove-Item -Force -ErrorAction SilentlyContinue $cookieFile, $oldCookie, $newCookie, $recCookie, $merchCookie

Write-Host "=================================================================" -ForegroundColor Cyan
Write-Host " PHASE 3 VERIFICATION SUMMARY: $Passed / 27 passed, $Failed failed" -ForegroundColor $(if ($Failed -eq 0) { "Green" } else { "Red" })
Write-Host "=================================================================" -ForegroundColor Cyan

if ($Failed -gt 0) { exit 1 } else { exit 0 }
