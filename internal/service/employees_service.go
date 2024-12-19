package service

import (
	employeesPkg "github.com/meli-fresh-products-api-backend-go-t2/internal/pkg"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

// EmployeeDefault is the default implementation of the employee service
// it handles business logic and delegates data operations to the repository
type EmployeeDefault struct {
	rp employeesPkg.EmployeeRepository
}

// NewEmployeeService creates a new instance of EmployeeDefault
// takes an EmployeeRepository as a parameter to handle data operations
func NewEmployeeService(rp employeesPkg.EmployeeRepository) *EmployeeDefault {
	return &EmployeeDefault{rp: rp}
}

// FindAll retrieves all employees from the repository
func (s *EmployeeDefault) FindAll() (employees map[int]employeesPkg.Employee, err error) {
	employees, err = s.rp.FindAll()
	return
}

// FindById retrieves an employee by ID from the repository
func (s *EmployeeDefault) FindById(id int) (employee employeesPkg.Employee, err error) {
	employee, err = s.rp.FindById(id)
	return
}

// CreateEmployee adds a new employee to the repository
func (s *EmployeeDefault) CreateEmployee(newEmployee employeesPkg.EmployeeAttributes) (employee employeesPkg.Employee, err error) {
	// validate required fields
	err = validateFields(newEmployee)
	if err != nil {
		return
	}

	// check for duplicates
	employees, _ := s.rp.FindAll()
	err = validateDuplicates(employees, newEmployee)
	if err != nil {
		return
	}

	// attempt to create the new employee
	return s.rp.CreateEmployee(newEmployee)
}

// UpdateEmployee updates an employee in the repository
func (s *EmployeeDefault) UpdateEmployee(inputEmployee employeesPkg.Employee) (employee employeesPkg.Employee, err error) {
	// find the existing employee
	internalEmployee, err := s.rp.FindById(inputEmployee.ID)
	if err != nil {
		return employeesPkg.Employee{}, utils.ErrNotFound
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
func validateFields(newEmployee employeesPkg.EmployeeAttributes) (err error) {
	if newEmployee.FirstName == "" || newEmployee.LastName == "" || newEmployee.CardNumberId == 0 || newEmployee.WarehouseId == 0 {
		return utils.ErrEmptyArguments
	}
	return
}

// validateDuplicates ensures that no existing employee has the same CardNumberId as the new employee
func validateDuplicates(employees map[int]employeesPkg.Employee, newEmployee employeesPkg.EmployeeAttributes) error {
	for _, employee := range employees {
		if employee.Attributes.CardNumberId == newEmployee.CardNumberId {
			return utils.ErrConflict
		}
	}
	return nil
}

// mergeEmployeeFields merges the fields of the input employee with the internal employee
func mergeEmployeeFields(inputEmployee, internalEmployee employeesPkg.Employee) employeesPkg.Employee {
	// copy existing fields to the new employee
	updatedEmployee := internalEmployee

	// update only the fields that are provided in the input
	if inputEmployee.Attributes.CardNumberId != 0 {
		updatedEmployee.Attributes.CardNumberId = inputEmployee.Attributes.CardNumberId
	}
	if inputEmployee.Attributes.FirstName != "" {
		updatedEmployee.Attributes.FirstName = inputEmployee.Attributes.FirstName
	}
	if inputEmployee.Attributes.LastName != "" {
		updatedEmployee.Attributes.LastName = inputEmployee.Attributes.LastName
	}
	if inputEmployee.Attributes.WarehouseId != 0 {
		updatedEmployee.Attributes.WarehouseId = inputEmployee.Attributes.WarehouseId
	}

	return updatedEmployee
}
