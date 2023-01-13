package handlers

import (
	"dvm.wallet/harsh/cmd/api/config"
	context_config "dvm.wallet/harsh/cmd/api/context"
	"dvm.wallet/harsh/cmd/api/errors"
	"dvm.wallet/harsh/ent/item"
	"dvm.wallet/harsh/ent/order"
	"dvm.wallet/harsh/ent/user"
	vendor "dvm.wallet/harsh/ent/vendorschema"
	"dvm.wallet/harsh/internal/helpers"
	"dvm.wallet/harsh/internal/request"
	"dvm.wallet/harsh/internal/response"
	"dvm.wallet/harsh/internal/validator"
	"dvm.wallet/harsh/service"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func GetVendorOrders(app *config.Application) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		usr := context_config.ContextGetAuthenticatedUser(r)
		if !(usr.Occupation == "VendorSchema") {
			errors.ErrorMessage(w, r, 403, "Requesting user is not a VendorSchema", nil, app)
			return
		}
		vendor := usr.Edges.VendorSchema
		vars := mux.Vars(r)
		status := vars["status"]

		if status == "" {
			orders := vendor.QueryOrders().AllX(r.Context())
			orderOps := service.NewOrderOps(r.Context(), app.Client)
			var data []service.OrderStruct
			for _, order := range orders {
				data = append(data, orderOps.ToDict(order))
			}
			err := response.JSON(w, http.StatusOK, &data)
			if err != nil {
				errors.ServerError(w, r, err, app)
				return
			}
		}
		conversionMap := map[string]helpers.Status{
			"pending":  helpers.PENDING,
			"accepted": helpers.ACCEPTED,
			"ready":    helpers.READY,
			"finished": helpers.FINISHED,
			"declined": helpers.DECLINED,
		}
		orders := vendor.QueryOrders().Where(order.StatusEQ(conversionMap[status])).AllX(r.Context())
		orderOps := service.NewOrderOps(r.Context(), app.Client)
		var data []service.OrderStruct
		for _, order := range orders {
			data = append(data, orderOps.ToDict(order))
		}
		err := response.JSON(w, http.StatusOK, &data)
		if err != nil {
			errors.ServerError(w, r, err, app)
			return
		}
	}
}

// GetOrderDetails is redundant, check with the app team and get it removed possibly.
func GetOrderDetails(app *config.Application) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		type itemStruct struct {
			ItemId     int    `json:"itemclass_id"`
			Name       string `json:"name"`
			UnitPrice  int    `json:"unit_price"`
			Quantity   int    `json:"quantity"`
			Veg        bool   `json:"is_veg"`
			TotalPrice int    `json:"total_price"`
		}

		type OrderDetailStruct struct {
			ShellId    int            `json:"shell_id"`
			VendorName string         `json:"vendor_name"`
			Status     helpers.Status `json:"status"`
			Otp        string         `json:"otp"`
			Items      []itemStruct   `json:"items"`
		}

		vars := mux.Vars(r)
		orderId, err := strconv.Atoi(vars["order_id"])
		if err != nil {
			errors.ErrorMessage(w, r, 400, "Order ID is not valid", nil, app)
			return
		}
		usr := context_config.ContextGetAuthenticatedUser(r)

		order, err := app.Client.Order.Query().Where(order.ID(orderId)).Only(r.Context())
		if err != nil {
			errors.ErrorMessage(w, r, 404, fmt.Sprintf("no orders of ID %d found in the database", orderId), nil, app)
			return
		}

		orderItems := order.QueryIteminstances().AllX(r.Context())
		//orderItemsList := make([]map[string]string, len(orderItems))

		var orderItemsList []itemStruct
		for _, item := range orderItems {
			orderItemsList = append(orderItemsList, itemStruct{
				ItemId:     item.Edges.Item.ID,
				Name:       item.Edges.Item.Name,
				UnitPrice:  item.PricePerQuantity,
				Quantity:   item.Quantity,
				Veg:        item.Edges.Item.Veg,
				TotalPrice: item.PricePerQuantity * item.Quantity,
			})
		}

		orderDetails := OrderDetailStruct{
			ShellId:    order.Edges.Shell.ID,
			VendorName: order.Edges.VendorSchema.Name,
			Status:     order.Status,
			Otp:        order.Otp,
			Items:      orderItemsList,
		}
		if usr.Occupation == "VendorSchema" {
			if order.Edges.VendorSchema.ID != usr.Edges.VendorSchema.ID {
				errors.ErrorMessage(w, r, 403, "The given order is not handled by requesting VendorSchema", nil, app)
				return
			}
		} else {
			if order.Edges.Shell.Edges.Wallet.Edges.User.ID != usr.ID {
				errors.ErrorMessage(w, r, 403, "This order was not placed by the requesting user", nil, app)
				return
			}
		}
		err = response.JSON(w, http.StatusOK, &orderDetails)
		if err != nil {
			errors.ServerError(w, r, err, app)
			return
		}
	}

}

