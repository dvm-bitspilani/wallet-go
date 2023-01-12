package service

import (
	"context"
	"dvm.wallet/harsh/ent"
)

type ItemInstanceOps struct {
	ctx    context.Context
	client *ent.Client
}

type ItemInstanceStruct struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Quantity  int    `json:"quantity"`
	UnitPrice int    `json:"unit_price"`
}

func NewItemInstance(ctx context.Context, client *ent.Client) *ItemOps {
	return &ItemOps{
		ctx:    ctx,
		client: client,
	}
}

func (r *ItemOps) CalculateTotalPrice(itemInstance *ent.ItemInstance) int {
	return itemInstance.PricePerQuantity * itemInstance.Quantity
}
