package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"

	"github.com/bootcamp-go/web/response"
	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

// PurchaseOrderDefault is the http handler for PurchaseOrder-related endpoints
// it communicates with the service layer to process requests
type PurchaseOrderDefault struct {
	sv internal.PurchaseOrderService
}

// NewPurchaseOrderHandler creates a new instance of PurchaseOrderDefault
func NewPurchaseOrdersHandler(sv internal.PurchaseOrderService) *PurchaseOrderDefault {
	return &PurchaseOrderDefault{sv: sv}
}

// GetAllPurchaseOrders handles the GET /PurchaseOrders route
func (h *PurchaseOrderDefault) GetAllPurchaseOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		PurchaseOrders, err := h.sv.FindAll()
		if err != nil {
			response.JSON(w, http.StatusInternalServerError, nil)
			return
		}

		data := make(map[int]internal.PurchaseOrder)
		for key, value := range PurchaseOrders {
			data[key] = internal.PurchaseOrder{
				ID: value.ID,
				Attributes: internal.PurchaseOrderAttributes{
					OrderNumber:     value.Attributes.OrderNumber,
					OrderDate:       value.Attributes.OrderDate,
					TrackingCode:    value.Attributes.TrackingCode,
					BuyerId:         value.Attributes.BuyerId,
					ProductRecordId: value.Attributes.ProductRecordId,
				},
			}
		}

		// returns status 200 and the data if all ok
		response.JSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

// GetPurchaseOrdersById handles the GET /PurchaseOrders/{id} route
func (h *PurchaseOrderDefault) GetPurchaseOrdersById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			handleError(w, utils.ErrInvalidFormat)
			return
		}

		PurchaseOrder, err := h.sv.FindById(id)
		if err != nil {
			handleError(w, utils.ErrNotFound)
			return
		}

		data := internal.PurchaseOrder{
			ID: PurchaseOrder.ID,
			Attributes: internal.PurchaseOrderAttributes{
				OrderNumber:     PurchaseOrder.Attributes.OrderNumber,
				OrderDate:       PurchaseOrder.Attributes.OrderDate,
				TrackingCode:    PurchaseOrder.Attributes.TrackingCode,
				BuyerId:         PurchaseOrder.Attributes.BuyerId,
				ProductRecordId: PurchaseOrder.Attributes.ProductRecordId,
			},
		}

		// returns status 200 and the data if all ok
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
			handleError(w, utils.ErrInvalidFormat)
			return
		}

		// create the PurchaseOrder
		PurchaseOrder, err := h.sv.CreatePurchaseOrder(newPurchaseOrder)
		if err != nil {
			if err == utils.ErrConflict {
				handleError(w, utils.ErrConflict)
			} else if err == utils.ErrEmptyArguments {
				handleError(w, utils.ErrEmptyArguments)
			} else if err == utils.ErrWarehouseDoesNotExists {
				handleError(w, utils.ErrWarehouseDoesNotExists)
			} else {
				handleError(w, utils.ErrInvalidArguments)
			}
			return
		}

		// returns status 201 and the data if all ok
		response.JSON(w, http.StatusCreated, map[string]any{
			"message": "success",
			"data":    PurchaseOrder,
		})
	}
}

// PatchPurchaseOrders handles the PATCH /PurchaseOrders route
func (h *PurchaseOrderDefault) PatchPurchaseOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			handleError(w, utils.ErrInvalidFormat)
			return
		}

		var inputPurchaseOrder internal.PurchaseOrder
		// decode the json request body into PurchaseOrder struct
		err = json.NewDecoder(r.Body).Decode(&inputPurchaseOrder)
		if err != nil {
			handleError(w, utils.ErrInvalidFormat)
			return
		}

		inputPurchaseOrder.ID = id
		// update the PurchaseOrder
		PurchaseOrder, err := h.sv.UpdatePurchaseOrder(inputPurchaseOrder)
		if err != nil {
			if err == utils.ErrNotFound {
				handleError(w, utils.ErrNotFound)
			} else {
				handleError(w, utils.ErrWarehouseDoesNotExists)
			}
			return
		}

		// returns status 200 and the data if all ok
		response.JSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    PurchaseOrder,
		})
	}
}

// DeletePurchaseOrders handles the DELETE /PurchaseOrders/{id} route
// it deletes an existing PurchaseOrder based on the provided ID
func (h *PurchaseOrderDefault) DeletePurchaseOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// extract the PurchaseOrder ID from the URL parameters and converts it to int
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			handleError(w, utils.ErrInvalidFormat)
			return
		}

		// delete the PurchaseOrder
		err = h.sv.DeletePurchaseOrder(id)
		if err != nil {
			if err == utils.ErrNotFound {
				handleError(w, utils.ErrNotFound)
			} else {
				handleError(w, utils.ErrInvalidArguments)
			}
			return
		}

		// returns status 204 and a success message if all ok
		response.JSON(w, http.StatusNoContent, map[string]any{
			"message": "PurchaseOrder deleted successfully",
		})
	}
}
