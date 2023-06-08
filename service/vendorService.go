package service

import (
	"context"
	"dvm.wallet/harsh/ent"
	"dvm.wallet/harsh/ent/order"
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

func (r *VendorOps) GetVendorArray(vendor *ent.VendorSchema) []int {
	var vendorIdArray []int
	for _, ItemObj := range vendor.QueryItems().AllX(r.ctx) {
		vendorIdArray = append(vendorIdArray, ItemObj.ID)
	}
	return vendorIdArray
}

func CalculateEarnings(vendor *ent.VendorSchema) int {
	ctx := context.Background()
	//total := 0
	//for _, order := range vendor.QueryOrders().Where(order.StatusIn(2, 3)).AllX(ctx) {
	//	total += order.Price
	//}
	return vendor.QueryOrders().Where(order.StatusIn(2, 3)).Aggregate(ent.Sum("price")).IntX(ctx)
}
