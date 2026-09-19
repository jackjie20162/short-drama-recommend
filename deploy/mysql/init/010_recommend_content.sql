USE short_drama;

-- V1.1 content/recommendation schema.
-- Business source remains MySQL; search/recommend layers derive from these fields.

ALTER TABLE dramas
  ADD COLUMN subtitle VARCHAR(255) NOT NULL DEFAULT '' AFTER title,
  ADD COLUMN popularity DECIMAL(12,6) NOT NULL DEFAULT 0 AFTER currency,
  ADD COLUMN completion_rate DECIMAL(8,6) NOT NULL DEFAULT 0 AFTER popularity,
  ADD COLUMN pay_rate DECIMAL(8,6) NOT NULL DEFAULT 0 AFTER completion_rate;

ALTER TABLE episodes
  ADD COLUMN description TEXT NOT NULL AFTER title;

CREATE TABLE IF NOT EXISTS genres (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  name VARCHAR(64) NOT NULL,
  language VARCHAR(16) NOT NULL DEFAULT 'en',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uk_genre_name_language (name, language)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS drama_genres (
  drama_id BIGINT UNSIGNED NOT NULL,
  genre_id BIGINT UNSIGNED NOT NULL,
  PRIMARY KEY (drama_id, genre_id),
  KEY idx_genre_drama (genre_id, drama_id)
) ENGINE=InnoDB;

ALTER TABLE behavior_events
  ADD COLUMN source VARCHAR(64) NOT NULL DEFAULT '' AFTER device,
  ADD COLUMN request_id VARCHAR(128) NOT NULL DEFAULT '' AFTER source,
  ADD KEY idx_request_id (request_id),
  ADD KEY idx_user_drama_event (user_id, drama_id, event_at);

CREATE TABLE IF NOT EXISTS drama_feature_snapshots (
  drama_id BIGINT UNSIGNED NOT NULL,
  feature_version VARCHAR(32) NOT NULL,
  title_text TEXT NOT NULL,
  description_text TEXT NOT NULL,
  sparse_json JSON NOT NULL,
  dense_json JSON NOT NULL,
  semantic_json JSON NOT NULL,
  built_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (drama_id, feature_version),
  KEY idx_built_at (built_at)
) ENGINE=InnoDB;
