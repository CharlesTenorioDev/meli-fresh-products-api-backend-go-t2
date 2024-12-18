package repository

import (
	employeesPkg "github.com/meli-fresh-products-api-backend-go-t2/internal/pkg"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

// EmployeeRepository represents the repository implementation for employee data storage using a map
type EmployeeRepository struct {
	// db is an in-memory map that stores employees by their ID
	db map[int]employeesPkg.Employee
}

// NewEmployeeRepository creates a new EmployeeRepository repository instance
// if a non-nil `db` is provided, it uses it as the initial database
// otherwise, it initializes an empty map
func NewEmployeeRepository(db map[int]employeesPkg.Employee) *EmployeeRepository {
	// initialize a default database map if no `db` is provided
	defaultDb := make(map[int]employeesPkg.Employee)
	if db != nil {
		defaultDb = db
	}
	return &EmployeeRepository{db: defaultDb}
}

// FindAll retrieves all employees from the repository
// it copies the data from the internal map to avoid direct manipulation of the repository state
func (r *EmployeeRepository) FindAll() (employees map[int]employeesPkg.Employee, err error) {
	// create a new map to hold the employees
	employees = make(map[int]employeesPkg.Employee)

	// copy each employee from the repository's internal map
	for key, value := range r.db {
		employees[key] = value
	}

	return
}

// FindById retrieves an employee by their ID
// if the employee is not found, it returns an error
func (r *EmployeeRepository) FindById(id int) (employee employeesPkg.Employee, err error) {
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
func (r *EmployeeRepository) CreateEmployee(newEmployee employeesPkg.EmployeeAttributes) (employee employeesPkg.Employee, err error) {
	// generate a new ID for the employee by finding the largest existing ID and adding 1
	newID := utils.GetBiggestId(r.db) + 1

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

// UpdateEmployee updates an employee's data in the repository
func (r *EmployeeRepository) UpdateEmployee(inputEmployee employeesPkg.Employee) (employee employeesPkg.Employee, err error) {
	// check if the employee exists in the internal map
	existingEmployee, exists := r.db[inputEmployee.ID]
	if !exists {
		return employeesPkg.Employee{}, utils.ErrNotFound
	}

	// merge fields from the inputEmployee
	if inputEmployee.Attributes.CardNumberId != 0 {
		existingEmployee.Attributes.CardNumberId = inputEmployee.Attributes.CardNumberId
	}
	if inputEmployee.Attributes.FirstName != "" {
		existingEmployee.Attributes.FirstName = inputEmployee.Attributes.FirstName
	}
	if inputEmployee.Attributes.LastName != "" {
		existingEmployee.Attributes.LastName = inputEmployee.Attributes.LastName
	}
	if inputEmployee.Attributes.WarehouseId != 0 {
		existingEmployee.Attributes.WarehouseId = inputEmployee.Attributes.WarehouseId
	}

	// update the map with the modified employee
	r.db[inputEmployee.ID] = existingEmployee

	return existingEmployee, nil
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
