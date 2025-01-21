package product

import (
	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-go-t2/cmd/server/handler"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
)

func NewProductRoutes(mux *chi.Mux, service internal.ProductService) error {
	productHandler := handler.NewProductHandler(service)

	mux.Route("/api/v1/products", func(router chi.Router) {
		router.Get("/", productHandler.GetProducts)
		router.Post("/", productHandler.CreateProduct)
		router.Get("/{id}", productHandler.GetProductByID)
		router.Patch("/{id}", productHandler.UpdateProduct)
		router.Delete("/{id}", productHandler.DeleteProduct)
	})

	return nil
}
