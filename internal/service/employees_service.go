package service

import (
	employeesPkg "github.com/yywatanabe_meli/api-produtos-frescos/internal/pkg"
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
