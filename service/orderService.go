package service

import (
	"context"
	"dvm.wallet/harsh/ent"
	"dvm.wallet/harsh/internal/helpers"
	"dvm.wallet/harsh/internal/validator"
	"errors"
	"fmt"
	"strconv"
	"time"
)

type OrderOps struct {
	ctx    context.Context
	client *ent.Client
}

func NewOrderOps(ctx context.Context, client *ent.Client) *OrderOps {
	return &OrderOps{
		ctx:    ctx,
		client: client,
	}
}

func (r *OrderOps) ChangeStatus(order *ent.Order, newStatus helpers.Status, usr *ent.User) (int, error) {
	if validator.In(order.Status, helpers.DECLINED, helpers.FINISHED) {
		err := fmt.Errorf("maximum/final status achieved")
		return 0, err
	}
	if order.Status == helpers.READY {
		if !order.OtpSeen {
			err := fmt.Errorf("user has not yet hit see otp")
			return 0, err
		}
	}
	order.Update().SetStatus(newStatus)
	walletOps := NewWalletOps(r.ctx, r.client)
	if order.Status == helpers.READY {
		transaction := r.client.Transactions.Create().
			SetUser(usr).
			SetAmount(order.Price).
			SetKind(helpers.PURCHASE).
			SetSource(usr.Edges.Wallet).
			SetDestination(usr.Edges.Vendor.Edges.User.Edges.Wallet).
			SaveX(r.ctx)
		order.Update().SetTransaction(transaction).SetTimestamp(time.Now()).SaveX(r.ctx)
		_ = walletOps.Add(usr.Edges.Wallet, order.Price, helpers.TRANSFER_BAL)
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
	return int(order.Status), nil // not sure if this direct conversion works
}

func (r *VendorOps) Decline(order *ent.Order) error {
	if order.Status == helpers.DECLINED {
		return errors.New("vendor has already declined the order, cannot re-decline an order")
	}
	if validator.In(order.Status, helpers.ACCEPTED, helpers.READY, helpers.FINISHED) {
		return errors.New("vendor has already accepted the order, cannot decline now")
	}
	order.Update().SetStatus(helpers.DECLINED).SaveX(r.ctx)
	// TODO:	update_order_status
	return nil
}

func (r *VendorOps) CalculateTotalPrice(order *ent.Order) int {
	price := 0
	items := order.QueryIteminstances().AllX(r.ctx)
	ItemOps := NewItemOps(r.ctx, r.client)
	for _, item := range items {
		price += ItemOps.CalculateTotalPrice(item)
	}
	return price
}

func (r *OrderOps) ToDict(order *ent.Order) map[string]string {
	//vendor := map[string]string{
	//	"id":        strconv.Itoa(order.Edges.Vendor.ID),
	//	"name":      order.Edges.Vendor.Name,
	//	"image_url": order.Edges.Vendor.ImageURL.String(),
	//}
	//items := make([]map[string]interface{}, len())
	return map[string]string{
		"order_id": strconv.Itoa(order.ID),
		"shell":    strconv.Itoa(order.Edges.Shell.ID),
		//"vendor":      vendor,
		//"items":       items,
		"transaction": strconv.Itoa(order.Edges.Transaction.ID),
		"price":       strconv.Itoa(order.Price),
		"status":      order.Status.String(),
		"otp":         strconv.Itoa(order.Otp), //TODO:	change otp to string
		"otp_seen":    strconv.FormatBool(order.OtpSeen),
	}
}

func (r *OrderOps) MakeOtpSeen(order *ent.Order) {
	order.Update().SetOtpSeen(true).SaveX(r.ctx)
	// TODO:	update_otp_seen
}
