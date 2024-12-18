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

// FindAll is a method that returns a map of all employees and an error if something goes wrong
func (s *EmployeeDefault) FindAll() (employees map[int]employeesPkg.Employee, err error) {
	employees, err = s.rp.FindAll()
	return
}

// FindById is a method that returns a map containing the employee and an error if employee is not found
func (s *EmployeeDefault) FindById(id int) (employees map[int]employeesPkg.Employee, err error) {
	employees, err = s.rp.FindById(id)
	return
}

// CreateEmployee adds a new employee to the repository after validating the input and checking for duplicates
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
	return s.rp.CreateEmployee(newEmployee)
}

// validateFields checks if the required fields of a new employee are not empty
func validateFields(newEmployee employeesPkg.EmployeeAttributes) (err error) {
	if newEmployee.FirstName == "" || newEmployee.LastName == "" || newEmployee.CardNumberId == 0 || newEmployee.WarehouseId == 0 {
		return utils.ErrEmptyArguments
	}
	return nil
}

// validateDuplicates ensures that no existing employee has the same CardNumberId as the new employee
func validateDuplicates(employees map[int]employeesPkg.Employee, newEmployee employeesPkg.EmployeeAttributes) (err error) {
	for _, employee := range employees {
		if employee.Attributes.CardNumberId == newEmployee.CardNumberId {
			return utils.ErrConflict
		}
	}
	return nil
}
