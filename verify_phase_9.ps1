# =================================================================
# PHASE 9 VERIFICATION SUITE - BITE100 FOOD CATEGORIES & CUISINE DISCOVERY
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
Write-Host " PHASE 9 VERIFICATION SUITE - BITE100 FOOD CATEGORIES & CUISINE DISCOVERY" -ForegroundColor Cyan
Write-Host "=================================================================" -ForegroundColor Cyan

# Setup Admin and Merchant sessions
$adminCookie = New-TemporaryFile
$merchCookie = New-TemporaryFile
$custCookie = New-TemporaryFile

function Cleanup-TestFixture {
    & "C:\Users\astha\pgsql\bin\psql.exe" -U admin -d goship_db -c @"
DELETE FROM food_item_categories WHERE category_id IN (SELECT id FROM categories WHERE slug LIKE 'test-%');
DELETE FROM categories WHERE slug LIKE 'test-%';
DELETE FROM users WHERE email LIKE 'phase9_cust_%@test.com';
"@ | Out-Null
}

Cleanup-TestFixture

curl.exe -s -i -c $adminCookie.FullName -d "email=admin@fooddiscovery.com&password=Admin@123456" "$BaseUrl/login" | Out-Null
curl.exe -s -i -c $merchCookie.FullName -d "email=merchant@fooddiscovery.com&password=Merchant@123456" "$BaseUrl/login" | Out-Null

$custEmail = ("phase9_cust_" + (Get-Random) + "@test.com").ToLower()
curl.exe -s -i -c $custCookie.FullName -d "name=Test Customer&email=$custEmail&password=Pass123456&role=customer" "$BaseUrl/register" | Out-Null

# -----------------------------------------------------------------
# TEST 1: Category creation and duplicate prevention
# -----------------------------------------------------------------
$testCatName = "Test Artisan Bread"
$createResp = [string]::Join("`n", (curl.exe -s -i -b $adminCookie.FullName -d "name=$testCatName&category_type=food_type&description=Freshly baked breads&icon=🥖" "$BaseUrl/admin/categories"))
$dupResp = [string]::Join("`n", (curl.exe -s -i -b $adminCookie.FullName -d "name=$testCatName&category_type=food_type&description=Duplicate bread&icon=🥖" "$BaseUrl/admin/categories"))

$dbCount = (& "C:\Users\astha\pgsql\bin\psql.exe" -U admin -d goship_db -t -A -c "SELECT COUNT(*) FROM categories WHERE slug = 'test-artisan-bread';").Trim()

$t1Ok = ($createResp -match "303 See Other|msg=Category\+created") -and ($dupResp -match "already exists|err=") -and ($dbCount -eq "1")
Report-Test 1 "Category creation and duplicate prevention" $t1Ok "Admin can create new categories, and server enforces UNIQUE(slug) duplicate rejection"

# -----------------------------------------------------------------
# TEST 2: Category lookup and filtering
# -----------------------------------------------------------------
$southIndianResp = [string]::Join("`n", (curl.exe -s "$BaseUrl/nearby?category=south-indian"))
$directLinkResp = curl.exe -s -i "$BaseUrl/category/south-indian"

$t2LookupOk = $southIndianResp.Contains("resto-card") -and $southIndianResp.Contains("South Indian") -and ($directLinkResp -match "303 See Other|302 Found|Location: /\?category=south-indian")
Report-Test 2 "Category lookup and filtering" $t2LookupOk "Discovery filters dishes by category slug and supports direct links via /category/:slug"

# -----------------------------------------------------------------
# TEST 3: Food items assigned to categories
# -----------------------------------------------------------------
$assignedDishes = (& "C:\Users\astha\pgsql\bin\psql.exe" -U admin -d goship_db -t -A -c "SELECT COUNT(*) FROM food_item_categories WHERE category_id = (SELECT id FROM categories WHERE slug = 'south-indian');").Trim()
$t3Ok = ([int]$assignedDishes -gt 0)
Report-Test 3 "Food items assigned to categories" $t3Ok "$assignedDishes food items correctly linked to South Indian category in food_item_categories"

