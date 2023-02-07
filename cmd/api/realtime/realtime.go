package realtime

import (
	"dvm.wallet/harsh/internal/helpers"
	"dvm.wallet/harsh/pkg/websocket"
	"log"
)

func UpdateBalance(m *websocket.Manager, userId int, balance int) {

	client := m.ClientUserIDList[userId]
	event := websocket.Event{}
	err := websocket.UpdateBalanceHandler(event, client, balance)
	if err != nil {
		log.Println(err)
		return
	}
}

func UpdateOrderStatus(m *websocket.Manager, userId int, status helpers.Status) {
	client := m.ClientUserIDList[userId]
	event := websocket.Event{}
	log.Println(m.ClientUserIDList)
	err := websocket.UpdateOrderStatusHandler(event, client, status)
	if err != nil {
		return
	}
}
