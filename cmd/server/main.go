package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/repository"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/routes"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/service"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

func main() {
	err := utils.LoadProperties("/Users/dfcarvalho/Documents/aulas-go-meli/meli-fresh-products-api-backend-go-t2/.env")
	if err != nil {
		panic(err)
	}

	router := chi.NewRouter()

	ld := internal.NewSellerJSONFile("/Users/dfcarvalho/Documents/aulas-go-meli/meli-fresh-products-api-backend-go-t2/internal/sellers.json")
	db, err := ld.Load()
	if err != nil {
		return
	}

	rp := repository.NewSellerDbRepository(db)
	sv := service.NewSellerService(rp)
	routes.RegisterSellerRoutes(router, sv)

	log.Printf("starting server at %s\n", os.Getenv("SERVER.PORT"))
	if err := http.ListenAndServe(os.Getenv("SERVER.PORT"), router); err != nil {
		panic(err)
	}
}
