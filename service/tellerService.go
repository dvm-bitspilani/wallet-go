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
		//err := exceptions.Exception{Message: fmt.Sprintf("teller %s is disabled", teller.Edges.User.Username), Status: 403}
		return nil, err
	}
	if teller.Edges.User.Occupation == "bitsian" {
		err := fmt.Errorf("cash additions to BITSian wallets is not allowed")
		//err := exceptions.Exception{Message: "cash additions to BITSian wallets is not allowed", Status: 403}
		return nil, err
	}
	transaction, err := GenerateAndPerform(amount, helpers.ADD_CASH, teller.Edges.User, user, r.ctx, r.client)
	if err != nil {
		return nil, err
	}
	teller.Update().AddCashCollected(amount).SaveX(r.ctx)
	// TODO:	update_balance
	return transaction, nil
}

func (r *TellerOps) AddBySwd(teller *ent.Teller, user *ent.User, amount int) (*ent.Transactions, error) {
	if user.Occupation != "bitsian" {
		err := fmt.Errorf("only bitsians can add money via SWD")
		//err := exceptions.Exception{Message: "Only bitsians can add money via SWD", Status: 403}
		return nil, err
	}
	if teller.Edges.User.Username != "SWD" {
		err := fmt.Errorf("only the SWD teller may add money via SWD")
		//err := exceptions.Exception{Message: "Only the SWD teller may add money via SWD", Status: 403}
		return nil, err
	}
	transaction, err := GenerateAndPerform(amount, helpers.ADD_SWD, teller.Edges.User, user, r.ctx, r.client)
	if err != nil {
		return nil, err
	}
	teller.Update().AddCashCollected(amount).SaveX(r.ctx)
	//TODO:		update_balance
	return transaction, nil
}

// TODO: 	Add by PG
