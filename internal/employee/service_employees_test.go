package employee_test

import (
	"errors"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/employee"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

type MockEmployeeRepository struct {
	mock.Mock
}

func (m *MockEmployeeRepository) FindAll() (employees map[int]internal.Employee, err error) {
	args := m.Called()
	return args.Get(0).(map[int]internal.Employee), args.Error(1)
}

func (m *MockEmployeeRepository) FindByID(id int) (employee internal.Employee, err error) {
	args := m.Called(id)
	return args.Get(0).(internal.Employee), args.Error(1)
}

func (m *MockEmployeeRepository) CreateEmployee(newEmployee internal.EmployeeAttributes) (employee internal.Employee, err error) {
	args := m.Called(newEmployee)
	return args.Get(0).(internal.Employee), args.Error(1)
}

func (m *MockEmployeeRepository) UpdateEmployee(inputEmployee internal.Employee) (employee internal.Employee, err error) {
	args := m.Called(inputEmployee)
	return args.Get(0).(internal.Employee), args.Error(1)
}

func (m *MockEmployeeRepository) DeleteEmployee(id int) (err error) {
	args := m.Called(id)
	return args.Error(0)
}

type MockWarehouseRepository struct {
	mock.Mock
}

func (m *MockWarehouseRepository) GetByID(int) (internal.Warehouse, error) {
	args := m.Called()
	return args.Get(0).(internal.Warehouse), args.Error(1)
}

func TestUnitEmployees_GetAllEmployees(t *testing.T) {
	type testCase struct {
		name              string
		mockEmployees     map[int]internal.Employee
		mockError         error
		expectedEmployees map[int]internal.Employee
		expectedError     error
	}

	employeesOk := map[int]internal.Employee{
		1: {
			ID: 1,
			Attributes: internal.EmployeeAttributes{
				CardNumberID: "E001",
				FirstName:    "Alice",
				LastName:     "Johnson",
				WarehouseID:  1,
			},
		},
		2: {
			ID: 2,
			Attributes: internal.EmployeeAttributes{
				CardNumberID: "E002",
				FirstName:    "Bob",
				LastName:     "Anderson",
				WarehouseID:  2,
			},
		},
	}

	tests := []testCase{
		{
			name:              "OK",
			mockEmployees:     employeesOk,
			mockError:         nil,
			expectedEmployees: employeesOk,
			expectedError:     nil,
		},
		{
			name:              "ERROR",
			mockEmployees:     nil,
			mockError:         errors.New("some error"),
			expectedEmployees: nil,
			expectedError:     errors.New("some error"),
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			repositoryEmployee := new(MockEmployeeRepository)
			repositoryWarehouse := new(MockWarehouseRepository)
			service := employee.NewEmployeeService(repositoryEmployee, repositoryWarehouse)

			repositoryEmployee.On("FindAll").Return(tc.mockEmployees, tc.mockError)

			result, err := service.FindAll()

			require.Equal(t, tc.expectedEmployees, result)

			if tc.expectedError == nil {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedError.Error())
			}

			repositoryEmployee.AssertExpectations(t)
			repositoryWarehouse.AssertExpectations(t)
		})
	}
}

func TestUnitEmployees_FindByID(t *testing.T) {
	type testCase struct {
		name             string
		paramID          int
		mockEmployee     internal.Employee
		mockError        error
		expectedEmployee internal.Employee
		expectedError    error
	}

	exampleEmployee := internal.Employee{
		ID: 1,
		Attributes: internal.EmployeeAttributes{
			CardNumberID: "E001",
			FirstName:    "Alice",
			LastName:     "Johnson",
			WarehouseID:  1,
		},
	}

	tests := []testCase{
		{
			name:             "OK",
			paramID:          1,
			mockEmployee:     exampleEmployee,
			mockError:        nil,
			expectedEmployee: exampleEmployee,
			expectedError:    nil,
		},
		{
			name:             "NOT_FOUND",
			paramID:          2,
			mockEmployee:     internal.Employee{},
			mockError:        utils.ErrNotFound,
			expectedEmployee: internal.Employee{},
			expectedError:    utils.ErrNotFound,
		},
		{
			name:             "REPO_ERROR",
			paramID:          3,
			mockEmployee:     internal.Employee{},
			mockError:        errors.New("some repo error"),
			expectedEmployee: internal.Employee{},
			expectedError:    errors.New("some repo error"),
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			repositoryEmployee := new(MockEmployeeRepository)
			repositoryWarehouse := new(MockWarehouseRepository)
			service := employee.NewEmployeeService(repositoryEmployee, repositoryWarehouse)

			repositoryEmployee.
				On("FindByID", tc.paramID).
				Return(tc.mockEmployee, tc.mockError)

			result, err := service.FindByID(tc.paramID)

			require.Equal(t, tc.expectedEmployee, result)

			if tc.expectedError == nil {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedError.Error())
			}

			repositoryEmployee.AssertExpectations(t)
			repositoryWarehouse.AssertExpectations(t)
		})
	}
}

