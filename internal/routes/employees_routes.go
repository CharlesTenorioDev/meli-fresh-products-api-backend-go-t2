package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/handler"
	employeesPkg "github.com/meli-fresh-products-api-backend-go-t2/internal/pkg"
)

// RegisterEmployeesRoutes is used to record the routes associated to the employee entity
func RegisterEmployeesRoutes(mux *chi.Mux, service employeesPkg.EmployeeService) error {
	handler := handler.NewEmployeeHandler(service)
	mux.Route("/api/v1/employees", func(router chi.Router) {
		router.Get("/", handler.GetAllEmployees())
		router.Get("/{id}", handler.GetEmployeesById())
		router.Post("/", handler.PostEmployees())
	})

	return nil
}