func GetOrderIdArrayDetails(app *config.Application) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			OrderIdList []int `json:"order_id_list"`
		}
		err := request.DecodeJSON(w, r, &input)
		if err != nil {
			errors.ErrorMessage(w, r, 400, "request body is not correct", nil, app)
			return
		}
		usr := context_config.ContextGetAuthenticatedUser(r)

		for _, orderId := range input.OrderIdList {

			type itemStruct struct {
				ItemId     int    `json:"itemclass_id"`
				Name       string `json:"name"`
				UnitPrice  int    `json:"unit_price"`
				Quantity   int    `json:"quantity"`
				Veg        bool   `json:"is_veg"`
				TotalPrice int    `json:"total_price"`
			}

			type OrderDetailStruct struct {
				ShellId    int            `json:"shell_id"`
				VendorName string         `json:"vendor_name"`
				Status     helpers.Status `json:"status"`
				Otp        string         `json:"otp"`
				Items      []itemStruct   `json:"items"`
			}

			order, err := app.Client.Order.Query().Where(order.ID(orderId)).Only(r.Context())
			if err != nil {
				errors.ErrorMessage(w, r, 404, fmt.Sprintf("no orders of ID %d found in the database", orderId), nil, app)
				return
			}

			orderItems := order.QueryIteminstances().AllX(r.Context())
			//orderItemsList := make([]map[string]string, len(orderItems))

			var orderItemsList []itemStruct
			for _, item := range orderItems {
				orderItemsList = append(orderItemsList, itemStruct{
					ItemId:     item.Edges.Item.ID,
					Name:       item.Edges.Item.Name,
					UnitPrice:  item.PricePerQuantity,
					Quantity:   item.Quantity,
					Veg:        item.Edges.Item.Veg,
					TotalPrice: item.PricePerQuantity * item.Quantity,
				})
			}

			orderDetails := OrderDetailStruct{
				ShellId:    order.Edges.Shell.ID,
				VendorName: order.Edges.VendorSchema.Name,
				Status:     order.Status,
				Otp:        order.Otp,
				Items:      orderItemsList,
			}
			if usr.Occupation == "VendorSchema" {
				if order.Edges.VendorSchema.ID != usr.Edges.VendorSchema.ID {
					errors.ErrorMessage(w, r, 403, "The given order is not handled by requesting VendorSchema", nil, app)
					return
				}
			} else {
				if order.Edges.Shell.Edges.Wallet.Edges.User.ID != usr.ID {
					errors.ErrorMessage(w, r, 403, "This order was not placed by the requesting user", nil, app)
					return
				}
			}
			err = response.JSON(w, http.StatusOK, &orderDetails)
			if err != nil {
				errors.ServerError(w, r, err, app)
				return
			}
		}
	}

}

