package schema

import (
	"dvm.wallet/harsh/internal/helpers"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"time"
)

// Order holds the schema definition for the Order entity.
type Order struct {
	ent.Schema
}

// Fields of the Order.
func (Order) Fields() []ent.Field {
	return []ent.Field{
		field.Int("price").Default(0).NonNegative(),
		field.Enum("status").GoType(helpers.Status(0)),
		field.Int("otp").DefaultFunc(func() int {
			return 0 // TODO: Implement a OTP generating function
		}).NonNegative(),
		field.Bool("otp_seen").Default(false),
		field.Time("timestamp").Default(time.Now),
		field.Time("ready_timestamp").Optional(),
		field.Time("accepted_timestamp").Optional(),
		field.Time("finished_timestamp").Optional(),
		field.Time("declined_timestamp").Optional(),
	}
}

// Edges of the Order.
func (Order) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("shell", OrderShell.Type).Ref("orders").Unique(),
		edge.From("vendor", Vendor.Type).Ref("orders").Unique(),
		edge.From("transaction", Transactions.Type).Ref("orders").Unique(),
		edge.To("iteminstances", ItemInstance.Type),
	}
}
