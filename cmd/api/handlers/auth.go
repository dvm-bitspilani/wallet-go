package handlers

import (
	"dvm.wallet/harsh/cmd/api/config"
	"dvm.wallet/harsh/cmd/api/errors"
	"dvm.wallet/harsh/ent"
	"dvm.wallet/harsh/ent/user"
	"dvm.wallet/harsh/internal/helpers"
	"dvm.wallet/harsh/internal/password"
	"dvm.wallet/harsh/internal/request"
	"dvm.wallet/harsh/internal/response"
	"dvm.wallet/harsh/internal/validator"
	"fmt"
	"google.golang.org/api/idtoken"
	"net/http"
	"strconv"
)

func Login(app *config.Application) func(http.ResponseWriter, *http.Request) {

	//	The auth endpoint services 3 different categories of users.
	//		Category 1:
	//	Bitsians - Google OAuth 2.0.
	//		Participants - Simple username/password based login
	//	both require many extra fields like qr_code, etc. to be sent.
	//		Category 2:
	//	Vendors and Tellers - Simple username/password based login.
	//		Category 3:
	//	Show Managers - Also simple username/password based login,
	//		however unlike the other kinds of people, they do not have a
	//	UserProfile instance associated with them.
	//
	//		Required keys for this endpoint:
	//		Mode 1: Simple username/password
	//			username: str
	//			password: str
	//		Mode 2: OAuth 2.0
	//			id_token: str
	//	NOTE: In mode 2, if the bitsian does not exist then we must create them.
	//
	//		Extra/optional keys:
	//			reg_token: str - used by the clients for firebase notifications.
	//			avatar: int - used by the android app.
	//
	//		Return:
	//			Categories 1:
	//				1. JWT: str
	//				2. qr_code: str
	//				3. user_id: int
	//				4. name: str
	//				5. email: str
	//				6. phone: str
	//				7. bitsian_id: str - Only for category 1.
	//			Category 2:
	//				1. JWT: str
	//				2. user_id: int
	//			Category 3:
	//				1. JWT: str
	//				2. shows: Dict[int, str]

	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Username  string              `json:"username"`
			Password  string              `json:"password"`
			IdToken   string              `json:"id_token"`
			Validator validator.Validator `json:"-"`
		}

		err := request.DecodeJSON(w, r, &input)
		if err != nil {
			errors.BadRequest(w, r, err, app)
			return
		}

		// TODO: check if these ever do not get assigned to the valid values
		var authMode int
		var category int
		var userObject *ent.User
		var jwtPayload map[string]string

		if input.Username != "" {
			authMode = 1
		} else if input.IdToken != "" {
			authMode = 2
		} else {
			errors.ErrorMessage(w, r, 400, "Insufficient authentication parameters", nil, app)
			return
		}

		if authMode == 1 {
			userObject, err = app.Client.User.Query().Where(user.Username(input.Username)).Only(r.Context())
			if err != nil {
				errors.ErrorMessage(w, r, 404, "User not found", nil, app)
				return
			}
			if userObject != nil {
				passwordMatches, err := password.Matches(input.Password, userObject.Password)
				if err != nil {
					errors.ServerError(w, r, err, app)
					return
				}
				if passwordMatches {
					switch userObject.Occupation {
					// no bitsian here because they're never gonna login through password auth,
					// but it's probably worth it to handle that case some time in the future.
					case "participant":
						category = 1
					case "vendor":
						category = 2
					case "teller":
						category = 2
					case "manager":
						category = 3
					default:
						errors.ErrorMessage(w, r, 403, "Undeterminable Occupation", nil, app)
						return
					}
				}
			}

		} else if authMode == 2 {
			category = 1
			payload, err := idtoken.Validate(r.Context(), input.IdToken, "") //LOW:	Can't ise audience here, look into why
			if err != nil {
				//errors.InvalidCredentials(w, r, app)
				app.Logger.Println(err)
				errors.InvalidCredentials(w, r, app)
				return
			}
			if validator.NotIn(payload.Claims["iss"].(string), "https://accounts.google.com", "https://accounts.google.com") {
				errors.ErrorMessage(w, r, 401, "Not a valid google account", nil, app)
			}
			email := payload.Claims["email"].(string)

			hd, ok := payload.Claims["hd"]
			if !ok {
				errors.ErrorMessage(w, r, 401, "No hosted domain in the token", nil, app)
				return
			}
			if hd.(string) != "pilani.bits-pilani.ac.in" {
				errors.ErrorMessage(w, r, 401, "Not a valid BITSian account", nil, app)
				return
			}
			userObject, err = app.Client.User.Query().Where(user.Email(email)).Only(r.Context())
			if err != nil {
				errors.ErrorMessage(w, r, 401, "User not in DVM's database", nil, app)
				return
			}

		}

		jwt, err := helpers.GenerateJwtToken(userObject.ID, app.Config.BaseURL, app.Config.Jwt.SecretKey)
		if err != nil {
			errors.ErrorMessage(w, r, 500, fmt.Sprintf("Could not generate a JWT token for user:%s", userObject.Username), nil, app)
			return
		}

		if category == 1 {
			// create payload
			jwtPayload = map[string]string{
				"AuthenticationToken":       jwt["AuthenticationToken"],
				"AuthenticationTokenExpiry": jwt["AuthenticationTokenExpiry"],
				"user_id":                   strconv.Itoa(userObject.ID),
				"name":                      userObject.Name,
				"email":                     userObject.Email,
			}
		} else if category == 2 {
			jwtPayload = map[string]string{
				"AuthenticationToken":       jwt["AuthenticationToken"],
				"AuthenticationTokenExpiry": jwt["AuthenticationTokenExpiry"],
				"user_id":                   strconv.Itoa(userObject.ID),
			}
			if userObject.Occupation == "vendor" {
				// TODO:	websocket implementation here
				// 			implement put_vendor_orders
				//			(Also check if its really required)
				app.Logger.Println("PUT_VENDOR_ORDERS")
			} else if userObject.Occupation == "teller" {
				// TODO:	implement websocket based
				//			put_teller_node here
				app.Logger.Println("PUT_TELLER_NODE")
			}
		} else if category == 3 {
			jwtPayload = map[string]string{
				"AuthenticationToken":       jwt["AuthenticationToken"],
				"AuthenticationTokenExpiry": jwt["AuthenticationTokenExpiry"],
				"shows":                     "SHOWS", // TODO:	implement showmanager workflow
			}
		} else {
			errors.ErrorMessage(w, r, 500, "Something went wrong", nil, app)
			return
		}

		err = response.JSON(w, http.StatusOK, jwtPayload)

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
