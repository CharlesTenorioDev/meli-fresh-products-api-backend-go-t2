package service

import (
	"errors"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

func NewSellerService(rp internal.SellerRepository, localityRp internal.SellerLocalityValidation) internal.SellerService {
	return &SellerService{rp: rp, localityRp: localityRp}
}

type SellerService struct {
	rp         internal.SellerRepository
	localityRp internal.SellerLocalityValidation
}

func (s *SellerService) GetAll() ([]internal.Seller, error) {
	sellers, err := s.rp.GetAll()

	if err != nil {
		return nil, err
	}

	return sellers, nil
}

func (s *SellerService) GetByID(id int) (internal.Seller, error) {
	seller, err := s.rp.GetById(id)

	if err != nil {
		return internal.Seller{}, err
	}

	if seller == (internal.Seller{}) {
		return internal.Seller{}, utils.ErrNotFound
	}

	return seller, nil
}

func (s *SellerService) Create(newSeller *internal.Seller) error {
	sellerValidation := s.verify(*newSeller)

	if sellerValidation != nil {
		return sellerValidation
	}

	_, err := s.localityRp.GetByID(newSeller.LocalityId)

	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			return errors.Join(utils.ErrInvalidArguments, errors.New("invalid locality_id"))
		}

		return err
	}

	err = s.rp.Create(newSeller)

	if err != nil {
		return err
	}

	return nil
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

	err = s.rp.Update(&existingSeller)
	if err != nil {
		return internal.Seller{}, err
	}

	return existingSeller, nil
}

func (s *SellerService) Delete(id int) error {
	existingSeller, err := s.rp.GetById(id)

	if err != nil {
		return err
	}

	if existingSeller == (internal.Seller{}) {
		return utils.ErrNotFound
	}

	err = s.rp.Delete(id)

	if err != nil {
		return err
	}

	return nil
}

func (s *SellerService) verify(newSeller internal.Seller) error {
	if newSeller.Cid <= 0 {
		return utils.ErrInvalidArguments
	}

	existingCid, err := s.rp.GetByCid(newSeller.Cid)
	if err != nil && !errors.Is(err, utils.ErrNotFound) {
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
