package repository

import (
	"context"
	"database/sql"

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

func (r *MySQLDramaRepository) GetByID(ctx context.Context, id uint64) (*model.Drama, error) {
	const q = `SELECT id,title,description,cover,country,language,total_episodes,is_paid,price_cents,currency,status
FROM dramas WHERE id=? LIMIT 1`
	var d model.Drama
	var paid int8
	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&d.ID, &d.Title, &d.Description, &d.Cover, &d.Country,
		&d.Language, &d.TotalEpisodes, &paid, &d.PriceCents, &d.Currency, &d.Status,
	)
	if err != nil {
		return nil, err
	}
	d.IsPaid = paid == 1
	return &d, nil
}

func (r *MySQLDramaRepository) ListPublished(ctx context.Context, country, language string, limit, offset int) ([]*model.Drama, error) {
	const q = `SELECT id,title,description,cover,country,language,total_episodes,is_paid,status
FROM dramas
WHERE status=1 AND (country=? OR country='') AND (language=? OR language='')
ORDER BY published_at DESC, id DESC
LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, q, country, language, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*model.Drama, 0, limit)
	for rows.Next() {
		var d model.Drama
		var paid int8
		if err := rows.Scan(&d.ID, &d.Title, &d.Description, &d.Cover, &d.Country, &d.Language, &d.TotalEpisodes, &paid, &d.Status); err != nil {
			return nil, err
		}
		d.IsPaid = paid == 1
		result = append(result, &d)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
