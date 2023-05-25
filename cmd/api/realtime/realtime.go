package realtime

import (
	"context"
	"dvm.wallet/harsh/ent"
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

func UpdateOrderStatus(m *websocket.Manager, userId int, orderId int, status helpers.Status) {
	client := m.ClientUserIDList[userId]
	event := websocket.Event{}
	err := websocket.UpdateOrderStatusHandler(event, client, orderId, status)
	if err != nil {
		return
	}
}

func PutUserOrders(m *websocket.Manager, userObj *ent.User) {
	ctx := context.Background()
	orderArray := userObj.QueryWallet().QueryShells().QueryOrders().AllX(ctx)
	orderIdArray := make([]int, len(orderArray))
	for _, orderObj := range orderArray {
		orderIdArray = append(orderIdArray, orderObj.ID)
	}
	client := m.ClientUserIDList[userObj.ID]
	event := websocket.Event{}
	err := websocket.PutUserOrderHandler(event, client, orderIdArray)
	if err != nil {
		return
	}
}

func PutVendorOrders(m *websocket.Manager, userId int, orderIdArray []int) {
	client := m.ClientUserIDList[userId]
	event := websocket.Event{}
	err := websocket.PutVendorOrdersHandler(event, client, orderIdArray)
	if err != nil {
		return
	}
}
