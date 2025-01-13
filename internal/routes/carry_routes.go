package routes

import (
	"errors"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"

	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/handler"
)

func CarryRoutes(mux *chi.Mux, service internal.CarryService) error {
	if mux == nil {
		return errors.New("mux router is nil")
	}

	handler := handler.NewCarryHandler(service)
	mux.Route("/api/v1/carries", func(router chi.Router) {
		router.Post("/", handler.SaveCarry())
		router.Get("/", handler.GetAllCarries())
		router.Get("/{id}", handler.GetCarryById())
		router.Patch("/{id}", handler.UpdateCarry())
		router.Delete("/{id}", handler.DeleteCarry())
	})
	return nil
}
