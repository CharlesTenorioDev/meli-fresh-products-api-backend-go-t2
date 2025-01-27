package carry

import (
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type MySQLCarryService struct {
	repo             internal.CarryRepository
	validateLocality internal.LocalityValidation
}

// NewMySQLCarryService creates a new instance of MySQLCarryService with the provided
// CarryRepository and LocalityValidation. It returns a pointer to the created MySQLCarryService.
//
// Parameters:
//   - repo: an implementation of the CarryRepository interface used for data access.
//   - validateLocality: an implementation of the LocalityValidation interface used for validating localities.
//
// Returns:
//   - A pointer to a MySQLCarryService instance initialized with the provided repository and locality validation.
func NewMySQLCarryService(repo internal.CarryRepository, validateLocality internal.LocalityValidation) *MySQLCarryService {
	return &MySQLCarryService{
		repo:             repo,
		validateLocality: validateLocality,
	}
}

// Save validates the given Carry object and saves it to the repository.
// It first checks if the locality ID exists using the validateLocality service.
// Then, it validates that there are no empty fields in the Carry object.
// If both validations pass, it saves the Carry object using the repository.
//
// Parameters:
//
//	carry (*internal.Carry): The Carry object to be validated and saved.
//
// Returns:
//
//	error: An error if any validation fails or if there is an issue saving the Carry object.
func (s *MySQLCarryService) Save(carry *internal.Carry) error {
	if _, err := s.validateLocality.GetByID(carry.LocalityID); err != nil {
		return utils.ENotFound("Locality")
	}

	if err := s.validateEmptyFields(carry); err != nil {
		return err
	}

	return s.repo.Save(carry)
}

// GetAll retrieves all Carry records from the repository.
// It returns a slice of Carry objects and an error if any issues occur during the retrieval process.
func (s *MySQLCarryService) GetAll() ([]internal.Carry, error) {
	return s.repo.GetAll()
}

// GetByID retrieves a Carry entity by its unique identifier.
// It takes an integer id as a parameter and returns the corresponding Carry entity
// and an error if something goes wrong during the retrieval process.
//
// Parameters:
//   - id: The unique identifier of the Carry entity to be retrieved.
//
// Returns:
//   - internal.Carry: The Carry entity corresponding to the provided id.
//   - error: An error object if there is an issue with the retrieval process, otherwise nil.
func (s *MySQLCarryService) GetByID(id int) (internal.Carry, error) {
	return s.repo.GetByID(id)
}

// Update updates an existing carry record in the database with the provided carry data.
// If any field in the provided carry is empty or zero, it retains the value from the existing carry record.
// It validates the LocalityID if it is provided.
//
// Parameters:
//   - carry: A pointer to the Carry struct containing the updated carry data.
//
// Returns:
//   - error: An error if the update operation fails or if the LocalityID validation fails.
func (s *MySQLCarryService) Update(carry *internal.Carry) error {
	existingCarry, err := s.repo.GetByID(carry.CID)

	if err != nil {
		return err
	}

	s.prepareCarryQuery(carry, existingCarry)

	if carry.LocalityID == 0 {
		(*carry).LocalityID = existingCarry.LocalityID
	} else if _, err := s.validateLocality.GetByID(carry.LocalityID); err != nil {
		return err
	}

	return s.repo.Update(carry)
}

// Delete removes a carry record from the repository based on the provided ID.
// It returns an error if the deletion process fails.
//
// Parameters:
//   - id: The unique identifier of the carry record to be deleted.
//
// Returns:
//   - error: An error object if the deletion fails, otherwise nil.
func (s *MySQLCarryService) Delete(id int) error {
	return s.repo.Delete(id)
}

// validateEmptyFields checks if the required fields in the Carry struct are empty.
// It returns an error if any of the fields CID, CompanyName, Address, or Telephone are empty or zero.
//
// Parameters:
//
//	carry (*internal.Carry): The Carry struct to validate.
//
// Returns:
//
//	error: Returns utils.ErrInvalidArguments if any required field is empty or zero, otherwise returns nil.
func (s *MySQLCarryService) validateEmptyFields(carry *internal.Carry) error {
	if carry.CID == 0 {
		return utils.EZeroValue("CID")
	}

	if carry.CompanyName == "" {
		return utils.EZeroValue("CompanyName")
	}

	if carry.Address == "" {
		return utils.EZeroValue("Address")
	}

	if carry.Telephone == "" {
		return utils.EZeroValue("Telephone")
	}

	return nil
}

// prepareCarryQuery updates the fields of the given carry object with the values
// from the existingCarry object if they are not already set. Specifically, it
// checks if the CID, CompanyName, Address, Telephone, and LocalityID fields of
// the carry object are zero values (e.g., 0 for integers, empty string for strings)
// and if so, assigns the corresponding values from the existingCarry object.
//
// Parameters:
//
//	carry - a pointer to the Carry object that needs to be updated.
//	existingCarry - the Carry object containing the existing values to be used
//	                for updating the carry object.
func (s *MySQLCarryService) prepareCarryQuery(carry *internal.Carry, existingCarry internal.Carry) {
	if carry.CID == 0 {
		(*carry).CID = existingCarry.CID
	}

	if carry.CompanyName == "" {
		(*carry).CompanyName = existingCarry.CompanyName
	}

	if carry.Address == "" {
		(*carry).Address = existingCarry.Address
	}

	if carry.Telephone == "" {
		(*carry).Telephone = existingCarry.Telephone
	}

	if carry.LocalityID == 0 {
		(*carry).LocalityID = existingCarry.LocalityID
	}
}
