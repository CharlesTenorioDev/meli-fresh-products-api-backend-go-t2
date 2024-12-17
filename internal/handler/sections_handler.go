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

type reqPostSection struct {
	SectionNumber      int     `json:"section_number"`
	CurrentCapacity    int     `json:"current_capacity"`
	MaximumCapacity    int     `json:"maximum_capacity"`
	MinimumCapacity    int     `json:"minimum_capacity"`
	CurrentTemperature float64 `json:"current_temperature"`
	MinimumTemperature float64 `json:"minimum_temperature"`
	ProductTypeID      int     `json:"warehouse_id"`
	WarehouseID        int     `json:"product_type_id"`
}

type SectionHandler struct {
	service pkg.SectionService
}

func NewSectionHandler(service pkg.SectionService) *SectionHandler {
	return &SectionHandler{service}
}

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
				utils.Error(w, http.StatusNotFound, fmt.Sprintf("no section with id %d", id))
				return
			}
			utils.Error(w, http.StatusInternalServerError, "Some error occurs")
			return
		}
		utils.JSON(w, http.StatusOK, section)
	}
}

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
			}
		}
		utils.JSON(w, http.StatusCreated, newSection)
	}
}

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
				utils.Error(w, http.StatusNotFound, fmt.Sprintf("no section with id %d", id))
				return
			}
			utils.Error(w, http.StatusInternalServerError, "Some error occurs")
			return
		}
		utils.JSON(w, http.StatusNoContent, nil)
	}
}
