package service

import (
	"github.com/meli-fresh-products-api-backend-go-t2/internal/pkg"
)

type BuyerService struct {
	repo pkg.BuyerRepository
}

func NewBuyer(repo pkg.BuyerRepository) *BuyerService {
	return &BuyerService{repo: repo}
}

func (service *BuyerService) GetAll() (buyer []pkg.Buyer, err error) {
	buyer, err = service.repo.GetAll()
	return buyer, err
}
