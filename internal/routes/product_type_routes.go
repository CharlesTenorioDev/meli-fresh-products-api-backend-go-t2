package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/handler"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/pkg"
)

func NewProductTypeRoutes(mux *chi.Mux, service pkg.ProductTypeService) error {
	handler := handler.NewProductTypeHandler(service)
	mux.Route("/api/v1/product_types", func(router chi.Router) {
		router.Get("/", handler.GetProductTypes)
		router.Post("/", handler.CreateProductType)
		router.Get("/{id}", handler.GetProductTypeByID)
		router.Patch("/{id}", handler.UpdateProductType)
		router.Delete("/{id}", handler.DeleteProductType)
	})
	return nil
}
