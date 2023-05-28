package realtime

import (
	"context"
	"dvm.wallet/harsh/cmd/api/config"
	"dvm.wallet/harsh/ent"
	"dvm.wallet/harsh/internal/helpers"
	"dvm.wallet/harsh/pkg/sse"
	"dvm.wallet/harsh/pkg/websocket"
	"encoding/json"
	"fmt"
)

//func UpdateBalance(m *websocket.Manager, userId int, balance int) {
//	client := m.ClientUserIDList[userId]
//	event := websocket.Event{}
//	err := websocket.UpdateBalanceHandler(event, client, balance)
//	if err != nil {
//		log.Println(err)
//		return
//	}
//}

func UpdateBalance(app *config.Application, ctx *context.Context, userId int, balance int) {
	var updateBalanceEvent sse.UpdateBalanceEvent
	updateBalanceEvent.TotalBalance = balance

	data, err := json.Marshal(updateBalanceEvent)
	if err != nil {
		//return fmt.Errorf("failed to marshal broadcast message: %v", err)
		app.Logger.Errorf("failed to marshal data. Error: %v", err)
		return
	}
	var outgoingEvent sse.Event
	outgoingEvent.Type = sse.EventUpdateBalance
	outgoingEvent.Payload = data
	outgoingEventData, err := json.Marshal(&outgoingEvent)

	userChannelName := fmt.Sprintf("user#%d", userId)
	err = app.Rdb.Publish(*ctx, userChannelName, outgoingEventData).Err()
	if err != nil {
		app.Logger.Errorf("failed to publish message to Redis. Error: %v", err)
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
