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

func (r *walletOps) Add(wallet *ent.Wallet, amount int, balanceType database.BalanceType) error {
	if amount < 0 {
		err := fmt.Errorf("amount to add to wallet cannot be negative")
		return err
	}
	if balanceType == database.SWD {
		wallet.Update().AddSwd(amount).SaveX(r.ctx)
		return nil
	} else if balanceType == database.CASH {
		wallet.Update().AddCash(amount).SaveX(r.ctx)
		return nil
	} else if balanceType == database.PG {
		wallet.Update().AddPg(amount).SaveX(r.ctx)
		return nil
	} else if balanceType == database.TRANSFER_BAL {
		wallet.Update().AddTransfers(amount).SaveX(r.ctx)
		return nil
	} else {
		return errors.New("invalid addition of funds")
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

func (r *walletOps) Deduct(wallet *ent.Wallet, amount int) error {
	if amount < 0 {
		return errors.New("amount to deduct from the wallet cannot be negative")
	}
	if r.Balance(wallet) < amount {
		return fmt.Errorf("%s's current balance is %d", wallet.Edges.User.Username, r.Balance(wallet))
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
				amount -= wallet.Swd
			} else {
				wallet.Update().AddPg(-amount).SaveX(r.ctx)
			}
		} else {
			wallet.Update().AddCash(-amount).SaveX(r.ctx)
		}
	} else {
		wallet.Update().SetTransfers(-amount).SaveX(r.ctx)
	}
	return nil
}

func (r *walletOps) ToDict(wallet *ent.Wallet) map[string]string {
	return map[string]string{
		"user":          wallet.Edges.User.Username,
		"swd":           strconv.Itoa(wallet.Swd),
		"cash":          strconv.Itoa(wallet.Cash),
		"pg":            strconv.Itoa(wallet.Pg),
		"transfers":     strconv.Itoa(wallet.Transfers),
		"total_balance": strconv.Itoa(r.Balance(wallet)),
	}
}
