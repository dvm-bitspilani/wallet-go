package service

import (
	"context"
	"dvm.wallet/harsh/ent"
	"dvm.wallet/harsh/internal/helpers"
	"fmt"
	"reflect"
	"strconv"
)

type TransactionOps struct {
	ctx    context.Context
	client *ent.Client
}

func NewTransactionOps(ctx context.Context, app *ent.Client) *TransactionOps {
	return &TransactionOps{
		ctx:    ctx,
		client: app,
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
		err := fmt.Errorf("reflexive transfers are not allowed")
		//err := exceptions.Exception{Message: "Reflexive transfers are not allowed", Status: 403}
		return nil, err
	}
	if src_user.Disabled {
		err := fmt.Errorf("%s has been disabled", src_user.Username)
		//err := exceptions.Exception{Message: fmt.Sprintf("%s has been disabled", src_user.Username), Status: 403}
		return nil, err
	}
	if dst_user.Disabled {
		err := fmt.Errorf("%s has been disabled", dst_user.Username)
		//err := exceptions.Exception{Message: fmt.Sprintf("%s has been disabled", src_user.Username), Status: 403}
		return nil, err
	}
	occupationPair := []string{src_user.Occupation.String(), dst_user.Occupation.String()}
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
			//err := exceptions.Exception{Message: fmt.Sprintf("Transaction forbidden: %s", occupationPair), Status: 403}
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

func (r *TransactionOps) To_dict(txn *ent.Transactions) map[string]string {
	return map[string]string{
		"id":        strconv.Itoa(txn.ID),
		"amount":    strconv.Itoa(txn.Amount),
		"kind":      txn.Kind.String(),
		"timestamp": txn.Timestamp.String(),
	}
}
