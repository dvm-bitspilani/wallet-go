package service

import (
	"context"
	"dvm.wallet/harsh/cmd/api/config"
	"dvm.wallet/harsh/ent"
)

type ItemInstanceOps struct {
	ctx context.Context
	app *config.Application
}

type ItemInstanceStruct struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Quantity  int    `json:"quantity"`
	UnitPrice int    `json:"unit_price"`
}

func NewItemInstance(ctx context.Context, app *config.Application) *ItemInstanceOps {
	return &ItemInstanceOps{
		ctx: ctx,
		app: app,
	}
}

func (r *ItemOps) CalculateTotalPrice(itemInstance *ent.ItemInstance) int {
	return itemInstance.PricePerQuantity * itemInstance.Quantity
}
