package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Restaurant holds the schema definition for the Restaurant entity.
type Restaurant struct {
	ent.Schema
}

func (Restaurant) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields of the Restaurant.
func (Restaurant) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty(),
		field.String("description").
			Optional(),
		field.String("address").
			NotEmpty(),
		field.Float("latitude"),
		field.Float("longitude"),
		field.String("phone").
			Optional(),
		field.String("website").
			Optional(),
		field.String("opening_hours").
			Optional(),
		field.Enum("status").
			Values("pending", "approved", "rejected").
			Default("pending"),
	}
}

// Edges of the Restaurant.
func (Restaurant) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("city", City.Type).
			Ref("restaurants").
			Unique(),
		edge.From("owner", User.Type).
			Ref("restaurants").
			Unique(),
		edge.To("food_items", FoodItem.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}
