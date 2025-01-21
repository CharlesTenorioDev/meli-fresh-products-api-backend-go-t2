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

// SellerHandler handles HTTP requests for sellers.
//
//	@Summary		Seller Handler
//	@Description	Handles HTTP requests for managing sellers
//	@Tags			sellers
type SellerHandler struct {
	service internal.SellerService
}

// GetAll retrieves all sellers.
//
//	@Summary		Get all sellers
//	@Description	Retrieve a list of all sellers
//	@Tags			sellers
//	@Produce		json
//	@Success		200	{array}		internal.Seller		"List of sellers"
//	@Failure		500	{object}	utils.ErrorResponse	"Internal server error"
//	@Router			/sellers [get]
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

// GetById retrieves a seller by ID.
//
//	@Summary		Get seller by ID
//	@Description	Retrieve a seller by its ID
//	@Tags			sellers
//	@Produce		json
//	@Param			id	path		int					true	"Seller ID"
//	@Success		200	{object}	internal.Seller		"Seller details"
//	@Failure		400	{object}	utils.ErrorResponse	"Invalid ID"
//	@Failure		404	{object}	utils.ErrorResponse	"Seller not found"
//	@Failure		500	{object}	utils.ErrorResponse	"Internal server error"
//	@Router			/api/v1/sellers/{id} [get]
func (h *SellerHandler) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "invalid id")
			return
		}

		seller, err := h.service.GetByID(id)
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

// Create handles the creation of a new seller.
//
//	@Summary		Create a new seller
//	@Description	Create a new seller with the provided details
//	@Tags			sellers
//	@Accept			json
//	@Produce		json
//	@Param			seller	body		internal.SellerRequest	true	"Seller details"
//	@Success		201		{object}	internal.Seller			"Created seller"
//	@Failure		400		{object}	utils.ErrorResponse		"Invalid request format"
//	@Failure		409		{object}	utils.ErrorResponse		"Seller already exists"
//	@Failure		422		{object}	utils.ErrorResponse		"Invalid arguments"
//	@Failure		500		{object}	utils.ErrorResponse		"Internal server error"
//	@Router			/sellers [post]
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
			LocalityID:  reqBody.LocalityID,
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

// Update handles the update of an existing seller.
//
//	@Summary		Update a seller
//	@Description	Update an existing seller with the provided details
//	@Tags			sellers
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int								true	"Seller ID"
//	@Param			seller	body		internal.SellerRequestPointer	true	"Updated seller details"
//	@Success		200		{object}	internal.Seller					"Updated seller"
//	@Failure		400		{object}	utils.ErrorResponse				"Invalid request format"
//	@Failure		404		{object}	utils.ErrorResponse				"Seller not found"
//	@Failure		409		{object}	utils.ErrorResponse				"Seller already exists"
//	@Failure		500		{object}	utils.ErrorResponse				"Internal server error"
//	@Router			/sellers/{id} [put]
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

// Delete handles the deletion of a seller.
//
//	@Summary		Delete a seller
//	@Description	Delete a seller by its ID
//	@Tags			sellers
//	@Produce		json
//	@Param			id	path	int	true	"Seller ID"
//	@Success		204	"No content"
//	@Failure		400	{object}	utils.ErrorResponse	"Invalid ID"
//	@Failure		404	{object}	utils.ErrorResponse	"Seller not found"
//	@Failure		500	{object}	utils.ErrorResponse	"Internal server error"
//	@Router			/sellers/{id} [delete]
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
