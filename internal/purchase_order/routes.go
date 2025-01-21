package purchase_order

import (
	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-go-t2/cmd/server/handler"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
)

// RegisterPurchaseOrdersRoutes is used to record the routes associated to the purchase orders entity
func RegisterPurchaseOrdersRoutes(mux *chi.Mux, service internal.PurchaseOrderService) error {
	purchaseOrdersHandler := handler.NewPurchaseOrdersHandler(service)

	mux.Route("/api/v1/purchaseOrders", func(router chi.Router) {
		// Post
		router.Post("/", purchaseOrdersHandler.PostPurchaseOrders())
	})
	mux.HandleFunc("/api/v1/buyers/reportPurchaseOrders", purchaseOrdersHandler.GetAllPurchaseOrders())

	return nil
}
