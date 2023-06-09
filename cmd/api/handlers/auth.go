package handlers

import (
	"dvm.wallet/harsh/cmd/api/config"
	"dvm.wallet/harsh/cmd/api/errors"
	"dvm.wallet/harsh/ent"
	"dvm.wallet/harsh/ent/user"
	"dvm.wallet/harsh/internal/helpers"
	"dvm.wallet/harsh/service"
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
			Username  string            `json:"username"`
			Password  string            `json:"password"`
			IdToken   string            `json:"id_token"`
			Validator helpers.Validator `json:"-"`
		}

		err := helpers.DecodeJSON(w, r, &input)
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
				passwordMatches, err := helpers.Matches(input.Password, userObject.Password)
				if err != nil {
					errors.ServerError(w, r, err, app)
					return
				}
				if passwordMatches {
					switch userObject.Occupation {
					// no bitsian here because they're never gonna login through password auth,
					// but it's probably worth it to handle that case some time in the future.
					case helpers.PARTICIPANT:
						category = 1
					case helpers.VENDOR:
						category = 2
					case helpers.TELLER:
						category = 2
					case helpers.MANAGER:
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
				errors.InvalidCredentials(w, r, app)
				return
			}
			if helpers.NotIn(payload.Claims["iss"].(string), "https://accounts.google.com", "https://accounts.google.com") {
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
			userOps := service.NewUserOps(r.Context(), app)
			_, err := userOps.GetOrCreateWallet(userObject)
			if err != nil {
				errors.ErrorMessage(w, r, 500, "Something went wrong and we were unable to create the user's wallet", nil, app)
				return
			}

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
			if userObject.Occupation == helpers.VENDOR {
				service.PutVendorOrders(userObject.QueryVendorSchema().OnlyIDX(r.Context()), app, app.FirestoreClient)
			} else if userObject.Occupation == helpers.TELLER {
				// TODO:	implement websocket based
				//			put_teller_node here
				app.Logger.Debugf("PUT_TELLER_NODE")
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

		err = helpers.JSON(w, http.StatusOK, jwtPayload)

	}
}