func TestUnitEmployees_CreateEmployee(t *testing.T) {
	type testCase struct {
		name             string
		inputAttributes  internal.EmployeeAttributes
		mockFindAllRes   map[int]internal.Employee
		mockFindAllErr   error
		mockCreateRes    internal.Employee
		mockCreateErr    error
		mockWarehouse    internal.Warehouse
		mockWarehouseErr error
		expectedOutput   internal.Employee
		expectedError    error
	}

	validAttr := internal.EmployeeAttributes{
		CardNumberID: "E003",
		FirstName:    "Charlie",
		LastName:     "Brown",
		WarehouseID:  1,
	}
	validEmployee := internal.Employee{
		ID:         3,
		Attributes: validAttr,
	}

	existingEmployees := map[int]internal.Employee{
		1: {
			ID: 1,
			Attributes: internal.EmployeeAttributes{
				CardNumberID: "E001",
				FirstName:    "Alice",
				LastName:     "Johnson",
				WarehouseID:  1,
			},
		},
		2: {
			ID: 2,
			Attributes: internal.EmployeeAttributes{
				CardNumberID: "E002",
				FirstName:    "Bob",
				LastName:     "Anderson",
				WarehouseID:  2,
			},
		},
	}

	tests := []testCase{
		{
			name:             "OK",
			inputAttributes:  validAttr,
			mockFindAllRes:   existingEmployees,
			mockFindAllErr:   nil,
			mockWarehouse:    internal.Warehouse{ID: 1},
			mockWarehouseErr: nil,
			mockCreateRes:    validEmployee,
			mockCreateErr:    nil,
			expectedOutput:   validEmployee,
			expectedError:    nil,
		},
		{
			name: "ERROR - empty required fields",
			inputAttributes: internal.EmployeeAttributes{
				CardNumberID: "",
				FirstName:    "NoCardNumber",
				LastName:     "NotAllowed",
				WarehouseID:  1,
			},
			mockFindAllRes:   nil,
			mockFindAllErr:   nil,
			mockWarehouse:    internal.Warehouse{},
			mockWarehouseErr: nil,
			mockCreateRes:    internal.Employee{},
			mockCreateErr:    nil,
			expectedOutput:   internal.Employee{},
			expectedError:    utils.ErrEmptyArguments,
		},
		{
			name:             "ERROR - repo FindAll fails",
			inputAttributes:  validAttr,
			mockFindAllRes:   nil,
			mockFindAllErr:   errors.New("repo FindAll error"),
			mockWarehouse:    internal.Warehouse{},
			mockWarehouseErr: nil,
			mockCreateRes:    internal.Employee{},
			mockCreateErr:    nil,
			expectedOutput:   internal.Employee{},
			expectedError:    errors.New("repo FindAll error"),
		},
		{
			name: "ERROR - conflict by CardNumberID",
			inputAttributes: internal.EmployeeAttributes{
				CardNumberID: "E001",
				FirstName:    "Dup",
				LastName:     "Conflict",
				WarehouseID:  1,
			},
			mockFindAllRes:   existingEmployees,
			mockFindAllErr:   nil,
			mockWarehouse:    internal.Warehouse{ID: 1},
			mockWarehouseErr: nil,
			mockCreateRes:    internal.Employee{},
			mockCreateErr:    nil,
			expectedOutput:   internal.Employee{},
			expectedError:    utils.ErrConflict,
		},
		{
			name:             "ERROR - warehouse not found",
			inputAttributes:  validAttr,
			mockFindAllRes:   existingEmployees,
			mockFindAllErr:   nil,
			mockWarehouse:    internal.Warehouse{},
			mockWarehouseErr: nil,
			mockCreateRes:    internal.Employee{},
			mockCreateErr:    nil,
			expectedOutput:   internal.Employee{},
			expectedError:    utils.EDependencyNotFound("warehouse", "id: 1"),
		},
		{
			name:             "ERROR - createEmployee fails on repo",
			inputAttributes:  validAttr,
			mockFindAllRes:   existingEmployees,
			mockFindAllErr:   nil,
			mockWarehouse:    internal.Warehouse{ID: 1},
			mockWarehouseErr: nil,
			mockCreateRes:    internal.Employee{},
			mockCreateErr:    errors.New("repo CreateEmployee error"),
			expectedOutput:   internal.Employee{},
			expectedError:    errors.New("repo CreateEmployee error"),
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			repositoryEmployee := new(MockEmployeeRepository)
			repositoryWarehouse := new(MockWarehouseRepository)
			service := employee.NewEmployeeService(repositoryEmployee, repositoryWarehouse)

			emptyFields := (tc.inputAttributes.CardNumberID == "" ||
				tc.inputAttributes.FirstName == "" ||
				tc.inputAttributes.LastName == "")

			if !emptyFields {
				repositoryEmployee.
					On("FindAll").
					Return(tc.mockFindAllRes, tc.mockFindAllErr).
					Maybe()

				if tc.mockFindAllErr == nil && tc.expectedError != utils.ErrConflict {
					repositoryWarehouse.
						On("GetByID").
						Return(tc.mockWarehouse, tc.mockWarehouseErr).
						Maybe()

					if tc.mockWarehouse.ID != 0 && tc.mockCreateErr != nil {
						repositoryEmployee.
							On("CreateEmployee", tc.inputAttributes).
							Return(tc.mockCreateRes, tc.mockCreateErr).
							Maybe()
					} else if tc.mockWarehouse.ID != 0 && tc.mockCreateErr == nil && tc.expectedError == nil {
						repositoryEmployee.
							On("CreateEmployee", tc.inputAttributes).
							Return(tc.mockCreateRes, tc.mockCreateErr).
							Maybe()
					}
				}
			}

			result, err := service.CreateEmployee(tc.inputAttributes)

			require.Equal(t, tc.expectedOutput, result)

			if tc.expectedError == nil {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedError.Error())
			}

			repositoryEmployee.AssertExpectations(t)
			repositoryWarehouse.AssertExpectations(t)
		})
	}
}

