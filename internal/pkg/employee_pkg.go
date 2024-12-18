package employeesPkg

type Employee struct {
	ID         int                `json:"id"`
	Attributes EmployeeAttributes `json:"attributes"`
}
type EmployeeAttributes struct {
	CardNumberId int    `json:"card_number_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	WarehouseId  int    `json:"warehouse_id"`
}

type EmployeeRepository interface {
	FindAll() (employees map[int]Employee, err error)
	FindById(id int) (employees map[int]Employee, err error)
	CreateEmployee(newEmployee EmployeeAttributes) (employee Employee, err error)
}

type EmployeeService interface {
	FindAll() (employees map[int]Employee, err error)
	FindById(id int) (employees map[int]Employee, err error)
	CreateEmployee(newEmployee EmployeeAttributes) (employee Employee, err error)
}
