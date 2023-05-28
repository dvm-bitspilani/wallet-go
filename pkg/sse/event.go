package sse

import "encoding/json"

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

const (
	EventUpdateBalance   = "update_balance"
	EventUpdateOrder     = "update_order_status"
	EventPutVendorOrders = "put_vendor_orders"
	EventPutUserOrder    = "put_user_order"
)

type UpdateBalanceEvent struct {
	TotalBalance int `json:"total_balance"`
}