// GetDayEarnings request has probably been altered a bit, let the app team know.
// This is also redundant, just get this removed
//func GetDayEarnings(app *config.Application) func(http.ResponseWriter, *http.Request) {
//	return func(w http.ResponseWriter, r *http.Request) {
//		//var input struct {
//		//	Time time.Time
//		//}
//		//err := request.DecodeJSON(w, r, &input)
//		//if err != nil {
//		//	errors.BadRequest(w, r, err, app)
//		//	return
//		//}
//		vars := mux.Vars(r)
//		//date.
//		usr := context_config.ContextGetAuthenticatedUser(r)
//		if usr.Occupation != "VendorSchema" {
//			errors.ErrorMessage(w, r, 403, "Requesting user is not a VendorSchema", nil, app)
//			return
//		}
//		var totalEarnings int
//		var dayEarnings int
//		var orderIdList []int
//		orders := usr.Edges.Vendor.QueryOrders().AllX(r.Context())
//		for _, order := range orders {
//			if order.Status == helpers.FINISHED {
//				totalEarnings += order.Price
//				if order.Edges.Shell.Timestamp.Day() == input.Time.Day() && order.Edges.Shell.Timestamp.Month() == input.Time.Month() {
//					dayEarnings += order.Price
//					orderIdList = append(orderIdList, order.ID)
//				}
//			}
//		}
//		var output struct {
//			DayEarnings   int   `json:"day_earnings"`
//			TotalEarnings int   `json:"total_earnings"`
//			Orders        []int `json:"orders"`
//		}
//		output.DayEarnings = dayEarnings
//		output.TotalEarnings = totalEarnings
//		output.Orders = orderIdList
//		err = response.JSON(w, http.StatusOK, &output)
//		if err != nil {
//			errors.ServerError(w, r, err, app)
//			return
//		}
//	}
//}

func GetDayListEarnings(app *config.Application) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			DateList []time.Time `json:"date_list"`
		}
		err := request.DecodeJSON(w, r, &input)
		if err != nil {
			errors.BadRequest(w, r, err, app)
			return
		}
		usr := context_config.ContextGetAuthenticatedUser(r)
		if usr.Occupation != "VendorSchema" {
			errors.ErrorMessage(w, r, 403, "Requesting user is not a VendorSchema", nil, app)
			return
		}

		for _, timestamp := range input.DateList {
			var totalEarnings int
			var dayEarnings int
			var orderIdList []int
			orders := usr.Edges.VendorSchema.QueryOrders().AllX(r.Context())
			for _, order := range orders {
				if order.Status == helpers.FINISHED {
					totalEarnings += order.Price
					if order.Edges.Shell.Timestamp.Day() == timestamp.Day() && order.Edges.Shell.Timestamp.Month() == timestamp.Month() {
						dayEarnings += order.Price
						orderIdList = append(orderIdList, order.ID)
					}
				}
			}
			var output struct {
				DayEarnings   int   `json:"day_earnings"`
				TotalEarnings int   `json:"total_earnings"`
				Orders        []int `json:"orders"`
			}
			output.DayEarnings = dayEarnings
			output.TotalEarnings = totalEarnings
			output.Orders = orderIdList
			err = response.JSON(w, http.StatusOK, &output)
			if err != nil {
				errors.ServerError(w, r, err, app)
				return
			}
		}
	}
}

