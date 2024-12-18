package service

import (
	"errors"
	"fmt"

	"github.com/meli-fresh-products-api-backend-go-t2/internal/pkg"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type BasicSectionService struct {
	repo        pkg.SectionRepository
	validations pkg.SectionValidations
}

func NewBasicSectionService(repo pkg.SectionRepository, validations pkg.SectionValidations) pkg.SectionService {
	return &BasicSectionService{
		repo:        repo,
		validations: validations,
	}
}

// Returns all the sections
func (r *BasicSectionService) GetAll() ([]pkg.Section, error) {
	return r.repo.GetAll()
}

func (r *BasicSectionService) GetById(id int) (pkg.Section, error) {
	// Check if section exists
	possibleSection, err := r.repo.GetById(id)
	if err != nil {
		return pkg.Section{}, err
	}

	// If does not exists, 404 error
	if possibleSection == (pkg.Section{}) {
		return pkg.Section{}, utils.ErrNotFound
	}

	return possibleSection, nil
}

func (r *BasicSectionService) Save(newSection pkg.Section) (pkg.Section, error) {
	// Zero value validation
	if newSection.SectionNumber == 0 {
		return pkg.Section{}, errors.Join(utils.ErrInvalidArguments, errors.New("section_number cannot be empty/null"))
	}
	if newSection.WarehouseID == 0 {
		return pkg.Section{}, errors.Join(utils.ErrInvalidArguments, errors.New("warehouse_id cannot be empty/null"))
	}
	if newSection.ProductTypeID == 0 {
		return pkg.Section{}, errors.Join(utils.ErrInvalidArguments, errors.New("product_type_id cannot be empty/null"))
	}

	// ToDo: Implement the warehouse_id validation and product_type_id validation
	if !r.validations.WarehouseExistsById(newSection.WarehouseID) {
		return pkg.Section{}, errors.Join(utils.ErrNotFound, fmt.Errorf("warehouse_id not found for id %d", newSection.WarehouseID))
	}
	if !r.validations.ProductTypeExistsById(newSection.ProductTypeID) {
		return pkg.Section{}, errors.Join(utils.ErrNotFound, fmt.Errorf("product_type_id not found for id %d", newSection.ProductTypeID))
	}

	if newSection.MinimumCapacity > newSection.MaximumCapacity {
		return pkg.Section{}, errors.Join(utils.ErrInvalidArguments, errors.New("minimum_capacity cannot be greater than maximum_capacity"))
	}

	// Check if a section already exists for section number
	possibleSection, err := r.repo.GetBySectionNumber(newSection.SectionNumber)
	if possibleSection != (pkg.Section{}) {
		return pkg.Section{}, utils.ErrConflict
	}
	if err != nil {
		return pkg.Section{}, err
	}

	// Save if ok
	newSection, err = r.repo.Save(newSection)
	if err != nil {
		return pkg.Section{}, err
	}

	return newSection, nil
}

func (r *BasicSectionService) Update(id int, sectionToUpdate pkg.SectionPointers) (pkg.Section, error) {
	section, err := r.repo.GetById(id)
	if err != nil {
		return pkg.Section{}, err
	}
	// If does not exists, 404 error
	if section == (pkg.Section{}) {
		return pkg.Section{}, utils.ErrNotFound
	}

	// Check which field will be updated
	if sectionToUpdate.SectionNumber != nil {
		section.SectionNumber = *sectionToUpdate.SectionNumber
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
	}
	if sectionToUpdate.WarehouseID != nil {
		section.WarehouseID = *sectionToUpdate.WarehouseID
	}

	// Update
	section, err = r.repo.Update(section)
	return section, nil
}

func (r *BasicSectionService) Delete(id int) error {
	possibleSection, err := r.repo.GetById(id)
	if err != nil {
		return err
	}
	if possibleSection == (pkg.Section{}) {
		return utils.ErrNotFound
	}
	err = r.repo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
