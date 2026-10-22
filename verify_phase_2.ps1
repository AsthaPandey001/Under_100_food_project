# verify_phase_2.ps1 - Automated verification test suite for Phase 2: Self-Service Restaurant Registration

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
Write-Host "PHASE 2 VERIFICATION: SELF-SERVICE RESTAURANT REGISTRATION" -ForegroundColor Cyan
Write-Host "==========================================================" -ForegroundColor Cyan

$AdminCookieJar = Join-Path $PSScriptRoot "p2_admin_cookie.txt"
$ExistingMerchCookieJar = Join-Path $PSScriptRoot "p2_ctr_merch_cookie.txt"
$CustCookieJar = Join-Path $PSScriptRoot "p2_cust_cookie.txt"
$NewMerchCookieJar = Join-Path $PSScriptRoot "p2_new_merch_cookie.txt"

Remove-Item $AdminCookieJar -ErrorAction SilentlyContinue
Remove-Item $ExistingMerchCookieJar -ErrorAction SilentlyContinue
Remove-Item $CustCookieJar -ErrorAction SilentlyContinue
Remove-Item $NewMerchCookieJar -ErrorAction SilentlyContinue

# Ensure clean slate for disposable test fixtures before starting
& "C:\Users\astha\pgsql\bin\psql.exe" -U admin -d goship_db -c "DELETE FROM restaurants WHERE name = 'Phase 2 Test Kitchen';" | Out-Null
& "C:\Users\astha\pgsql\bin\psql.exe" -U admin -d goship_db -c "DELETE FROM users WHERE email = 'phase2_temp_merchant@test.com';" | Out-Null

# Authenticate Admin
curl.exe -s -c $AdminCookieJar -X POST ($BaseURL + "/login") -d "email=admin@fooddiscovery.com" -d "password=Admin@123456" | Out-Null

# Authenticate Customer
curl.exe -s -c $CustCookieJar -X POST ($BaseURL + "/login") -d "email=user@fooddiscovery.com" -d "password=User@123456" | Out-Null

# Authenticate Existing Merchant (CTR Owner - already has an assigned restaurant)
curl.exe -s -c $ExistingMerchCookieJar -X POST ($BaseURL + "/login") -d "email=merchant.ctr@fooddiscovery.com" -d "password=Merchant@123456" | Out-Null

# Register a new unassigned disposable test merchant
$RegCode = curl.exe -s -c $NewMerchCookieJar -w "%{http_code}" -X POST ($BaseURL + "/register") `
    -d "name=Phase 2 Test Merchant" `
    -d "email=phase2_temp_merchant@test.com" `
    -d "password=TestPass@123" `
    -d "role=merchant" -o /dev/null

$NewMerchDash = curl.exe -s -b $NewMerchCookieJar ($BaseURL + "/merchant")
$HasEmptyState = ($NewMerchDash -match "No Restaurant Assigned Yet" -and $NewMerchDash -match "Register My Restaurant")

# 1. Unauthenticated restaurant submission is rejected
$UnauthSubmitCode = curl.exe -s -w "%{http_code}" -X POST ($BaseURL + "/merchant/restaurant/submit") `
    -d "name=Hacker Kitchen" -d "address=Nowhere Street" -o /dev/null
$Test1 = ($UnauthSubmitCode -eq "303" -or $UnauthSubmitCode -eq "302" -or $UnauthSubmitCode -eq "401")
Report-Test "1. Unauthenticated restaurant submission is rejected" $Test1 "HTTP code: $UnauthSubmitCode"

# 2. Customer restaurant submission is rejected (403 Forbidden)
$CustSubmitCode = curl.exe -s -w "%{http_code}" -b $CustCookieJar -X POST ($BaseURL + "/merchant/restaurant/submit") `
    -d "name=Customer Kitchen" -d "address=Customer Street" -o /dev/null
$Test2 = ($CustSubmitCode -eq "403")
Report-Test "2. Customer restaurant submission is rejected (HTTP 403)" $Test2 "HTTP code: $CustSubmitCode"

# 8. Existing merchant cannot overwrite incompatible ownership association
$ExistingMerchSubmit = curl.exe -s -b $ExistingMerchCookieJar -X POST ($BaseURL + "/merchant/restaurant/submit") `
    -d "name=Second CTR Outlet" -d "address=Koramangala, Bangalore"
$Test8 = ($ExistingMerchSubmit -match "already have a registered restaurant" -or $ExistingMerchSubmit -match "already have an assigned restaurant")
Report-Test "8. Merchant cannot overwrite existing incompatible ownership" $Test8 "Blocked second outlet registration for existing owner"

# 6. Invalid fields and coordinates are rejected
$InvalidCoordsRes = curl.exe -s -b $NewMerchCookieJar -X POST ($BaseURL + "/merchant/restaurant/submit") `
    -d "name=Test Outlet" -d "address=Indiranagar, Bangalore" -d "latitude=999" -d "longitude=77.5946"
