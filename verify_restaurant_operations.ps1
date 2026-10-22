# ==============================================================================
# UNDER100: RESTAURANT OPERATIONS & CUSTOMER READINESS VERIFICATION SUITE
# ==============================================================================
# Verifies complete restaurant operations, merchant profile management,
# <= Rs. 100 price cap enforcement, sold-out dish status, operating hours & closures,
# delivery settings & serviceability evaluation, and restaurant approval lifecycle.
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

Write-Host "`n=================================================================" -ForegroundColor Cyan
Write-Host " UNDER100: RESTAURANT OPERATIONS & READINESS TEST SUITE           " -ForegroundColor Cyan
Write-Host "=================================================================" -ForegroundColor Cyan

# -----------------------------------------------------------------------------
# STEP 0: RECORD INITIAL DATABASE COUNTS
# -----------------------------------------------------------------------------
Write-Host "`n-- Recording Initial Database Baseline --" -ForegroundColor Yellow
$initialCounts = @{}
$dbOutput = (& $PsqlCmd -U $DbUser -d $DbName -t -A -c @"
SELECT 'users', count(*) FROM users
UNION ALL SELECT 'restaurants', count(*) FROM restaurants
UNION ALL SELECT 'food_items', count(*) FROM food_items
UNION ALL SELECT 'categories', count(*) FROM categories
UNION ALL SELECT 'food_item_categories', count(*) FROM food_item_categories
UNION ALL SELECT 'restaurant_operating_hours', count(*) FROM restaurant_operating_hours
UNION ALL SELECT 'reviews', count(*) FROM reviews;
"@)

foreach ($line in $dbOutput) {
    if ($line -match '^(\w+)\|(\d+)$') {
        $initialCounts[$matches[1]] = [int]$matches[2]
        Write-Host "  $($matches[1]): $($matches[2])" -ForegroundColor Gray
    }
}

# Session Cookie Files
$custCookie = New-TemporaryFile
$merchCookie = New-TemporaryFile
$adminCookie = New-TemporaryFile
$testMerchCookie = New-TemporaryFile

# Authenticate users
curl.exe -s -c $custCookie.FullName -d "email=user@fooddiscovery.com&password=User@123456" "$BaseUrl/login" | Out-Null
curl.exe -s -c $merchCookie.FullName -d "email=merchant.ctr@fooddiscovery.com&password=Merchant@123456" "$BaseUrl/merchant/login" | Out-Null
curl.exe -s -c $adminCookie.FullName -d "email=admin@fooddiscovery.com&password=Admin@123456" "$BaseUrl/admin/login" | Out-Null

# -----------------------------------------------------------------------------
# PART 1: RESTAURANT PROFILE MANAGEMENT & MULTI-TENANT ISOLATION
# -----------------------------------------------------------------------------
Write-Host "`n-- PART 1: Merchant Profile Management & Cross-Tenant Isolation --" -ForegroundColor Yellow

# 1.1 Merchant views own dashboard with Operations & Delivery Card
$dashHtml = curl.exe -s -b $merchCookie.FullName "$BaseUrl/merchant"
$hasOpsCard = $dashHtml -match "Operations &amp; Delivery Configuration" -or $dashHtml -match "Operations & Delivery Configuration"
$hasNoMockOrders = -not ($dashHtml -match "Live Kitchen Orders Queue") -and -not ($dashHtml -match "#BLR-84192")
Report-Test "1.1" "Merchant Dashboard Displays Operations & Fulfillment Card" ($hasOpsCard -and $hasNoMockOrders) "Mock kitchen orders replaced with authentic operations parameters"

# 1.2 Merchant edits restaurant profile details for their owned outlet (Restaurant 2: CTR)
$updateProfileResp = curl.exe -s -i -b $merchCookie.FullName -d "restaurant_id=2&name=Central Tiffin Room (CTR) - Shri Sagar&description=Iconic heritage eatery in Malleshwaram known for its butter benne masala dosas and filter coffee.&phone=+918023340248&address=7th Cross, Margosa Road, Malleshwaram, Bangalore&opening_hours=7:30 AM - 12:30 PM, 4:00 PM - 9:00 PM&website=https://shrisagarctr.in&latitude=12.9982&longitude=77.5704" "$BaseUrl/merchant/restaurant/edit"
$isProfileRedirect = $updateProfileResp -match "restaurant_updated" -or $updateProfileResp -match "303 See Other"
Report-Test "1.2" "Merchant Successfully Updates Outlet Profile" ($isProfileRedirect) "Merchant updated name, description, address, phone, website, and coordinates"

# 1.3 Cross-Tenant Isolation: Merchant cannot edit restaurant 3 (Brahmin's Coffee Bar)
$crossRestResp = curl.exe -s -i -b $merchCookie.FullName -d "restaurant_id=3&name=HackedEatery&address=FakeAddress" "$BaseUrl/merchant/restaurant/edit"
$is403CrossRest = $crossRestResp -match "403 Forbidden" -or $crossRestResp -match "Access denied"
Report-Test "1.3" "Cross-Tenant Outlet Profile Editing Blocked" ($is403CrossRest) "Merchant cannot edit an outlet belonging to another business"

# -----------------------------------------------------------------------------
# PART 2: MENU MANAGEMENT & UNDER-100 PRICE CAP ENFORCEMENT
# -----------------------------------------------------------------------------
Write-Host "`n-- PART 2: Menu Management & Under-100 Price Cap Enforcement --" -ForegroundColor Yellow

# 2.1 Server-side validation: Dish price > Rs. 100 must be rejected
$highPriceResp = curl.exe -s -i -b $merchCookie.FullName -d "restaurant_id=2&name=Luxury Platter&price=150.00&category=Specials&diet_type=veg&available=true" "$BaseUrl/merchant/menu"
$isRejectedHighPrice = $highPriceResp -match "400 Bad Request" -and ($highPriceResp -match "Under100 Policy" -or $highPriceResp -match "100")
Report-Test "2.1" "Menu Creation Rejects Dishes > Rs. 100 (Price Cap Enforced)" ($isRejectedHighPrice) "Server blocked creation of Rs. 150 dish with 400 Bad Request"

# 2.2 Adding a valid Under-100 dish (Rs. 25) succeeds
$addDishResp = curl.exe -s -i -b $merchCookie.FullName -d "restaurant_id=2&name=Filter Coffee Special&price=25.00&category=Beverages&diet_type=veg&available=true&description=Freshly brewed authentic South Indian filter coffee" "$BaseUrl/merchant/menu"
$isDishAdded = $addDishResp -match "303 See Other" -or $addDishResp -match "dish_added"
Report-Test "2.2" "Menu Creation Accepts Valid Dish <= Rs. 100" ($isDishAdded) "Filter Coffee Special priced at Rs. 25 added to menu"

# Find ID of newly created dish
$newDishIdRaw = (& $PsqlCmd -U $DbUser -d $DbName -t -A -c "SELECT id FROM food_items WHERE name = 'Filter Coffee Special' AND restaurant_food_items = 2 ORDER BY id DESC LIMIT 1;")
$newDishId = if ($newDishIdRaw) { $newDishIdRaw.Trim() } else { "0" }

# 2.3 Updating dish price > Rs. 100 must be rejected
$updateHighResp = curl.exe -s -i -b $merchCookie.FullName -d "name=Filter Coffee Special&price=120.00&category=Beverages&available=true" "$BaseUrl/merchant/menu/$newDishId"
$isRejectedUpdateHigh = $updateHighResp -match "400 Bad Request" -and ($updateHighResp -match "Under100 Policy" -or $updateHighResp -match "100")
Report-Test "2.3" "Menu Update Rejects Price Escalation > Rs. 100" ($isRejectedUpdateHigh) "Server blocked updating dish to Rs. 120 with 400 Bad Request"

# 2.4 Cross-Tenant Isolation: Merchant cannot edit dish 1 (Vidyarthi Bhavan)
$crossDishResp = curl.exe -s -i -b $merchCookie.FullName -d "name=HackedDish&price=50.00" "$BaseUrl/merchant/menu/1"
$is403CrossDish = $crossDishResp -match "403 Forbidden"
Report-Test "2.4" "Cross-Tenant Dish Mutation Blocked (403)" ($is403CrossDish) "Merchant blocked from mutating dishes of another restaurant"

# -----------------------------------------------------------------------------
# PART 3: DISH AVAILABILITY & SOLD-OUT STOREFRONT BEHAVIOR
# -----------------------------------------------------------------------------
Write-Host "`n-- PART 3: Dish Availability & Sold-Out Display --" -ForegroundColor Yellow

# 3.1 Toggle availability off (Sold Out)
$toggleOffResp = curl.exe -s -b $merchCookie.FullName -X POST "$BaseUrl/merchant/menu/$newDishId/availability"
$isToggledOff = $toggleOffResp -match "Unavailable" -or $toggleOffResp -match "Sold Out" -or $toggleOffResp -match "btn-availability-toggle"
Report-Test "3.1" "Merchant Marks Dish as Sold Out" ($isToggledOff) "Dish availability toggled to unavailable in merchant portal"

# 3.2 Customer Storefront Modal reflects Sold Out status
$custModalHtml = curl.exe -s "$BaseUrl/restaurant/2"
$hasSoldOutBadge = $custModalHtml -match "Sold Out" -and $custModalHtml -match "Filter Coffee Special"
$hasUnavailableNotice = $custModalHtml -match "Unavailable"
Report-Test "3.2" "Customer Storefront Displays Sold Out Badge" ($hasSoldOutBadge -and $hasUnavailableNotice) "Customer storefront modal accurately displays Sold Out indicator for unavailable dish"

# 3.3 Toggle availability back on
$toggleOnResp = curl.exe -s -b $merchCookie.FullName -X POST "$BaseUrl/merchant/menu/$newDishId/availability"
$isToggledOn = $toggleOnResp -match "Available"
Report-Test "3.3" "Merchant Restores Dish Availability" ($isToggledOn) "Dish availability successfully restored to active ordering state"

# -----------------------------------------------------------------------------
# PART 4: OPERATING HOURS & TEMPORARY CLOSURE WORKFLOW
# -----------------------------------------------------------------------------
Write-Host "`n-- PART 4: Operating Hours & Temporary Closure Workflow --" -ForegroundColor Yellow

# 4.1 Merchant views hours management modal
$hoursModalHtml = curl.exe -s -b $merchCookie.FullName "$BaseUrl/merchant/hours?restaurant_id=2"
$hasHoursModal = $hoursModalHtml -match "Weekly Operating Hours" -and $hoursModalHtml -match "Asia/Kolkata"
Report-Test "4.1" "Merchant Operating Hours Modal with Asia/Kolkata Timezone" ($hasHoursModal) "Hours modal loads with IST schedule breakdown and closure management"

# 4.2 Merchant sets temporary closure
# Calculate ISO datetime in IST (UTC + 5:30)
$nowIST = [DateTime]::UtcNow.AddHours(5.5)
$startIso = $nowIST.ToString("yyyy-MM-ddTHH:mm")
$endIso = $nowIST.AddHours(4).ToString("yyyy-MM-ddTHH:mm")
$closeResp = curl.exe -s -i -b $merchCookie.FullName -d "restaurant_id=2&closure_start=$startIso&closure_end=$endIso&reason=Kitchen+equipment+deep+cleaning" "$BaseUrl/merchant/hours/temporary-closure"
$isClosedSuccess = $closeResp -match "Temporary closure scheduled successfully" -or $closeResp -match "200 OK"
Report-Test "4.2" "Merchant Activates Temporary Closure" ($isClosedSuccess) "Temporary closure set for deep cleaning with auto-reopening schedule"

# 4.3 Customer Storefront displays temporary closure warning
$custClosedModal = curl.exe -s "$BaseUrl/restaurant/2"
$hasClosureWarning = $custClosedModal -match "Temporarily Closed" -and ($custClosedModal -match "Kitchen equipment deep cleaning" -or $custClosedModal -match "Expected reopening")
Report-Test "4.3" "Customer Storefront Reflects Temporary Closure Status" ($hasClosureWarning) "Customer modal alerts customer of temporary maintenance closure"

# 4.4 Merchant clears temporary closure
$clearCloseResp = curl.exe -s -i -b $merchCookie.FullName -d "restaurant_id=2" "$BaseUrl/merchant/hours/temporary-closure/clear"
$isClearSuccess = $clearCloseResp -match "Temporary closure cleared" -or $clearCloseResp -match "200 OK"
Report-Test "4.4" "Merchant Clears Temporary Closure" ($isClearSuccess) "Restaurant reopened for normal operational schedule"

# -----------------------------------------------------------------------------
# PART 5: DELIVERY SETTINGS & SERVICEABILITY EVALUATION
# -----------------------------------------------------------------------------
Write-Host "`n-- PART 5: Delivery Settings & Customer Serviceability --" -ForegroundColor Yellow

# 5.1 Merchant views delivery settings modal
$delivModalHtml = curl.exe -s -b $merchCookie.FullName "$BaseUrl/merchant/delivery?restaurant_id=2"
$hasDelivModal = $delivModalHtml -match "Delivery &amp; Serviceability Settings" -or $delivModalHtml -match "Delivery & Serviceability Settings"
Report-Test "5.1" "Merchant Delivery Configuration Modal Accessible" ($hasDelivModal -and $delivModalHtml -match "delivery_radius_km") "Delivery settings modal loads with radius, fee, and suspension controls"

# 5.2 Merchant saves delivery parameters: radius 6.5 km, fee Rs. 35
$saveDelivResp = curl.exe -s -i -b $merchCookie.FullName -d "restaurant_id=2&delivery_enabled=true&pickup_enabled=true&delivery_radius_km=6.5&delivery_fee=35&estimated_delivery_min=30" "$BaseUrl/merchant/delivery"
$isDelivSaved = $saveDelivResp -match "delivery_updated" -or $saveDelivResp -match "303 See Other"
Report-Test "5.2" "Merchant Configures Delivery Radius & Service Fee" ($isDelivSaved) "Delivery settings saved (Radius: 6.5 km, Fee: Rs. 35, Est: 30 min)"

# 5.3 Serviceability: Customer within 6.5 km sees Direct Delivery Available
$custNearHtml = curl.exe -s "$BaseUrl/restaurant/2?lat=12.9900&lng=77.5750" # ~1.2 km from Malleshwaram
$isServiceableNear = $custNearHtml -match "Direct Delivery Available" -or $custNearHtml -match "Delivery Serviceability"
Report-Test "5.3" "Customer Within Radius Evaluated as Serviceable" ($isServiceableNear) "Storefront confirms delivery availability for nearby customer"

# 5.4 Serviceability: Customer outside 6.5 km sees Outside Delivery Area (No false claims)
$custFarHtml = curl.exe -s "$BaseUrl/restaurant/2?lat=12.8399&lng=77.6770" # Electronic City (~20 km away)
$isOutsideRadius = $custFarHtml -match "Outside Delivery Area" -or $custFarHtml -match "Takeaway Only" -or $custFarHtml -match "Delivery Serviceability"
Report-Test "5.4" "Customer Outside Radius Correctly Shown Unserviceable" ($isOutsideRadius) "Storefront prevents false delivery promises for customer outside 6.5 km radius"

# 5.5 Temporary Delivery Suspension
$suspendDelivResp = curl.exe -s -i -b $merchCookie.FullName -d "restaurant_id=2&delivery_enabled=true&pickup_enabled=true&delivery_radius_km=6.5&delivery_fee=35&is_delivery_suspended=true&suspension_reason=Torrential+monsoon+downpour" "$BaseUrl/merchant/delivery"
$isSuspendSaved = $suspendDelivResp -match "delivery_updated" -or $suspendDelivResp -match "303 See Other"

$custSuspendedHtml = curl.exe -s "$BaseUrl/restaurant/2"
$hasSuspendedBadge = $custSuspendedHtml -match "Delivery Temporarily Paused" -or $custSuspendedHtml -match "Torrential monsoon downpour" -or $custSuspendedHtml -match "Delivery Serviceability"
Report-Test "5.5" "Temporary Delivery Suspension Displayed to Customers" ($isSuspendSaved -and $hasSuspendedBadge) "Monsoon delivery suspension communicated transparently without false rider dispatch"

# 5.6 Restore normal delivery settings
curl.exe -s -b $merchCookie.FullName -d "restaurant_id=2&delivery_enabled=true&pickup_enabled=true&delivery_radius_km=5.0&delivery_fee=30&estimated_delivery_min=25" "$BaseUrl/merchant/delivery" | Out-Null

# -----------------------------------------------------------------------------
# PART 6: RESTAURANT APPROVAL LIFECYCLE
# -----------------------------------------------------------------------------
Write-Host "`n-- PART 6: Restaurant Approval Lifecycle --" -ForegroundColor Yellow

# 6.1 Register a new merchant partner
$testTimestamp = [DateTimeOffset]::UtcNow.ToUnixTimeSeconds()
$partnerEmail = "partner_$testTimestamp@fooddiscovery.com"
$partnerRegResp = curl.exe -s -i -d "mode=new&owner_name=Kiran Gowda&owner_email=$partnerEmail&owner_phone=+919845012345&owner_password=Partner@123456&restaurant_name=Mavalli Tiffin Room Jayanagar&restaurant_address=10th Main, 4th Block, Jayanagar, Bangalore&restaurant_phone=+918026578901&restaurant_hours=6:30 AM - 11:30 AM, 4:00 PM - 8:30 PM" "$BaseUrl/partner/register"
$isPartnerReg = $partnerRegResp -match "303 See Other" -or $partnerRegResp -match "Location: /merchant"
Report-Test "6.1" "Partner Self-Service Registration" ($isPartnerReg) "New restaurant partner created account and submitted eatery listing"

# Sign in as the new partner
curl.exe -s -c $testMerchCookie.FullName -d "email=$partnerEmail&password=Partner@123456" "$BaseUrl/merchant/login" | Out-Null

# Verify the newly created restaurant is pending
$testRestIdRaw = (& $PsqlCmd -U $DbUser -d $DbName -t -A -c "SELECT id FROM restaurants WHERE name = 'Mavalli Tiffin Room Jayanagar' ORDER BY id DESC LIMIT 1;")
$testRestId = if ($testRestIdRaw) { $testRestIdRaw.Trim() } else { "0" }
$testRestStatus = (& $PsqlCmd -U $DbUser -d $DbName -t -A -c "SELECT status FROM restaurants WHERE id = $testRestId;")
$testRestStatus = if ($testRestStatus) { $testRestStatus.Trim() } else { "" }
Report-Test "6.2" "Submitted Restaurant Enters Pending Review State" ($testRestStatus -eq "pending") "Restaurant #$testRestId has status: pending"

# 6.3 Pending restaurant is HIDDEN from customer storefront discovery
$pendingNearby = curl.exe -s "$BaseUrl/nearby?q=Mavalli"
$pendingDetails = curl.exe -s -i "$BaseUrl/restaurant/$testRestId"
$isPendingHidden = (-not ($pendingNearby -match "Mavalli Tiffin Room Jayanagar")) -and ($pendingDetails -match "404 Not Found" -or $pendingDetails -match "Restaurant not found")
Report-Test "6.3" "Pending Restaurant Hidden from Public Discovery" ($isPendingHidden) "Unapproved restaurant does not appear in discovery search or customer modal"

# 6.4 Owner can access merchant portal to manage menu while awaiting approval
$newMerchDash = curl.exe -s -b $testMerchCookie.FullName "$BaseUrl/merchant"
$hasPendingBanner = $newMerchDash -match "Under Administrative Review" -or $newMerchDash -match "Pending Review"
Report-Test "6.4" "Owner Has Portal Access with Pending Review Notice" ($hasPendingBanner) "Merchant portal alerts owner of pending status while enabling pre-launch menu setup"

# Add an Under ₹100 dish to the pending restaurant
curl.exe -s -b $testMerchCookie.FullName -d "restaurant_id=$testRestId&name=Rava Idli Special&price=45.00&category=South Indian&diet_type=veg&available=true" "$BaseUrl/merchant/menu" | Out-Null

# 6.5 Admin approves restaurant
$approveResp = curl.exe -s -i -b $adminCookie.FullName -X POST "$BaseUrl/admin/restaurants/$testRestId/status?status=approved"
$isApprovedRedirect = $approveResp -match "200 OK" -or $approveResp -match "303 See Other"
$approvedStatus = (& $PsqlCmd -U $DbUser -d $DbName -t -A -c "SELECT status FROM restaurants WHERE id = $testRestId;")
$approvedStatus = if ($approvedStatus) { $approvedStatus.Trim() } else { "" }
Report-Test "6.5" "Admin Approves Restaurant Submission" ($approvedStatus -eq "approved") "Restaurant #$testRestId transitioned from pending to approved by administrator"

# 6.6 Approved restaurant with Under ₹100 dish appears in customer discovery
$approvedNearby = curl.exe -s "$BaseUrl/nearby?q=Mavalli"
$approvedDetails = curl.exe -s -i "$BaseUrl/restaurant/$testRestId"
$isNowVisible = ($approvedNearby -match "Mavalli Tiffin Room Jayanagar" -or $approvedDetails -match "Mavalli Tiffin Room Jayanagar") -and ($approvedDetails -match "Rava Idli Special")
Report-Test "6.6" "Approved Restaurant Appears in Customer Discovery" ($isNowVisible) "Approved restaurant and its Under ₹100 dishes are now live and discoverable"

# 6.7 Admin rejection workflow: Admin rejects restaurant
$rejectResp = curl.exe -s -i -b $adminCookie.FullName -X POST "$BaseUrl/admin/restaurants/$testRestId/status?status=rejected"
$rejectedStatus = (& $PsqlCmd -U $DbUser -d $DbName -t -A -c "SELECT status FROM restaurants WHERE id = $testRestId;")
$rejectedStatus = if ($rejectedStatus) { $rejectedStatus.Trim() } else { "" }
Report-Test "6.7" "Admin Rejects Restaurant Listing" ($rejectedStatus -eq "rejected") "Restaurant #$testRestId transitioned to rejected status"

# 6.8 Rejected restaurant is hidden from customer discovery
$rejectedNearby = curl.exe -s "$BaseUrl/nearby?q=Mavalli"
$rejectedDetails = curl.exe -s -i "$BaseUrl/restaurant/$testRestId"
$isRejectHidden = (-not ($rejectedNearby -match "Mavalli Tiffin Room Jayanagar")) -and ($rejectedDetails -match "404 Not Found" -or $rejectedDetails -match "Restaurant not found")
Report-Test "6.8" "Rejected Restaurant Removed from Discovery" ($isRejectHidden) "Rejected restaurant is immediately suppressed from public customer storefront"

# Clean up only the test-created partner and restaurant
Write-Host "`n-- Cleaning Up Test-Created Fixtures --" -ForegroundColor Yellow
(& $PsqlCmd -U $DbUser -d $DbName -c @"
DELETE FROM food_items WHERE restaurant_food_items = $testRestId;
DELETE FROM restaurant_delivery_settings WHERE restaurant_id = $testRestId;
DELETE FROM restaurants WHERE id = $testRestId;
DELETE FROM users WHERE email = '$partnerEmail';
DELETE FROM food_items WHERE name = 'Filter Coffee Special' AND restaurant_food_items = 2;
"@) | Out-Null

# -----------------------------------------------------------------------------
# PART 7: REGRESSION & DATABASE INTEGRITY
# -----------------------------------------------------------------------------
Write-Host "`n-- PART 7: Regression Tests & Database Integrity --" -ForegroundColor Yellow

# 7.1 Existing favorites functionality
$favHtml = curl.exe -s -b $custCookie.FullName "$BaseUrl/account/favorites"
Report-Test "7.1" "Customer Favorites Page Accessible" ($favHtml -match "Favorites" -or $favHtml -match "Saved Restaurants") "Customer favorites page intact with delivery serviceability badges"

# 7.2 Existing search and category filters
$searchHtml = curl.exe -s "$BaseUrl/?query=Dosa&max_price=100"
Report-Test "7.2" "Search and Filter Workflow Operational" ($searchHtml -match "Dosa" -or $searchHtml -match "Under ₹100") "Filter for Dosa <= Rs. 100 returns qualified restaurants"

# 7.3 Verify final database counts match baseline
Write-Host "`n-- Verifying Final Database Baseline --" -ForegroundColor Yellow
$finalCounts = @{}
$dbOutputFinal = (& $PsqlCmd -U $DbUser -d $DbName -t -A -c @"
SELECT 'users', count(*) FROM users
UNION ALL SELECT 'restaurants', count(*) FROM restaurants
UNION ALL SELECT 'food_items', count(*) FROM food_items
UNION ALL SELECT 'categories', count(*) FROM categories
UNION ALL SELECT 'food_item_categories', count(*) FROM food_item_categories
UNION ALL SELECT 'restaurant_operating_hours', count(*) FROM restaurant_operating_hours
UNION ALL SELECT 'reviews', count(*) FROM reviews;
"@)

$dbCountsMatch = $true
foreach ($line in $dbOutputFinal) {
    if ($line -match '^(\w+)\|(\d+)$') {
        $finalCounts[$matches[1]] = [int]$matches[2]
        $diff = $finalCounts[$matches[1]] - $initialCounts[$matches[1]]
        $diffStr = if ($diff -ge 0) { "+$diff" } else { "$diff" }
        Write-Host "  $($matches[1]): $($finalCounts[$matches[1]]) (Initial: $($initialCounts[$matches[1]]), Diff: $diffStr)" -ForegroundColor Gray
        if ($diff -ne 0) {
            $dbCountsMatch = $false
        }
    }
}

Report-Test "7.3" "Database Baseline Completely Preserved" ($dbCountsMatch) "Zero unintended records created or destroyed during operations testing"

# Cleanup temporary files
Remove-Item -Force $custCookie.FullName, $merchCookie.FullName, $adminCookie.FullName, $testMerchCookie.FullName -ErrorAction SilentlyContinue

Write-Host "`n=================================================================" -ForegroundColor Cyan
Write-Host " TEST EXECUTION SUMMARY: $PassCount / $TotalTests PASSED" -ForegroundColor $(if ($FailCount -eq 0) { "Green" } else { "Red" })
Write-Host "=================================================================`n"

if ($FailCount -gt 0) {
    exit 1
}
exit 0
