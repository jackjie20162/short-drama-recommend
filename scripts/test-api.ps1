$ErrorActionPreference = "Stop"
$BaseUrl = if ($env:BASE_URL) { $env:BASE_URL } else { "http://127.0.0.1:8080" }
Write-Host "== short-drama API smoke test =="
Write-Host "BASE_URL=$BaseUrl"
function Get-Api($Path) { Invoke-RestMethod -Uri "$BaseUrl$Path" -Method Get }
Write-Host "[1] drama list"; Get-Api "/api/v1/dramas?page=1&page_size=20" | ConvertTo-Json -Depth 8
Write-Host "[2] drama detail"; Get-Api "/api/v1/dramas/1" | ConvertTo-Json -Depth 8
Write-Host "[3] recommendation feed"; Get-Api "/api/v1/feed?user_id=1&country=US&language=en&page_size=10" | ConvertTo-Json -Depth 8
Write-Host "[4] behavior event"
$body = @{ user_id=1; drama_id=1; episode_id=1; event_type="PLAY"; watch_seconds=5; duration_seconds=60; country="US"; language="en"; device="web" } | ConvertTo-Json
Invoke-RestMethod -Uri "$BaseUrl/api/v1/behaviors" -Method Post -ContentType "application/json" -Body $body | ConvertTo-Json -Depth 8
Write-Host "Smoke test completed."
