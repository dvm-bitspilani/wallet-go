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
	EventUpdateBalance = "update_balance"
	EventUpdateOrder   = "update_order_status"
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

	for client := range c.manager.Clients {
		client.egress <- outgoingEvent
	}

	return nil
}
