package routes

import (
	"errors"

	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/handler"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/pkg"
)

func BuyerRoutes(mux *chi.Mux, service pkg.BuyerService) error {
	if mux == nil {
		return errors.New("mux router is nil")
	}

	handler := handler.NewBuyerHandler(service)
	mux.Route("/api/v1/buyers", func(router chi.Router) {
		router.Get("/", handler.GetAll())
	})
	return nil
}
