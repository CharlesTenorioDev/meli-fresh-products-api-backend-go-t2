package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"

	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

func NewSellerHandler(service internal.SellerService) *SellerHandler {
	return &SellerHandler{service}
}

type SellerHandler struct {
	service internal.SellerService
}

func (h *SellerHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sellers, err := h.service.GetAll()
		if err != nil {
			fmt.Println(err.Error())
			utils.JSON(w, http.StatusInternalServerError, nil)
			return
		}
		utils.JSON(w, http.StatusOK, sellers)
	}
}

func (h *SellerHandler) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "invalid id")
			return
		}

		seller, err := h.service.GetById(id)
		if err != nil {
			if errors.Is(err, utils.ErrNotFound) {
				utils.Error(w, http.StatusNotFound, fmt.Sprintln("id:", id, "not found"))
				return
			}
			utils.Error(w, http.StatusInternalServerError, "Internal error")
			return
		}
		utils.JSON(w, http.StatusOK, seller)
	}
}

func (h *SellerHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqBody internal.SellerRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			utils.JSON(w, http.StatusBadRequest, utils.ErrInvalidFormat)
		}
		newSeller := internal.Seller{
			Cid:         reqBody.Cid,
			CompanyName: reqBody.CompanyName,
			Address:     reqBody.Address,
			Telephone:   reqBody.Telephone,
			LocalityId:  reqBody.LocalityId,
		}
		err := h.service.Create(&newSeller)
		if err != nil {
			fmt.Println(err.Error())
			if errors.Is(err, utils.ErrConflict) {
				utils.Error(w, http.StatusConflict, err.Error())
				return
			}
			if errors.Is(err, utils.ErrInvalidArguments) {
				utils.Error(w, http.StatusUnprocessableEntity, err.Error())
				return
			}

			utils.Error(w, http.StatusInternalServerError, "Internal error")
			return

		}
		utils.JSON(w, http.StatusCreated, newSeller)
	}
}

func (h *SellerHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "invalid id")
			return
		}

		var reqBody internal.SellerRequestPointer

		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			utils.JSON(w, http.StatusBadRequest, utils.ErrInvalidFormat)
		}

		seller, err := h.service.Update(id, reqBody)
		if err != nil {

			if errors.Is(err, utils.ErrConflict) {
				utils.Error(w, http.StatusConflict, err.Error())
				return
			}
			if errors.Is(err, utils.ErrNotFound) {
				utils.Error(w, http.StatusNotFound, err.Error())
				return
			}

			utils.Error(w, http.StatusInternalServerError, "Internal error")
			return

		}
		utils.JSON(w, http.StatusOK, seller)
	}
}

func (h *SellerHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "invalid id")
			return
		}

		err = h.service.Delete(id)
		if err != nil {
			if errors.Is(err, utils.ErrNotFound) {
				utils.Error(w, http.StatusNotFound, err.Error())
				return
			}
			utils.Error(w, http.StatusInternalServerError, "Internal error")
			return
		}
		utils.JSON(w, http.StatusNoContent, nil)
	}
}
