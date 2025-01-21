package employee

import (
	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-go-t2/cmd/server/handler"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
)

// RegisterEmployeesRoutes is used to record the routes associated to the employee entity
func RegisterEmployeesRoutes(mux *chi.Mux, service internal.EmployeeService) error {
	employeeHandler := handler.NewEmployeeHandler(service)

	mux.Route("/api/v1/employees", func(router chi.Router) {
		// Get
		router.Get("/", employeeHandler.GetAllEmployees())
		router.Get("/{id}", employeeHandler.GetEmployeesByID())
		// Post
		router.Post("/", employeeHandler.PostEmployees())
		// Patch
		router.Patch("/{id}", employeeHandler.PatchEmployees())
		// Delete
		router.Delete("/{id}", employeeHandler.DeleteEmployees())
	})

	return nil
}
