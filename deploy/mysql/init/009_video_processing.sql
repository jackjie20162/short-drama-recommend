USE short_drama;

ALTER TABLE episodes ADD COLUMN IF NOT EXISTS video_processing_error TEXT NOT NULL DEFAULT '';
