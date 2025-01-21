package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/handler"
)

func RegisterSellerRoutes(mux *chi.Mux, service internal.SellerService) error {
	sellerHandler := handler.NewSellerHandler(service)

	mux.Route("/api/v1/sellers", func(router chi.Router) {
		router.Get("/", sellerHandler.GetAll())
		router.Get("/{id}", sellerHandler.GetByID())
		router.Post("/", sellerHandler.Create())
		router.Patch("/{id}", sellerHandler.Update())
		router.Delete("/{id}", sellerHandler.Delete())
	})

	return nil
}
