$ErrorActionPreference = "Stop"

$Root = (Resolve-Path (Join-Path $PSScriptRoot "..")).Path
Set-Location $Root

function Require-Command([string]$Name) {
    $command = Get-Command $Name -ErrorAction SilentlyContinue
    if (-not $command) {
        throw "$Name is required. Please install it and make sure it is available in PATH."
    }
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

function Generate-Rpc([string]$Proto, [string]$Out) {
    Write-Host ""
    Write-Host "Generating RPC: $Proto -> $Out"
    New-Item -ItemType Directory -Force -Path (Join-Path $Out "pb") | Out-Null
    & $goctl rpc protoc $Proto `
        "--go_out=$Out/pb" `
        "--go-grpc_out=$Out/pb" `
        "--zrpc_out=$Out"
    if ($LASTEXITCODE -ne 0) {
        throw "RPC generation failed: $Proto"
    }
}

Generate-Rpc "proto/user.proto" "rpc/user-rpc"
Generate-Rpc "proto/drama.proto" "rpc/drama-rpc"
Generate-Rpc "proto/drama_admin.proto" "rpc/drama-rpc"
Generate-Rpc "proto/behavior.proto" "rpc/behavior-rpc"
Generate-Rpc "proto/recommend.proto" "rpc/recommend-rpc"
Generate-Rpc "proto/payment.proto" "rpc/payment-rpc"
Generate-Rpc "proto/media.proto" "rpc/media-rpc"

Write-Host ""
Write-Host "go-zero RPC generation completed"
