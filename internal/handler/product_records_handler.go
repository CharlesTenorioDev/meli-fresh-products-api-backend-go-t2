package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"

	"github.com/bootcamp-go/web/response"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type ProductRecordsHandler struct {
	service internal.ProductRecordsService
}

func NewProductRecordsHandler(service internal.ProductRecordsService) *ProductRecordsHandler {
	return &ProductRecordsHandler{service: service}
}

func (p *ProductRecordsHandler) GetProductRecords(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	var id int

	if idStr != "" {
		var err error
		id, err = strconv.Atoi(idStr)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid 'id' format")
			return
		}
	}

	products, err := p.service.GetProductRecords(id)
	if err != nil {
		response.Error(w, http.StatusNotFound, utils.ErrNotFound.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data": products,
	})
}

func (p *ProductRecordsHandler) CreateProductRecord(w http.ResponseWriter, r *http.Request) {
	var newProduct internal.ProductRecords
	err := json.NewDecoder(r.Body).Decode(&newProduct)
	if err != nil {
		response.Error(w, http.StatusBadRequest, utils.ErrInvalidFormat.Error())
		return
	}

	product, err := p.service.CreateProductRecord(newProduct)
	if err != nil {
		if errors.Is(err, utils.ErrConflict) {
			response.Error(w, http.StatusConflict, utils.ErrConflict.Error())
			return

		} else if errors.Is(err, utils.ErrInvalidArguments) {
			response.Error(w, http.StatusUnprocessableEntity, utils.ErrInvalidArguments.Error())
			return
		} else {
			response.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	response.JSON(w, http.StatusCreated, map[string]any{
		"data": product,
	})
}
