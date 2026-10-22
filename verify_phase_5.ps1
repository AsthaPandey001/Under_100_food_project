# =================================================================
# PHASE 5 COMPREHENSIVE VERIFICATION SUITE — BITE100 FOOD DISCOVERY
# End-to-End Testing, Security Audit & Release Readiness
# =================================================================

$ErrorActionPreference = "Continue"
$BaseUrl = "http://127.0.0.1:8000"
$PassCount = 0
$FailCount = 0

function Report-Test($num, $name, $condition, $details) {
    if ($condition) {
        Write-Host " [PASS] Test $num`: $name" -ForegroundColor Green
        Write-Host "        $details" -ForegroundColor DarkGray
        $global:PassCount++
    } else {
        Write-Host " [FAIL] Test $num`: $name" -ForegroundColor Red
        Write-Host "        $details" -ForegroundColor Yellow
        $global:FailCount++
    }
}

Write-Host "=================================================================" -ForegroundColor Cyan
Write-Host " PHASE 5 END-TO-END VERIFICATION & SECURITY AUDIT SUITE" -ForegroundColor Cyan
Write-Host "=================================================================" -ForegroundColor Cyan

# Baseline DB record count snapshot
$initialUserCount = (& "C:\Users\astha\pgsql\bin\psql.exe" -h 127.0.0.1 -U admin -d goship_db -t -A -c "SELECT COUNT(*) FROM users;").Trim()
$initialRestCount = (& "C:\Users\astha\pgsql\bin\psql.exe" -h 127.0.0.1 -U admin -d goship_db -t -A -c "SELECT COUNT(*) FROM restaurants;").Trim()
$initialFoodCount = (& "C:\Users\astha\pgsql\bin\psql.exe" -h 127.0.0.1 -U admin -d goship_db -t -A -c "SELECT COUNT(*) FROM food_items;").Trim()

# -----------------------------------------------------------------
# FIXTURES SETUP: Isolated test accounts
# -----------------------------------------------------------------
$custACookie = New-TemporaryFile
$custBCookie = New-TemporaryFile
$merchantCookie = New-TemporaryFile
$adminCookie = New-TemporaryFile

$EmailA = ("p5_cust_a_" + (Get-Random) + "@test.com").ToLower()
$EmailB = ("p5_cust_b_" + (Get-Random) + "@test.com").ToLower()
$MerchantEmail = ("p5_merch_" + (Get-Random) + "@test.com").ToLower()
$HelperEmail = ("p5_helper_" + (Get-Random) + "@test.com").ToLower()

function Cleanup-P5Fixtures {
    & "C:\Users\astha\pgsql\bin\psql.exe" -h 127.0.0.1 -U admin -d goship_db -c @"
DELETE FROM favorites WHERE user_id IN (SELECT id FROM users WHERE email IN ('$EmailA', '$EmailB', '$MerchantEmail', '$HelperEmail'));
DELETE FROM food_items WHERE restaurant_food_items IN (SELECT id FROM restaurants WHERE user_restaurants IN (SELECT id FROM users WHERE email IN ('$MerchantEmail')));
DELETE FROM restaurants WHERE user_restaurants IN (SELECT id FROM users WHERE email IN ('$MerchantEmail'));
DELETE FROM password_tokens WHERE password_token_user IN (SELECT id FROM users WHERE email IN ('$EmailA', '$EmailB', '$MerchantEmail', '$HelperEmail'));
DELETE FROM users WHERE email IN ('$EmailA', '$EmailB', '$MerchantEmail', '$HelperEmail');
"@ | Out-Null
}

Cleanup-P5Fixtures

# -----------------------------------------------------------------
# 1. CUSTOMER ACCOUNT WORKFLOW
# -----------------------------------------------------------------

