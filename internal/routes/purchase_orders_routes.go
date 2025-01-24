package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/handler"
)

// RegisterPurchaseOrdersRoutes is used to record the routes associated to the purchase orders entity
func RegisterPurchaseOrdersRoutes(mux *chi.Mux, service internal.PurchaseOrderService) error {
	handler := handler.NewPurchaseOrdersHandler(service)
	mux.Route("/api/v1/purchaseOrders", func(router chi.Router) {
		// Post
		router.Post("/", handler.PostPurchaseOrders())
	})
	mux.HandleFunc("/api/v1/buyers/reportPurchaseOrders", handler.GetAllPurchaseOrders())

	return nil
}
