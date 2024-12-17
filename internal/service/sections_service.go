package service

import "github.com/meli-fresh-products-api-backend-go-t2/internal/pkg"

type BasicSectionService struct {
	repo pkg.SectionRepository
}

func NewBasicSectionService(repo pkg.SectionRepository) pkg.SectionService {
	return BasicSectionService{repo}
}

func (r BasicSectionService) GetAll() ([]pkg.Section, error) {
	return r.repo.GetAll()
}