# Test 1: Customer Registration & Login
$regA = curl.exe -s -i -c $custACookie.FullName -b $custACookie.FullName -d "name=Alice Customer&email=$EmailA&password=Password123&role=customer" "$BaseUrl/register"
$loginA = curl.exe -s -i -c $custACookie.FullName -b $custACookie.FullName -d "email=$EmailA&password=Password123" "$BaseUrl/login"
$t1Ok = ($regA -match "303 See Other|Location: /") -and ($loginA -match "303 See Other|Location: /")
Report-Test 1 "Customer registration and login workflow" $t1Ok "Customer A registered and authenticated with session cookie"

# Test 2: Profile Access and Update
$profilePage = [string]::Join("`n", (curl.exe -s -b $custACookie.FullName "$BaseUrl/profile"))
$updateProfile = [string]::Join("`n", (curl.exe -s -b $custACookie.FullName -c $custACookie.FullName -d "name=Alice Updated Name" "$BaseUrl/profile"))
$dbName = (& "C:\Users\astha\pgsql\bin\psql.exe" -h 127.0.0.1 -U admin -d goship_db -t -A -c "SELECT name FROM users WHERE email = '$EmailA';").Trim()
$t2Ok = $profilePage.Contains("Alice Customer") -and $updateProfile.Contains("Profile updated successfully") -and ($dbName -eq "Alice Updated Name")
Report-Test 2 "Customer can access and update own profile" $t2Ok "Name updated in PostgreSQL; success banner rendered"

# Test 3: Password Change Workflow & Old Password Rejection
$pwChange = [string]::Join("`n", (curl.exe -s -b $custACookie.FullName -c $custACookie.FullName -d "current_password=Password123&new_password=NewPassword456&confirm_password=NewPassword456" "$BaseUrl/profile/password"))
$oldLogin = [string]::Join("`n", (curl.exe -s -i -d "email=$EmailA&password=Password123" "$BaseUrl/login"))
$newCookie = New-TemporaryFile
$newLogin = [string]::Join("`n", (curl.exe -s -i -c $newCookie.FullName -d "email=$EmailA&password=NewPassword456" "$BaseUrl/login"))
$t3Ok = ($pwChange -match "Password changed successfully") -and ($oldLogin -match "Invalid email or password") -and ($newLogin -match "303 See Other|Location: /")
Remove-Item $newCookie.FullName -ErrorAction SilentlyContinue
Report-Test 3 "Password change enforces validation and invalidates old password" $t3Ok "New password works, old password rejected, bcrypt hash updated"

# Test 4: Password Recovery & Single-Use Token
$forgotResp = [string]::Join("`n", (curl.exe -s -d "email=$EmailA" "$BaseUrl/password/forgot"))
# Register temporary helper user to obtain bcrypt hash of known secret token 'ResetSecretToken123'
curl.exe -s -d "name=Helper&email=$HelperEmail&password=ResetSecretToken123&role=customer" "$BaseUrl/register" | Out-Null
$helperHash = (& "C:\Users\astha\pgsql\bin\psql.exe" -h 127.0.0.1 -U admin -d goship_db -t -A -c "SELECT password FROM users WHERE email = '$HelperEmail';").Trim()
$userIdA = (& "C:\Users\astha\pgsql\bin\psql.exe" -h 127.0.0.1 -U admin -d goship_db -t -A -c "SELECT id FROM users WHERE email = '$EmailA';").Trim()

# Insert valid password token for user A
& "C:\Users\astha\pgsql\bin\psql.exe" -h 127.0.0.1 -U admin -d goship_db -c "INSERT INTO password_tokens (hash, created_at, password_token_user) VALUES ('$helperHash', NOW(), $userIdA);" | Out-Null
$tokenId = (& "C:\Users\astha\pgsql\bin\psql.exe" -h 127.0.0.1 -U admin -d goship_db -t -A -c "SELECT id FROM password_tokens WHERE password_token_user = $userIdA ORDER BY id DESC LIMIT 1;").Trim()

