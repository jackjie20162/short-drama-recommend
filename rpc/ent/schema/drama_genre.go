package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type DramaGenre struct { ent.Schema }

func (DramaGenre) Annotations() []schema.Annotation {
	return []schema.Annotation{field.ID("drama_id", "genre_id")}
}

func (DramaGenre) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("drama_id"),
		field.Uint64("genre_id"),
	}
}

func (DramaGenre) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("drama", Drama.Type).Required().Unique().Field("drama_id"),
		edge.To("genre", Genre.Type).Required().Unique().Field("genre_id"),
	}
}

func (DramaGenre) Indexes() []ent.Index {
	return []ent.Index{index.Fields("drama_id", "genre_id").Unique()}
}
