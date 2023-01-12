package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"net/url"
)

// VendorSchema holds the schema definition for the VendorSchema entity.
type VendorSchema struct {
	ent.Schema
}

// Fields of the VendorSchema.
func (VendorSchema) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.Text("address").Optional(),
		field.Bool("closed").Default(true),
		field.Text("Description").Optional(),
		field.JSON("image_url", &url.URL{}).Optional(),
	}
}

// Edges of the VendorSchema.
func (VendorSchema) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("VendorSchema").Required().Unique(),
		edge.To("items", Item.Type),
		edge.To("orders", Order.Type),
	}
}
