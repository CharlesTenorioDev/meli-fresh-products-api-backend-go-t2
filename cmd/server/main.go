package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/loader"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/repository"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/routes"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/service"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

func main() {
	err := utils.LoadProperties(".env")
	if err != nil {
		panic(err)
	}

	router := chi.NewRouter()

	// Requisito 5 - Employees
	filePath := "docs/db/employees.json"
	ld := loader.NewEmployeeJsonFile(filePath)
	db, err := ld.Load()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create the routes and deps
	repo := repository.NewWarehouseDB(nil)
	warehouseService := service.NewWarehouseService(repo)
	err = routes.NewWarehouseRoutes(router, warehouseService)
	if err != nil {
		panic(err)
	}

	rp := repository.NewEmployeeRepository(db)
	sv := service.NewEmployeeService(rp, warehouseService)
	routes.RegisterEmployeesRoutes(router, sv)

	log.Printf("starting server at %s\n", os.Getenv("SERVER.PORT"))
	if err := http.ListenAndServe(os.Getenv("SERVER.PORT"), router); err != nil {
		panic(err)
	}
}
