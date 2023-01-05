package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// ItemInstance holds the schema definition for the ItemInstance entity.
type ItemInstance struct {
	ent.Schema
}

// Fields of the ItemInstance.
func (ItemInstance) Fields() []ent.Field {
	return []ent.Field{
		field.Int("quantity").Default(1).Positive(),
		field.Int("price_per_quantity").Default(0).NonNegative(),
	}
}

// Edges of the ItemInstance.
func (ItemInstance) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("item", Item.Type).Ref("iteminstances").Unique(),
		edge.From("order", Order.Type).Ref("iteminstances").Unique(),
	}
}
