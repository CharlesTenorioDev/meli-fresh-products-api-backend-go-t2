package service

import (
	"github.com/meli-fresh-products-api-backend-go-t2/internal/pkg"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/repository"
)

func GetAllBuyers() ([]pkg.Buyer, error) {
	return repository.GetAllBuyers()
}
