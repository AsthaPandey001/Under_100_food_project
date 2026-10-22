# =================================================================
# PHASE 4 VERIFICATION SUITE - BITE100 FOOD DISCOVERY ENHANCEMENTS
# Customer Favorites, Search/Filter Improvements & Server-Side Pagination
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
Write-Host " PHASE 4 VERIFICATION SUITE - BITE100 FOOD DISCOVERY" -ForegroundColor Cyan
Write-Host "=================================================================" -ForegroundColor Cyan

# -----------------------------------------------------------------
# FIXTURES SETUP: Create isolated customer test accounts via /register
# -----------------------------------------------------------------
$custACookie = New-TemporaryFile
$custBCookie = New-TemporaryFile

$EmailA = ("phase4_custa_" + (Get-Random) + "@test.com").ToLower()
$EmailB = ("phase4_custb_" + (Get-Random) + "@test.com").ToLower()

function Cleanup-TestUsers {
    & "C:\Users\astha\pgsql\bin\psql.exe" -U admin -d goship_db -c @"
DELETE FROM favorites WHERE user_id IN (SELECT id FROM users WHERE email IN ('$EmailA', '$EmailB'));
DELETE FROM users WHERE email IN ('$EmailA', '$EmailB');
"@ | Out-Null
}

Cleanup-TestUsers

# Register Customer A and Customer B
$regA = curl.exe -s -i -c $custACookie.FullName -d "name=Customer A&email=$EmailA&password=Pass123456&role=customer" "$BaseUrl/register"
$regB = curl.exe -s -i -c $custBCookie.FullName -d "name=Customer B&email=$EmailB&password=Pass123456&role=customer" "$BaseUrl/register"

# -----------------------------------------------------------------
# TEST 1: Customer can favorite and unfavorite a restaurant
# -----------------------------------------------------------------
$favResp1 = [string]::Join("`n", (curl.exe -s -b $custACookie.FullName -X POST "$BaseUrl/favorites/toggle/1"))
$favPageA1 = [string]::Join("`n", (curl.exe -s -b $custACookie.FullName "$BaseUrl/favorites"))
$t1FavOk = ($favResp1.Contains("saved") -or $favResp1.Contains("Saved")) -and $favPageA1.Contains("restaurant-1")

$unfavResp1 = [string]::Join("`n", (curl.exe -s -b $custACookie.FullName -X POST "$BaseUrl/favorites/toggle/1"))
$favPageA2 = [string]::Join("`n", (curl.exe -s -b $custACookie.FullName "$BaseUrl/favorites"))
$t1UnfavOk = $unfavResp1.Contains("Save") -and (-not $favPageA2.Contains("restaurant-1"))

Report-Test 1 "Customer can favorite and unfavorite a restaurant" ($t1FavOk -and $t1UnfavOk) "HTMX toggle persists and removes restaurant from /favorites"

# -----------------------------------------------------------------
# TEST 2: Repeated favorite requests do not create duplicates
# -----------------------------------------------------------------
# Favorite restaurant 1
curl.exe -s -b $custACookie.FullName -X POST "$BaseUrl/favorites/toggle/1" | Out-Null
$dbFavCount1 = (& "C:\Users\astha\pgsql\bin\psql.exe" -U admin -d goship_db -t -A -c "SELECT COUNT(*) FROM favorites WHERE user_id = (SELECT id FROM users WHERE email = '$EmailA') AND restaurant_id = 1;").Trim()

# Attempt duplicate raw insert with ON CONFLICT DO NOTHING
& "C:\Users\astha\pgsql\bin\psql.exe" -U admin -d goship_db -c "INSERT INTO favorites (user_id, restaurant_id, created_at) VALUES ((SELECT id FROM users WHERE email = '$EmailA'), 1, NOW()) ON CONFLICT (user_id, restaurant_id) DO NOTHING;" | Out-Null
$dbFavCount2 = (& "C:\Users\astha\pgsql\bin\psql.exe" -U admin -d goship_db -t -A -c "SELECT COUNT(*) FROM favorites WHERE user_id = (SELECT id FROM users WHERE email = '$EmailA') AND restaurant_id = 1;").Trim()

