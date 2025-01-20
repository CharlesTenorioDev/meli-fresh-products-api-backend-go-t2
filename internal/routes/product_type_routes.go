package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/handler"
)

func NewProductTypeRoutes(mux *chi.Mux, service internal.ProductTypeService) error {
	productTypeHandler := handler.NewProductTypeHandler(service)

	mux.Route("/api/v1/product_types", func(router chi.Router) {
		router.Get("/", productTypeHandler.GetProductTypes)
		router.Post("/", productTypeHandler.CreateProductType)
		router.Get("/{id}", productTypeHandler.GetProductTypeByID)
		router.Patch("/{id}", productTypeHandler.UpdateProductType)
		router.Delete("/{id}", productTypeHandler.DeleteProductType)
	})

	return nil
}
