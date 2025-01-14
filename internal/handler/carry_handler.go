package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type CarryHandler struct {
	service internal.CarryService
}

// NewCarryHandler creates a new instance of CarryHandler with the provided CarryService.
// It returns a pointer to the newly created CarryHandler.
//
// Parameters:
//   - service: an instance of CarryService that will be used by the CarryHandler.
//
// Returns:
//   - A pointer to the newly created CarryHandler.
func NewCarryHandler(service internal.CarryService) *CarryHandler {
	return &CarryHandler{service: service}
}

// SaveCarry handles the HTTP request to save a carry object.
// It decodes the request body into a Carry struct and calls the service layer to save it.
// If the decoding fails, it responds with a 400 Bad Request status.
// If the service layer returns an error, it responds with the appropriate status code:
// - 409 Conflict if the carry ID already exists
// - 422 Unprocessable Entity if the carry data is invalid
// - 404 Not Found if the locality is not found
// - 500 Internal Server Error for any other errors
// On success, it responds with a 201 Created status and the saved carry object.
func (handler *CarryHandler) SaveCarry() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var carry *internal.Carry
		err := json.NewDecoder(r.Body).Decode(&carry)
		if err != nil {
			http.Error(w, "Failed to decode carry: "+err.Error(), http.StatusBadRequest)
			return
		}
		if err := handler.service.Save(carry); err != nil {
			if errors.Is(err, utils.ErrConflict) {
				utils.Error(w, http.StatusConflict, "CID already exists: "+utils.ErrConflict.Error())
				return
			}
			if errors.Is(err, utils.ErrInvalidArguments) {
				utils.Error(w, http.StatusUnprocessableEntity, "Invalid carry: "+utils.ErrInvalidArguments.Error())
				return
			}
			if errors.Is(err, utils.ErrNotFound) {
				utils.Error(w, http.StatusConflict, "Locality: "+err.Error())
				return
			}
			utils.Error(w, http.StatusInternalServerError, "Failed to save carry: "+err.Error())
			return
		}
		utils.JSON(w, http.StatusCreated, carry)
	}
}

// GetAllCarries handles the HTTP request to retrieve all carries.
// It returns an HTTP handler function that writes the carries data as a JSON response.
// If an error occurs while retrieving the carries, it responds with a 404 status code and an error message.
func (handler *CarryHandler) GetAllCarries() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		carries, err := handler.service.GetAll()
		if err != nil {
			utils.Error(w, http.StatusNotFound, "Failed to get all carries")
			return
		}

		utils.JSON(w, http.StatusOK, carries)
	}
}

// GetCarryById handles the HTTP request to retrieve a carry by its ID.
// It extracts the ID from the URL parameters, validates it, and then
// calls the service layer to fetch the carry. If the ID is invalid or
// the carry is not found, it responds with the appropriate HTTP error
// status and message. On success, it responds with the carry data in
// JSON format.
//
// Returns an http.HandlerFunc that can be used to handle the request.
func (handler *CarryHandler) GetCarryById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "Invalid ID")
			return
		}

		carry, err := handler.service.GetById(id)
		if err != nil {
			utils.Error(w, http.StatusNotFound, "Failed to get carry")
			return
		}

		utils.JSON(w, http.StatusOK, carry)
	}
}

// UpdateCarry handles the HTTP request for updating a carry item.
// It extracts the carry ID from the URL parameters, decodes the request body into a Carry struct,
// and calls the service layer to update the carry item in the database.
// If the ID is invalid, the request body cannot be decoded, or the update fails,
// it responds with the appropriate HTTP status code and error message.
// On success, it responds with the updated carry item and a 200 OK status.
func (handler *CarryHandler) UpdateCarry() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "Invalid ID")
			return
		}

		carry := &internal.Carry{}
		err = json.NewDecoder(r.Body).Decode(carry)
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "Failed to decode carry: "+err.Error())
			return
		}

		carry.ID = id

		if err := handler.service.Update(carry); err != nil {
			if errors.Is(err, utils.ErrConflict) {
				utils.Error(w, http.StatusConflict, err.Error())
				return
			}
			if errors.Is(err, utils.ErrNotFound) {
				utils.Error(w, http.StatusNotFound, err.Error())
				return
			}
			utils.Error(w, http.StatusInternalServerError, "Failed to update carry: "+err.Error())
			return
		}

		utils.JSON(w, http.StatusOK, carry)
	}
}

// DeleteCarry handles the HTTP request for deleting a carry item by its ID.
// It extracts the ID from the URL parameters, validates it, and calls the service layer to delete the carry item.
// If the ID is invalid, it responds with a 400 Bad Request status.
// If the carry item is not found, it responds with a 404 Not Found status.
// If there is an internal server error during deletion, it responds with a 500 Internal Server Error status.
// On successful deletion, it responds with a 204 No Content status.
func (handler *CarryHandler) DeleteCarry() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "Invalid ID")
			return
		}

		if err := handler.service.Delete(id); err != nil {
			if errors.Is(err, utils.ErrNotFound) {
				utils.Error(w, http.StatusNotFound, "Carry not found")
				return
			}
			utils.Error(w, http.StatusInternalServerError, "Failed to delete carry")
			return
		}

		utils.JSON(w, http.StatusNoContent, "Carry deleted successfully")
	}
}
