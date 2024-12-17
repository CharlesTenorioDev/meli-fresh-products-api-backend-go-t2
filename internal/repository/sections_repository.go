package repository

import (
	"slices"

	"github.com/meli-fresh-products-api-backend-go-t2/internal/pkg"
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
