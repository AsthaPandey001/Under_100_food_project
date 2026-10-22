Write-Host "=== TEST 1: Guest homepage ==="
$r1 = Invoke-WebRequest -Uri "http://localhost:8000/" -UseBasicParsing
$hasSignIn = $r1.Content.Contains("Sign In")
$hasSignUp = $r1.Content.Contains("Sign Up")
$hasAdminBtn = $r1.Content.Contains('href="/admin"')
Write-Host "Guest Home - Has Sign In: $hasSignIn, Has Sign Up: $hasSignUp, Has Admin Link: $hasAdminBtn"

Write-Host "`n=== TEST 2: Guest accessing /admin ==="
try {
    $r2 = Invoke-WebRequest -Uri "http://localhost:8000/admin" -MaximumRedirection 0 -ErrorAction Stop -UseBasicParsing
    Write-Host "Unexpected status: $($r2.StatusCode)"
} catch {
    $resp = $_.Exception.Response
    Write-Host "Redirect Status: $($resp.StatusCode), Location: $($resp.Headers['Location'])"
}

Write-Host "`n=== TEST 3: Dedicated /login page ==="
$r3 = Invoke-WebRequest -Uri "http://localhost:8000/login?msg=admin_required" -UseBasicParsing
$hasAdminReq = $r3.Content.Contains("Administrator Access Required")
$hasDemoAdmin = $r3.Content.Contains("Demo Admin")
$hasDemoUser = $r3.Content.Contains("Demo User")
Write-Host "Login Page - Has Admin Warning: $hasAdminReq, Has Demo Admin: $hasDemoAdmin, Has Demo User: $hasDemoUser"

Write-Host "`n=== TEST 4: Normal User Authentication & Forbidden Check ==="
$userSession = New-Object Microsoft.PowerShell.Commands.WebRequestSession
$r4 = Invoke-WebRequest -Uri "http://localhost:8000/api/quick-login" -Method POST -Body @{role="user"} -WebSession $userSession -MaximumRedirection 0 -ErrorAction SilentlyContinue
Write-Host "User Quick Login Status: $($r4.StatusCode)"

# Check user homepage
$r4Home = Invoke-WebRequest -Uri "http://localhost:8000/" -WebSession $userSession -UseBasicParsing
$userHasAdminLink = $r4Home.Content.Contains('href="/admin"')
$userHasLogout = $r4Home.Content.Contains('href="/logout"')
Write-Host "User Home - Has Admin Link: $userHasAdminLink, Has Logout: $userHasLogout"

# User accessing /admin (MUST RETURN 403 Forbidden)
try {
    $r4Admin = Invoke-WebRequest -Uri "http://localhost:8000/admin" -WebSession $userSession -UseBasicParsing -ErrorAction Stop
    Write-Host "UNEXPECTED: Non-admin accessed /admin with status: $($r4Admin.StatusCode)"
} catch {
    $resp = $_.Exception.Response
    Write-Host "User accessing /admin - Forbidden Status Code: $($resp.StatusCode)"
    $sr = New-Object System.IO.StreamReader($resp.GetResponseStream())
    $body = $sr.ReadToEnd()
    $has403Text = $body.Contains("Administrator Access Required") -or $body.Contains("403 Access Denied")
    Write-Host "User accessing /admin - Received 403 page: $has403Text"
}

Write-Host "`n=== TEST 5: Admin Authentication & Dashboard Access ==="
$adminSession = New-Object Microsoft.PowerShell.Commands.WebRequestSession
$r5 = Invoke-WebRequest -Uri "http://localhost:8000/api/quick-login" -Method POST -Body @{role="admin"} -WebSession $adminSession -MaximumRedirection 0 -ErrorAction SilentlyContinue
Write-Host "Admin Quick Login Status: $($r5.StatusCode)"

# Check admin homepage
$r5Home = Invoke-WebRequest -Uri "http://localhost:8000/" -WebSession $adminSession -UseBasicParsing
$adminHasAdminLink = $r5Home.Content.Contains('href="/admin"')
$adminHasBadge = $r5Home.Content.Contains('badge-role')
Write-Host "Admin Home - Has Admin Link: $adminHasAdminLink, Has Admin Badge: $adminHasBadge"

# Admin accessing /admin (MUST RETURN 200 OK)
$r5Admin = Invoke-WebRequest -Uri "http://localhost:8000/admin" -WebSession $adminSession -UseBasicParsing
$adminDashboardLoaded = $r5Admin.Content.Contains("Food Discovery Admin Panel")
Write-Host "Admin accessing /admin - Status: $($r5Admin.StatusCode), Dashboard Loaded: $adminDashboardLoaded"
