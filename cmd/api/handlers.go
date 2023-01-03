package main

import (
	"context"
	"dvm.wallet/harsh/ent/user"
	"dvm.wallet/harsh/internal/password"
	"dvm.wallet/harsh/internal/request"
	"dvm.wallet/harsh/internal/response"
	"dvm.wallet/harsh/internal/validator"
	"net/http"
)

func (app *application) status(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"Status": "OK",
	}

	err := response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

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

func (app *application) createAuthenticationToken(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username  string              `json:"username"`
		Password  string              `json:"password"`
		Validator validator.Validator `json:"-"`
	}

	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}
	ctx := context.Background()
	user, err := app.client.User.Query().Where(user.Username(input.Username)).Only(ctx)
	if err != nil {
		app.invalidCredentials(w, r)
		return
	}

	input.Validator.CheckField(input.Username != "", "username", "Username is required")
	input.Validator.CheckField(user != nil, "username", "Username could not be found")

	if user != nil {
		passwordMatches, err := password.Matches(input.Password, user.Password)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		input.Validator.CheckField(input.Password != "", "Password", "Password is required")
		input.Validator.CheckField(passwordMatches, "Password", "Password is incorrect")
	}

	if input.Validator.HasErrors() {
		app.failedValidation(w, r, input.Validator)
		return
	}

	data, err := generate_jwt_token(user.ID, app.config.baseURL, app.config.jwt.secretKey)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) protected(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is a protected handler"))
}
