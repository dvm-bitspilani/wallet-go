package service

import (
	"context"
	"dvm.wallet/harsh/ent"
)

type VendorOps struct {
	ctx    context.Context
	client *ent.Client
}

type VendorStruct struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	ImageUrl string `json:"image_url"`
}

func NewVendorOps(ctx context.Context, client *ent.Client) *VendorOps {
	return &VendorOps{
		ctx:    ctx,
		client: client,
	}
}

func (r *VendorOps) Open(vendor *ent.VendorSchema) {
	vendor.Update().SetClosed(false).SaveX(r.ctx)
}

func (r *VendorOps) Close(vendor *ent.VendorSchema) {
	vendor.Update().SetClosed(true).SaveX(r.ctx)
}
