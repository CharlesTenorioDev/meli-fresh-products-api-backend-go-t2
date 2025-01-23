package section

import (
	"errors"
	"strconv"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"

	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

const (
	MinCelsiusTemperature = -273.15
)

type DefaultSectionService struct {
	repo               internal.SectionRepository
	warehouseService   internal.SectionWarehouseValidation
	productTypeService internal.SectionProductTypeValidation
}

func NewBasicSectionService(repo internal.SectionRepository, warehouseService internal.SectionWarehouseValidation, productTypeService internal.SectionProductTypeValidation) internal.SectionService {
	return &DefaultSectionService{
		repo:               repo,
		warehouseService:   warehouseService,
		productTypeService: productTypeService,
	}
}

// GetAll Returns all the sections
func (s DefaultSectionService) GetAll() ([]internal.Section, error) {
	return s.repo.GetAll()
}

// GetByID Get the section by id, if sections does not exist, utils.ErrNotFound is returned
func (s DefaultSectionService) GetByID(id int) (internal.Section, error) {
	// Check if section exists
	possibleSection, err := s.repo.GetByID(id)

	if err != nil {
		return internal.Section{}, err
	}

	// If does not exists, 404 error
	if possibleSection == (internal.Section{}) {
		return internal.Section{}, utils.ENotFound("section")
	}

	return possibleSection, nil
}

func (s *DefaultSectionService) warehouseExistsByID(id int) error {
	possibleWarehouse, err := s.warehouseService.GetByID(id)
	// When internal server error
	if err != nil && !errors.Is(err, utils.ErrNotFound) {
		return err
	}

	if possibleWarehouse == (internal.Warehouse{}) {
		return utils.EDependencyNotFound("warehouse", "id: "+strconv.Itoa(id))
	}

	return nil
}

func (s *DefaultSectionService) productTypeExistsByID(id int) error {
	possibleProductType, err := s.productTypeService.GetProductTypeByID(id)
	// When internal server error
	if err != nil && !errors.Is(err, utils.ErrNotFound) {
		return err
	}

	if possibleProductType == (internal.ProductType{}) {
		return utils.EDependencyNotFound("product_type", "id: "+strconv.Itoa(id))
	}

	return nil
}

func (s *DefaultSectionService) sectionExistsBySectionNumber(sectionNumber int) error {
	possibleSection, err := s.repo.GetBySectionNumber(sectionNumber)
	if possibleSection != (internal.Section{}) {
		return utils.EConflict("section", "id: "+strconv.Itoa(sectionNumber))
	}

	if err != nil && !errors.Is(err, utils.ErrNotFound) {
		return err
	}

	return nil
}

func (s *DefaultSectionService) validateLogicRules(section internal.Section) error {
	if section.MinimumCapacity > section.MaximumCapacity {
		return utils.EBR("minimum_capacity cannot be greater than maximum_capacity")
	}

	if section.MinimumTemperature < MinCelsiusTemperature {
		return utils.EBR("minimum_temperature cannot be less than -273.15 Celsius")
	}

	if section.CurrentTemperature < MinCelsiusTemperature {
		return utils.EBR("current_temperature cannot be less than -273.15 Celsius")
	}

	return nil
}

// Save a section, check the relations, zero value when applicable, and basic logic
func (s *DefaultSectionService) Save(newSection internal.Section) (internal.Section, error) {
	// Zero value validation
	if newSection.SectionNumber <= 0 {
		return internal.Section{}, utils.EZeroValue("section_number")
	}

	if newSection.WarehouseID <= 0 {
		return internal.Section{}, utils.EZeroValue("warehouse_id")
	}

	if newSection.ProductTypeID <= 0 {
		return internal.Section{}, utils.EZeroValue("product_type_id")
	}

	if err := s.warehouseExistsByID(newSection.WarehouseID); err != nil {
		return internal.Section{}, err
	}

	if err := s.productTypeExistsByID(newSection.ProductTypeID); err != nil {
		return internal.Section{}, err
	}

	if err := s.validateLogicRules(newSection); err != nil {
		return internal.Section{}, err
	}

	if err := s.sectionExistsBySectionNumber(newSection.SectionNumber); err != nil {
		return internal.Section{}, err
	}

	// Save if ok
	err := s.repo.Save(&newSection)
	if err != nil {
		return internal.Section{}, err
	}

	return newSection, nil
}

func (s *DefaultSectionService) Update(id int, sectionToUpdate internal.SectionPointers) (internal.Section, error) {
	section, err := s.repo.GetByID(id)

	if err != nil && !errors.Is(err, utils.ErrNotFound) {
		return internal.Section{}, err
	}

	// If does not exists, 404 error
	if err != nil && errors.Is(err, utils.ErrNotFound) {
		return internal.Section{}, utils.ENotFound("section")
	}

	// Check which field will be updated
	if sectionToUpdate.SectionNumber != nil && *sectionToUpdate.SectionNumber != section.SectionNumber {
		section.SectionNumber = *sectionToUpdate.SectionNumber
		if section.SectionNumber <= 0 {
			return internal.Section{}, utils.EZeroValue("section_number")
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
			return internal.Section{}, utils.EZeroValue("product_type_id")
		}

		if err := s.productTypeExistsByID(section.ProductTypeID); err != nil {
			return internal.Section{}, err
		}
	}

	if sectionToUpdate.WarehouseID != nil {
		section.WarehouseID = *sectionToUpdate.WarehouseID
		if section.WarehouseID == 0 {
			return internal.Section{}, errors.Join(utils.ErrInvalidArguments, errors.New("warehouse_id cannot be empty/null"))
		}

		if err := s.warehouseExistsByID(section.WarehouseID); err != nil {
			return internal.Section{}, err
		}
	}

	if err := s.validateLogicRules(section); err != nil {
		return internal.Section{}, err
	}

	// Update
	err = s.repo.Update(&section)

	if err != nil {
		return internal.Section{}, err
	}

	return section, nil
}

func (s *DefaultSectionService) Delete(id int) error {
	possibleSection, err := s.repo.GetByID(id)

	if err != nil {
		return err
	}

	if possibleSection == (internal.Section{}) {
		return utils.ENotFound("section")
	}

	err = s.repo.Delete(id)

	if err != nil {
		return err
	}

	return nil
}

func (s *DefaultSectionService) GetSectionProductsReport(id int) ([]internal.SectionProductsReport, error) {
	var report []internal.SectionProductsReport

	var err error

	if id == 0 {
		report, err = s.repo.GetSectionProductsReport()

		if err != nil {
			return nil, err
		}

		return report, nil
	} else {
		sectionExists, err := s.repo.GetByID(id)
		if err != nil {
			return nil, err
		}

		if sectionExists == (internal.Section{}) {
			return nil, utils.ENotFound("section")
		}

		report, err = s.repo.GetSectionProductsReportByID(id)

		if err != nil {
			return nil, err
		}

		return report, nil
	}
}
