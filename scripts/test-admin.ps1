$ErrorActionPreference = "Stop"
$BaseUrl = if ($env:ADMIN_BASE_URL) { $env:ADMIN_BASE_URL } else { "http://127.0.0.1:8081" }
Write-Host "== short-drama admin smoke test =="`n
$drama = @{ title="Local Test Drama"; description="Created by admin smoke test"; cover=""; country="US"; language="en"; total_episodes=3; is_paid=$false } | ConvertTo-Json
$created = Invoke-RestMethod -Uri "$BaseUrl/api/v1/admin/dramas" -Method Post -ContentType "application/json" -Body $drama
$created | ConvertTo-Json -Depth 8
$id = $created.drama.id
Write-Host "Created drama id=$id"
$status = @{ status=1 } | ConvertTo-Json
Invoke-RestMethod -Uri "$BaseUrl/api/v1/admin/dramas/$id/status" -Method Post -ContentType "application/json" -Body $status | ConvertTo-Json -Depth 8
$episode = @{ episode_no=1; title="Episode 1"; duration_seconds=90; video_url="https://example.com/test/1.m3u8"; poster_url=""; is_paid=$false } | ConvertTo-Json
Invoke-RestMethod -Uri "$BaseUrl/api/v1/admin/dramas/$id/episodes" -Method Post -ContentType "application/json" -Body $episode | ConvertTo-Json -Depth 8
Invoke-RestMethod -Uri "$BaseUrl/api/v1/admin/dramas/$id/episodes" -Method Get | ConvertTo-Json -Depth 8
Write-Host "Admin smoke test completed. Drama id=$id"
