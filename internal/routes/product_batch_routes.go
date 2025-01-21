package routes

import (
	"errors"
	"github.com/meli-fresh-products-api-backend-go-t2/cmd/server/handler"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"

	"github.com/go-chi/chi/v5"
)

func ProductBatchRoutes(mux *chi.Mux, service internal.ProductBatchService) error {
	if mux == nil {
		return errors.New("mux router is nil")
	}

	batchHandler := handler.NewProductBatchHandler(service)

	mux.Route("/api/v1/productBatches", func(router chi.Router) {
		router.Post("/", batchHandler.Create())
	})

	return nil
}