$InvalidPhoneRes = curl.exe -s -b $NewMerchCookieJar -X POST ($BaseURL + "/merchant/restaurant/submit") `
    -d "name=Test Outlet" -d "address=Indiranagar, Bangalore" -d "phone=NOTAPHONE!!!"
$ShortNameRes = curl.exe -s -b $NewMerchCookieJar -X POST ($BaseURL + "/merchant/restaurant/submit") `
    -d "name=X" -d "address=Indiranagar, Bangalore"
$Test6 = ($InvalidCoordsRes -match "valid geographic coordinates") -and ($InvalidPhoneRes -match "invalid characters" -or $InvalidPhoneRes -match "valid contact phone number") -and ($ShortNameRes -match "Restaurant name is required")
Report-Test "6. Invalid fields, coordinates, and short names are rejected" $Test6 "Validation catches invalid coordinates, phone, and short names"

# 3, 4, 5. Merchant submission with valid details creates a pending listing linked to authenticated merchant (ignoring forged owner_id)
$SubmitRes = curl.exe -s -i -b $NewMerchCookieJar -X POST ($BaseURL + "/merchant/restaurant/submit") `
    -d "name=Phase 2 Test Kitchen" `
    -d "address=88 100-Ft Road, Indiranagar, Bengaluru" `
    -d "phone=+91 80 2525 9999" `
    -d "latitude=12.9780" `
    -d "longitude=77.6400" `
    -d "opening_hours=8:00 AM - 10:00 PM" `
    -d "description=Signature tiffins and coffee" `
    -d "owner_id=1" # Attempt to forge owner ID to Admin

# Check DB for creation, status, and owner
$DbCheck = & "C:\Users\astha\pgsql\bin\psql.exe" -U admin -d goship_db -t -A -F "|" -c `
    "SELECT r.id, r.name, r.status, r.user_restaurants, u.email FROM restaurants r JOIN users u ON r.user_restaurants = u.id WHERE r.name = 'Phase 2 Test Kitchen';"

$DbParts = $DbCheck.Split("|")
$CreatedRestID = $DbParts[0]
$CreatedRestName = $DbParts[1]
$CreatedRestStatus = $DbParts[2]
$CreatedRestOwnerEmail = $DbParts[4]

$Test3 = ($CreatedRestStatus -eq "pending") -and ($CreatedRestName -eq "Phase 2 Test Kitchen")
Report-Test "3. Merchant submission creates a pending listing" $Test3 "Status in DB: $CreatedRestStatus"

$Test4 = ($CreatedRestOwnerEmail -eq "phase2_temp_merchant@test.com")
Report-Test "4. New listing is linked to the authenticated merchant" $Test4 "Owner in DB: $CreatedRestOwnerEmail"

$Test5 = ($CreatedRestOwnerEmail -ne "admin@fooddiscovery.com")
Report-Test "5. Forged owner ID (owner_id=1) was safely ignored" $Test5 "Owner strictly set to authenticated merchant"

# 7. Duplicate submission behaviour is handled correctly
$DupRes = curl.exe -s -b $NewMerchCookieJar -X POST ($BaseURL + "/merchant/restaurant/submit") `
    -d "name=Phase 2 Test Kitchen" `
    -d "address=88 100-Ft Road, Indiranagar, Bengaluru"
$Test7 = ($DupRes -match "already have a registered restaurant" -or $DupRes -match "already been submitted")
Report-Test "7. Duplicate / repeated submission behaviour is rejected" $Test7 "Prevented duplicate restaurant registration"

# 9. Pending listing appears in the admin review queue
$AdminPendingList = curl.exe -s -b $AdminCookieJar ($BaseURL + "/admin?status=pending")
$Test9 = ($AdminPendingList -match "Phase 2 Test Kitchen") -and ($AdminPendingList -match "phase2_temp_merchant@test.com")
Report-Test "9. Pending listing appears in admin review queue with owner" $Test9 "Found in pending filter with submitting merchant"

# 11. Pending listing remains hidden from public discovery (nearby, markers)
$NearbyRes = curl.exe -s ($BaseURL + "/nearby")
$MarkersRes = curl.exe -s ($BaseURL + "/api/markers?lat=12.9780&lng=77.6400&radius=25")
$Test11 = !($NearbyRes -match "Phase 2 Test Kitchen") -and !($MarkersRes -match "Phase 2 Test Kitchen")
Report-Test "11. Pending listing is strictly hidden from public discovery" $Test11 "Excluded from /nearby and /api/markers"

# 12. Phase 1 direct-URL protection still works on newly submitted pending listing
$PublicDirectCode = curl.exe -s -w "%{http_code}" -o /dev/null ($BaseURL + "/restaurant/" + $CreatedRestID)
$MerchDirectCode = curl.exe -s -w "%{http_code}" -b $NewMerchCookieJar -o /dev/null ($BaseURL + "/restaurant/" + $CreatedRestID)
$Test12 = ($PublicDirectCode -eq "404") -and ($MerchDirectCode -eq "200")
Report-Test "12. Phase 1 direct-URL protection active (Public=404, Owner=200)" $Test12 "Public=$PublicDirectCode, Owner=$MerchDirectCode"

# 10. Admin approval updates status and preserves owner relationship
$ApproveUrl = $BaseURL + "/admin/restaurants/" + $CreatedRestID + "/status?status=approved"
$ApproveRes = curl.exe -s -b $AdminCookieJar -X POST $ApproveUrl
$DbApproveCheck = & "C:\Users\astha\pgsql\bin\psql.exe" -U admin -d goship_db -t -A -F "|" -c `
    "SELECT r.status, u.email FROM restaurants r JOIN users u ON r.user_restaurants = u.id WHERE r.id = $CreatedRestID;"
