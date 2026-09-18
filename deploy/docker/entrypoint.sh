#!/bin/sh
set -eu
SERVICE="$1"
CONFIG="$2"
shift 2
sed -i -e 's/127\.0\.0\.1:3307/mysql:3306/g' -e 's/127\.0\.0\.1:6380/redis:6379/g' -e 's/127\.0\.0\.1:9002/drama-rpc:9002/g' -e 's/127\.0\.0\.1:9003/behavior-rpc:9003/g' -e 's/127\.0\.0\.1:9004/recommend-rpc:9004/g' -e 's/127\.0\.0\.1:9005/payment-rpc:9005/g' -e "s|__STRIPE_SECRET_KEY__|${STRIPE_SECRET_KEY:-}|g" -e "s|__STRIPE_WEBHOOK_SECRET__|${STRIPE_WEBHOOK_SECRET:-}|g" -e "s|__PAYPAL_CLIENT_ID__|${PAYPAL_CLIENT_ID:-}|g" -e "s|__PAYPAL_CLIENT_SECRET__|${PAYPAL_CLIENT_SECRET:-}|g" -e "s|__PAYPAL_ENV__|${PAYPAL_ENV:-sandbox}|g" -e "s|__PAYPAL_WEBHOOK_ID__|${PAYPAL_WEBHOOK_ID:-}|g" "$CONFIG"
exec "/app/$SERVICE" -f "$CONFIG" "$@"
