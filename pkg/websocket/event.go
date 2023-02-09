package websocket

import (
	"dvm.wallet/harsh/internal/helpers"
	"encoding/json"
	"fmt"
)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type EventHandler func(event Event, c *Client) error

const (
	EventUpdateBalance   = "update_balance"
	EventUpdateOrder     = "update_order_status"
	EventPutVendorOrders = "put_vendor_orders"
	EventPutUserOrder    = "put_user_order"
)

type UpdateBalanceEvent struct {
	TotalBalance int `json:"total_balance"`
}

func UpdateBalanceHandler(e Event, c *Client, amount int) error {
	//var BalanceEvent UpdateBalanceEvent
	//if err := json.Unmarshal(e.Payload, &BalanceEvent); err != nil {
	//	return fmt.Errorf("bad payload in request: %v", err)
	//}
	//fmt.Println(BalanceEvent.TotalBalance)
	//return nil
	var updateBalanceEvent UpdateBalanceEvent
	updateBalanceEvent.TotalBalance = amount

	data, err := json.Marshal(updateBalanceEvent)
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}

	var outgoingEvent Event
	outgoingEvent.Type = EventUpdateBalance
	outgoingEvent.Payload = data

	c.egress <- outgoingEvent
	//for client := range c.manager.Clients {
	//	client.egress <- outgoingEvent
	//}

	return nil
}

type UpdateOrderStatus struct {
	OrderId int    `json:"order_id"`
	Status  string `json:"status"`
}

func UpdateOrderStatusHandler(e Event, c *Client, orderId int, status helpers.Status) error {
	//var BalanceEvent UpdateBalanceEvent
	//if err := json.Unmarshal(e.Payload, &BalanceEvent); err != nil {
	//	return fmt.Errorf("bad payload in request: %v", err)
	//}
	//fmt.Println(BalanceEvent.TotalBalance)
	//return nil
	var updateOrderStatusEvent UpdateOrderStatus
	updateOrderStatusEvent.OrderId = orderId
	updateOrderStatusEvent.Status = status.String()

	data, err := json.Marshal(updateOrderStatusEvent)
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}

	var outgoingEvent Event
	outgoingEvent.Type = EventUpdateOrder
	outgoingEvent.Payload = data

	c.egress <- outgoingEvent

	return nil
}

type PutVendorOrders struct {
	OrderIdArray []int `json:"order_ids"`
}

func PutVendorOrdersHandler(e Event, c *Client, orderIdArray []int) error {
	var putOrdersEvent PutVendorOrders
	putOrdersEvent.OrderIdArray = orderIdArray

	data, err := json.Marshal(putOrdersEvent)
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}

	var outgoingEvent Event
	outgoingEvent.Type = EventPutVendorOrders
	outgoingEvent.Payload = data

	c.egress <- outgoingEvent

	return nil
}

type PutUserOrder struct {
	UserOrderIdArray []int `json:"user_order_id_array"`
}

func PutUserOrderHandler(e Event, c *Client, orderIdArray []int) error {
	var putUserOrderEvent PutUserOrder
	putUserOrderEvent.UserOrderIdArray = orderIdArray

	data, err := json.Marshal(putUserOrderEvent)
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}

	var outgoingEvent Event
	outgoingEvent.Type = EventPutUserOrder
	outgoingEvent.Payload = data

	c.egress <- outgoingEvent

	return nil
}
