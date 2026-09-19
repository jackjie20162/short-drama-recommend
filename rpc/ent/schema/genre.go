package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Genre struct { ent.Schema }

func (Genre) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id").Immutable(),
		field.String("name"),
		field.String("language").Default("en"),
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}

func (Genre) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("dramas", Drama.Type).Ref("genres").Through("drama_genres", DramaGenre.Type),
	}
}

func (Genre) Indexes() []ent.Index {
	return []ent.Index{index.Fields("name", "language").Unique()}
}