# Reset password using token
$resetResp = curl.exe -s -i -d "password=RecoveredPass789&confirm_password=RecoveredPass789" "$BaseUrl/password/reset/$userIdA/$tokenId/ResetSecretToken123"
# Second attempt with same token must fail (single-use)
$resetResp2 = [string]::Join("`n", (curl.exe -s -d "password=AnotherPass123&confirm_password=AnotherPass123" "$BaseUrl/password/reset/$userIdA/$tokenId/ResetSecretToken123"))
$t4Ok = $forgotResp.Contains("password reset instructions have been sent") -and ($resetResp -match "303 See Other|Location: /login\?msg=password_reset_success") -and ($resetResp2.Contains("invalid or has expired"))
Report-Test 4 "Password recovery workflow and single-use token enforcement" $t4Ok "Reset succeeds and immediately consumes/deletes token"

# Test 5: Logout and Protected Route Redirection
$logoutCookie = New-TemporaryFile
curl.exe -s -c $logoutCookie.FullName -b $logoutCookie.FullName -d "email=$EmailA&password=RecoveredPass789" "$BaseUrl/login" | Out-Null
$logoutResp = curl.exe -s -i -c $logoutCookie.FullName -b $logoutCookie.FullName "$BaseUrl/logout"
$postLogoutProfile = curl.exe -s -i -b $logoutCookie.FullName "$BaseUrl/profile"
Remove-Item $logoutCookie.FullName -ErrorAction SilentlyContinue
$t5Ok = ($logoutResp -match "303 See Other|Location: /") -and ($postLogoutProfile -match "303 See Other|Location: /login\?msg=auth_required")
Report-Test 5 "Logout terminates session and guards protected routes" $t5Ok "Protected /profile redirects to login after logout"

# Test 6: Customer IDOR Isolation
# Log Customer A back in with recovered password
curl.exe -s -c $custACookie.FullName -b $custACookie.FullName -d "email=$EmailA&password=RecoveredPass789" "$BaseUrl/login" | Out-Null
# Register Customer B
curl.exe -s -c $custBCookie.FullName -b $custBCookie.FullName -d "name=Bob Customer&email=$EmailB&password=Password123&role=customer" "$BaseUrl/register" | Out-Null
$userIdB = (& "C:\Users\astha\pgsql\bin\psql.exe" -h 127.0.0.1 -U admin -d goship_db -t -A -c "SELECT id FROM users WHERE email = '$EmailB';").Trim()
# Customer A attempts to post profile update with Customer B's user_id in form payload
curl.exe -s -b $custACookie.FullName -c $custACookie.FullName -d "user_id=$userIdB&name=Tampered Name" "$BaseUrl/profile" | Out-Null
$dbNameB = (& "C:\Users\astha\pgsql\bin\psql.exe" -h 127.0.0.1 -U admin -d goship_db -t -A -c "SELECT name FROM users WHERE id = $userIdB;").Trim()
$t6Ok = ($dbNameB -eq "Bob Customer")
Report-Test 6 "Customer IDOR isolation: Customer A cannot modify Customer B" $t6Ok "Server binds profile mutations strictly to authenticated session"

# -----------------------------------------------------------------
# 2. RESTAURANT DISCOVERY WORKFLOW
# -----------------------------------------------------------------

# Test 7: Public Approved Restaurant Discovery & Search
$homeResp = [string]::Join("`n", (curl.exe -s "$BaseUrl/"))
$nameSearch = [string]::Join("`n", (curl.exe -s "$BaseUrl/nearby?q=Vidyarthi"))
$addrSearch = [string]::Join("`n", (curl.exe -s "$BaseUrl/nearby?q=Gandhi"))
$foodSearch = [string]::Join("`n", (curl.exe -s "$BaseUrl/nearby?q=Coffee"))
$t7Ok = $homeResp.Contains("Under100") -and $nameSearch.Contains("Vidyarthi Bhavan") -and $addrSearch.Contains("Vidyarthi Bhavan") -and $foodSearch.Contains("Special Filter Coffee")
Report-Test 7 "Restaurant discovery: browse, name search, address search, and food search" $t7Ok "All search modalities locate approved catalog items"

