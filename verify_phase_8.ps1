# ==============================================================================
# verify_phase_8.ps1 - Phase 8 Restaurant Ratings, Customer Reviews & Moderation Suite
# ==============================================================================

$ErrorActionPreference = "Continue"
$BaseUrl = "http://127.0.0.1:8000"
$PassCount = 0
$FailCount = 0

function Report-Assert($TestName, $Condition, $Detail) {
    if ($Condition) {
        Write-Host "  [PASS] $TestName" -ForegroundColor Green
        $script:PassCount++
    } else {
        Write-Host "  [FAIL] $TestName - $Detail" -ForegroundColor Red
        $script:FailCount++
    }
}

Write-Host ""
Write-Host "=== Phase 8 Verification: Restaurant Ratings, Reviews and Moderation ===" -ForegroundColor Cyan

# STEP 1: Baseline Database Integrity Check
Write-Host ""
Write-Host "-- Step 1: Checking Baseline Database Record Counts --" -ForegroundColor Yellow

$psqlCmd = "C:\Users\astha\pgsql\bin\psql.exe"
$dbCounts = & $psqlCmd -U postgres -d goship_db -t -A -c "SELECT 'users:' || count(*) FROM users UNION ALL SELECT 'restaurants:' || count(*) FROM restaurants UNION ALL SELECT 'food_items:' || count(*) FROM food_items UNION ALL SELECT 'favorites:' || count(*) FROM favorites UNION ALL SELECT 'reviews:' || count(*) FROM reviews UNION ALL SELECT 'review_reports:' || count(*) FROM review_reports;"

$countMap = @{}
foreach ($line in $dbCounts) {
    if ($line -match "([^:]+):(\d+)") {
        $countMap[$matches[1]] = [int]$matches[2]
    }
}

Report-Assert "Baseline Users count is 3" ($countMap["users"] -eq 3) "Found: $($countMap['users'])"
Report-Assert "Baseline Restaurants count is 42" ($countMap["restaurants"] -eq 42) "Found: $($countMap['restaurants'])"
Report-Assert "Baseline Food Items count is 166" ($countMap["food_items"] -eq 166) "Found: $($countMap['food_items'])"
Report-Assert "Baseline Favorites count is 0" ($countMap["favorites"] -eq 0) "Found: $($countMap['favorites'])"
Report-Assert "Baseline Reviews count is 0" ($countMap["reviews"] -eq 0) "Found: $($countMap['reviews'])"
Report-Assert "Baseline Review Reports count is 0" ($countMap["review_reports"] -eq 0) "Found: $($countMap['review_reports'])"

# STEP 2: Authenticate Personas with Cookie Jars
Write-Host ""
Write-Host "-- Step 2: Authenticating Personas --" -ForegroundColor Yellow

$adminCookie = New-TemporaryFile
$cust1Cookie = New-TemporaryFile
$cust2Cookie = New-TemporaryFile

$Email2 = ("phase8_custb_" + (Get-Random) + "@test.com").ToLower()

# Cleanup function
function Cleanup-TestData {
    & $psqlCmd -U postgres -d goship_db -c @"
DELETE FROM review_reports;
DELETE FROM reviews;
DELETE FROM users WHERE email = '$Email2';
"@ | Out-Null
}

Cleanup-TestData

# Log in Admin via quick-login
$adminLogin = curl.exe -s -i -c $adminCookie.FullName -d "role=admin" "$BaseUrl/api/quick-login"
Report-Assert "Admin login successful" ($adminLogin -match "302|303|200") ""

# Log in Customer 1 via quick-login
$cust1Login = curl.exe -s -i -c $cust1Cookie.FullName -d "role=user" "$BaseUrl/api/quick-login"
Report-Assert "Customer 1 login successful" ($cust1Login -match "302|303|200") ""

# Register Customer 2
$cust2Reg = curl.exe -s -i -c $cust2Cookie.FullName -d "name=Customer B&email=$Email2&password=Pass123456&confirm_password=Pass123456" "$BaseUrl/register"
Report-Assert "Customer 2 registration successful" ($cust2Reg -match "302|303|200") ""

