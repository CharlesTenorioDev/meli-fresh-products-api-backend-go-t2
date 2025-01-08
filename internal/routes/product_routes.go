package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/handler"
)

func NewProductRoutes(mux *chi.Mux, service internal.ProductService) error {
	handler := handler.NewProductHandler(service)
	mux.Route("/api/v1/products", func(router chi.Router) {
		router.Get("/", handler.GetProducts)
		router.Post("/", handler.CreateProduct)
		router.Get("/{id}", handler.GetProductByID)
		router.Patch("/{id}", handler.UpdateProduct)
		router.Delete("/{id}", handler.DeleteProduct)
	})
	return nil
}
