package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"

	"github.com/go-chi/chi/v5"
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
	service internal.SectionService
}

// NewSectionHandler Get a new instance of SectionHandler
func NewSectionHandler(service internal.SectionService) *SectionHandler {
	return &SectionHandler{service}
}

// GetAll handles the HTTP request to retrieve all sections.
// @Summary Get all sections
// @Description Retrieve a list of all sections
// @Tags sections
// @Produce json
// @Success 200 {array} internal.Section "List of sections"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/sections [get]
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

// GetById godoc
// @Summary Get section by ID
// @Description Get a section by its ID
// @Tags sections
// @Accept json
// @Produce json
// @Param id path int true "Section ID"
// @Success 200 {object} internal.Section
// @Failure 400 {object} utils.ErrorResponse "Invalid ID"
// @Failure 404 {object} utils.ErrorResponse "Section not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/sections/{id} [get]
func (h *SectionHandler) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "invalid id")
			return
		}

		section, err := h.service.GetByID(id)
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

// Post handles the creation of a new section.
//
// @Summary Create a new section
// @Description Create a new section with the provided details
// @Tags sections
// @Accept json
// @Produce json
// @Param section body reqPostSection true "Section details"
// @Success 201 {object} internal.Section
// @Failure 400 {object} utils.ErrorResponse "Invalid request format"
// @Failure 409 {object} utils.ErrorResponse "Section conflict"
// @Failure 422 {object} utils.ErrorResponse "Invalid arguments"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/sections [post]
func (h *SectionHandler) Post() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body reqPostSection
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			utils.Error(w, http.StatusBadRequest, utils.ErrInvalidFormat.Error())
			return
		}

		newSection := internal.Section{
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

// Update godoc
// @Summary Update a section
// @Description Update a section by ID
// @Tags sections
// @Accept json
// @Produce json
// @Param id path int true "Section ID"
// @Param section body internal.SectionPointers true "Section data"
// @Success 200 {object} internal.Section
// @Failure 400 {object} utils.ErrorResponse "Invalid ID or request body"
// @Failure 409 {object} utils.ErrorResponse "Conflict error"
// @Failure 422 {object} utils.ErrorResponse "Unprocessable entity"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/sections/{id} [put]
func (h *SectionHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "invalid id")
			return
		}

		var body internal.SectionPointers
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

// Delete handles the deletion of a section by its ID.
//
// @Summary Delete a section
// @Description Delete a section by its ID
// @Tags sections
// @Accept json
// @Produce json
// @Param id path int true "Section ID"
// @Success 204 {object} nil "No Content"
// @Failure 400 {object} utils.ErrorResponse "Invalid ID"
// @Failure 404 {object} utils.ErrorResponse "Section not found"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /api/v1/sections/{id} [delete]
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

// GetSectionProductsReport godoc
// @Summary Get section products report
// @Description Retrieves a report of products for a given section ID
// @Tags sections
// @Accept json
// @Produce json
// @Param id query int false "Section ID"
// @Success 200 {object} internal.SectionProductsReport
// @Failure 400 {object} utils.ErrorResponse "Invalid ID"
// @Failure 404 {object} utils.ErrorResponse "Section not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /api/v1/sections/products/report [get]
func (h *SectionHandler) GetSectionProductsReport() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idReq := strings.TrimSpace(r.URL.Query().Get("id"))
		id := 0

		var err error

		if idReq != "" {
			id, err = strconv.Atoi(idReq)
			if err != nil {
				utils.Error(w, http.StatusBadRequest, "invalid id")
				return
			}
		}

		sectionProductReport, err := h.service.GetSectionProductsReport(id)
		if err != nil {
			if errors.Is(err, utils.ErrNotFound) {
				utils.Error(w, http.StatusNotFound, fmt.Sprintf("no section for id %d", id))
				return
			}

			utils.Error(w, http.StatusInternalServerError, "Some error occurs")

			return
		}

		utils.JSON(w, http.StatusOK, sectionProductReport)
	}
}
