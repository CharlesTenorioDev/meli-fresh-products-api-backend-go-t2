package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"

	"github.com/bootcamp-go/web/response"
	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

// EmployeeDefault is the http handler for employee-related endpoints
// it communicates with the service layer to process requests
type EmployeeDefault struct {
	sv internal.EmployeeService
}

// NewEmployeeHandler creates a new instance of EmployeeDefault
func NewEmployeeHandler(sv internal.EmployeeService) *EmployeeDefault {
	return &EmployeeDefault{sv: sv}
}

// GetAllEmployees handles the GET /employees route
func (h *EmployeeDefault) GetAllEmployees() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		employees, err := h.sv.FindAll()
		if err != nil {
			response.JSON(w, http.StatusInternalServerError, nil)
			return
		}

		data := make(map[int]internal.Employee)
		for key, value := range employees {
			data[key] = internal.Employee{
				ID: value.ID,
				Attributes: internal.EmployeeAttributes{
					CardNumberId: value.Attributes.CardNumberId,
					FirstName:    value.Attributes.FirstName,
					LastName:     value.Attributes.LastName,
					WarehouseId:  value.Attributes.WarehouseId,
				},
			}
		}

		// returns status 200 and the data if all ok
		response.JSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

// GetEmployeesById handles the GET /employees/{id} route
func (h *EmployeeDefault) GetEmployeesById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			utils.HandleError(w, utils.ErrInvalidFormat)
			return
		}

		employee, err := h.sv.FindById(id)
		if err != nil {
			utils.HandleError(w, utils.ErrNotFound)
			return
		}

		data := internal.Employee{
			ID: employee.ID,
			Attributes: internal.EmployeeAttributes{
				CardNumberId: employee.Attributes.CardNumberId,
				FirstName:    employee.Attributes.FirstName,
				LastName:     employee.Attributes.LastName,
				WarehouseId:  employee.Attributes.WarehouseId,
			},
		}

		// returns status 200 and the data if all ok
		response.JSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

// PostEmployees handles the POST /employees route
func (h *EmployeeDefault) PostEmployees() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newEmployee internal.EmployeeAttributes

		// decode the json request body
		err := json.NewDecoder(r.Body).Decode(&newEmployee)
		if err != nil {
			utils.HandleError(w, utils.ErrInvalidFormat)
			return
		}

		// create the employee
		employee, err := h.sv.CreateEmployee(newEmployee)
		if err != nil {
			if err == utils.ErrConflict {
				utils.HandleError(w, utils.ErrConflict)
			} else if err == utils.ErrEmptyArguments {
				utils.HandleError(w, utils.ErrEmptyArguments)
			} else if err == utils.ErrWarehouseDoesNotExists {
				utils.HandleError(w, utils.ErrWarehouseDoesNotExists)
			} else {
				utils.HandleError(w, utils.ErrInvalidArguments)
			}
			return
		}

		// returns status 201 and the data if all ok
		response.JSON(w, http.StatusCreated, map[string]any{
			"message": "success",
			"data":    employee,
		})
	}
}

// PatchEmployees handles the PATCH /employees route
func (h *EmployeeDefault) PatchEmployees() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			utils.HandleError(w, utils.ErrInvalidFormat)
			return
		}

		var inputEmployee internal.Employee
		// decode the json request body into Employee struct
		err = json.NewDecoder(r.Body).Decode(&inputEmployee)
		if err != nil {
			utils.HandleError(w, utils.ErrInvalidFormat)
			return
		}

		inputEmployee.ID = id
		// update the employee
		employee, err := h.sv.UpdateEmployee(inputEmployee)
		if err != nil {
			if err == utils.ErrNotFound {
				utils.HandleError(w, utils.ErrNotFound)
			} else {
				utils.HandleError(w, utils.ErrWarehouseDoesNotExists)
			}
			return
		}

		// returns status 200 and the data if all ok
		response.JSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    employee,
		})
	}
}

// DeleteEmployees handles the DELETE /employees/{id} route
// it deletes an existing employee based on the provided ID
func (h *EmployeeDefault) DeleteEmployees() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// extract the employee ID from the URL parameters and converts it to int
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			utils.HandleError(w, utils.ErrInvalidFormat)
			return
		}

		// delete the employee
		err = h.sv.DeleteEmployee(id)
		if err != nil {
			if err == utils.ErrNotFound {
				utils.HandleError(w, utils.ErrNotFound)
			} else {
				utils.HandleError(w, utils.ErrInvalidArguments)
			}
			return
		}

		// returns status 204 and a success message if all ok
		response.JSON(w, http.StatusNoContent, map[string]any{
			"message": "employee deleted successfully",
		})
	}
}
