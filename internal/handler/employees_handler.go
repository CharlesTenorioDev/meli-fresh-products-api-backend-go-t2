package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/bootcamp-go/web/response"
	"github.com/go-chi/chi/v5"
	employeesPkg "github.com/meli-fresh-products-api-backend-go-t2/internal/pkg"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type EmployeeDefault struct {
	sv employeesPkg.EmployeeService
}

func NewEmployeeHandler(sv employeesPkg.EmployeeService) *EmployeeDefault {
	return &EmployeeDefault{sv: sv}
}

// GetAllEmployees is a method that returns a handler for the route GET /employees
func (h *EmployeeDefault) GetAllEmployees() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		employees, err := h.sv.FindAll()
		if err != nil {
			response.JSON(w, http.StatusInternalServerError, nil)
			return
		}

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
		response.JSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

// GetEmployeesById is a method that returns a handler for the route GET /employees/{id}
func (h *EmployeeDefault) GetEmployeesById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			response.Error(w, http.StatusBadRequest, utils.ErrInvalidFormat.Error())
		}
		employee, err := h.sv.FindById(id)
		if err != nil {
			response.Error(w, http.StatusNotFound, utils.ErrInvalidArguments.Error())
		}

		e := employee[id]

		data := employeesPkg.Employee{
			ID: e.ID,
			Attributes: employeesPkg.EmployeeAttributes{
				CardNumberId: e.Attributes.CardNumberId,
				FirstName:    e.Attributes.FirstName,
				LastName:     e.Attributes.LastName,
				WarehouseId:  e.Attributes.WarehouseId,
			},
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

func (h *EmployeeDefault) PostEmployees() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newEmployee employeesPkg.EmployeeAttributes

		err := json.NewDecoder(r.Body).Decode(&newEmployee)
		if err != nil {
			response.Error(w, http.StatusBadRequest, utils.ErrInvalidFormat.Error())
			return
		}

		employee, err := h.sv.CreateEmployee(newEmployee)
		if err != nil {
			if err == utils.ErrConflict {
				response.Error(w, http.StatusConflict, utils.ErrConflict.Error())
				return

			} else {
				response.Error(w, http.StatusUnprocessableEntity, utils.ErrInvalidArguments.Error())
				fmt.Println(err)
				return
			}
		}

		response.JSON(w, http.StatusCreated, map[string]any{
			"message": "success",
			"data":    employee,
		})
	}
}