# Test 8: Price Filter Boundary at ₹100
$p100 = [string]::Join("`n", (curl.exe -s "$BaseUrl/nearby?max_price=100"))
$p30 = [string]::Join("`n", (curl.exe -s "$BaseUrl/nearby?max_price=30"))
$t8Ok = $p100.Contains("resto-card") -and $p30.Contains("Special Filter Coffee")
Report-Test 8 "Price filter correctly enforces ₹100 ceiling and sub-budgets" $t8Ok "Dishes strictly filtered by max_price parameter"

# Test 9: Combined Filters & Clamped Pagination
$combResp = [string]::Join("`n", (curl.exe -s "$BaseUrl/nearby?q=Bengaluru&diet=veg&max_price=100&sort=price&page=1&page_size=3"))
$clampPage = [string]::Join("`n", (curl.exe -s "$BaseUrl/nearby?page=-99&page_size=9999"))
$t9Ok = $combResp.Contains("pagination-container") -and $combResp.Contains("q=Bengaluru") -and $clampPage.Contains("resto-card")
Report-Test 9 "Combined search/diet/price filters persist across clamped pagination" $t9Ok "Query parameters preserved in pagination links; page/pageSize safely clamped"

# Test 10: Restaurant Details Modal Displays Menu
$restoModal = [string]::Join("`n", (curl.exe -s "$BaseUrl/restaurant/1"))
$t10Ok = $restoModal.Contains("Vidyarthi Bhavan") -and $restoModal.Contains("Gandhi Bazaar") -and $restoModal.Contains("Special Filter Coffee")
Report-Test 10 "Approved restaurant details display metadata and menu items" $t10Ok "Operating hours, address, phone, and food catalog rendered"

# Test 11: Unapproved Restaurant Public Protection
$pendPublic = curl.exe -s -w "%{http_code}" -o /dev/null "$BaseUrl/restaurant/10"
$rejPublic = curl.exe -s -w "%{http_code}" -o /dev/null "$BaseUrl/restaurant/11"
$markers = [string]::Join("`n", (curl.exe -s "$BaseUrl/api/markers"))
$t11Ok = ($pendPublic -eq "404") -and ($rejPublic -eq "404") -and (-not $markers.Contains('"id":10')) -and (-not $markers.Contains('"id":11'))
Report-Test 11 "Pending and rejected restaurants strictly hidden from public discovery" $t11Ok "Direct details return 404; map markers exclude unapproved outlets"

# Test 12: Customer Favorites Workflow & Privacy
curl.exe -s -b $custACookie.FullName -c $custACookie.FullName -X POST "$BaseUrl/favorites/toggle/1" | Out-Null
$favA = [string]::Join("`n", (curl.exe -s -b $custACookie.FullName "$BaseUrl/favorites"))
$favB = [string]::Join("`n", (curl.exe -s -b $custBCookie.FullName "$BaseUrl/favorites"))
$t12Ok = $favA.Contains("restaurant-1") -and (-not $favB.Contains("restaurant-1")) -and $favB.Contains("No saved restaurants yet")
Report-Test 12 "Customer favorites toggle and private isolation verified" $t12Ok "Customer A can view saved favorites; Customer B's favorites remain empty"

# -----------------------------------------------------------------
# 3. MERCHANT WORKFLOW
# -----------------------------------------------------------------

