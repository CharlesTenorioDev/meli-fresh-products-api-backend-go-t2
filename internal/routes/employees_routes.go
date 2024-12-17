package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/yywatanabe_meli/api-produtos-frescos/internal/handler"
	employeesPkg "github.com/yywatanabe_meli/api-produtos-frescos/internal/pkg"
	"github.com/yywatanabe_meli/api-produtos-frescos/internal/repository"
)

func RegisterEmployeesRoutes(mux *chi.Mux, service employeesPkg.EmployeeService) error {
	handler := handler.NewEmployeeHandler(service)
	mux.Route("/employees", func(router chi.Router) {
		router.Get("/", handler.GetAllEmployees())
	})

	router := chi.NewRouter()

	// - repository
	rp := repository.NewEmployeeRepository(db)
	// - service
	sv := service.NewEmployeeDefault(rp)
	// - handler
	hd := handler.NewEmployeeS(sv)
	// Create the routes and deps
	router.Route("/api/v1", func(rt chi.Router) {
		// Employees
		rt.Get("/employees", hd.GetAllEmployees())
	})
	return nil
}
