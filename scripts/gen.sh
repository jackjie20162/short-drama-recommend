#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

command -v goctl >/dev/null 2>&1 || { echo "goctl is required"; exit 1; }
command -v protoc >/dev/null 2>&1 || { echo "protoc is required"; exit 1; }
command -v protoc-gen-go >/dev/null 2>&1 || { echo "protoc-gen-go is required"; exit 1; }
command -v protoc-gen-go-grpc >/dev/null 2>&1 || { echo "protoc-gen-go-grpc is required"; exit 1; }

generate_rpc() {
  local proto="$1"
  local out="$2"
  mkdir -p "$out/pb"
  goctl rpc protoc "$proto" \\
    --go_out="$out/pb" \\
    --go-grpc_out="$out/pb" \\
    --zrpc_out="$out"
}

generate_rpc proto/user.proto rpc/user-rpc
generate_rpc proto/drama.proto rpc/drama-rpc
generate_rpc proto/drama_admin.proto rpc/drama-rpc
generate_rpc proto/behavior.proto rpc/behavior-rpc
generate_rpc proto/recommend.proto rpc/recommend-rpc
generate_rpc proto/payment.proto rpc/payment-rpc
generate_rpc proto/media.proto rpc/media-rpc

echo "go-zero RPC generation completed"
