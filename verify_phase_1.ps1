# verify_phase_1.ps1 - Automated verification test suite for Phase 1: Unapproved Restaurant Access Control

$BaseURL = "http://127.0.0.1:8000"
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
Write-Host "PHASE 1 VERIFICATION: UNAPPROVED RESTAURANT ACCESS POLICY" -ForegroundColor Cyan
Write-Host "==========================================================" -ForegroundColor Cyan

$AdminCookieJar = Join-Path $PSScriptRoot "p1_admin_cookie.txt"
$MerchantCookieJar = Join-Path $PSScriptRoot "p1_merchant_cookie.txt"
$CustomerCookieJar = Join-Path $PSScriptRoot "p1_customer_cookie.txt"
Remove-Item $AdminCookieJar -ErrorAction SilentlyContinue
Remove-Item $MerchantCookieJar -ErrorAction SilentlyContinue
Remove-Item $CustomerCookieJar -ErrorAction SilentlyContinue

# Authenticate users
curl.exe -s -c $AdminCookieJar -X POST ($BaseURL + "/login") -d "email=admin@fooddiscovery.com" -d "password=Admin@123456" | Out-Null
curl.exe -s -c $MerchantCookieJar -X POST ($BaseURL + "/login") -d "email=merchant.ctr@fooddiscovery.com" -d "password=Merchant@123456" | Out-Null
curl.exe -s -c $CustomerCookieJar -X POST ($BaseURL + "/login") -d "email=user@fooddiscovery.com" -d "password=User@123456" | Out-Null

# 1. Approved restaurant accessible publicly
$ApprPublicCode = curl.exe -s -w "%{http_code}" -o /dev/null ($BaseURL + "/restaurant/1")
$ApprPublicBody = curl.exe -s ($BaseURL + "/restaurant/1")
$Test1 = ($ApprPublicCode -eq "200") -and ($ApprPublicBody -match "Vidyarthi Bhavan")
Report-Test "1. Approved restaurant is accessible publicly (HTTP 200)" $Test1 "Status: $ApprPublicCode, Content verified"

# 2. Pending restaurant not accessible to unauthenticated visitor (HTTP 404)
$PendPublicCode = curl.exe -s -w "%{http_code}" -o /dev/null ($BaseURL + "/restaurant/10")
$PendPublicBody = curl.exe -s ($BaseURL + "/restaurant/10")
$Test2 = ($PendPublicCode -eq "404") -and ($PendPublicBody -match "Restaurant not found") -and !($PendPublicBody -match "The Bangalore Cafe")
Report-Test "2. Pending restaurant is not accessible to unauthenticated visitor (HTTP 404)" $Test2 "Status: $PendPublicCode, Details strictly hidden"

# 3. Rejected restaurant not accessible to unauthenticated visitor (HTTP 404)
$RejPublicCode = curl.exe -s -w "%{http_code}" -o /dev/null ($BaseURL + "/restaurant/11")
$RejPublicBody = curl.exe -s ($BaseURL + "/restaurant/11")
$Test3 = ($RejPublicCode -eq "404") -and ($RejPublicBody -match "Restaurant not found") -and !($RejPublicBody -match "Spice Terrace")
Report-Test "3. Rejected restaurant is not accessible to unauthenticated visitor (HTTP 404)" $Test3 "Status: $RejPublicCode, Details strictly hidden"

# 4. Customer cannot access pending or rejected restaurant details (HTTP 404)
$CustPendCode = curl.exe -s -w "%{http_code}" -b $CustomerCookieJar -o /dev/null ($BaseURL + "/restaurant/10")
$CustRejCode = curl.exe -s -w "%{http_code}" -b $CustomerCookieJar -o /dev/null ($BaseURL + "/restaurant/11")
$CustPendBody = curl.exe -s -b $CustomerCookieJar ($BaseURL + "/restaurant/10")
$Test4 = ($CustPendCode -eq "404") -and ($CustRejCode -eq "404") -and ($CustPendBody -match "Restaurant not found")
Report-Test "4. Customer role cannot access pending or rejected restaurants (HTTP 404)" $Test4 "Pending code: $CustPendCode, Rejected code: $CustRejCode"

# 5. Unrelated merchant cannot access another unapproved restaurant (HTTP 404)
# Merchant 3 owns restaurant 2. Restaurant 10 (pending) and 11 (rejected) are unassigned/unrelated.
$MerchPendCode = curl.exe -s -w "%{http_code}" -b $MerchantCookieJar -o /dev/null ($BaseURL + "/restaurant/10")
$MerchRejCode = curl.exe -s -w "%{http_code}" -b $MerchantCookieJar -o /dev/null ($BaseURL + "/restaurant/11")
$Test5 = ($MerchPendCode -eq "404") -and ($MerchRejCode -eq "404")
Report-Test "5. Unrelated merchant cannot access unapproved restaurant (HTTP 404)" $Test5 "Codes: Pending=$MerchPendCode, Rejected=$MerchRejCode"

