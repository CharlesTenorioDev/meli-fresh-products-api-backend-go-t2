package routes

import (
	"errors"
	"github.com/meli-fresh-products-api-backend-go-t2/cmd/server/handler"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"

	"github.com/go-chi/chi/v5"
)

func BuyerRoutes(mux *chi.Mux, service internal.BuyerService) error {
	if mux == nil {
		return errors.New("mux router is nil")
	}

	buyerHandler := handler.NewBuyerHandler(service)

	mux.Route("/api/v1/buyers", func(router chi.Router) {
		router.Get("/", buyerHandler.GetAll())
		router.Get("/{id}", buyerHandler.GetOne())
		router.Post("/", buyerHandler.CreateBuyer())
		router.Patch("/{id}", buyerHandler.UpdateBuyer())
		router.Delete("/{id}", buyerHandler.DeleteBuyer())
	})

	return nil
}
