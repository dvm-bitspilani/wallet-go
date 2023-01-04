package handlers

import (
	"context"
	wallet "dvm.wallet/harsh/cmd/api"
	"dvm.wallet/harsh/cmd/api/config"
	"dvm.wallet/harsh/cmd/api/errors"
	"dvm.wallet/harsh/ent/user"
	"dvm.wallet/harsh/internal/password"
	"dvm.wallet/harsh/internal/request"
	"dvm.wallet/harsh/internal/response"
	"dvm.wallet/harsh/internal/validator"
	"net/http"
)

func Login(app *config.Application) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Username  string              `json:"username"`
			Password  string              `json:"password"`
			Validator validator.Validator `json:"-"`
		}

		err := request.DecodeJSON(w, r, &input)
		if err != nil {
			errors.BadRequest(w, r, err, app)
			return
		}
		ctx := context.Background()
		user, err := app.Client.User.Query().Where(user.Username(input.Username)).Only(ctx)
		if err != nil {
			errors.InvalidCredentials(w, r, app)
			return
		}

		input.Validator.CheckField(input.Username != "", "username", "Username is required")
		input.Validator.CheckField(user != nil, "username", "Username could not be found")

		if user != nil {
			passwordMatches, err := password.Matches(input.Password, user.Password)
			if err != nil {
				errors.ServerError(w, r, err, app)
				return
			}

			input.Validator.CheckField(input.Password != "", "Password", "Password is required")
			input.Validator.CheckField(passwordMatches, "Password", "Password is incorrect")
		}

		if input.Validator.HasErrors() {
			errors.FailedValidation(w, r, input.Validator, app)
			return
		}

		data, err := wallet.Generate_jwt_token(user.ID, app.Config.BaseURL, app.Config.Jwt.SecretKey)
		if err != nil {
			errors.ServerError(w, r, err, app)
			return
		}

		err = response.JSON(w, http.StatusOK, data)
		if err != nil {
			errors.ServerError(w, r, err, app)
		}
	}
}

// func (app *application) protected(w http.ResponseWriter, r *http.Request) {
//	w.Write([]byte("This is a protected handler"))
//}

//func (app *application) createUser(w http.ResponseWriter, r *http.Request) {
//	var input struct {
//		Email     string              `json:"Email"`
//		Password  string              `json:"Password"`
//		Validator validator.Validator `json:"-"`
//	}
//
//	err := request.DecodeJSON(w, r, &input)
//	if err != nil {
//		app.badRequest(w, r, err)
//		return
//	}
//
//	existingUser, err := app.db.GetUserByEmail(input.Email)
//	if err != nil {
//		app.serverError(w, r, err)
//		return
//	}
//
//	input.Validator.CheckField(input.Email != "", "Email", "Email is required")
//	input.Validator.CheckField(validator.Matches(input.Email, validator.RgxEmail), "Email", "Must be a valid email address")
//	input.Validator.CheckField(existingUser == nil, "Email", "Email is already in use")
//
//	input.Validator.CheckField(input.Password != "", "Password", "Password is required")
//	input.Validator.CheckField(len(input.Password) >= 8, "Password", "Password is too short")
//	input.Validator.CheckField(len(input.Password) <= 72, "Password", "Password is too long")
//	input.Validator.CheckField(validator.NotIn(input.Password, password.CommonPasswords...), "Password", "Password is too common")
//
//	if input.Validator.HasErrors() {
//		app.failedValidation(w, r, input.Validator)
//		return
//	}
//
//	hashedPassword, err := password.Hash(input.Password)
//	if err != nil {
//		app.serverError(w, r, err)
//		return
//	}
//
//	_, err = app.db.InsertUser(input.Email, hashedPassword)
//	if err != nil {
//		app.serverError(w, r, err)
//		return
//	}
//
//	w.WriteHeader(http.StatusNoContent)
//}
