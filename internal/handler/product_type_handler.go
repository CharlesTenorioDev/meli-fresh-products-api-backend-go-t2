package handler

import (
	"encoding/json"
	"errors"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"net/http"
	"strconv"

	"github.com/bootcamp-go/web/response"
	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type ProductTypeHandler struct {
	service internal.ProductTypeService
}

func NewProductTypeHandler(service internal.ProductTypeService) *ProductTypeHandler {
	return &ProductTypeHandler{service: service}
}

func (h *ProductTypeHandler) GetProductTypes(w http.ResponseWriter, _ *http.Request) {
	productTypes, err := h.service.GetProductTypes()
	if err != nil {
		response.Error(w, http.StatusNotFound, utils.ErrNotFound.Error())
		return
	}
	response.JSON(w, http.StatusOK, map[string]interface{}{
		"data": productTypes,
	})

}

func (h *ProductTypeHandler) GetProductTypeByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, utils.ErrInvalidFormat.Error())
		return
	}
	productType, err := h.service.GetProductTypeByID(id)
	if err != nil {
		response.Error(w, http.StatusNotFound, utils.ErrNotFound.Error())
		return
	}
	response.JSON(w, http.StatusOK, map[string]interface{}{
		"data": productType,
	})

}

func (h *ProductTypeHandler) CreateProductType(w http.ResponseWriter, r *http.Request) {
	var newProductType internal.ProductType
	err := json.NewDecoder(r.Body).Decode(&newProductType)
	if err != nil {
		response.Error(w, http.StatusBadRequest, utils.ErrInvalidFormat.Error())
		return
	}
	productType, err := h.service.CreateProductType(newProductType)
	if err != nil {
		if errors.Is(err, utils.ErrConflict) {
			response.Error(w, http.StatusConflict, utils.ErrConflict.Error())
			return
		} else {
			response.Error(w, http.StatusUnprocessableEntity, utils.ErrInvalidArguments.Error())
			return
		}
	}
	response.JSON(w, http.StatusCreated, map[string]interface{}{
		"data": productType,
	})

}

func (h *ProductTypeHandler) UpdateProductType(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, utils.ErrInvalidFormat.Error())
		return
	}
	var inputProductType internal.ProductType
	err = json.NewDecoder(r.Body).Decode(&inputProductType)
	if err != nil {
		response.Error(w, http.StatusBadRequest, utils.ErrInvalidFormat.Error())
		return
	}
	inputProductType.ID = id
	productType, err := h.service.UpdateProductType(inputProductType)
	if err != nil {
		response.Error(w, http.StatusNotFound, utils.ErrNotFound.Error())
		return
	}
	response.JSON(w, http.StatusOK, map[string]interface{}{
		"data": productType,
	})

}

func (h *ProductTypeHandler) DeleteProductType(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, utils.ErrInvalidFormat.Error())
		return
	}
	err = h.service.DeleteProductType(id)
	if err != nil {
		response.Error(w, http.StatusNotFound, utils.ErrNotFound.Error())
		return
	}
	response.JSON(w, http.StatusNoContent, nil)

}
