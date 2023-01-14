package service

import (
	"context"
	"dvm.wallet/harsh/ent"
	"dvm.wallet/harsh/internal/database"
	"dvm.wallet/harsh/internal/helpers"
	"fmt"
	"reflect"
	"time"
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

type TransactionStruct struct {
	Id        int              `json:"id"`
	Amount    int              `json:"amount"`
	Kind      helpers.Txn_type `json:"kind"`
	Timestamp time.Time        `json:"timestamp"`
}

func GenerateAndPerform(amt int, kind helpers.Txn_type, srcUser *ent.User, dstUser *ent.User, ctx context.Context, client *ent.Client) (*ent.Transactions, error, int) {
	var statusCode int

	userOps := NewUserOps(ctx, client)
	src, err := userOps.GetOrCreateWallet(srcUser)
	if err != nil {
		return nil, err, 403 // 403
	}
	dst, err := userOps.GetOrCreateWallet(dstUser)
	if err != nil {
		return nil, err, 403 // 403
	}

	if src == dst {
		err := fmt.Errorf("reflexive transfers are not allowed")
		//err := exceptions.Exception{Message: "Reflexive transfers are not allowed", Status: 403}
		return nil, err, 403 // 403
	}
	if srcUser.Disabled {
		err := fmt.Errorf("%s has been disabled", srcUser.Username)
		//err := exceptions.Exception{Message: fmt.Sprintf("%s has been disabled", src_user.Username), Status: 403}
		return nil, err, 403 // 403
	}
	if dstUser.Disabled {
		err := fmt.Errorf("%s has been disabled", dstUser.Username)
		//err := exceptions.Exception{Message: fmt.Sprintf("%s has been disabled", src_user.Username), Status: 403}
		return nil, err, 403 // 403
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
			return nil, err, 403 // 403
		}
	}
	walletOps := NewWalletOps(ctx, client)
	if !(srcUser.Occupation == "teller") {
		err, statusCode = walletOps.Deduct(src, amt)
		if err != nil {
			return nil, err, statusCode
		} // 400, 412
		err, statusCode = walletOps.Add(dst, amt, database.GetBalanceFromTransactionType(kind))
		if err != nil {
			return nil, err, statusCode
		} // 400
	}
	return client.Transactions.Create().
		SetUser(dstUser).
		SetAmount(amt).
		SetKind(kind).
		SetSource(src).
		SetDestination(dst).
		SaveX(ctx), nil, 0
}

func (r *TransactionOps) ToDict(txn *ent.Transactions) *TransactionStruct {
	//return map[string]string{
	//	"id":        strconv.Itoa(txn.ID),
	//	"amount":    strconv.Itoa(txn.Amount),
	//	"kind":      txn.Kind.String(),
	//	"timestamp": txn.Timestamp.String(),
	//}
	return &TransactionStruct{
		Id:        txn.ID,
		Amount:    txn.Amount,
		Kind:      txn.Kind,
		Timestamp: txn.Timestamp,
	}
}
