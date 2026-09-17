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

INSERT INTO dramas (title, description, cover, country, language, total_episodes, is_paid, status, published_at)
VALUES
('The Contract Wife', 'A fictional short drama used for local development and recommendation testing.', '', 'US', 'en', 20, 0, 1, NOW()),
('Revenge of the Hidden Heiress', 'A fictional revenge romance drama for development seed data.', '', 'US', 'en', 24, 0, 1, NOW()),
('CEO Next Door', 'A fictional modern romance drama for feed and ranking tests.', '', 'GB', 'en', 18, 0, 1, NOW())
ON DUPLICATE KEY UPDATE title = VALUES(title);
