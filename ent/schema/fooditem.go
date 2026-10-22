package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// FoodItem holds the schema definition for the FoodItem entity.
type FoodItem struct {
	ent.Schema
}

func (FoodItem) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields of the FoodItem.
func (FoodItem) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty(),
		field.String("description").
			Optional(),
		field.String("category").
			Default("General"),
		field.Float("price").
			Positive(),
		field.Bool("available").
			Default(true),
	}
}

// Edges of the FoodItem.
func (FoodItem) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("restaurant", Restaurant.Type).
			Ref("food_items").
			Unique().
			Required(),
	}
}
