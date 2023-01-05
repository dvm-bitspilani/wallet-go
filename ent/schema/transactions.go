package schema

import (
	"database/sql/driver"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"time"
)

type txn_type int

// ADD_SWD = 0
// ADD_CASH = 1
// ADD_PG = 2
// TRANSFER = 3
// PURCHASE = 4
const (
	ADD_SWD  txn_type = 0
	ADD_CASH txn_type = 1
	ADD_PG   txn_type = 2
	TRANSFER txn_type = 3
	PURCHASE txn_type = 4
)

func (t txn_type) String() string {
	switch t {
	case ADD_SWD:
		return "Add from SWD"
	case ADD_CASH:
		return "Add from Cash"
	case ADD_PG:
		return "Add from Payment Gateway"
	case TRANSFER:
		return "Transfer"
	case PURCHASE:
		return "Purchase"
	}
	return ""
}

// Values provides list valid values for Enum.
func (t txn_type) Values() []string {
	return []string{ADD_SWD.String(), ADD_CASH.String(), ADD_PG.String(), TRANSFER.String(), PURCHASE.String()}
}

// Value provides the DB a string from int.
func (t txn_type) Value() (driver.Value, error) {
	return t.String(), nil
}

// Scan tells our code how to read the enum into our type.
func (t *txn_type) Scan(val any) error {
	var s string
	switch v := val.(type) {
	case nil:
		return nil
	case string:
		s = v
	case []uint8:
		s = string(v)
	}
	switch s {
	case "Add from SWD":
		*t = ADD_SWD
	case "Add from Cash":
		*t = ADD_CASH
	case "Add from Payment Gateway":
		*t = ADD_PG
	case "Transfer":
		*t = PURCHASE
	case "Purchase":
		*t = TRANSFER
	}
	return nil
}

// Transactions holds the schema definition for the Transactions entity.
type Transactions struct {
	ent.Schema
}

// Fields of the Transactions.
func (Transactions) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("kind").GoType(txn_type(0)), //TODO: check if this default to 0 or not
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
