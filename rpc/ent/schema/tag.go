package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Tag struct { ent.Schema }

func (Tag) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id").Immutable(),
		field.String("name"),
		field.String("language").Default("en"),
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}

func (Tag) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("dramas", Drama.Type).Ref("tags").Through("drama_tags", DramaTag.Type),
	}
}

func (Tag) Indexes() []ent.Index {
	return []ent.Index{index.Fields("name", "language").Unique()}
}
