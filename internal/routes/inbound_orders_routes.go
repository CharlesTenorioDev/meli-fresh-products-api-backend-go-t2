package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/handler"
)

// RegisterInboundOrderRoutes is used to register the routes associated with the inboundOrder entity
func RegisterInboundOrderRoutes(mux *chi.Mux, service internal.InboundOrderService) error {
	handler := handler.NewInboundOrderHandler(service)

	// POST /api/v1/inboundOrders
	mux.Route("/api/v1/inboundOrders", func(router chi.Router) {
		router.Post("/", handler.CreateInboundOrder())
	})

	mux.Route("/api/v1/employees/reportInboundOrders", func(router chi.Router) {
		router.Get("/", handler.GetInboundOrdersReport())
	})

	return nil
}
