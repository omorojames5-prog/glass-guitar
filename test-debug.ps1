Write-Host "=== Testing Login with Debug Logging ===" -ForegroundColor Cyan

$timestamp = [DateTimeOffset]::UtcNow.ToUnixTimeMilliseconds()
$email = "debug$timestamp@example.com"
$password = "password123"

Write-Host "Email: $email" -ForegroundColor Yellow
Write-Host "Password: $password" -ForegroundColor Yellow
Write-Host ""

# Register
$body = @{
    name = "Debug Test"
    email = $email
    password = $password
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/api/register" -Method POST -Body $body -ContentType "application/json"
    Write-Host "✅ Registration successful!" -ForegroundColor Green
} catch {
    Write-Host "❌ Registration failed: $($_.Exception.Message)" -ForegroundColor Red
    exit
}

# Wait a moment
Start-Sleep -Seconds 1

# Login
$loginBody = @{
    email = $email
    password = $password
} | ConvertTo-Json

try {
    $login = Invoke-RestMethod -Uri "http://localhost:8080/api/login" -Method POST -Body $loginBody -ContentType "application/json"
    Write-Host "✅ LOGIN SUCCESSFUL! 🎉" -ForegroundColor Green
    Write-Host "Token: $($login.token)" -ForegroundColor Gray
} catch {
    Write-Host "❌ Login failed: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.Exception.Response) {
        $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
        $reader.BaseStream.Position = 0
        $reader.DiscardBufferedData()
        $errorBody = $reader.ReadToEnd()
        Write-Host "Error: $errorBody" -ForegroundColor Red
    }
}
