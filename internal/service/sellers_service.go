package service

import (
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

func NewSellerService(rp internal.SellerRepository) *SellerService {
	return &SellerService{rp: rp}
}

type SellerService struct {
	rp internal.SellerRepository
}

func (s *SellerService) GetAll() (map[int]internal.Seller, error) {
	sellers, err := s.rp.GetAll()
	if err != nil {
		return nil, err
	}
	return sellers, nil
}

func (s *SellerService) GetById(id int) (internal.Seller, error) {
	seller, err := s.rp.GetById(id)
	if err != nil {
		return internal.Seller{}, err
	}
	if seller == (internal.Seller{}) {
		return internal.Seller{}, utils.ErrNotFound
	}
	return seller, nil
}

func (s *SellerService) Create(newSeller internal.SellerRequest) (internal.Seller, error) {
	sellerValidation := s.verify(newSeller)

	if sellerValidation != nil {
		return internal.Seller{}, sellerValidation
	}

	createdSeller, err := s.rp.Create(newSeller)
	if err != nil {
		return internal.Seller{}, err
	}

	return createdSeller, nil

}

func (s *SellerService) Update(id int, newSeller internal.SellerRequestPointer) (internal.Seller, error) {

	existingSeller, err := s.rp.GetById(id)
	if err != nil {
		return internal.Seller{}, err
	}
	if existingSeller == (internal.Seller{}) {
		return internal.Seller{}, utils.ErrNotFound
	}
	existingCid, err := s.rp.GetByCid(*newSeller.Cid)
	if err != nil {
		return internal.Seller{}, err
	}
	if existingCid.Cid != 0 && existingCid.ID != id {
		return internal.Seller{}, utils.ErrConflict
	}
	if *newSeller.Cid != 0 {
		existingSeller.Cid = *newSeller.Cid
	}

	if newSeller.CompanyName != nil {
		existingSeller.CompanyName = *newSeller.CompanyName
	}
	if newSeller.Address != nil {
		existingSeller.Address = *newSeller.Address
	}
	if newSeller.Telephone != nil {
		existingSeller.Telephone = *newSeller.Telephone
	}

	updatedSeller, err := s.rp.Update(existingSeller)
	if err != nil {
		return internal.Seller{}, err
	}

	return updatedSeller, nil

}

func (s *SellerService) Delete(id int) (bool, error) {

	existingSeller, err := s.rp.GetById(id)
	if err != nil {
		return false, err
	}
	if existingSeller == (internal.Seller{}) {
		return false, utils.ErrNotFound
	}
	result, err := s.rp.Delete(id)
	if err != nil {
		return false, err
	}
	return result, nil
}

func (s *SellerService) verify(newSeller internal.SellerRequest) error {

	if newSeller.Cid <= 0 {
		return utils.ErrInvalidArguments
	}

	existingCid, err := s.rp.GetByCid(newSeller.Cid)
	if err != nil {
		return err
	}
	if existingCid.Cid != 0 {
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
