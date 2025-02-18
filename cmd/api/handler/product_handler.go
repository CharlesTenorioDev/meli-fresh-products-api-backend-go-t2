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

type ProductHandler struct {
	service internal.ProductService
}

func NewProductHandler(service internal.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (p *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	r = logger.GetContext(r, "PRODUCT:GET_ALL")
	logger.Start(r)

	products, err := p.service.GetProducts()
	if err != nil {
		utils.HandleErrorContext(r.Context(), w, err)
		return
	}

	utils.JSONContext(r.Context(), w, http.StatusOK, products)
}

func (p *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	r = logger.GetContext(r, "PRODUCT:GET_BY_ID")
	logger.Start(r)

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.HandleErrorContext(r.Context(), w, utils.EBadRequest("ID"))
		return
	}

	product, err := p.service.GetProductByID(id)
	if err != nil {
		utils.HandleErrorContext(r.Context(), w, err)
		return
	}

	utils.JSONContext(r.Context(), w, http.StatusOK, product)
}

func (p *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	r = logger.GetContext(r, "PRODUCT:CREATE")
	logger.Start(r)

	var newProduct internal.ProductAttributes

	err := json.NewDecoder(r.Body).Decode(&newProduct)
	if err != nil {
		utils.HandleErrorContext(r.Context(), w, err)
		return
	}
	logger.Info(r.Context(), "processing", newProduct)

	product, err := p.service.CreateProduct(newProduct)
	if err != nil {
		utils.HandleErrorContext(r.Context(), w, err)
		return
	}

	utils.JSONContext(r.Context(), w, http.StatusCreated, product)
}

func (p *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	r = logger.GetContext(r, "PRODUCT:UPDATE")
	logger.Start(r)

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.HandleErrorContext(r.Context(), w, utils.EBadRequest("ID"))
		return
	}

	var inputProduct internal.Product

	err = json.NewDecoder(r.Body).Decode(&inputProduct)
	if err != nil {
		utils.HandleErrorContext(r.Context(), w, utils.EBadRequest("body"))
		return
	}

	logger.Info(r.Context(), "processing", inputProduct)

	inputProduct.ID = id

	product, err := p.service.UpdateProduct(inputProduct)
	if err != nil {
		utils.HandleErrorContext(r.Context(), w, err)
		return
	}

	utils.JSONContext(r.Context(), w, http.StatusOK, product)
}

func (p *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	r = logger.GetContext(r, "PRODUCT:DELETE")
	logger.Start(r)

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.HandleErrorContext(r.Context(), w, utils.EBadRequest("ID"))
		return
	}

	err = p.service.DeleteProduct(id)
	if err != nil {
		utils.HandleErrorContext(r.Context(), w, err)
		return
	}

	utils.JSONContext(r.Context(), w, http.StatusNoContent, nil)
}
