package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/pkg/logger"

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
		r = logger.GetContext(r, "BUYER:GET_ALL")
		logger.Start(r)

		buyers, err := handler.service.GetAll()
		if err != nil {
			utils.HandleErrorContext(r.Context(), w, err)
			// http.Error(w, "500 Erro Internal api error", http.StatusInternalServerError)
			return
		}

		utils.JSONContext(r.Context(), w, http.StatusOK, buyers)
	}
}

func (handler *BuyerHandler) GetOne() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r = logger.GetContext(r, "BUYER:GET_BY_ID")
		logger.Start(r)

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			utils.HandleErrorContext(r.Context(), w, utils.EBadRequest("id"))
		}

		buyer, err := handler.service.GetOne(id)

		if err != nil {
			utils.HandleErrorContext(r.Context(), w, err)
		}

		utils.JSONContext(r.Context(), w, http.StatusOK, buyer)
	}
}

func (handler *BuyerHandler) CreateBuyer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r = logger.GetContext(r, "BUYER:CREATE")
		logger.Start(r)

		var newBuyer internal.BuyerAttributes
		if err := json.NewDecoder(r.Body).Decode(&newBuyer); err != nil {
			utils.HandleErrorContext(r.Context(), w, utils.EBadRequest("body"))
			return
		}
		logger.Info(r.Context(), "processing", newBuyer)

		buyer, err := handler.service.CreateBuyer(newBuyer)
		if err != nil {
			utils.HandleErrorContext(r.Context(), w, err)

			return
		}

		utils.JSONContext(r.Context(), w, http.StatusCreated, buyer)
	}
}

func (handler *BuyerHandler) UpdateBuyer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r = logger.GetContext(r, "BUYER:UPDATE")
		logger.Start(r)

		var newBuyer internal.BuyerAttributes

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			utils.HandleErrorContext(r.Context(), w, utils.EBadRequest("id"))
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&newBuyer); err != nil {
			utils.HandleErrorContext(r.Context(), w, utils.EBadRequest("body"))
			return
		}

		logger.Info(r.Context(), "processing", newBuyer)

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
			utils.HandleErrorContext(r.Context(), w, err)

			return
		}

		utils.JSONContext(r.Context(), w, http.StatusCreated, buyer)
	}
}

func (handler *BuyerHandler) DeleteBuyer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r = logger.GetContext(r, "BUYER:DELETE")
		logger.Start(r)

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			utils.HandleErrorContext(r.Context(), w, utils.EBadRequest("id"))
			return
		}

		err = handler.service.DeleteBuyer(id)

		if err != nil {
			utils.HandleErrorContext(r.Context(), w, err)
			return
		}

		utils.JSONContext(r.Context(), w, http.StatusNoContent, nil)
	}
}
