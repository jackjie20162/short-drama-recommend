$ErrorActionPreference = "Stop"
$Root = (Resolve-Path (Join-Path $PSScriptRoot "..")).Path
Set-Location $Root

function Require-Command([string]$Name) {
    $command = Get-Command $Name -ErrorAction SilentlyContinue
    if (-not $command) { throw "$Name is required. Please install it and make sure it is available in PATH." }
    return $command.Source
}

$goctl = Require-Command "goctl"
$protoc = Require-Command "protoc"
$protocGenGo = Require-Command "protoc-gen-go"
$protocGenGoGrpc = Require-Command "protoc-gen-go-grpc"

Write-Host "goctl: $goctl"
Write-Host "protoc: $protoc"
Write-Host "protoc-gen-go: $protocGenGo"
Write-Host "protoc-gen-go-grpc: $protocGenGoGrpc"

$ProtoDir = Join-Path $Root "proto"

function Flatten-Pb([string]$Out) {
    $outDir = Join-Path $Root $Out
    $pbDir = Join-Path $outDir "pb"
    $relativeOut = $Out -replace "^rpc[\\/]", ""
    $nestedPackageDir = Join-Path (Join-Path $pbDir "short-drama-recommend") $relativeOut
    if (-not (Test-Path $nestedPackageDir)) { return }
    Write-Host "Flattening generated PB: $nestedPackageDir -> $pbDir"
    Get-ChildItem -Path $nestedPackageDir -File | ForEach-Object { Move-Item -Force $_.FullName (Join-Path $pbDir $_.Name) }
    Remove-Item -Recurse -Force (Join-Path $pbDir "short-drama-recommend")
}

function Generate-Rpc([string]$Proto, [string]$Out) {
    Write-Host ""
    Write-Host "Generating RPC: $Proto -> $Out"
    $outDir = Join-Path $Root $Out
    New-Item -ItemType Directory -Force -Path (Join-Path $outDir "pb") | Out-Null
    Push-Location $ProtoDir
    try {
        & $goctl rpc protoc $Proto `
            "--go_out=../$Out/pb" `
            "--go_opt=module=short-drama-recommend" `
            "--go-grpc_out=../$Out/pb" `
            "--go-grpc_opt=module=short-drama-recommend" `
            "--zrpc_out=../$Out"
        if ($LASTEXITCODE -ne 0) { throw "RPC generation failed: $Proto" }
    } finally { Pop-Location }
    Flatten-Pb $Out
}

Generate-Rpc "user.proto" "rpc/user-rpc"
Generate-Rpc "drama.proto" "rpc/drama-rpc"
Generate-Rpc "drama_admin.proto" "rpc/drama-rpc"
Generate-Rpc "behavior.proto" "rpc/behavior-rpc"
Generate-Rpc "recommend.proto" "rpc/recommend-rpc"
Generate-Rpc "payment.proto" "rpc/payment-rpc"
Generate-Rpc "media.proto" "rpc/media-rpc"

Write-Host ""
Write-Host "go-zero RPC generation completed"