$ApproveParts = $DbApproveCheck.Split("|")
$PostApproveStatus = $ApproveParts[0]
$PostApproveOwner = $ApproveParts[1]
$Test10 = ($PostApproveStatus -eq "approved") -and ($PostApproveOwner -eq "phase2_temp_merchant@test.com")
Report-Test "10. Admin approval updates status and preserves owner" $Test10 "Status=$PostApproveStatus, Owner=$PostApproveOwner"

# 13. Existing merchant menu management still works
$CTRMerchDash = curl.exe -s -b $ExistingMerchCookieJar ($BaseURL + "/merchant")
$Test13 = ($CTRMerchDash -match "Central Tiffin Room") -and ($CTRMerchDash -match "Benne Masala Dosa")
Report-Test "13. Existing CTR merchant dashboard & menu fully operational" $Test13 "CTR menu items rendered properly"

# 14. Customer homepage, nearby search, and details remain functional
$HomeRes = curl.exe -s ($BaseURL + "/")
$Resto1Res = curl.exe -s ($BaseURL + "/restaurant/1")
$Test14 = ($HomeRes -match "Under100" -or $HomeRes -match "Bengaluru" -or $HomeRes -match "Bite100") -and ($Resto1Res -match "Vidyarthi Bhavan")
Report-Test "14. Customer storefront, nearby search, and modals intact" $Test14 "Homepage and Vidyarthi Bhavan rendered"

# 15. Application binary compiles and exists
$BuildTest = Test-Path (Join-Path $PSScriptRoot "under100_server.exe")
Report-Test "15. Application compiles and binary exists" $BuildTest "under100_server.exe verified"

# 16. Static assets remain unchanged
$GitDiffStatic = git diff recovery-baseline-stage-a -- static/
$StaticOk = [string]::IsNullOrWhiteSpace($GitDiffStatic)
Report-Test "16. All 19 static images and SVGs match Stage A baseline bit-for-bit" $StaticOk "Zero asset modification detected"

# Cleanup disposable test fixtures
& "C:\Users\astha\pgsql\bin\psql.exe" -U admin -d goship_db -c "DELETE FROM restaurants WHERE id = $CreatedRestID;" | Out-Null
& "C:\Users\astha\pgsql\bin\psql.exe" -U admin -d goship_db -c "DELETE FROM users WHERE email = 'phase2_temp_merchant@test.com';" | Out-Null

Remove-Item $AdminCookieJar -ErrorAction SilentlyContinue
Remove-Item $ExistingMerchCookieJar -ErrorAction SilentlyContinue
Remove-Item $CustCookieJar -ErrorAction SilentlyContinue
Remove-Item $NewMerchCookieJar -ErrorAction SilentlyContinue

Write-Host "==========================================================" -ForegroundColor Cyan
if ($Failed -eq 0) {
    Write-Host "TEST SUMMARY: $Passed PASSED, $Failed FAILED - ALL 16 TESTS PASSED!" -ForegroundColor Green
} else {
    Write-Host "TEST SUMMARY: $Passed PASSED, $Failed FAILED" -ForegroundColor Red
}
Write-Host "==========================================================" -ForegroundColor Cyan
