package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
	"github.com/meli-fresh-products-api-backend-go-t2/pkg/logger"
)

type InboundOrderHandler struct {
	service internal.InboundOrderService
}

func NewInboundOrderHandler(service internal.InboundOrderService) *InboundOrderHandler {
	return &InboundOrderHandler{service: service}
}

// CreateInboundOrder Handle POST /api/v1/inboundOrders
func (h *InboundOrderHandler) CreateInboundOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r = logger.GetContext(r, "INBOUND_ORDER:CREATE")
		logger.Start(r)

		var request struct {
			Data internal.InboundOrderAttributes `json:"data"`
		}

		// Decodifica o JSON
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			utils.HandleErrorContext(r.Context(), w, utils.EBadRequest("body"))
			return
		}

		logger.Info(r.Context(), "processing", request)

		newOrder := request.Data

		// Cria a ordem usando o serviço
		order, err := h.service.CreateInboundOrder(newOrder)
		if err != nil {
			utils.HandleErrorContext(r.Context(), w, err)

			return
		}

		// Retorna o JSON com apenas os dados diretamente na resposta
		utils.JSONContext(r.Context(), w, http.StatusCreated, map[string]any{
			"order_date":       order.Attributes.OrderDate,
			"order_number":     order.Attributes.OrderNumber,
			"employee_id":      order.Attributes.EmployeeID,
			"product_batch_id": order.Attributes.ProductBatchID,
			"warehouse_id":     order.Attributes.WarehouseID,
		})
	}
}

// GenerateInboundOrdersReport sHandle GET /api/v1/employees/reportInboundOrders
//
//	handles the generation of the report.
func (h *InboundOrderHandler) GenerateInboundOrdersReport() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r = logger.GetContext(r, "INBOUND_ORDER:INBOUND_ORDERS_REPORT")
		logger.Start(r)

		idsParam := r.URL.Query().Get("id")

		var ids []int

		if idsParam != "" {
			idsStrings := strings.Split(idsParam, ",")
			for _, idStr := range idsStrings {
				id, err := strconv.Atoi(strings.TrimSpace(idStr))
				if err != nil {
					utils.HandleErrorContext(r.Context(), w, utils.EBadRequest("id"))
					return
				}

				ids = append(ids, id)
			}
		}

		report, err := h.service.GenerateInboundOrdersReport(ids)
		if err != nil {
			utils.HandleErrorContext(r.Context(), w, err)

			return
		}

		utils.JSONContext(r.Context(), w, http.StatusOK, report)
	}
}
