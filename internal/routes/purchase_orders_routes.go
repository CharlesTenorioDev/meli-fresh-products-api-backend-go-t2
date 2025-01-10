package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/handler"
)

// RegisterPurchaseOrdersRoutes is used to record the routes associated to the purchase orders entity
func RegisterPurchaseOrdersRoutes(mux *chi.Mux, service internal.PurchaseOrdersService) error {
	handler := handler.NewEmployeeHandler(service)
	mux.Route("/api/v1/PurchaseOrders", func(router chi.Router) {
		// Get
		router.Get("/", handler.GetAllPurchaseOrders())
		router.Get("/{id}", handler.GetPurchaseOrdersById())
		// Post
		router.Post("/", handler.PostPurchaseOrders())
		// Patch
		router.Patch("/{id}", handler.PatchPurchaseOrders())
		// Delete
		router.Delete("/{id}", handler.DeletePurchaseOrders())
	})

	return nil
}
