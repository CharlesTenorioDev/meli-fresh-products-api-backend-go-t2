package warehouse

import (
	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-go-t2/cmd/server/handler"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
)

// NewWarehouseRoutes sets up the routes for warehouse-related endpoints.
// It registers the following routes:
// - GET /api/v1/warehouses: Retrieves all warehouses.
// - POST /api/v1/warehouses: Creates a new warehouse.
// - GET /api/v1/warehouses/{id}: Retrieves a warehouse by its ID.
// - PATCH /api/v1/warehouses/{id}: Updates a warehouse by its ID.
// - DELETE /api/v1/warehouses/{id}: Deletes a warehouse by its ID.
//
// Parameters:
// - mux: The HTTP request multiplexer from the chi package.
// - service: The service layer that provides warehouse-related operations.
//
// Returns:
// - An error if there is an issue setting up the routes, otherwise nil.

func NewWarehouseRoutes(mux *chi.Mux, service internal.WarehouseService) error {
	warehouseHandler := handler.NewWarehouseHandler(service)

	mux.Route("/api/v1/warehouses", func(router chi.Router) {
		router.Get("/", warehouseHandler.GetAll())
		router.Post("/", warehouseHandler.Post())
		router.Get("/{id}", warehouseHandler.GetByID())
		router.Patch("/{id}", warehouseHandler.Update())
		router.Delete("/{id}", warehouseHandler.Delete())
	})

	return nil
}
