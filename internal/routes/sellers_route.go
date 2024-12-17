package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/handler"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/pkg"
)

func RegisterSellerRoutes(mux *chi.Mux, service pkg.SellerService) error {
	handler := handler.NewSellerHandler(service)
	mux.Route("/api/v1/sellers", func(router chi.Router) {
		router.Get("/", handler.GetAll())
		router.Get("/{id}", handler.GetById())
	})
	return nil
}