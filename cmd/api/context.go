package wallet

import (
	"context"
	"dvm.wallet/harsh/ent"
	"net/http"
)

type contextKey string

const (
	authenticatedUserContextKey = contextKey("authenticatedUser")
)

func contextSetAuthenticatedUser(r *http.Request, user *ent.User) *http.Request {
	ctx := context.WithValue(r.Context(), authenticatedUserContextKey, user)
	return r.WithContext(ctx)
}

func contextGetAuthenticatedUser(r *http.Request) *ent.User {
	user, ok := r.Context().Value(authenticatedUserContextKey).(*ent.User)
	if !ok {
		return nil
	}

	return user
}
