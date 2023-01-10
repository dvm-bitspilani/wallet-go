package helpers

import (
	"context"
	"dvm.wallet/harsh/cmd/api/config"
	"dvm.wallet/harsh/ent"
	"dvm.wallet/harsh/ent/user"
)

// GetValidTransactionPairs check with seniors if this is the combo they're looking for
func GetValidTransactionPairs() [][]string {
	return [][]string{
		{"bitsian", "bitsian"},
		{"bitsian", "vendor"},
		{"participant", "participant"},
		{"participant", "vendor"},
		{"teller", "participant"},
	}
}

func GetOrCreateSwdTeller(app *config.Application, ctx context.Context) *ent.Teller {
	var swdUser *ent.User
	var swdTeller *ent.Teller
	swdUser, err := app.Client.User.Query().Where(user.Username("SWD")).Only(ctx)
	if err != nil {
		swdUser = app.Client.User.Create().
			SetUsername("SWD").
			SetPassword("swdgivememymoneybackwtf"). // TODO:	write a random password generator and add it to user's default password func
			SetName("SWD").
			SetEmail("swd@example.com").
			SaveX(ctx)
	}
	swdTeller, err = swdUser.Edges.TellerOrErr()
	if err != nil {
		swdTeller = app.Client.Teller.Create().
			SetUser(swdUser).
			SaveX(ctx)
	}
	return swdTeller
}
