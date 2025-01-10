package service

import (
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
)

type MysqlLocalityService struct {
	Repo internal.LocalityRepository
}

func NewMysqlLocalityService(repo internal.LocalityRepository) *MysqlLocalityService {
	return &MysqlLocalityService{Repo: repo}
}

func (s *MysqlLocalityService) Save(*internal.Locality, *internal.Province, *internal.Country) error {
	return nil
}

func (s *MysqlLocalityService) GetById(id int) (internal.Locality, error) {
	return s.Repo.GetById(id)
}

func (s *MysqlLocalityService) GetSellersByLocalityId(localityId int) ([]internal.SellersByLocality, error) {
	return s.Repo.GetSellersByLocalityId(localityId)
}

func (s *MysqlLocalityService) GetCarriesByLocalityId(localityId int) ([]internal.CarriesByLocality, error) {
	return s.Repo.GetCarriesByLocalityId(localityId)
}
