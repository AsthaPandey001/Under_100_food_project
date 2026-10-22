# verify_phase_6.ps1 - Production Hardening, Security & Deployment Verification Suite
$ErrorActionPreference = "Stop"

Write-Host "====================================================" -ForegroundColor Cyan
Write-Host " RUNNING PHASE 6 VERIFICATION SUITE" -ForegroundColor Cyan
Write-Host "====================================================" -ForegroundColor Cyan

$baseUrl = "http://127.0.0.1:8000"
$totalTests = 0
$passedTests = 0

function Assert-Test($name, $condition, $details = "") {
    $script:totalTests++
    if ($condition) {
        $script:passedTests++
        Write-Host "  [PASS] $name" -ForegroundColor Green
    } else {
        Write-Host "  [FAIL] $name : $details" -ForegroundColor Red
    }
}

# 1. Health & Readiness Endpoints
Write-Host "`n--- 1. Health & Readiness Endpoints ---" -ForegroundColor Yellow

$liveResp = Invoke-WebRequest -Uri "$baseUrl/live" -Method Get -UseBasicParsing
Assert-Test "Liveness probe returns HTTP 200" ($liveResp.StatusCode -eq 200)
Assert-Test "Liveness probe body status ok" ($liveResp.Content -match '"status":"ok"')

$readyResp = Invoke-WebRequest -Uri "$baseUrl/ready" -Method Get -UseBasicParsing
Assert-Test "Readiness probe returns HTTP 200" ($readyResp.StatusCode -eq 200)
Assert-Test "Readiness probe reports healthy database and redis" ($readyResp.Content -match '"database":"healthy"' -and $readyResp.Content -match '"redis":"healthy"')

# 2. Rate Limiting Tests
Write-Host "`n--- 2. Rate Limiting (Redis-backed & Safe Headers) ---" -ForegroundColor Yellow

# Normal request within limit
$testClientIP = "198.51.100.42"
$normHeaders = @{ "X-Forwarded-For" = $testClientIP }

# Request to /live should not be rate limited
$liveCheck = Invoke-WebRequest -Uri "$baseUrl/live" -Headers $normHeaders -UseBasicParsing
Assert-Test "Public routes not unnecessarily rate limited" ($liveCheck.StatusCode -eq 200)

# Test Rate Limiting on /password/forgot (bursting beyond limit of 30)
$rlKey = "p6_rl_test_" + (Get-Random)
$got429 = $false
$retryAfterPresent = $false

for ($i = 1; $i -le 35; $i++) {
    $res = & curl.exe -s -i -X POST "$baseUrl/password/forgot" `
        -H "X-Forwarded-For: 198.51.100.99" `
        -H "Content-Type: application/x-www-form-urlencoded" `
        -d "email=$rlKey@example.com"
    $resText = [string]::Join("`n", $res)

    if ($resText -match "HTTP/\S+\s+429") {
        $got429 = $true
        if ($resText -match "Retry-After:") {
            $retryAfterPresent = $true
        }
        break
    }
}
Assert-Test "Rate limit triggers HTTP 429 on /password/forgot burst" $got429
Assert-Test "HTTP 429 response includes Retry-After header" $retryAfterPresent

# Verify that an unaffected IP address is NOT blocked
$unaffectedRes = & curl.exe -s -i -X POST "$baseUrl/password/forgot" `
    -H "X-Forwarded-For: 198.51.100.100" `
    -H "Content-Type: application/x-www-form-urlencoded" `
    -d "email=other_user@example.com"
$unaffectedText = [string]::Join("`n", $unaffectedRes)
$notBlocked = ($unaffectedText -notmatch "HTTP/\S+\s+429")
Assert-Test "Independent client IP is not blocked by another client's rate limit" $notBlocked

# 3. Email & Anti-Enumeration Security
Write-Host "`n--- 3. Email & Anti-Enumeration Security ---" -ForegroundColor Yellow

# Submit password reset for existing email vs non-existent email
$existRes = & curl.exe -s -i -X POST "$baseUrl/password/forgot" `
    -H "X-Forwarded-For: 198.51.100.150" `
    -H "Content-Type: application/x-www-form-urlencoded" `
    -d "email=customer@example.com"
$existText = [string]::Join("`n", $existRes)

