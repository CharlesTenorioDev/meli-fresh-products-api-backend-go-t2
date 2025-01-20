package internal

// Employee represents an employee entity with its unique ID and attributes
type Employee struct {
	ID         int                `json:"id"`
	Attributes EmployeeAttributes `json:"attributes"`
}

// EmployeeAttributes defines the details associated with an employee
type EmployeeAttributes struct {
	CardNumberID string `json:"card_number_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	WarehouseID  int    `json:"warehouse_id"`
}

// EmployeeJSON defines the structure of the employee data as it appears in a json file
type EmployeeJSON struct {
	ID           int    `json:"id"`
	CardNumberId string `json:"card_number_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	WarehouseID  int    `json:"warehouse_id"`
}

// EmployeeRepository defines the interface for employee data persistence
// it specifies methods for fetching and creating employee data
type EmployeeRepository interface {
	FindAll() (employees map[int]Employee, err error)
	FindByID(id int) (employee Employee, err error)
	CreateEmployee(newEmployee EmployeeAttributes) (employee Employee, err error)
	UpdateEmployee(inputEmployee Employee) (employee Employee, err error)
	DeleteEmployee(id int) (err error)
}

// EmployeeService defines the interface for employee-related business logic
// it includes methods for fetching and creating employees
type EmployeeService interface {
	FindAll() (employees map[int]Employee, err error)
	FindByID(id int) (employee Employee, err error)
	CreateEmployee(newEmployee EmployeeAttributes) (employee Employee, err error)
	UpdateEmployee(inputEmployee Employee) (employee Employee, err error)
	DeleteEmployee(id int) (err error)
}

type EmployeesWarehouseValidation interface {
	GetByID(int) (Warehouse, error)
}
