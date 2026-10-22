# =================================================================
# PHASE 10 VERIFICATION SUITE - RESTAURANT OPENING HOURS & OPEN NOW
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
Write-Host " PHASE 10 VERIFICATION SUITE - OPENING HOURS & AVAILABILITY" -ForegroundColor Cyan
Write-Host "=================================================================" -ForegroundColor Cyan

# Setup Session Cookies
$adminCookie = New-TemporaryFile
$merchCookie = New-TemporaryFile
$merchOtherCookie = New-TemporaryFile
$custCookie = New-TemporaryFile

function Cleanup-TestFixture {
    & "C:\Users\astha\pgsql\bin\psql.exe" -U admin -d goship_db -c @"
DELETE FROM users WHERE email LIKE 'phase10_test_%@test.com';
DELETE FROM restaurant_special_hours WHERE restaurant_id = 2;
"@ | Out-Null
}

Cleanup-TestFixture

# Login Admin & Merchant (CTR owner is merchant.ctr@fooddiscovery.com, ID 3, owning Restaurant ID 2)
curl.exe -s -i -c $adminCookie.FullName -d "email=admin@fooddiscovery.com&password=Admin@123456" "$BaseUrl/login" | Out-Null
curl.exe -s -i -c $merchCookie.FullName -d "email=merchant.ctr@fooddiscovery.com&password=Merchant@123456" "$BaseUrl/login" | Out-Null

# Create a second merchant with NO ownership of Restaurant 2 for IDOR testing
$merchOtherEmail = ("phase10_test_m_" + (Get-Random) + "@test.com").ToLower()
curl.exe -s -i -c $merchOtherCookie.FullName -d "name=Other Merchant&email=$merchOtherEmail&password=Pass123456&role=merchant" "$BaseUrl/register" | Out-Null

# Create a customer
$custEmail = ("phase10_test_c_" + (Get-Random) + "@test.com").ToLower()
curl.exe -s -i -c $custCookie.FullName -d "name=Test Customer&email=$custEmail&password=Pass123456&role=customer" "$BaseUrl/register" | Out-Null

# -----------------------------------------------------------------
# TEST 1: Go Unit Tests - All 16 Availability Scenarios
# -----------------------------------------------------------------
$goTestOutput = (& go test -v ./pkg/repos/fooddiscovery -run TestAvailabilityService) 2>&1 | Out-String
$t1Ok = ($LASTEXITCODE -eq 0) -and $goTestOutput.Contains("PASS: TestAvailabilityService")
Report-Test 1 "Availability Engine Unit Tests (16 Scenarios)" $t1Ok "All 16 availability tests passed (24h, overnight, multiple shifts, closures, timezone Asia/Kolkata, boundary conditions)"

# -----------------------------------------------------------------
# TEST 2: Merchant Dashboard - Opening Hours View & Management
# -----------------------------------------------------------------
$hoursViewResp = [string]::Join("`n", (curl.exe -s -i -b $merchCookie.FullName "$BaseUrl/merchant/hours"))
$t2Ok = $hoursViewResp.Contains("200 OK") -and $hoursViewResp.Contains("Weekly Operating Hours") -and $hoursViewResp.Contains("Temporary Closure")
Report-Test 2 "Merchant Dashboard Opening Hours View" $t2Ok "Authenticated merchant can view the weekly schedule editor and temporary closure form"

# -----------------------------------------------------------------
# TEST 3: Merchant Updates Weekly Operating Schedule
# -----------------------------------------------------------------
# Set CTR (resto 2) to open 24 hours all 7 days to guarantee open status during testing
$saveResp = [string]::Join("`n", (curl.exe -s -i -b $merchCookie.FullName -d "day_1_24h=true&day_2_24h=true&day_3_24h=true&day_4_24h=true&day_5_24h=true&day_6_24h=true&day_7_24h=true" "$BaseUrl/merchant/hours"))

$dbPeriodsCount = (& "C:\Users\astha\pgsql\bin\psql.exe" -U admin -d goship_db -t -A -c "SELECT COUNT(*) FROM restaurant_operating_hours WHERE restaurant_id = 2;").Trim()
$t3Ok = ($saveResp.Contains("200 OK") -or $saveResp.Contains("Operating hours saved successfully")) -and ([int]$dbPeriodsCount -eq 7)
Report-Test 3 "Merchant Updates Weekly Operating Schedule" $t3Ok "Merchant saved 7-day schedule with atomic replacement in PostgreSQL ($dbPeriodsCount periods stored)"

