package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-go-t2/cmd/server/handler"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
)

func NewProductRecordsRoutes(mux *chi.Mux, service internal.ProductRecordsService) error {
	recordsHandler := handler.NewProductRecordsHandler(service)

	mux.Route("/api/v1/productRecords", func(router chi.Router) {
		router.Post("/", recordsHandler.CreateProductRecord)
	})
	mux.HandleFunc("/api/v1/products/reportRecords", recordsHandler.GetProductRecords)

	return nil
}
