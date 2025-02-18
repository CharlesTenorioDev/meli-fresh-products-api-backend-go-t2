package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/pkg/logger"

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
		r = logger.GetContext(r, "EMPLOYEE:GET_ALL")
		logger.Start(r)

		employees, err := h.sv.FindAll()
		if err != nil {
			utils.HandleErrorContext(r.Context(), w, err)
			return
		}
		// returns status 200 and the data if all ok
		utils.JSONContext(r.Context(), w, http.StatusOK, employees)
	}
}

// GetEmployeesByID handles the GET /employees/{id} route
func (h *EmployeeDefault) GetEmployeesByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r = logger.GetContext(r, "EMPLOYEE:GET_BY_ID")
		logger.Start(r)

		idStr := chi.URLParam(r, "id")

		id, err := strconv.Atoi(idStr)
		if err != nil {
			utils.HandleErrorContext(r.Context(), w, utils.EBadRequest("id"))
			return
		}

		employee, err := h.sv.FindByID(id)
		if err != nil {
			utils.HandleErrorContext(r.Context(), w, err)
			return
		}

		data := internal.Employee{
			ID: employee.ID,
			Attributes: internal.EmployeeAttributes{
				CardNumberID: employee.Attributes.CardNumberID,
				FirstName:    employee.Attributes.FirstName,
				LastName:     employee.Attributes.LastName,
				WarehouseID:  employee.Attributes.WarehouseID,
			},
		}

		// returns status 200 and the data if all ok
		utils.JSONContext(r.Context(), w, http.StatusOK, data)
	}
}

// PostEmployees handles the POST /employees route
func (h *EmployeeDefault) PostEmployees() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r = logger.GetContext(r, "EMPLOYEE:CREATE")
		logger.Start(r)

		var newEmployee internal.EmployeeAttributes

		// decode the json request body
		err := json.NewDecoder(r.Body).Decode(&newEmployee)
		if err != nil {
			utils.HandleErrorContext(r.Context(), w, utils.EBadRequest("body"))
			return
		}

		logger.Info(r.Context(), "processing", newEmployee)

		// create the employee
		employee, err := h.sv.CreateEmployee(newEmployee)
		if err != nil {
			utils.HandleErrorContext(r.Context(), w, err)
			return
		}

		// returns status 201 and the data if all ok
		utils.JSONContext(r.Context(), w, http.StatusCreated, employee)
	}
}

// PatchEmployees handles the PATCH /employees route
func (h *EmployeeDefault) PatchEmployees() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r = logger.GetContext(r, "EMPLOYEE:UPDATE")
		logger.Start(r)

		idStr := chi.URLParam(r, "id")

		id, err := strconv.Atoi(idStr)
		if err != nil {
			utils.HandleErrorContext(r.Context(), w, utils.EBadRequest("id"))
			return
		}

		var inputEmployee internal.Employee
		// decode the json request body into Employee struct
		err = json.NewDecoder(r.Body).Decode(&inputEmployee)
		if err != nil {
			utils.HandleErrorContext(r.Context(), w, utils.EBadRequest("body"))
			return
		}

		logger.Info(r.Context(), "processing", inputEmployee)

		inputEmployee.ID = id
		// update the employee
		employee, err := h.sv.UpdateEmployee(inputEmployee)
		if err != nil {
			utils.HandleErrorContext(r.Context(), w, err)
			return
		}

		// returns status 200 and the data if all ok
		utils.JSONContext(r.Context(), w, http.StatusOK, employee)
	}
}

// DeleteEmployees handles the DELETE /employees/{id} route
// it deletes an existing employee based on the provided ID
func (h *EmployeeDefault) DeleteEmployees() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r = logger.GetContext(r, "EMPLOYEE:DELETE")
		logger.Start(r)

		// extract the employee ID from the URL parameters and converts it to int
		idStr := chi.URLParam(r, "id")

		id, err := strconv.Atoi(idStr)
		if err != nil {
			utils.HandleErrorContext(r.Context(), w, utils.EBadRequest("id"))
			return
		}

		// delete the employee
		err = h.sv.DeleteEmployee(id)
		if err != nil {
			utils.HandleErrorContext(r.Context(), w, err)

			return
		}

		// returns status 204 and a success message if all ok
		utils.JSONContext(r.Context(), w, http.StatusNoContent, nil)
	}
}
