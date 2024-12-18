package repository

import (
	"fmt"

	employeesPkg "github.com/meli-fresh-products-api-backend-go-t2/internal/pkg"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
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

func (r *EmployeeMap) FindAll() (employees map[int]employeesPkg.Employee, err error) {
	employees = make(map[int]employeesPkg.Employee)

	// copy db
	for key, value := range r.db {
		employees[key] = value
		fmt.Printf("key: %v, value: %v\n", key, value)
	}

	return
}

// FindById is a method to find a employee by its ID
func (r *EmployeeMap) FindById(id int) (employees map[int]employeesPkg.Employee, err error) {
	employees = make(map[int]employeesPkg.Employee)

	e, exists := r.db[id]
	if !exists {
		err = utils.ErrNotFound
		return
	}

	employees[id] = e
	return
}

func (r *EmployeeMap) CreateEmployee(newEmployee employeesPkg.EmployeeAttributes) (employee employeesPkg.Employee, err error) {
	newID := utils.GetBiggestId(map[int]employeesPkg.Employee(r.db)) + 1

	employee = employeesPkg.Employee{
		ID: newID,
		Attributes: employeesPkg.EmployeeAttributes{
			CardNumberId: newEmployee.CardNumberId,
			FirstName:    newEmployee.FirstName,
			LastName:     newEmployee.LastName,
			WarehouseId:  newEmployee.WarehouseId,
		},
	}

	r.db[employee.ID] = employee
	return
}
