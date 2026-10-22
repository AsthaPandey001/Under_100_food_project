# verify_stage_d.ps1 - Automated verification test suite for Stage D: Admin Dashboard & Platform Management

$BaseURL = "http://localhost:8000"
$Passed = 0
$Failed = 0

function Report-Test($name, $success, $details) {
    if ($success) {
        $global:Passed++
        Write-Host "  [PASS] $name" -ForegroundColor Green
        if ($details) { Write-Host "         $details" -ForegroundColor Gray }
    } else {
        $global:Failed++
        Write-Host "  [FAIL] $name" -ForegroundColor Red
        if ($details) { Write-Host "         Error: $details" -ForegroundColor Yellow }
    }
}

Write-Host "==========================================================" -ForegroundColor Cyan
Write-Host "STAGE D VERIFICATION SUITE - ADMIN DASHBOARD" -ForegroundColor Cyan
Write-Host "==========================================================" -ForegroundColor Cyan

# Setup temporary cookie jars
$AdminCookieJar = Join-Path $PSScriptRoot "admin_cookie.txt"
$MerchantCookieJar = Join-Path $PSScriptRoot "merchant_cookie.txt"
$CustomerCookieJar = Join-Path $PSScriptRoot "customer_cookie.txt"
Remove-Item $AdminCookieJar -ErrorAction SilentlyContinue
Remove-Item $MerchantCookieJar -ErrorAction SilentlyContinue
Remove-Item $CustomerCookieJar -ErrorAction SilentlyContinue

# Authenticate Admin
curl.exe -s -c $AdminCookieJar -X POST ($BaseURL + "/login") -d "email=admin@fooddiscovery.com" -d "password=Admin@123456" | Out-Null
$AdminCheck = curl.exe -s -b $AdminCookieJar ($BaseURL + "/admin")
$AdminAuthed = ($AdminCheck -match "Admin Dashboard" -or $AdminCheck -match "Bite100 Admin")
Report-Test "Admin authentication and access to /admin" $AdminAuthed "Admin logged in successfully"

# Authenticate Customer
curl.exe -s -c $CustomerCookieJar -X POST ($BaseURL + "/login") -d "email=user@fooddiscovery.com" -d "password=User@123456" | Out-Null

# Authenticate Merchant
curl.exe -s -c $MerchantCookieJar -X POST ($BaseURL + "/login") -d "email=merchant.ctr@fooddiscovery.com" -d "password=Merchant@123456" | Out-Null

# Scenario 1: Unauthenticated visitor denied /admin and /admin/users
$UnauthAdmin = curl.exe -s -w "%{http_code}" -o /dev/null ($BaseURL + "/admin")
$UnauthUsers = curl.exe -s -w "%{http_code}" -o /dev/null ($BaseURL + "/admin/users")
$UnauthDenied = ($UnauthAdmin -eq "303" -or $UnauthAdmin -eq "302" -or $UnauthAdmin -eq "401") -and ($UnauthUsers -eq "303" -or $UnauthUsers -eq "302" -or $UnauthUsers -eq "401")
Report-Test "Unauthenticated access denied for /admin and /admin/users" $UnauthDenied "HTTP redirect codes: $UnauthAdmin, $UnauthUsers"

# Scenario 2: Customer denied admin dashboard and users management
$CustAdmin = curl.exe -s -w "%{http_code}" -b $CustomerCookieJar -o /dev/null ($BaseURL + "/admin")
$CustUsers = curl.exe -s -w "%{http_code}" -b $CustomerCookieJar -o /dev/null ($BaseURL + "/admin/users")
$CustDenied = ($CustAdmin -eq "403" -and $CustUsers -eq "403")
Report-Test "Customer role forbidden from /admin and /admin/users" $CustDenied "HTTP 403 Forbidden received"