func AdvanceOrders(app *config.Application) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		orderId, err := strconv.Atoi(vars["order_id"])
		if err != nil {
			errors.BadRequest(w, r, err, app)
			return
		}

		var input struct {
			NewStatus int `json:"new_status"`
		}
		err = request.DecodeJSON(w, r, &input)
		if err != nil {
			errors.BadRequest(w, r, err, app)
			return
		}

		usr := context_config.ContextGetAuthenticatedUser(r)
		if !(usr.Occupation == "VendorSchema") {
			usr.Update().SetDisabled(true).SaveX(r.Context())
		}

		order, err := app.Client.Order.Query().Where(order.ID(orderId)).Only(r.Context())
		if err != nil {
			errors.ErrorMessage(w, r, 404, fmt.Sprintf("Order %d not found", orderId), nil, app)
			return
		}

		if validator.In(input.NewStatus-int(order.Status), 0, 1) {
			errors.ErrorMessage(w, r, 403, "Invalid action", nil, app)
			return
		}

		orderOps := service.NewOrderOps(r.Context(), app.Client)
		_, err = orderOps.ChangeStatus(order, helpers.FromInt(input.NewStatus), usr)
		if err != nil {
			errors.ErrorMessage(w, r, 403, err.Error(), nil, app)
			return
		}
		//TODO:	Disable vendors if they're trying to access orders that do not belong to them
		err = response.JSON(w, http.StatusOK, "Successfully Updated!")
		if err != nil {
			errors.ServerError(w, r, err, app)
			return
		}
	}
}

func DeclineOrders(app *config.Application) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		orderId, err := strconv.Atoi(vars["order_id"])
		if err != nil {
			errors.BadRequest(w, r, err, app)
			return
		}

		usr := context_config.ContextGetAuthenticatedUser(r)
		if !(usr.Occupation == "VendorSchema") {
			usr.Update().SetDisabled(true).SaveX(r.Context())
		}

		order, err := app.Client.Order.Query().Where(order.ID(orderId)).Only(r.Context())
		if err != nil {
			errors.ErrorMessage(w, r, 404, fmt.Sprintf("Order %d not found", orderId), nil, app)
			return
		}
		orderOps := service.NewOrderOps(r.Context(), app.Client)

		err = orderOps.Decline(order)
		if err != nil {
			errors.ErrorMessage(w, r, 403, err.Error(), nil, app)
			return
		}
		err = response.JSON(w, http.StatusOK, "Successfully Declined!")
		if err != nil {
			errors.ServerError(w, r, err, app)
			return
		}
	}
}

func ToggleAvailability(app *config.Application) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		usr := context_config.ContextGetAuthenticatedUser(r)
		if usr.Occupation != "VendorSchema" {
			usr.Update().SetDisabled(true).SaveX(r.Context())
			errors.ErrorMessage(w, r, 403, "Requesting user is not a VendorSchema", nil, app)
			return
		}

		var input struct {
			ItemObjList []struct {
				ItemId               int `json:"item_id"`
				NewAvailabilityState int `json:"new_availability_state"`
			} `json:"item_id_list"`
		}
		err := request.DecodeJSON(w, r, &input)
		if err != nil {
			errors.BadRequest(w, r, err, app)
			return
		}
		type outputItem struct {
			ItemId    int  `json:"item_id"`
			Available bool `json:"is_available"`
		}
		var availabilityData struct {
			Items []outputItem `json:"items"`
		}
		for _, itemStruct := range input.ItemObjList {
			itemObject, err := app.Client.Item.Query().Where(item.ID(itemStruct.ItemId)).Only(r.Context())
			if err != nil {
				errors.ErrorMessage(w, r, 404, fmt.Sprintf("Item with ID %d does not exist", itemObject.ID), nil, app)
				return
			}
			if !validator.In(itemObject, usr.Edges.VendorSchema.QueryItems().AllX(r.Context())...) {
				usr.Update().SetDisabled(true).SaveX(r.Context())
				errors.ErrorMessage(w, r, 403, "Vendor has been disabled for trying to toggle the availibility of an item not belonging to them", nil, app)
				return
			}
			if !validator.In(itemStruct.NewAvailabilityState, 0, 1) {
				errors.ErrorMessage(w, r, 400, fmt.Sprintf("Invalid valye of new_availability state for item_id %d", itemStruct.ItemId), nil, app)
				return
			}
			var updatedItem outputItem
			updatedItem.ItemId = itemObject.ID
			if itemStruct.NewAvailabilityState == 0 {
				itemObject.Update().SetAvailable(false).SaveX(r.Context())
				updatedItem.Available = false
				availabilityData.Items = append(availabilityData.Items, updatedItem)
			} else if itemStruct.NewAvailabilityState == 1 {
				itemObject.Update().SetAvailable(true).SaveX(r.Context())
				updatedItem.Available = true
				availabilityData.Items = append(availabilityData.Items, updatedItem)
			}
			err = response.JSON(w, http.StatusOK, &availabilityData)
			if err != nil {
				errors.ServerError(w, r, err, app)
				return
			}
		}

	}
}

