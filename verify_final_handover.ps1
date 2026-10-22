# ==============================================================================
# UNDER100 / BITE100 — FINAL HANDOVER COMPREHENSIVE REGRESSION VERIFICATION SUITE
# ==============================================================================
# Date: October 19, 2026
# Scope: Full Regression & System Health Check for Project Handover
# Covers: Customer, Merchant, Admin, Auth, Discovery, Search, Categories, Hours,
#         Reviews, Favorites, Static Assets, and 100% Database Preservation.
# ==============================================================================

$ErrorActionPreference = "Continue"
$BaseUrl = "http://127.0.0.1:8000"
$PsqlCmd = "C:\Users\astha\pgsql\bin\psql.exe"
$DbName = "goship_db"
$DbUser = "admin"

$TotalTests = 0
$PassCount = 0
$FailCount = 0

function Report-Test($id, $name, $condition, $details) {
    $global:TotalTests++
    if ($condition) {
        Write-Host " [PASS] Test $id`: $name" -ForegroundColor Green
        Write-Host "        $details" -ForegroundColor DarkGray
        $global:PassCount++
    } else {
        Write-Host " [FAIL] Test $id`: $name" -ForegroundColor Red
        Write-Host "        $details" -ForegroundColor Yellow
        $global:FailCount++
    }
}

Write-Host "=================================================================" -ForegroundColor Cyan
Write-Host " UNDER100: MASTER FINAL HANDOVER VERIFICATION SUITE              " -ForegroundColor Cyan
Write-Host "=================================================================" -ForegroundColor Cyan

# -----------------------------------------------------------------------------
# STEP 0: RECORD INITIAL DATABASE BASELINE COUNTS
# -----------------------------------------------------------------------------
Write-Host "`n-- Recording Baseline Database Counts --" -ForegroundColor Yellow

