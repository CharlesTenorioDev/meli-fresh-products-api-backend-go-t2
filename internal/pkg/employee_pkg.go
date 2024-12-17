package employeesPkg

type Employee struct {
	ID           int
	CardNumberId int
	FirstName    string
	LastName     string
	WarehouseId  int
}

type EmployeeJson struct {
	ID           int    `json:"id"`
	CardNumberId int    `json:"card_number_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	WarehouseId  int    `json:"warehouse_id"`
}

type EmployeeRepository interface {
	FindAll() (employees map[int]Employee, err error)
}

type EmployeeService interface {
	FindAll() (employees map[int]Employee, err error)
}
