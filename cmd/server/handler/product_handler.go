package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"

	"github.com/bootcamp-go/web/response"
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
	products, err := p.service.GetProducts()
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data": products,
	})
}

func (p *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.HandleError(w, utils.EBadRequest("Invalid ID"))
		return
	}

	product, err := p.service.GetProductByID(id)
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data": product,
	})
}

func (p *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var newProduct internal.ProductAttributes

	err := json.NewDecoder(r.Body).Decode(&newProduct)
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	product, err := p.service.CreateProduct(newProduct)
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, map[string]any{
		"data": product,
	})
}

func (p *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.HandleError(w, utils.EBadRequest("Invalid ID"))
		return
	}

	var inputProduct internal.Product

	err = json.NewDecoder(r.Body).Decode(&inputProduct)
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	inputProduct.ID = id

	product, err := p.service.UpdateProduct(inputProduct)
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data": product,
	})
}

func (p *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.HandleError(w, utils.EBadRequest("Invalid ID"))
		return
	}

	err = p.service.DeleteProduct(id)
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	response.JSON(w, http.StatusNoContent, map[string]any{})
}
