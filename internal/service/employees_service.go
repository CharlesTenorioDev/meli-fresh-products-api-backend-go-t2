package service

import (
	pkg "github.com/meli-fresh-products-api-backend-go-t2/internal/pkg"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

// EmployeeDefault is the default implementation of the employee service
// it handles business logic and delegates data operations to the repository
type EmployeeDefault struct {
	rp               pkg.EmployeeRepository
	warehouseService pkg.EmployeesWarehouseValidation
}

// NewEmployeeService creates a new instance of EmployeeDefault
// takes an EmployeeRepository as a parameter to handle data operations
func NewEmployeeService(rp pkg.EmployeeRepository, warehouseService pkg.EmployeesWarehouseValidation) *EmployeeDefault {
	return &EmployeeDefault{rp: rp, warehouseService: warehouseService}
}

// FindAll retrieves all employees from the repository
func (s *EmployeeDefault) FindAll() (employees map[int]pkg.Employee, err error) {
	employees, err = s.rp.FindAll()
	return
}

// FindById retrieves an employee by ID from the repository
func (s *EmployeeDefault) FindById(id int) (employee pkg.Employee, err error) {
	employee, err = s.rp.FindById(id)
	return
}

// CreateEmployee adds a new employee to the repository
func (s *EmployeeDefault) CreateEmployee(newEmployee pkg.EmployeeAttributes) (employee pkg.Employee, err error) {
	// validate required fields
	err = s.validateFields(newEmployee)
	if err != nil {
		return
	}

	// check for duplicates
	employees, _ := s.rp.FindAll()
	err = s.validateDuplicates(employees, newEmployee)
	if err != nil {
		return
	}

	// verify if warehouse_id exists
	err = s.warehouseExistsById(newEmployee.WarehouseId)
	if err != nil {
		return
	}

	// attempt to create the new employee
	return s.rp.CreateEmployee(newEmployee)
}

// UpdateEmployee updates an employee in the repository
func (s *EmployeeDefault) UpdateEmployee(inputEmployee pkg.Employee) (employee pkg.Employee, err error) {
	// find the existing employee
	internalEmployee, err := s.rp.FindById(inputEmployee.ID)
	if err != nil {
		err = utils.ErrNotFound
		return
	}

	// verify if warehouse_id exists
	err = s.warehouseExistsById(inputEmployee.Attributes.WarehouseId)
	if err != nil {
		return
	}

	// merge input fields with the existing employee
	updatedEmployee := mergeEmployeeFields(inputEmployee, internalEmployee)

	// update the employee in the repository
	employee, err = s.rp.UpdateEmployee(updatedEmployee)
	return
}

// DeleteEmployee deletes an employee from the repository based on the provided ID
func (s *EmployeeDefault) DeleteEmployee(id int) (err error) {
	// find the employee to ensure it exists
	employee, err := s.rp.FindById(id)
	if err != nil {
		return utils.ErrNotFound
	}

	// delete the employee by passing only the ID to the repository
	err = s.rp.DeleteEmployee(employee.ID)
	if err != nil {
		return utils.ErrInvalidArguments
	}

	return nil
}

// validateFields checks if the required fields of a new employee are not empty
func (s *EmployeeDefault) validateFields(newEmployee pkg.EmployeeAttributes) (err error) {
	if newEmployee.FirstName == "" || newEmployee.LastName == "" || newEmployee.CardNumberId == 0 {
		return utils.ErrEmptyArguments
	}
	return
}

// validateDuplicates ensures that no existing employee has the same CardNumberId as the new employee
func (s *EmployeeDefault) validateDuplicates(employees map[int]pkg.Employee, newEmployee pkg.EmployeeAttributes) error {
	for _, employee := range employees {
		if employee.Attributes.CardNumberId == newEmployee.CardNumberId {
			return utils.ErrConflict
		}
	}
	return nil
}

// mergeEmployeeFields merges the fields of the input employee with the internal employee
func mergeEmployeeFields(inputEmployee, internalEmployee pkg.Employee) (updatedEmployee pkg.Employee) {
	updatedEmployee.ID = internalEmployee.ID

	if inputEmployee.Attributes.FirstName != "" {
		updatedEmployee.Attributes.FirstName = inputEmployee.Attributes.FirstName
	} else {
		updatedEmployee.Attributes.FirstName = internalEmployee.Attributes.FirstName
	}

	if inputEmployee.Attributes.LastName != "" {
		updatedEmployee.Attributes.LastName = inputEmployee.Attributes.LastName
	} else {
		updatedEmployee.Attributes.LastName = internalEmployee.Attributes.LastName
	}

	if inputEmployee.Attributes.CardNumberId != 0 {
		updatedEmployee.Attributes.CardNumberId = inputEmployee.Attributes.CardNumberId
	} else {
		updatedEmployee.Attributes.CardNumberId = internalEmployee.Attributes.CardNumberId
	}

	if inputEmployee.Attributes.WarehouseId != 0 {
		updatedEmployee.Attributes.WarehouseId = inputEmployee.Attributes.WarehouseId
	} else {
		updatedEmployee.Attributes.WarehouseId = internalEmployee.Attributes.WarehouseId
	}

	return updatedEmployee
}

func (s *EmployeeDefault) warehouseExistsById(id int) error {
	possibleWarehouse, err := s.warehouseService.GetById(id)
	// When internal server error
	if err != nil && err != utils.ErrNotFound {
		return err
	}
	if possibleWarehouse == (pkg.Warehouse{}) {
		return utils.ErrWarehouseDoesNotExists
	}
	return nil
}
