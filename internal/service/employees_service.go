package service

import (
	employeesPkg "github.com/meli-fresh-products-api-backend-go-t2/internal/pkg"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type EmployeeDefault struct {
	rp employeesPkg.EmployeeRepository
}

func NewEmployeeService(rp employeesPkg.EmployeeRepository) *EmployeeDefault {
	return &EmployeeDefault{rp: rp}
}

// FindAll is a method that returns a map of all employees
func (s *EmployeeDefault) FindAll() (employees map[int]employeesPkg.Employee, err error) {
	employees, err = s.rp.FindAll()
	return
}

// FindAll is a method that returns a map of all employees
func (s *EmployeeDefault) FindById(id int) (employees map[int]employeesPkg.Employee, err error) {
	employees, err = s.rp.FindById(id)
	return
}

func (s *EmployeeDefault) CreateEmployee(newEmployee employeesPkg.EmployeeAttributes) (employee employeesPkg.Employee, err error) {
	err = validateFields(newEmployee)
	if err != nil {
		return
	}

	employees, _ := s.rp.FindAll()
	err = validateDuplicates(employees, newEmployee)
	if err != nil {
		return
	}
	return s.rp.CreateEmployee(newEmployee)
}

func validateFields(newEmployee employeesPkg.EmployeeAttributes) (err error) {
	if newEmployee.FirstName == "" || newEmployee.LastName == "" || newEmployee.CardNumberId == 0 || newEmployee.WarehouseId == 0 {
		return utils.ErrEmptyArguments
	}
	return nil
}

func validateDuplicates(employees map[int]employeesPkg.Employee, newEmployee employeesPkg.EmployeeAttributes) (err error) {
	for _, employee := range employees {
		if employee.Attributes.CardNumberId == newEmployee.CardNumberId {
			return utils.ErrConflict
		}
	}
	return nil
}