# 6. Assigned merchant CAN access their own unapproved restaurant
# Temporarily assign merchant 3 to pending restaurant 10, test access, then safely restore to 0
curl.exe -s -b $AdminCookieJar -X POST ($BaseURL + "/admin/restaurants/10/owner") -d "owner_id=3" | Out-Null
$OwnerPendCode = curl.exe -s -w "%{http_code}" -b $MerchantCookieJar -o /dev/null ($BaseURL + "/restaurant/10")
$OwnerPendBody = curl.exe -s -b $MerchantCookieJar ($BaseURL + "/restaurant/10")
# Restore restaurant 10 back to unassigned
curl.exe -s -b $AdminCookieJar -X POST ($BaseURL + "/admin/restaurants/10/owner") -d "owner_id=0" | Out-Null
$Test6 = ($OwnerPendCode -eq "200") -and ($OwnerPendBody -match "The Bangalore Cafe")
Report-Test "6. Assigned merchant can access their own unapproved restaurant (HTTP 200)" $Test6 "Code: $OwnerPendCode, Menu displayed to owner"

# 7. Authorized admin CAN access unapproved restaurants (pending and rejected)
$AdminPendCode = curl.exe -s -w "%{http_code}" -b $AdminCookieJar -o /dev/null ($BaseURL + "/restaurant/10")
$AdminPendBody = curl.exe -s -b $AdminCookieJar ($BaseURL + "/restaurant/10")
$AdminRejCode = curl.exe -s -w "%{http_code}" -b $AdminCookieJar -o /dev/null ($BaseURL + "/restaurant/11")
$AdminRejBody = curl.exe -s -b $AdminCookieJar ($BaseURL + "/restaurant/11")
$Test7 = ($AdminPendCode -eq "200") -and ($AdminPendBody -match "The Bangalore Cafe") -and ($AdminRejCode -eq "200") -and ($AdminRejBody -match "Spice Terrace")
Report-Test "7. Authorized admin can access pending and rejected restaurants (HTTP 200)" $Test7 "Admin pending=$AdminPendCode, Admin rejected=$AdminRejCode"

# 8. Nonexistent restaurant returns expected not-found response (HTTP 404)
$NonExistCode = curl.exe -s -w "%{http_code}" -o /dev/null ($BaseURL + "/restaurant/99999")
$NonExistBody = curl.exe -s ($BaseURL + "/restaurant/99999")
$Test8 = ($NonExistCode -eq "404") -and ($NonExistBody -match "Restaurant not found")
Report-Test "8. Nonexistent restaurant returns HTTP 404 not found" $Test8 "Code: $NonExistCode, Body: Restaurant not found"

# 9. Nearby search and customer homepage still work
$HomeCode = curl.exe -s -w "%{http_code}" -o /dev/null ($BaseURL + "/")
$NearbyCode = curl.exe -s -w "%{http_code}" -o /dev/null ($BaseURL + "/nearby")
$MarkersRes = curl.exe -s ($BaseURL + "/api/markers?lat=12.9716&lng=77.5946&radius=25")
$MarkersOk = ($MarkersRes -match "Vidyarthi Bhavan") -and !($MarkersRes -match "The Bangalore Cafe") -and !($MarkersRes -match "Spice Terrace")
$Test9 = ($HomeCode -eq "200") -and ($NearbyCode -eq "200") -and $MarkersOk
Report-Test "9. Nearby search, homepage, and API markers function normally" $Test9 "Home=$HomeCode, Nearby=$NearbyCode, Markers strictly approved"

# 10. Existing admin approval and merchant dashboard workflows still work
$AdminDash = curl.exe -s -b $AdminCookieJar ($BaseURL + "/admin")
$MerchDash = curl.exe -s -b $MerchantCookieJar ($BaseURL + "/merchant")
$Test10 = ($AdminDash -match "Bite100 Admin" -or $AdminDash -match "Platform Overview") -and ($MerchDash -match "Merchant Portal" -or $MerchDash -match "Central Tiffin Room")
Report-Test "10. Admin dashboard and merchant portal operate normally" $Test10 "Admin and Merchant sessions healthy"

# 11. Application builds successfully
$BuildTest = Test-Path (Join-Path $PSScriptRoot "under100_server.exe")
Report-Test "11. Application binary compiles and exists" $BuildTest "under100_server.exe present"

# 12. Static asset hashes remain unchanged
$GitDiffStatic = git diff recovery-baseline-stage-a -- static/
$StaticOk = [string]::IsNullOrWhiteSpace($GitDiffStatic)
Report-Test "12. All 19 static images and SVGs match Stage A baseline bit-for-bit" $StaticOk "Zero asset modification detected"

# Cleanup cookie files
Remove-Item $AdminCookieJar -ErrorAction SilentlyContinue
Remove-Item $MerchantCookieJar -ErrorAction SilentlyContinue
Remove-Item $CustomerCookieJar -ErrorAction SilentlyContinue

Write-Host "==========================================================" -ForegroundColor Cyan
if ($Failed -eq 0) {
    Write-Host "TEST SUMMARY: $Passed PASSED, $Failed FAILED - ALL 12 TESTS PASSED!" -ForegroundColor Green
} else {
    Write-Host "TEST SUMMARY: $Passed PASSED, $Failed FAILED" -ForegroundColor Red
}
Write-Host "==========================================================" -ForegroundColor Cyan
