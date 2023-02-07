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

	// Health Routes
	mux.HandleFunc("/status", handlers.Status(app)).Methods("GET")

	//mux.HandleFunc("/users", app.createUser).Methods("POST")

	// Auth Routes
	mux.HandleFunc("/auth", handlers.Login(app)).Methods("POST")

	authenticatedRoutes := mux.NewRoute().Subrouter()
	requireAuthenticatedUser := newRequireAuthenticatedUser(app)
	disallowDisabledUser := newDisallowDisabledUser(app)

	// using middlewares on our subroute
	authenticatedRoutes.Use(requireAuthenticatedUser)
	authenticatedRoutes.Use(disallowDisabledUser)

	// For websocket
	authenticatedRoutes.HandleFunc("/ws", app.Manager.ServeWs)

	// Monetary Routes
	authenticatedRoutes.HandleFunc("/monetary/add/cash", handlers.AddCash(app)).Methods("POST")
	authenticatedRoutes.HandleFunc("/monetary/add/swd", handlers.AddSwd(app)).Methods("POST")
	authenticatedRoutes.HandleFunc("/monetary/transfers", handlers.Transfer(app)).Methods("POST")
	authenticatedRoutes.HandleFunc("/monetary/getqr", handlers.GetUserQR(app)).Methods("GET")
	authenticatedRoutes.HandleFunc("/monetary/balance", handlers.GetBalance(app)).Methods("GET")
	authenticatedRoutes.HandleFunc("/monetary/transactions/history", handlers.TransactionHistory(app)).Methods("GET")

	// Ordering Routes
	authenticatedRoutes.HandleFunc("/order/{shell_id:(?:|[0-9]+)}", handlers.Order(app)).Methods("GET", "POST")
	authenticatedRoutes.HandleFunc("/order/make_otp_seen", handlers.MakeOtpSeen(app)).Methods("POST")

	// Vendor Routes
	authenticatedRoutes.HandleFunc("/vendors/", handlers.GetAllVendorsWithMenu(app)).Methods("GET")
	authenticatedRoutes.HandleFunc("/vendor/{vendor_id}", handlers.GetVendor(app)).Methods("GET")
	authenticatedRoutes.HandleFunc("/vendor/{vendor_id}/items", handlers.GetMenu(app)).Methods("GET")
	authenticatedRoutes.HandleFunc("/vendor/orders/{status:(?:pending|accepted|ready|finished|declined|)}", handlers.GetVendorOrders(app)).Methods("GET")
	authenticatedRoutes.HandleFunc("/vendor/orders/idlist", handlers.GetOrderIdArrayDetails(app)).Methods("POST")
	authenticatedRoutes.HandleFunc("/vendor/order/{order_id}", handlers.GetOrderDetails(app)).Methods("GET")
	authenticatedRoutes.HandleFunc("/vendor/orders/{order_id}/change_status", handlers.AdvanceOrders(app)).Methods("POST")
	authenticatedRoutes.HandleFunc("/vendor/orders/{order_id}/decline", handlers.DeclineOrders(app)).Methods("GET")
	authenticatedRoutes.HandleFunc("/vendor/items/toggle_availability", handlers.ToggleAvailability(app)).Methods("POST")
	authenticatedRoutes.HandleFunc("/vendor/earnings-list", handlers.GetDayListEarnings(app)).Methods("POST")

	//authenticatedRoutes.HandleFunc("/protected", app.protected).Methods("GET")

	return mux
}
