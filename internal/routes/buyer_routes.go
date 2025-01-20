package routes

import (
	"errors"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"

	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/handler"
)

func BuyerRoutes(mux *chi.Mux, service internal.BuyerService) error {
	if mux == nil {
		return errors.New("mux router is nil")
	}

	handler := handler.NewBuyerHandler(service)

	mux.Route("/api/v1/buyers", func(router chi.Router) {
		router.Get("/", handler.GetAll())
		router.Get("/{id}", handler.GetOne())
		router.Post("/", handler.CreateBuyer())
		router.Patch("/{id}", handler.UpdateBuyer())
		router.Delete("/{id}", handler.DeleteBuyer())
	})

	return nil
}
