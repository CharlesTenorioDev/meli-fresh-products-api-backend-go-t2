package repository

import (
	"cmp"
	"slices"

	"github.com/meli-fresh-products-api-backend-go-t2/internal/pkg"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type MemorySectionRepository struct {
	db     map[int]pkg.Section
	nextId int
}

// Initialize a MemorySectionRepository and set the next id to 1
func NewMemorySectionRepository(load map[int]pkg.Section) pkg.SectionRepository {
	repo := &MemorySectionRepository{}
	if len(load) == 0 {
		repo.db = make(map[int]pkg.Section)
		repo.nextId = 1
	} else {
		repo.db = load
		repo.nextId = utils.GetBiggestId(load) + 1
	}
	return repo
}

// Get all the sections and return in asc order
func (r MemorySectionRepository) GetAll() ([]pkg.Section, error) {
	sections := []pkg.Section{}
	for _, section := range r.db {
		sections = append(sections, section)
	}
	slices.SortFunc(sections, func(a pkg.Section, b pkg.Section) int {
		return cmp.Compare(a.ID, b.ID)
	})
	return sections, nil
}

// Return the r.db[id]
func (r *MemorySectionRepository) GetById(id int) (pkg.Section, error) {
	return r.db[id], nil
}

// Finds the section by its sectionNumber
// If not found, pkg.Section{} is returned
func (r *MemorySectionRepository) GetBySectionNumber(sectionNumber int) (pkg.Section, error) {
	for _, section := range r.db {
		if section.SectionNumber == sectionNumber {
			return section, nil
		}
	}
	return pkg.Section{}, nil
}

// Generate a new ID and save the entity
// All validatinos should be made on service layer
func (r *MemorySectionRepository) Save(newSection pkg.Section) (pkg.Section, error) {
	newSection.ID = r.nextId
	r.nextId++
	r.db[newSection.ID] = newSection
	return r.db[newSection.ID], nil
}

// Update the map
// If a section does not exist for the id, a new one will be created
func (r *MemorySectionRepository) Update(section pkg.Section) (pkg.Section, error) {
	r.db[section.ID] = section
	return r.db[section.ID], nil
}

// Delete the section by its id
// If no sections exists by id attribute, utils.ErrNotFound is returned
func (r *MemorySectionRepository) Delete(id int) error {
	if r.db[id] == (pkg.Section{}) {
		return utils.ErrNotFound
	}
	delete(r.db, id)
	return nil
}