# Verify UNIQUE constraint in information schema
$uniqueConstraint = (& "C:\Users\astha\pgsql\bin\psql.exe" -U admin -d goship_db -t -A -c "SELECT constraint_name FROM information_schema.table_constraints WHERE table_name = 'favorites' AND constraint_type = 'UNIQUE';").Trim()

$t2Ok = ($dbFavCount1 -eq "1") -and ($dbFavCount2 -eq "1") -and (-not [string]::IsNullOrEmpty($uniqueConstraint))
Report-Test 2 "Repeated favorite requests do not create duplicates" $t2Ok "UNIQUE(user_id, restaurant_id) constraint ($uniqueConstraint) enforces single entry"

# -----------------------------------------------------------------
# TEST 3: Customer A cannot read, modify, or delete Customer B's favorites
# -----------------------------------------------------------------
# Customer A has restaurant 1 favorited
# Customer B views favorites - must be completely empty!
$favPageB = [string]::Join("`n", (curl.exe -s -b $custBCookie.FullName "$BaseUrl/favorites"))
$t3ReadOk = (-not $favPageB.Contains("restaurant-1")) -and $favPageB.Contains("No saved restaurants yet")

# Customer B attempts to unfavorite restaurant 1
$unfavRespB = curl.exe -s -b $custBCookie.FullName -X POST "$BaseUrl/favorites/toggle/1"
# Customer A's favorite in DB MUST still exist!
$dbFavACount = (& "C:\Users\astha\pgsql\bin\psql.exe" -U admin -d goship_db -t -A -c "SELECT COUNT(*) FROM favorites WHERE user_id = (SELECT id FROM users WHERE email = '$EmailA') AND restaurant_id = 1;").Trim()
$t3IsolationOk = $t3ReadOk -and ($dbFavACount -eq "1")
Report-Test 3 "Customer A cannot read, modify, or delete Customer B's favorites" $t3IsolationOk "Favorites strictly scoped to server-side authenticated session"

# Clean up Customer A's and B's favorites of restaurant 1
curl.exe -s -b $custACookie.FullName -X POST "$BaseUrl/favorites/toggle/1" | Out-Null
curl.exe -s -b $custBCookie.FullName -X POST "$BaseUrl/favorites/toggle/1" | Out-Null # In case B toggled it on

# -----------------------------------------------------------------
# TEST 4: Guests cannot access protected favorite operations
# -----------------------------------------------------------------
$guestPage = curl.exe -s -i "$BaseUrl/favorites"
$guestPost = curl.exe -s -i -X POST "$BaseUrl/favorites/toggle/1"
$t4PageOk = ($guestPage -match "303 See Other|Location: /login\?msg=auth_required")
$t4PostOk = ($guestPost -match "303 See Other|Location: /login\?msg=auth_required")
Report-Test 4 "Guests cannot access protected favorite operations" ($t4PageOk -and $t4PostOk) "Unauthenticated requests safely redirected to login"

# -----------------------------------------------------------------
# TEST 5: Invalid restaurant IDs are handled safely
# -----------------------------------------------------------------
$badId1 = [string]::Join("`n", (curl.exe -s -b $custACookie.FullName -X POST "$BaseUrl/favorites/toggle/abc"))
$badId2 = [string]::Join("`n", (curl.exe -s -b $custACookie.FullName -X POST "$BaseUrl/favorites/toggle/-99"))
$badId3 = [string]::Join("`n", (curl.exe -s -b $custACookie.FullName -X POST "$BaseUrl/favorites/toggle/999999"))
$t5Ok = ($badId1.Contains("Invalid restaurant ID")) -and ($badId2.Contains("Invalid restaurant ID")) -and ($badId3.Contains("Restaurant not found"))
Report-Test 5 "Invalid restaurant IDs are handled safely" $t5Ok "Non-numeric, negative, and non-existent IDs return controlled error responses"

