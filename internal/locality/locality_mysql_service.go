package locality

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

func (s *MysqlLocalityService) GetByID(id int) (internal.Locality, error) {
	return s.Repo.GetByID(id)
}

func (s *MysqlLocalityService) GetSellersByLocalityID(localityID int) ([]internal.SellersByLocality, error) {
	return s.Repo.GetSellersByLocalityID(localityID)
}

func (s *MysqlLocalityService) GetCarriesByLocalityID(localityID int) ([]internal.CarriesByLocality, error) {
	return s.Repo.GetCarriesByLocalityID(localityID)
}
