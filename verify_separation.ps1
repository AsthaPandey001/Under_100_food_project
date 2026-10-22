# ==============================================================================
# UNDER100: CUSTOMER, MERCHANT & ADMIN SEPARATION VERIFICATION SUITE
# ==============================================================================
# Verifies three distinct, access-controlled experiences:
# 1. Customer-facing website (clean storefront, zero merchant/admin links)
# 2. Restaurant owner portal (/merchant/login, /merchant, role checks, cross-tenant isolation)
# 3. Platform administrator panel (/admin/login, /admin, 403 for unauthorized users)
# 4. Role-based redirects and open-redirect prevention
# 5. Zero database corruption / count preservation
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
Write-Host " UNDER100: THREE-WAY INTERFACE SEPARATION & SECURITY TEST SUITE   " -ForegroundColor Cyan
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
$guestCookie = New-TemporaryFile

# -----------------------------------------------------------------------------
# PART 1: CUSTOMER-FACING STOREFRONT AUDIT
# -----------------------------------------------------------------------------
Write-Host "`n-- PART 1: Customer Website Cleanliness & Storefront Functionality --" -ForegroundColor Yellow

# 1.1 Homepage loads with HTTP 200
$homeHtml = curl.exe -s -b $guestCookie.FullName -c $guestCookie.FullName "$BaseUrl/"
Report-Test "1.1" "Customer Homepage Loads" ($homeHtml -match "Under100" -or $homeHtml -match "under100") "Homepage returned valid HTML markup with Under100 branding"

# 1.2 Desktop Header contains no merchant or admin links
$hasPartnerNav = $homeHtml -match "btn-partner-nav" -or $homeHtml -match "Partner With Us"
$hasMerchantDirect = $homeHtml -match "btn-nav-merchant-direct" -or $homeHtml -match "My Outlets"
$hasAdminDirect = $homeHtml -match "btn-nav-admin-direct" -or $homeHtml -match "Admin Console"
Report-Test "1.2" "Customer Header Free of Management Links" (-not $hasPartnerNav -and -not $hasMerchantDirect -and -not $hasAdminDirect) "Desktop header contains zero merchant/admin entry points"

# 1.3 Footer contains no merchant or admin links
$hasFooterPartner = $homeHtml -match "FOR RESTAURANTS" -or $homeHtml -match "/partner/register"
$hasFooterMerchant = $homeHtml -match "Merchant Owner Portal"
$hasFooterAdmin = $homeHtml -match "Admin Database Console"
Report-Test "1.3" "Customer Footer Free of Management Controls" (-not $hasFooterPartner -and -not $hasFooterMerchant -and -not $hasFooterAdmin) "Public footer replaced with customer-facing CUSTOMER CARE links"

# 1.4 Homepage Login Modal has no role switcher or partner registration
$hasRoleSwitch = $homeHtml -match "Fast Role Access" -or $homeHtml -match "loginDirectRole"
$hasOwnerBanner = $homeHtml -match "Partner with us &amp; add your restaurant"
Report-Test "1.4" "Customer Login Modal Clean of Role Shortcuts" (-not $hasRoleSwitch -and -not $hasOwnerBanner) "Login modal contains only standard phone/email entry"

# 1.5 Restaurant Details Modal contains no Eatery Claiming Controls
$restHtml = curl.exe -s "$BaseUrl/restaurant/1"
$hasClaimBanner = $restHtml -match "Claim this listing" -or $restHtml -match "Claim This Eatery"
Report-Test "1.5" "Restaurant Detail Modal Free of Claim Banner" (-not $hasClaimBanner -and $restHtml -match "Menu Items") "Restaurant detail page shows food items and prices without management actions"