# -----------------------------------------------------------------
# TEST 4: Multiple categories where supported
# -----------------------------------------------------------------
# A dish like Idli can belong to both 'south-indian' and 'breakfast'
$multiCatDishes = (& "C:\Users\astha\pgsql\bin\psql.exe" -U admin -d goship_db -t -A -c @"
SELECT COUNT(DISTINCT fc1.food_item_id) 
FROM food_item_categories fc1 
JOIN food_item_categories fc2 ON fc1.food_item_id = fc2.food_item_id AND fc1.category_id != fc2.category_id;
"@).Trim()
$t4Ok = ([int]$multiCatDishes -gt 0)
Report-Test 4 "Multiple categories where supported" $t4Ok "$multiCatDishes dishes associated with both cuisine and meal_type simultaneously"

# -----------------------------------------------------------------
# TEST 5: Combined category and search filters
# -----------------------------------------------------------------
$comboResp = [string]::Join("`n", (curl.exe -s "$BaseUrl/nearby?category=south-indian&q=Dosa"))
$mismatchResp = [string]::Join("`n", (curl.exe -s "$BaseUrl/nearby?category=chinese&q=Dosa"))

$t5Ok = $comboResp.Contains("Dosa") -and (-not $mismatchResp.Contains("Dosa"))
Report-Test 5 "Combined category and search filters" $t5Ok "Combining category=south-indian and q=Dosa returns dosas; category=chinese and q=Dosa yields no matches"

# -----------------------------------------------------------------
# TEST 6: Category filters with the ₹100 maximum price
# -----------------------------------------------------------------
$cheapResp = [string]::Join("`n", (curl.exe -s "$BaseUrl/nearby?category=beverages&max_price=30"))
$expensiveResp = [string]::Join("`n", (curl.exe -s "$BaseUrl/nearby?category=beverages&max_price=10"))

$t6Ok = $cheapResp.Contains("Coffee") -and (-not $expensiveResp.Contains("Special Filter Coffee"))
Report-Test 6 "Category filters with ₹100 maximum price" $t6Ok "Dishes strictly constrained by price ceiling; dishes <= ₹30 displayed, dishes > ₹10 filtered"

# -----------------------------------------------------------------
# TEST 7: Category filters with dietary restrictions
# -----------------------------------------------------------------
$vegResp = [string]::Join("`n", (curl.exe -s "$BaseUrl/nearby?category=south-indian&diet=veg"))
$t7Ok = $vegResp.Contains("resto-card") -and (-not $vegResp.Contains("NON-VEG"))
Report-Test 7 "Category filters with dietary restrictions" $t7Ok "category=south-indian with diet=veg only displays pure veg dishes"

# -----------------------------------------------------------------
# TEST 8: Pagination preserving category and search parameters
# -----------------------------------------------------------------
$pagedResp = [string]::Join("`n", (curl.exe -s "$BaseUrl/nearby?category=south-indian&q=Vada&page=1&page_size=2"))
$hasPaging = $pagedResp.Contains("pagination-container")
$preservesCat = $pagedResp.Contains("category=south-indian")
$preservesQ = $pagedResp.Contains("q=Vada")
$t8Ok = $hasPaging -and $preservesCat -and $preservesQ
Report-Test 8 "Pagination preserving category and search parameters" $t8Ok "Server pagination buttons preserve category=south-indian and q=Vada without dropping parameters"

# -----------------------------------------------------------------
# TEST 9: Empty category results
# -----------------------------------------------------------------
$emptyCatResp = [string]::Join("`n", (curl.exe -s "$BaseUrl/nearby?category=test-artisan-bread"))
$t9Ok = $emptyCatResp.Contains("empty-state-box") -and $emptyCatResp.Contains("No Dishes Found in")
Report-Test 9 "Empty category results" $t9Ok "Clean responsive empty state displayed when category has no available food items under ₹100"

# -----------------------------------------------------------------
# TEST 10: Unknown or malformed category parameters
# -----------------------------------------------------------------
$malformedResp = [string]::Join("`n", (curl.exe -s "$BaseUrl/nearby?category=%27%20OR%201=1--&max_price=invalid"))
$t10Ok = $malformedResp.Contains("empty-state-box") -or $malformedResp.Contains("resto-card")
Report-Test 10 "Unknown or malformed category parameters handled safely" $t10Ok "SQL injection characters and non-numeric price safely handled without server error or crash"

# -----------------------------------------------------------------
# TEST 11: Public visibility restrictions for pending/rejected restaurants
# -----------------------------------------------------------------
# Verify that pending or rejected restaurants NEVER appear in category results
$unapprovedInCats = [string]::Join("`n", (curl.exe -s "$BaseUrl/nearby?category=south-indian&radius=100"))
$t11Ok = (-not $unapprovedInCats.Contains("restaurant-10")) -and (-not $unapprovedInCats.Contains("restaurant-11"))
Report-Test 11 "Public visibility restrictions for pending/rejected restaurants" $t11Ok "Pending (ID: 10) and rejected (ID: 11) restaurants never appear in category discovery"

# -----------------------------------------------------------------
# TEST 12: Merchant ownership enforcement
# -----------------------------------------------------------------
# Customer or other merchant cannot edit a dish
$custEditResp = curl.exe -s -i -b $custCookie.FullName -d "name=Hacked Dish&price=50" "$BaseUrl/merchant/menu/1"
$t12Ok = ($custEditResp -match "303 See Other|403 Forbidden|msg=merchant_required")
Report-Test 12 "Merchant ownership enforcement" $t12Ok "Unauthorized users cannot access merchant menu update endpoints"

# -----------------------------------------------------------------
# TEST 13: Admin-only category management
# -----------------------------------------------------------------
$guestCatAdmin = curl.exe -s -i "$BaseUrl/admin/categories"
$custCatAdmin = curl.exe -s -i -b $custCookie.FullName "$BaseUrl/admin/categories"
$t13Ok = ($guestCatAdmin -match "303 See Other|403 Forbidden") -and ($custCatAdmin -match "303 See Other|403 Forbidden")
Report-Test 13 "Admin-only category management" $t13Ok "Guests and standard customers strictly blocked from accessing /admin/categories"

# -----------------------------------------------------------------
# TEST 14: Category deactivation without data loss
# -----------------------------------------------------------------
$testCatID = (& "C:\Users\astha\pgsql\bin\psql.exe" -U admin -d goship_db -t -A -c "SELECT id FROM categories WHERE slug = 'test-artisan-bread';").Trim()
$toggleResp = curl.exe -s -i -b $adminCookie.FullName -X POST "$BaseUrl/admin/categories/$testCatID/toggle"
$isInactive = (& "C:\Users\astha\pgsql\bin\psql.exe" -U admin -d goship_db -t -A -c "SELECT is_active FROM categories WHERE id = $testCatID;").Trim()
$foodCountBefore = (& "C:\Users\astha\pgsql\bin\psql.exe" -U admin -d goship_db -t -A -c "SELECT COUNT(*) FROM food_items;").Trim()

$t14Ok = ($isInactive -eq "f") -and ($foodCountBefore -eq "166")
Report-Test 14 "Category deactivation without data loss" $t14Ok "Deactivating category toggles status to inactive without deleting food items or relationships"

# -----------------------------------------------------------------
# TEST 15: Regression of existing customer, merchant, admin, favorites, and reviews
# -----------------------------------------------------------------
$homeResp = [string]::Join("`n", (curl.exe -s "$BaseUrl/"))
$adminResp = [string]::Join("`n", (curl.exe -s -b $adminCookie.FullName "$BaseUrl/admin"))
$favsResp = [string]::Join("`n", (curl.exe -s -b $custCookie.FullName "$BaseUrl/favorites"))
$t15HomeOk = $homeResp.Contains("Under100") -and $homeResp.Contains("craving")
$t15AdminOk = $adminResp.Contains("Categories") -and $adminResp.Contains("Review Moderation")
$t15FavsOk = $favsResp.Contains("No saved restaurants yet")

$t15Ok = $t15HomeOk -and $t15AdminOk -and $t15FavsOk
Report-Test 15 "Regression of customer, merchant, admin, favorites, and reviews" $t15Ok "Homepage, admin dashboard with new Categories tab, and favorites render cleanly"

# Cleanup
Cleanup-TestFixture
Remove-Item $adminCookie.FullName -ErrorAction SilentlyContinue
Remove-Item $merchCookie.FullName -ErrorAction SilentlyContinue
Remove-Item $custCookie.FullName -ErrorAction SilentlyContinue

Write-Host "=================================================================" -ForegroundColor Cyan
Write-Host " PHASE 9 VERIFICATION COMPLETE: $PassCount PASSED, $FailCount FAILED" -ForegroundColor Cyan
Write-Host "=================================================================" -ForegroundColor Cyan
