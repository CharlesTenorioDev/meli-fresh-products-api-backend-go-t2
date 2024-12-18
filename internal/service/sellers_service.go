package service

import (

	"github.com/meli-fresh-products-api-backend-go-t2/internal/pkg"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
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
	if seller == (pkg.Seller{}) {
		return pkg.Seller{}, utils.ErrNotFound
	}
	return seller, nil
}

func (s *SellerService) Create(newSeller pkg.SellerRequest) (pkg.Seller, error) {
	sellerValidation := s.verify(newSeller)

	if sellerValidation != nil {
		return pkg.Seller{}, sellerValidation
	}

	createdSeller, err := s.rp.Create(newSeller)
	if err != nil {
		return pkg.Seller{}, err
	}

	return createdSeller, nil

}

func (s *SellerService) verify(newSeller pkg.SellerRequest) (error) {

	if newSeller.Cid <= 0 {
		return utils.ErrInvalidArguments
	}

	existintgCid, err := s.rp.GetByCid(newSeller.Cid)
	if err != nil {
		return err
	}
	if existintgCid.Cid != 0 {
		return utils.ErrConflict
	}

	if len(newSeller.CompanyName) == 0 {
		return utils.ErrInvalidArguments
	}
	if len(newSeller.Telephone) == 0 {
		return utils.ErrInvalidArguments
	}
	if len(newSeller.Address) == 0 {
		return utils.ErrInvalidArguments
	}

	return nil

}

	// existintgSeller, err := s.rp.GetById(newSeller.ID)
	// if err != nil {
	// 	return pkg.Seller{}, err
	// }
	// if existintgSeller.ID != 0 {
	// 	return pkg.Seller{}, utils.ErrConflict
	// }


