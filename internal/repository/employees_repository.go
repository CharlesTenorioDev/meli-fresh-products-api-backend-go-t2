package repository

import (
	pkg "github.com/meli-fresh-products-api-backend-go-t2/internal/pkg"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

// EmployeeRepository represents the repository implementation for employee data storage using a map
type EmployeeRepository struct {
	// db is an in-memory map that stores employees by their ID
	db map[int]pkg.Employee
}

// NewEmployeeRepository creates a new EmployeeRepository repository instance
// if a non-nil `db` is provided, it uses it as the initial database
// otherwise, it initializes an empty map
func NewEmployeeRepository(db map[int]pkg.Employee) *EmployeeRepository {
	// initialize a default database map if no `db` is provided
	defaultDb := make(map[int]pkg.Employee)
	if db != nil {
		defaultDb = db
	}
	return &EmployeeRepository{db: defaultDb}
}

// FindAll retrieves all employees from the repository
// it copies the data from the internal map to avoid direct manipulation of the repository state
func (r *EmployeeRepository) FindAll() (employees map[int]pkg.Employee, err error) {
	// create a new map to hold the employees
	employees = make(map[int]pkg.Employee)

	// copy each employee from the repository's internal map
	for key, value := range r.db {
		employees[key] = value
	}

	return
}

// FindById retrieves an employee by their ID
// if the employee is not found, it returns an error
func (r *EmployeeRepository) FindById(id int) (employee pkg.Employee, err error) {
	// check if the employee exists in the internal map
	employee, exists := r.db[id]
	if !exists {
		// return an error if the employee is not found
		err = utils.ErrNotFound
		return
	}

	return employee, nil
}

// CreateEmployee adds a new employee to the repository
// it generates a new unique ID for the employee and adds them to the map
func (r *EmployeeRepository) CreateEmployee(newEmployee pkg.EmployeeAttributes) (employee pkg.Employee, err error) {
	// generate a new ID for the employee by finding the largest existing ID and adding 1
	newID := utils.GetBiggestId(r.db) + 1

	// create the new employee struct with the provided attributes and the new ID
	employee = pkg.Employee{
		ID: newID,
		Attributes: pkg.EmployeeAttributes{
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

// UpdateEmployee updates an employee's data in the repository
func (r *EmployeeRepository) UpdateEmployee(inputEmployee pkg.Employee) (employee pkg.Employee, err error) {
	// update the modified employee

	r.db[inputEmployee.ID] = inputEmployee
	return inputEmployee, nil
}

// DeleteEmployee removes an employee from the in-memory repository
func (r *EmployeeRepository) DeleteEmployee(id int) error {
	// check if the employee exists
	if _, exists := r.db[id]; !exists {
		return utils.ErrNotFound
	}

	// delete the employee from the map using the ID
	delete(r.db, id)
	return nil
}
