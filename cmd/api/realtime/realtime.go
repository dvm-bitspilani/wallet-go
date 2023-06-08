// Package realtime: responsible for handling firestore related functionality,
package realtime

import (
	"cloud.google.com/go/firestore"
	"context"
	"dvm.wallet/harsh/cmd/api/config"
	"dvm.wallet/harsh/ent/order"
	"dvm.wallet/harsh/ent/user"
	"dvm.wallet/harsh/internal/helpers"
	"dvm.wallet/harsh/internal/validator"
	"dvm.wallet/harsh/service"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	"strconv"
)

// NewFirestoreClient Creates and returns a new firestore client
func NewFirestoreClient(app *config.Application) *firestore.Client {
	ctx := context.Background()
	sa := option.WithCredentialsFile("internal/firebase-keyconfig.json")
	clientApp, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		app.Logger.Errorf("Error occuered while creating firebase app %s", err)
	}
	db, err := clientApp.Firestore(ctx)
	if err != nil {
		app.Logger.Errorf("Error occuered while creating firestore client %s", err)
	}
	//defer db.Close()
	return db
}

// PutUserOrders The functionality of uploading user Orders to Firebase is baked in tightly within this function.
func PutUserOrders(userId int, app *config.Application, db *firestore.Client) {
	ctx := context.Background()
	batch := db.Batch()
	usr := app.Client.User.Query().Where(user.ID(userId)).OnlyX(ctx)
	orderShells := usr.QueryWallet().QueryShells().AllX(ctx)
	for _, ordershell := range orderShells {
		for key, order := range ordershell.QueryOrders().AllX(ctx) {
			if key%498 == 0 {
				_, err := batch.Commit(ctx)
				if err != nil {
					app.Logger.Errorf("Could not commit batch: %s", err)
				}
			}
			orderRef := db.Collection("orders").Doc(strconv.Itoa(order.ID))
			batch.Set(orderRef, map[string]interface{}{
				"status":       order.Status,
				"userid":       usr.ID,
				"vendorid":     order.QueryShell().QueryOrders().QueryVendorSchema().OnlyX(ctx).ID,
				"vendoruserid": order.QueryShell().QueryWallet().QueryUser().OnlyX(ctx).ID,
				"otp_seen":     order.OtpSeen,
			}, firestore.MergeAll)
		}
	}
	_, err := batch.Commit(ctx)
	if err != nil {
		return
	}
}

func UpdateBalance(userId int, app *config.Application, db *firestore.Client) {
	ctx := context.Background()
	usr := app.Client.User.Query().Where(user.ID(userId)).OnlyX(ctx)

	if usr.Occupation == helpers.VENDOR {
		if usr.Username == helpers.PROFSHOW_USERNAME {
			return
		}
		balanceRef := db.Collection("vendors").Doc(strconv.Itoa(usr.ID))
		balanceRef.Set(ctx, map[string]interface{}{
			"earnings": service.CalculateEarnings(usr.QueryVendorSchema().OnlyX(ctx)),
		}, firestore.MergeAll)

	} else if usr.Occupation == helpers.TELLER {
		balanceRef := db.Collection("tellers").Doc(strconv.Itoa(usr.ID))
		balanceRef.Set(ctx, map[string]interface{}{
			"cash_collected": usr.QueryTeller().OnlyX(ctx).CashCollected,
		}, firestore.MergeAll)

	} else if validator.In(usr.Occupation, helpers.BITSIAN, helpers.PARTICIPANT) {
		balanceRef := db.Collection("users").Doc(strconv.Itoa(usr.ID))
		userOps := service.NewUserOps(ctx, app)
		wallet, err := userOps.GetOrCreateWallet(usr)
		if err != nil {
			app.Logger.Errorf("Could not create wallet for user: %s", err)
		}
		walletOps := service.NewWalletOps(ctx, app)
		balanceRef.Set(ctx, map[string]interface{}{
			"total_balance":      walletOps.Balance(wallet),
			"refundable_balance": wallet.Swd,
		}, firestore.MergeAll)
	}
}

func UpdateOrderStatus(orderId int, status helpers.Status, app *config.Application, db *firestore.Client) {
	ctx := context.Background()
	orderObj := app.Client.Order.Query().Where(order.ID(orderId)).OnlyX(ctx)
	orderRef := db.Collection("orders").Doc(strconv.Itoa(orderObj.ID))
	orderRef.Set(ctx, map[string]interface{}{
		"status": orderObj.Status,
	}, firestore.MergeAll)
}
