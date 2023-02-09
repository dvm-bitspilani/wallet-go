package service

import (
	"context"
	"dvm.wallet/harsh/cmd/api/config"
	"dvm.wallet/harsh/ent"
	"time"
)

type OrderShellOps struct {
	ctx context.Context
	app *config.Application
}

type OrderShellStruct struct {
	Id        int           `json:"id"`
	Timestamp time.Time     `json:"timestamp"`
	Orders    []OrderStruct `json:"orders"`
}

func NewOrderShellOps(ctx context.Context, app *config.Application) *OrderShellOps {
	return &OrderShellOps{
		ctx: ctx,
		app: app,
	}
}

func (r *OrderShellOps) CalculateTotalPrice(orderShell *ent.OrderShell) int {
	price := 0
	orders := orderShell.QueryOrders().AllX(r.ctx)
	OrderOps := NewOrderOps(r.ctx, r.app)
	for _, order := range orders {
		if order.Price == 0 {
			order.Update().SetPrice(OrderOps.CalculateTotalPrice(order)).SaveX(r.ctx)
		}
		price += order.Price
	}
	return price
}

func (r *OrderShellOps) ToDict(orderShell *ent.OrderShell) *OrderShellStruct {
	orderOps := NewOrderOps(r.ctx, r.app)
	var orders []OrderStruct
	for _, order := range orderShell.QueryOrders().AllX(r.ctx) {
		orders = append(orders, orderOps.ToDict(order))
	}
	return &OrderShellStruct{
		Id:        orderShell.ID,
		Timestamp: orderShell.Timestamp,
		Orders:    orders,
	}
}