func TestUnitEmployees_UpdateEmployee(t *testing.T) {
	type testCase struct {
		name             string
		inputEmployee    internal.Employee
		mockFindByIDRes  internal.Employee
		mockFindByIDErr  error
		mockWarehouse    internal.Warehouse
		mockWarehouseErr error
		mockUpdateRes    internal.Employee
		mockUpdateErr    error
		expectedEmployee internal.Employee
		expectedError    error
	}

	existingEmployee := internal.Employee{
		ID: 1,
		Attributes: internal.EmployeeAttributes{
			CardNumberID: "E001",
			FirstName:    "Alice",
			LastName:     "Johnson",
			WarehouseID:  1,
		},
	}

	updatedEmployee := internal.Employee{
		ID: 1,
		Attributes: internal.EmployeeAttributes{
			CardNumberID: "E005",
			FirstName:    "AliceUpdated",
			LastName:     "Johnson",
			WarehouseID:  2,
		},
	}

	tests := []testCase{
		{
			name: "OK - atualiza parcialmente (CardNumberID, WarehouseID, FirstName)",
			inputEmployee: internal.Employee{
				ID: 1,
				Attributes: internal.EmployeeAttributes{
					CardNumberID: "E005",
					FirstName:    "AliceUpdated",
					WarehouseID:  2,
				},
			},
			mockFindByIDRes:  existingEmployee,
			mockFindByIDErr:  nil,
			mockWarehouse:    internal.Warehouse{ID: 2},
			mockWarehouseErr: nil,
			mockUpdateRes:    updatedEmployee,
			mockUpdateErr:    nil,
			expectedEmployee: updatedEmployee,
			expectedError:    nil,
		},
		{
			name: "WAREHOUSE INTERNAL ERR",
			inputEmployee: internal.Employee{
				ID:         1,
				Attributes: internal.EmployeeAttributes{WarehouseID: 2},
			},
			mockFindByIDRes:  existingEmployee,
			mockFindByIDErr:  nil,
			mockWarehouse:    internal.Warehouse{},
			mockWarehouseErr: errors.New("db connection fail"),
			mockUpdateRes:    internal.Employee{},
			mockUpdateErr:    nil,
			expectedEmployee: internal.Employee{},
			expectedError:    errors.New("db connection fail"),
		},
		{
			name: "WAREHOUSE NOT FOUND",
			inputEmployee: internal.Employee{
				ID:         1,
				Attributes: internal.EmployeeAttributes{WarehouseID: 55},
			},
			mockFindByIDRes:  existingEmployee,
			mockFindByIDErr:  nil,
			mockWarehouse:    internal.Warehouse{},
			mockWarehouseErr: nil,
			mockUpdateRes:    internal.Employee{},
			mockUpdateErr:    nil,
			expectedEmployee: internal.Employee{},
			expectedError:    utils.EDependencyNotFound("warehouse", "id: 55"),
		},
		{
			name: "ERROR - UpdateEmployee falha no repositório",
			inputEmployee: internal.Employee{
				ID: 1,
				Attributes: internal.EmployeeAttributes{
					CardNumberID: "E010",
					FirstName:    "Alice2",
					WarehouseID:  9,
				},
			},
			mockFindByIDRes:  existingEmployee,
			mockFindByIDErr:  nil,
			mockWarehouse:    internal.Warehouse{ID: 9},
			mockWarehouseErr: nil,
			mockUpdateRes:    internal.Employee{},
			mockUpdateErr:    errors.New("repo update failed"),
			expectedEmployee: internal.Employee{},
			expectedError:    errors.New("repo update failed"),
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			repositoryEmployee := new(MockEmployeeRepository)
			repositoryWarehouse := new(MockWarehouseRepository)
			service := employee.NewEmployeeService(repositoryEmployee, repositoryWarehouse)

			repositoryEmployee.
				On("FindByID", tc.inputEmployee.ID).
				Return(tc.mockFindByIDRes, tc.mockFindByIDErr).
				Maybe()

			if tc.mockFindByIDErr == nil {
				repositoryWarehouse.
					On("GetByID").
					Return(tc.mockWarehouse, tc.mockWarehouseErr).
					Maybe()

				if tc.mockWarehouse.ID != 0 && tc.mockWarehouseErr == nil {
					repositoryEmployee.
						On("UpdateEmployee", mock.Anything).
						Return(tc.mockUpdateRes, tc.mockUpdateErr).
						Maybe()
				}
			}

			actual, err := service.UpdateEmployee(tc.inputEmployee)

			require.Equal(t, tc.expectedEmployee, actual)

			if tc.expectedError == nil {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedError.Error())
			}

			repositoryEmployee.AssertExpectations(t)
			repositoryWarehouse.AssertExpectations(t)
		})
	}
}

