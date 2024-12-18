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

// Simple structure to hold the data when POST request
type reqPostSection struct {
	SectionNumber      int     `json:"section_number"`
	CurrentCapacity    int     `json:"current_capacity"`
	MaximumCapacity    int     `json:"maximum_capacity"`
	MinimumCapacity    int     `json:"minimum_capacity"`
	CurrentTemperature float64 `json:"current_temperature"`
	MinimumTemperature float64 `json:"minimum_temperature"`
	ProductTypeID      int     `json:"product_type_id"`
	WarehouseID        int     `json:"warehouse_id"`
}

type SectionHandler struct {
	service pkg.SectionService
}

// Get a new instance of SectionHandler
func NewSectionHandler(service pkg.SectionService) *SectionHandler {
	return &SectionHandler{service}
}

// Get all the sections - 200
// An error not mapped - 500
func (h *SectionHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sections, err := h.service.GetAll()
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, "Some error occurs")
			return
		}
		utils.JSON(w, http.StatusOK, sections)
	}
}

// Get the section by Id - 200
// If the id is in the wrong format - 400
// If the section doesn't exist - 404
// An error not mapped - 500
func (h *SectionHandler) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "invalid id")
			return
		}
		section, err := h.service.GetById(id)
		if err != nil {
			if errors.Is(err, utils.ErrNotFound) {
				utils.Error(w, http.StatusNotFound, fmt.Sprintf("no section for id %d", id))
				return
			}
			utils.Error(w, http.StatusInternalServerError, "Some error occurs")
			return
		}
		utils.JSON(w, http.StatusOK, section)
	}
}

// Save the section - 201
// If payload is in the wrong format - 400
// If a section already exists for section_number - 409
// If the payload contains invalid or empty fields for mandatory data - 422
// An error not mapped - 500
func (h *SectionHandler) Post() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body reqPostSection
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			utils.Error(w, http.StatusBadRequest, utils.ErrInvalidFormat.Error())
			return
		}
		newSection := pkg.Section{
			SectionNumber:      body.SectionNumber,
			CurrentCapacity:    body.CurrentCapacity,
			MaximumCapacity:    body.MaximumCapacity,
			MinimumCapacity:    body.MinimumCapacity,
			CurrentTemperature: body.CurrentTemperature,
			MinimumTemperature: body.MinimumTemperature,
			ProductTypeID:      body.ProductTypeID,
			WarehouseID:        body.WarehouseID,
		}
		newSection, err := h.service.Save(newSection)
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
		utils.JSON(w, http.StatusCreated, newSection)
	}
}

// Update the section - 200
// If payload or Id is in a incorrect format - 400
// If a section already exists for section_number - 409
// If the payload contains invalid or empty fields for mandatory data - 422
// An error not mapped - 500
func (h *SectionHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "invalid id")
			return
		}
		var body pkg.SectionPointers
		if err = json.NewDecoder(r.Body).Decode(&body); err != nil {
			utils.Error(w, http.StatusBadRequest, utils.ErrInvalidFormat.Error())
			return
		}
		updatedSection, err := h.service.Update(id, body)
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
		utils.JSON(w, http.StatusOK, updatedSection)
	}
}

// Delete a section by id - 204
// If the id is in the wrong format - 400
// If the section doesn't exist - 404
// An error not mapped - 500
func (h *SectionHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "invalid id")
			return
		}
		err = h.service.Delete(id)
		if err != nil {
			if errors.Is(err, utils.ErrNotFound) {
				utils.Error(w, http.StatusNotFound, fmt.Sprintf("no section for id %d", id))
				return
			}
			utils.Error(w, http.StatusInternalServerError, "Some error occurs")
			return
		}
		utils.JSON(w, http.StatusNoContent, nil)
	}
}
