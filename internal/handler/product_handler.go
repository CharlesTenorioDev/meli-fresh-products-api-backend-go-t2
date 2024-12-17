package handler

import (
	"net/http"

	"github.com/bootcamp-go/web/response"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/pkg"
)

type ProductHandler struct {
	service pkg.ProductService
}

func NewProductHandler(service pkg.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (p *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := p.service.GetProducts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.JSON(w, http.StatusOK, products)

}
