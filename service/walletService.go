package service

import (
	"context"
	"dvm.wallet/harsh/ent"
	"dvm.wallet/harsh/internal/database"
	"errors"
	"fmt"
	"strconv"
)

type walletOps struct {
	ctx    context.Context
	client *ent.Client
}

func NewWalletOps(ctx context.Context, client *ent.Client) *walletOps {
	return &walletOps{
		ctx:    ctx,
		client: client,
	}
}

func (r *walletOps) Balance(wallet *ent.Wallet) int {
	return wallet.Swd + wallet.Cash + wallet.Pg + wallet.Transfers
}

// TODO: update_balance by overiding the save method

func (r *walletOps) Add(wallet *ent.Wallet, amount int, balanceType database.BalanceType) (error, int) {
	if amount < 0 {
		err := fmt.Errorf("amount to add to wallet cannot be negative")
		return err, 400 // 400
	}
	if balanceType == database.SWD {
		wallet.Update().AddSwd(amount).SaveX(r.ctx)
		return nil, 0
	} else if balanceType == database.CASH {
		wallet.Update().AddCash(amount).SaveX(r.ctx)
		return nil, 0
	} else if balanceType == database.PG {
		wallet.Update().AddPg(amount).SaveX(r.ctx)
		return nil, 0
	} else if balanceType == database.TRANSFER_BAL {
		wallet.Update().AddTransfers(amount).SaveX(r.ctx)
		return nil, 0
	} else {
		return errors.New("invalid addition of funds"), 0
	}
}

func (r *walletOps) AddAll(wallet *ent.Wallet, addDict map[string]int) {
	swd, ok := addDict["swd"]
	if !ok {
		swd = 0
	}
	cash, ok := addDict["cash"]
	if !ok {
		swd = 0
	}
	pg, ok := addDict["pg"]
	if !ok {
		swd = 0
	}
	transfers, ok := addDict["transfers"]
	if !ok {
		swd = 0
	}
	wallet.Update().AddSwd(swd).AddCash(cash).AddPg(pg).AddTransfers(transfers).SaveX(r.ctx)
}

func (r *walletOps) Deduct(wallet *ent.Wallet, amount int) (error, int) {
	if amount < 0 {
		return errors.New("amount to deduct from the wallet cannot be negative"), 400 // 400
	}
	if r.Balance(wallet) < amount {
		return fmt.Errorf("%s's current balance is %d", wallet.QueryUser().OnlyX(r.ctx).Username, r.Balance(wallet)), 412 // 412
	}
	if wallet.Transfers < amount {
		amount -= wallet.Transfers
		wallet.Update().SetTransfers(0).SaveX(r.ctx)
		if wallet.Cash < amount {
			amount -= wallet.Cash
			wallet.Update().SetCash(0).SaveX(r.ctx)
			if wallet.Pg < amount {
				amount -= wallet.Pg
				wallet.Update().SetPg(0).SaveX(r.ctx)
				wallet.Update().AddSwd(-amount).SaveX(r.ctx)
			} else {
				wallet.Update().AddPg(-amount).SaveX(r.ctx)
			}
		} else {
			wallet.Update().AddCash(-amount).SaveX(r.ctx)
		}
	} else {
		wallet.Update().AddTransfers(-amount).SaveX(r.ctx)
	}
	return nil, 0
}

func (r *walletOps) ToDict(wallet *ent.Wallet) map[string]string {
	return map[string]string{
		"user":          wallet.QueryUser().OnlyX(r.ctx).Username,
		"swd":           strconv.Itoa(wallet.Swd),
		"cash":          strconv.Itoa(wallet.Cash),
		"pg":            strconv.Itoa(wallet.Pg),
		"transfers":     strconv.Itoa(wallet.Transfers),
		"total_balance": strconv.Itoa(r.Balance(wallet)),
	}
}