# 1.6 Customer Login Page (/login) has no demo box or account type dropdown
$loginHtml = curl.exe -s "$BaseUrl/login"
$hasDemoBox = $loginHtml -match "Quick 1-Click Role Login" -or $loginHtml -match "demo-box"
$hasRoleSelect = $loginHtml -match 'name="role"'
Report-Test "1.6" "Customer Login & Sign-up Page Free of Role Escalation" (-not $hasDemoBox -and -not $hasRoleSelect) "/login presents focused customer sign-in/sign-up without role dropdown"

# -----------------------------------------------------------------------------
# PART 2: RESTAURANT OWNER PORTAL & ACCESS CONTROLS
# -----------------------------------------------------------------------------
Write-Host "`n-- PART 2: Restaurant Owner Portal & Multi-Tenant Security --" -ForegroundColor Yellow

# 2.1 Dedicated Merchant Login Page (/merchant/login)
$merchLoginHtml = curl.exe -s "$BaseUrl/merchant/login"
$hasMerchBrand = $merchLoginHtml -match "Partner Portal" -or $merchLoginHtml -match "Restaurant Owner"
Report-Test "2.1" "Dedicated Merchant Login Page Accessible" ($hasMerchBrand -and $merchLoginHtml -match 'action="/merchant/login"') "/merchant/login presents dedicated partner portal sign-in"

# 2.2 Unauthenticated Access to /merchant Redirects to /merchant/login
$unauthMerchResp = curl.exe -s -i "$BaseUrl/merchant"
$isRedirectMerch = $unauthMerchResp -match "303 See Other" -and $unauthMerchResp -match "Location: /merchant/login\?msg=merchant_required"
Report-Test "2.2" "Unauthenticated Access to /merchant Protected" ($isRedirectMerch) "Unauthenticated user redirected to /merchant/login"

# 2.3 Customer Account Blocked from /merchant (403 Forbidden)
# Sign in as normal customer
$custLoginResp = curl.exe -s -i -c $custCookie.FullName -d "email=user@fooddiscovery.com&password=User@123456" "$BaseUrl/login"
$custMerchAccess = curl.exe -s -i -b $custCookie.FullName "$BaseUrl/merchant"
$is403CustMerch = $custMerchAccess -match "403 Forbidden"
Report-Test "2.3" "Customer Blocked from Merchant Portal (403)" ($is403CustMerch) "Customer account receives 403 Forbidden when requesting /merchant"

# 2.4 Customer Cannot Call Merchant Mutation Endpoints
$custMutation = curl.exe -s -i -b $custCookie.FullName -d "name=HackedDish&price=50" "$BaseUrl/merchant/menu"
$is403CustMutation = $custMutation -match "403 Forbidden"
Report-Test "2.4" "Customer Blocked from Merchant Mutations (403)" ($is403CustMutation) "Direct POST to /merchant/menu blocked with 403 Forbidden"

# 2.5 Authorized Merchant Login via /merchant/login
$merchLoginResp = curl.exe -s -i -c $merchCookie.FullName -d "email=merchant.ctr@fooddiscovery.com&password=Merchant@123456" "$BaseUrl/merchant/login"
$isMerchAuthRedirect = $merchLoginResp -match "Location: /merchant"
$merchDashHtml = curl.exe -s -b $merchCookie.FullName "$BaseUrl/merchant"
$hasMerchDash = $merchDashHtml -match "Merchant Owner Portal" -and ($merchDashHtml -match "Shri Sagar" -or $merchDashHtml -match "Central Tiffin Room")
Report-Test "2.5" "Authorized Merchant Signs In & Accesses Dashboard" ($isMerchAuthRedirect -and $hasMerchDash) "Merchant authenticated via /merchant/login and dashboard loaded"

# 2.6 Merchant Cannot Access Platform Admin Console (/admin)
$merchAdminAccess = curl.exe -s -i -b $merchCookie.FullName "$BaseUrl/admin"
$is403MerchAdmin = $merchAdminAccess -match "403 Forbidden"
Report-Test "2.6" "Merchant Blocked from Admin Panel (403)" ($is403MerchAdmin) "Merchant account receives 403 Forbidden when attempting to access /admin"

