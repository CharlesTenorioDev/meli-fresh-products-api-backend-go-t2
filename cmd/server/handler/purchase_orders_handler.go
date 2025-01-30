package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"

	"github.com/bootcamp-go/web/response"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

// PurchaseOrderDefault is the http handler for PurchaseOrder-related endpoints
// it communicates with the service layer to process requests
type PurchaseOrderDefault struct {
	sv internal.PurchaseOrderService
}

// NewPurchaseOrdersHandler creates a new instance of PurchaseOrderDefault
func NewPurchaseOrdersHandler(sv internal.PurchaseOrderService) *PurchaseOrderDefault {
	return &PurchaseOrderDefault{sv: sv}
}

func (h *PurchaseOrderDefault) GetAllPurchaseOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		buyerIDParam := queryParams.Get("id")

		var buyerID int

		if buyerIDParam != "" {
			var err error

			buyerID, err = strconv.Atoi(buyerIDParam)
			if err != nil {
				utils.HandleError(w, utils.ErrInvalidFormat)
				return
			}
		}

		PurchaseOrdersSummary, err := h.sv.FindAllByBuyerID(buyerID)
		if err != nil {
			utils.HandleError(w, err)

			return
		}

		data := make(map[int]map[string]any)
		for _, value := range PurchaseOrdersSummary {
			data[value.BuyerID] = map[string]any{
				"total_orders": value.TotalOrders,
				"order_codes":  value.OrderCodes,
			}
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

// PostPurchaseOrders handles the POST /PurchaseOrders route
func (h *PurchaseOrderDefault) PostPurchaseOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newPurchaseOrder internal.PurchaseOrderAttributes

		// decode the json request body
		err := json.NewDecoder(r.Body).Decode(&newPurchaseOrder)
		if err != nil {
			utils.HandleError(w, utils.ErrInvalidFormat)
			return
		}

		// create the PurchaseOrder
		PurchaseOrder, err := h.sv.CreatePurchaseOrder(newPurchaseOrder)
		if err != nil {
			utils.HandleError(w, err)

			return
		}

		// returns status 201 and the data if all ok
		response.JSON(w, http.StatusCreated, map[string]any{
			"message": "success",
			"data":    PurchaseOrder,
		})
	}
}
