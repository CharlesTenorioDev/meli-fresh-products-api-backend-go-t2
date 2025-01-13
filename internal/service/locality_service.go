package service

import (
	"errors"
	"fmt"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type BasicLocalityService struct {
	localityRepo internal.LocalityRepository
	provinceRepo internal.ProvinceRepository
	countryRepo  internal.CountryRepository
}

func NewBasicLocalityService(
	lr internal.LocalityRepository,
	pr internal.ProvinceRepository,
	cr internal.CountryRepository) internal.LocalityService {
	return &BasicLocalityService{
		lr, pr, cr,
	}
}

func (s *BasicLocalityService) Save(locality *internal.Locality, province *internal.Province, country *internal.Country) error {
	if locality.LocalityName == "" {
		return utils.ErrInvalidArguments
	}
	if locality.ID == 0 {
		return utils.ErrInvalidArguments
	}
	if province.ProvinceName == "" {
		return utils.ErrInvalidArguments
	}
	if country.CountryName == "" {
		return utils.ErrInvalidArguments
	}
	// If locality exists by id
	possibleLocality, err := s.localityRepo.GetById(locality.ID)
	if err != nil && !errors.Is(err, utils.ErrNotFound) {
		return err
	}
	if possibleLocality != (internal.Locality{}) {
		return utils.ErrConflict
	}

	// Check if we find a country by its name
	possibleCountry, err := s.countryRepo.GetByName(country.CountryName)
	if err != nil {
		// If not exists, we create on
		if errors.Is(err, utils.ErrNotFound) {
			if err = s.countryRepo.Save(country); err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		country.ID = possibleCountry.ID
	}
	(*province).CountryID = country.ID

	// Check if we find a province by its name
	possibleProvince, err := s.provinceRepo.GetByName(province.ProvinceName)
	if err != nil {
		// If not exists, we create on
		if errors.Is(err, utils.ErrNotFound) {
			if err = s.provinceRepo.Save(province); err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		province.ID = possibleProvince.ID
	}
	(*locality).ProvinceID = province.ID
	fmt.Printf("%+v %+v %+v\n\n", locality, province, country)
	if err = s.localityRepo.Save(locality); err != nil {
		return err
	}

	return nil
}

// Get the sellers quantity by ther location
// if id == 0, then all location are returned
func (s *BasicLocalityService) GetSellersByLocalityId(localityId int) ([]internal.SellersByLocality, error) {
	// Of id != 0, check if locality exists
	if localityId != 0 {
		_, err := s.localityRepo.GetById(localityId)
		if err != nil {
			return []internal.SellersByLocality{}, err
		}
	}
	return s.localityRepo.GetSellersByLocalityId(localityId)
}

func (s *BasicLocalityService) GetCarriesByLocalityId(localityId int) ([]internal.CarriesByLocality, error) {
	// Of id != 0, check if locality exists
	if localityId != 0 {
		_, err := s.localityRepo.GetById(localityId)
		if err != nil {
			return []internal.CarriesByLocality{}, err
		}
	}
	return s.localityRepo.GetCarriesByLocalityId(localityId)
}
