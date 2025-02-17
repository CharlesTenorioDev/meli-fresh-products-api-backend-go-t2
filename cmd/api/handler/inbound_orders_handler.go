package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
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
		var request struct {
			Data internal.InboundOrderAttributes `json:"data"`
		}

		// Decodifica o JSON
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			utils.Error(w, http.StatusBadRequest, utils.ErrInvalidFormat.Error())
			return
		}

		newOrder := request.Data

		// Cria a ordem usando o serviço
		order, err := h.service.CreateInboundOrder(newOrder)
		if err != nil {
			if errors.Is(err, utils.ErrConflict) {
				utils.Error(w, http.StatusConflict, err.Error())
				return
			}

			if errors.Is(err, utils.ErrNotFound) {
				log.Println(newOrder)
				utils.Error(w, http.StatusNotFound, err.Error())
				return
			}

			if errors.Is(err, utils.ErrInvalidArguments) {
				utils.Error(w, http.StatusUnprocessableEntity, err.Error())
				return
			}

			utils.Error(w, http.StatusInternalServerError, "Failed to create order")

			return
		}

		// Retorna o JSON com apenas os dados diretamente na resposta
		utils.JSON(w, http.StatusCreated, map[string]any{
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
		idsParam := r.URL.Query().Get("id")

		var ids []int

		if idsParam != "" {
			idsStrings := strings.Split(idsParam, ",")
			for _, idStr := range idsStrings {
				id, err := strconv.Atoi(strings.TrimSpace(idStr))
				if err != nil {
					utils.HandleError(w, utils.ErrInvalidFormat)
					return
				}

				ids = append(ids, id)
			}
		}

		report, err := h.service.GenerateInboundOrdersReport(ids)
		if err != nil {
			if errors.Is(err, utils.ErrConflict) {
				utils.Error(w, http.StatusConflict, err.Error())
				return
			}

			if errors.Is(err, utils.ErrNotFound) {
				utils.Error(w, http.StatusNotFound, err.Error())
				return
			}

			if errors.Is(err, utils.ErrInvalidArguments) {
				utils.Error(w, http.StatusUnprocessableEntity, err.Error())
				return
			}

			utils.Error(w, http.StatusInternalServerError, "Failed to create order")

			return
		}

		utils.JSON(w, http.StatusOK, report)
	}
}