# 2.7 Cross-Tenant Isolation: Merchant Cannot Mutate Another Outlet's Dish
# CTR Shri Sagar owns restaurant 1. Try modifying dish 4 belonging to Brahmin's Coffee Bar (restaurant 2)
$crossTenantResp = curl.exe -s -i -b $merchCookie.FullName -d "name=IllegalEdit&price=50" "$BaseUrl/merchant/menu/4"
$is403CrossTenant = $crossTenantResp -match "403 Forbidden"
Report-Test "2.7" "Cross-Tenant Dish Mutation Blocked (403)" ($is403CrossTenant) "Merchant cannot edit dishes belonging to another restaurant"

# -----------------------------------------------------------------------------
# PART 3: PLATFORM ADMINISTRATOR PANEL & ACCESS CONTROLS
# -----------------------------------------------------------------------------
Write-Host "`n-- PART 3: Administrator Panel & Access Controls --" -ForegroundColor Yellow

# 3.1 Dedicated Admin Login Page (/admin/login)
$adminLoginHtml = curl.exe -s "$BaseUrl/admin/login"
$hasAdminBrand = $adminLoginHtml -match "Admin Console" -or $adminLoginHtml -match "Platform Administration"
Report-Test "3.1" "Dedicated Admin Login Page Accessible" ($hasAdminBrand -and $adminLoginHtml -match 'action="/admin/login"') "/admin/login presents dedicated administrative sign-in"

# 3.2 Unauthenticated Access to /admin Redirects to /admin/login
$unauthAdminResp = curl.exe -s -i "$BaseUrl/admin"
$isRedirectAdmin = $unauthAdminResp -match "303 See Other" -and $unauthAdminResp -match "Location: /admin/login\?msg=admin_required"
Report-Test "3.2" "Unauthenticated Access to /admin Protected" ($isRedirectAdmin) "Unauthenticated user redirected to /admin/login"

# 3.3 Customer Blocked from Admin Panel (403 Forbidden)
$custAdminAccess = curl.exe -s -i -b $custCookie.FullName "$BaseUrl/admin"
$is403CustAdmin = $custAdminAccess -match "403 Forbidden"
Report-Test "3.3" "Customer Blocked from Admin Panel (403)" ($is403CustAdmin) "Customer account receives 403 Forbidden when requesting /admin"

# 3.4 Customer Blocked from Admin Mutation Endpoints
$custAdminMut = curl.exe -s -i -b $custCookie.FullName -X POST "$BaseUrl/admin/restaurants/1/status?status=rejected"
$is403CustAdminMut = $custAdminMut -match "403 Forbidden"
Report-Test "3.4" "Customer Blocked from Admin Mutations (403)" ($is403CustAdminMut) "Direct POST to /admin/restaurants/1/status blocked with 403 Forbidden"

# 3.5 Authorized Admin Login via /admin/login
$adminLoginResp = curl.exe -s -i -c $adminCookie.FullName -d "email=admin@fooddiscovery.com&password=Admin@123456" "$BaseUrl/admin/login"
$isAdminAuthRedirect = $adminLoginResp -match "Location: /admin"
$adminDashHtml = curl.exe -s -b $adminCookie.FullName "$BaseUrl/admin"
$hasAdminDash = $adminDashHtml -match "Admin Dashboard" -and $adminDashHtml -match "Platform Users"
Report-Test "3.5" "Authorized Admin Signs In & Accesses Dashboard" ($isAdminAuthRedirect -and $hasAdminDash) "Admin authenticated via /admin/login and admin panel loaded"

