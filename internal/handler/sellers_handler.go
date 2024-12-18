package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/pkg"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type SellerRequest struct {
	Cid           int     `json:"cid"`
	CompanyName   string  `json:"company_name"`
	Address       string  `json:"adress"`
	Telephone     string  `json:"telephone"`
}


func NewSellerHandler(service pkg.SellerService) *SellerHandler {
	return &SellerHandler{service}
}

type SellerHandler struct {
	service pkg.SellerService
}

func (h *SellerHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sellers, err := h.service.GetAll()
		if err != nil {
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
		var reqBody pkg.Seller
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			utils.JSON(w, http.StatusInternalServerError, utils.ErrInvalidFormat)
		}

		seller, err := h.service.Create(reqBody)
		if err != nil {
			utils.JSON(w, http.StatusInternalServerError, "internal error")
			return
		}

		utils.JSON(w, http.StatusCreated, seller)


	}

	
}