package handler

import (
	"net/http"

	"github.com/bootcamp-go/web/response"
	employeesPkg "github.com/meli-fresh-products-api-backend-go-t2/internal/pkg"
)

type EmployeeDefault struct {
	sv employeesPkg.EmployeeService
}

func NewEmployeeHandler(sv employeesPkg.EmployeeService) *EmployeeDefault {
	return &EmployeeDefault{sv: sv}
}

// GetAll is a method that returns a handler for the route GET /Employees
func (h *EmployeeDefault) GetAllEmployees() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		employees, err := h.sv.FindAll()
		if err != nil {
			response.JSON(w, http.StatusInternalServerError, nil)
			return
		}

		data := make(map[int]employeesPkg.EmployeeJson)
		for key, value := range employees {
			data[key] = employeesPkg.EmployeeJson{
				ID:           value.ID,
				CardNumberId: value.CardNumberId,
				FirstName:    value.FirstName,
				LastName:     value.LastName,
				WarehouseId:  value.WarehouseId,
			}
		}
		response.JSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}
