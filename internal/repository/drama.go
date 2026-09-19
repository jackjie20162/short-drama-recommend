package repository

import (
	"context"
	"database/sql"
	"strings"

	"short-drama-recommend/internal/model"
)

type DramaRepository interface {
	GetByID(ctx context.Context, id uint64) (*model.Drama, error)
	ListPublished(ctx context.Context, country, language string, limit, offset int) ([]*model.Drama, error)
}

type MySQLDramaRepository struct {
	db *sql.DB
}

func NewMySQLDramaRepository(db *sql.DB) *MySQLDramaRepository {
	return &MySQLDramaRepository{db: db}
}

const dramaSelect = `
SELECT
	d.id, d.title, d.subtitle, d.description, d.cover, d.country, d.language,
	d.total_episodes, d.is_paid, d.price_cents, d.currency, d.status,
	d.popularity, d.completion_rate, d.pay_rate, d.published_at,
	COALESCE((SELECT GROUP_CONCAT(g.name ORDER BY g.id SEPARATOR '||')
		FROM drama_genres dg JOIN genres g ON g.id=dg.genre_id
		WHERE dg.drama_id=d.id AND g.language=d.language), '') AS genres,
	COALESCE((SELECT GROUP_CONCAT(t.name ORDER BY t.id SEPARATOR '||')
		FROM drama_tags dt JOIN tags t ON t.id=dt.tag_id
		WHERE dt.drama_id=d.id AND t.language=d.language), '') AS tags
FROM dramas d
`

func scanDrama(row interface{ Scan(...any) error }) (*model.Drama, error) {
	var d model.Drama
	var paid int8
	var published sql.NullTime
	var genres, tags string
	if err := row.Scan(
		&d.ID, &d.Title, &d.Subtitle, &d.Description, &d.Cover, &d.Country, &d.Language,
		&d.TotalEpisodes, &paid, &d.PriceCents, &d.Currency, &d.Status,
		&d.Popularity, &d.CompletionRate, &d.PayRate, &published, &genres, &tags,
	); err != nil {
		return nil, err
	}
	d.IsPaid = paid == 1
	if published.Valid {
		d.PublishedAt = published.Time
	}
	d.Genres = splitContentValues(genres)
	d.Tags = splitContentValues(tags)
	return &d, nil
}

func splitContentValues(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, "||")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

func (r *MySQLDramaRepository) GetByID(ctx context.Context, id uint64) (*model.Drama, error) {
	row := r.db.QueryRowContext(ctx, dramaSelect+" WHERE d.id=? LIMIT 1", id)
	return scanDrama(row)
}

func (r *MySQLDramaRepository) ListPublished(ctx context.Context, country, language string, limit, offset int) ([]*model.Drama, error) {
	rows, err := r.db.QueryContext(ctx,
		dramaSelect+`
WHERE d.status=1
  AND (d.country=? OR d.country='')
  AND (d.language=? OR d.language='')
ORDER BY d.published_at DESC, d.id DESC
LIMIT ? OFFSET ?`, country, language, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*model.Drama, 0, limit)
	for rows.Next() {
		d, err := scanDrama(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, d)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

