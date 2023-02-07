package realtime

import (
	"dvm.wallet/harsh/ent"
	"dvm.wallet/harsh/pkg/websocket"
	"log"
)

func UpdateBalance(m *websocket.Manager, user *ent.User, balance int) {

	client := m.ClientUserIDList[user.ID]
	event := websocket.Event{}
	err := websocket.UpdateBalanceHandler(event, client, balance)
	if err != nil {
		log.Println(err)
		return
	}
}
