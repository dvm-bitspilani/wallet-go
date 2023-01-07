package service

import (
	"context"
	"dvm.wallet/harsh/ent"
	"dvm.wallet/harsh/internal/helpers"
	"fmt"
)

type TellerOps struct {
	ctx    context.Context
	client *ent.Client
}

func NewTellerOps(ctx context.Context, client *ent.Client) *TellerOps {
	return &TellerOps{
		ctx:    ctx,
		client: client,
	}
}

func (r *TellerOps) AddByCash(teller *ent.Teller, user *ent.User, amount int) (*ent.Transactions, error) {
	if teller.Edges.User.Disabled == true {
		err := fmt.Errorf("teller %s is disabled", teller.Edges.User.Username)
		return nil, err
	}
	if teller.Edges.User.Occupation == "bitsian" {
		err := fmt.Errorf("cash additions to BITSian wallets is not allowed")
		return nil, err
	}

	// TODO:	Transaction generate_and_perform
	transaction, err := GenerateAndPerform(amount, helpers.TRANSFER, teller.Edges.User, user, r.ctx, r.client)
	if err != nil {
		return nil, err
	}
	teller.Update().AddCashCollected(amount).SaveX(r.ctx)
	// update_balance
	return transaction, nil
}
