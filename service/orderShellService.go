package service

import (
	"context"
	"dvm.wallet/harsh/ent"
	"strconv"
)

type OrderShellOps struct {
	ctx    context.Context
	client *ent.Client
}

func NewOrderShellOps(ctx context.Context, client *ent.Client) *OrderShellOps {
	return &OrderShellOps{
		ctx:    ctx,
		client: client,
	}
}

func (r *OrderShellOps) CalculateTotalPrice(orderShell *ent.OrderShell) int {
	price := 0
	orders := orderShell.QueryOrders().AllX(r.ctx)
	OrderOps := NewOrderOps(r.ctx, r.client)
	for _, order := range orders {
		if order.Price == 0 {
			order.Update().SetPrice(OrderOps.CalculateTotalPrice(order)).SaveX(r.ctx)
		}
		price += order.Price
	}
	return price
}

func (r *OrderShellOps) ToDict(orderShell *ent.OrderShell) map[string]string {
	return map[string]string{
		"id":        strconv.Itoa(orderShell.ID),
		"timestamp": orderShell.Timestamp.String(),
		//"orders": //TODO: fix map/string incompatibility
	}
}
