package service

import (
	"context"
	"dvm.wallet/harsh/ent"
	"dvm.wallet/harsh/internal/database"
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

func GenerateAndPerform(amt int, kind helpers.Txn_type, srcUser *ent.User, dstUser *ent.User, ctx context.Context, client *ent.Client) (*ent.Transactions, error) {
	userOps := NewUserOps(ctx, client)

	src, err := userOps.GetOrCreateWallet(srcUser)
	if err != nil {
		return nil, err
	}
	dst, err := userOps.GetOrCreateWallet(dstUser)
	if err != nil {
		return nil, err
	}

	if src == dst {
		err := fmt.Errorf("reflexive transfers are not allowed")
		//err := exceptions.Exception{Message: "Reflexive transfers are not allowed", Status: 403}
		return nil, err
	}
	if srcUser.Disabled {
		err := fmt.Errorf("%s has been disabled", srcUser.Username)
		//err := exceptions.Exception{Message: fmt.Sprintf("%s has been disabled", src_user.Username), Status: 403}
		return nil, err
	}
	if dstUser.Disabled {
		err := fmt.Errorf("%s has been disabled", dstUser.Username)
		//err := exceptions.Exception{Message: fmt.Sprintf("%s has been disabled", src_user.Username), Status: 403}
		return nil, err
	}
	occupationPair := []string{srcUser.Occupation.String(), dstUser.Occupation.String()}
	validOccupationPair := false // It's probably a good idea to prevent any transaction than to allow *any* transaction
	for _, pair := range database.GetValidTransactionPairs() {
		if reflect.DeepEqual(occupationPair, pair) {
			validOccupationPair = true
			break
		}
	}
	if !validOccupationPair {
		if !(srcUser.Username == "SWD" && dstUser.Occupation == "bitsian") {
			err := fmt.Errorf("transaction forbidden: %s", occupationPair)
			//err := exceptions.Exception{Message: fmt.Sprintf("Transaction forbidden: %s", occupationPair), Status: 403}
			return nil, err
		}
	}
	walletOps := NewWalletOps(ctx, client)
	if !(srcUser.Occupation == "teller") {
		err = walletOps.Deduct(src, amt)
		if err != nil {
			return nil, err
		}
		err = walletOps.Add(dst, amt, database.GetBalanceFromTransactionType(kind))
		if err != nil {
			return nil, err
		}
	}
	return client.Transactions.Create().
		SetUser(dstUser).
		SetAmount(amt).
		SetKind(kind).
		SetSource(src).
		SetDestination(dst).
		SaveX(ctx), nil
}

func (r *TransactionOps) ToDict(txn *ent.Transactions) map[string]string {
	return map[string]string{
		"id":        strconv.Itoa(txn.ID),
		"amount":    strconv.Itoa(txn.Amount),
		"kind":      txn.Kind.String(),
		"timestamp": txn.Timestamp.String(),
	}
}
