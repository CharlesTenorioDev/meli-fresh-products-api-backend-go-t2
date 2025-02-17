package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

func RegisterHealhcheck(mux *chi.Mux) {
	mux.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		utils.JSON(w, 200, "pong")
	})
}
