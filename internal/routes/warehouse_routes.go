package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/handler"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/pkg"
)

func NewWarehouseRoutes(mux *chi.Mux, service pkg.WarehouseService) error {
	warehouseHandler := handler.NewWarehouseHandler(service)
	mux.Route("/api/v1/warehouses", func(router chi.Router) {
		router.Get("/", warehouseHandler.GetAll())
		router.Post("/", warehouseHandler.Post())
		router.Get("/{id}", warehouseHandler.GetById())
		router.Patch("/{id}", warehouseHandler.Update())
		router.Delete("/{id}", warehouseHandler.Delete())
	})
	return nil
}
