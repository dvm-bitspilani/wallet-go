package service

import (
	"context"
	"dvm.wallet/harsh/cmd/api/config"
	"dvm.wallet/harsh/ent"
	"dvm.wallet/harsh/internal/helpers"
	"fmt"
	"reflect"
)

type TransactionOps struct {
	ctx    context.Context
	client *ent.Client
}

func NewTransactionOps(ctx context.Context, app *config.Application) *TransactionOps {
	return &TransactionOps{
		ctx:    ctx,
		client: app.Client,
	}
}

func GenerateAndPerform(amt int, kind helpers.Txn_type, src_user *ent.User, dst_user *ent.User, ctx context.Context, client *ent.Client) (*ent.Transactions, error) {
	userops := NewUserOps(ctx, client)

	src, err := userops.GetOrCreateWallet(src_user)
	if err != nil {
		return nil, err
	}
	dst, err := userops.GetOrCreateWallet(dst_user)
	if err != nil {
		return nil, err
	}

	if src == dst {
		err := fmt.Errorf("Reflexive transfers are not allowed")
		return nil, err
	}
	if src_user.Disabled {
		err := fmt.Errorf("%s has been disabled", src_user.Username)
		return nil, err
	}
	if dst_user.Disabled {
		err := fmt.Errorf("%s has been disabled", dst_user.Username)
		return nil, err
	}
	// TODO:
	occupationPair := []string{src_user.Occupation.String(), dst_user.Occupation.String()}
	//if !(validator.In(occupation_pair, helpers.GetValidTransactionPairs()...)) {
	//
	//}
	validOccupationPair := false // It's probably a good idea to prevent any transaction than to allow *any* transaction
	for _, pair := range helpers.GetValidTransactionPairs() {
		if reflect.DeepEqual(occupationPair, pair) {
			validOccupationPair = true
			break
		}
	}
	if !validOccupationPair {
		if !(src_user.Username == "SWD" && dst_user.Occupation == "bitsian") {
			err := fmt.Errorf("Transaction forbidden: %s", occupationPair)
			return nil, err
		}
	}

	if !(src_user.Occupation == "teller") {
		// src.deduct()
		// dst.add()
	}
	return client.Transactions.Create().
		SetUser(dst_user).
		SetAmount(amt).
		SetKind(kind).
		SetSource(src).
		SetDestination(dst).
		SaveX(ctx), nil
}
