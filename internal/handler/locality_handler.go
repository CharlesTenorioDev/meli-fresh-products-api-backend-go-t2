package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"
)

type LocalityHandler struct {
	service internal.LocalityService
}

func NewLocalityHandler(service internal.LocalityService) *LocalityHandler {
	return &LocalityHandler{service: service}
}

func (handler *LocalityHandler) GetCarriesByLocalityId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			idInt = 0
		}
		buyers, err := handler.service.GetCarriesByLocalityId(idInt)
		if err != nil {
			http.Error(w, "500 Erro Internal server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(buyers); err != nil {
			http.Error(w, "Failed to encode buyers: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
