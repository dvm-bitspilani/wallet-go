package websocket

import (
	"encoding/json"
	"fmt"
	"golang.org/x/exp/rand"
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

func UpdateBalanceHandler(e Event, c *Client) error {
	//var BalanceEvent UpdateBalanceEvent
	//if err := json.Unmarshal(e.Payload, &BalanceEvent); err != nil {
	//	return fmt.Errorf("bad payload in request: %v", err)
	//}
	//fmt.Println(BalanceEvent.TotalBalance)
	//return nil
	var updateBalanceEvent UpdateBalanceEvent
	updateBalanceEvent.TotalBalance = rand.Int()

	data, err := json.Marshal(updateBalanceEvent)
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}

	var outgoingEvent Event
	outgoingEvent.Type = EventUpdateBalance
	outgoingEvent.Payload = data

	for client := range c.manager.clients {
		client.egress <- outgoingEvent
	}

	return nil
}

type UpdateOrderStatus struct {
	Status string `json:"status"`
}

func UpdateOrderStatusHandler(e Event, c *Client) error {
	//var BalanceEvent UpdateBalanceEvent
	//if err := json.Unmarshal(e.Payload, &BalanceEvent); err != nil {
	//	return fmt.Errorf("bad payload in request: %v", err)
	//}
	//fmt.Println(BalanceEvent.TotalBalance)
	//return nil
	var updateOrderStatusEvent UpdateOrderStatus
	updateOrderStatusEvent.Status = "Bruh"

	data, err := json.Marshal(updateOrderStatusEvent)
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}

	var outgoingEvent Event
	outgoingEvent.Type = EventUpdateOrder
	outgoingEvent.Payload = data

	for client := range c.manager.clients {
		client.egress <- outgoingEvent
	}

	return nil
}
