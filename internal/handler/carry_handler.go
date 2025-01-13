package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

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
				utils.Error(w, http.StatusConflict, "CID already exists: "+err.Error())
				return
			}
			if errors.Is(err, utils.ErrInvalidArguments) {
				utils.Error(w, http.StatusUnprocessableEntity, "Invalid carry: "+err.Error())
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
			utils.Error(w, http.StatusNotFound, "Failed to get all carries")
			return
		}

		utils.JSON(w, http.StatusOK, carries)
	}
}

func (handler *CarryHandler) GetCarryById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "Invalid ID")
			return
		}

		carry, err := handler.service.GetById(id)
		if err != nil {
			utils.Error(w, http.StatusNotFound, "Failed to get carry")
			return
		}

		utils.JSON(w, http.StatusOK, carry)
	}
}

func (handler *CarryHandler) UpdateCarry() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "Invalid ID")
			return
		}

		carry := &internal.Carry{}
		err = json.NewDecoder(r.Body).Decode(carry)
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "Failed to decode carry: "+err.Error())
			return
		}

		carry.ID = id

		if err := handler.service.Update(carry); err != nil {
			if errors.Is(err, utils.ErrConflict) {
				utils.Error(w, http.StatusConflict, err.Error())
				return
			}
			if errors.Is(err, utils.ErrNotFound) {
				utils.Error(w, http.StatusNotFound, err.Error())
				return
			}
			utils.Error(w, http.StatusInternalServerError, "Failed to update carry: "+err.Error())
			return
		}

		utils.JSON(w, http.StatusOK, carry)
	}
}

func (handler *CarryHandler) DeleteCarry() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "Invalid ID")
			return
		}

		if err := handler.service.Delete(id); err != nil {
			if errors.Is(err, utils.ErrNotFound) {
				utils.Error(w, http.StatusNotFound, "Carry not found")
				return
			}
			utils.Error(w, http.StatusInternalServerError, "Failed to delete carry")
			return
		}

		utils.JSON(w, http.StatusNoContent, "Carry deleted successfully")
	}
}
