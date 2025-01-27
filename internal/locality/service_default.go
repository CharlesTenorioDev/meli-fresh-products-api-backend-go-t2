package locality

import (
	"errors"
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

// Save validates and saves a locality, province, and country to the repository.
// It first checks if the locality, province, and country have valid names and IDs.
// If the locality already exists, it returns an ErrConflict error.
// If the country does not exist, it creates a new country entry.
// If the province does not exist, it creates a new province entry.
// Finally, it saves the locality with the associated province and country IDs.
//
// Parameters:
//
//	locality - a pointer to the Locality struct to be saved
//	province - a pointer to the Province struct to be saved
//	country - a pointer to the Country struct to be saved
//
// Returns:
//
//	error - an error if any validation or save operation fails, otherwise nil
func (s *BasicLocalityService) Save(locality *internal.Locality, province *internal.Province, country *internal.Country) error {
	if locality.LocalityName == "" {
		return utils.EZeroValue("locality_name")
	}

	if locality.ID == 0 {
		return utils.EZeroValue("id")
	}

	if province.ProvinceName == "" {
		return utils.EZeroValue("province_name")
	}

	if country.CountryName == "" {
		return utils.EZeroValue("country_name")
	}

	// If locality exists by id
	// Check for error 500
	possibleLocality, err := s.localityRepo.GetByID(locality.ID)
	if err != nil && !errors.Is(err, utils.ErrNotFound) {
		return err
	}

	if possibleLocality != (internal.Locality{}) {
		return utils.EConflict("id", "locality")
	}

	// Check if we find a country by its name
	possibleCountry, err := s.countryRepo.GetByName(country.CountryName)

	// We find the country
	if err == nil {
		country.ID = possibleCountry.ID
	} else if errors.Is(err, utils.ErrNotFound) { // We need to create the country
		if err = s.countryRepo.Save(country); err != nil {
			// Internal error
			return err
		}
	} else { // Internal error
		return err
	}

	(*province).CountryID = country.ID

	// Check if we find a province by its name
	possibleProvince, err := s.provinceRepo.GetByName(province.ProvinceName)

	if err == nil {
		province.ID = possibleProvince.ID
	} else if errors.Is(err, utils.ErrNotFound) {
		if err = s.provinceRepo.Save(province); err != nil {
			// Internal error
			return err
		}
	} else { // Internal error
		return err
	}

	(*locality).ProvinceID = province.ID

	if err = s.localityRepo.Save(locality); err != nil {
		return err
	}

	return nil
}

// GetSellersByLocalityID the sellers quantity by there location
// if id == 0, then all location are returned
func (s *BasicLocalityService) GetSellersByLocalityID(localityID int) ([]internal.SellersByLocality, error) {
	if localityID <= 0 {
		return []internal.SellersByLocality{}, utils.EZeroValue("locality_id")
	}

	if _, err := s.localityRepo.GetByID(localityID); err != nil {
		return []internal.SellersByLocality{}, err
	}

	return s.localityRepo.GetSellersByLocalityID(localityID)
}

// GetCarriesByLocalityID retrieves a list of carriers associated with a given locality ID.
// If the provided locality ID is not zero, it first checks if the locality exists.
// If the locality does not exist, it returns an error.
// If the locality exists or the locality ID is zero, it returns the list of carriers.
//
// Parameters:
//   - localityId: The ID of the locality to retrieve carriers for.
//
// Returns:
//   - []internal.CarriesByLocality: A slice of carriers associated with the locality.
//   - error: An error if the locality does not exist or if there is an issue retrieving the carriers.
func (s *BasicLocalityService) GetCarriesByLocalityID(localityID int) ([]internal.CarriesByLocality, error) {
	if localityID <= 0 {
		return []internal.CarriesByLocality{}, utils.EZeroValue("locality_id")
	}

	if _, err := s.localityRepo.GetByID(localityID); err != nil {
		return []internal.CarriesByLocality{}, err
	}

	return s.localityRepo.GetCarriesByLocalityID(localityID)
}