# STEP 3: Security & Validation Tests
Write-Host ""
Write-Host "-- Step 3: Security and Input Validation Tests --" -ForegroundColor Yellow

# Unauthenticated review submission
$unauthResp = curl.exe -s -i -d "rating=5&comment=Great dosa!" "$BaseUrl/restaurant/1/reviews"
Report-Assert "Unauthenticated review submission blocked" ($unauthResp -match "302|303|401") ""

# Invalid rating < 1
$invRatingLow = curl.exe -s -i -b $cust1Cookie.FullName -d "rating=0&comment=Invalid rating" "$BaseUrl/restaurant/1/reviews"
Report-Assert "Rating 0 rejected with 400" ($invRatingLow -match "400 Bad Request") ""

# Invalid rating > 5
$invRatingHigh = curl.exe -s -i -b $cust1Cookie.FullName -d "rating=6&comment=Invalid rating" "$BaseUrl/restaurant/1/reviews"
Report-Assert "Rating 6 rejected with 400" ($invRatingHigh -match "400 Bad Request") ""

# Empty comment rejected
$emptyComment = curl.exe -s -i -b $cust1Cookie.FullName -d "rating=4&comment=   " "$BaseUrl/restaurant/1/reviews"
Report-Assert "Empty comment rejected with 400" ($emptyComment -match "400 Bad Request") ""

# Oversized comment rejected
$oversized = "A" * 1005
$oversizedResp = curl.exe -s -i -b $cust1Cookie.FullName -d "rating=4&comment=$oversized" "$BaseUrl/restaurant/1/reviews"
Report-Assert "Oversized comment over 1000 chars rejected with 400" ($oversizedResp -match "400 Bad Request") ""

# Non-existent restaurant rejected
$nonExistentResp = curl.exe -s -i -b $cust1Cookie.FullName -d "rating=5&comment=Valid comment" "$BaseUrl/restaurant/999999/reviews"
Report-Assert "Non-existent restaurant rejected with 404" ($nonExistentResp -match "404 Not Found") ""

# STEP 4: Review Submission & Duplicate Prevention
Write-Host ""
Write-Host "-- Step 4: Review Submission and Duplicate Prevention --" -ForegroundColor Yellow

$reviewResp = curl.exe -s -i -b $cust1Cookie.FullName -d "rating=5&comment=The legendary Benne Masala Dosa is crispy, golden, and worth every rupee!" "$BaseUrl/restaurant/1/reviews"
Report-Assert "Customer 1 review submission succeeded with 200" ($reviewResp -match "200 OK") ""
Report-Assert "Response HTML contains Your Review & Rating" ($reviewResp -match "Your Review") ""
Report-Assert "Response HTML contains 5.0 (1 review)" ($reviewResp -match "5\.0") ""

# Duplicate review prevention
$dupResp = curl.exe -s -i -b $cust1Cookie.FullName -d "rating=4&comment=Trying to submit duplicate" "$BaseUrl/restaurant/1/reviews"
Report-Assert "Duplicate review prevented with 409 Conflict" ($dupResp -match "409 Conflict") ""

# Restaurant modal check
$modalResp = curl.exe -s -b $cust1Cookie.FullName "$BaseUrl/restaurant/1"
Report-Assert "Modal reflects updated rating stats" ($modalResp -match "5\.0") ""
Report-Assert "Modal shows customer existing review with Edit option" ($modalResp -match "Edit") ""

# Retrieve Review ID from DB
$revIDStr = (& $psqlCmd -U postgres -d goship_db -t -A -c "SELECT id FROM reviews WHERE restaurant_id = 1 LIMIT 1").Trim()
$revID = [int]$revIDStr
Report-Assert "Review ID retrieved from database" ($revID -gt 0) "Review ID: $revID"

# STEP 5: IDOR Security Tests
Write-Host ""
Write-Host "-- Step 5: IDOR Protection and Cross-Customer Isolation --" -ForegroundColor Yellow

$idorEditResp = curl.exe -s -i -b $cust2Cookie.FullName -d "rating=1&comment=Hacked by Customer 2" "$BaseUrl/reviews/$revID/edit"
Report-Assert "Cross-customer edit blocked with 403 Forbidden" ($idorEditResp -match "403 Forbidden") ""

