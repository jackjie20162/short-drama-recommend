package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Episode struct { ent.Schema }

func (Episode) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id").Immutable(),
		field.Uint64("drama_id"),
		field.Uint32("episode_no"),
		field.String("title").Default(""),
		field.String("description").Default(""),
		field.Uint32("duration_seconds").Default(0),
		field.String("video_url").Default(""),
		field.String("poster_url").Default(""),
		field.Bool("is_paid").Default(false),
		field.Int8("status").Default(1),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (Episode) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("drama", Drama.Type).Ref("episodes").Field("drama_id").Unique().Required(),
	}
}

func (Episode) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("drama_id", "episode_no").Unique(),
		index.Fields("drama_id", "status"),
	}
}