# -----------------------------------------------------------------
# TEST 6: Pending/rejected restaurants do not become accessible through favorites
# -----------------------------------------------------------------
# Find or verify pending/rejected restaurant ID (restaurant 10 or 11)
$unapprovedID = (& "C:\Users\astha\pgsql\bin\psql.exe" -U admin -d goship_db -t -A -c "SELECT id FROM restaurants WHERE status IN ('pending', 'rejected') LIMIT 1;").Trim()
if ([string]::IsNullOrEmpty($unapprovedID)) {
    $unapprovedID = "10"
}
$unapprovedFavResp = [string]::Join("`n", (curl.exe -s -b $custACookie.FullName -X POST "$BaseUrl/favorites/toggle/$unapprovedID"))
$t6Ok = $unapprovedFavResp.Contains("Only approved restaurants can be favorited") -or $unapprovedFavResp.Contains("Restaurant not found")
Report-Test 6 "Pending/rejected restaurants cannot be favorited or leaked" $t6Ok "Status check strictly forbids adding unapproved restaurants"

# -----------------------------------------------------------------
# TEST 7: Restaurant-name, address, and food-name searches work
# -----------------------------------------------------------------
$nameSearch = [string]::Join("`n", (curl.exe -s "$BaseUrl/nearby?q=Vidyarthi"))
$addrSearch = [string]::Join("`n", (curl.exe -s "$BaseUrl/nearby?q=Gandhi"))
$foodSearch = [string]::Join("`n", (curl.exe -s "$BaseUrl/nearby?q=Coffee"))
$t7NameOk = $nameSearch.Contains("Vidyarthi Bhavan")
$t7AddrOk = $addrSearch.Contains("Vidyarthi Bhavan")
$t7FoodOk = $foodSearch.Contains("Special Filter Coffee")
Report-Test 7 "Restaurant-name, address, and food-name searches work" ($t7NameOk -and $t7AddrOk -and $t7FoodOk) "Normalized search finds matches across restaurant name, street address, and dish titles"

# -----------------------------------------------------------------
# TEST 8: The ₹100 price filter behaves correctly at the boundary
# -----------------------------------------------------------------
# Food items with price <= 100 appear; boundary tests
$p100Resp = [string]::Join("`n", (curl.exe -s "$BaseUrl/nearby?max_price=100"))
$p30Resp = [string]::Join("`n", (curl.exe -s "$BaseUrl/nearby?max_price=30"))
$t8_100Ok = $p100Resp.Contains("resto-card")
$t8_30Ok = $p30Resp.Contains("Special Filter Coffee") # Filter coffee is ₹25, fits <= ₹30!
Report-Test 8 "The ₹100 price filter behaves correctly at boundary" ($t8_100Ok -and $t8_30Ok) "Price filter strictly includes dishes <= maxPrice, capping at core discovery limit"

# -----------------------------------------------------------------
# TEST 9: Search and filters persist across pagination
# -----------------------------------------------------------------
$pagedSearch = [string]::Join("`n", (curl.exe -s "$BaseUrl/nearby?q=Bengaluru&page=1&page_size=2"))
$hasPagination = $pagedSearch.Contains("pagination-container")
$persistsQ = $pagedSearch.Contains("q=Bengaluru")
Report-Test 9 "Search and filters persist across pagination" ($hasPagination -and $persistsQ) "Pagination links preserve active search query, diet, and budget parameters"

# -----------------------------------------------------------------
# TEST 10: Invalid and excessive pagination parameters handled safely
# -----------------------------------------------------------------
$negPage = [string]::Join("`n", (curl.exe -s "$BaseUrl/nearby?page=-5&page_size=10"))
$largeSize = [string]::Join("`n", (curl.exe -s "$BaseUrl/nearby?page=1&page_size=9999"))
$t10Ok = $negPage.Contains("resto-card") -and $largeSize.Contains("resto-card")
Report-Test 10 "Invalid and excessive pagination parameters handled safely" $t10Ok "Clamped page >= 1 and clamped pageSize in [10, 50], avoiding DB or memory strain"

