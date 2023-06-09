// Package realtime: responsible for handling firestore related functionality,
package service

import (
	"cloud.google.com/go/firestore"
	"context"
	"dvm.wallet/harsh/cmd/api/config"
	"dvm.wallet/harsh/ent/order"
	"dvm.wallet/harsh/ent/user"
	"dvm.wallet/harsh/ent/vendorschema"
	"dvm.wallet/harsh/internal/helpers"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	"strconv"
)

// NewFirestoreClient Creates and returns a new firestore client
func NewFirestoreClient(filePath string) (*firestore.Client, error) {
	ctx := context.Background()
	sa := option.WithCredentialsFile(filePath)
	clientApp, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		return nil, err
	}
	db, err := clientApp.Firestore(ctx)
	if err != nil {
		return nil, err
	}
	//defer db.Close()
	return db, nil
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
				"vendorid":     order.QueryShell().QueryOrders().QueryVendorSchema().OnlyIDX(ctx),
				"vendoruserid": order.QueryShell().QueryWallet().QueryUser().OnlyIDX(ctx),
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
			"earnings": CalculateEarnings(usr.QueryVendorSchema().OnlyX(ctx)),
		}, firestore.MergeAll)

	} else if usr.Occupation == helpers.TELLER {
		balanceRef := db.Collection("tellers").Doc(strconv.Itoa(usr.ID))
		balanceRef.Set(ctx, map[string]interface{}{
			"cash_collected": usr.QueryTeller().OnlyX(ctx).CashCollected,
		}, firestore.MergeAll)

	} else if helpers.In(usr.Occupation, helpers.BITSIAN, helpers.PARTICIPANT) {
		balanceRef := db.Collection("users").Doc(strconv.Itoa(usr.ID))
		userOps := NewUserOps(ctx, app)
		wallet, err := userOps.GetOrCreateWallet(usr)
		if err != nil {
			app.Logger.Errorf("Could not create wallet for user: %s", err)
		}
		walletOps := NewWalletOps(ctx, app)
		balanceRef.Set(ctx, map[string]interface{}{
			"total_balance":      walletOps.Balance(wallet),
			"refundable_balance": wallet.Swd,
		}, firestore.MergeAll)
	}
}

func UpdateOrderStatus(orderId int, app *config.Application, db *firestore.Client) {
	ctx := context.Background()
	orderObj := app.Client.Order.Query().Where(order.ID(orderId)).OnlyX(ctx)
	orderRef := db.Collection("orders").Doc(strconv.Itoa(orderObj.ID))
	orderRef.Set(ctx, map[string]interface{}{
		"status": orderObj.Status,
	}, firestore.MergeAll)
}

func PutVendorOrders(vendorId int, app *config.Application, db *firestore.Client) {
	ctx := context.Background()
	batch := db.Batch()
	vendor := app.Client.VendorSchema.Query().Where(vendorschema.ID(vendorId)).OnlyX(ctx)

	for key, order := range vendor.QueryOrders().AllX(ctx) {
		if key%498 == 0 {
			_, err := batch.Commit(ctx)
			if err != nil {
				app.Logger.Errorf("Could not commit batch: %s", err)
			}
		}
		orderRef := db.Collection("orders").Doc(strconv.Itoa(order.ID))
		batch.Set(orderRef, map[string]interface{}{
			"status":   order.Status,
			"userid":   order.QueryShell().QueryWallet().QueryUser().OnlyIDX(ctx),
			"vendorid": order.QueryVendorSchema().OnlyIDX(ctx),
			"otp_seen": order.OtpSeen,
		})
	}
}

func UpdateOtpSeen(orderID int, app *config.Application, db *firestore.Client) {
	ctx := context.Background()
	orderObj := app.Client.Order.Query().Where(order.ID(orderID)).OnlyX(ctx)
	orderRef := db.Collection("orders").Doc(strconv.Itoa(orderObj.ID))
	orderRef.Set(ctx, map[string]interface{}{
		"otp_seen": orderObj.OtpSeen,
	}, firestore.MergeAll)
}