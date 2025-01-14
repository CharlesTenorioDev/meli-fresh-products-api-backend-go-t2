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
		var request struct {
			Data internal.InboundOrderAttributes `json:"data"`
		}

		// Decodifica o JSON
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			utils.Error(w, http.StatusBadRequest, "Invalid JSON format")
			return
		}

		newOrder := request.Data

		// Valida os campos obrigatórios
		if newOrder.OrderDate == "" || newOrder.OrderNumber == "" || newOrder.EmployeeID == 0 || newOrder.ProductBatchID == 0 || newOrder.WarehouseID == 0 {
			utils.Error(w, http.StatusUnprocessableEntity, "Missing required fields")
			return
		}

		// Cria a ordem usando o serviço
		order, err := h.service.CreateOrder(newOrder)
		if err != nil {
			if err == utils.ErrConflict {
				utils.Error(w, http.StatusConflict, "Order number already exists or employee ID is invalid")
			} else {
				utils.Error(w, http.StatusInternalServerError, "Failed to create order")
			}
			return
		}

		// Retorna o JSON com apenas os dados diretamente na resposta
		utils.JSON(w, http.StatusCreated, map[string]any{
			"message":          "Order created successfully",
			"order_date":       order.Attributes.OrderDate,
			"order_number":     order.Attributes.OrderNumber,
			"employee_id":      order.Attributes.EmployeeID,
			"product_batch_id": order.Attributes.ProductBatchID,
			"warehouse_id":     order.Attributes.WarehouseID,
		})
	}
}

// Handle GET /api/v1/employees/reportInboundOrders
func (h *InboundOrderHandler) GetInboundOrdersReport() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		employeeIDStr := r.URL.Query().Get("id")
		if employeeIDStr == "" {
			utils.Error(w, http.StatusBadRequest, "Missing 'id' parameter")
			return
		}

		employeeID, err := strconv.Atoi(employeeIDStr)
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "Invalid 'id' parameter")
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
