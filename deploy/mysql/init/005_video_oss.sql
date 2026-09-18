USE short_drama;

-- Video files are stored in OSS/CDN, never in MySQL.
-- episode.video_url remains backward-compatible; new fields describe the OSS asset.
ALTER TABLE episodes ADD COLUMN IF NOT EXISTS video_format VARCHAR(16) NOT NULL DEFAULT 'm3u8';
ALTER TABLE episodes ADD COLUMN IF NOT EXISTS video_storage VARCHAR(16) NOT NULL DEFAULT 'oss';
ALTER TABLE episodes ADD COLUMN IF NOT EXISTS video_object_key VARCHAR(1024) NOT NULL DEFAULT '';
ALTER TABLE episodes ADD COLUMN IF NOT EXISTS video_playback_url VARCHAR(2048) NOT NULL DEFAULT '';
ALTER TABLE episodes ADD COLUMN IF NOT EXISTS video_status VARCHAR(24) NOT NULL DEFAULT 'READY';
ALTER TABLE episodes ADD COLUMN IF NOT EXISTS video_size_bytes BIGINT UNSIGNED NOT NULL DEFAULT 0;
ALTER TABLE episodes ADD COLUMN IF NOT EXISTS video_checksum VARCHAR(128) NOT NULL DEFAULT '';
ALTER TABLE episodes ADD COLUMN IF NOT EXISTS video_updated_at DATETIME NULL;

CREATE INDEX IF NOT EXISTS idx_episode_video_status ON episodes(video_status);
