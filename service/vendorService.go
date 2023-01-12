package service

import (
	"context"
	"dvm.wallet/harsh/ent"
	"net/url"
)

type VendorOps struct {
	ctx    context.Context
	client *ent.Client
}

type VendorStruct struct {
	Id       int     `json:"id"`
	Name     string  `json:"name"`
	ImageUrl url.URL `json:"image_url"`
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
