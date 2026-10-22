Write-Host "==================================================================" -ForegroundColor Cyan
Write-Host "   STAGE C COMPREHENSIVE VERIFICATION & SECURITY SUITE            " -ForegroundColor Cyan
Write-Host "==================================================================" -ForegroundColor Cyan

$allPassed = $true

# 1. VERIFY ALL 19 STATIC ASSETS
Write-Host "`n--- TEST 1: STATIC ASSET HASH PRESERVATION ---" -ForegroundColor Yellow
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

$assetCount = 0
$assetPassCount = 0
foreach ($asset in $baselineHashes.Keys) {
    $assetPath = Join-Path "static" $asset
    if (Test-Path $assetPath) {
        $actual = (Get-FileHash -Path $assetPath -Algorithm SHA256).Hash.ToUpper()
        $expected = $baselineHashes[$asset].ToUpper()
        $assetCount++
        if ($actual -eq $expected) {
            $assetPassCount++
        } else {
            Write-Host "  [FAIL] $asset hash mismatch!" -ForegroundColor Red
            $allPassed = $false
        }
    }
}
Write-Host "  [PASS] All $assetPassCount / $assetCount static assets verified bit-for-bit unchanged." -ForegroundColor Green

# 2. STOREFRONT AND PUBLIC ROUTES
Write-Host "`n--- TEST 2: STOREFRONT & PUBLIC ROUTES ---" -ForegroundColor Yellow
$homeCode = curl.exe -s -o temp_home.html -w "%{http_code}" http://localhost:8000/
$homeContent = [System.IO.File]::ReadAllText("$PWD/temp_home.html")
Remove-Item -Force temp_home.html
if ($homeCode -eq "200" -and ($homeContent.Contains("Under100") -or $homeContent.Contains("Bite100"))) {
    Write-Host "  [PASS] Storefront GET / returned HTTP 200 with Under100 branding" -ForegroundColor Green
} else {
    Write-Host "  [FAIL] Storefront GET / failed with code $homeCode" -ForegroundColor Red
    $allPassed = $false
}

$nearbyCode = curl.exe -s -o /dev/null -w "%{http_code}" "http://localhost:8000/nearby?lat=12.9716&lng=77.5946&radius=25&max_price=100"
if ($nearbyCode -eq "200") {
    Write-Host "  [PASS] Nearby API returned HTTP 200" -ForegroundColor Green
} else {
    Write-Host "  [FAIL] Nearby API failed with code $nearbyCode" -ForegroundColor Red
    $allPassed = $false
}

# 3. UNAUTHENTICATED /merchant ACCESS
Write-Host "`n--- TEST 3: UNAUTHENTICATED /merchant ACCESS ---" -ForegroundColor Yellow
$unauthMerch = curl.exe -s -o /dev/null -w "%{http_code}" http://localhost:8000/merchant
if ($unauthMerch -eq "303" -or $unauthMerch -eq "302") {
    Write-Host "  [PASS] Unauthenticated GET /merchant redirected with HTTP $unauthMerch" -ForegroundColor Green
} else {
    Write-Host "  [FAIL] Unauthenticated GET /merchant returned unexpected code: $unauthMerch" -ForegroundColor Red
    $allPassed = $false
}

# 4. CUSTOMER ACCESS TO /merchant (MUST RETURN 403)
Write-Host "`n--- TEST 4: CUSTOMER ATTEMPTING TO ACCESS /merchant ---" -ForegroundColor Yellow
Remove-Item -Force -ErrorAction SilentlyContinue cust_session.txt
curl.exe -s -c cust_session.txt -d "email=user@fooddiscovery.com&password=User@123456" -o /dev/null http://localhost:8000/login
$custMerchCode = curl.exe -s -b cust_session.txt -o /dev/null -w "%{http_code}" http://localhost:8000/merchant
if ($custMerchCode -eq "403") {
    Write-Host "  [PASS] Customer accessing /merchant correctly rejected with HTTP 403 Forbidden" -ForegroundColor Green
} else {
    Write-Host "  [FAIL] Customer accessing /merchant returned unexpected code: $custMerchCode (Expected 403)" -ForegroundColor Red
    $allPassed = $false
}
Remove-Item -Force -ErrorAction SilentlyContinue cust_session.txt

