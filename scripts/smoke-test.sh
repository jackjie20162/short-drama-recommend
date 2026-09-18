#!/usr/bin/env bash
set -euo pipefail
BASE_URL="${BASE_URL:-http://127.0.0.1:8080}"
echo "== short-drama smoke test =="
echo "BASE_URL=$BASE_URL"
wait_http(){ local url="$1"; for i in $(seq 1 30); do if curl -fsS "$url" >/dev/null 2>&1; then return 0; fi; sleep 2; done; echo "timeout waiting for $url" >&2; exit 1; }
wait_http "$BASE_URL/api/v1/dramas?page=1&page_size=20"
echo "[1] drama list"; curl -fsS "$BASE_URL/api/v1/dramas?page=1&page_size=20"; echo
echo "[2] drama detail"; curl -fsS "$BASE_URL/api/v1/dramas/1"; echo
echo "[3] recommendation feed"; curl -fsS "$BASE_URL/api/v1/feed?user_id=1&country=US&language=en&page_size=10"; echo
echo "[4] behavior event"; curl -fsS -X POST "$BASE_URL/api/v1/behaviors" -H 'Content-Type: application/json' -d '{"user_id":1,"drama_id":1,"episode_id":1,"event_type":"PLAY","watch_seconds":5,"duration_seconds":60,"country":"US","language":"en","device":"web"}'; echo
echo "Smoke test completed."