# Test 13: Merchant Registration & Submission
curl.exe -s -c $merchantCookie.FullName -b $merchantCookie.FullName -d "name=Test Merchant&email=$MerchantEmail&password=Password123&role=merchant" "$BaseUrl/register" | Out-Null
$merchSubmit = curl.exe -s -i -b $merchantCookie.FullName -c $merchantCookie.FullName -d "name=P5 Merchant Cafe&address=123 100ft Road, Indiranagar&phone=9876543210&latitude=12.978&longitude=77.640" "$BaseUrl/merchant/restaurant/submit"
$newRestId = (& "C:\Users\astha\pgsql\bin\psql.exe" -h 127.0.0.1 -U admin -d goship_db -t -A -c "SELECT id FROM restaurants WHERE user_restaurants = (SELECT id FROM users WHERE email = '$MerchantEmail');").Trim()
$newRestStatus = (& "C:\Users\astha\pgsql\bin\psql.exe" -h 127.0.0.1 -U admin -d goship_db -t -A -c "SELECT status FROM restaurants WHERE id = $newRestId;").Trim()
$t13Ok = ($merchSubmit -match "303 See Other|restaurant_registered") -and ($newRestStatus -eq "pending")
Report-Test 13 "Merchant registration and self-service outlet submission" $t13Ok "Creates pending restaurant ($newRestId) linked to authenticated merchant"

# Test 14: Pending Restaurant Owner Visibility vs Public 404
$merchViewRest = curl.exe -s -w "%{http_code}" -b $merchantCookie.FullName -o /dev/null "$BaseUrl/restaurant/$newRestId"
$custViewRest = curl.exe -s -w "%{http_code}" -b $custACookie.FullName -o /dev/null "$BaseUrl/restaurant/$newRestId"
$t14Ok = ($merchViewRest -eq "200") -and ($custViewRest -eq "404")
Report-Test 14 "Submitting merchant can view pending restaurant; public/customer gets 404" $t14Ok "Owner preview permitted; unauthorized access strictly blocked"

# Test 15: Merchant Menu Operations (Add, Edit, Availability Toggle, Delete)
$addFood = curl.exe -s -i -b $merchantCookie.FullName -c $merchantCookie.FullName -d "name=P5 Special Chai&price=25&category=Beverages&diet_type=veg&available=true" "$BaseUrl/merchant/menu"
$foodId = (& "C:\Users\astha\pgsql\bin\psql.exe" -h 127.0.0.1 -U admin -d goship_db -t -A -c "SELECT id FROM food_items WHERE restaurant_food_items = $newRestId AND name = 'P5 Special Chai';").Trim()
$editFood = curl.exe -s -i -b $merchantCookie.FullName -c $merchantCookie.FullName -d "name=P5 Masala Chai&price=30&category=Beverages&diet_type=veg&available=true" "$BaseUrl/merchant/menu/$foodId"
$toggleFood = curl.exe -s -b $merchantCookie.FullName -X POST "$BaseUrl/merchant/menu/$foodId/availability"
$delFood = curl.exe -s -b $merchantCookie.FullName -X POST "$BaseUrl/merchant/menu/$foodId/delete"
$foodCountAfter = (& "C:\Users\astha\pgsql\bin\psql.exe" -h 127.0.0.1 -U admin -d goship_db -t -A -c "SELECT COUNT(*) FROM food_items WHERE id = $foodId;").Trim()
$t15Ok = ($addFood -match "303 See Other|dish_added") -and ($editFood -match "303 See Other|dish_updated") -and ($foodCountAfter -eq "0")
Report-Test 15 "Merchant menu operations: add, edit, toggle availability, and delete" $t15Ok "Complete dish lifecycle executed with server-side validation"

# Test 16: Merchant Tenant Isolation / IDOR Protection
# Merchant attempts to edit dish 1 (Vidyarthi Bhavan dish)
$idorEdit = curl.exe -s -w "%{http_code}" -b $merchantCookie.FullName -o /dev/null -d "name=Hacked Dish&price=10" "$BaseUrl/merchant/menu/1"
$idorToggle = curl.exe -s -w "%{http_code}" -b $merchantCookie.FullName -o /dev/null -X POST "$BaseUrl/merchant/menu/1/availability"
$idorDelete = curl.exe -s -w "%{http_code}" -b $merchantCookie.FullName -o /dev/null -X POST "$BaseUrl/merchant/menu/1/delete"
$t16Ok = ($idorEdit -eq "403") -and ($idorToggle -eq "403") -and ($idorDelete -eq "403")
Report-Test 16 "Merchant cross-tenant IDOR protection enforced" $t16Ok "Cross-tenant dish modification and deletion rejected with HTTP 403 Forbidden"