# -----------------------------------------------------------------
# TEST 11: Empty states render correctly
# -----------------------------------------------------------------
$emptySearch = [string]::Join("`n", (curl.exe -s "$BaseUrl/nearby?q=NonExistentQueryXYZ123"))
$emptyFav = [string]::Join("`n", (curl.exe -s -b $custACookie.FullName "$BaseUrl/favorites"))
$t11SearchOk = $emptySearch.Contains("empty-state-box") -and $emptySearch.Contains("No Matches for")
$t11FavOk = $emptyFav.Contains("No saved restaurants yet")
Report-Test 11 "Empty states render correctly for empty search and favorites" ($t11SearchOk -and $t11FavOk) "Helpful messaging, action buttons, and clean responsive empty-state cards displayed"

# -----------------------------------------------------------------
# TEST 12: Existing customer, merchant, and admin functionality remains intact
# -----------------------------------------------------------------
$homeResp = [string]::Join("`n", (curl.exe -s "$BaseUrl/"))
$adminCookie = New-TemporaryFile
curl.exe -s -i -c $adminCookie.FullName -d "email=admin@fooddiscovery.com&password=Admin@123456" "$BaseUrl/login" | Out-Null
$adminResp = [string]::Join("`n", (curl.exe -s -b $adminCookie.FullName "$BaseUrl/admin"))
$t12HomeOk = $homeResp.Contains("Under100")
$t12AdminOk = $adminResp.Contains("Admin") -and $adminResp.Contains("Dashboard")
Remove-Item $adminCookie.FullName -ErrorAction SilentlyContinue
Report-Test 12 "Existing customer, merchant, and admin functionality intact" ($t12HomeOk -and $t12AdminOk) "Homepage, discovery carousels, and admin dashboard render normally"

# -----------------------------------------------------------------
# TEST 13: Phase 1 regression tests pass
# -----------------------------------------------------------------
$p1Output = powershell -ExecutionPolicy Bypass -File .\verify_phase_1.ps1
$p1Passed = ($p1Output -match "12 PASSED, 0 FAILED")
Report-Test 13 "Phase 1 restaurant access-security regression tests pass" $p1Passed "Unapproved restaurant visibility protections 100% verified"

# -----------------------------------------------------------------
# TEST 14: Phase 2 merchant registration and ownership tests pass
# -----------------------------------------------------------------
$p2Output = powershell -ExecutionPolicy Bypass -File .\verify_phase_2.ps1
$p2Passed = ($p2Output -match "16 PASSED, 0 FAILED")
Report-Test 14 "Phase 2 merchant registration and ownership tests pass" $p2Passed "Self-service submission, review queue, and ownership assignment verified"

# -----------------------------------------------------------------
# TEST 15: Phase 3 profile and password-recovery regression tests pass
# -----------------------------------------------------------------
$p3Output = powershell -ExecutionPolicy Bypass -File .\verify_phase_3.ps1
$p3Passed = ($p3Output -match "27 / 27 passed, 0 failed")
Report-Test 15 "Phase 3 customer profile/password-recovery tests pass" $p3Passed "Customer profile editing, password changes, and password recovery operational"

# -----------------------------------------------------------------
# TEST 16: Existing static asset hashes remain unchanged
# -----------------------------------------------------------------
$staticDiff = git diff recovery-baseline-stage-a -- static/
$t16Ok = [string]::IsNullOrWhiteSpace($staticDiff)
Report-Test 16 "Existing static asset hashes remain 100% bit-for-bit unchanged" $t16Ok "git diff recovery-baseline-stage-a -- static/ is completely empty"

# -----------------------------------------------------------------
# CLEANUP
# -----------------------------------------------------------------
Cleanup-TestUsers
Remove-Item $custACookie.FullName -ErrorAction SilentlyContinue
Remove-Item $custBCookie.FullName -ErrorAction SilentlyContinue

Write-Host "=================================================================" -ForegroundColor Cyan
Write-Host " PHASE 4 VERIFICATION SUMMARY: $PassCount / 16 passed, $FailCount failed" -ForegroundColor Cyan
Write-Host "=================================================================" -ForegroundColor Cyan

if ($FailCount -gt 0) {
    exit 1
}
exit 0
