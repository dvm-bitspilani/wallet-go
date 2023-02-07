package service

import (
	"context"
	"dvm.wallet/harsh/cmd/api/config"
	"dvm.wallet/harsh/cmd/api/realtime"
	"dvm.wallet/harsh/ent"
	"dvm.wallet/harsh/internal/helpers"
	"fmt"
)

type TellerOps struct {
	ctx context.Context
	//client *ent.Client
	app *config.Application
}

func NewTellerOps(ctx context.Context, app *config.Application) *TellerOps {
	return &TellerOps{
		ctx: ctx,
		app: app,
		//client: client,
	}
}

func (r *TellerOps) AddByCash(teller *ent.Teller, user *ent.User, amount int) (*ent.Transactions, error, int) {
	var statusCode int
	tellerUsr := teller.QueryUser().OnlyX(r.ctx)
	if tellerUsr.Disabled == true {
		err := fmt.Errorf("teller %s is disabled", teller.QueryUser().OnlyX(r.ctx).Username)
		//err := exceptions.Exception{Message: fmt.Sprintf("teller %s is disabled", teller.Edges.User.Username), Status: 403}
		return nil, err, 403 // 403
	}
	if tellerUsr.Occupation == "bitsian" {
		err := fmt.Errorf("cash additions to BITSian wallets is not allowed")
		//err := exceptions.Exception{Message: "cash additions to BITSian wallets is not allowed", Status: 403}
		return nil, err, 403 // 403
	}
	transaction, err, statusCode := GenerateAndPerform(amount, helpers.ADD_CASH, tellerUsr, user, r.ctx, r.app.Client)
	if err != nil {
		return nil, err, statusCode
	}
	teller.Update().AddCashCollected(amount).SaveX(r.ctx)
	// TODO:	update_balance
	return transaction, nil, 0
}

func (r *TellerOps) AddBySwd(teller *ent.Teller, user *ent.User, amount int) (*ent.Transactions, error, int) {
	var statusCode int

	if user.Occupation != "bitsian" {
		err := fmt.Errorf("only bitsians can add money via SWD")
		//err := exceptions.Exception{Message: "Only bitsians can add money via SWD", Status: 403}
		return nil, err, 403
	}
	tellerUsr := teller.QueryUser().OnlyX(r.ctx)
	if tellerUsr.Username != "SWD" {
		err := fmt.Errorf("only the SWD teller may add money via SWD")
		//err := exceptions.Exception{Message: "Only the SWD teller may add money via SWD", Status: 403}
		return nil, err, 403
	}
	transaction, err, statusCode := GenerateAndPerform(amount, helpers.ADD_SWD, tellerUsr, user, r.ctx, r.app.Client)
	if err != nil {
		return nil, err, statusCode
	}
	teller.Update().AddCashCollected(amount).SaveX(r.ctx)
	//TODO:		update_balance
	walletOps := NewWalletOps(r.ctx, r.app.Client)
	realtime.UpdateBalance(r.app.Manager, user, walletOps.Balance(user.QueryWallet().OnlyX(r.ctx)))
	return transaction, nil, 0
}

// TODO: 	Add by PG