# 3.6 Non-Admin Login Rejected at /admin/login
$badAdminLogin = curl.exe -s -i -d "email=user@fooddiscovery.com&password=User@123456" "$BaseUrl/admin/login"
$is403BadAdmin = $badAdminLogin -match "403 Forbidden"
Report-Test "3.6" "Customer Rejected at /admin/login Entry Point (403)" ($is403BadAdmin) "Non-admin credentials rejected with 403 Forbidden on /admin/login"

# -----------------------------------------------------------------------------
# PART 4: ROLE-BASED AUTHENTICATION, REDIRECTS & OPEN REDIRECT PREVENTION
# -----------------------------------------------------------------------------
Write-Host "`n-- PART 4: Role-Based Auto-Redirects & Open Redirect Prevention --" -ForegroundColor Yellow

# 4.1 Customer Login via /login Redirects to Storefront
$cLogin = curl.exe -s -i -d "email=user@fooddiscovery.com&password=User@123456" "$BaseUrl/login"
$isCustRedirect = $cLogin -match "Location: /\s*$"
Report-Test "4.1" "Customer Login Redirects to Storefront (/)" ($isCustRedirect) "Customer authenticated and redirected to /"

# 4.2 Merchant Login via /login Redirects to /merchant
$mLogin = curl.exe -s -i -d "email=merchant.ctr@fooddiscovery.com&password=Merchant@123456" "$BaseUrl/login"
$isMerchRedirect = $mLogin -match "Location: /merchant"
Report-Test "4.2" "Merchant Login via Shared /login Redirects to /merchant" ($isMerchRedirect) "Merchant role detected and auto-routed to /merchant"

# 4.3 Admin Login via /login Redirects to /admin
$aLogin = curl.exe -s -i -d "email=admin@fooddiscovery.com&password=Admin@123456" "$BaseUrl/login"
$isAdminRedirect = $aLogin -match "Location: /admin"
Report-Test "4.3" "Admin Login via Shared /login Redirects to /admin" ($isAdminRedirect) "Admin role detected and auto-routed to /admin"

# 4.4 Legitimate Internal Return URL Honored for Customer
$safeReturn = curl.exe -s -i -d "email=user@fooddiscovery.com&password=User@123456&return_to=/favorites" "$BaseUrl/login"
$isSafeReturn = $safeReturn -match "Location: /favorites"
Report-Test "4.4" "Safe Internal Return URL Honored" ($isSafeReturn) "User redirected to requested customer return URL (/favorites)"

# 4.5 Open Redirect Attempt (https://evil.com) Blocked
$evilReturn1 = curl.exe -s -i -d "email=user@fooddiscovery.com&password=User@123456&return_to=https://evil.com" "$BaseUrl/login"
$isBlocked1 = $evilReturn1 -match "Location: /\s*$"
Report-Test "4.5" "Open Redirect (https://evil.com) Blocked" ($isBlocked1) "External absolute URL rejected, safely redirected to /"

# 4.6 Open Redirect Attempt (//evil.com) Blocked
$evilReturn2 = curl.exe -s -i -d "email=user@fooddiscovery.com&password=User@123456&return_to=//evil.com" "$BaseUrl/login"
$isBlocked2 = $evilReturn2 -match "Location: /\s*$"
Report-Test "4.6" "Open Redirect (//evil.com) Blocked" ($isBlocked2) "Protocol-relative URL rejected, safely redirected to /"

# 4.7 Open Redirect Attempt (/\\evil.com) Blocked
$evilReturn3 = curl.exe -s -i -d "email=user@fooddiscovery.com&password=User@123456&return_to=/\evil.com" "$BaseUrl/login"
$isBlocked3 = $evilReturn3 -match "Location: /\s*$"
Report-Test "4.7" "Open Redirect (/\\evil.com) Blocked" ($isBlocked3) "Backslash bypass URL rejected, safely redirected to /"