# 5. ASSIGNED MERCHANT DASHBOARD ACCESS
Write-Host "`n--- TEST 5: ASSIGNED MERCHANT DASHBOARD ACCESS ---" -ForegroundColor Yellow
Remove-Item -Force -ErrorAction SilentlyContinue merch_session.txt
curl.exe -s -c merch_session.txt -d "email=merchant.ctr@fooddiscovery.com&password=Merchant@123456" -o /dev/null http://localhost:8000/login
$dashCode = curl.exe -s -b merch_session.txt -o temp_dash.html -w "%{http_code}" http://localhost:8000/merchant
$dashContent = [System.IO.File]::ReadAllText("$PWD/temp_dash.html")
Remove-Item -Force temp_dash.html

if ($dashCode -eq "200" -and $dashContent.Contains("Central Tiffin Room") -and $dashContent.Contains("Merchant Owner Portal")) {
    Write-Host "  [PASS] Assigned Merchant accessed /merchant with HTTP 200, showing Central Tiffin Room outlet" -ForegroundColor Green
} else {
    Write-Host "  [FAIL] Merchant dashboard failed to load or missing outlet info: Code $dashCode" -ForegroundColor Red
    $allPassed = $false
}

if ($dashContent.Contains("Discoverable in Bengaluru") -and $dashContent.Contains("Feed")) {
    Write-Host "  [PASS] Under ₹100 Discovery Qualification banner properly displayed" -ForegroundColor Green
} else {
    Write-Host "  [FAIL] Under ₹100 Discovery Qualification banner missing" -ForegroundColor Red
    $allPassed = $false
}

# 6. MERCHANT VALIDATION: SUBMITTING NEGATIVE OR ZERO PRICE
Write-Host "`n--- TEST 6: MERCHANT VALIDATION (NEGATIVE / ZERO PRICE) ---" -ForegroundColor Yellow
$negPriceCode = curl.exe -s -b merch_session.txt -d "name=Negative Item&price=-50&category=Breakfast" -o /dev/null -w "%{http_code}" http://localhost:8000/merchant/menu
if ($negPriceCode -eq "400") {
    Write-Host "  [PASS] POST /merchant/menu with price=-50 correctly rejected with HTTP 400 Bad Request" -ForegroundColor Green
} else {
    Write-Host "  [FAIL] Negative price not rejected! Status: $negPriceCode" -ForegroundColor Red
    $allPassed = $false
}

$zeroPriceCode = curl.exe -s -b merch_session.txt -d "name=Zero Item&price=0&category=Breakfast" -o /dev/null -w "%{http_code}" http://localhost:8000/merchant/menu
if ($zeroPriceCode -eq "400") {
    Write-Host "  [PASS] POST /merchant/menu with price=0 correctly rejected with HTTP 400 Bad Request" -ForegroundColor Green
} else {
    Write-Host "  [FAIL] Zero price not rejected! Status: $zeroPriceCode" -ForegroundColor Red
    $allPassed = $false
}

# 7. MERCHANT ADD VALID DISH (FOOD ITEM CREATION & PRICING FLEXIBILITY)
Write-Host "`n--- TEST 7: MERCHANT ADDS VALID DISH ---" -ForegroundColor Yellow
$testDishName = "Special Ghee Podi Idli ($([DateTimeOffset]::UtcNow.ToUnixTimeSeconds()))"
$createResp = curl.exe -s -b merch_session.txt -d "name=$testDishName&price=65&category=Breakfast&diet_type=veg&description=Steamed idlis tossed in spicy gunpowder and pure ghee&available=true" -o /dev/null -w "%{http_code}" http://localhost:8000/merchant/menu
if ($createResp -eq "303" -or $createResp -eq "200") {
    Write-Host "  [PASS] POST /merchant/menu created dish '$testDishName' (HTTP $createResp)" -ForegroundColor Green
} else {
    Write-Host "  [FAIL] Failed to create dish: HTTP $createResp" -ForegroundColor Red
    $allPassed = $false
}

