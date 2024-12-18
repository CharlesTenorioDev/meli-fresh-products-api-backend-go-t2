package service

import (
	"errors"
	"strings"

	"github.com/meli-fresh-products-api-backend-go-t2/internal/pkg"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type WarehouseService struct {
	repo pkg.WarehouseRepository
}

func NewWarehouseService(repo pkg.WarehouseRepository) *WarehouseService {
	return &WarehouseService{repo: repo}
}

func (s *WarehouseService) GetAll() ([]pkg.Warehouse, error) {
	warehouses, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	if len(warehouses) == 0 {
		return nil, utils.ErrNotFound
	}

	return warehouses, nil
}

func (s *WarehouseService) GetById(id int) (pkg.Warehouse, error) {

	warehouse, err := s.repo.GetById(id)
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			return pkg.Warehouse{}, utils.ErrNotFound
		}
		return pkg.Warehouse{}, err
	}

	return warehouse, nil
}

func (s *WarehouseService) Save(newWarehouse pkg.Warehouse) (pkg.Warehouse, error) {
	if err := s.validateWarehouse(newWarehouse); err != nil {
		return pkg.Warehouse{}, err
	}

	if err := s.existingWarehouseCode(newWarehouse, false); err != nil {
		return pkg.Warehouse{}, err
	}

	createdWarehouse, err := s.repo.Save(newWarehouse)
	if err != nil {
		return pkg.Warehouse{}, err
	}

	return createdWarehouse, nil
}

func (s *WarehouseService) existingWarehouseCode(newWarehouse pkg.Warehouse, isUpdate bool) error {
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

func (s *WarehouseService) Update(id int, updatedWarehouse pkg.WarehousePointers) (pkg.Warehouse, error) {
	warehouse, err := s.repo.GetById(id)
	if err != nil {
		return pkg.Warehouse{}, err
	}
	if warehouse == (pkg.Warehouse{}) {
		return pkg.Warehouse{}, utils.ErrNotFound
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
	if updatedWarehouse.MinimumCapacity != nil {
		warehouse.MinimumCapacity = *updatedWarehouse.MinimumCapacity
	}
	if updatedWarehouse.MinimumTemperature != nil {
		warehouse.MinimumTemperature = *updatedWarehouse.MinimumTemperature
	}

	if err := s.validateWarehouse(warehouse); err != nil {
		return pkg.Warehouse{}, err
	}

	if err := s.existingWarehouseCode(warehouse, true); err != nil {
		return pkg.Warehouse{}, err
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

func (s *WarehouseService) validateWarehouse(warehouse pkg.Warehouse) error {
	if warehouse.WarehouseCode == "" {
		return utils.ErrInvalidArguments
	}
	if warehouse.Address == "" {
		return utils.ErrInvalidArguments
	}
	if warehouse.Telephone == "" {
		return utils.ErrInvalidArguments
	}
	if warehouse.MinimumCapacity <= 0 {
		return utils.ErrInvalidArguments
	}
	if warehouse.MinimumTemperature < -273 {
		return utils.ErrInvalidArguments
	}
	return nil
}