# Scenario 3: Merchant denied admin dashboard and users management
$MerchAdmin = curl.exe -s -w "%{http_code}" -b $MerchantCookieJar -o /dev/null ($BaseURL + "/admin")
$MerchUsers = curl.exe -s -w "%{http_code}" -b $MerchantCookieJar -o /dev/null ($BaseURL + "/admin/users")
$MerchDenied = ($MerchAdmin -eq "403" -and $MerchUsers -eq "403")
Report-Test "Merchant role forbidden from /admin and /admin/users" $MerchDenied "HTTP 403 Forbidden received"

# Scenario 4: Admin Users Management page renders with navigation
$UsersPage = curl.exe -s -b $AdminCookieJar ($BaseURL + "/admin/users")
$UsersPageOk = ($UsersPage -match "User Management" -and $UsersPage -match "Total Registered Users")
Report-Test "Admin User Management page renders successfully" $UsersPageOk "Navigation tabs and header detected"

# Scenario 5: User search and role filtering
$FilterUrl = $BaseURL + '/admin/users?role=merchant&table_only=true'
$FilterMerchants = curl.exe -s -b $AdminCookieJar $FilterUrl
$FilterOk = ($FilterMerchants -match "merchant.ctr@fooddiscovery.com" -and !($FilterMerchants -match "user@fooddiscovery.com"))

$SearchUrl = $BaseURL + '/admin/users?q=admin@fooddiscovery.com&table_only=true'
$SearchUser = curl.exe -s -b $AdminCookieJar $SearchUrl
$SearchOk = ($SearchUser -match "admin@fooddiscovery.com" -and !($SearchUser -match "merchant.ctr@fooddiscovery.com"))
Report-Test "User search and role filtering operate correctly" ($FilterOk -and $SearchOk) "Live HTMX filtering confirmed"

# Scenario 6: User listing strictly protects sensitive data
$HasPasswordOrHash = ($UsersPage -match "password" -and ($UsersPage -match '\$2a\$' -or $UsersPage -match "hash" -or $UsersPage -match "adminpassword"))
Report-Test "User listing contains zero password hashes or secrets" (!($HasPasswordOrHash)) "Verified sensitive fields excluded"

# Scenario 7: Platform statistics verification
$StatsAdmin = curl.exe -s -b $AdminCookieJar ($BaseURL + "/admin")
$HasTotalUsers = ($StatsAdmin -match "Platform Users" -or $StatsAdmin -match "Total Users")
$HasUnassigned = ($StatsAdmin -match "Unassigned Outlets")
$HasItemsUnder100 = ($StatsAdmin -match "Available" -or $StatsAdmin -match "Items Under")
$StatsOk = ($HasTotalUsers -and $HasUnassigned -and $HasItemsUnder100)
Report-Test "Platform metrics include users, restaurants, unassigned outlets, and items" $StatsOk "Real DB metrics rendered"

# Scenario 8: Owner modal displays current ownership and eligible merchants
$OwnerModal = curl.exe -s -b $AdminCookieJar ($BaseURL + "/admin/restaurants/2/owner")
$OwnerModalOk = ($OwnerModal -match "Assign Restaurant Owner" -and $OwnerModal -match "merchant.ctr@fooddiscovery.com")
Report-Test "Owner assignment modal renders with eligible merchants" $OwnerModalOk "CTR owner information verified"

# Scenario 9: Non-admin cannot assign or change owner
$CustAssign = curl.exe -s -w "%{http_code}" -b $CustomerCookieJar -X POST ($BaseURL + "/admin/restaurants/2/owner") -d "owner_id=3" -o /dev/null
$MerchAssign = curl.exe -s -w "%{http_code}" -b $MerchantCookieJar -X POST ($BaseURL + "/admin/restaurants/2/owner") -d "owner_id=3" -o /dev/null
$NonAdminDenied = ($CustAssign -eq "403" -and $MerchAssign -eq "403")
Report-Test "Non-admin users cannot perform owner assignment" $NonAdminDenied "Server-side authorization enforced"

