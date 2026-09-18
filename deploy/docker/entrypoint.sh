#!/bin/sh
set -eu
SERVICE="$1"
CONFIG="$2"
shift 2
sed -i -e 's/127\.0\.0\.1:3307/mysql:3306/g' -e 's/127\.0\.0\.1:6380/redis:6379/g' -e 's/127\.0\.0\.1:9002/drama-rpc:9002/g' -e 's/127\.0\.0\.1:9003/behavior-rpc:9003/g' -e 's/127\.0\.0\.1:9004/recommend-rpc:9004/g' "$CONFIG"
exec "/app/$SERVICE" -f "$CONFIG" "$@"
