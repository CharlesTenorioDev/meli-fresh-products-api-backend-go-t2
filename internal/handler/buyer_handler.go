package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/pkg"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
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

func (handler *BuyerHandler) GetOne() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			log.Println("Error in parse param to int")
			utils.Error(w, http.StatusBadRequest, err.Error())
		}

		buyer, err := handler.service.GetOne(id)

		if err != nil {
			log.Println("Error to get an user - ", err)
			utils.Error(w, http.StatusInternalServerError, err.Error())
		}

		utils.JSON(w, http.StatusOK, buyer)
	}

}
