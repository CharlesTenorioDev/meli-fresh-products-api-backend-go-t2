package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type InboundOrderHandler struct {
	service internal.InboundOrderService
}

func NewInboundOrderHandler(service internal.InboundOrderService) *InboundOrderHandler {
	return &InboundOrderHandler{service: service}
}

// Handle POST /api/v1/inboundOrders
func (h *InboundOrderHandler) CreateInboundOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newOrder internal.InboundOrderAttributes

		if err := json.NewDecoder(r.Body).Decode(&newOrder); err != nil {
			utils.Error(w, http.StatusBadRequest, "Invalid JSON format")
			return
		}

		order, err := h.service.CreateOrder(newOrder)
		if err != nil {
			if err == utils.ErrConflict {
				utils.Error(w, http.StatusConflict, "Order number already exists")
			} else if err == utils.ErrInvalidArguments {
				utils.Error(w, http.StatusUnprocessableEntity, "Invalid arguments provided")
			} else {
				utils.Error(w, http.StatusInternalServerError, "Failed to create order")
			}
			return
		}

		utils.JSON(w, http.StatusCreated, map[string]any{
			"message": "Order created successfully",
			"data":    order,
		})

	}
}

// Handle GET /api/v1/employees/reportInboundOrders
func (h *InboundOrderHandler) GetInboundOrdersReport() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		employeeID, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "Invalid employee ID")
			return
		}

		report, err := h.service.GetInboundOrdersReport(employeeID)
		if err != nil {
			if err == utils.ErrNotFound {
				utils.Error(w, http.StatusNotFound, "No data found for the given employee ID")
			} else {
				utils.Error(w, http.StatusInternalServerError, "Failed to retrieve report")
			}
			return
		}

		utils.JSON(w, http.StatusOK, map[string]any{
			"message": "Report generated successfully",
			"data":    report,
		})

	}
}
