package loader

import (
	"encoding/json"
	"os"

	pkg "github.com/meli-fresh-products-api-backend-go-t2/internal/pkg"
)

// EmployeeJsonFile represents a loader for JSON files containing employee data
type EmployeeJsonFile struct {
	// path is the path to the file that contains the employees in json format
	path string
}

// NewEmployeeJsonFile creates a new instance of EmployeeJsonFile with the given file path
// it takes the path to the JSON file as a parameter and returns an instance of EmployeeJsonFile
func NewEmployeeJsonFile(path string) *EmployeeJsonFile {
	return &EmployeeJsonFile{
		path: path,
	}
}

// Load reads and loads employees data from the json file
// and returns a map with employee IDs as keys and an error in case of failure
func (l *EmployeeJsonFile) Load() (employee map[int]pkg.Employee, err error) {
	// open file
	file, err := os.Open(l.path)
	if err != nil {
		return
	}
	defer file.Close()

	// decode file
	var EmployeesJson []pkg.EmployeeJson
	err = json.NewDecoder(file).Decode(&EmployeesJson)
	if err != nil {
		return
	}

	// convert the decoded data to the internal Employee format
	employee = make(map[int]pkg.Employee)
	for _, employees := range EmployeesJson {
		employee[employees.ID] = pkg.Employee{
			ID: employees.ID,
			Attributes: pkg.EmployeeAttributes{
				CardNumberId: employees.CardNumberId,
				FirstName:    employees.FirstName,
				LastName:     employees.LastName,
				WarehouseId:  employees.WarehouseId,
			},
		}
	}

	return
}