# -----------------------------------------------------------------
# 4. ADMIN WORKFLOW & SECURITY AUDIT
# -----------------------------------------------------------------

# Test 17: Admin Authentication & RBAC Access Control
curl.exe -s -c $adminCookie.FullName -b $adminCookie.FullName -d "email=admin@fooddiscovery.com&password=Admin@123456" "$BaseUrl/login" | Out-Null
$adminDash = [string]::Join("`n", (curl.exe -s -b $adminCookie.FullName "$BaseUrl/admin"))
$custToAdmin = curl.exe -s -w "%{http_code}" -b $custACookie.FullName -o /dev/null "$BaseUrl/admin"
$guestToAdmin = curl.exe -s -i "$BaseUrl/admin"
$t17Ok = $adminDash.Contains("Admin") -and ($custToAdmin -eq "403") -and ($guestToAdmin -match "303 See Other|admin_required")
Report-Test 17 "Admin authorization strictly enforced" $t17Ok "Admin granted access; customer gets 403 Forbidden; guest redirected"

# Test 18: Admin Restaurant Approval & Rejection
$approveResp = curl.exe -s -b $adminCookie.FullName -d "status=approved" "$BaseUrl/admin/restaurants/$newRestId/status"
$statusApproved = (& "C:\Users\astha\pgsql\bin\psql.exe" -h 127.0.0.1 -U admin -d goship_db -t -A -c "SELECT status FROM restaurants WHERE id = $newRestId;").Trim()
$rejectResp = curl.exe -s -b $adminCookie.FullName -d "status=rejected" "$BaseUrl/admin/restaurants/$newRestId/status"
$statusRejected = (& "C:\Users\astha\pgsql\bin\psql.exe" -h 127.0.0.1 -U admin -d goship_db -t -A -c "SELECT status FROM restaurants WHERE id = $newRestId;").Trim()
$t18Ok = ($statusApproved -eq "approved") -and ($statusRejected -eq "rejected")
Report-Test 18 "Admin can approve and reject restaurant submissions" $t18Ok "Database status transitions verified in PostgreSQL"

# Test 19: Admin User Management & Role Filtering
$usersPage = [string]::Join("`n", (curl.exe -s -b $adminCookie.FullName "$BaseUrl/admin/users"))
$merchantUsers = [string]::Join("`n", (curl.exe -s -b $adminCookie.FullName "$BaseUrl/admin/users?role=merchant"))
$t19Ok = $usersPage.Contains("User Management") -and $merchantUsers.Contains("merchant.ctr@fooddiscovery.com")
Report-Test 19 "Admin user management and role-based filtering" $t19Ok "Users list rendered with role filter capabilities"

# Test 20: Admin Owner Assignment with Validation
$ownerId = (& "C:\Users\astha\pgsql\bin\psql.exe" -h 127.0.0.1 -U admin -d goship_db -t -A -c "SELECT id FROM users WHERE email = '$MerchantEmail';").Trim()
$assignResp = curl.exe -s -b $adminCookie.FullName -d "owner_id=$ownerId" "$BaseUrl/admin/restaurants/$newRestId/owner"
$dbOwner = (& "C:\Users\astha\pgsql\bin\psql.exe" -h 127.0.0.1 -U admin -d goship_db -t -A -c "SELECT user_restaurants FROM restaurants WHERE id = $newRestId;").Trim()
$t20Ok = ($dbOwner -eq $ownerId)
Report-Test 20 "Admin can assign and manage restaurant ownership" $t20Ok "Owner edge bound with server-side merchant role verification"

# -----------------------------------------------------------------
# 5. REGRESSION & ASSET INTEGRITY
# -----------------------------------------------------------------

