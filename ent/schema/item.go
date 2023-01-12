package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Item holds the schema definition for the Item entity.
type Item struct {
	ent.Schema
}

// Fields of the Item.
func (Item) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.Int("base_price").Default(0).NonNegative(),
		field.Text("description").Optional(),
		field.Bool("available").Default(false),
		field.Bool("veg").Default(true),
	}
}

// Edges of the Item.
func (Item) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("VendorSchema", VendorSchema.Type).Ref("items").Unique(),
		edge.To("iteminstances", ItemInstance.Type),
	}
}
