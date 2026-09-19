package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Drama struct { ent.Schema }

func (Drama) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id").Immutable(),
		field.String("title"),
		field.String("subtitle").Default(""),
		field.String("description"),
		field.String("cover").Default(""),
		field.String("country").Default(""),
		field.String("language").Default("en"),
		field.Uint32("total_episodes").Default(0),
		field.Bool("is_paid").Default(false),
		field.Uint64("price_cents").Default(0),
		field.String("currency").Default("USD"),
		field.Int8("status").Default(1),
		field.Float64("popularity").Default(0),
		field.Float64("completion_rate").Default(0),
		field.Float64("pay_rate").Default(0),
		field.Time("published_at").Optional().Nillable(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (Drama) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("episodes", Episode.Type),
		edge.To("genres", Genre.Type).Through("drama_genres", DramaGenre.Type),
		edge.To("tags", Tag.Type).Through("drama_tags", DramaTag.Type),
	}
}

func (Drama) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("title", "country", "language").Unique(),
		index.Fields("status", "published_at"),
		index.Fields("country", "language"),
	}
}
