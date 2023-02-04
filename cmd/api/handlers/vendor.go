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
	"strconv"
	"time"
)

func GetVendorOrders(app *config.Application) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		usr := context_config.ContextGetAuthenticatedUser(r)
		if usr.Occupation != "vendor" {
			errors.ErrorMessage(w, r, 403, "Requesting user is not a VendorSchema", nil, app)
			return
		}
		vendorObj := usr.QueryVendorSchema().OnlyX(r.Context())
		vars := mux.Vars(r)
		status := vars["status"]
		//check what empty vars does here
		if status == "" {
			orders := vendorObj.QueryOrders().AllX(r.Context())
			orderOps := service.NewOrderOps(r.Context(), app.Client)
			var data []service.OrderStruct
			for _, orderObj := range orders {
				data = append(data, orderOps.ToDict(orderObj))
			}
			err := response.JSON(w, http.StatusOK, &data)
			return
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
		orders := vendorObj.QueryOrders().Where(order.StatusEQ(conversionMap[status])).AllX(r.Context())
		orderOps := service.NewOrderOps(r.Context(), app.Client)
		var data []service.OrderStruct
		for _, orderObj := range orders {
			data = append(data, orderOps.ToDict(orderObj))
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

		orderObj, err := app.Client.Order.Query().Where(order.ID(orderId)).Only(r.Context())
		if err != nil {
			errors.ErrorMessage(w, r, 404, fmt.Sprintf("no orders of ID %d found in the database", orderId), nil, app)
			return
		}

		orderItems := orderObj.QueryIteminstances().AllX(r.Context())
		//orderItemsList := make([]map[string]string, len(orderItems))

		var orderItemsList []itemStruct
		for _, orderItem := range orderItems {
			menuItem := orderItem.QueryItem().OnlyX(r.Context())
			orderItemsList = append(orderItemsList, itemStruct{
				ItemId:     menuItem.ID,
				Name:       menuItem.Name,
				UnitPrice:  orderItem.PricePerQuantity,
				Quantity:   orderItem.Quantity,
				Veg:        menuItem.Veg,
				TotalPrice: orderItem.PricePerQuantity * orderItem.Quantity,
			})
		}
		orderShell := orderObj.QueryShell().OnlyX(r.Context())
		orderVendorSchema := orderObj.QueryVendorSchema().OnlyX(r.Context())
		orderDetails := OrderDetailStruct{
			ShellId:    orderShell.ID,
			VendorName: orderVendorSchema.Name,
			Status:     orderObj.Status,
			Otp:        orderObj.Otp,
			Items:      orderItemsList,
		}
		if usr.Occupation == "vendor" {
			if orderVendorSchema.ID != usr.QueryVendorSchema().OnlyX(r.Context()).ID {
				errors.ErrorMessage(w, r, 403, "The given order is not handled by requesting VendorSchema", nil, app)
				return
			}
		} else {
			if orderShell.QueryWallet().QueryUser().OnlyX(r.Context()).ID != usr.ID {
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

			orderObj, err := app.Client.Order.Query().Where(order.ID(orderId)).Only(r.Context())
			if err != nil {
				errors.ErrorMessage(w, r, 404, fmt.Sprintf("no orders of ID %d found in the database", orderId), nil, app)
				return
			}

			orderItems := orderObj.QueryIteminstances().AllX(r.Context())
			//orderItemsList := make([]map[string]string, len(orderItems))

			var orderItemsList []itemStruct
			for _, orderItem := range orderItems {
				menuItem := orderItem.QueryItem().OnlyX(r.Context())
				orderItemsList = append(orderItemsList, itemStruct{
					ItemId:     menuItem.ID,
					Name:       menuItem.Name,
					UnitPrice:  orderItem.PricePerQuantity,
					Quantity:   orderItem.Quantity,
					Veg:        menuItem.Veg,
					TotalPrice: orderItem.PricePerQuantity * orderItem.Quantity,
				})
			}
			orderShell := orderObj.QueryShell().OnlyX(r.Context())
			orderVendorSchema := orderObj.QueryVendorSchema().OnlyX(r.Context())
			orderDetails := OrderDetailStruct{
				ShellId:    orderShell.ID,
				VendorName: orderVendorSchema.Name,
				Status:     orderObj.Status,
				Otp:        orderObj.Otp,
				Items:      orderItemsList,
			}
			if usr.Occupation == "VendorSchema" {
				if orderVendorSchema.ID != usr.QueryVendorSchema().OnlyX(r.Context()).ID {
					errors.ErrorMessage(w, r, 403, "The given order is not handled by requesting VendorSchema", nil, app)
					return
				}
			} else {
				if orderShell.QueryWallet().QueryUser().OnlyX(r.Context()).ID != usr.ID {
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
		if usr.Occupation != "vendor" {
			errors.ErrorMessage(w, r, 403, "Requesting user is not a VendorSchema", nil, app)
			return
		}

		for _, timestamp := range input.DateList {
			var totalEarnings int
			var dayEarnings int
			var orderIdList []int
			orders := usr.QueryVendorSchema().QueryOrders().AllX(r.Context())
			for _, orderObj := range orders {
				if orderObj.Status == helpers.FINISHED {
					totalEarnings += orderObj.Price
					if orderObj.QueryShell().OnlyX(r.Context()).Timestamp.Day() == timestamp.Day() && orderObj.QueryShell().OnlyX(r.Context()).Timestamp.Month() == timestamp.Month() {
						dayEarnings += orderObj.Price
						orderIdList = append(orderIdList, orderObj.ID)
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
		if usr.Occupation != "vendor" {
			usr.Update().SetDisabled(true).SaveX(r.Context())
			errors.ErrorMessage(w, r, 403, "User is not a vendor, Disabled", nil, app)
			return
		}

		orderObj, err := app.Client.Order.Query().Where(order.ID(orderId)).Only(r.Context())
		if err != nil {
			errors.ErrorMessage(w, r, 404, fmt.Sprintf("Order %d not found", orderId), nil, app)
			return
		}
		if !validator.In(input.NewStatus-int(orderObj.Status), 0, 1) {
			errors.ErrorMessage(w, r, 403, "Invalid action", nil, app)
			return
		}

		orderOps := service.NewOrderOps(r.Context(), app.Client)
		_, err, statusCode := orderOps.ChangeStatus(orderObj, helpers.FromInt(input.NewStatus), usr)
		if err != nil {
			errors.ErrorMessage(w, r, statusCode, err.Error(), nil, app)
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
		if usr.Occupation != "vendor" {
			usr.Update().SetDisabled(true).SaveX(r.Context())
			errors.ErrorMessage(w, r, 403, "User is not a vendor, disabled", nil, app)
			return
		}

		orderObj, err := app.Client.Order.Query().Where(order.ID(orderId)).Only(r.Context())
		if err != nil {
			errors.ErrorMessage(w, r, 404, fmt.Sprintf("Order %d not found", orderId), nil, app)
			return
		}
		orderOps := service.NewOrderOps(r.Context(), app.Client)

		err, statusCode := orderOps.Decline(orderObj)
		if err != nil {
			errors.ErrorMessage(w, r, statusCode, err.Error(), nil, app)
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
		if usr.Occupation != "vendor" {
			usr.Update().SetDisabled(true).SaveX(r.Context())
			errors.ErrorMessage(w, r, 403, "Requesting user is not a Vendor", nil, app)
			return
		}

		var input struct {
			ItemObjList []struct {
				ItemId               int  `json:"item_id"`
				NewAvailabilityState bool `json:"new_availability_state"`
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

		vendorOps := service.NewVendorOps(r.Context(), app.Client)
		vendorItemsIdArray := vendorOps.GetVendorArray(usr.QueryVendorSchema().OnlyX(r.Context()))
		for _, itemStruct := range input.ItemObjList {
			itemObject, err := app.Client.Item.Query().Where(item.ID(itemStruct.ItemId)).Only(r.Context())
			if err != nil {
				errors.ErrorMessage(w, r, 404, fmt.Sprintf("Item with ID %d does not exist", itemObject.ID), nil, app)
				return
			}
			if !validator.In(itemObject.ID, vendorItemsIdArray...) {
				usr.Update().SetDisabled(true).SaveX(r.Context())
				errors.ErrorMessage(w, r, 403, "Vendor has been disabled for trying to toggle the availibility of an item not belonging to them", nil, app)
				return
			}
			//if !validator.In(itemStruct.NewAvailabilityState, 0, 1) {
			//	errors.ErrorMessage(w, r, 400, fmt.Sprintf("Invalid valye of new_availability state for item_id %d", itemStruct.ItemId), nil, app)
			//	return
			//}
			var updatedItem outputItem
			updatedItem.ItemId = itemObject.ID
			if !itemStruct.NewAvailabilityState {
				itemObject.Update().SetAvailable(false).SaveX(r.Context())
				updatedItem.Available = false
				availabilityData.Items = append(availabilityData.Items, updatedItem)
			} else if itemStruct.NewAvailabilityState {
				itemObject.Update().SetAvailable(true).SaveX(r.Context())
				updatedItem.Available = true
				availabilityData.Items = append(availabilityData.Items, updatedItem)
			}
		}
		err = response.JSON(w, http.StatusOK, &availabilityData)
		if err != nil {
			errors.ServerError(w, r, err, app)
			return
		}
	}
}

func GetMenu(app *config.Application) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		var vendorId int
		var err error
		if _, ok := vars["vendor_id"]; ok {
			vendorId, err = strconv.Atoi(vars["vendor_id"])
			if err != nil {
				errors.BadRequest(w, r, err, app)
				return
			}
		}
		vendorObject, err := app.Client.VendorSchema.Query().Where(vendor.ID(vendorId)).Only(r.Context())
		if err != nil {
			errors.ErrorMessage(w, r, 404, fmt.Sprintf("Vendor with ID %d does not exist", vendorId), nil, app)
			return
		}
		if vendorObject.QueryUser().OnlyX(r.Context()).Username == "PROF_SHOW" {
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
		for _, itemObj := range vendorObject.QueryItems().AllX(r.Context()) {
			data = append(data, menuItem{
				Id:          itemObj.ID,
				Name:        itemObj.Name,
				Price:       itemObj.BasePrice,
				Description: itemObj.Description,
				VendorId:    itemObj.QueryVendorSchema().OnlyX(r.Context()).ID,
				IsVeg:       itemObj.Veg,
				//IsCombo:     item,
				IsAvailable: itemObj.Available,
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
			ImageUrl    string     `json:"image_url"`
			Description string     `json:"description"`
			Closed      bool       `json:"closed"`
			Menu        []menuItem `json:"menu"`
			Address     string     `json:"address"`
		}

		var data []vendorStruct

		for _, vendorObj := range app.Client.VendorSchema.Query().Where(vendor.HasUserWith(user.UsernameNEQ("PROF_SHOW"))).AllX(r.Context()) {
			if vendorObj.Closed {
				continue
			}

			var menu []menuItem
			for _, itemObj := range vendorObj.QueryItems().AllX(r.Context()) {
				menu = append(menu, menuItem{
					Id:          itemObj.ID,
					Name:        itemObj.Name,
					Price:       itemObj.BasePrice,
					Description: itemObj.Description,
					VendorId:    itemObj.QueryVendorSchema().OnlyX(r.Context()).ID,
					IsVeg:       itemObj.Veg,
					IsAvailable: itemObj.Available,
				})
			}

			data = append(data, vendorStruct{
				ID:          vendorObj.ID,
				Name:        vendorObj.Name,
				ImageUrl:    vendorObj.ImageURL,
				Description: vendorObj.Description,
				Closed:      vendorObj.Closed,
				Menu:        menu,
				Address:     vendorObj.Address,
			})
		}
		err := response.JSON(w, http.StatusOK, &data)
		if err != nil {
			errors.ServerError(w, r, err, app)
			return
		}
	}
}

func GetVendor(app *config.Application) func(http.ResponseWriter, *http.Request) {
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
			ImageUrl    string     `json:"image_url"`
			Description string     `json:"description"`
			Closed      bool       `json:"closed"`
			Menu        []menuItem `json:"menu"`
			Address     string     `json:"address"`
		}
		vars := mux.Vars(r)
		var vendorId int
		var err error
		if _, ok := vars["vendor_id"]; ok {
			vendorId, err = strconv.Atoi(vars["vendor_id"])
			if err != nil {
				errors.BadRequest(w, r, err, app)
				return
			}
		}
		vendorObject, err := app.Client.VendorSchema.Query().Where(vendor.ID(vendorId)).Only(r.Context())
		if err != nil {
			errors.ErrorMessage(w, r, 404, fmt.Sprintf("Vendor with ID %d does not exist", vendorId), nil, app)
			return
		}
		if vendorObject.QueryUser().OnlyX(r.Context()).Username == "PROF_SHOW" {
			errors.ErrorMessage(w, r, 403, "Vendor is a Prof Show", nil, app)
			return
		}
		if vendorObject.Closed {
			errors.ErrorMessage(w, r, 412, "Vendor is closed", nil, app)
			return
		}

		var menu []menuItem
		for _, itemObj := range vendorObject.QueryItems().AllX(r.Context()) {
			menu = append(menu, menuItem{
				Id:          itemObj.ID,
				Name:        itemObj.Name,
				Price:       itemObj.BasePrice,
				Description: itemObj.Description,
				VendorId:    itemObj.QueryVendorSchema().OnlyX(r.Context()).ID,
				IsVeg:       itemObj.Veg,
				IsAvailable: itemObj.Available,
			})
		}

		data := vendorStruct{
			ID:          vendorObject.ID,
			Name:        vendorObject.Name,
			ImageUrl:    vendorObject.ImageURL,
			Description: vendorObject.Description,
			Closed:      vendorObject.Closed,
			Menu:        menu,
			Address:     vendorObject.Address,
		}
		err = response.JSON(w, http.StatusOK, &data)
		if err != nil {
			errors.ServerError(w, r, err, app)
			return
		}
	}
}
