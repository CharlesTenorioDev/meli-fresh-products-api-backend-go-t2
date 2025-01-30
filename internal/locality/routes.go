package locality

import (
	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-go-t2/cmd/server/handler"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
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
