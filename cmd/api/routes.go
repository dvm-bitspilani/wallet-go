package main

import (
	"dvm.wallet/harsh/cmd/api/config"
	"dvm.wallet/harsh/cmd/api/errors"
	"dvm.wallet/harsh/cmd/api/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

func routes(app *config.Application) http.Handler {
	mux := mux.NewRouter()

	mux.NotFoundHandler = http.HandlerFunc(errors.NotFound(app))
	mux.MethodNotAllowedHandler = http.HandlerFunc(errors.MethodNotAllowed(app))

	recoverPanic := newRecoverPanic(app)
	authenticate := newAuthenticate(app)
	mux.Use(recoverPanic)
	mux.Use(authenticate)

	mux.HandleFunc("/status", handlers.Status(app)).Methods("GET")

	// For websockets
	//mux.HandleFunc("/ws", helpers.WsEndpoint)

	//mux.HandleFunc("/users", app.createUser).Methods("POST")

	// Auth Routes
	mux.HandleFunc("/auth", handlers.Login(app)).Methods("POST")

	authenticatedRoutes := mux.NewRoute().Subrouter()

	requireAuthenticatedUser := newRequireAuthenticatedUser(app)
	authenticatedRoutes.Use(requireAuthenticatedUser)

	// Montary Routes
	authenticatedRoutes.HandleFunc("/monetory/add/cash", handlers.AddCash(app)).Methods("POST")

	//authenticatedRoutes.HandleFunc("/protected", app.protected).Methods("GET")

	return mux
}
