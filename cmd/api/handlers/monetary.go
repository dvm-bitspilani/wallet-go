package handlers

import (
	"dvm.wallet/harsh/cmd/api/config"
	context_config "dvm.wallet/harsh/cmd/api/context"
	"dvm.wallet/harsh/cmd/api/errors"
	"dvm.wallet/harsh/ent/user"
	"dvm.wallet/harsh/internal/request"
	"dvm.wallet/harsh/service"
	"github.com/google/uuid"
	"net/http"
)

// TODO:	implement maintenance mode

func AddCash(app *config.Application) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Amount int    `json:"amount"`
			QrCode string `json:"qr_code"`
		}

		err := request.DecodeJSON(w, r, &input)
		if err != nil {
			errors.BadRequest(w, r, err, app)
			return
		}
		usrUuid, err := uuid.Parse(input.QrCode)
		if err != nil {
			errors.BadRequest(w, r, err, app)
			return
		}
		targetUser, err := app.Client.User.Query().Where(user.QrCode(usrUuid)).Only(r.Context())
		if err != nil {
			errors.ErrorMessage(w, r, 404, "User not found with this QR code", nil, app)
			return
		}

		tellerUser := context_config.ContextGetAuthenticatedUser(r)
		if tellerUser.Occupation != "teller" {
			tellerUser.Update().SetDisabled(true)
			errors.ErrorMessage(w, r, 403, "Requesting user must be a teller", nil, app)
			return
		}

		if input.Amount < 0 {
			tellerUser.Update().SetDisabled(true)
			errors.ErrorMessage(w, r, 400, "Amount cannot be negative, You've been disabled", nil, app)
			return
		}
		if targetUser.Disabled {
			errors.ErrorMessage(w, r, 403, "Target user has been disabled. Money was not added", nil, app)
			return
		}

		tellerOps := service.NewTellerOps(r.Context(), app.Client)
		_, err = tellerOps.AddByCash(tellerUser.Edges.Teller, targetUser, input.Amount)
		if err != nil {
			errors.ErrorMessage(w, r, 403, err.Error(), nil, app)
			return
		}

	}
}
