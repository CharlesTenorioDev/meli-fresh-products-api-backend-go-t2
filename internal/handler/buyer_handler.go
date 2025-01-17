package handler

import (
	"encoding/json"
	"errors"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type BuyerHandler struct {
	service internal.BuyerService
}

func NewBuyerHandler(service internal.BuyerService) *BuyerHandler {
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

func (handler *BuyerHandler) CreateBuyer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newBuyer internal.BuyerAttributes
		if err := json.NewDecoder(r.Body).Decode(&newBuyer); err != nil {
			utils.JSON(w, http.StatusInternalServerError, utils.ErrInvalidFormat)
		}

		buyer, err := handler.service.CreateBuyer(newBuyer)
		if err != nil {
			if errors.Is(err, utils.ErrConflict) {
				utils.Error(w, http.StatusConflict, err.Error())
				return
			}
			if errors.Is(err, utils.ErrConflict) {
				utils.Error(w, http.StatusConflict, err.Error())
				return
			}
			utils.Error(w, http.StatusInternalServerError, "500")
			return
		}

		utils.JSON(w, http.StatusCreated, buyer)
	}
}

func (handler *BuyerHandler) UpdateBuyer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newBuyer internal.BuyerAttributes
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			utils.JSON(w, http.StatusBadRequest, utils.ErrInvalidFormat)
		}

		if err := json.NewDecoder(r.Body).Decode(&newBuyer); err != nil {
			utils.JSON(w, http.StatusInternalServerError, utils.ErrInvalidFormat)
		}

		updatedBuyer := internal.Buyer{
			ID: int64(id),
			BuyerAttributes: internal.BuyerAttributes{
				CardNumberID: newBuyer.CardNumberID,
				FirstName:    newBuyer.FirstName,
				LastName:     newBuyer.LastName,
			},
		}

		buyer, err := handler.service.UpdateBuyer(&updatedBuyer)

		if err != nil {
			if errors.Is(err, utils.ErrConflict) {
				utils.Error(w, http.StatusConflict, err.Error())
				return
			}
			if errors.Is(err, utils.ErrConflict) {
				utils.Error(w, http.StatusConflict, err.Error())
				return
			}
			utils.Error(w, http.StatusInternalServerError, "500")
			return
		}

		utils.JSON(w, http.StatusCreated, buyer)
	}
}

func (handler *BuyerHandler) DeleteBuyer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			log.Println("Error in parse param to int")
			utils.Error(w, http.StatusBadRequest, err.Error())
		}

		err = handler.service.DeleteBuyer(id)

		if err != nil {
			log.Println("Error to  an user - ", err)
			utils.Error(w, http.StatusInternalServerError, err.Error())
		}

		utils.JSON(w, http.StatusNoContent, nil)
	}

}
