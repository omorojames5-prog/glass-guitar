Write-Host "=== Testing API with Unique Emails ===" -ForegroundColor Cyan

# Generate unique email
$randomNum = Get-Random -Minimum 1000 -Maximum 999999
$testEmail = "test$randomNum@example.com"
$testPassword = "password123"

Write-Host "Using email: $testEmail" -ForegroundColor Yellow

# 1. Health Check
Write-Host "`n1. Health Check:" -ForegroundColor Yellow
try {
    $health = Invoke-RestMethod -Uri "http://localhost:8080/health" -Method GET
    Write-Host "  ✓ Health: $($health | ConvertTo-Json)" -ForegroundColor Green
} catch {
    Write-Host "  ✗ Health check failed: $($_.Exception.Message)" -ForegroundColor Red
    exit
}

# 2. Register a new user with unique email
Write-Host "`n2. Register User:" -ForegroundColor Yellow
$registerBody = @{
    name = "Test User $randomNum"
    email = $testEmail
    password = $testPassword
} | ConvertTo-Json

try {
    $register = Invoke-RestMethod -Uri "http://localhost:8080/api/register" -Method POST -Body $registerBody -ContentType "application/json"
    Write-Host "  ✓ Registration successful!" -ForegroundColor Green
    Write-Host "  User: $($register.user | ConvertTo-Json)" -ForegroundColor Gray
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

# 3. Login with the new user
Write-Host "`n3. Login:" -ForegroundColor Yellow
$loginBody = @{
    email = $testEmail
    password = $testPassword
} | ConvertTo-Json

try {
    $login = Invoke-RestMethod -Uri "http://localhost:8080/api/login" -Method POST -Body $loginBody -ContentType "application/json"
    Write-Host "  ✓ Login successful!" -ForegroundColor Green
    Write-Host "  Token: $($login.token)" -ForegroundColor Gray
    Write-Host "  User: $($login.user | ConvertTo-Json)" -ForegroundColor Gray
    $token = $login.token
} catch {
    Write-Host "  ✗ Login failed: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.Exception.Response) {
        $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
        $reader.BaseStream.Position = 0
        $reader.DiscardBufferedData()
        $errorBody = $reader.ReadToEnd()
        Write-Host "  Error Response: $errorBody" -ForegroundColor Red
    }
    exit
}

# 4. Get Users (Protected)
Write-Host "`n4. Get Users (Protected):" -ForegroundColor Yellow
$headers = @{
    "Authorization" = "Bearer $token"
}

try {
    $users = Invoke-RestMethod -Uri "http://localhost:8080/api/users" -Method GET -Headers $headers
    Write-Host "  ✓ Users retrieved successfully!" -ForegroundColor Green
    Write-Host "  Total users: $($users.Count)" -ForegroundColor Gray
} catch {
    Write-Host "  ✗ Failed to get users: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.Exception.Response) {
        $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
        $reader.BaseStream.Position = 0
        $reader.DiscardBufferedData()
        $errorBody = $reader.ReadToEnd()
        Write-Host "  Error Response: $errorBody" -ForegroundColor Red
    }
}

Write-Host "`n=== Testing Complete ===" -ForegroundColor Cyan
