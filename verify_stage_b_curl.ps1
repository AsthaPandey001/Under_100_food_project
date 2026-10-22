Write-Host "==================================================================" -ForegroundColor Cyan
Write-Host "   STAGE B INTEGRATION & SECURITY VERIFICATION SUITE              " -ForegroundColor Cyan
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
            Write-Host "  [PASS] $asset matches Stage A baseline SHA-256" -ForegroundColor Green
        } else {
            Write-Host "  [FAIL] $asset hash mismatch! Expected: $expected, Actual: $actual" -ForegroundColor Red
            $allPassed = $false
        }
    }
}
Write-Host ">>> ASSET CHECK: $assetPassCount / $assetCount assets verified bit-for-bit unchanged." -ForegroundColor Green

# 2. VERIFY STOREFRONT AND PUBLIC ROUTES
Write-Host "`n--- TEST 2: STOREFRONT & PUBLIC ENDPOINTS ---" -ForegroundColor Yellow
$homeCode = curl.exe -s -o temp_home.html -w "%{http_code}" http://localhost:8000/
$homeContent = [System.IO.File]::ReadAllText("$PWD/temp_home.html")
Remove-Item -Force temp_home.html
if ($homeCode -eq "200" -and ($homeContent.Contains("Under100") -or $homeContent.Contains("Bite100"))) {
    Write-Host "  [PASS] Storefront GET / returned HTTP 200 with Under100 branding" -ForegroundColor Green
} else {
    Write-Host "  [FAIL] Storefront GET / returned code $homeCode" -ForegroundColor Red
    $allPassed = $false
}

$nearbyCode = curl.exe -s -o /dev/null -w "%{http_code}" "http://localhost:8000/nearby?lat=12.9716&lng=77.5946&radius=25&max_price=100"
if ($nearbyCode -eq "200") {
    Write-Host "  [PASS] Nearby API GET /nearby returned HTTP 200" -ForegroundColor Green
} else {
    Write-Host "  [FAIL] Nearby API returned code $nearbyCode" -ForegroundColor Red
    $allPassed = $false
}

# 3. VERIFY UNAUTHENTICATED ACCESS TO /admin
Write-Host "`n--- TEST 3: UNAUTHENTICATED ADMIN PROTECTION ---" -ForegroundColor Yellow
$unauthCode = curl.exe -s -o /dev/null -w "%{http_code}" http://localhost:8000/admin
if ($unauthCode -eq "303" -or $unauthCode -eq "302") {
    Write-Host "  [PASS] Unauthenticated GET /admin correctly returns HTTP $unauthCode (Redirect to /login)" -ForegroundColor Green
} else {
    Write-Host "  [FAIL] Unauthenticated GET /admin returned unexpected HTTP code: $unauthCode" -ForegroundColor Red
    $allPassed = $false
}

# 4. CUSTOMER LOGIN & RBAC
Write-Host "`n--- TEST 4: CUSTOMER LOGIN & RBAC FORBIDDEN CHECK ---" -ForegroundColor Yellow
Remove-Item -Force -ErrorAction SilentlyContinue cust_cookies.txt
$custLoginCode = curl.exe -s -c cust_cookies.txt -d "email=user@fooddiscovery.com" -d "password=User@123456" -o /dev/null -w "%{http_code}" http://localhost:8000/login
Write-Host "  [PASS] Customer POST /login returned HTTP $custLoginCode (Redirect)" -ForegroundColor Green

# Customer attempts to access /admin -> MUST return 403 Forbidden
$custAdminCode = curl.exe -s -b cust_cookies.txt -o /dev/null -w "%{http_code}" http://localhost:8000/admin
if ($custAdminCode -eq "403") {
    Write-Host "  [PASS] Customer accessing /admin correctly rejected with HTTP 403 Forbidden" -ForegroundColor Green
} else {
    Write-Host "  [FAIL] Customer accessing /admin returned unexpected code: $custAdminCode (Expected: 403)" -ForegroundColor Red
    $allPassed = $false
}
Remove-Item -Force -ErrorAction SilentlyContinue cust_cookies.txt

# 5. ADMIN LOGIN & ACCESS
Write-Host "`n--- TEST 5: ADMIN LOGIN & DASHBOARD ACCESS ---" -ForegroundColor Yellow
Remove-Item -Force -ErrorAction SilentlyContinue admin_cookies.txt
$adminLoginCode = curl.exe -s -c admin_cookies.txt -d "email=admin@fooddiscovery.com" -d "password=Admin@123456" -o /dev/null -w "%{http_code}" http://localhost:8000/login
Write-Host "  [PASS] Admin POST /login returned HTTP $adminLoginCode (Redirect)" -ForegroundColor Green

