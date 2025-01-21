package inbound_order

import (
	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-go-t2/cmd/server/handler"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
)

// RegisterInboundOrderRoutes is used to register the routes associated with the inboundOrder entity
func RegisterInboundOrderRoutes(mux *chi.Mux, service internal.InboundOrderService) error {
	orderHandler := handler.NewInboundOrderHandler(service)

	// POST /api/v1/inboundOrders
	mux.Route("/api/v1/inboundOrders", func(router chi.Router) {
		router.Post("/", orderHandler.CreateInboundOrder())
	})

	mux.Route("/api/v1/employees/reportInboundOrders", func(router chi.Router) {
		router.Get("/", orderHandler.GetInboundOrdersReport())
	})

	return nil
}