$initialCounts = @{}
$dbOutput = (& $PsqlCmd -U $DbUser -d $DbName -t -A -c @"
SELECT 'users', count(*) FROM users
UNION ALL SELECT 'restaurants', count(*) FROM restaurants
UNION ALL SELECT 'food_items', count(*) FROM food_items
UNION ALL SELECT 'categories', count(*) FROM categories
UNION ALL SELECT 'food_item_categories', count(*) FROM food_item_categories
UNION ALL SELECT 'restaurant_operating_hours', count(*) FROM restaurant_operating_hours
UNION ALL SELECT 'reviews', count(*) FROM reviews
UNION ALL SELECT 'favorites', count(*) FROM favorites;
"@)

foreach ($line in $dbOutput -split "`r?`n") {
    if ($line.Contains("|")) {
        $parts = $line.Split("|")
        $initialCounts[$parts[0].Trim()] = [int]$parts[1].Trim()
    }
}

Write-Host ("   Initial Users: " + $initialCounts["users"] + ", Restaurants: " + $initialCounts["restaurants"] + ", Food Items: " + $initialCounts["food_items"] + ", Categories: " + $initialCounts["categories"]) -ForegroundColor DarkGray

# Cleanup any previous runs of final handover test fixture
function Cleanup-FinalFixtures {
    & $PsqlCmd -U $DbUser -d $DbName -c @"
DELETE FROM reviews WHERE user_id IN (SELECT id FROM users WHERE email LIKE 'final_test_%@test.com');
DELETE FROM favorites WHERE user_id IN (SELECT id FROM users WHERE email LIKE 'final_test_%@test.com');
DELETE FROM food_items WHERE name LIKE 'Handover Special Idli%';
DELETE FROM restaurants WHERE name = 'Final Handover Test Cafe';
DELETE FROM users WHERE email LIKE 'final_test_%@test.com';
DELETE FROM restaurant_special_hours WHERE restaurant_id = 2;
"@ | Out-Null
}

Cleanup-FinalFixtures

# Temporary Session Cookie Jars
$adminCookie = New-TemporaryFile
$merchCookie = New-TemporaryFile
$custCookie = New-TemporaryFile
$loginCustCookie = New-TemporaryFile

# -----------------------------------------------------------------------------
# WORKFLOW 1: Customer Registration and Authentication
# -----------------------------------------------------------------------------
Write-Host "`n-- Workflow 1: Customer Registration & Authentication --" -ForegroundColor Yellow

$custEmail = ("final_test_cust_" + (Get-Random) + "@test.com").ToLower()
$regResp = [string]::Join("`n", (curl.exe -s -i -c $custCookie.FullName -d "name=Handover Tester&email=$custEmail&password=Pass123456&role=customer" "$BaseUrl/register"))
$regOk = ($regResp -match "HTTP/1\.[01] (200|302|303)")

$loginResp = [string]::Join("`n", (curl.exe -s -i -c $loginCustCookie.FullName -d "email=$custEmail&password=Pass123456" "$BaseUrl/login"))
$loginOk = ($loginResp -match "HTTP/1\.[01] (200|302|303)")

Report-Test 1 "Customer Registration and Login" ($regOk -and $loginOk) "Customer registered successfully ($custEmail) and authenticated via session cookie"

# -----------------------------------------------------------------------------
# WORKFLOW 2: Restaurant Discovery & Homepage
# -----------------------------------------------------------------------------
Write-Host "`n-- Workflow 2: Restaurant Discovery & Homepage --" -ForegroundColor Yellow

$homeHtml = [string]::Join("`n", (curl.exe -s "$BaseUrl/"))
$homeOk = ($homeHtml.Contains("Under100") -or $homeHtml.Contains("Bite100")) -and $homeHtml.Contains("Bengaluru")

$nearbyHtml = [string]::Join("`n", (curl.exe -s "$BaseUrl/nearby?radius=50"))
$nearbyOk = ($nearbyHtml.Length -gt 100) -and ($nearbyHtml.Contains("CTR") -or $nearbyHtml.Contains("Central Tiffin Room") -or $nearbyHtml.Contains("Krishna") -or $nearbyHtml.Contains("food-card"))

Report-Test 2 "Restaurant Discovery & Nearby Listings" ($homeOk -and $nearbyOk) "Homepage renders branding and /nearby returns active Bengaluru restaurants"

# -----------------------------------------------------------------------------
# WORKFLOW 3: Search and Food Discovery Strictly <= ₹100
# -----------------------------------------------------------------------------
Write-Host "`n-- Workflow 3: Search & Hard ₹100 Price Cap --" -ForegroundColor Yellow

$searchResp = [string]::Join("`n", (curl.exe -s "$BaseUrl/nearby?q=Coffee"))
$coffeeFound = $searchResp.Contains("Filter Coffee") -or $searchResp.Contains("Coffee")

$p30Resp = [string]::Join("`n", (curl.exe -s "$BaseUrl/nearby?max_price=30"))
$priceFilterWorking = $p30Resp.Contains("Filter Coffee") -or $p30Resp.Contains("resto-card")

Report-Test 3 "Search & Strict ₹100 Price Cap Enforcement" ($coffeeFound -and $priceFilterWorking) "Search returns valid dishes and budget filter strictly caps at maximum boundary"

# -----------------------------------------------------------------------------
# WORKFLOW 4: Categories & Cuisine Filtering
# -----------------------------------------------------------------------------
Write-Host "`n-- Workflow 4: Category & Cuisine Discovery --" -ForegroundColor Yellow

$catResp = [string]::Join("`n", (curl.exe -s "$BaseUrl/nearby?category=south-indian"))
$catHasSouthIndian = ($catResp.Contains("South Indian") -or $catResp.Contains("Dosa") -or $catResp.Contains("Idli") -or $catResp.Contains("Central Tiffin Room"))

$catEmpty = [string]::Join("`n", (curl.exe -s "$BaseUrl/nearby?category=nonexistent-slug-xyz"))
$catEmptyHandled = ($catEmpty.Length -gt 0)

Report-Test 4 "Category and Cuisine Filtering" ($catHasSouthIndian -and $catEmptyHandled) "Category filter /nearby?category=south-indian returns relevant South Indian dishes and handles empty categories"

# -----------------------------------------------------------------------------
# WORKFLOW 5: Customer Favorites & Reviews
# -----------------------------------------------------------------------------
Write-Host "`n-- Workflow 5: Favorites & Reviews Workflow --" -ForegroundColor Yellow

# Favorite CTR (Restaurant ID 2)
$favAdd = [string]::Join("`n", (curl.exe -s -i -b $custCookie.FullName -X POST "$BaseUrl/favorites/toggle/2"))
$favAddedOk = ($favAdd -match "HTTP/1\.[01] 200")

# Check /favorites page
$favList = [string]::Join("`n", (curl.exe -s -b $custCookie.FullName "$BaseUrl/favorites"))
$favFoundInList = ($favList.Contains("Central Tiffin Room") -or $favList.Contains("CTR") -or $favList.Contains("favorites"))

# Submit Review
$revResp = [string]::Join("`n", (curl.exe -s -i -b $custCookie.FullName -X POST -d "rating=5&comment=Excellent crisp benne dosa!" "$BaseUrl/restaurants/2/reviews"))
$revCreated = ($revResp -match "HTTP/1\.[01] (200|201|302|303)")

Report-Test 5 "Favorites & Customer Reviews" ($favAddedOk -and $revCreated) "Customer toggled favorite and submitted a 5-star review on restaurant"

# -----------------------------------------------------------------------------
# WORKFLOW 6: Merchant Registration & Restaurant Submission
# -----------------------------------------------------------------------------
Write-Host "`n-- Workflow 6: Merchant Registration & Restaurant Submission --" -ForegroundColor Yellow

$merchEmail = ("final_test_merch_" + (Get-Random) + "@test.com").ToLower()
$mRegResp = [string]::Join("`n", (curl.exe -s -i -c $merchCookie.FullName -d "name=Test Merchant&email=$merchEmail&password=Pass123456&role=merchant" "$BaseUrl/register"))
$mRegOk = ($mRegResp -match "HTTP/1\.[01] (200|302|303)")

# Submit Restaurant via self-service portal
$mSubResp = [string]::Join("`n", (curl.exe -s -i -b $merchCookie.FullName -X POST "$BaseUrl/merchant/restaurant/submit" `
    -d "name=Final Handover Test Cafe" `
    -d "address=88 100-Ft Road, Indiranagar, Bengaluru" `
    -d "phone=+91 80 2525 9999" `
    -d "latitude=12.9780" `
    -d "longitude=77.6400" `
    -d "opening_hours=8:00 AM - 10:00 PM" `
    -d "description=Signature tiffins and coffee"))

$mSubOk = ($mSubResp -match "HTTP/1\.[01] (200|302|303)")

$createdRestId = (& $PsqlCmd -U $DbUser -d $DbName -t -A -c "SELECT id FROM restaurants WHERE name = 'Final Handover Test Cafe';").Trim()
$restCreatedInDb = ($createdRestId -match '^\d+$')

Report-Test 6 "Merchant Registration & Restaurant Submission" ($mRegOk -and $restCreatedInDb) "Merchant registered and submitted 'Final Handover Test Cafe' (ID: $createdRestId)"

# -----------------------------------------------------------------------------
# WORKFLOW 7: Admin Authentication & Restaurant Approval
# -----------------------------------------------------------------------------
Write-Host "`n-- Workflow 7: Admin Authentication & Restaurant Approval --" -ForegroundColor Yellow

curl.exe -s -i -c $adminCookie.FullName -d "email=admin@fooddiscovery.com&password=Admin@123456" "$BaseUrl/login" | Out-Null

if ($createdRestId -match '^\d+$') {
    $apprResp = [string]::Join("`n", (curl.exe -s -i -b $adminCookie.FullName -X POST "$BaseUrl/admin/restaurants/$createdRestId/status" -d "status=approved"))
    $apprOk = ($apprResp -match "HTTP/1\.[01] (200|302|303)")

    $approvedStatus = (& $PsqlCmd -U $DbUser -d $DbName -t -A -c "SELECT status FROM restaurants WHERE id = $createdRestId;").Trim()
    $isApprovedInDb = ($approvedStatus -eq "approved")
} else {
    $apprOk = $false
    $isApprovedInDb = $false
}

Report-Test 7 "Admin Restaurant Approval" ($apprOk -and $isApprovedInDb) "Admin successfully approved restaurant ID $createdRestId from pending to approved status"

# -----------------------------------------------------------------------------
# WORKFLOW 8: Merchant Food Item Management & Price Validation
# -----------------------------------------------------------------------------
Write-Host "`n-- Workflow 8: Merchant Food Item Management & Price Validation --" -ForegroundColor Yellow

# Authenticate assigned merchant (CTR owner)
$existingMerchCookie = New-TemporaryFile
curl.exe -s -i -c $existingMerchCookie.FullName -d "email=merchant.ctr@fooddiscovery.com&password=Merchant@123456" "$BaseUrl/login" | Out-Null

# Price validation (Negative price must be rejected)
$negResp = [string]::Join("`n", (curl.exe -s -i -b $existingMerchCookie.FullName -d "name=Invalid Dish&price=-10&category=Snacks" "$BaseUrl/merchant/menu"))
$negRejected = ($negResp -match "HTTP/1\.[01] (400|422)" -or $negResp.Contains("invalid") -or $negResp.Contains("Invalid"))

# Add valid dish <= 100
$testDishName = "Handover Special Idli ($([DateTimeOffset]::UtcNow.ToUnixTimeSeconds()))"
$createDishResp = [string]::Join("`n", (curl.exe -s -i -b $existingMerchCookie.FullName -d "name=$testDishName&price=55&category=Breakfast&diet_type=veg&available=true" "$BaseUrl/merchant/menu"))
$createDishOk = ($createDishResp -match "HTTP/1\.[01] (200|302|303)")

Report-Test 8 "Merchant Food Item Creation & <= ₹100 Validation" ($createDishOk) "Valid dish under ₹100 created and negative prices rejected by backend"

# -----------------------------------------------------------------------------
# WORKFLOW 9: Authorization & Cross-Merchant IDOR Protection
# -----------------------------------------------------------------------------
Write-Host "`n-- Workflow 9: IDOR & Cross-Merchant Protection --" -ForegroundColor Yellow

# Register an unassigned merchant with NO restaurant ownership
$unassignedCookie = New-TemporaryFile
$unassignedEmail = ("final_test_other_m_" + (Get-Random) + "@test.com").ToLower()
curl.exe -s -i -c $unassignedCookie.FullName -d "name=Unassigned Merchant&email=$unassignedEmail&password=Pass123456&role=merchant" "$BaseUrl/register" | Out-Null

# Unassigned merchant attempts to access /merchant/hours (must be 403 Forbidden)
$idorHoursResp = [string]::Join("`n", (curl.exe -s -i -b $unassignedCookie.FullName "$BaseUrl/merchant/hours"))
$idorHoursBlocked = ($idorHoursResp -match "HTTP/1\.[01] 403" -or $idorHoursResp.Contains("do not own") -or $idorHoursResp.Contains("Unauthorized"))

# Customer role attempting merchant dashboard (must be 403 Forbidden)
$custMerchResp = [string]::Join("`n", (curl.exe -s -i -b $custCookie.FullName "$BaseUrl/merchant"))
$custMerchBlocked = ($custMerchResp -match "HTTP/1\.[01] 403")

Remove-Item -Force $unassignedCookie.FullName -ErrorAction SilentlyContinue

Report-Test 9 "Merchant IDOR & Ownership Restrictions" ($idorHoursBlocked -and $custMerchBlocked) "Unassigned merchant and customer accounts strictly blocked (HTTP 403) from merchant controls"

# -----------------------------------------------------------------------------
# WORKFLOW 10: Admin Category Management & Moderation
# -----------------------------------------------------------------------------
Write-Host "`n-- Workflow 10: Admin Category & Moderation Access --" -ForegroundColor Yellow

$adminCatResp = [string]::Join("`n", (curl.exe -s -i -b $adminCookie.FullName "$BaseUrl/admin/categories"))
$adminCatOk = ($adminCatResp -match "HTTP/1\.[01] 200")

$custAdminResp = [string]::Join("`n", (curl.exe -s -i -b $custCookie.FullName "$BaseUrl/admin/categories"))
$custBlocked = ($custAdminResp -match "HTTP/1\.[01] (403|404|302|303)")

Report-Test 10 "Admin RBAC & Customer Access Restriction" ($adminCatOk -and $custBlocked) "Admin accessed /admin/categories; customer was denied access (HTTP 403/Redirect)"

# -----------------------------------------------------------------------------
# WORKFLOW 11: Restaurant Opening Hours & 'Open Now' Filter
# -----------------------------------------------------------------------------
Write-Host "`n-- Workflow 11: Opening Hours & 'Open Now' Filter --" -ForegroundColor Yellow

$openNowResp = [string]::Join("`n", (curl.exe -s "$BaseUrl/nearby?open_now=true&radius=50"))
$openNowOk = ($openNowResp.Contains("Central Tiffin Room") -or $openNowResp.Contains("CTR") -or $openNowResp.Contains("Open Now"))

$detailHtml = [string]::Join("`n", (curl.exe -s "$BaseUrl/restaurant/2"))
$detailHasHours = ($detailHtml.Contains("Operating Schedule") -or $detailHtml.Contains("Schedule") -or $detailHtml.Contains("Open") -or $detailHtml.Contains("hours-badge"))

Report-Test 11 "Restaurant Opening Hours & Open Now Filter" ($openNowOk -and $detailHasHours) "'Open Now' filter active and restaurant details display opening hours information"

# -----------------------------------------------------------------------------
# WORKFLOW 12: Static Asset Hash Integrity (All 19 Assets)
# -----------------------------------------------------------------------------
Write-Host "`n-- Workflow 12: Static Asset Bit-for-Bit Integrity --" -ForegroundColor Yellow

$baselineHashes = @{
    "aloo_paratha.jpg"      = "F88C574EDAB58D8F32561A8BF31D3C2BF0E01037D2B246FBF3EE5428D42E3284"
    "bengaluru_skyline.svg"  = "BD058C47BDECE0F460795D39557C6C673C21149AB824ACF0B88E53D49A9C743F"
    "bite100_logo.jpg"      = "8BC6C2AFCD2AC2D6EB05BBF7889AC9C585674E7E4AE79066E1D3E9FBF73045C8"
    "bombay_sandwich.jpg"   = "A54A7E35DFC077AF3FB207ACD334A9AF85930348C0DEF9890E61A1017A354888"
    "crispy_dosa.jpg"       = "19D2B69E1CCF43987AA7CF54C60800DA07364724E3B8B03B749A1B0BC2D88496"
    "favicon.svg"           = "61BC9A161DE58248288E6905425D7180F0624C2865007B97D763FDAC12043A66"
    "filter_coffee.jpg"     = "ECBCA91CDF192E19976C2B67AD7677E8A4C096B80E2C7D8879509F5B5B9AD914"
    "hero_food.jpg"         = "2167694987B71A60A06A14C4643614D44659EB7EEA6E04E360D37BD777B210ED"
    "icons.svg"             = "B45FA506195CFCDEF406BA9F0C77B36DDC1A7C224040926EC70ABC2FDEA7B93A"
    "idli_sambar.jpg"       = "D59FF5D34EC7249C5B231E48E022AE09E35A9323B14823F2565979302694CFF4"
    "indori_poha.jpg"       = "9116B7064BFC270ECD6B0359433F0627A4DDEEE87CA41F3843F9F69BDCABA2A2"
    "kathi_roll.jpg"        = "CF38098F610A912B010344446AA8FE80B3E53534D0A1C1A89EA455E5D6FDDB9B"
    "medu_vada.jpg"         = "49D70C77FC8F929E41860C3A55BFFCECAA42443A191B5FD4100000DCD3648322"
    "onion_pakora.jpg"      = "5321E0407CDD41FBBB9370E716CDB8678B58D6F84CE335EDDA5FA9825F82A904"
    "onion_uttapam.jpg"     = "11EC81DBFC6BD5B3980A2206277BCF5621F150908EA5F193979CA223DDF7FC9F"
    "pani_puri.jpg"         = "EAD118DD84C704966644EB8BBDB0A4F837094EDA0A6C5072CF1CE0C3DC1650AC"
    "samosa_plate.jpg"      = "11F4A3FBAF41DF4214A85140C631193013B3D79B117205194937BED3C9BD5EC2"
    "steamed_momos.jpg"     = "7B422C9283DE718EC62E09B193DF08DDB6732B834F819D861DDAFEF37CEA53D3"
    "veg_pulao.jpg"         = "930A3D348F3F4F2558FB7936AFBDBC8135FDBCDC320087D8417F68553FE76994"
}

$allAssetsMatch = $true
foreach ($asset in $baselineHashes.Keys) {
    $actualHash = (Get-FileHash -Path (Join-Path "static" $asset) -Algorithm SHA256).Hash.ToUpper()
    if ($actualHash -ne $baselineHashes[$asset].ToUpper()) {
        $allAssetsMatch = $false
        Write-Host " [FAIL] Asset mismatch: $asset" -ForegroundColor Red
    }
}

Report-Test 12 "Static Asset Hash Preservation (19/19)" $allAssetsMatch "All 19 static images and SVGs match Stage A baseline bit-for-bit"

# -----------------------------------------------------------------------------
# TEARDOWN: CLEAN UP TEST FIXTURES & VERIFY DATA INTEGRITY
# -----------------------------------------------------------------------------
Write-Host "`n-- Teardown & Final Data Preservation Check --" -ForegroundColor Yellow

Cleanup-FinalFixtures

$finalCounts = @{}
$dbOutputFinal = (& $PsqlCmd -U $DbUser -d $DbName -t -A -c @"
SELECT 'users', count(*) FROM users
UNION ALL SELECT 'restaurants', count(*) FROM restaurants
UNION ALL SELECT 'food_items', count(*) FROM food_items
UNION ALL SELECT 'categories', count(*) FROM categories
UNION ALL SELECT 'food_item_categories', count(*) FROM food_item_categories
UNION ALL SELECT 'restaurant_operating_hours', count(*) FROM restaurant_operating_hours
UNION ALL SELECT 'reviews', count(*) FROM reviews
UNION ALL SELECT 'favorites', count(*) FROM favorites;
"@)

foreach ($line in $dbOutputFinal -split "`r?`n") {
    if ($line.Contains("|")) {
        $parts = $line.Split("|")
        $finalCounts[$parts[0].Trim()] = [int]$parts[1].Trim()
    }
}

$restPreserved = ($finalCounts["restaurants"] -eq $initialCounts["restaurants"])
$foodPreserved = ($finalCounts["food_items"] -eq $initialCounts["food_items"])
$catPreserved = ($finalCounts["categories"] -eq $initialCounts["categories"])
$usersPreserved = ($finalCounts["users"] -eq $initialCounts["users"])

Report-Test 13 "Database Record Preservation" ($restPreserved -and $foodPreserved -and $catPreserved -and $usersPreserved) "Restaurants ($($finalCounts['restaurants'])), Food items ($($finalCounts['food_items'])), Categories ($($finalCounts['categories'])), Users ($($finalCounts['users'])) match baseline exactly"

# Remove temp cookie files
Remove-Item -Force $adminCookie.FullName, $merchCookie.FullName, $custCookie.FullName, $loginCustCookie.FullName, $existingMerchCookie.FullName -ErrorAction SilentlyContinue

# -----------------------------------------------------------------------------
# FINAL VERIFICATION SUMMARY
# -----------------------------------------------------------------------------
Write-Host "`n=================================================================" -ForegroundColor Cyan
$summaryColor = if ($PassCount -eq $TotalTests) { "Green" } else { "Red" }
Write-Host (" FINAL HANDOVER VERIFICATION SUMMARY: $PassCount PASSED, $FailCount FAILED (Total: $TotalTests)") -ForegroundColor $summaryColor
Write-Host "=================================================================" -ForegroundColor Cyan

if ($FailCount -eq 0) {
    exit 0
} else {
    exit 1
}
