package repository

import (
	"fmt"

	employeesPkg "github.com/meli-fresh-products-api-backend-go-t2/internal/pkg"
)

type EmployeeMap struct {
	// db is a map of Employees
	db map[int]employeesPkg.Employee
}

func NewEmployeeRepository(db map[int]employeesPkg.Employee) *EmployeeMap {
	// default db
	defaultDb := make(map[int]employeesPkg.Employee)
	if db != nil {
		defaultDb = db
	}
	return &EmployeeMap{db: defaultDb}
}

func (r *EmployeeMap) FindAll() (v map[int]employeesPkg.Employee, err error) {
	v = make(map[int]employeesPkg.Employee)

	// copy db
	for key, value := range r.db {
		v[key] = value
		fmt.Printf("key: %v, value: %v\n", key, value)
	}

	return
}
