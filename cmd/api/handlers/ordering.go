package handlers

import (
	"dvm.wallet/harsh/cmd/api/config"
	context_config "dvm.wallet/harsh/cmd/api/context"
	"dvm.wallet/harsh/cmd/api/errors"
	"dvm.wallet/harsh/ent"
	"dvm.wallet/harsh/ent/order"
	"dvm.wallet/harsh/ent/ordershell"
	"dvm.wallet/harsh/internal/helpers"
	"dvm.wallet/harsh/internal/request"
	"dvm.wallet/harsh/internal/response"
	"dvm.wallet/harsh/internal/validator"
	"dvm.wallet/harsh/service"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// Order : The request body of this endpoint is being changed to be a bit more verbose and intuitive
func Order(app *config.Application) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		usr := context_config.ContextGetAuthenticatedUser(r)
		vars := mux.Vars(r)
		var shellId int

		shellId, err := strconv.Atoi(vars["shell_id"])
		if err != nil {
			shellId = 0
		}

		if r.Method == "GET" {
			orderShellOps := service.NewOrderShellOps(r.Context(), app.Client)
			if shellId == 0 { // since initialized ints have 0 as their default value
				var data []service.OrderShellStruct
				for _, shell := range usr.QueryWallet().QueryShells().Order(ent.Asc(ordershell.FieldTimestamp)).AllX(r.Context()) {
					data = append(data, *orderShellOps.ToDict(shell))
				}
				err := response.JSON(w, http.StatusOK, &data)
				if err != nil {
					errors.ServerError(w, r, err, app)
					return
				}
			} else {
				shell, err := app.Client.OrderShell.Query().Where(ordershell.ID(shellId)).Only(r.Context())
				if err != nil {
					errors.ErrorMessage(w, r, 404, fmt.Sprintf("OrderShell not found for ID %d", shellId), nil, app)
					return
				}
				data := orderShellOps.ToDict(shell)
				err = response.JSON(w, http.StatusOK, &data)
				if err != nil {
					errors.ServerError(w, r, err, app)
					return
				}
			}
		} else if r.Method == "POST" {

			//{
			//	{
			//		vendor_id: 1
			//		order: [1,2,3,4,5]
			//	}
			//	{
			//		vendor_id: 3
			//		orders: [2, 5]
			//	}

			var input struct {
				Vendor []helpers.OrderActionVendorStruct `json:"user_order"`
			}

			err := request.DecodeJSON(w, r, &input)
			if err != nil {
				errors.BadRequest(w, r, err, app)
				return
			}

			if !(validator.In(usr.Occupation, "bitsian", "participant")) {
				errors.ErrorMessage(w, r, 403, "Only bitsians or participants may place orders", nil, app)
				return
			}
			userOps := service.NewUserOps(r.Context(), app.Client)
			data, err, statusCode := userOps.PlaceOrder(usr, input.Vendor)
			if err != nil {
				errors.ErrorMessage(w, r, statusCode, err.Error(), nil, app)
				return
			}
			err = response.JSON(w, http.StatusOK, &data)
			if err != nil {
				errors.ServerError(w, r, err, app)
				return
			}
		}
	}
}

func MakeOtpSeen(app *config.Application) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			OrderId int
		}
		err := request.DecodeJSON(w, r, &input)
		if err != nil {
			errors.BadRequest(w, r, err, app)
			return
		}
		orderObj, err := app.Client.Order.Query().Where(order.ID(input.OrderId)).Only(r.Context())
		if err != nil {
			errors.ErrorMessage(w, r, 404, "This order does not exist", nil, app)
			return
		}
		usr := context_config.ContextGetAuthenticatedUser(r)
		if orderObj.QueryShell().QueryWallet().OnlyX(r.Context()) != usr.QueryWallet().OnlyX(r.Context()) {
			errors.ErrorMessage(w, r, 403, "User did not place this order", nil, app)
			return
		}
		OrderOps := service.NewOrderOps(r.Context(), app.Client)
		if validator.In(orderObj.Status, helpers.FINISHED, helpers.READY) {
			OrderOps.MakeOtpSeen(orderObj)
			err := response.JSON(w, http.StatusOK, "OTP has been successfully seen!")
			if err != nil {
				errors.ServerError(w, r, err, app)
				return
			}
		} else {
			err := response.JSON(w, http.StatusPreconditionFailed, "Order is not yet Ready")
			if err != nil {
				errors.ServerError(w, r, err, app)
				return
			}
		}
	}
}
