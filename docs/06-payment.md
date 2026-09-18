# Stripe + PayPal 支付

## 目标

海外短剧购买统一走后端支付服务。Tauri + Vue3 客户端只负责发起订单、展示支付界面和查询订单状态，不保存 Stripe Secret Key 或 PayPal Client Secret。

## Stripe

V1 采用 PaymentIntent：
1. 创建业务订单。
2. 服务端创建 PaymentIntent。
3. 客户端拿到 client_secret 完成支付。
4. 服务端处理 payment_intent.succeeded / payment_intent.payment_failed。
5. 成功后写入 user_entitlements，解锁剧集。

## PayPal

V1 采用 Orders v2：
1. 服务端创建 PayPal Order，intent=CAPTURE。
2. 返回 approve URL。
3. 客户端打开 PayPal 授权页。
4. 授权后服务端 Capture。
5. 监听 PAYMENT.CAPTURE.COMPLETED / PAYMENT.CAPTURE.DENIED。
6. 成功后写入 user_entitlements。

## 配置

STRIPE_SECRET_KEY
STRIPE_WEBHOOK_SECRET
PAYPAL_CLIENT_ID
PAYPAL_CLIENT_SECRET
PAYPAL_ENV=sandbox
PAYPAL_WEBHOOK_ID

## 本地验收

创建订单 -> Stripe 测试支付或 PayPal Sandbox -> webhook -> 订单 PAID -> entitlement -> 付费剧集解锁。

生产环境 webhook 必须使用公网 HTTPS，并进行签名/消息验证。