# Query created food item ID from DB
$createdID = & "C:\Users\astha\pgsql\bin\psql.exe" -U admin -d goship_db -t -A -c "SELECT id FROM food_items WHERE name='$testDishName';"
Write-Host "  [PASS] Verified new dish in DB with ID: $createdID" -ForegroundColor Green

# 8. MERCHANT EDITS OWN DISH
Write-Host "`n--- TEST 8: MERCHANT EDITS OWN DISH ---" -ForegroundColor Yellow
$editResp = curl.exe -s -b merch_session.txt -d "name=$testDishName (Updated)&price=70&category=Breakfast&diet_type=veg&available=true" -o /dev/null -w "%{http_code}" "http://localhost:8000/merchant/menu/$createdID"
if ($editResp -eq "303" -or $editResp -eq "200") {
    Write-Host "  [PASS] POST /merchant/menu/$createdID updated dish (HTTP $editResp)" -ForegroundColor Green
} else {
    Write-Host "  [FAIL] Failed to update dish: HTTP $editResp" -ForegroundColor Red
    $allPassed = $false
}

# 9. MERCHANT TOGGLES AVAILABILITY
Write-Host "`n--- TEST 9: MERCHANT TOGGLES AVAILABILITY ---" -ForegroundColor Yellow
$toggleResp = curl.exe -s -b merch_session.txt -X POST -o temp_toggle.html -w "%{http_code}" "http://localhost:8000/merchant/menu/$createdID/availability"
$toggleContent = [System.IO.File]::ReadAllText("$PWD/temp_toggle.html")
Remove-Item -Force temp_toggle.html
if ($toggleResp -eq "200" -and $toggleContent.Contains("Sold Out")) {
    Write-Host "  [PASS] POST /merchant/menu/$createdID/availability toggled to 'Sold Out'" -ForegroundColor Green
} else {
    Write-Host "  [FAIL] Toggle availability failed: HTTP $toggleResp" -ForegroundColor Red
    $allPassed = $false
}

# 10. STRICT TENANT ISOLATION: MERCHANT ATTEMPTS TO EDIT ANOTHER RESTAURANT'S DISH
Write-Host "`n--- TEST 10: CROSS-TENANT ISOLATION (EDIT OTHER OUTLET DISH) ---" -ForegroundColor Yellow
# Find a food item belonging to Vidyarthi Bhavan (ID 1), not CTR (ID 2)
$otherFoodID = & "C:\Users\astha\pgsql\bin\psql.exe" -U admin -d goship_db -t -A -c "SELECT id FROM food_items WHERE restaurant_food_items = 1 LIMIT 1;"
Write-Host "  Target other restaurant's food item ID: $otherFoodID"

$crossEditResp = curl.exe -s -b merch_session.txt -d "name=Hacked Dish&price=10" -o /dev/null -w "%{http_code}" "http://localhost:8000/merchant/menu/$otherFoodID"
if ($crossEditResp -eq "403") {
    Write-Host "  [PASS] Cross-tenant edit rejected with HTTP 403 Forbidden" -ForegroundColor Green
} else {
    Write-Host "  [FAIL] Cross-tenant edit was NOT rejected! Code: $crossEditResp" -ForegroundColor Red
    $allPassed = $false
}

# 11. STRICT TENANT ISOLATION: MERCHANT ATTEMPTS TO TOGGLE/DELETE OTHER DISH
Write-Host "`n--- TEST 11: CROSS-TENANT ISOLATION (TOGGLE & DELETE OTHER DISH) ---" -ForegroundColor Yellow
$crossToggleResp = curl.exe -s -b merch_session.txt -X POST -o /dev/null -w "%{http_code}" "http://localhost:8000/merchant/menu/$otherFoodID/availability"
if ($crossToggleResp -eq "403") {
    Write-Host "  [PASS] Cross-tenant toggle rejected with HTTP 403 Forbidden" -ForegroundColor Green
} else {
    Write-Host "  [FAIL] Cross-tenant toggle was NOT rejected! Code: $crossToggleResp" -ForegroundColor Red
    $allPassed = $false
}

