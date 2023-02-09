package handlers

import (
	"dvm.wallet/harsh/cmd/api/config"
	context_config "dvm.wallet/harsh/cmd/api/context"
	"dvm.wallet/harsh/cmd/api/errors"
	"dvm.wallet/harsh/ent"
	"dvm.wallet/harsh/ent/user"
	"dvm.wallet/harsh/internal/database"
	"dvm.wallet/harsh/internal/request"
	"dvm.wallet/harsh/internal/response"
	"dvm.wallet/harsh/service"
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

// TODO:	implement maintenance mode

func AddCash(app *config.Application) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var statusCode int
		var input struct {
			Amount int    `json:"amount"`
			QrCode string `json:"qr_code"`
		}
		tellerUser := context_config.ContextGetAuthenticatedUser(r)

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

		tellerOps := service.NewTellerOps(r.Context(), app)
		_, err, statusCode = tellerOps.AddByCash(tellerUser.QueryTeller().OnlyX(r.Context()), targetUser, input.Amount)
		if err != nil {
			errors.ErrorMessage(w, r, statusCode, err.Error(), nil, app)
			return
		}
		err = response.JSON(w, http.StatusOK, "Funds added to wallet successfully!")
		if err != nil {
			errors.ServerError(w, r, err, app)
			return
		}
	}
}

func AddSwd(app *config.Application) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var statusCode int
		var input struct {
			Amount int `json:"amount"`
		}
		usr := context_config.ContextGetAuthenticatedUser(r)
		err := request.DecodeJSON(w, r, &input)
		if err != nil {
			errors.BadRequest(w, r, err, app)
			return
		}

		if usr.Occupation != "bitsian" {
			usr.Update().SetDisabled(true).SaveX(r.Context())
			errors.ErrorMessage(w, r, 403, "Non bitsians cannot add money via. SWD. You have been disabled", nil, app)
			return
		}

		if input.Amount <= 0 {
			usr.Update().SetDisabled(true).SaveX(r.Context())
			errors.ErrorMessage(w, r, 400, "Amount cannot be negative or zero. You've been disabled", nil, app)
			return
		}

		swdTeller := database.GetOrCreateSwdTeller(app, r.Context())
		tellerOps := service.NewTellerOps(r.Context(), app)
		_, err, statusCode = tellerOps.AddBySwd(swdTeller, usr, input.Amount)
		if err != nil {
			errors.ErrorMessage(w, r, statusCode, err.Error(), nil, app)
			return
		}
		err = response.JSON(w, http.StatusOK, "Funds added to wallet successfully")
		if err != nil {
			errors.ServerError(w, r, err, app)
			return
		}
	}
}

func Transfer(app *config.Application) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			QrCode uuid.UUID `json:"qr_code"`
			UserId int       `json:"id"`
			Amount int       `json:"amount"`
		}
		var transferMode int
		var targetUser *ent.User

		err := request.DecodeJSON(w, r, &input)

		if input.Amount == 0 {
			errors.ErrorMessage(w, r, 400, "Amount cannot be 0", nil, app)
			return
		} else if input.Amount < 0 {
			errors.ErrorMessage(w, r, 412, "Insufficient funds", nil, app)
			return
		}
		transferMode = 1
		if input.QrCode == uuid.Nil {
			if input.UserId == 0 {
				errors.ErrorMessage(w, r, 400, "Missing key in request body", nil, app)
				return
			}
			transferMode = 2
		}
		if transferMode == 1 {
			var err error
			targetUser, err = app.Client.User.Query().Where(user.QrCode(input.QrCode)).Only(r.Context())
			if err != nil {
				errors.ErrorMessage(w, r, 404, fmt.Sprintf("User not found with QR %d", input.QrCode), nil, app)
				return
			}
		} else if transferMode == 2 {
			var err error
			targetUser, err = app.Client.User.Query().Where(user.ID(input.UserId)).Only(r.Context())
			if err != nil {
				errors.ErrorMessage(w, r, 404, fmt.Sprintf("User not found with ID %d", input.UserId), nil, app)
				return
			}
		}
		srcUser := context_config.ContextGetAuthenticatedUser(r)
		userOps := service.NewUserOps(r.Context(), app)
		_, err, statusCode := userOps.Transfer(srcUser, targetUser, input.Amount)
		if err != nil {
			errors.ErrorMessage(w, r, statusCode, err.Error(), nil, app)
			return
		}
		err = response.JSON(w, http.StatusOK, "Funds transferred from wallet successfully")
		if err != nil {
			errors.ServerError(w, r, err, app)
		}
	}
}

// TODO:	get_paytm_checksum
// TODO:	confirm_pg_payment

func GetUserQR(app *config.Application) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		usr := context_config.ContextGetAuthenticatedUser(r)
		var data struct {
			UserId int       `json:"user_id"`
			QrCode uuid.UUID `json:"qr_code"`
		}
		data.UserId = usr.ID
		data.QrCode = usr.QrCode

		err := response.JSON(w, http.StatusOK, &data)
		if err != nil {
			errors.ServerError(w, r, err, app)
		}
	}
}

func GetBalance(app *config.Application) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		usr := context_config.ContextGetAuthenticatedUser(r)
		wallet, err := usr.QueryWallet().Only(r.Context())
		if err != nil {
			errors.ErrorMessage(w, r, 404, "User does not have a wallet", nil, app)
			return
		}
		data := struct {
			Swd       int `json:"swd"`
			Cash      int `json:"cash"`
			Pg        int `json:"pg"`
			Transfers int `json:"transfers"`
		}{
			Swd:       wallet.Swd,
			Cash:      wallet.Cash,
			Pg:        wallet.Pg,
			Transfers: wallet.Transfers,
		}
		err = response.JSON(w, http.StatusOK, &data)
		if err != nil {
			errors.ServerError(w, r, err, app)
			return
		}
	}
}

func TransactionHistory(app *config.Application) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		usr := context_config.ContextGetAuthenticatedUser(r)
		//_, err := usr.Edges.WalletOrErr()
		_, err := usr.QueryWallet().Only(r.Context())
		if err != nil {
			errors.ErrorMessage(w, r, 404, "User does not have a wallet", nil, app)
			return
		}
		transactions := usr.QueryTransactions().AllX(r.Context())

		var txns []service.TransactionStruct
		transactionOps := service.NewTransactionOps(r.Context(), app)
		for _, txn := range transactions {
			txnStruct := transactionOps.ToDict(txn)
			txns = append(txns, *txnStruct)
		}

		err = response.JSON(w, http.StatusOK, &txns) // does this even work, need to verify
		if err != nil {
			errors.ServerError(w, r, err, app)
			return
		}
	}
}
