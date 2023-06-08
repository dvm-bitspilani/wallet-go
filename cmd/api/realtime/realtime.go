// Package realtime: responsible for handling firestore related functionality,
package realtime

import (
	"cloud.google.com/go/firestore"
	"context"
	"dvm.wallet/harsh/cmd/api/config"
	"dvm.wallet/harsh/ent/user"
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
			order_ref := db.Collection("orders").Doc(strconv.Itoa(order.ID))
			batch.Set(order_ref, map[string]interface{}{
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