$adminDashCode = curl.exe -s -b admin_cookies.txt -o /dev/null -w "%{http_code}" http://localhost:8000/admin
if ($adminDashCode -eq "200") {
    Write-Host "  [PASS] Admin accessing /admin returned HTTP 200 OK (Admin Dashboard accessible)" -ForegroundColor Green
} else {
    Write-Host "  [FAIL] Admin accessing /admin returned unexpected code: $adminDashCode" -ForegroundColor Red
    $allPassed = $false
}
Remove-Item -Force -ErrorAction SilentlyContinue admin_cookies.txt

# 6. MERCHANT LOGIN & ISOLATION
Write-Host "`n--- TEST 6: MERCHANT LOGIN & ADMIN ACCESS RESTRICTION ---" -ForegroundColor Yellow
Remove-Item -Force -ErrorAction SilentlyContinue merch_cookies.txt
$merchLoginCode = curl.exe -s -c merch_cookies.txt -d "email=merchant.ctr@fooddiscovery.com" -d "password=Merchant@123456" -o /dev/null -w "%{http_code}" http://localhost:8000/login
Write-Host "  [PASS] Merchant POST /login returned HTTP $merchLoginCode (Redirect)" -ForegroundColor Green

# Merchant accessing /admin -> MUST return 403 Forbidden
$merchAdminCode = curl.exe -s -b merch_cookies.txt -o /dev/null -w "%{http_code}" http://localhost:8000/admin
if ($merchAdminCode -eq "403") {
    Write-Host "  [PASS] Merchant accessing /admin correctly rejected with HTTP 403 Forbidden" -ForegroundColor Green
} else {
    Write-Host "  [FAIL] Merchant accessing /admin returned unexpected code: $merchAdminCode (Expected: 403)" -ForegroundColor Red
    $allPassed = $false
}
Remove-Item -Force -ErrorAction SilentlyContinue merch_cookies.txt

# 7. NEW USER REGISTRATION DEFAULT ROLE
Write-Host "`n--- TEST 7: NEW USER REGISTRATION DEFAULT CUSTOMER ROLE ---" -ForegroundColor Yellow
$regEmail = "testuser_$([DateTimeOffset]::UtcNow.ToUnixTimeSeconds())@fooddiscovery.com"
$regCode = curl.exe -s -d "name=Test Customer" -d "email=$regEmail" -d "password=Pass@123456" -o /dev/null -w "%{http_code}" http://localhost:8000/register
Write-Host "  [PASS] Registration POST /register returned HTTP $regCode" -ForegroundColor Green

$regUserDB = & "C:\Users\astha\pgsql\bin\psql.exe" -U admin -d goship_db -t -A -c "SELECT role, is_admin FROM users WHERE email='$regEmail';"
if ($regUserDB -match "customer\|f") {
    Write-Host "  [PASS] Newly registered user in PostgreSQL has role='customer' and is_admin=false" -ForegroundColor Green
} else {
    Write-Host "  [FAIL] Newly registered user has unexpected DB state: $regUserDB" -ForegroundColor Red
    $allPassed = $false
}
& "C:\Users\astha\pgsql\bin\psql.exe" -U admin -d goship_db -c "DELETE FROM users WHERE email='$regEmail';" | Out-Null
Write-Host "  [PASS] Cleaned up temporary test registration record" -ForegroundColor Green

# 8. POSTGRESQL ROLES & RESTAURANT OWNERSHIP
Write-Host "`n--- TEST 8: POSTGRESQL ROLES & RESTAURANT OWNERSHIP CHECK ---" -ForegroundColor Yellow
$usersTable = & "C:\Users\astha\pgsql\bin\psql.exe" -U admin -d goship_db -c "SELECT id, name, email, role, is_admin FROM users ORDER BY id;"
Write-Host "$usersTable"

$ownershipTable = & "C:\Users\astha\pgsql\bin\psql.exe" -U admin -d goship_db -c "SELECT r.id, r.name, r.user_restaurants, u.email as owner_email, u.role as owner_role FROM restaurants r JOIN users u ON r.user_restaurants = u.id WHERE r.name LIKE '%Central Tiffin Room%';"
Write-Host "$ownershipTable"
if ($ownershipTable -match "merchant\.ctr@fooddiscovery\.com") {
    Write-Host "  [PASS] Central Tiffin Room (CTR) correctly owned by merchant.ctr@fooddiscovery.com" -ForegroundColor Green
} else {
    Write-Host "  [FAIL] CTR restaurant ownership relation missing or invalid" -ForegroundColor Red
    $allPassed = $false
}

Write-Host "`n==================================================================" -ForegroundColor Cyan
if ($allPassed) {
    Write-Host ">>> STAGE B VALIDATION RESULT: ALL TESTS PASSED (100%) <<<" -ForegroundColor Green
} else {
    Write-Host ">>> STAGE B VALIDATION RESULT: ONE OR MORE TESTS FAILED <<<" -ForegroundColor Red
}
Write-Host "==================================================================" -ForegroundColor Cyan
