package repository

import (
	"cmp"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"slices"

	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type MemorySectionRepository struct {
	db     map[int]internal.Section
	nextId int
}

// Initialize a MemorySectionRepository and set the next id to 1
func NewMemorySectionRepository(load map[int]internal.Section) internal.SectionRepository {
	repo := &MemorySectionRepository{}
	if len(load) == 0 {
		repo.db = make(map[int]internal.Section)
		repo.nextId = 1
	} else {
		repo.db = load
		repo.nextId = utils.GetBiggestId(load) + 1
	}
	return repo
}

// Get all the sections and return in asc order
func (r MemorySectionRepository) GetAll() ([]internal.Section, error) {
	sections := []internal.Section{}
	for _, section := range r.db {
		sections = append(sections, section)
	}
	slices.SortFunc(sections, func(a internal.Section, b internal.Section) int {
		return cmp.Compare(a.ID, b.ID)
	})
	return sections, nil
}

// Return the r.db[id]
func (r *MemorySectionRepository) GetById(id int) (internal.Section, error) {
	return r.db[id], nil
}

// Finds the section by its sectionNumber
// If not found, pkg.Section{} is returned
func (r *MemorySectionRepository) GetBySectionNumber(sectionNumber int) (internal.Section, error) {
	for _, section := range r.db {
		if section.SectionNumber == sectionNumber {
			return section, nil
		}
	}
	return internal.Section{}, nil
}

// Generate a new ID and save the entity
// All validatinos should be made on service layer
func (r *MemorySectionRepository) Save(newSection internal.Section) (internal.Section, error) {
	newSection.ID = r.nextId
	r.nextId++
	r.db[newSection.ID] = newSection
	return r.db[newSection.ID], nil
}

// Update the map
// If a section does not exist for the id, a new one will be created
func (r *MemorySectionRepository) Update(section internal.Section) (internal.Section, error) {
	r.db[section.ID] = section
	return r.db[section.ID], nil
}

// Delete the section by its id
// If no sections exists by id attribute, utils.ErrNotFound is returned
func (r *MemorySectionRepository) Delete(id int) error {
	if r.db[id] == (internal.Section{}) {
		return utils.ErrNotFound
	}
	delete(r.db, id)
	return nil
}
