package main

import (
	"github.com/meli-fresh-products-api-backend-go-t2/internal/repository"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/routes"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/service"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

func main() {
	err := utils.LoadProperties(".env")
	if err != nil {
		panic(err)
	}
	router := chi.NewRouter()

	// Create the routes and deps
	repo := repository.NewWarehouseDB(nil)
	warehouseService := service.NewWarehouseService(repo)
	err = routes.NewWarehouseRoutes(router, warehouseService)
	if err != nil {
		panic(err)
	}

	log.Printf("starting server at %s\n", os.Getenv("SERVER.PORT"))
	if err := http.ListenAndServe(os.Getenv("SERVER.PORT"), router); err != nil {
		panic(err)
	}
}
