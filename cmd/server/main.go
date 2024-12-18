package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
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

	// ProductType
	repotypes := repository.NewProductTypeDB(nil)
	serviceType := service.NewProductTypeService(repotypes)
	if err := routes.NewProductTypeRoutes(router, serviceType); err != nil {
		panic(err)
	}

	// Product
	repo := repository.NewProductDB(nil)
	service := service.NewProductService(repo, serviceType)
	err = routes.NewProductRoutes(router, service)
	if err != nil {
		panic(err)
	}

	log.Printf("starting server at %s\n", os.Getenv("SERVER.PORT"))
	if err := http.ListenAndServe(os.Getenv("SERVER.PORT"), router); err != nil {
		panic(err)
	}
}
