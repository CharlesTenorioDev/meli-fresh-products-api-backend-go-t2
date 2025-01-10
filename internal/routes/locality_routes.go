package routes

import (
	"errors"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"

	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/handler"
)

func LocalityRoutes(mux *chi.Mux, service internal.LocalityService) error {
	if mux == nil {
		return errors.New("mux router is nil")
	}

	handler := handler.NewLocalityHandler(service)
	mux.Route("/api/v1/localities", func(router chi.Router) {
		router.Get("/reportCarries", handler.GetCarriesByLocalityId())
	})
	return nil
}
