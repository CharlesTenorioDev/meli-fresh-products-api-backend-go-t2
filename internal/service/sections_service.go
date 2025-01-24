package service

import (
	"errors"
	"fmt"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"

	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

const (
	MinCelsiusTemperature = -273.15
)

type BasicSectionService struct {
	repo               internal.SectionRepository
	warehouseService   internal.SectionWarehouseValidation
	productTypeService internal.SectionProductTypeValidation
}

func NewBasicSectionService(repo internal.SectionRepository, warehouseService internal.SectionWarehouseValidation, productTypeService internal.SectionProductTypeValidation) internal.SectionService {
	return &BasicSectionService{
		repo:               repo,
		warehouseService:   warehouseService,
		productTypeService: productTypeService,
	}
}

// Returns all the sections
func (s BasicSectionService) GetAll() ([]internal.Section, error) {
	return s.repo.GetAll()
}

// Get the section by id, if sections does not exist, utils.ErrNotFound is returned
func (s BasicSectionService) GetById(id int) (internal.Section, error) {
	// Check if section exists
	possibleSection, err := s.repo.GetById(id)
	if err != nil {
		return internal.Section{}, err
	}

	// If does not exists, 404 error
	if possibleSection == (internal.Section{}) {
		return internal.Section{}, utils.ErrNotFound
	}

	return possibleSection, nil
}

func (s *BasicSectionService) warehouseExistsById(id int) error {
	possibleWarehouse, err := s.warehouseService.GetById(id)
	// When internal server error
	if err != nil && !errors.Is(err, utils.ErrNotFound) {
		return err
	}
	if possibleWarehouse == (internal.Warehouse{}) {
		return errors.Join(utils.ErrInvalidArguments, fmt.Errorf("warehouse not found for id %d", id))
	}
	return nil
}

func (s *BasicSectionService) productTypeExistsById(id int) error {
	possibleProductType, err := s.productTypeService.GetProductTypeByID(id)
	// When internal server error
	if err != nil && !errors.Is(err, utils.ErrNotFound) {
		return err
	}
	if possibleProductType == (internal.ProductType{}) {
		return errors.Join(utils.ErrInvalidArguments, fmt.Errorf("product_type not found for id %d", id))
	}
	return nil
}

func (s *BasicSectionService) sectionExistsBySectionNumber(sectionNumber int) error {
	possibleSection, err := s.repo.GetBySectionNumber(sectionNumber)
	if possibleSection != (internal.Section{}) {
		return utils.ErrConflict
	}
	if err != nil && !errors.Is(err, utils.ErrNotFound) {
		return err
	}
	return nil
}

func (s *BasicSectionService) validateLogicRules(section internal.Section) error {
	if section.MinimumCapacity > section.MaximumCapacity {
		return errors.Join(utils.ErrInvalidArguments, errors.New("minimum_capacity cannot be greater than maximum_capacity"))
	}
	if section.MinimumTemperature < MinCelsiusTemperature {
		return errors.Join(utils.ErrInvalidArguments, errors.New("minimum_temperature cannot be less than -273.15 Celsius"))
	}
	if section.CurrentTemperature < MinCelsiusTemperature {
		return errors.Join(utils.ErrInvalidArguments, errors.New("current_temperature cannot be less than -273.15 Celsius"))
	}
	return nil
}

// Save a section, check the relations, zero value when applicable, and basic logic
func (s *BasicSectionService) Save(newSection internal.Section) (internal.Section, error) {
	// Zero value validation
	if newSection.SectionNumber <= 0 {
		return internal.Section{}, errors.Join(utils.ErrInvalidArguments, errors.New("section_number cannot be empty/null"))
	}
	if newSection.WarehouseID <= 0 {
		return internal.Section{}, errors.Join(utils.ErrInvalidArguments, errors.New("warehouse_id cannot be empty/null"))
	}
	if newSection.ProductTypeID <= 0 {
		return internal.Section{}, errors.Join(utils.ErrInvalidArguments, errors.New("product_type_id cannot be empty/null"))
	}

	if err := s.warehouseExistsById(newSection.WarehouseID); err != nil {
		return internal.Section{}, err
	}
	if err := s.productTypeExistsById(newSection.ProductTypeID); err != nil {
		return internal.Section{}, err
	}
	if err := s.validateLogicRules(newSection); err != nil {
		return internal.Section{}, err
	}
	if err := s.sectionExistsBySectionNumber(newSection.SectionNumber); err != nil {
		return internal.Section{}, err
	}

	// Save if ok
	newSection, err := s.repo.Save(&newSection)
	if err != nil {
		return internal.Section{}, err
	}

	return newSection, nil
}

func (s *BasicSectionService) Update(id int, sectionToUpdate internal.SectionPointers) (internal.Section, error) {
	section, err := s.repo.GetById(id)
	if err != nil {
		return internal.Section{}, err
	}
	// If does not exists, 404 error
	if section == (internal.Section{}) {
		return internal.Section{}, utils.ErrNotFound
	}

	// Check which field will be updated
	if sectionToUpdate.SectionNumber != nil && *sectionToUpdate.SectionNumber != section.SectionNumber {
		section.SectionNumber = *sectionToUpdate.SectionNumber
		if section.SectionNumber == 0 {
			return internal.Section{}, errors.Join(utils.ErrInvalidArguments, errors.New("section_number cannot be empty/null"))
		}

		if err := s.sectionExistsBySectionNumber(section.SectionNumber); err != nil {
			return internal.Section{}, err
		}
	}
	if sectionToUpdate.CurrentCapacity != nil {
		section.CurrentCapacity = *sectionToUpdate.CurrentCapacity
	}
	if sectionToUpdate.MaximumCapacity != nil {
		section.MaximumCapacity = *sectionToUpdate.MaximumCapacity
	}
	if sectionToUpdate.MinimumCapacity != nil {
		section.MinimumCapacity = *sectionToUpdate.MinimumCapacity
	}
	if sectionToUpdate.CurrentTemperature != nil {
		section.CurrentTemperature = *sectionToUpdate.CurrentTemperature
	}
	if sectionToUpdate.MinimumTemperature != nil {
		section.MinimumTemperature = *sectionToUpdate.MinimumTemperature
	}
	if sectionToUpdate.ProductTypeID != nil {
		section.ProductTypeID = *sectionToUpdate.ProductTypeID
		if section.ProductTypeID == 0 {
			return internal.Section{}, errors.Join(utils.ErrInvalidArguments, errors.New("product_type_id cannot be empty/null"))
		}
		if err := s.productTypeExistsById(section.ProductTypeID); err != nil {
			return internal.Section{}, err
		}
	}
	if sectionToUpdate.WarehouseID != nil {
		section.WarehouseID = *sectionToUpdate.WarehouseID
		if section.WarehouseID == 0 {
			return internal.Section{}, errors.Join(utils.ErrInvalidArguments, errors.New("warehouse_id cannot be empty/null"))
		}
		if err := s.warehouseExistsById(section.WarehouseID); err != nil {
			return internal.Section{}, err
		}
	}
	if err := s.validateLogicRules(section); err != nil {
		return internal.Section{}, err
	}
	// Update
	section, err = s.repo.Update(&section)
	if err != nil {
		return internal.Section{}, err
	}
	return section, nil
}

func (s *BasicSectionService) Delete(id int) error {
	possibleSection, err := s.repo.GetById(id)
	if err != nil {
		return err
	}
	if possibleSection == (internal.Section{}) {
		return utils.ErrNotFound
	}
	err = s.repo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *BasicSectionService) GetSectionProductsReport(id int) ([]internal.SectionProductsReport, error) {
	var report []internal.SectionProductsReport
	var err error

	if id == 0 {
		report, err = s.repo.GetSectionProductsReport()
		if err != nil {
			return nil, err
		}
		return report, nil

	} else {
		sectionExists, err := s.repo.GetById(id)
		if err != nil {
			return nil, err
		}
		if sectionExists == (internal.Section{}) {
			return nil, utils.ErrNotFound
		}
		report, err = s.repo.GetSectionProductsReportById(id)
		if err != nil {
			return nil, err
		}

		return report, nil
	}
}
