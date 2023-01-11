package service

import (
	"context"
	"dvm.wallet/harsh/ent"
)

type ItemInstanceOps struct {
	ctx    context.Context
	client *ent.Client
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
