package schema

import (
	"dvm.wallet/harsh/internal/helpers"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"time"
)

// Transactions holds the schema definition for the Transactions entity.
type Transactions struct {
	ent.Schema
}

// Fields of the Transactions.
func (Transactions) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("kind").GoType(helpers.Txn_type(0)), //TODO: check if this default to 0 or not
		field.Int("amount").Default(0).NonNegative(),
		field.Time("timestamp").Default(time.Now),
	}
}

// Edges of the Transactions.
func (Transactions) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("transactions").Unique(),
		edge.To("orders", Order.Type),
		edge.From("source", Wallet.Type).Ref("source_transactions").Unique(),
		edge.From("destination", Wallet.Type).Ref("destination_transactions").Unique(),
	}
}
