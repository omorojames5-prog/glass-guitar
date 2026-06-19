$randomNum = Get-Random -Minimum 100000 -Maximum 999999
$testEmail = "test$randomNum@example.com"
$testPassword = "password123"

Write-Host "=== Testing with NEW User ===" -ForegroundColor Cyan
Write-Host "Email: $testEmail" -ForegroundColor Yellow
Write-Host "Password: $testPassword" -ForegroundColor Yellow

# Register
$registerBody = @{
    name = "Test User $randomNum"
    email = $testEmail
    password = $testPassword
} | ConvertTo-Json

Write-Host "`n1. Registering..." -ForegroundColor Yellow
try {
    $register = Invoke-RestMethod -Uri "http://localhost:8080/api/register" -Method POST -Body $registerBody -ContentType "application/json"
    Write-Host "  ✓ Registration successful!" -ForegroundColor Green
} catch {
    Write-Host "  ✗ Registration failed: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.Exception.Response) {
        $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
        $reader.BaseStream.Position = 0
        $reader.DiscardBufferedData()
        $errorBody = $reader.ReadToEnd()
        Write-Host "  Error: $errorBody" -ForegroundColor Red
    }
    exit
}

# Login
Write-Host "`n2. Logging in..." -ForegroundColor Yellow
$loginBody = @{
    email = $testEmail
    password = $testPassword
} | ConvertTo-Json

try {
    $login = Invoke-RestMethod -Uri "http://localhost:8080/api/login" -Method POST -Body $loginBody -ContentType "application/json"
    Write-Host "  ✓ Login successful!" -ForegroundColor Green
    Write-Host "  Token: $($login.token)" -ForegroundColor Gray
    Write-Host "  User: $($login.user.email)" -ForegroundColor Gray
} catch {
    Write-Host "  ✗ Login failed: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.Exception.Response) {
        $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
        $reader.BaseStream.Position = 0
        $reader.DiscardBufferedData()
        $errorBody = $reader.ReadToEnd()
        Write-Host "  Error: $errorBody" -ForegroundColor Red
    }
}
