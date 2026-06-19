Write-Host "=== Testing API ===" -ForegroundColor Cyan

# 1. Health Check
Write-Host "`n1. Health Check:" -ForegroundColor Yellow
$health = Invoke-RestMethod -Uri "http://localhost:8080/health" -Method GET
Write-Host "  ✓ Health: $($health | ConvertTo-Json)" -ForegroundColor Green

# 2. Register a new user
Write-Host "`n2. Register User:" -ForegroundColor Yellow
$registerBody = @{
    name = "Test User"
    email = "test$(Get-Random)@example.com"
    password = "password123"
} | ConvertTo-Json

try {
    $register = Invoke-RestMethod -Uri "http://localhost:8080/api/register" -Method POST -Body $registerBody -ContentType "application/json"
    Write-Host "  ✓ Registration successful!" -ForegroundColor Green
    Write-Host "  User: $($register.user | ConvertTo-Json)" -ForegroundColor Gray
    $userEmail = $register.user.email
    $userPassword = "password123"
} catch {
    Write-Host "  ✗ Registration failed: $($_.Exception.Message)" -ForegroundColor Red
    # Use existing test user
    $userEmail = "test@example.com"
    $userPassword = "password123"
}

# 3. Login
Write-Host "`n3. Login:" -ForegroundColor Yellow
$loginBody = @{
    email = $userEmail
    password = $userPassword
} | ConvertTo-Json

try {
    $login = Invoke-RestMethod -Uri "http://localhost:8080/api/login" -Method POST -Body $loginBody -ContentType "application/json"
    Write-Host "  ✓ Login successful!" -ForegroundColor Green
    $token = $login.token
    Write-Host "  Token: $token" -ForegroundColor Gray
} catch {
    Write-Host "  ✗ Login failed: $($_.Exception.Message)" -ForegroundColor Red
    exit
}

# 4. Get Users (Protected endpoint)
Write-Host "`n4. Get Users (Protected):" -ForegroundColor Yellow
$headers = @{
    "Authorization" = "Bearer $token"
}

try {
    $users = Invoke-RestMethod -Uri "http://localhost:8080/api/users" -Method GET -Headers $headers
    Write-Host "  ✓ Users retrieved successfully!" -ForegroundColor Green
    Write-Host "  Users: $($users | ConvertTo-Json)" -ForegroundColor Gray
} catch {
    Write-Host "  ✗ Failed to get users: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n=== Testing Complete ===" -ForegroundColor Cyan
