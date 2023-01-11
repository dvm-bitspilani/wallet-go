package service

import (
	"context"
	"dvm.wallet/harsh/ent"
)

type VendorOps struct {
	ctx    context.Context
	client *ent.Client
}

func NewVendorOps(ctx context.Context, client *ent.Client) *VendorOps {
	return &VendorOps{
		ctx:    ctx,
		client: client,
	}
}

func (r *VendorOps) Open(vendor *ent.Vendor) {
	vendor.Update().SetClosed(false).SaveX(r.ctx)
}

func (r *VendorOps) Close(vendor *ent.Vendor) {
	vendor.Update().SetClosed(true).SaveX(r.ctx)
}
