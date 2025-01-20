package service

import (
	"errors"
	"strings"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"

	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type WarehouseService struct {
	repo             internal.WarehouseRepository
	validateLocality internal.LocalityValidation
}

// NewWarehouseService creates a new instance of WarehouseService with the provided WarehouseRepository.
// It returns a pointer to the created WarehouseService.
//
// Parameters:
//   - repo: an implementation of the WarehouseRepository interface.
//
// Returns:
//   - A pointer to the newly created WarehouseService.
func NewWarehouseService(repo internal.WarehouseRepository, validateLocality internal.LocalityValidation) *WarehouseService {
	return &WarehouseService{
		repo:             repo,
		validateLocality: validateLocality,
	}
}

// GetAll retrieves all warehouses from the repository.
// It returns a slice of Warehouse and an error if any occurred during the retrieval process.
// If no warehouses are found, it returns an ErrNotFound error.
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

// GetByID retrieves a warehouse by its ID.
// If the warehouse is not found, it returns an ErrNotFound error.
// If any other error occurs during the retrieval, it returns that error.
// Parameters:
//   - id: the ID of the warehouse to retrieve.
//
// Returns:
//   - internal.Warehouse: the warehouse with the specified ID.
//   - error: an error if the warehouse is not found or if any other error occurs.
func (s *WarehouseService) GetByID(id int) (internal.Warehouse, error) {
	warehouse, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			return internal.Warehouse{}, utils.ErrNotFound
		}

		return internal.Warehouse{}, err
	}

	return warehouse, nil
}

// Save validates and saves a new warehouse to the repository.
// It first validates the warehouse data, then checks for existing warehouse codes.
// If the warehouse is valid and the code is unique, it saves the warehouse to the repository.
// Returns the created warehouse and any error encountered during the process.
//
// Parameters:
//
//	newWarehouse - the warehouse data to be saved.
//
// Returns:
//
//	internal.Warehouse - the saved warehouse data.
//	error - any error encountered during the validation or saving process.
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

// existingWarehouseCode checks if a warehouse code already exists in the repository.
// If `isUpdate` is false, it checks for any warehouse with the same code.
// If `isUpdate` is true, it checks for any warehouse with the same code but a different ID.
// Returns an error if a conflict is found, otherwise returns nil.
//
// Parameters:
//   - newWarehouse: the warehouse to check for conflicts.
//   - isUpdate: a flag indicating if the operation is an update.
//
// Returns:
//   - error: utils.ErrConflict if a conflict is found, otherwise nil.
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

// Update updates an existing warehouse with the provided updatedWarehouse data.
// It retrieves the warehouse by its ID, applies the updates, validates the warehouse,
// checks for existing warehouse codes, and then saves the updated warehouse.
//
// Parameters:
//   - id: The ID of the warehouse to be updated.
//   - updatedWarehouse: A struct containing pointers to the fields to be updated.
//
// Returns:
//   - internal.Warehouse: The updated warehouse.
//   - error: An error if the update process fails, or nil if successful.
func (s *WarehouseService) Update(id int, updatedWarehouse internal.WarehousePointers) (internal.Warehouse, error) {
	warehouse, err := s.repo.GetByID(id)
	if err != nil {
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

	if err != nil {
		return internal.Warehouse{}, err
	}

	return warehouse, nil
}

// Delete removes a warehouse entry from the repository by its ID.
// It first checks if the warehouse exists by calling GetByID method.
// If the warehouse does not exist or any error occurs during the check, it returns the error.
// If the warehouse exists, it proceeds to delete it by calling the Delete method of the repository.
// If any error occurs during the deletion, it returns the error.
// Returns nil if the deletion is successful.
//
// Parameters:
//   - id: the ID of the warehouse to be deleted.
//
// Returns:
//   - error: an error if the warehouse does not exist or if there is an issue during deletion, otherwise nil.
func (s *WarehouseService) Delete(id int) error {
	_, err := s.repo.GetByID(id)

	if err != nil {
		return err
	}

	if err := s.repo.Delete(id); err != nil {
		return err
	}

	return nil
}

// validateWarehouse validates the given warehouse object.
// It checks if the WarehouseCode, Address, and Telephone fields are not empty.
// If any of these fields are empty, it returns an ErrInvalidArguments error.
// TODO: Add more validations to localities.
//
// Parameters:
//   - warehouse: the warehouse object to be validated.
//
// Returns:
//   - error: an error if any validation fails, otherwise nil.
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

	if warehouse.MinimumCapacity < 0 {
		return utils.ErrInvalidArguments
	}

	if warehouse.MinimumTemperature < -273 {
		return utils.ErrInvalidArguments
	}

	if _, err := s.validateLocality.GetByID(warehouse.LocalityID); err != nil {
		return err
	}

	return nil
}
