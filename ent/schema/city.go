package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// City holds the schema definition for the City entity.
type City struct {
	ent.Schema
}

func (City) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields of the City.
func (City) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			Unique(),
		field.String("state").
			Default("Karnataka"),
		field.Float("latitude").
			Default(12.9716),
		field.Float("longitude").
			Default(77.5946),
	}
}

// Edges of the City.
func (City) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("restaurants", Restaurant.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}
