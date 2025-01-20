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

type reqPostWarehouse struct {
	Code               string `json:"warehouse_code"`
	Address            string `json:"address"`
	Telephone          string `json:"telephone"`
	LocalityID         int    `json:"locality_id"`
	MinimumCapacity    int    `json:"minimum_capacity"`
	MinimumTemperature int    `json:"minimum_temperature"`
}

type WarehouseHandler struct {
	service internal.WarehouseService
}

// NewWarehouseHandler creates a new instance of WarehouseHandler with the provided WarehouseService.
// It returns a pointer to the created WarehouseHandler.
//
// Parameters:
//   - service: an implementation of the WarehouseService interface.
//
// Returns:
//   - A pointer to the newly created WarehouseHandler.
func NewWarehouseHandler(service internal.WarehouseService) *WarehouseHandler {
	return &WarehouseHandler{service}
}

// GetAll handles the HTTP request to retrieve all warehouses.
// It returns a JSON response with the list of warehouses or an error message if no warehouses are found or an internal error occurs.
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

// GetByID handles the HTTP request to retrieve a warehouse by its ID.
func (h *WarehouseHandler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "Invalid ID format")
			return
		}

		warehouse, err := h.service.GetByID(id)
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

// Post handles the creation of a new warehouse.
// @Summary Create a new warehouse
// @Description Create a new warehouse with the provided details
// @Tags warehouses
// @Accept json
// @Produce json
// @Param warehouse body reqPostWarehouse true "Warehouse details"
// @Success 201 {object} internal.Warehouse "Created warehouse"
// @Failure 400 {object} utils.ErrorResponse "Invalid request format"
// @Failure 409 {object} utils.ErrorResponse "Warehouse code conflict"
// @Failure 422 {object} utils.ErrorResponse "Invalid arguments"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /warehouses [post]
func (h *WarehouseHandler) Post() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body reqPostWarehouse
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			utils.Error(w, http.StatusBadRequest, utils.ErrInvalidFormat.Error())
			return
		}

		newWarehouse := internal.Warehouse{
			WarehouseCode:      body.Code,
			Address:            body.Address,
			Telephone:          body.Telephone,
			LocalityID:         body.LocalityID,
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

// Update handles the HTTP request for updating a warehouse.
// It extracts the warehouse ID from the URL parameters and decodes the request body into a WarehousePointers struct.
// If the ID is invalid or the request body is improperly formatted, it returns a 400 Bad Request error.
// If the update service returns specific errors, it handles them accordingly:
// - 409 Conflict if there is a conflict error.
// - 404 Not Found if the warehouse is not found.
// - 422 Unprocessable Entity if there are invalid arguments.
// For any other errors, it returns a 500 Internal Server Error.
// On success, it returns the updated warehouse with a 200 OK status.
func (h *WarehouseHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "Invalid ID format")
			return
		}

		var body internal.WarehousePointers

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

			utils.Error(w, http.StatusInternalServerError, "An error occurred while updating the warehouse: "+err.Error())

			return
		}

		utils.JSON(w, http.StatusOK, updatedWarehouse)
	}
}

// Delete handles the HTTP request to delete a warehouse by its ID.
// It extracts the ID from the URL parameters, validates it, and calls the service layer to delete the warehouse.
// If the ID is invalid, it responds with a 400 Bad Request status.
// If the warehouse is not found, it responds with a 404 Not Found status.
// If an error occurs during deletion, it responds with a 500 Internal Server Error status.
// On successful deletion, it responds with a 204 No Content status.
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
