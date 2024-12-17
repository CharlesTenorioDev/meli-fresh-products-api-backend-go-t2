package service

import (

	"github.com/meli-fresh-products-api-backend-go-t2/internal/pkg"
)


func NewSellerService(rp pkg.SellerRepository) *SellerService {
	return &SellerService{rp: rp}
}

type SellerService struct {
	rp pkg.SellerRepository
}

func (s *SellerService) GetAll() (map[int]pkg.Seller, error) {
	sellers, err := s.rp.GetAll()
	if err != nil {
		return nil, err
	}
	return sellers, nil
}

func (s *SellerService) GetById(id int) (pkg.Seller, error) {
	seller, err := s.rp.GetById(id)
	if err != nil {
		return pkg.Seller{}, err
	}
	return seller, nil
}