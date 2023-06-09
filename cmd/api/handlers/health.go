package handlers

import (
	"dvm.wallet/harsh/cmd/api/config"
	"dvm.wallet/harsh/cmd/api/errors"
	"dvm.wallet/harsh/internal/helpers"
	"net/http"
)

func Status(app *config.Application) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data := map[string]string{
			"Status": "OK",
		}

		err := helpers.JSON(w, http.StatusOK, data)
		if err != nil {
			errors.ServerError(w, r, err, app)
		}
	}
}
