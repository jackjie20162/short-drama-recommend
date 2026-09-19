# 后台配置说明

系统设置已接入 admin-tauri → drama-admin-api → MySQL system_settings。

## 支付
后台可以维护 Stripe 与 PayPal：
- Stripe Secret Key / Webhook Secret
- PayPal Client ID / Client Secret / Webhook ID
- PayPal Sandbox / Live
- Provider 开关

支付 RPC 优先读取后台配置；如果数据库没有配置，则兼容读取 payment-rpc 的旧环境变量。

## 视频存储
后台配置：
- 默认 Provider：Aliyun OSS / AWS S3
- OSS Endpoint / Region / Bucket / AccessKey / SecretKey
- S3 Region / Bucket / Endpoint / AccessKey / SecretKey
- CDN Base URL

当前设置页已完成配置持久化和参数检查；实际 presigned upload、FFmpeg/HLS、CDN 签名由后续 media-rpc 接入。

## 密钥安全
Secret 字段使用 AES-256-GCM 加密存储。服务启动必须设置 SETTINGS_ENCRYPTION_KEY。

必须是 32 字节随机值。该值只放服务端环境变量，不进入 Git，不下发客户端。

## 数据库
新数据库会自动执行 007_system_settings.sql。
已有 Docker MySQL 数据卷不会重复执行 init 脚本，需要手动执行对应 SQL。