$idorDelResp = curl.exe -s -i -b $cust2Cookie.FullName -X DELETE "$BaseUrl/reviews/$revID"
Report-Assert "Cross-customer delete blocked with 403 Forbidden" ($idorDelResp -match "403 Forbidden") ""

# STEP 6: Review Owner Editing
Write-Host ""
Write-Host "-- Step 6: Review Owner Editing --" -ForegroundColor Yellow

$editFormResp = curl.exe -s -b $cust1Cookie.FullName "$BaseUrl/reviews/$revID/edit"
Report-Assert "Customer 1 can load review edit form" ($editFormResp -match "Edit Your Rating") ""

$updateResp = curl.exe -s -i -b $cust1Cookie.FullName -d "rating=4&comment=Updated: Still amazing, though coffee was sweet today." "$BaseUrl/reviews/$revID/edit"
Report-Assert "Customer 1 review update succeeded with 200" ($updateResp -match "200 OK") ""
Report-Assert "Updated rating reflects 4.0" ($updateResp -match "4\.0") ""

# STEP 7: Second Customer Review & Reporting
Write-Host ""
Write-Host "-- Step 7: Second Customer Review and Reporting --" -ForegroundColor Yellow

$cust2RevResp = curl.exe -s -i -b $cust2Cookie.FullName -d "rating=2&comment=Wait time was too long, but food is cheap." "$BaseUrl/restaurant/1/reviews"
Report-Assert "Customer 2 review submitted successfully with 200" ($cust2RevResp -match "200 OK") ""

$statsCheck = curl.exe -s "$BaseUrl/restaurant/1"
Report-Assert "Average rating computed accurately (3.0 across 2 reviews)" ($statsCheck -match "3\.0") ""

$rev2IDStr = (& $psqlCmd -U postgres -d goship_db -t -A -c "SELECT id FROM reviews WHERE restaurant_id = 1 ORDER BY id DESC LIMIT 1").Trim()
$rev2ID = [int]$rev2IDStr

$reportResp = curl.exe -s -i -b $cust1Cookie.FullName -d "reason=Rude language in review" "$BaseUrl/reviews/$rev2ID/report"
Report-Assert "Customer 1 can report a review" ($reportResp -match "200 OK" -and $reportResp -match "Reported") ""

$dupReportResp = curl.exe -s -i -b $cust1Cookie.FullName -d "reason=Spam" "$BaseUrl/reviews/$rev2ID/report"
Report-Assert "Duplicate report handled gracefully" ($dupReportResp -match "200 OK" -and $dupReportResp -match "Already reported") ""

# STEP 8: Admin Moderation Workflow
Write-Host ""
Write-Host "-- Step 8: Admin Moderation Workflow --" -ForegroundColor Yellow

$nonAdminMod = curl.exe -s -i -b $cust1Cookie.FullName "$BaseUrl/admin/reviews"
Report-Assert "Non-admin blocked from /admin/reviews" ($nonAdminMod -match "403 Forbidden|302|303") ""

$adminReviewsResp = curl.exe -s -b $adminCookie.FullName "$BaseUrl/admin/reviews"
Report-Assert "Admin can access Review Moderation page" ($adminReviewsResp -match "Review Moderation") ""
Report-Assert "Admin page displays reported badge" ($adminReviewsResp -match "badge-reports") ""

$adminReportedResp = curl.exe -s -b $adminCookie.FullName "$BaseUrl/admin/reviews?filter=reported"
Report-Assert "Admin reported filter shows reported review" ($adminReportedResp -match "#$rev2ID") ""

$hideResp = curl.exe -s -i -b $adminCookie.FullName -d "reason=Inappropriate complaint" "$BaseUrl/admin/reviews/$rev2ID/status?status=hidden"
Report-Assert "Admin can hide a review" ($hideResp -match "200 OK" -and $hideResp -match "badge-hidden") ""