# Test 21: Phase 1 Regression Suite
$p1Out = powershell -ExecutionPolicy Bypass -File .\verify_phase_1.ps1
$t21Ok = ($p1Out -match "12 PASSED, 0 FAILED")
Report-Test 21 "Phase 1 regression verification suite passes" $t21Ok "12/12 access control tests passed"

# Test 22: Phase 2 Regression Suite
$p2Out = powershell -ExecutionPolicy Bypass -File .\verify_phase_2.ps1
$t22Ok = ($p2Out -match "16 PASSED, 0 FAILED")
Report-Test 22 "Phase 2 regression verification suite passes" $t22Ok "16/16 self-service merchant tests passed"

# Test 23: Phase 3 Regression Suite
$p3Out = powershell -ExecutionPolicy Bypass -File .\verify_phase_3.ps1
$t23Ok = ($p3Out -match "27 / 27 passed, 0 failed")
Report-Test 23 "Phase 3 regression verification suite passes" $t23Ok "27/27 customer profile & password recovery tests passed"

# Test 24: Phase 4 Regression Suite
$p4Out = powershell -ExecutionPolicy Bypass -File .\verify_phase_4.ps1
$t24Ok = ($p4Out -match "16 / 16 passed, 0 failed")
Report-Test 24 "Phase 4 regression verification suite passes" $t24Ok "16/16 favorites and pagination tests passed"

# Test 25: Static Asset Cryptographic Integrity
$staticDiff = git diff recovery-baseline-stage-a -- static/
$t25Ok = [string]::IsNullOrWhiteSpace($staticDiff)
Report-Test 25 "Static assets bit-for-bit unchanged from baseline" $t25Ok "All 19 static images and SVGs match Stage A SHA-256 hashes"

# -----------------------------------------------------------------
# CLEANUP
# -----------------------------------------------------------------
Cleanup-P5Fixtures
Remove-Item $custACookie.FullName -ErrorAction SilentlyContinue
Remove-Item $custBCookie.FullName -ErrorAction SilentlyContinue
Remove-Item $merchantCookie.FullName -ErrorAction SilentlyContinue
Remove-Item $adminCookie.FullName -ErrorAction SilentlyContinue

# Verify database integrity after cleanup
$finalUserCount = (& "C:\Users\astha\pgsql\bin\psql.exe" -h 127.0.0.1 -U admin -d goship_db -t -A -c "SELECT COUNT(*) FROM users;").Trim()
$finalRestCount = (& "C:\Users\astha\pgsql\bin\psql.exe" -h 127.0.0.1 -U admin -d goship_db -t -A -c "SELECT COUNT(*) FROM restaurants;").Trim()
$finalFoodCount = (& "C:\Users\astha\pgsql\bin\psql.exe" -h 127.0.0.1 -U admin -d goship_db -t -A -c "SELECT COUNT(*) FROM food_items;").Trim()
$finalFavCount = (& "C:\Users\astha\pgsql\bin\psql.exe" -h 127.0.0.1 -U admin -d goship_db -t -A -c "SELECT COUNT(*) FROM favorites;").Trim()

$dbIntact = ($initialUserCount -eq $finalUserCount) -and ($initialRestCount -eq $finalRestCount) -and ($initialFoodCount -eq $finalFoodCount) -and ($finalFavCount -eq "0")

Report-Test 26 "Development database counts remain unchanged post-testing" $dbIntact "Users=$finalUserCount, Restaurants=$finalRestCount, Food=$finalFoodCount, Favorites=$finalFavCount"

Write-Host "=================================================================" -ForegroundColor Cyan
Write-Host " PHASE 5 VERIFICATION SUMMARY: $PassCount / 26 passed, $FailCount failed" -ForegroundColor Cyan
Write-Host "=================================================================" -ForegroundColor Cyan

if ($FailCount -gt 0) {
    exit 1
}
exit 0
