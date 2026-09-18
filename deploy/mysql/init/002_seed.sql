USE short_drama;

INSERT INTO tags (name, language) VALUES
('Romance', 'en'),
('Billionaire', 'en'),
('Revenge', 'en'),
('Contract Marriage', 'en'),
('CEO', 'en'),
('Family', 'en'),
('Fantasy', 'en'),
('Thriller', 'en')
ON DUPLICATE KEY UPDATE name = VALUES(name);

INSERT INTO dramas (title, description, cover, country, language, total_episodes, is_paid, price_cents, currency, status, published_at)
VALUES
('The Contract Wife', 'A fictional short drama used for local development and recommendation testing.', '', 'US', 'en', 20, 1, 499, 'USD', 1, NOW()),
('Revenge of the Hidden Heiress', 'A fictional revenge romance drama for development seed data.', '', 'US', 'en', 24, 1, 699, 'USD', 1, NOW()),
('CEO Next Door', 'A fictional modern romance drama for feed and ranking tests.', '', 'GB', 'en', 18, 0, 0, 'USD', 1, NOW())
ON DUPLICATE KEY UPDATE title = VALUES(title), is_paid=VALUES(is_paid), price_cents=VALUES(price_cents), currency=VALUES(currency), status=VALUES(status);


INSERT INTO episodes (drama_id, episode_no, title, duration_seconds, video_url, poster_url, is_paid, status)
SELECT d.id, x.episode_no, CONCAT('Episode ', x.episode_no), 90, CONCAT('https://example.com/video/', d.id, '/', x.episode_no, '.m3u8'), '', 0, 1
FROM dramas d
JOIN (
  SELECT 1 episode_no UNION ALL SELECT 2 UNION ALL SELECT 3
) x
WHERE d.title IN ('The Contract Wife','Revenge of the Hidden Heiress','CEO Next Door')
ON DUPLICATE KEY UPDATE title=VALUES(title), video_url=VALUES(video_url);
