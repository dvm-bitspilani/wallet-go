package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"time"
)

// OrderShell holds the schema definition for the OrderShell entity.
type OrderShell struct {
	ent.Schema
}

// Fields of the OrderShell.
func (OrderShell) Fields() []ent.Field {
	return []ent.Field{
		field.Int("price").Default(0).NonNegative(),
		field.Time("timestamp").Default(time.Now),
	}
}

// Edges of the OrderShell.
func (OrderShell) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("wallet", Wallet.Type).Ref("shells").Unique(),
		edge.To("orders", Order.Type),
	}
}
