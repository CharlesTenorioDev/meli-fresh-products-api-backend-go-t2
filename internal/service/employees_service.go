package service

import employeesPkg "github.com/meli-fresh-products-api-backend-go-t2/internal/pkg"

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
