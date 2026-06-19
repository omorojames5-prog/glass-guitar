Write-Host "=== Testing Protected Routes ===" -ForegroundColor Cyan

# 1. Login
Write-Host "`n1. Logging in..." -ForegroundColor Yellow
$loginBody = @{
    email = "fix1781864991128@example.com"
    password = "password123"
} | ConvertTo-Json

try {
    $login = Invoke-RestMethod -Uri "http://localhost:8080/api/login" -Method POST -Body $loginBody -ContentType "application/json"
    $token = $login.token
    Write-Host "  ✅ Login successful!" -ForegroundColor Green
    Write-Host "  Token: $($token.Substring(0, 30))..." -ForegroundColor Gray
} catch {
    Write-Host "  ❌ Login failed: $($_.Exception.Message)" -ForegroundColor Red
    exit
}

# 2. Get all users (protected)
Write-Host "`n2. Getting all users..." -ForegroundColor Yellow
$headers = @{
    "Authorization" = "Bearer $token"
}

try {
    $users = Invoke-RestMethod -Uri "http://localhost:8080/api/users" -Method GET -Headers $headers
    Write-Host "  ✅ Users retrieved successfully!" -ForegroundColor Green
    Write-Host "  Total users: $($users.Count)" -ForegroundColor Gray
    Write-Host "  Users: $($users | ConvertTo-Json)" -ForegroundColor Gray
} catch {
    Write-Host "  ❌ Failed to get users: $($_.Exception.Message)" -ForegroundColor Red
}

# 3. Get a specific user (protected)
Write-Host "`n3. Getting specific user..." -ForegroundColor Yellow
try {
    $userId = 1
    $user = Invoke-RestMethod -Uri "http://localhost:8080/api/users/$userId" -Method GET -Headers $headers
    Write-Host "  ✅ User retrieved successfully!" -ForegroundColor Green
    Write-Host "  User: $($user | ConvertTo-Json)" -ForegroundColor Gray
} catch {
    Write-Host "  ❌ Failed to get user: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n=== Test Complete ===" -ForegroundColor Cyan
