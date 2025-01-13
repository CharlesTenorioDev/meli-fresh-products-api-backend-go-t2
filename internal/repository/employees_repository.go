package repository

import (
	"database/sql"
	"log"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type EmployeeRepository struct {
	db *sql.DB
}

func NewEmployeeRepository(db *sql.DB) *EmployeeRepository {
	return &EmployeeRepository{db: db}
}

// FindAll retrieves all employees from the repository
func (r *EmployeeRepository) FindAll() (map[int]internal.Employee, error) {
	rows, err := r.db.Query("SELECT id, id_card_number, first_name, last_name, warehouse_id FROM employees")
	if err != nil {
		log.Printf("Error in FindAll Query: %v", err)
		return nil, err
	}
	defer rows.Close()

	employees := make(map[int]internal.Employee)
	for rows.Next() {
		var emp internal.Employee
		emp.Attributes = internal.EmployeeAttributes{}
		err := rows.Scan(&emp.ID, &emp.Attributes.CardNumberId, &emp.Attributes.FirstName, &emp.Attributes.LastName, &emp.Attributes.WarehouseId)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}
		employees[emp.ID] = emp
	}

	return employees, nil
}

// FindById retrieves an employee by their ID
func (r *EmployeeRepository) FindById(id int) (internal.Employee, error) {
	var employee internal.Employee
	employee.Attributes = internal.EmployeeAttributes{}

	err := r.db.QueryRow("SELECT id, id_card_number, first_name, last_name, warehouse_id FROM employees WHERE id = ?", id).
		Scan(&employee.ID, &employee.Attributes.CardNumberId, &employee.Attributes.FirstName, &employee.Attributes.LastName, &employee.Attributes.WarehouseId)
	if err == sql.ErrNoRows {
		return internal.Employee{}, utils.ErrNotFound
	}
	if err != nil {
		log.Printf("Error in FindById Query: %v", err)
		return internal.Employee{}, err
	}

	return employee, nil
}

// CreateEmployee adds a new employee
func (r *EmployeeRepository) CreateEmployee(newEmployee internal.EmployeeAttributes) (internal.Employee, error) {
	result, err := r.db.Exec("INSERT INTO employees (id_card_number, first_name, last_name, warehouse_id) VALUES (?, ?, ?, ?)",
		newEmployee.CardNumberId, newEmployee.FirstName, newEmployee.LastName, newEmployee.WarehouseId)
	if err != nil {
		log.Printf("Error in CreateEmployee Query: %v", err)
		return internal.Employee{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error getting LastInsertId: %v", err)
		return internal.Employee{}, err
	}

	return r.FindById(int(id))
}

// UpdateEmployee updates an employee's data
func (r *EmployeeRepository) UpdateEmployee(inputEmployee internal.Employee) (internal.Employee, error) {
	_, err := r.db.Exec("UPDATE employees SET id_card_number = ?, first_name = ?, last_name = ?, warehouse_id = ? WHERE id = ?",
		inputEmployee.Attributes.CardNumberId, inputEmployee.Attributes.FirstName, inputEmployee.Attributes.LastName, inputEmployee.Attributes.WarehouseId, inputEmployee.ID)
	if err != nil {
		log.Printf("Error in UpdateEmployee Query: %v", err)
		return internal.Employee{}, err
	}

	return r.FindById(inputEmployee.ID)
}

// DeleteEmployee removes an employee
func (r *EmployeeRepository) DeleteEmployee(id int) error {
	_, err := r.db.Exec("DELETE FROM employees WHERE id = ?", id)
	if err != nil {
		log.Printf("Error in DeleteEmployee Query: %v", err)
		return err
	}

	return nil
}
