package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"

	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

const (
	INVALID = "Invalid ID"
)

// reqPostWarehouse represents the request payload for creating a new warehouse.
// swagger:model
type reqPostWarehouse struct {
	// The unique code of the warehouse
	// required: true
	// example: WH-001
	Code string `json:"warehouse_code"`

	// The address of the warehouse
	// required: true
	// example: 1234 Warehouse St.
	Address string `json:"address"`

	// The telephone number of the warehouse
	// required: true
	// example: +1-800-555-5555
	Telephone string `json:"telephone"`

	// The ID of the locality where the warehouse is located
	// required: true
	// example: 101
	LocalityID int `json:"locality_id"`

	// The minimum capacity of the warehouse
	// required: true
	// example: 1000
	MinimumCapacity int `json:"minimum_capacity"`

	// The minimum temperature that the warehouse can maintain
	// required: true
	// example: -5
	MinimumTemperature int `json:"minimum_temperature"`
}

// WarehouseHandler handles HTTP requests related to warehouse operations.
//
//	@Summary		Handles warehouse operations
//	@Description	This handler provides endpoints to manage warehouse operations such as creating, updating, and retrieving warehouse information.
//	@Tags			warehouse
type WarehouseHandler struct {
	service internal.WarehouseService
}

func NewWarehouseHandler(service internal.WarehouseService) *WarehouseHandler {
	return &WarehouseHandler{service}
}

// GetAll handles the HTTP request to retrieve all warehouses.
//
//	@Summary		Get all warehouses
//	@Description	Retrieve a list of all warehouses
//	@Tags			warehouses
//	@Produce		json
//	@Success		200	{array}		internal.Warehouse	"List of warehouses"
//	@Failure		404	{object}	utils.ErrorResponse	"No warehouses found"
//	@Failure		500	{object}	utils.ErrorResponse	"An error occurred while retrieving warehouses"
//	@Router			/warehouses [get]
func (h *WarehouseHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		warehouses, err := h.service.GetAll()
		if err != nil {
			utils.HandleError(w, err)
			return
		}

		utils.JSON(w, http.StatusOK, warehouses)
	}
}

// GetByID handles the HTTP request to retrieve a warehouse by its ID.
//
//	@Summary		Get warehouse by ID
//	@Description	Get a warehouse by its ID
//	@Tags			warehouses
//	@Produce		json
//	@Param			id	path		int	true	"Warehouse ID"
//	@Success		200	{object}	internal.Warehouse
//	@Failure		400	{object}	utils.ErrorResponse	"Invalid ID format"
//	@Failure		404	{object}	utils.ErrorResponse	"No warehouse found with ID"
//	@Failure		500	{object}	utils.ErrorResponse	"An error occurred while retrieving the warehouse"
//	@Router			/warehouses/{id} [get]
func (h *WarehouseHandler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			utils.HandleError(w, utils.EBadRequest(INVALID))
			return
		}

		warehouse, err := h.service.GetByID(id)
		if err != nil {
			utils.HandleError(w, err)

			return
		}

		utils.JSON(w, http.StatusOK, warehouse)
	}
}

// Post handles the creation of a new warehouse.
//
//	@Summary		Create a new warehouse
//	@Description	Create a new warehouse with the provided details
//	@Tags			warehouses
//	@Accept			json
//	@Produce		json
//	@Param			warehouse	body		reqPostWarehouse	true	"Warehouse details"
//	@Success		201			{object}	internal.Warehouse	"Created warehouse"
//	@Failure		404			{object}	utils.ErrorResponse	"No warehouse found"
//	@Failure		400			{object}	utils.ErrorResponse	"Invalid request format"
//	@Failure		409			{object}	utils.ErrorResponse	"Warehouse code conflict"
//	@Failure		422			{object}	utils.ErrorResponse	"Invalid arguments"
//	@Failure		500			{object}	utils.ErrorResponse	"Internal server error"
//	@Router			/warehouses [post]
func (h *WarehouseHandler) Post() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body reqPostWarehouse
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			utils.HandleError(w, err)
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
			utils.HandleError(w, err)
			return
		}

		utils.JSON(w, http.StatusCreated, newWarehouse)
	}
}

// Update godoc
//
//	@Summary		Update a warehouse
//	@Description	Update the details of an existing warehouse by ID
//	@Tags			warehouses
//	@Accept			json
//	@Produce		json
//	@Param			id			path		int							true	"Warehouse ID"
//	@Param			warehouse	body		internal.WarehousePointers	true	"Warehouse data"
//	@Success		200			{object}	internal.Warehouse
//	@Failure		400			{object}	utils.ErrorResponse	"Invalid ID format or request body"
//	@Failure		404			{object}	utils.ErrorResponse	"Warehouse not found"
//	@Failure		409			{object}	utils.ErrorResponse	"Conflict error"
//	@Failure		422			{object}	utils.ErrorResponse	"Invalid arguments"
//	@Failure		500			{object}	utils.ErrorResponse	"Internal server error"
//	@Router			/warehouses/{id} [put]
func (h *WarehouseHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			utils.HandleError(w, utils.EBadRequest(INVALID))
			return
		}

		var body internal.WarehousePointers
		if err = json.NewDecoder(r.Body).Decode(&body); err != nil {
			utils.HandleError(w, err)
			return
		}

		updatedWarehouse, err := h.service.Update(id, body)
		if err != nil {
			utils.HandleError(w, err)
			return
		}

		utils.JSON(w, http.StatusOK, updatedWarehouse)
	}
}

// Delete handles the deletion of a warehouse by its ID.
//
//	@Summary		Delete a warehouse
//	@Description	Deletes a warehouse by its ID
//	@Tags			warehouses
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int					true	"Warehouse ID"
//	@Success		204	{object}	nil					"No Content"
//	@Failure		400	{object}	utils.ErrorResponse	"Invalid ID format"
//	@Failure		404	{object}	utils.ErrorResponse	"No warehouse found with the given ID"
//	@Failure		500	{object}	utils.ErrorResponse	"An error occurred while deleting the warehouse"
//	@Router			/warehouses/{id} [delete]
func (h *WarehouseHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			utils.HandleError(w, utils.EBadRequest(INVALID))
			return
		}

		err = h.service.Delete(id)

		if err != nil {
			utils.HandleError(w, err)
			return
		}

		utils.JSON(w, http.StatusNoContent, nil)
	}
}
