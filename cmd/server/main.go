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
	sectionRepo := repository.NewMemorySectionRepository(nil)
	sectionValidations := service.SectionsValidation{}
	sectionService := service.NewBasicSectionService(sectionRepo, &sectionValidations)
	routes.RegisterSectionRoutes(router, sectionService)

	log.Printf("starting server at %s\n", os.Getenv("SERVER.PORT"))
	if err := http.ListenAndServe(os.Getenv("SERVER.PORT"), router); err != nil {
		panic(err)
	}
}
