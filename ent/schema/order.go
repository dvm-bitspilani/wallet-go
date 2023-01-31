package schema

import (
	"crypto/rand"
	"dvm.wallet/harsh/internal/helpers"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"io"
	"time"
)

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

// Order holds the schema definition for the Order entity.
type Order struct {
	ent.Schema
}

// Fields of the Order.
func (Order) Fields() []ent.Field {
	return []ent.Field{
		field.Int("price").Default(0).NonNegative(),
		field.Enum("status").GoType(helpers.Status(0)),
		field.String("otp").DefaultFunc(func() string {
			b := make([]byte, 6)
			n, err := io.ReadAtLeast(rand.Reader, b, 6)
			if n != 6 {
				panic(err)
			}
			for i := 0; i < len(b); i++ {
				b[i] = table[int(b[i])%len(table)]
			}
			return string(b)
		}),
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
		edge.From("VendorSchema", VendorSchema.Type).Ref("orders").Unique(),
		edge.From("transaction", Transactions.Type).Ref("orders").Unique(),
		edge.To("iteminstances", ItemInstance.Type),
	}
}
