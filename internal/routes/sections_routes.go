package routes

import (
	"errors"

	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/handler"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/pkg"
)

func RegisterSectionRoutes(mux *chi.Mux, service pkg.SectionService) error {
	if mux == nil {
		return errors.New("mux router is nil")
	}
	handler := handler.NewSectionHandler(service)
	mux.Route("/api/v1/sections", func(router chi.Router) {
		router.Get("/", handler.GetAll())
		router.Get("/{id}", handler.GetById())
		router.Post("/", handler.Post())
		router.Patch("/{id}", handler.Update())
		router.Delete("/{id}", handler.Delete())
	})
	return nil
}