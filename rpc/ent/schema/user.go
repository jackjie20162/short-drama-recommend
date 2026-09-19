package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type User struct { ent.Schema }

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id").Immutable(),
		field.String("external_id").Default(""),
		field.String("country").Default(""),
		field.String("language").Default("en"),
		field.String("locale").Default("en-US"),
		field.String("timezone").Default("UTC"),
		field.Int8("status").Default(1),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("external_id").Unique(),
		index.Fields("country", "language"),
	}
}
