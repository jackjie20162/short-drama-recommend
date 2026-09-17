#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

command -v goctl >/dev/null 2>&1 || { echo "goctl is required"; exit 1; }
command -v protoc >/dev/null 2>&1 || { echo "protoc is required"; exit 1; }

mkdir -p rpc/user-rpc/pb rpc/drama-rpc/pb rpc/behavior-rpc/pb rpc/recommend-rpc/pb

for proto in proto/*.proto; do
  protoc \
    --go_out=. \
    --go_opt=paths=source_relative \
    --go-grpc_out=. \
    --go-grpc_opt=paths=source_relative \
    "$proto"
done

echo "protobuf generation completed"
