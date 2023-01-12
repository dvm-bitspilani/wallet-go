package helpers

import "database/sql/driver"

type Status int

// PENDING = 0
// ACCEPTED = 1
// READY = 2
// FINISHED = 3
// DECLINED = 4
const (
	PENDING  Status = 0
	ACCEPTED Status = 1
	READY    Status = 2
	FINISHED Status = 3
	DECLINED Status = 4
)

func (s Status) String() string {
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

func FromInt(num int) Status {
	switch num {
	case 0:
		return PENDING
	case 1:
		return ACCEPTED
	case 2:
		return READY
	case 3:
		return FINISHED
	case 4:
		return DECLINED
	}
	return 0 // TODO:	check if this is right or not
}

// Values provides list valid values for Enum.
func (s Status) Values() []string {
	return []string{PENDING.String(), ACCEPTED.String(), READY.String(), FINISHED.String(), DECLINED.String()}
}

// Value provides the DB a string from int.
func (s Status) Value() (driver.Value, error) {
	return s.String(), nil
}

// Scan tells our code how to read the enum into our type.
func (s *Status) Scan(val any) error {
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

// ==============================================================

type Txn_type int

// ADD_SWD = 0
// ADD_CASH = 1
// ADD_PG = 2
// TRANSFER = 3
// PURCHASE = 4
const (
	ADD_SWD  Txn_type = 0
	ADD_CASH Txn_type = 1
	ADD_PG   Txn_type = 2
	TRANSFER Txn_type = 3
	PURCHASE Txn_type = 4
)

func (t Txn_type) String() string {
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
func (t Txn_type) Values() []string {
	return []string{ADD_SWD.String(), ADD_CASH.String(), ADD_PG.String(), TRANSFER.String(), PURCHASE.String()}
}

// Value provides the DB a string from int.
func (t Txn_type) Value() (driver.Value, error) {
	return t.String(), nil
}

// Scan tells our code how to read the enum into our type.
func (t *Txn_type) Scan(val any) error {
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
