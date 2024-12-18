package repository

import (
	"fmt"

	employeesPkg "github.com/meli-fresh-products-api-backend-go-t2/internal/pkg"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

// EmployeeMap represents the repository implementation for employee data storage using a map
type EmployeeMap struct {
	// db is an in-memory map that stores employees by their ID
	db map[int]employeesPkg.Employee
}

// NewEmployeeRepository creates a new EmployeeMap repository instance
// if a non-nil `db` is provided, it uses it as the initial database
// otherwise, it initializes an empty map
func NewEmployeeRepository(db map[int]employeesPkg.Employee) *EmployeeMap {
	// initialize a default database map if no `db` is provided
	defaultDb := make(map[int]employeesPkg.Employee)
	if db != nil {
		defaultDb = db
	}
	return &EmployeeMap{db: defaultDb}
}

// FindAll retrieves all employees from the repository
// it copies the data from the internal map to avoid direct manipulation of the repository state
func (r *EmployeeMap) FindAll() (employees map[int]employeesPkg.Employee, err error) {
	// create a new map to hold the employees
	employees = make(map[int]employeesPkg.Employee)

	// copy each employee from the repository's internal map
	for key, value := range r.db {
		employees[key] = value
		fmt.Printf("key: %v, value: %v\n", key, value)
	}

	return
}

// FindById retrieves an employee by their ID
// if the employee is not found, it returns an error
func (r *EmployeeMap) FindById(id int) (employees map[int]employeesPkg.Employee, err error) {
	// create a map to hold the result
	employees = make(map[int]employeesPkg.Employee)

	// check if the employee exists in the internal map
	e, exists := r.db[id]
	if !exists {
		// return an error if the employee is not found
		err = utils.ErrNotFound
		return
	}

	// add the employee to the result map
	employees[id] = e
	return
}

// CreateEmployee adds a new employee to the repository
// it generates a new unique ID for the employee and adds them to the map
func (r *EmployeeMap) CreateEmployee(newEmployee employeesPkg.EmployeeAttributes) (employee employeesPkg.Employee, err error) {
	// generate a new ID for the employee by finding the largest existing ID and adding 1
	newID := utils.GetBiggestId(map[int]employeesPkg.Employee(r.db)) + 1

	// create the new employee struct with the provided attributes and the new ID
	employee = employeesPkg.Employee{
		ID: newID,
		Attributes: employeesPkg.EmployeeAttributes{
			CardNumberId: newEmployee.CardNumberId,
			FirstName:    newEmployee.FirstName,
			LastName:     newEmployee.LastName,
			WarehouseId:  newEmployee.WarehouseId,
		},
	}

	// add the new employee to the repository's internal map
	r.db[employee.ID] = employee
	return
}
