package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/handler"
)

func NewLocalityRoutes(mux *chi.Mux, service internal.LocalityService) error {
	localityHandler := handler.NewLocalityHandler(service)

	mux.Route("/api/v1/localities", func(router chi.Router) {
		router.Get("/reportSellers", localityHandler.GetSellersByLocalityID())
		router.Post("/", localityHandler.CreateLocality())
		router.Get("/reportCarries", localityHandler.GetCarriesByLocalityID())
	})

	return nil
}
