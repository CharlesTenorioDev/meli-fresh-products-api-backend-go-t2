package handler

import (
	"encoding/json"
	"net/http"

	"github.com/meli-fresh-products-api-backend-go-t2/internal/pkg"
)

type BuyerHandler struct {
	service pkg.BuyerService
}

func NewBuyerHandler(service pkg.BuyerService) *BuyerHandler {
	return &BuyerHandler{service: service}
}

func (handler *BuyerHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		buyers, err := handler.service.GetAll()
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