# 4.8 Customer Attempting return_to=/admin Blocked
$privReturn = curl.exe -s -i -d "email=user@fooddiscovery.com&password=User@123456&return_to=/admin" "$BaseUrl/login"
$isPrivBlocked = $privReturn -match "Location: /\s*$"
Report-Test "4.8" "Customer Cannot Use return_to to Escalate to /admin" ($isPrivBlocked) "Return URL pointing to /admin ignored for customer, redirected to /"

# 4.9 Public Registration Privilege Escalation Prevented
# Attempt registering with form value role=merchant or role=admin
$secEmail = "sec_test_audit_$(Get-Random)@example.com"
$regResp = curl.exe -s -i -d "name=AuditUser&email=$secEmail&password=Passw0rd123&role=merchant" "$BaseUrl/register"
$isRegRedirect = $regResp -match "Location: /\s*$"
$userRoleInDb = (& $PsqlCmd -U $DbUser -d $DbName -t -A -c "SELECT role FROM users WHERE email = '$secEmail';").Trim()
$isAdminInDb = (& $PsqlCmd -U $DbUser -d $DbName -t -A -c "SELECT is_admin FROM users WHERE email = '$secEmail';").Trim()

Report-Test "4.9" "Registration Role Parameter Cannot Escalate Privileges" ($userRoleInDb -eq "customer" -and $isAdminInDb -eq "f") "Submitted role=merchant ignored: user created as role=customer, is_admin=false"

# Clean up test user record
& $PsqlCmd -U $DbUser -d $DbName -c "DELETE FROM users WHERE email = '$secEmail';" | Out-Null

# -----------------------------------------------------------------------------
# STEP 5: FINAL DATABASE INTEGRITY CHECK
# -----------------------------------------------------------------------------
Write-Host "`n-- Recording Final Database Verification Counts --" -ForegroundColor Yellow
$finalCounts = @{}
$dbFinalOutput = (& $PsqlCmd -U $DbUser -d $DbName -t -A -c @"
SELECT 'users', count(*) FROM users
UNION ALL SELECT 'restaurants', count(*) FROM restaurants
UNION ALL SELECT 'food_items', count(*) FROM food_items
UNION ALL SELECT 'categories', count(*) FROM categories
UNION ALL SELECT 'food_item_categories', count(*) FROM food_item_categories
UNION ALL SELECT 'restaurant_operating_hours', count(*) FROM restaurant_operating_hours
UNION ALL SELECT 'reviews', count(*) FROM reviews;
"@)

$dbPreserved = $true
foreach ($line in $dbFinalOutput) {
    if ($line -match '^(\w+)\|(\d+)$') {
        $tbl = $matches[1]
        $cnt = [int]$matches[2]
        $finalCounts[$tbl] = $cnt
        $initial = $initialCounts[$tbl]
        $match = ($cnt -eq $initial)
        if (-not $match) { $dbPreserved = $false }
        Write-Host "  $tbl`: Final=$cnt (Initial=$initial) $(if($match){'[PRESERVED]'}else{'[MISMATCH]'})" -ForegroundColor $(if($match){'Green'}else{'Red'})
    }
}

Report-Test "5.1" "100% Database Records Preserved" ($dbPreserved) "All database tables match exact pre-test baseline counts"

# Clean up session temp files
Remove-Item $custCookie.FullName -Force -ErrorAction SilentlyContinue
Remove-Item $merchCookie.FullName -Force -ErrorAction SilentlyContinue
Remove-Item $adminCookie.FullName -Force -ErrorAction SilentlyContinue
Remove-Item $guestCookie.FullName -Force -ErrorAction SilentlyContinue

Write-Host "`n=================================================================" -ForegroundColor Cyan
Write-Host " TEST SUITE COMPLETE: $PassCount passed, $FailCount failed out of $TotalTests tests" -ForegroundColor $(if($FailCount -eq 0){'Green'}else{'Red'})
Write-Host "=================================================================`n" -ForegroundColor Cyan

if ($FailCount -gt 0) {
    exit 1
} else {
    exit 0
}
