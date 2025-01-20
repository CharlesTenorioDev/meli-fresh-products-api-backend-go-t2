package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/handler"
)

func RegisterSellerRoutes(mux *chi.Mux, service internal.SellerService) error {
	handler := handler.NewSellerHandler(service)

	mux.Route("/api/v1/sellers", func(router chi.Router) {
		router.Get("/", handler.GetAll())
		router.Get("/{id}", handler.GetById())
		router.Post("/", handler.Create())
		router.Patch("/{id}", handler.Update())
		router.Delete("/{id}", handler.Delete())
	})

	return nil
}
