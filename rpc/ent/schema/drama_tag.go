package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type DramaTag struct { ent.Schema }

func (DramaTag) Annotations() []schema.Annotation {
	return []schema.Annotation{field.ID("drama_id", "tag_id")}
}

func (DramaTag) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("drama_id"),
		field.Uint64("tag_id"),
	}
}

func (DramaTag) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("drama", Drama.Type).Required().Unique().Field("drama_id"),
		edge.To("tag", Tag.Type).Required().Unique().Field("tag_id"),
	}
}

func (DramaTag) Indexes() []ent.Index {
	return []ent.Index{index.Fields("drama_id", "tag_id").Unique()}
}
