package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Wallet holds the schema definition for the Wallet entity.
type Wallet struct {
	ent.Schema
}

// Fields of the Wallet.
func (Wallet) Fields() []ent.Field {
	return []ent.Field{
		field.Int("cash").NonNegative().Default(0),
		field.Int("pg").NonNegative().Default(0),
		field.Int("swd").NonNegative().Default(0),
		field.Int("transfers").NonNegative().Default(0),
	}
}

// Edges of the Wallet.
func (Wallet) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("wallet").Unique().Required(),
		//edge.To("source", Transactions.Type).,
		edge.To("shells", OrderShell.Type),
	}
}