$publicAfterHide = curl.exe -s "$BaseUrl/restaurant/1"
Report-Assert "Hidden review excluded from public view" (-not ($publicAfterHide -match "Wait time was too long")) ""
Report-Assert "Public rating recalculates to 4.0 after hiding" ($publicAfterHide -match "4\.0") ""

$restoreResp = curl.exe -s -i -b $adminCookie.FullName -d "reason=" "$BaseUrl/admin/reviews/$rev2ID/status?status=approved"
Report-Assert "Admin can restore review" ($restoreResp -match "200 OK" -and $restoreResp -match "badge-approved") ""

$publicAfterRestore = curl.exe -s "$BaseUrl/restaurant/1"
Report-Assert "Restored review visible again in public modal" ($publicAfterRestore -match "Wait time was too long") ""

# STEP 9: Customer Deletion & Reset to Empty
Write-Host ""
Write-Host "-- Step 9: Customer Deletion and Reset to Empty --" -ForegroundColor Yellow

$del2Resp = curl.exe -s -i -b $cust2Cookie.FullName -X DELETE "$BaseUrl/reviews/$rev2ID"
Report-Assert "Customer 2 can delete their own review" ($del2Resp -match "200 OK") ""

$del1Resp = curl.exe -s -i -b $cust1Cookie.FullName -X DELETE "$BaseUrl/reviews/$revID"
Report-Assert "Customer 1 can delete their own review" ($del1Resp -match "200 OK") ""

$emptyModal = curl.exe -s "$BaseUrl/restaurant/1"
Report-Assert "Modal returns to New when all reviews deleted" ($emptyModal -match "New" -and $emptyModal -match "No reviews yet") ""
Report-Assert "Modal shows empty reviews prompt" ($emptyModal -match "No customer reviews yet") ""

# STEP 10: Teardown & Database Record Preservation
Write-Host ""
Write-Host "-- Step 10: Teardown and Database Record Preservation --" -ForegroundColor Yellow

Cleanup-TestData

Remove-Item -Force $adminCookie.FullName, $cust1Cookie.FullName, $cust2Cookie.FullName -ErrorAction SilentlyContinue

$finalCounts = & $psqlCmd -U postgres -d goship_db -t -A -c "SELECT 'users:' || count(*) FROM users UNION ALL SELECT 'restaurants:' || count(*) FROM restaurants UNION ALL SELECT 'food_items:' || count(*) FROM food_items UNION ALL SELECT 'favorites:' || count(*) FROM favorites UNION ALL SELECT 'reviews:' || count(*) FROM reviews UNION ALL SELECT 'review_reports:' || count(*) FROM review_reports;"

$finalCountMap = @{}
foreach ($line in $finalCounts) {
    if ($line -match "([^:]+):(\d+)") {
        $finalCountMap[$matches[1]] = [int]$matches[2]
    }
}

Report-Assert "Final Users count matches baseline (3)" ($finalCountMap["users"] -eq 3) "Found: $($finalCountMap['users'])"
Report-Assert "Final Restaurants count matches baseline (42)" ($finalCountMap["restaurants"] -eq 42) "Found: $($finalCountMap['restaurants'])"
Report-Assert "Final Food Items count matches baseline (166)" ($finalCountMap["food_items"] -eq 166) "Found: $($finalCountMap['food_items'])"
Report-Assert "Final Favorites count matches baseline (0)" ($finalCountMap["favorites"] -eq 0) "Found: $($finalCountMap['favorites'])"
Report-Assert "Final Reviews count is 0" ($finalCountMap["reviews"] -eq 0) "Found: $($finalCountMap['reviews'])"
Report-Assert "Final Review Reports count is 0" ($finalCountMap["review_reports"] -eq 0) "Found: $($finalCountMap['review_reports'])"

# SUMMARY
Write-Host ""
Write-Host "==============================================================================" -ForegroundColor Cyan
Write-Host " Phase 8 Verification Summary: $PassCount PASSED, $FailCount FAILED" -ForegroundColor $(if ($FailCount -eq 0) { "Green" } else { "Red" })
Write-Host "==============================================================================" -ForegroundColor Cyan
Write-Host ""

if ($FailCount -gt 0) {
    exit 1
} else {
    exit 0
}
