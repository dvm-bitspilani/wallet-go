package service

import (
	"context"
	"dvm.wallet/harsh/ent"
	"dvm.wallet/harsh/internal/database"
	"dvm.wallet/harsh/internal/helpers"
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
	if !user.Disabled {
		err := fmt.Errorf("cannot create or return wallet of as %s is disabled", user.Username)
		//err := exceptions.UserDisabledException{Exception: exceptions.Exception{Message: fmt.Sprintf("cannot create or return wallet of as %s is disabled", user.Username)}}
		return nil, err
	}
	wallet, err := user.QueryWallet().Only(r.ctx)
	if err != nil {
		wallet = r.client.Wallet.Create().SetUser(user).SaveX(r.ctx)
	}
	return wallet, nil
}

func (r *UserOps) Transfer(user *ent.User, target *ent.User, amount int) (*ent.Transactions, error) {
	if user.Disabled {
		err := fmt.Errorf("requesting User is disabled")
		//err := exceptions.UserDisabledException{Exception: exceptions.Exception{Message: "requesting User is disabled"}}
		return nil, err
	}
	if target.Disabled {
		err := fmt.Errorf("target User is disabled")
		//err := exceptions.UserDisabledException{Exception: exceptions.Exception{Message: "target User is disabled"}}
		return nil, err
	}
	if amount <= 0 {
		err := fmt.Errorf("amount cannot be negative")
		//err := exceptions.Exception{Message: "amount cannot be negative", Status: 400}
		return nil, err
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
		transaction, err := GenerateAndPerform(amount, helpers.TRANSFER, user, target, r.ctx, r.client)
		if err != nil {
			return nil, err
		}
		return transaction, nil
	} else {
		err := fmt.Errorf("transaction forbidden: %s", occupationPair)
		//err := exceptions.Exception{Message: fmt.Sprintf("transaction forbidden: %s", occupationPair), Status: 403}
		return nil, err
	}
}
