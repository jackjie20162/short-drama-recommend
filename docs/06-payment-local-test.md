# Payment local test

## Environment

Copy `.env.example` to a local environment file or export variables before Docker Compose:
- `STRIPE_SECRET_KEY`
- `STRIPE_WEBHOOK_SECRET`
- `PAYPAL_CLIENT_ID`
- `PAYPAL_CLIENT_SECRET`
- `PAYPAL_ENV=sandbox`
- `PAYPAL_WEBHOOK_ID`
- `VITE_STRIPE_PUBLISHABLE_KEY` in `client-tauri/.env.local`

## Backend

Start:
```bash
docker compose up -d --build
```

Payment RPC listens on `9005`.

REST payment endpoints:
- `POST /api/v1/payments/orders`
- `POST /api/v1/payments/capture`
- `GET /api/v1/payments/orders`
- `GET /api/v1/payments/orders/:id`
- `POST /api/v1/payments/webhook`
- `GET /api/v1/dramas/:id/episodes?user_id=1`

## Stripe

The client creates a PaymentIntent through the backend and confirms it with Stripe.js/Card Element. The server marks the order paid from a verified `payment_intent.succeeded` webhook and creates `user_entitlements`.

Webhook: `https://<your-api-domain>/api/v1/payments/webhook?provider=stripe`
Use the webhook signing secret as `STRIPE_WEBHOOK_SECRET`.

## PayPal

The backend creates a PayPal Orders v2 order. The client opens the PayPal approval URL. The return page sends the PayPal token back to the backend for Capture. The server also accepts verified `PAYMENT.CAPTURE.COMPLETED` webhook events and grants the entitlement.

Webhook: `https://<your-api-domain>/api/v1/payments/webhook?provider=paypal`
Set the PayPal webhook ID in `PAYPAL_WEBHOOK_ID`.

## Client

```bash
cd client-tauri
npm install
npm run build
npm run tauri dev
```

The client contains Home, Discover, Drama Detail, Checkout, Orders and Profile views. Paid episode video URLs are withheld by the API until the user has a matching entitlement.