func TestUnitEmployees_DeleteEmployee(t *testing.T) {
	type testCase struct {
		name                 string
		paramID              int
		mockFindByIDEmployee internal.Employee
		mockFindByIDError    error
		mockDeleteError      error
		expectedError        error
	}

	existingEmployee := internal.Employee{
		ID: 1,
		Attributes: internal.EmployeeAttributes{
			CardNumberID: "E001",
			FirstName:    "Alice",
			LastName:     "Johnson",
			WarehouseID:  1,
		},
	}

	tests := []testCase{
		{
			name:                 "OK",
			paramID:              1,
			mockFindByIDEmployee: existingEmployee,
			mockFindByIDError:    nil,
			mockDeleteError:      nil,
			expectedError:        nil,
		},
		{
			name:                 "NOT_FOUND",
			paramID:              2,
			mockFindByIDEmployee: internal.Employee{},
			mockFindByIDError:    errors.New("some repo error"),
			mockDeleteError:      nil,
			expectedError:        utils.ErrNotFound,
		},
		{
			name:                 "ERROR_DELETE - invalid arguments",
			paramID:              3,
			mockFindByIDEmployee: existingEmployee,
			mockFindByIDError:    nil,
			mockDeleteError:      errors.New("delete error"),
			expectedError:        utils.ErrInvalidArguments,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			repositoryEmployee := new(MockEmployeeRepository)
			repositoryWarehouse := new(MockWarehouseRepository)
			service := employee.NewEmployeeService(repositoryEmployee, repositoryWarehouse)

			repositoryEmployee.
				On("FindByID", tc.paramID).
				Return(tc.mockFindByIDEmployee, tc.mockFindByIDError).
				Maybe()

			if tc.mockFindByIDError == nil {
				repositoryEmployee.
					On("DeleteEmployee", tc.mockFindByIDEmployee.ID).
					Return(tc.mockDeleteError).
					Maybe()
			}

			err := service.DeleteEmployee(tc.paramID)

			if tc.expectedError == nil {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedError.Error())
			}

			repositoryEmployee.AssertExpectations(t)
			repositoryWarehouse.AssertExpectations(t)
		})
	}
}