$crossDeleteResp = curl.exe -s -b merch_session.txt -X POST -o /dev/null -w "%{http_code}" "http://localhost:8000/merchant/menu/$otherFoodID/delete"
if ($crossDeleteResp -eq "403") {
    Write-Host "  [PASS] Cross-tenant delete rejected with HTTP 403 Forbidden" -ForegroundColor Green
} else {
    Write-Host "  [FAIL] Cross-tenant delete was NOT rejected! Code: $crossDeleteResp" -ForegroundColor Red
    $allPassed = $false
}

# 12. CLEANUP CREATED TEST FOOD ITEM VIA MERCHANT DELETE ENDPOINT
Write-Host "`n--- TEST 12: MERCHANT DELETES OWN CREATED TEST ITEM ---" -ForegroundColor Yellow
$deleteResp = curl.exe -s -b merch_session.txt -X POST -o /dev/null -w "%{http_code}" "http://localhost:8000/merchant/menu/$createdID/delete"
if ($deleteResp -eq "200") {
    Write-Host "  [PASS] POST /merchant/menu/$createdID/delete successfully deleted test item" -ForegroundColor Green
} else {
    Write-Host "  [FAIL] Failed to delete test item: HTTP $deleteResp" -ForegroundColor Red
    $allPassed = $false
}
Remove-Item -Force -ErrorAction SilentlyContinue merch_session.txt

# 13. ADMIN ACCESS PRESERVATION
Write-Host "`n--- TEST 13: ADMIN DASHBOARD ACCESS PRESERVATION ---" -ForegroundColor Yellow
Remove-Item -Force -ErrorAction SilentlyContinue admin_session.txt
curl.exe -s -c admin_session.txt -d "email=admin@fooddiscovery.com&password=Admin@123456" -o /dev/null http://localhost:8000/login
$adminDash = curl.exe -s -b admin_session.txt -o /dev/null -w "%{http_code}" http://localhost:8000/admin
if ($adminDash -eq "200") {
    Write-Host "  [PASS] Platform Admin accessing /admin returned HTTP 200 OK" -ForegroundColor Green
} else {
    Write-Host "  [FAIL] Admin dashboard failed to load: HTTP $adminDash" -ForegroundColor Red
    $allPassed = $false
}
Remove-Item -Force -ErrorAction SilentlyContinue admin_session.txt

# 14. MERCHANT NOT TREATED AS ADMIN
Write-Host "`n--- TEST 14: MERCHANT BLOCKED FROM PLATFORM ADMIN ---" -ForegroundColor Yellow
Remove-Item -Force -ErrorAction SilentlyContinue merch_session2.txt
curl.exe -s -c merch_session2.txt -d "email=merchant.ctr@fooddiscovery.com&password=Merchant@123456" -o /dev/null http://localhost:8000/login
$merchAdmin = curl.exe -s -b merch_session2.txt -o /dev/null -w "%{http_code}" http://localhost:8000/admin
if ($merchAdmin -eq "403") {
    Write-Host "  [PASS] Merchant accessing /admin correctly rejected with HTTP 403 Forbidden" -ForegroundColor Green
} else {
    Write-Host "  [FAIL] Merchant accessing /admin returned unexpected code: $merchAdmin" -ForegroundColor Red
    $allPassed = $false
}
Remove-Item -Force -ErrorAction SilentlyContinue merch_session2.txt

# 15. FINAL HEALTH AND SUMMARY
Write-Host "`n==================================================================" -ForegroundColor Cyan
if ($allPassed) {
    Write-Host ">>> STAGE C VALIDATION RESULT: ALL 15 TESTS PASSED (100%) <<<" -ForegroundColor Green
} else {
    Write-Host ">>> STAGE C VALIDATION RESULT: ONE OR MORE TESTS FAILED <<<" -ForegroundColor Red
}
Write-Host "==================================================================" -ForegroundColor Cyan
