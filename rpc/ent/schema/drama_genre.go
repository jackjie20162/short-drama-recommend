package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type DramaGenre struct { ent.Schema }

func (DramaGenre) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("drama_id"),
		field.Uint64("genre_id"),
	}
}

func (DramaGenre) Indexes() []ent.Index {
	return []ent.Index{index.Fields("drama_id", "genre_id").Unique()}
}
