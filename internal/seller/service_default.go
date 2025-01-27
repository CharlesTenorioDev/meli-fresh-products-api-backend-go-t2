package seller

import (
	"errors"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type DefaultSellerService struct {
	rp         internal.SellerRepository
	localityRp internal.SellerLocalityValidation
}

func NewSellerService(rp internal.SellerRepository, localityRp internal.SellerLocalityValidation) internal.SellerService {
	return &DefaultSellerService{rp: rp, localityRp: localityRp}
}

func (s *DefaultSellerService) GetAll() ([]internal.Seller, error) {
	sellers, err := s.rp.GetAll()

	if err != nil {
		return nil, err
	}

	return sellers, nil
}

func (s *DefaultSellerService) GetByID(id int) (internal.Seller, error) {
	seller, err := s.rp.GetByID(id)

	if err != nil {
		return internal.Seller{}, err
	}

	if seller == (internal.Seller{}) {
		return internal.Seller{}, utils.ENotFound("Seller")
	}

	return seller, nil
}

func (s *DefaultSellerService) Create(newSeller *internal.Seller) error {
	sellerValidation := s.verify(*newSeller)

	if sellerValidation != nil {
		return sellerValidation
	}

	_, err := s.localityRp.GetByID(newSeller.LocalityID)

	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			return utils.EDependencyNotFound("Seller", "locality ID")
		}

		return err
	}

	err = s.rp.Create(newSeller)

	if err != nil {
		return err
	}

	return nil
}

func (s *DefaultSellerService) Update(id int, newSeller *internal.Seller) (internal.Seller, error) {
	existingSeller, err := s.rp.GetByID(id)
	if err != nil {
		return internal.Seller{}, err
	}

	if existingSeller == (internal.Seller{}) {
		return internal.Seller{}, utils.ENotFound("Seller")
	}

	existingCid, err := s.rp.GetByCid(newSeller.Cid)

	if err != nil {
		return internal.Seller{}, err
	}

	if existingCid.Cid != 0 && existingCid.ID != id {
		return internal.Seller{}, utils.EConflict("Cid", "Seller")
	}

	if newSeller.Cid != 0 {
		existingSeller.Cid = newSeller.Cid
	}

	if len(newSeller.CompanyName) != 0 {
		existingSeller.CompanyName = newSeller.CompanyName
	}

	if len(newSeller.Address) != 0 {
		existingSeller.Address = newSeller.Address
	}

	if len(newSeller.Telephone) != 0 {
		existingSeller.Telephone = newSeller.Telephone
	}

	err = s.rp.Update(&existingSeller)
	if err != nil {
		return internal.Seller{}, err
	}

	return existingSeller, nil
}

func (s *DefaultSellerService) Delete(id int) error {
	existingSeller, err := s.rp.GetByID(id)

	if err != nil {
		return err
	}

	if existingSeller == (internal.Seller{}) {
		return utils.ENotFound("Seller")
	}

	err = s.rp.Delete(id)

	if err != nil {
		return err
	}

	return nil
}

func (s *DefaultSellerService) verify(newSeller internal.Seller) error {
	if newSeller.Cid <= 0 {
		return utils.EZeroValue("Cid")
	}

	existingCid, err := s.rp.GetByCid(newSeller.Cid)
	if err != nil && !errors.Is(err, utils.ErrNotFound) {
		return err
	}

	if existingCid.Cid != 0 {
		return utils.EConflict("Cid", "Seller")
	}

	if len(newSeller.CompanyName) == 0 {
		return utils.EZeroValue("Company name")
	}

	if len(newSeller.Telephone) == 0 {
		return utils.EZeroValue("Telephone")
	}

	if len(newSeller.Address) == 0 {
		return utils.EZeroValue("Address")
	}

	return nil
}
