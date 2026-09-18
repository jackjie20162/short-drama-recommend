USE short_drama;

ALTER TABLE dramas
  ADD COLUMN price_cents BIGINT NOT NULL DEFAULT 0 AFTER is_paid,
  ADD COLUMN currency VARCHAR(8) NOT NULL DEFAULT 'USD' AFTER price_cents;

CREATE TABLE IF NOT EXISTS orders (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  order_no VARCHAR(64) NOT NULL UNIQUE,
  user_id BIGINT NOT NULL,
  drama_id BIGINT NOT NULL,
  provider VARCHAR(16) NOT NULL,
  provider_order_id VARCHAR(128) NULL,
  amount_cents BIGINT NOT NULL,
  currency VARCHAR(8) NOT NULL,
  status VARCHAR(32) NOT NULL DEFAULT 'PENDING',
  client_secret VARCHAR(255) NULL,
  approve_url VARCHAR(1000) NULL,
  paid_at DATETIME NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY idx_orders_user (user_id, created_at),
  KEY idx_orders_provider (provider, provider_order_id),
  KEY idx_orders_drama (drama_id, created_at)
);

CREATE TABLE IF NOT EXISTS user_entitlements (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT NOT NULL,
  drama_id BIGINT NOT NULL,
  order_id BIGINT NOT NULL,
  source VARCHAR(32) NOT NULL DEFAULT 'PURCHASE',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY uk_user_drama (user_id, drama_id),
  KEY idx_entitlement_order (order_id)
);

UPDATE dramas SET price_cents=499, currency='USD', is_paid=1 WHERE title='The Contract Wife';
UPDATE dramas SET price_cents=699, currency='USD', is_paid=1 WHERE title='Revenge of the Hidden Heiress';
UPDATE dramas SET price_cents=399, currency='USD', is_paid=1 WHERE title='CEO Next Door';
