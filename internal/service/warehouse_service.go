package service

import (
	"errors"
	"log"
	"strings"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"

	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type WarehouseService struct {
	repo internal.WarehouseRepository
}

func NewWarehouseService(repo internal.WarehouseRepository) *WarehouseService {
	return &WarehouseService{repo: repo}
}

func (s *WarehouseService) GetAll() ([]internal.Warehouse, error) {
	warehouses, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	if len(warehouses) == 0 {
		return nil, utils.ErrNotFound
	}

	return warehouses, nil
}

func (s *WarehouseService) GetById(id int) (internal.Warehouse, error) {

	warehouse, err := s.repo.GetById(id)
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			return internal.Warehouse{}, utils.ErrNotFound
		}
		return internal.Warehouse{}, err
	}

	return warehouse, nil
}

func (s *WarehouseService) Save(newWarehouse internal.Warehouse) (internal.Warehouse, error) {
	if err := s.validateWarehouse(newWarehouse); err != nil {
		return internal.Warehouse{}, err
	}

	if err := s.existingWarehouseCode(newWarehouse, false); err != nil {
		return internal.Warehouse{}, err
	}

	createdWarehouse, err := s.repo.Save(newWarehouse)
	if err != nil {
		return internal.Warehouse{}, err
	}

	return createdWarehouse, nil
}

func (s *WarehouseService) existingWarehouseCode(newWarehouse internal.Warehouse, isUpdate bool) error {
	existing, _ := s.repo.GetAll()
	for _, w := range existing {
		if !isUpdate {
			if strings.EqualFold(w.WarehouseCode, newWarehouse.WarehouseCode) {
				return utils.ErrConflict
			}
		} else {
			if strings.EqualFold(w.WarehouseCode, newWarehouse.WarehouseCode) && w.ID != newWarehouse.ID {
				return utils.ErrConflict
			}
		}

	}
	return nil
}

func (s *WarehouseService) Update(id int, updatedWarehouse internal.WarehousePointers) (internal.Warehouse, error) {
	warehouse, err := s.repo.GetById(id)
	if err != nil {
		log.Println("error here")
		return internal.Warehouse{}, err
	}
	if warehouse == (internal.Warehouse{}) {
		return internal.Warehouse{}, utils.ErrNotFound
	}
	if updatedWarehouse.Address != nil {
		warehouse.Address = *updatedWarehouse.Address
	}
	if updatedWarehouse.Telephone != nil {
		warehouse.Telephone = *updatedWarehouse.Telephone
	}
	if updatedWarehouse.WarehouseCode != nil {
		warehouse.WarehouseCode = *updatedWarehouse.WarehouseCode
	}
	if updatedWarehouse.LocalityID != nil {
		warehouse.LocalityID = *updatedWarehouse.LocalityID
	}
	if updatedWarehouse.MinimumCapacity != nil {
		warehouse.MinimumCapacity = *updatedWarehouse.MinimumCapacity
	}
	if updatedWarehouse.MinimumTemperature != nil {
		warehouse.MinimumTemperature = *updatedWarehouse.MinimumTemperature
	}

	if err := s.validateWarehouse(warehouse); err != nil {
		return internal.Warehouse{}, err
	}

	if err := s.existingWarehouseCode(warehouse, true); err != nil {
		return internal.Warehouse{}, err
	}

	warehouse, err = s.repo.Update(warehouse)
	return warehouse, nil
}

func (s *WarehouseService) Delete(id int) error {
	_, err := s.repo.GetById(id)

	if err != nil {
		return err
	}

	if err := s.repo.Delete(id); err != nil {
		return err
	}

	return nil
}

func (s *WarehouseService) validateWarehouse(warehouse internal.Warehouse) error {
	if warehouse.WarehouseCode == "" {
		return utils.ErrInvalidArguments
	}
	if warehouse.Address == "" {
		return utils.ErrInvalidArguments
	}
	if warehouse.Telephone == "" {
		return utils.ErrInvalidArguments
	}
	// TODO: Add more validations to localities
	return nil
}