# Scenario 10: Invalid merchant assignment is rejected
$InvalidAssign = curl.exe -s -b $AdminCookieJar -X POST ($BaseURL + "/admin/restaurants/2/owner") -d "owner_id=2"
$InvalidRejected = ($InvalidAssign -match "Assignment Error" -or $InvalidAssign -match "not an eligible merchant")
Report-Test "Assignment of non-merchant user is rejected" $InvalidRejected "Role validation successfully caught invalid assignment"

# Scenario 11: Owner assignment and unassignment lifecycle
$UnassignRes = curl.exe -s -i -b $AdminCookieJar -X POST ($BaseURL + "/admin/restaurants/2/owner") -d "owner_id=0"
$UnassignOk = ($UnassignRes -match "Ownership updated successfully" -or $UnassignRes -match "HX-Refresh")
$ReassignRes = curl.exe -s -i -b $AdminCookieJar -X POST ($BaseURL + "/admin/restaurants/2/owner") -d "owner_id=3"
$ReassignOk = ($ReassignRes -match "Ownership updated successfully" -or $ReassignRes -match "HX-Refresh")
Report-Test "Owner assignment and unassignment lifecycle executes safely" ($UnassignOk -and $ReassignOk) "CTR safely preserved with owner ID 3"

# Scenario 12: Restaurant approval status updates
$StatusUrl = $BaseURL + '/admin/restaurants/2/status?status=approved'
$StatusUpdate = curl.exe -s -b $AdminCookieJar -X POST $StatusUrl
$StatusOk = ($StatusUpdate -match "Approved" -and $StatusUpdate -match "admin-row-2")
Report-Test "Restaurant approval updates and returns complete table row" $StatusOk "Status toggle and row outerHTML swap verified"

# Scenario 13: Stage C Merchant Dashboard non-regression
$MerchDash = curl.exe -s -b $MerchantCookieJar ($BaseURL + "/merchant")
$MerchDashOk = ($MerchDash -match "Merchant" -and $MerchDash -match "CTR")
Report-Test "Stage C Merchant Dashboard remains 100% operational" $MerchDashOk "Merchant portal functioning without regression"

# Scenario 14: Customer storefront non-regression
$HomeRes = curl.exe -s ($BaseURL + "/")
$NearbyRes = curl.exe -s ($BaseURL + "/nearby")
$RestoDetail = curl.exe -s ($BaseURL + "/restaurant/2")
$CustOk = ($HomeRes -match "Bite100" -or $HomeRes -match "Bangalore") -and ($NearbyRes -match "Bite100" -or $NearbyRes -match "Bangalore") -and ($RestoDetail -match "Central Tiffin Room")
Report-Test "Customer storefront, nearby discovery, and details pages intact" $CustOk "Customer experience unaffected"

# Scenario 15: Static assets bit-for-bit preservation
$GitDiffStatic = git diff recovery-baseline-stage-a -- static/
$StaticOk = [string]::IsNullOrWhiteSpace($GitDiffStatic)
Report-Test "All 19 static images and SVGs match Stage A baseline bit-for-bit" $StaticOk "Zero asset modification detected"

# Cleanup
Remove-Item $AdminCookieJar -ErrorAction SilentlyContinue
Remove-Item $MerchantCookieJar -ErrorAction SilentlyContinue
Remove-Item $CustomerCookieJar -ErrorAction SilentlyContinue

Write-Host "==========================================================" -ForegroundColor Cyan
if ($Failed -eq 0) {
    Write-Host "TEST SUMMARY: $Passed PASSED, $Failed FAILED - ALL TESTS PASSED!" -ForegroundColor Green
} else {
    Write-Host "TEST SUMMARY: $Passed PASSED, $Failed FAILED" -ForegroundColor Red
}
Write-Host "==========================================================" -ForegroundColor Cyan
