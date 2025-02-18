package handler

import (
	"encoding/json"
	"net/http"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
	"github.com/meli-fresh-products-api-backend-go-t2/pkg/logger"
)

type ProductBatchHandler struct {
	service internal.ProductBatchService
}

func NewProductBatchHandler(service internal.ProductBatchService) *ProductBatchHandler {
	return &ProductBatchHandler{service}
}

func (h *ProductBatchHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r = logger.GetContext(r, "PRODUCT_BATCH:CREATE")
		logger.Start(r)

		var body internal.ProductBatchRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			utils.HandleErrorContext(r.Context(), w, utils.EBadRequest("body"))
			return
		}
		logger.Info(r.Context(), "processing", body)

		newBatch, err := h.service.Save(&body)
		if err != nil {
			utils.HandleErrorContext(r.Context(), w, err)
			return
		}

		utils.JSONContext(r.Context(), w, http.StatusCreated, newBatch)
	}
}
