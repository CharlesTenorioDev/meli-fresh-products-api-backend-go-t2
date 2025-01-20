package routes

import (
	"errors"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"

	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/handler"
)

func RegisterSectionRoutes(mux *chi.Mux, service internal.SectionService) error {
	if mux == nil {
		return errors.New("mux router is nil")
	}

	sectionHandler := handler.NewSectionHandler(service)

	mux.Route("/api/v1/sections", func(router chi.Router) {
		router.Get("/", sectionHandler.GetAll())
		router.Get("/{id}", sectionHandler.GetById())
		router.Get("/reportProducts", sectionHandler.GetSectionProductsReport())
		router.Post("/", sectionHandler.Post())
		router.Patch("/{id}", sectionHandler.Update())
		router.Delete("/{id}", sectionHandler.Delete())
	})

	return nil
}
