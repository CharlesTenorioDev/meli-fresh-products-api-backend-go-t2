package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type LocalityHandler struct {
	service internal.LocalityService
}

func NewLocalityHandler(service internal.LocalityService) *LocalityHandler {
	return &LocalityHandler{
		service: service,
	}
}

type reqPostLocality struct {
	Data struct {
		ID           int    `json:"id"`
		LocalityName string `json:"locality_name"`
		ProvinceName string `json:"province_name"`
		CountryName  string `json:"country_name"`
	} `json:"data"`
}

func (h *LocalityHandler) CreateLocality() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body reqPostLocality
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			utils.Error(w, http.StatusBadRequest, utils.ErrInvalidFormat.Error())
			return
		}
		newLocality := internal.Locality{
			ID:           body.Data.ID,
			LocalityName: body.Data.LocalityName,
		}
		province := internal.Province{
			ProvinceName: body.Data.ProvinceName,
		}
		country := internal.Country{
			CountryName: body.Data.CountryName,
		}
		err := h.service.Save(&newLocality, &province, &country)
		if err != nil {
			if errors.Is(err, utils.ErrConflict) {
				utils.Error(w, http.StatusConflict, err.Error())
				return
			}
			if errors.Is(err, utils.ErrInvalidArguments) {
				utils.Error(w, http.StatusUnprocessableEntity, err.Error())
				return
			}
			utils.Error(w, http.StatusInternalServerError, "Some error occurs")
			return
		}
		utils.JSON(w, http.StatusCreated, newLocality)
	}
}

func (h *LocalityHandler) GetSellersByLocalityId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := 0
		var err error
		if strings.TrimSpace(r.URL.Query().Get("id")) != "" {
			id, err = strconv.Atoi(r.URL.Query().Get("id"))
			if err != nil {
				utils.Error(w, http.StatusBadRequest, "invalid id")
				return
			}
		}
		locality, err := h.service.GetSellersByLocalityId(id)
		if err != nil {
			if errors.Is(err, utils.ErrNotFound) {
				utils.Error(w, http.StatusNotFound, fmt.Sprintf("no locality for id %d", id))
				return
			}
			utils.Error(w, http.StatusInternalServerError, "Some error occurs")
			return
		}
		utils.JSON(w, http.StatusOK, locality)
	}
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