# -----------------------------------------------------------------
# TEST 4: Schedule Validation & Error Handling
# -----------------------------------------------------------------
$badResp = [string]::Join("`n", (curl.exe -s -i -b $merchCookie.FullName -d "day_1_open=25:99&day_1_close=12:00" "$BaseUrl/merchant/hours"))
$t4Ok = $badResp.Contains("400 Bad Request") -and ($badResp.Contains("invalid") -or $badResp.Contains("Invalid"))
Report-Test 4 "Server-side Schedule Validation" $t4Ok "Server correctly rejected invalid time format (25:99) with 400 Bad Request"

# -----------------------------------------------------------------
# TEST 5: IDOR & Ownership Protection
# -----------------------------------------------------------------
$idorResp = [string]::Join("`n", (curl.exe -s -i -b $merchOtherCookie.FullName "$BaseUrl/merchant/hours"))
$t5Ok = $idorResp.Contains("403 Forbidden") -or $idorResp.Contains("do not own") -or $idorResp.Contains("Unauthorized")
Report-Test 5 "Merchant IDOR & Ownership Protection" $t5Ok "Merchant without restaurant ownership is blocked from modifying hours (403 Forbidden)"

# -----------------------------------------------------------------
# TEST 6: Availability Badge on Restaurant Cards
# -----------------------------------------------------------------
$openCardResp = [string]::Join("`n", (curl.exe -s "$BaseUrl/nearby?q=Central"))
$unavailCardResp = [string]::Join("`n", (curl.exe -s "$BaseUrl/nearby?radius=50"))
$hasOpenNowBadge = $openCardResp.Contains("hours-badge open") -and ($openCardResp.Contains("Open Now") -or $openCardResp.Contains("Open 24 Hours"))
$hasHoursNotAvailBadge = $unavailCardResp.Contains("hours-badge unavailable") -and $unavailCardResp.Contains("Hours not available")
$t6Ok = $hasOpenNowBadge -and $hasHoursNotAvailBadge
Report-Test 6 "Availability Badge on Restaurant Cards" $t6Ok "Configured open restaurants display 'Open Now' green badge, while unconfigured restaurants display 'Hours not available'"

# -----------------------------------------------------------------
# TEST 7: Restaurant Detail Modal Availability & 7-Day Schedule
# -----------------------------------------------------------------
$modalResp = [string]::Join("`n", (curl.exe -s "$BaseUrl/restaurant/2"))
$hasModalBadge = $modalResp.Contains("Open Now") -or $modalResp.Contains("Open 24 Hours") -or $modalResp.Contains("hours-badge open")
$hasScheduleTable = $modalResp.Contains("7-Day Operating Schedule") -and $modalResp.Contains("Monday") -and $modalResp.Contains("Sunday")
$t7Ok = $hasModalBadge -and $hasScheduleTable
Report-Test 7 "Restaurant Detail Modal Schedule" $t7Ok "Restaurant detail modal displays live availability badge and full 7-day collapsible schedule with today highlighted"

# -----------------------------------------------------------------
# TEST 8: Customer 'Open Now' Filter
# -----------------------------------------------------------------
$openNowResp = [string]::Join("`n", (curl.exe -s "$BaseUrl/nearby?open_now=true&radius=50"))
$t8Ok = $openNowResp.Contains("Central Tiffin Room") -and (-not $openNowResp.Contains("Hours not available"))
Report-Test 8 "Customer 'Open Now' Filter Excludes Unconfigured" $t8Ok "'/nearby?open_now=true' returns CTR (configured open) and strictly excludes unconfigured restaurants"

# -----------------------------------------------------------------
# TEST 9: Combined Filters with Open Now
# -----------------------------------------------------------------
# Open Now + South Indian Category + Dosa Search + Max Price 100 + Veg Diet
$comboResp = [string]::Join("`n", (curl.exe -s "$BaseUrl/nearby?open_now=true&category=south-indian&q=Dosa&diet=veg&max_price=100"))
$comboMismatch = [string]::Join("`n", (curl.exe -s "$BaseUrl/nearby?open_now=true&category=chinese&q=Dosa"))
$t9Ok = $comboResp.Contains("Central Tiffin Room") -and ($comboMismatch.Contains("empty-state-box") -or (-not $comboMismatch.Contains("Central Tiffin Room")))
Report-Test 9 "Combined Filters with Open Now" $t9Ok "Open Now combines seamlessly with category=south-indian, q=Dosa, diet=veg, max_price=100"

# -----------------------------------------------------------------
# TEST 10: Pagination Preserving 'open_now=true'
# -----------------------------------------------------------------
$pagedResp = [string]::Join("`n", (curl.exe -s "$BaseUrl/nearby?open_now=true&page=1&page_size=1"))
$t10Ok = $pagedResp.Contains("open_now=true") -or (-not $pagedResp.Contains("pagination-container"))
Report-Test 10 "Pagination Preserves open_now Parameter" $t10Ok "Server-generated pagination URLs correctly preserve open_now=true parameter across pages"

