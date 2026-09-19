#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

command -v goctl >/dev/null 2>&1 || { echo "goctl is required"; exit 1; }
command -v protoc >/dev/null 2>&1 || { echo "protoc is required"; exit 1; }
command -v protoc-gen-go >/dev/null 2>&1 || { echo "protoc-gen-go is required"; exit 1; }
command -v protoc-gen-go-grpc >/dev/null 2>&1 || { echo "protoc-gen-go-grpc is required"; exit 1; }

PROTO_DIR="$ROOT/proto"
cd "$PROTO_DIR"

flatten_pb() {
  local out="$1"
  local pb_dir="../$out/pb"
  local nested_dir="$pb_dir/short-drama-recommend"
  local relative_out="${out#rpc/}"
  local nested_package_dir="$nested_dir/$relative_out"
  if [[ ! -d "$nested_package_dir" ]]; then return; fi
  echo "Flattening generated PB: $nested_package_dir -> $pb_dir"
  find "$nested_package_dir" -maxdepth 1 -type f -exec mv -f {} "$pb_dir/" \;
  rm -rf "$nested_dir"
}

generate_rpc() {
  local proto="$1"
  local out="$2"
  echo "Generating RPC: $proto -> $out"
  mkdir -p "../$out/pb"
  goctl rpc protoc "$proto" \
    --go_out="../$out/pb" \
    --go_opt=module=short-drama-recommend \
    --go-grpc_out="../$out/pb" \
    --go-grpc_opt=module=short-drama-recommend \
    --zrpc_out="../$out"
  flatten_pb "$out"
}

generate_rpc user.proto rpc/user-rpc
generate_rpc drama.proto rpc/drama-rpc
generate_rpc drama_admin.proto rpc/drama-rpc
generate_rpc behavior.proto rpc/behavior-rpc
generate_rpc recommend.proto rpc/recommend-rpc
generate_rpc payment.proto rpc/payment-rpc
generate_rpc media.proto rpc/media-rpc

echo "go-zero RPC generation completed"