func GetMenu(app *config.Application) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		vendorId, err := strconv.Atoi(vars["vendor_id"])
		if err != nil {
			errors.BadRequest(w, r, err, app)
			return
		}
		vendorObject, err := app.Client.VendorSchema.Query().Where(vendor.ID(vendorId)).Only(r.Context())
		if err != nil {
			errors.ErrorMessage(w, r, 404, fmt.Sprintf("Vendor with ID %d does not exist", vendorId), nil, app)
			return
		}
		if vendorObject.Edges.User.Username == "PROF_SHOW" {
			errors.ErrorMessage(w, r, 403, "Vendor is a Prof Show", nil, app)
			return
		}
		type menuItem struct {
			Id          int    `json:"id"`
			Name        string `json:"name"`
			Price       int    `json:"price"`
			Description string `json:"description"`
			VendorId    int    `json:"vendor_id"`
			IsVeg       bool   `json:"is_veg"`
			//IsCombo     bool   `json:"is_combo"`
			IsAvailable bool `json:"is_available"`
		}
		var data []menuItem
		for _, item := range vendorObject.QueryItems().AllX(r.Context()) {
			data = append(data, menuItem{
				Id:          item.ID,
				Name:        item.Name,
				Price:       item.BasePrice,
				Description: item.Description,
				VendorId:    item.Edges.VendorSchema.ID,
				IsVeg:       item.Veg,
				//IsCombo:     item,
				IsAvailable: item.Available,
			})
		}
		err = response.JSON(w, http.StatusOK, &data)
		if err != nil {
			errors.ServerError(w, r, err, app)
			return
		}
	}
}

func GetAllVendorsWithMenu(app *config.Application) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		type menuItem struct {
			Id          int    `json:"id"`
			Name        string `json:"name"`
			Price       int    `json:"price"`
			Description string `json:"description"`
			VendorId    int    `json:"vendor_id"`
			IsVeg       bool   `json:"is_veg"`
			//IsCombo     bool   `json:"is_combo"`
			IsAvailable bool `json:"is_available"`
		}
		type vendorStruct struct {
			ID          int        `json:"id"`
			Name        string     `json:"name"`
			ImageUrl    url.URL    `json:"image_url"`
			Description string     `json:"description"`
			Closed      bool       `json:"closed"`
			Menu        []menuItem `json:"menu"`
			Address     string     `json:"address"`
		}

		var data []vendorStruct

		for _, vendor := range app.Client.VendorSchema.Query().Where(vendor.HasUserWith(user.UsernameNEQ("PROF_SHOW"))).AllX(r.Context()) {
			if vendor.Closed {
				continue
			}

			var menu []menuItem
			for _, item := range vendor.QueryItems().AllX(r.Context()) {
				menu = append(menu, menuItem{
					Id:          item.ID,
					Name:        item.Name,
					Price:       item.BasePrice,
					Description: item.Description,
					VendorId:    item.Edges.VendorSchema.ID,
					IsVeg:       item.Veg,
					IsAvailable: item.Available,
				})
			}

			data = append(data, vendorStruct{
				ID:          vendor.ID,
				Name:        vendor.Name,
				ImageUrl:    *vendor.ImageURL,
				Description: vendor.Description,
				Closed:      vendor.Closed,
				Menu:        menu,
				Address:     vendor.Address,
			})
		}
		err := response.JSON(w, http.StatusOK, &data)
		if err != nil {
			errors.ServerError(w, r, err, app)
			return
		}
	}
}