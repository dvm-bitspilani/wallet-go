package main

import (
	"context"
	"dvm.wallet/harsh/cmd/api/config"
	"dvm.wallet/harsh/cmd/api/context"
	"dvm.wallet/harsh/cmd/api/errors"
	"dvm.wallet/harsh/ent/user"
	"fmt"
	"github.com/pascaldekloe/jwt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// 	TODO:	Implement disallow disabled users middleware

// We take advantage of factory functions here to create middlewares for us. This helps us to inject dependencies while
// following the middleware signature constraints.

func newRecoverPanic(app *config.Application) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					errors.ServerError(w, r, fmt.Errorf("%s", err), app)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func newAuthenticate(app *config.Application) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Vary", "Authorization")

			authorizationHeader := r.Header.Get("Authorization")

			if authorizationHeader != "" {
				headerParts := strings.Split(authorizationHeader, " ")

				if len(headerParts) == 2 && headerParts[0] == "Bearer" {
					token := headerParts[1]

					claims, err := jwt.HMACCheck([]byte(token), []byte(app.Config.Jwt.SecretKey))
					if err != nil {
						errors.InvalidAuthenticationToken(w, r, app)
						return
					}

					if !claims.Valid(time.Now()) {
						errors.InvalidAuthenticationToken(w, r, app)
						return
					}

					if claims.Issuer != app.Config.BaseURL {
						errors.InvalidAuthenticationToken(w, r, app)
						return
					}

					if !claims.AcceptAudience(app.Config.BaseURL) {
						errors.InvalidAuthenticationToken(w, r, app)
						return
					}

					userID, err := strconv.Atoi(claims.Subject)
					if err != nil {
						errors.ServerError(w, r, err, app)
						return
					}
					ctx := context.Background()
					user, err := app.Client.User.Query().Where(user.ID(userID)).Only(ctx)
					if err != nil {
						errors.ServerError(w, r, err, app)
						return
					}

					if user != nil {
						r = context_config.ContextSetAuthenticatedUser(r, user)
					}
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

func newRequireAuthenticatedUser(app *config.Application) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authenticatedUser := context_config.ContextGetAuthenticatedUser(r)

			if authenticatedUser == nil {
				errors.AuthenticationRequired(w, r, app)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
