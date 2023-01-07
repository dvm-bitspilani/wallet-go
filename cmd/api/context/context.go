package context_config

import (
	"context"
	"dvm.wallet/harsh/ent"
	"net/http"
)

type contextKey string

const (
	authenticatedUserContextKey = contextKey("authenticatedUser")
)

func ContextSetAuthenticatedUser(r *http.Request, user *ent.User) *http.Request {
	ctx := context.WithValue(r.Context(), authenticatedUserContextKey, user)
	return r.WithContext(ctx)
}

func ContextGetAuthenticatedUser(r *http.Request) *ent.User {
	user, ok := r.Context().Value(authenticatedUserContextKey).(*ent.User)
	if !ok {
		return nil
	}

	return user
}
