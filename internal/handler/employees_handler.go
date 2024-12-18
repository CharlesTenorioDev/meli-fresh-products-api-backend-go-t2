package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bootcamp-go/web/response"
	"github.com/go-chi/chi/v5"
	employeesPkg "github.com/meli-fresh-products-api-backend-go-t2/internal/pkg"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

// EmployeeDefault is the http handler for employee-related endpoints
// it communicates with the service layer to process requests
type EmployeeDefault struct {
	sv employeesPkg.EmployeeService
}

// NewEmployeeHandler creates a new instance of EmployeeDefault
// takes an EmployeeService implementation as a parameter
func NewEmployeeHandler(sv employeesPkg.EmployeeService) *EmployeeDefault {
	return &EmployeeDefault{sv: sv}
}

// GetAllEmployees handles the GET /employees route
// it retrieves all employees and sends them as a json response
func (h *EmployeeDefault) GetAllEmployees() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// fetch all employees from the service
		employees, err := h.sv.FindAll()
		if err != nil {
			// returns status 500 if an error occurs
			response.JSON(w, http.StatusInternalServerError, nil)
			return
		}

		// transform the data into the appropriate format
		data := make(map[int]employeesPkg.Employee)
		for key, value := range employees {
			data[key] = employeesPkg.Employee{
				ID: value.ID,
				Attributes: employeesPkg.EmployeeAttributes{
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
// it retrieves a specific employee by its ID and sends it as a json response
func (h *EmployeeDefault) GetEmployeesById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// extract the employee ID from the url parameters and converts it to int
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			// returns status 400 if the ID is invalid
			response.Error(w, http.StatusBadRequest, utils.ErrInvalidFormat.Error())
		}

		// fetch the employee from the service
		employee, err := h.sv.FindById(id)
		if err != nil {
			// returns status 404 if the employee is not found
			response.Error(w, http.StatusNotFound, utils.ErrInvalidArguments.Error())
		}

		// extract the employee from the returned map
		e := employee[id]

		// prepare the response data
		data := employeesPkg.Employee{
			ID: e.ID,
			Attributes: employeesPkg.EmployeeAttributes{
				CardNumberId: e.Attributes.CardNumberId,
				FirstName:    e.Attributes.FirstName,
				LastName:     e.Attributes.LastName,
				WarehouseId:  e.Attributes.WarehouseId,
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
// it creates a new employee based on the provided json body
func (h *EmployeeDefault) PostEmployees() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newEmployee employeesPkg.EmployeeAttributes

		// decode the json request body into an EmployeeAttributes struct
		err := json.NewDecoder(r.Body).Decode(&newEmployee)
		if err != nil {
			// returns status 400 if the json format is invalid
			response.Error(w, http.StatusBadRequest, utils.ErrInvalidFormat.Error())
			return
		}

		employee, err := h.sv.CreateEmployee(newEmployee)
		if err != nil {
			if err == utils.ErrConflict {
				// returns status 409 if a duplicate employee already exists
				response.Error(w, http.StatusConflict, utils.ErrConflict.Error())
				return

			} else {
				// returns status 422 if another validation error occurs
				response.Error(w, http.StatusUnprocessableEntity, utils.ErrInvalidArguments.Error())
				return
			}
		}

		// returns status 201 and the created employee if all ok
		response.JSON(w, http.StatusCreated, map[string]any{
			"message": "success",
			"data":    employee,
		})
	}
}