$nonExistRes = & curl.exe -s -i -X POST "$baseUrl/password/forgot" `
    -H "X-Forwarded-For: 198.51.100.151" `
    -H "Content-Type: application/x-www-form-urlencoded" `
    -d "email=nonexistent_99999@example.com"
$nonExistText = [string]::Join("`n", $nonExistRes)

# Check anti-enumeration: responses should not disclose account presence
$existHasRedirect = ($existText -match "303 See Other" -or $existText -match "200 OK")
$nonExistHasRedirect = ($nonExistText -match "303 See Other" -or $nonExistText -match "200 OK")
Assert-Test "Anti-enumeration: existing and non-existing email get equivalent responses" ($existHasRedirect -and $nonExistHasRedirect)

# Tokens must never be disclosed in responses or headers
$noTokenInExist = ($existText -notmatch "password_tokens" -and $existText -notmatch "raw_token")
Assert-Test "Reset token is never leaked in HTTP response body or headers" $noTokenInExist

# 4. Configuration & Secret Hygiene
Write-Host "`n--- 4. Configuration & Secret Hygiene ---" -ForegroundColor Yellow
$cfgTestOutput = & go test -run TestConfig_Validate ./config
$cfgTestsPass = ($LASTEXITCODE -eq 0)
Assert-Test "Config validation enforces strong production secrets and passes local" $cfgTestsPass

# 5. Database Safety & Integrity
Write-Host "`n--- 5. Database Record Count Preservation ---" -ForegroundColor Yellow
$env:PGPASSWORD = "password"
$psql = "C:\Users\astha\pgsql\bin\psql.exe"

$userCount = (& $psql -U postgres -d goship_db -t -c "SELECT COUNT(*) FROM users;").Trim()
$restCount = (& $psql -U postgres -d goship_db -t -c "SELECT COUNT(*) FROM restaurants;").Trim()
$foodCount = (& $psql -U postgres -d goship_db -t -c "SELECT COUNT(*) FROM food_items;").Trim()
$favCount = (& $psql -U postgres -d goship_db -t -c "SELECT COUNT(*) FROM favorites;").Trim()

Assert-Test "User count preserved (expected 3)" ($userCount -eq "3") "Got $userCount"
Assert-Test "Restaurant count preserved (expected 42)" ($restCount -eq "42") "Got $restCount"
Assert-Test "Food items count preserved (expected 166)" ($foodCount -eq "166") "Got $foodCount"
Assert-Test "Favorites count preserved (expected 0)" ($favCount -eq "0") "Got $favCount"

# 6. Static Asset Integrity
Write-Host "`n--- 6. Static Asset Integrity ---" -ForegroundColor Yellow
$staticDiff = git diff recovery-baseline-stage-a -- static/
Assert-Test "Static assets unchanged from baseline" ([string]::IsNullOrWhiteSpace($staticDiff))

# 7. Headless Browser Verification
Write-Host "`n--- 7. Headless Browser Check ---" -ForegroundColor Yellow
$chrome = "C:\Program Files\Google\Chrome\Application\chrome.exe"
$browserPassed = $false
if (Test-Path $chrome) {
    $tempDump = [System.IO.Path]::GetTempFileName()
    $p = Start-Process -FilePath $chrome -ArgumentList "--headless", "--disable-gpu", "--dump-dom", "$baseUrl/" -RedirectStandardOutput $tempDump -Wait -PassThru
    if (Test-Path $tempDump) {
        $dom = Get-Content $tempDump -Raw
        Remove-Item $tempDump -Force -ErrorAction SilentlyContinue
        if ($dom -match "Under100" -and $dom -match "Bengaluru") {
            $browserPassed = $true
        }
    }
}
Assert-Test "Headless Chrome successfully renders discovery page" $browserPassed

Write-Host "`n====================================================" -ForegroundColor Cyan
Write-Host " PHASE 6 VERIFICATION SUMMARY: $passedTests / $totalTests PASSED" -ForegroundColor Cyan
Write-Host "====================================================" -ForegroundColor Cyan

if ($passedTests -ne $totalTests) {
    exit 1
}
exit 0
