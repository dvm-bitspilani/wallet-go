package database

import (
	"context"
	"dvm.wallet/harsh/cmd/api/config"
	"dvm.wallet/harsh/ent"
	"dvm.wallet/harsh/ent/user"
	"dvm.wallet/harsh/internal/helpers"
)

type BalanceType int

const (
	SWD          BalanceType = 1
	CASH         BalanceType = 2
	PG           BalanceType = 3
	TRANSFER_BAL BalanceType = 4
)

// GetValidTransactionPairs check with seniors if this is the combo they're looking for
func GetValidTransactionPairs() [][]string {
	return [][]string{
		{helpers.BITSIAN, helpers.BITSIAN},
		{helpers.BITSIAN, helpers.VENDOR},
		{helpers.PARTICIPANT, helpers.PARTICIPANT},
		{helpers.PARTICIPANT, helpers.VENDOR},
		{helpers.TELLER, helpers.PARTICIPANT},
	}
}

func GetOrCreateSwdTeller(app *config.Application, ctx context.Context) *ent.Teller {
	var swdUser *ent.User
	var swdTeller *ent.Teller
	swdUser, err := app.Client.User.Query().Where(user.Username("SWD")).Only(ctx)
	if err != nil {
		swdUser = app.Client.User.Create().
			SetUsername(helpers.SWD_USERNAME).
			SetPassword(helpers.SWD_PASSWORD). // TODO:	write a random password generator and add it to user's default password func
			SetName("SWD").
			SetEmail("swd@example.com").
			SetOccupation(helpers.TELLER).
			SaveX(ctx)
	}
	swdTeller, err = swdUser.QueryTeller().Only(ctx)
	if err != nil {
		swdTeller = app.Client.Teller.Create().
			SetUser(swdUser).
			SaveX(ctx)
	}
	return swdTeller
}

func GetBalanceFromTransactionType(txn helpers.Txn_type) BalanceType {
	switch txn {
	case helpers.ADD_SWD:
		return SWD
	case helpers.ADD_CASH:
		return CASH
	case helpers.ADD_PG:
		return PG
	case helpers.TRANSFER:
		return TRANSFER_BAL
	case helpers.PURCHASE:
		return TRANSFER_BAL
	}
	return 0 // can be potentially dangerous
}
