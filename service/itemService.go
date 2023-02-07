package service

import (
	"context"
	"dvm.wallet/harsh/cmd/api/config"
	"dvm.wallet/harsh/ent"
)

type ItemOps struct {
	ctx context.Context
	app *config.Application
}

func NewItemOps(ctx context.Context, app *config.Application) *ItemOps {
	return &ItemOps{
		ctx: ctx,
		app: app,
	}
}

func (r *ItemOps) MakeAvailable(item *ent.Item) {
	item.Update().SetAvailable(true).SaveX(r.ctx)
}

func (r *ItemOps) MakeUnavailable(item *ent.Item) {
	item.Update().SetAvailable(false).SaveX(r.ctx)
}
