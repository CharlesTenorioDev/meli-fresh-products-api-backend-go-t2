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

// CreateLocality handles the creation of a new locality.
// It decodes the request body into a reqPostLocality struct, validates it,
// and then creates a new Locality, Province, and Country based on the provided data.
// If the locality already exists, it returns a 409 Conflict status.
// If the provided arguments are invalid, it returns a 422 Unprocessable Entity status.
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

// GetSellersByLocalityID handles HTTP requests to retrieve sellers by locality ID.
// It extracts the 'id' parameter from the query string, validates it, and calls the service layer to get the sellers.
// If the 'id' parameter is invalid, it responds with a 400 Bad Request status.
// If no sellers are found for the given ID, it responds with a 404 Not Found status.
// For any other errors, it responds with a 500 Internal Server Error status.
// On success, it responds with a 200 OK status and the sellers data in JSON format.
func (h *LocalityHandler) GetSellersByLocalityID() http.HandlerFunc {
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

		locality, err := h.service.GetSellersByLocalityID(id)
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

// GetCarriesByLocalityID handles HTTP requests to retrieve carriers by locality ID.
// It extracts the "id" parameter from the query string, converts it to an integer,
// and calls the service layer to get the carriers associated with the given locality ID.
// If the "id" parameter is missing or invalid, it defaults to 0.
// If an error occurs during the service call or while encoding the response, it returns
// an appropriate HTTP error response.
// The response is returned as a JSON-encoded list of carriers with a status code of 200 OK.
func (handler *LocalityHandler) GetCarriesByLocalityID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")

		idInt, err := strconv.Atoi(id)
		if err != nil {
			idInt = 0
		}

		buyers, err := handler.service.GetCarriesByLocalityID(idInt)
		if err != nil {
			utils.Error(w, http.StatusNotFound, "Locality: "+utils.ErrNotFound.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(buyers); err != nil {
			utils.Error(w, http.StatusInternalServerError, "Failed to encode buyers: "+err.Error())
			return
		}
	}
}
