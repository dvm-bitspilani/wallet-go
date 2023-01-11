package service

import (
	"context"
	"dvm.wallet/harsh/ent"
)

type ItemOps struct {
	ctx    context.Context
	client *ent.Client
}

func NewItemOps(ctx context.Context, client *ent.Client) *ItemOps {
	return &ItemOps{
		ctx:    ctx,
		client: client,
	}
}

func (r *ItemOps) MakeAvailable(item *ent.Item) {
	item.Update().SetAvailable(true).SaveX(r.ctx)
}

func (r *ItemOps) MakeUnavailable(item *ent.Item) {
	item.Update().SetAvailable(false).SaveX(r.ctx)
}
