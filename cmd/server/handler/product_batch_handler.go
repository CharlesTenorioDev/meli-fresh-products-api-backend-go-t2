package handler

import (
	"encoding/json"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
	"net/http"
)

type ProductBatchHandler struct {
	service internal.ProductBatchService
}

func NewProductBatchHandler(service internal.ProductBatchService) *ProductBatchHandler {
	return &ProductBatchHandler{service}
}

func (h *ProductBatchHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body internal.ProductBatchRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			utils.Error(w, http.StatusBadRequest, utils.ErrInvalidFormat.Error())
			return
		}

		newBatch, err := h.service.Save(&body)
		if err != nil {
			utils.HandleError(w, err)
			return
		}

		utils.JSON(w, http.StatusCreated, newBatch)
	}
}
