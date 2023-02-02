package service

import (
	"context"
	"dvm.wallet/harsh/ent"
	"dvm.wallet/harsh/internal/database"
	"dvm.wallet/harsh/internal/helpers"
	"dvm.wallet/harsh/internal/validator"
	"errors"
	"fmt"
	"time"
)

type OrderOps struct {
	ctx    context.Context
	client *ent.Client
}

type OrderStruct struct {
	OrderId     int                  `json:"order_id"`
	Shell       int                  `json:"shell"`
	Vendor      VendorStruct         `json:"VendorSchema"`
	Items       []ItemInstanceStruct `json:"items"`
	Transaction int                  `json:"transaction"`
	Price       int                  `json:"price"`
	Status      helpers.Status       `json:"status"`
	Otp         string               `json:"otp,"`
	OtpSeen     bool                 `json:"otp_seen"`
}

func NewOrderOps(ctx context.Context, client *ent.Client) *OrderOps {
	return &OrderOps{
		ctx:    ctx,
		client: client,
	}
}

func (r *OrderOps) ChangeStatus(order *ent.Order, newStatus helpers.Status, usr *ent.User) (int, error, int) {
	if validator.In(order.Status, helpers.DECLINED, helpers.FINISHED) {
		err := fmt.Errorf("maximum/final status achieved")
		return 0, err, 412
	}
	if order.Status == helpers.READY {
		if !order.OtpSeen {
			err := fmt.Errorf("user has not yet hit see otp")
			return 0, err, 412
		}
	}
	order.Update().SetStatus(newStatus)
	walletOps := NewWalletOps(r.ctx, r.client)
	if order.Status == helpers.READY {
		transaction := r.client.Transactions.Create().
			SetUser(usr).
			SetAmount(order.Price).
			SetKind(helpers.PURCHASE).
			SetSource(usr.QueryWallet().OnlyX(r.ctx)).
			SetDestination(usr.QueryVendorSchema().QueryUser().QueryWallet().OnlyX(r.ctx)).
			SaveX(r.ctx)
		order.Update().SetTransaction(transaction).SetTimestamp(time.Now()).SaveX(r.ctx)
		err, statusCode := walletOps.Add(usr.QueryWallet().OnlyX(r.ctx), order.Price, database.TRANSFER_BAL)
		if err != nil {
			return 0, err, statusCode
		}
	} else {
		if order.Status == helpers.DECLINED {
			order.Update().SetDeclinedTimestamp(time.Now()).SaveX(r.ctx)
		} else if order.Status == helpers.FINISHED {
			order.Update().SetAcceptedTimestamp(time.Now()).SaveX(r.ctx)
		} else if order.Status == helpers.ACCEPTED {
			order.Update().SetAcceptedTimestamp(time.Now())
		}
	}
	// TODO:	update_order_status
	return int(order.Status), nil, 0 // not sure if this direct conversion works
}

func (r *OrderOps) Decline(order *ent.Order) (error, int) {
	if order.Status == helpers.DECLINED {
		return errors.New("vendor has already declined the order, cannot re-decline an order"), 412
	}
	if validator.In(order.Status, helpers.ACCEPTED, helpers.READY, helpers.FINISHED) {
		return errors.New("vendor has already accepted the order, cannot decline now"), 412
	}
	order.Update().SetStatus(helpers.DECLINED).SaveX(r.ctx)
	// TODO:	update_order_status
	return nil, 0
}

func (r *OrderOps) CalculateTotalPrice(order *ent.Order) int {
	price := 0
	items := order.QueryIteminstances().AllX(r.ctx)
	ItemOps := NewItemOps(r.ctx, r.client)
	for _, item := range items {
		price += ItemOps.CalculateTotalPrice(item)
	}
	return price
}

func (r *OrderOps) ToDict(order *ent.Order) OrderStruct {
	var items []ItemInstanceStruct
	for _, item := range order.QueryIteminstances().AllX(r.ctx) {
		items = append(items, ItemInstanceStruct{
			Id:        item.ID,
			Name:      item.QueryItem().OnlyX(r.ctx).Name,
			Quantity:  item.Quantity,
			UnitPrice: item.PricePerQuantity,
		})
	}
	orderVendor := order.QueryVendorSchema().OnlyX(r.ctx)
	vendor := VendorStruct{
		Id:       orderVendor.ID,
		Name:     orderVendor.Name,
		ImageUrl: orderVendor.ImageURL,
	}

	return OrderStruct{
		OrderId:     order.ID,
		Shell:       order.QueryShell().OnlyX(r.ctx).ID,
		Vendor:      vendor,
		Items:       items,
		Transaction: order.QueryTransaction().OnlyX(r.ctx).ID,
		Price:       order.Price,
		Status:      order.Status,
		Otp:         order.Otp,
		OtpSeen:     order.OtpSeen,
	}
}

func (r *OrderOps) MakeOtpSeen(order *ent.Order) {
	order.Update().SetOtpSeen(true).SaveX(r.ctx)
	// TODO:	update_otp_seen
}
