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

type reqPostWarehouse struct {
	Code               string `json:"warehouse_code"`
	Address            string `json:"address"`
	Telephone          string `json:"telephone"`
	MinimumCapacity    int    `json:"minimum_capacity"`
	MinimumTemperature int    `json:"minimum_temperature"`
}

type WarehouseHandler struct {
	service pkg.WarehouseService
}

func NewWarehouseHandler(service pkg.WarehouseService) *WarehouseHandler {
	return &WarehouseHandler{service}
}

func (h *WarehouseHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		warehouses, err := h.service.GetAll()
		if err != nil {
			if errors.Is(err, utils.ErrNotFound) {
				utils.Error(w, http.StatusNotFound, "No warehouses found")
				return
			}
			utils.Error(w, http.StatusInternalServerError, "An error occurred while retrieving warehouses")
			return
		}
		utils.JSON(w, http.StatusOK, warehouses)
	}
}

func (h *WarehouseHandler) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "Invalid ID format")
			return
		}
		warehouse, err := h.service.GetById(id)
		if err != nil {
			if errors.Is(err, utils.ErrNotFound) {
				utils.Error(w, http.StatusNotFound, fmt.Sprintf("No warehouse found with ID %d", id))
				return
			}
			utils.Error(w, http.StatusInternalServerError, "An error occurred while retrieving the warehouse")
			return
		}
		utils.JSON(w, http.StatusOK, warehouse)
	}
}

func (h *WarehouseHandler) Post() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body reqPostWarehouse
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			utils.Error(w, http.StatusBadRequest, utils.ErrInvalidFormat.Error())
			return
		}
		newWarehouse := pkg.Warehouse{
			WarehouseCode:      body.Code,
			Address:            body.Address,
			Telephone:          body.Telephone,
			MinimumCapacity:    body.MinimumCapacity,
			MinimumTemperature: body.MinimumTemperature,
		}
		newWarehouse, err := h.service.Save(newWarehouse)
		if err != nil {
			if errors.Is(err, utils.ErrConflict) {
				utils.Error(w, http.StatusConflict, err.Error())
				return
			}
			if errors.Is(err, utils.ErrInvalidArguments) {
				utils.Error(w, http.StatusUnprocessableEntity, err.Error())
				return
			}

			utils.Error(w, http.StatusInternalServerError, "An error occurred while saving the warehouse")
			return
		}
		utils.JSON(w, http.StatusCreated, newWarehouse)
	}
}

func (h *WarehouseHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "Invalid ID format")
			return
		}
		var body pkg.WarehousePointers
		if err = json.NewDecoder(r.Body).Decode(&body); err != nil {
			utils.Error(w, http.StatusBadRequest, utils.ErrInvalidFormat.Error())
			return
		}
		updatedWarehouse, err := h.service.Update(id, body)
		if err != nil {
			if errors.Is(err, utils.ErrConflict) {
				utils.Error(w, http.StatusConflict, err.Error())
				return
			}
			if errors.Is(err, utils.ErrNotFound) {
				utils.Error(w, http.StatusNotFound, fmt.Sprintf("No warehouse found with ID %d", id))
				return
			}
			if errors.Is(err, utils.ErrInvalidArguments) {
				utils.Error(w, http.StatusUnprocessableEntity, err.Error())
				return
			}
			utils.Error(w, http.StatusInternalServerError, "An error occurred while updating the warehouse")
			return
		}
		utils.JSON(w, http.StatusOK, updatedWarehouse)
	}
}

func (h *WarehouseHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "Invalid ID format")
			return
		}
		err = h.service.Delete(id)
		if err != nil {
			if errors.Is(err, utils.ErrNotFound) {
				utils.Error(w, http.StatusNotFound, fmt.Sprintf("No warehouse found with ID %d", id))
				return
			}
			utils.Error(w, http.StatusInternalServerError, "An error occurred while deleting the warehouse")
			return
		}
		utils.JSON(w, http.StatusNoContent, nil)
	}
}
