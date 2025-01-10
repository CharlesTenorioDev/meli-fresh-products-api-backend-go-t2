package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/bootcamp-go/web/response"
	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type CarryHandler struct {
	service internal.CarryService
}

func NewCarryHandler(service internal.CarryService) *CarryHandler {
	return &CarryHandler{service: service}
}

func (handler *CarryHandler) SaveCarry() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var carry *internal.Carry
		err := json.NewDecoder(r.Body).Decode(&carry)
		if err != nil {
			http.Error(w, "Failed to decode carry: "+err.Error(), http.StatusBadRequest)
			return
		}
		if err := handler.service.Save(carry); err != nil {
			if errors.Is(err, utils.ErrConflict) {
				http.Error(w, "CID already exists: "+err.Error(), http.StatusBadRequest)
				return
			}
			if errors.Is(err, utils.ErrInvalidArguments) {
				http.Error(w, "Invalid carry: "+err.Error(), http.StatusUnprocessableEntity)
				return
			}
			if errors.Is(err, utils.ErrNotFound) {
				utils.Error(w, http.StatusNotFound, "Locality: "+err.Error())
				return
			}
			utils.Error(w, http.StatusInternalServerError, "Failed to save carry: "+err.Error())
			return
		}
		utils.JSON(w, http.StatusCreated, carry)
	}
}

func (handler *CarryHandler) GetAllCarries() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		carries, err := handler.service.GetAll()
		if err != nil {
			http.Error(w, "Failed to get all carries", http.StatusNotFound)
			return
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"data": carries,
		})
	}
}

func (handler *CarryHandler) GetCarryById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		carry, err := handler.service.GetById(id)
		if err != nil {
			http.Error(w, "Failed to get carry", http.StatusNotFound)
			return
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"data": carry,
		})
	}
}

func (handler *CarryHandler) UpdateCarry() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		var carry *internal.Carry
		err = json.NewDecoder(r.Body).Decode(carry)
		if err != nil {
			http.Error(w, "Failed to decode carry: "+err.Error(), http.StatusBadRequest)
			return
		}

		carry.ID = id

		if err := handler.service.Update(carry); err != nil {
			http.Error(w, "Failed to update carry: "+err.Error(), http.StatusInternalServerError)
			return
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"data": carry,
		})
	}
}

func (handler *CarryHandler) DeleteCarry() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		if err := handler.service.Delete(id); err != nil {
			http.Error(w, "Failed to delete carry", http.StatusInternalServerError)
			return
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"message": "Carry deleted successfully",
		})
	}
}
