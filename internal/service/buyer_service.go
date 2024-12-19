package service

import (
	"log"

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

func (service *BuyerService) GetOne(id int) (*pkg.Buyer, error) {
	buyers, err := service.GetAll()

	if err != nil {
		log.Println("Error in GetAll - ", err)
		return nil, err
	}

	for _, buyer := range buyers {
		if int(buyer.ID) == id {
			return &buyer, err
		}
	}

	return nil, err
}
