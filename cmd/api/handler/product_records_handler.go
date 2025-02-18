package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/pkg/logger"

	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type ProductRecordsHandler struct {
	service internal.ProductRecordsService
}

func NewProductRecordsHandler(service internal.ProductRecordsService) *ProductRecordsHandler {
	return &ProductRecordsHandler{service: service}
}

func (p *ProductRecordsHandler) GetProductRecords(w http.ResponseWriter, r *http.Request) {
	r = logger.GetContext(r, "PRODUCT_RECORD:GET_BY_ID")
	logger.Start(r)

	idStr := r.URL.Query().Get("id")

	var id int

	if idStr != "" {
		var err error

		id, err = strconv.Atoi(idStr)
		if err != nil {
			utils.HandleErrorContext(r.Context(), w, utils.EBadRequest("id"))
			return
		}
	}

	products, err := p.service.GetProductRecords(id)
	if err != nil {
		utils.HandleErrorContext(r.Context(), w, err)
		return
	}

	utils.JSONContext(r.Context(), w, http.StatusOK, products)
}

func (p *ProductRecordsHandler) CreateProductRecord(w http.ResponseWriter, r *http.Request) {
	r = logger.GetContext(r, "PRODUCT_RECORD:CREATE")
	logger.Start(r)

	var newProduct internal.ProductRecords

	err := json.NewDecoder(r.Body).Decode(&newProduct)
	if err != nil {
		utils.HandleErrorContext(r.Context(), w, utils.EBadRequest("body"))
		return
	}
	logger.Info(r.Context(), "processing", newProduct)

	product, err := p.service.CreateProductRecord(newProduct)
	if err != nil {
		utils.HandleErrorContext(r.Context(), w, err)
		return
	}

	utils.JSONContext(r.Context(), w, http.StatusCreated, product)
}
