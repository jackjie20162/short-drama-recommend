package repository

import (
	"context"
	"database/sql"

	"short-drama-recommend/internal/model"
)

type BehaviorRepository interface {
	Record(ctx context.Context, event *model.BehaviorEvent) error
}

type MySQLBehaviorRepository struct { db *sql.DB }

func NewMySQLBehaviorRepository(db *sql.DB) BehaviorRepository {
	return &MySQLBehaviorRepository{db: db}
}

func (r *MySQLBehaviorRepository) Record(ctx context.Context, event *model.BehaviorEvent) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO behavior_events
		(user_id, drama_id, episode_id, event_type, watch_seconds, duration_seconds, country, language, device, event_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		event.UserID, event.DramaID, event.EpisodeID, event.EventType,
		event.WatchSeconds, event.DurationSeconds, event.Country, event.Language, event.Device, event.EventAt)
	return err
}
