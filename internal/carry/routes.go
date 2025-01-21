package carry

import (
	"errors"
	"github.com/meli-fresh-products-api-backend-go-t2/cmd/server/handler"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"

	"github.com/go-chi/chi/v5"
)

// CarryRoutes sets up the routes for carry-related operations on the provided mux router.
func CarryRoutes(mux *chi.Mux, service internal.CarryService) error {
	if mux == nil {
		return errors.New("mux router is nil")
	}

	carryHandler := handler.NewCarryHandler(service)

	mux.Route("/api/v1/carries", func(router chi.Router) {
		router.Post("/", carryHandler.SaveCarry())
		router.Get("/", carryHandler.GetAllCarries())
		router.Get("/{id}", carryHandler.GetCarryByID())
		router.Patch("/{id}", carryHandler.UpdateCarry())
		router.Delete("/{id}", carryHandler.DeleteCarry())
	})

	return nil
}
