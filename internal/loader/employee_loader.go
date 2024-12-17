package loader

import (
	"encoding/json"
	"os"

	employeesPkg "github.com/yywatanabe_meli/api-produtos-frescos/internal/pkg"
)

type EmployeeJsonFile struct {
	// path is the path to the file that contains the employees in json format
	path string
}

type EmployeeJson struct {
	ID           int    `json:"id"`
	CardNumberId int    `json:"card_number_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	WarehouseId  int    `json:"warehouse_id"`
}

func NewEmployeeJsonFile(path string) *EmployeeJsonFile {
	return &EmployeeJsonFile{
		path: path,
	}
}

func (l *EmployeeJsonFile) Load() (v map[int]employeesPkg.Employee, err error) {
	// open file
	file, err := os.Open(l.path)
	if err != nil {
		return
	}
	defer file.Close()

	// decode file
	var EmployeesJson []EmployeeJson
	err = json.NewDecoder(file).Decode(&EmployeesJson)
	if err != nil {
		return
	}

	// serialize Employees
	v = make(map[int]employeesPkg.Employee)
	for _, employees := range EmployeesJson {
		v[employees.ID] = employeesPkg.Employee{
			ID:           employees.ID,
			CardNumberId: employees.CardNumberId,
			FirstName:    employees.FirstName,
			LastName:     employees.LastName,
			WarehouseId:  employees.WarehouseId,
		}
	}

	return
}