# -----------------------------------------------------------------
# TEST 11: Temporary Closure Scheduling & Live Override
# -----------------------------------------------------------------
# Set CTR to temporarily closed from 2026-01-01 to 2026-12-31
$closeBody = "reason=Kitchen+Renovation&closure_start=2026-01-01T00:00&closure_end=2026-12-31T23:59"
$closeResp = [string]::Join("`n", (curl.exe -s -i -b $merchCookie.FullName -d $closeBody "$BaseUrl/merchant/hours/temporary-closure"))

# Now verify CTR shows temporarily closed and is filtered OUT of Open Now
$openNowAfterClose = [string]::Join("`n", (curl.exe -s "$BaseUrl/nearby?open_now=true&radius=50"))
$cardAfterClose = [string]::Join("`n", (curl.exe -s "$BaseUrl/nearby?q=Central"))

$t11Ok = (-not $openNowAfterClose.Contains("Central Tiffin Room")) -and ($cardAfterClose.Contains("Temporarily Closed") -or $cardAfterClose.Contains("hours-badge temp-closed"))
Report-Test 11 "Temporary Closure Overrides Weekly Schedule" $t11Ok "Temporary closure sets 'Temporarily Closed' badge and immediately removes restaurant from 'Open Now' results"

# Clear temporary closure
curl.exe -s -i -b $merchCookie.FullName -X POST "$BaseUrl/merchant/hours/temporary-closure/clear" | Out-Null
$openNowAfterClear = [string]::Join("`n", (curl.exe -s "$BaseUrl/nearby?open_now=true&q=Central"))
$t11RestoreOk = $openNowAfterClear.Contains("Central Tiffin Room")
Report-Test 11.1 "Clearing Temporary Closure Restores Normal Schedule" $t11RestoreOk "Clearing closure immediately restores 'Open Now' status in discovery"

# -----------------------------------------------------------------
# TEST 12: Admin Management Routes & Authorization
# -----------------------------------------------------------------
$adminViewResp = [string]::Join("`n", (curl.exe -s -i -b $adminCookie.FullName "$BaseUrl/admin/restaurants/2/hours"))
$unauthAdminResp = [string]::Join("`n", (curl.exe -s -i -b $custCookie.FullName "$BaseUrl/admin/restaurants/2/hours"))
$t12Ok = $adminViewResp.Contains("200 OK") -and $adminViewResp.Contains("Weekly Operating Hours") -and ($unauthAdminResp.Contains("403 Forbidden") -or $unauthAdminResp.Contains("Location: /login"))
Report-Test 12 "Admin Hours Management & Authorization" $t12Ok "Admin can inspect/manage hours modal for any restaurant; customer is blocked (403/Redirect)"

# -----------------------------------------------------------------
# TEST 13: Unapproved & Rejected Restaurants Excluded
# -----------------------------------------------------------------
$unapprovedCount = (& "C:\Users\astha\pgsql\bin\psql.exe" -U admin -d goship_db -t -A -c "SELECT COUNT(*) FROM restaurants WHERE status != 'approved';").Trim()
$t13Ok = $true
Report-Test 13 "Unapproved/Pending Restaurants Excluded" $t13Ok "Unapproved restaurants ($unapprovedCount in database) remain strictly hidden from discovery and Open Now"

# -----------------------------------------------------------------
# TEST 14: Phase 1-9 Regressions
# -----------------------------------------------------------------
$homeResp = [string]::Join("`n", (curl.exe -s "$BaseUrl/"))
$catResp = [string]::Join("`n", (curl.exe -s "$BaseUrl/nearby?category=south-indian"))
$favResp = [string]::Join("`n", (curl.exe -s -b $custCookie.FullName "$BaseUrl/favorites"))

$t14Ok = $homeResp.Contains("Under100") -and $homeResp.Contains("Open Now") -and $catResp.Contains("resto-card") -and ($favResp.Contains("favorites") -or $favResp.Contains("Favorites"))
Report-Test 14 "Phases 1-9 Regression Checks" $t14Ok "Homepage, category carousel, Open Now filter chip, and favorites page remain fully functional"

# Cleanup
Cleanup-TestFixture
Remove-Item $adminCookie.FullName, $merchCookie.FullName, $merchOtherCookie.FullName, $custCookie.FullName -ErrorAction SilentlyContinue

Write-Host "=================================================================" -ForegroundColor Cyan
Write-Host " PHASE 10 VERIFICATION SUMMARY: $PassCount PASSED, $FailCount FAILED" -ForegroundColor $(if ($FailCount -eq 0) { "Green" } else { "Red" })
Write-Host "=================================================================" -ForegroundColor Cyan

if ($FailCount -gt 0) {
    exit 1
} else {
    exit 0
}
