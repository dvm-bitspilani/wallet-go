package service

import (
	"context"
	"dvm.wallet/harsh/ent"
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
	wallet, err := user.QueryWallet().Only(r.ctx)
	if err != nil {
		wallet, err = r.client.Wallet.Create().SetUser(user).Save(r.ctx)
		if err != nil {
			return nil, err
		}
	}
	return wallet, nil
}
