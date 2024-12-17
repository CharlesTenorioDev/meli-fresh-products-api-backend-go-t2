package service

import (
	"github.com/meli-fresh-products-api-backend-go-t2/internal/pkg"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type BasicSectionService struct {
	repo pkg.SectionRepository
}

func NewBasicSectionService(repo pkg.SectionRepository) pkg.SectionService {
	return BasicSectionService{repo}
}

func (r BasicSectionService) GetAll() ([]pkg.Section, error) {
	return r.repo.GetAll()
}

func (r BasicSectionService) GetById(id int) (pkg.Section, error) {
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
func (r BasicSectionService) Save(newSection pkg.Section) (pkg.Section, error) {
	// Zero value validation
	if newSection.WarehouseID == 0 {
		return pkg.Section{}, utils.ErrInvalidArguments
	}
	if newSection.ProductTypeID == 0 {
		return pkg.Section{}, utils.ErrInvalidArguments
	}
	possibleSection, err := r.repo.GetBySectionNumber(newSection.SectionNumber)
	if possibleSection != (pkg.Section{}) {
		return pkg.Section{}, utils.ErrConflict
	}
	if err != nil {
		return pkg.Section{}, err
	}
	newSection, err = r.repo.Save(newSection)
	if err != nil {
		return pkg.Section{}, err
	}

	return newSection, nil
}
