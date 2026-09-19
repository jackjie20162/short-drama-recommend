package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type BehaviorEvent struct { ent.Schema }

func (BehaviorEvent) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id").Immutable(),
		field.Uint64("user_id").Default(0),
		field.Uint64("drama_id").Default(0),
		field.Uint64("episode_id").Default(0),
		field.String("event_type"),
		field.Uint32("watch_seconds").Default(0),
		field.Uint32("duration_seconds").Default(0),
		field.String("country").Default(""),
		field.String("language").Default(""),
		field.String("device").Default(""),
		field.String("source").Default(""),
		field.String("request_id").Default(""),
		field.Time("event_at").Default(time.Now),
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}

func (BehaviorEvent) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "event_at"),
		index.Fields("drama_id", "event_type", "event_at"),
		index.Fields("event_at"),
		index.Fields("request_id"),
		index.Fields("user_id", "drama_id", "event_at"),
	}
}
