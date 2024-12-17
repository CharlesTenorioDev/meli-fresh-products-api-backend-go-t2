package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/yywatanabe_meli/api-produtos-frescos/internal/loader"
	"github.com/yywatanabe_meli/api-produtos-frescos/internal/repository"
	"github.com/yywatanabe_meli/api-produtos-frescos/internal/routes"
	"github.com/yywatanabe_meli/api-produtos-frescos/internal/service"
	"github.com/yywatanabe_meli/api-produtos-frescos/internal/utils"
)

func main() {
	err := utils.LoadProperties(".env")
	if err != nil {
		panic(err)
	}

	router := chi.NewRouter()

	// - loader
	filePath := "docs/db/employees.json"
	ld := loader.NewEmployeeJsonFile(filePath)
	db, err := ld.Load()
	if err != nil {
		fmt.Println(err)
		return
	}

	rp := repository.NewEmployeeRepository()
	sv := service.NewEmployeeService(rp)
	routes.RegisterEmployeesRoutes(router, sv)

	log.Printf("starting server at %s\n", os.Getenv("SERVER.PORT"))
	if err := http.ListenAndServe(os.Getenv("SERVER.PORT"), router); err != nil {
		panic(err)
	}
}
