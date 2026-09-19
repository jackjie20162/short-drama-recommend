#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

command -v protoc >/dev/null 2>&1 || { echo "protoc is required"; exit 1; }
command -v protoc-gen-go >/dev/null 2>&1 || { echo "protoc-gen-go is required"; exit 1; }
command -v protoc-gen-go-grpc >/dev/null 2>&1 || { echo "protoc-gen-go-grpc is required"; exit 1; }

generate_proto() {
  local proto="$1"
  local out="$2"
  mkdir -p "$out"
  protoc -I proto --go_out="$out" --go_opt=paths=source_relative \
    --go-grpc_out="$out" --go-grpc_opt=paths=source_relative "$(basename "$proto")"
}

generate_proto proto/user.proto rpc/user-rpc/pb
generate_proto proto/drama.proto rpc/drama-rpc/pb
generate_proto proto/drama_admin.proto rpc/drama-rpc/pb
generate_proto proto/behavior.proto rpc/behavior-rpc/pb
generate_proto proto/recommend.proto rpc/recommend-rpc/pb
generate_proto proto/payment.proto rpc/payment-rpc/pb
generate_proto proto/media.proto rpc/media-rpc/pb

echo "protobuf generation completed; API gateway code is preserved"
