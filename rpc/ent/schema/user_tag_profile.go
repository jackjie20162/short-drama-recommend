package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type UserTagProfile struct { ent.Schema }

func (UserTagProfile) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("user_id"),
		field.Uint64("tag_id"),
		field.Float64("weight").Default(0),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (UserTagProfile) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "tag_id").Unique(),
		index.Fields("user_id", "weight"),
	}
}
