package service

import (
	"context"
	"dvm.wallet/harsh/ent"
	"dvm.wallet/harsh/ent/item"
	vendor "dvm.wallet/harsh/ent/vendorschema"
	"dvm.wallet/harsh/internal/database"
	"dvm.wallet/harsh/internal/helpers"
	"dvm.wallet/harsh/internal/validator"
	"errors"
	"fmt"
	"reflect"
)

type UserOps struct {
	ctx    context.Context
	client *ent.Client
}

func NewUserOps(ctx context.Context, client *ent.Client) *UserOps {
	return &UserOps{
		ctx:    ctx,
		client: client,
	}
}

func (r *UserOps) Disable(user *ent.User) error {
	_, err := user.Update().SetDisabled(true).Save(r.ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserOps) Enable(user *ent.User) error {
	_, err := user.Update().SetDisabled(false).Save(r.ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserOps) GetOrCreateWallet(user *ent.User) (*ent.Wallet, error) {
	if user.Disabled {
		err := fmt.Errorf("cannot create or return wallet of as %s is disabled", user.Username)
		//err := exceptions.UserDisabledException{Exception: exceptions.Exception{Message: fmt.Sprintf("cannot create or return wallet of as %s is disabled", user.Username)}}
		return nil, err //403
	}
	wallet, err := user.QueryWallet().Only(r.ctx)
	if err != nil {
		wallet = r.client.Wallet.Create().SetUser(user).SaveX(r.ctx)
	}
	return wallet, nil
}

func (r *UserOps) Transfer(user *ent.User, target *ent.User, amount int) (*ent.Transactions, error, int) {
	if user.Disabled {
		err := fmt.Errorf("requesting User is disabled")
		//err := exceptions.UserDisabledException{Exception: exceptions.Exception{Message: "requesting User is disabled"}}
		return nil, err, 403 //403
	}
	if target.Disabled {
		err := fmt.Errorf("target User is disabled")
		//err := exceptions.UserDisabledException{Exception: exceptions.Exception{Message: "target User is disabled"}}
		return nil, err, 403 //403
	}
	if amount <= 0 {
		err := fmt.Errorf("amount cannot be negative")
		//err := exceptions.Exception{Message: "amount cannot be negative", Status: 400}
		return nil, err, 400 //400
	}
	occupationPair := []string{user.Occupation.String(), target.Occupation.String()}
	validOccupationPair := false // It's probably a good idea to prevent any transaction than to allow *any* transaction
	for _, pair := range database.GetValidTransactionPairs() {
		if reflect.DeepEqual(occupationPair, pair) {
			validOccupationPair = true
			break
		}
	}
	if validOccupationPair {
		transaction, err, statusCode := GenerateAndPerform(amount, helpers.TRANSFER, user, target, r.ctx, r.client)
		if err != nil {
			return nil, err, statusCode
		}
		return transaction, nil, 0
	} else {
		err := fmt.Errorf("transaction forbidden: %s", occupationPair)
		//err := exceptions.Exception{Message: fmt.Sprintf("transaction forbidden: %s", occupationPair), Status: 403}
		return nil, err, 403 //403
	}
}

func (r *UserOps) PlaceOrder(usr *ent.User, orderList []helpers.OrderActionVendorStruct) (*OrderShellStruct, error, int) {
	var statusCode int

	if !validator.In(usr.Occupation, "bitsian", "participant") {
		return nil, errors.New("only bitsians and participants may place orders"), 403
	}
	var totalPrice int
	if usr.Disabled {
		return nil, errors.New("cannot place order, user's account has been disabled"), 403
	}
	// TODO:	refactor this so we're not running two for loops
	for _, vendorStruct := range orderList {
		vendorObj, err := r.client.VendorSchema.Query().Where(vendor.ID(vendorStruct.VendorId)).Only(r.ctx)
		if err != nil {
			return nil, fmt.Errorf("vendor with ID %d does not exist", vendorStruct.VendorId), 403 //404
		}
		if vendorObj.Closed {
			return nil, fmt.Errorf("Vendor %s is closed", vendorObj.Name), 412 //412 ??
		}
		if vendorObj.QueryUser().OnlyX(r.ctx).Disabled {
			return nil, fmt.Errorf("Vendor %s is disabled", vendorObj.Name), 412 // 403
		}

		for _, itemStruct := range vendorStruct.Order {
			itemObj, err := r.client.Item.Query().Where(item.ID(itemStruct.ItemId)).Only(r.ctx)
			if err != nil {
				return nil, fmt.Errorf("item with ID %d does not exist", itemStruct.ItemId), 404
			}
			if itemObj.QueryVendorSchema().OnlyX(r.ctx).ID != vendorObj.ID {
				//fmt.Println(itemObj.QueryVendorSchema().OnlyX(r.ctx) == vendorObj)      // TODO: This returns false, find out why
				return nil, errors.New("cannot order items from the wrong vendor"), 403 // 403
			}
			if !itemObj.Available {
				return nil, fmt.Errorf("%s item is currently unavailable", itemObj.Name), 412 // 412
			}

			if itemStruct.Quantity <= 0 {
				return nil, errors.New("cannot order a negative or 0 quantity of items"), 400 // 400
			}
			totalPrice += itemObj.BasePrice * itemStruct.Quantity
		}
	}
	walletOps := NewWalletOps(r.ctx, r.client)
	if totalPrice > walletOps.Balance(usr.QueryWallet().OnlyX(r.ctx)) {
		return nil, fmt.Errorf("order price: %d, current balance: %d", totalPrice, walletOps.Balance(usr.QueryWallet().OnlyX(r.ctx))), 412
	}

	// creation and saving phase
	shell := r.client.OrderShell.Create().SetWallet(usr.QueryWallet().OnlyX(r.ctx)).SetPrice(totalPrice).SaveX(r.ctx)
	orderOps := NewOrderOps(r.ctx, r.client)
	for _, vendorStruct := range orderList {
		vendorObj := r.client.VendorSchema.Query().Where(vendor.ID(vendorStruct.VendorId)).OnlyX(r.ctx)
		order := r.client.Order.Create().
			SetShell(shell).
			SetVendorSchema(vendorObj).
			SetStatus(helpers.PENDING).
			SaveX(r.ctx)

		for _, itemStruct := range vendorStruct.Order {
			itemObj := r.client.Item.Query().Where(item.ID(itemStruct.ItemId)).OnlyX(r.ctx)
			r.client.ItemInstance.Create().
				SetItem(itemObj).
				SetQuantity(itemStruct.Quantity).
				SetOrder(order).
				SetPricePerQuantity(itemObj.BasePrice).
				SaveX(r.ctx)
			order.Update().SetPrice(orderOps.CalculateTotalPrice(order)).SaveX(r.ctx)
		}
	}
	err, statusCode := walletOps.Deduct(usr.QueryWallet().OnlyX(r.ctx), totalPrice)
	if err != nil {
		return nil, err, statusCode
	}
	//TODO:		put_orders
	OrderShellOps := NewOrderShellOps(r.ctx, r.client)
	return OrderShellOps.ToDict(shell), nil, 0
}
