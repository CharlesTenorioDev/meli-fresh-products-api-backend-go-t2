package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/handler"
)

func BuyerRoutes(r *chi.Mux) {
	r.Route("/api/v1/buyer", func(r chi.Router) {
		r.Get("/", handler.GetBuyersHandler)
	})

}
