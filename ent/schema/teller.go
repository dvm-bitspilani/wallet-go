package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Teller holds the schema definition for the Teller entity.
type Teller struct {
	ent.Schema
}

// Fields of the Teller.
func (Teller) Fields() []ent.Field {
	return []ent.Field{
		field.Int("cash_collected").Default(0).NonNegative(),
	}
}

// Edges of the Teller.
func (Teller) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("teller").Unique().Required(),
	}
}
