package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type DramaTag struct { ent.Schema }

func (DramaTag) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("drama_id"),
		field.Uint64("tag_id"),
	}
}

func (DramaTag) Indexes() []ent.Index {
	return []ent.Index{index.Fields("drama_id", "tag_id").Unique()}
}
