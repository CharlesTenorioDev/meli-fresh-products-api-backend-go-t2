package service

import (
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type MySQLCarryService struct {
	repo             internal.CarryRepository
	validateLocality internal.LocalityValidation
}

func NewMySQLCarryService(repo internal.CarryRepository, validateLocality internal.LocalityValidation) *MySQLCarryService {
	return &MySQLCarryService{
		repo:             repo,
		validateLocality: validateLocality,
	}
}

func (s *MySQLCarryService) Save(carry *internal.Carry) error {
	if _, err := s.validateLocality.GetById(carry.LocalityID); err != nil {
		return err
	}
	if err := s.validateEmptyFields(carry); err != nil {
		return err
	}
	return s.repo.Save(carry)
}

func (s *MySQLCarryService) GetAll() ([]internal.Carry, error) {
	return s.repo.GetAll()
}

func (s *MySQLCarryService) GetById(id int) (internal.Carry, error) {
	return s.repo.GetById(id)
}

func (s *MySQLCarryService) Update(carry *internal.Carry) error {
	return s.repo.Update(carry)
}

func (s *MySQLCarryService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *MySQLCarryService) validateEmptyFields(carry *internal.Carry) error {
	if carry.CID == 0 {
		return utils.ErrInvalidArguments
	}
	if carry.CompanyName == "" {
		return utils.ErrInvalidArguments
	}
	if carry.Address == "" {
		return utils.ErrInvalidArguments
	}
	if carry.Telephone == "" {
		return utils.ErrInvalidArguments
	}
	return nil
}
