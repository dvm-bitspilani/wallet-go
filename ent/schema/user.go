package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"time"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("username").
			Unique(),
		field.String("password").
			Sensitive(),
		field.String("email"),
		field.String("name"),
		field.UUID("qr_code", uuid.UUID{}).
			Default(uuid.New),
		field.Time("created").
			Default(time.Now),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("wallet", Wallet.Type).Unique(),
		edge.To("teller", Teller.Type).Unique(),
		edge.To("transactions", Transactions.Type),
		//edge.To("pg_transactions"),
		edge.To("vendor", Vendor.Type).Unique(),
	}
}
