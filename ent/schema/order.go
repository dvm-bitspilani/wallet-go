package schema

import (
	"database/sql/driver"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"time"
)

type status int

// PENDING = 0
// ACCEPTED = 1
// READY = 2
// FINISHED = 3
// DECLINED = 4
const (
	PENDING  status = 0
	ACCEPTED status = 1
	READY    status = 2
	FINISHED status = 3
	DECLINED status = 4
)

func (s status) String() string {
	switch s {
	case PENDING:
		return "Pending"
	case ACCEPTED:
		return "Accepted"
	case READY:
		return "Ready"
	case FINISHED:
		return "Finished"
	case DECLINED:
		return "Declined"
	}
	return ""
}

// Values provides list valid values for Enum.
func (s status) Values() []string {
	return []string{PENDING.String(), ACCEPTED.String(), READY.String(), FINISHED.String(), DECLINED.String()}
}

// Value provides the DB a string from int.
func (s status) Value() (driver.Value, error) {
	return s.String(), nil
}

// Scan tells our code how to read the enum into our type.
func (s *status) Scan(val any) error {
	var x string
	switch v := val.(type) {
	case nil:
		return nil
	case string:
		x = v
	case []uint8:
		x = string(v)
	}
	switch x {
	case "Pending":
		*s = PENDING
	case "Accepted":
		*s = ACCEPTED
	case "Ready":
		*s = READY
	case "Finished":
		*s = FINISHED
	case "Declined":
		*s = DECLINED
	}
	return nil
}

// Order holds the schema definition for the Order entity.
type Order struct {
	ent.Schema
}

// Fields of the Order.
func (Order) Fields() []ent.Field {
	return []ent.Field{
		field.Int("price").Default(0).NonNegative(),
		field.Enum("status").GoType(status(0)),
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
