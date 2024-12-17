package repository

import (
	"slices"

	"github.com/meli-fresh-products-api-backend-go-t2/internal/pkg"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type MemorySectionRepository struct {
	db     map[int]pkg.Section
	nextId int
}

func NewMemorySectionRepository() pkg.SectionRepository {
	repo := &MemorySectionRepository{}
	repo.db = make(map[int]pkg.Section)
	repo.nextId = 1
	return repo
}

func (r MemorySectionRepository) GetAll() ([]pkg.Section, error) {
	sections := []pkg.Section{}
	for _, section := range r.db {
		sections = append(sections, section)
	}
	slices.SortFunc(sections, func(a pkg.Section, b pkg.Section) int {
		if a.ID > b.ID {
			return 0
		}
		return 1
	})
	return sections, nil
}

func (r *MemorySectionRepository) GetById(id int) (pkg.Section, error) {
	return r.db[id], nil
}

func (r *MemorySectionRepository) GetBySectionNumber(sectionNumber int) (pkg.Section, error) {
	for _, section := range r.db {
		if section.SectionNumber == sectionNumber {
			return section, nil
		}
	}
	return pkg.Section{}, nil
}

func (r *MemorySectionRepository) Save(newSection pkg.Section) (pkg.Section, error) {
	newSection.ID = r.nextId
	r.nextId++
	r.db[newSection.ID] = newSection
	return r.db[newSection.ID], nil
}

func (r *MemorySectionRepository) Delete(id int) error {
	if r.db[id] == (pkg.Section{}) {
		return utils.ErrNotFound
	}
	delete(r.db, id)
	return nil
}
