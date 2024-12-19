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
	// Create the routes and deps

	// Requisito 4 - ProductType
	productTypeRepo := repository.NewProductTypeDB(nil)
	productTypeService := service.NewProductTypeService(productTypeRepo)
	if err := routes.NewProductTypeRoutes(router, productTypeService); err != nil {
		panic(err)
	}

	// Requisito 4 - Product
	productRepo := repository.NewProductDB(nil)
	productService := service.NewProductService(productRepo, productTypeService)
	err = routes.NewProductRoutes(router, productService)
	if err != nil {
		panic(err)
	}

	// Requisito 2 - Warehouses
	warehouseRepo := repository.NewWarehouseDB(nil)
	warehouseService := service.NewWarehouseService(warehouseRepo)
	err = routes.NewWarehouseRoutes(router, warehouseService)
	if err != nil {
		panic(err)
	}

	// Requisito 3 - Section
	sectionRepo := repository.NewMemorySectionRepository(nil)
	sectionService := service.NewBasicSectionService(sectionRepo, warehouseService, productTypeService)
	err = routes.RegisterSectionRoutes(router, sectionService)
	if err != nil {
		panic(err)
	}

	// Requisito 5 - Employees
	filePath := "docs/db/employees.json"
	ld := loader.NewEmployeeJsonFile(filePath)
	db, err := ld.Load()
	if err != nil {
		fmt.Println(err)
		return
	}
	employeesRepo := repository.NewEmployeeRepository(db)
	employeesService := service.NewEmployeeService(employeesRepo, warehouseService)
	routes.RegisterEmployeesRoutes(router, employeesService)

	log.Printf("starting server at %s\n", os.Getenv("SERVER.PORT"))
	if err := http.ListenAndServe(os.Getenv("SERVER.PORT"), router); err != nil {
		panic(err)
	}
}
