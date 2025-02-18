package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/pkg/logger"

	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type ProductTypeHandler struct {
	service internal.ProductTypeService
}

func NewProductTypeHandler(service internal.ProductTypeService) *ProductTypeHandler {
	return &ProductTypeHandler{service: service}
}

func (h *ProductTypeHandler) GetProductTypes(w http.ResponseWriter, r *http.Request) {
	r = logger.GetContext(r, "PRODUCT_TYPE:GET_ALL")
	logger.Start(r)

	productTypes, err := h.service.GetProductTypes()
	if err != nil {
		utils.HandleErrorContext(r.Context(), w, utils.ErrNotFound)
		return
	}

	utils.JSONContext(r.Context(), w, http.StatusOK, productTypes)
}

func (h *ProductTypeHandler) GetProductTypeByID(w http.ResponseWriter, r *http.Request) {
	r = logger.GetContext(r, "PRODUCT_TYPE:GET_BY_ID")
	logger.Start(r)

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.HandleErrorContext(r.Context(), w, utils.EBadRequest("id"))
		return
	}

	productType, err := h.service.GetProductTypeByID(id)
	if err != nil {
		utils.HandleErrorContext(r.Context(), w, err)
		return
	}

	utils.JSONContext(r.Context(), w, http.StatusOK, productType)
}

func (h *ProductTypeHandler) CreateProductType(w http.ResponseWriter, r *http.Request) {
	r = logger.GetContext(r, "PRODUCT_TYPE:CREATE")
	logger.Start(r)

	var newProductType internal.ProductType

	err := json.NewDecoder(r.Body).Decode(&newProductType)
	if err != nil {
		utils.HandleErrorContext(r.Context(), w, utils.EBadRequest("body"))
		return
	}

	logger.Info(r.Context(), "processing", newProductType)

	productType, err := h.service.CreateProductType(newProductType)
	if err != nil {
		utils.HandleErrorContext(r.Context(), w, err)
		return
	}

	utils.JSONContext(r.Context(), w, http.StatusCreated, productType)
}

func (h *ProductTypeHandler) UpdateProductType(w http.ResponseWriter, r *http.Request) {
	r = logger.GetContext(r, "PRODUCT_TYPE:UPDATE")
	logger.Start(r)

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.HandleErrorContext(r.Context(), w, utils.EBadRequest("ID"))
		return
	}

	var inputProductType internal.ProductType

	err = json.NewDecoder(r.Body).Decode(&inputProductType)
	if err != nil {
		utils.HandleErrorContext(r.Context(), w, utils.EBadRequest("body"))
		return
	}
	logger.Info(r.Context(), "processing", inputProductType)

	inputProductType.ID = id

	productType, err := h.service.UpdateProductType(inputProductType)
	if err != nil {
		utils.HandleErrorContext(r.Context(), w, err)
		return
	}

	utils.JSONContext(r.Context(), w, http.StatusOK, productType)
}

func (h *ProductTypeHandler) DeleteProductType(w http.ResponseWriter, r *http.Request) {
	r = logger.GetContext(r, "PRODUCT_TYPE:DELETE")
	logger.Start(r)

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.HandleErrorContext(r.Context(), w, utils.EBadRequest("id"))
		return
	}

	err = h.service.DeleteProductType(id)
	if err != nil {
		utils.HandleErrorContext(r.Context(), w, err)
		return
	}

	utils.JSONContext(r.Context(), w, http.StatusNoContent, nil)
}
