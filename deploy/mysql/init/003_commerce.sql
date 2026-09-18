-- Pricing and commerce migration.
USE short_drama;

ALTER TABLE dramas
  ADD COLUMN price_cents BIGINT UNSIGNED NOT NULL DEFAULT 0 AFTER is_paid,
  ADD COLUMN currency VARCHAR(8) NOT NULL DEFAULT 'USD' AFTER price_cents;

CREATE TABLE IF NOT EXISTS orders (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  order_no VARCHAR(64) NOT NULL,
  user_id BIGINT UNSIGNED NOT NULL,
  drama_id BIGINT UNSIGNED NOT NULL,
  provider VARCHAR(16) NOT NULL,
  provider_order_id VARCHAR(128) NOT NULL DEFAULT '',
  client_secret VARCHAR(512) NOT NULL DEFAULT '',
  approve_url VARCHAR(1024) NOT NULL DEFAULT '',
  amount_cents BIGINT UNSIGNED NOT NULL,
  currency VARCHAR(8) NOT NULL DEFAULT 'USD',
  status VARCHAR(24) NOT NULL DEFAULT 'PENDING',
  paid_at DATETIME NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uk_order_no (order_no),
  KEY idx_user_status (user_id,status,created_at),
  KEY idx_provider_order (provider,provider_order_id)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS user_entitlements (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  user_id BIGINT UNSIGNED NOT NULL,
  drama_id BIGINT UNSIGNED NOT NULL,
  order_id BIGINT UNSIGNED NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uk_user_drama (user_id,drama_id),
  KEY idx_order (order_id)
) ENGINE=InnoDB;
