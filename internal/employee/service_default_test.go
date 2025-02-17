package employee

import (
	"testing"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockEmployeeRepository struct {
	mock.Mock
}

func (m *mockEmployeeRepository) FindAll() (map[int]internal.Employee, error) {
	args := m.Called()
	return args.Get(0).(map[int]internal.Employee), args.Error(1)
}

func (m *mockEmployeeRepository) FindByID(id int) (internal.Employee, error) {
	args := m.Called(id)
	return args.Get(0).(internal.Employee), args.Error(1)
}

func (m *mockEmployeeRepository) CreateEmployee(inputEmployee internal.EmployeeAttributes) (employee internal.Employee, err error) {
	args := m.Called(inputEmployee)
	return args.Get(0).(internal.Employee), args.Error(1)
}

func (m *mockEmployeeRepository) UpdateEmployee(newEmployee internal.Employee) (internal.Employee, error) {
	args := m.Called(newEmployee)
	return args.Get(0).(internal.Employee), args.Error(1)
}

func (m *mockEmployeeRepository) DeleteEmployee(id int) (err error) {
	args := m.Called(id)
	return args.Error(0)
}

type mockWarehouseValidation struct {
	mock.Mock
}

func (m *mockWarehouseValidation) GetByID(id int) (internal.Warehouse, error) {
	args := m.Called(id)
	return args.Get(0).(internal.Warehouse), args.Error(1)
}

var (
	mockEmployee = internal.Employee{
		ID: 1,
		Attributes: internal.EmployeeAttributes{
			CardNumberID: "12345",
			FirstName:    "Aelin",
			LastName:     "Galanthynius",
			WarehouseID:  1,
		},
	}
	mockEmployee2 = internal.Employee{
		ID: 2,
		Attributes: internal.EmployeeAttributes{
			CardNumberID: "67890",
			FirstName:    "Rowan",
			LastName:     "Withethorn",
			WarehouseID:  1,
		},
	}
	mockEmployee3 = internal.Employee{
		ID: 1,
		Attributes: internal.EmployeeAttributes{
			CardNumberID: "12345",
			FirstName:    "Celaena",
			LastName:     "Sardothien",
			WarehouseID:  1,
		},
	}
	mockEmployeeAttr = internal.EmployeeAttributes{
		CardNumberID: "67890",
		FirstName:    "Rowan",
		LastName:     "Withethorn",
		WarehouseID:  1,
	}
	mockInputEmployee = internal.Employee{
		ID: 1,
		Attributes: internal.EmployeeAttributes{
			CardNumberID: "12345",
			FirstName:    "Celaena",
			LastName:     "Sardothien",
			WarehouseID:  1,
		},
	}
	mockInputEmployeeInvalidID = internal.Employee{
		ID: 99,
		Attributes: internal.EmployeeAttributes{
			CardNumberID: "12345",
			FirstName:    "Celaena",
			LastName:     "Sardothien",
			WarehouseID:  1,
		},
	}
	mockWarehouse = internal.Warehouse{
		ID:                 1,
		Address:            "Terrasen 452",
		Telephone:          "0123456789",
		WarehouseCode:      "XYZ",
		MinimumCapacity:    10,
		MinimumTemperature: 10,
	}
	mockErr1Employee = internal.Employee{
		ID: 1,
		Attributes: internal.EmployeeAttributes{
			CardNumberID: "",
			FirstName:    "",
			LastName:     "",
			WarehouseID:  1,
		},
	}
	mockNilInputEmployee = internal.Employee{
		ID: 1,
		Attributes: internal.EmployeeAttributes{
			WarehouseID: 0,
		},
	}
	mockEmployeeInvalidWarehouse = internal.Employee{
		ID: 1,
		Attributes: internal.EmployeeAttributes{
			CardNumberID: "12345",
			FirstName:    "Aelin",
			LastName:     "Galanthynius",
			WarehouseID:  99,
		},
	}
)

func TestEmployeeService_FindAll(t *testing.T) {
	t.Run("FindAll - Success", func(t *testing.T) {
		mockRepo := new(mockEmployeeRepository)
		mockRepo.On("FindAll").Return(map[int]internal.Employee{1: mockEmployee, 2: mockEmployee2}, nil)
		service := NewEmployeeService(mockRepo, nil)
		result, err := service.FindAll()

		assert.Equal(t, map[int]internal.Employee{1: mockEmployee, 2: mockEmployee2}, result)
		assert.Nil(t, err)
	})

	t.Run("FindAll - Error", func(t *testing.T) {
		mockRepo := new(mockEmployeeRepository)
		mockRepo.On("FindAll").Return(map[int]internal.Employee{}, assert.AnError)
		service := NewEmployeeService(mockRepo, nil)
		result, err := service.FindAll()

		assert.Equal(t, map[int]internal.Employee{}, result)
		assert.Error(t, err)
	})
}

func TestEmployeeService_FindById(t *testing.T) {
	t.Run("FindByID - Valid ID", func(t *testing.T) {
		mockRepo := new(mockEmployeeRepository)
		mockRepo.On("FindByID", 1).Return(mockEmployee, nil)
		service := NewEmployeeService(mockRepo, nil)
		result, err := service.FindByID(1)

		assert.Equal(t, mockEmployee, result)
		assert.Nil(t, err)
	})

	t.Run("FindByID - Invalid ID", func(t *testing.T) {
		mockRepo := new(mockEmployeeRepository)
		mockRepo.On("FindByID", 99).Return(internal.Employee{}, utils.ErrNotFound)
		service := NewEmployeeService(mockRepo, nil)
		result, err := service.FindByID(99)

		assert.Equal(t, internal.Employee{}, result)
		assert.Equal(t, utils.ErrNotFound, err)
	})

	t.Run("FindByID - Internal Error", func(t *testing.T) {
		mockRepo := new(mockEmployeeRepository)
		mockRepo.On("FindByID", 1).Return(internal.Employee{}, assert.AnError)
		service := NewEmployeeService(mockRepo, nil)
		result, err := service.FindByID(1)

		assert.Equal(t, internal.Employee{}, result)
		assert.NotNil(t, err)
	})
}

func TestEmployeeService_Delete(t *testing.T) {
	t.Run("Delete - Valid ID", func(t *testing.T) {
		mockRepo := new(mockEmployeeRepository)
		mockRepo.On("FindByID", 1).Return(mockEmployee, nil)
		mockRepo.On("DeleteEmployee", 1).Return(nil)
		service := NewEmployeeService(mockRepo, nil)
		err := service.DeleteEmployee(1)

		assert.Nil(t, err)
	})

	t.Run("Delete - Invalid ID", func(t *testing.T) {
		mockRepo := new(mockEmployeeRepository)
		mockRepo.On("FindByID", 99).Return(internal.Employee{}, assert.AnError)
		mockRepo.On("DeleteEmployee", 99).Return(utils.ErrNotFound)
		service := NewEmployeeService(mockRepo, nil)
		err := service.DeleteEmployee(99)

		assert.Equal(t, utils.ErrNotFound, err)
	})

	t.Run("Delete - Internal Error", func(t *testing.T) {
		mockRepo := new(mockEmployeeRepository)
		mockRepo.On("FindByID", 99).Return(internal.Employee{}, assert.AnError)
		mockRepo.On("DeleteEmployee", 99).Return(assert.AnError)
		service := NewEmployeeService(mockRepo, nil)
		err := service.DeleteEmployee(99)

		assert.NotNil(t, err)
	})

	t.Run("Delete - Invalid Arguments", func(t *testing.T) {
		mockRepo := new(mockEmployeeRepository)
		mockRepo.On("FindByID", 1).Return(mockEmployee, nil)
		mockRepo.On("DeleteEmployee", 1).Return(assert.AnError)
		service := NewEmployeeService(mockRepo, nil)
		err := service.DeleteEmployee(1)

		assert.NotNil(t, err)
		assert.Equal(t, utils.ErrInvalidArguments, err)
	})
}

func TestEmployeeService_Update(t *testing.T) {
	t.Run("Update - Valid ID", func(t *testing.T) {
		mockRepo := new(mockEmployeeRepository)
		mockWV := new(mockWarehouseValidation)
		mockWV.On("GetByID", 1).Return(mockWarehouse, nil)
		mockRepo.On("FindByID", 1).Return(mockEmployee, nil)
		mockRepo.On("UpdateEmployee", mockInputEmployee).Return(mockInputEmployee, nil)
		service := NewEmployeeService(mockRepo, mockWV)
		result, err := service.UpdateEmployee(mockInputEmployee)

		assert.Equal(t, mockInputEmployee, result)
		assert.Nil(t, err)
	})

	t.Run("Update - Invalid ID", func(t *testing.T) {
		mockRepo := new(mockEmployeeRepository)
		mockWV := new(mockWarehouseValidation)
		mockWV.On("GetByID", 1).Return(mockWarehouse, nil)
		mockRepo.On("FindByID", 99).Return(internal.Employee{}, utils.ErrNotFound)
		mockRepo.On("UpdateEmployee", mockInputEmployeeInvalidID).Return(internal.Employee{}, utils.ErrNotFound)
		service := NewEmployeeService(mockRepo, mockWV)
		result, err := service.UpdateEmployee(mockInputEmployeeInvalidID)

		assert.Equal(t, internal.Employee{}, result)
		assert.NotNil(t, err)
		assert.Equal(t, utils.ErrNotFound, err)
	})

	t.Run("Update - Warehouse Not Found", func(t *testing.T) {
		mockRepo := new(mockEmployeeRepository)
		mockWV := new(mockWarehouseValidation)
		mockWV.On("GetByID", 99).Return(internal.Warehouse{}, utils.ErrNotFound)
		mockRepo.On("FindByID", 99).Return(internal.Employee{}, assert.AnError)
		service := NewEmployeeService(mockRepo, mockWV)

		result, err := service.UpdateEmployee(mockInputEmployeeInvalidID)

		assert.Equal(t, internal.Employee{}, result)
		assert.NotNil(t, err)
		assert.Equal(t, utils.ErrNotFound, err)
	})

	t.Run("Update - Merge Employee Fields", func(t *testing.T) {
		mockRepo := new(mockEmployeeRepository)
		mockWV := new(mockWarehouseValidation)
		mockWV.On("GetByID", 0).Return(mockWarehouse, nil)
		mockRepo.On("FindByID", 1).Return(mockEmployee3, nil)
		mockRepo.On("UpdateEmployee", mockInputEmployee).Return(mockEmployee3, nil)
		service := NewEmployeeService(mockRepo, mockWV)
		result, err := service.UpdateEmployee(mockNilInputEmployee)

		assert.Equal(t, mockEmployee3, result)
		assert.Nil(t, err)

		assert.Equal(t, mockInputEmployee.ID, mockEmployee3.ID)
		assert.Equal(t, "Celaena", mockEmployee3.Attributes.FirstName)
		assert.Equal(t, "Sardothien", mockEmployee3.Attributes.LastName)
		assert.Equal(t, "12345", mockEmployee3.Attributes.CardNumberID)
		assert.Equal(t, 1, mockEmployee3.Attributes.WarehouseID)
	})

	t.Run("Update - Error Warehouse Nil", func(t *testing.T) {
		mockRepo := new(mockEmployeeRepository)
		mockWV := new(mockWarehouseValidation)
		mockWV.On("GetByID", 1).Return(internal.Warehouse{}, utils.EDependencyNotFound("warehouse", "id: "+"1"))
		mockRepo.On("FindByID", 1).Return(mockEmployee, nil)
		mockRepo.On("UpdateEmployee", mockEmployee).Return(internal.Employee{}, assert.AnError)
		service := NewEmployeeService(mockRepo, mockWV)
		result, err := service.UpdateEmployee(mockEmployee)

		assert.Equal(t, internal.Employee{}, result)
		assert.NotNil(t, err)
	})
}

func TestEmployeeService_Create(t *testing.T) {
	t.Run("Create - Success", func(t *testing.T) {
		mockRepo := new(mockEmployeeRepository)
		mockWV := new(mockWarehouseValidation)
		mockWV.On("GetByID", 1).Return(mockWarehouse, nil)
		mockRepo.On("FindAll").Return(map[int]internal.Employee{1: mockEmployee}, nil)
		mockRepo.On("CreateEmployee", mockEmployeeAttr).Return(mockEmployee2, nil)
		service := NewEmployeeService(mockRepo, mockWV)
		result, err := service.CreateEmployee(mockEmployeeAttr)

		assert.Equal(t, mockEmployee2, result)
		assert.Nil(t, err)
	})

	t.Run("Create - Conflict CardNumberID", func(t *testing.T) {
		mockRepo := new(mockEmployeeRepository)
		mockWV := new(mockWarehouseValidation)
		mockWV.On("GetByID", 1).Return(mockWarehouse, nil)
		mockRepo.On("FindAll").Return(map[int]internal.Employee{1: mockEmployee}, nil)
		mockRepo.On("CreateEmployee", mockEmployeeAttr).Return(internal.Employee{}, utils.ErrConflict)
		service := NewEmployeeService(mockRepo, mockWV)
		result, err := service.CreateEmployee(mockEmployeeAttr)

		assert.Equal(t, internal.Employee{}, result)
		assert.Equal(t, utils.ErrConflict, err)
	})

	t.Run("Create - Error on Validate Fiels", func(t *testing.T) {
		mockRepo := new(mockEmployeeRepository)
		mockWV := new(mockWarehouseValidation)
		mockWV.On("GetByID", 1).Return(mockWarehouse, nil)
		mockRepo.On("FindAll").Return(map[int]internal.Employee{1: mockEmployee2}, nil)
		mockRepo.On("CreateEmployee", mockEmployeeAttr).Return(internal.Employee{}, utils.ErrConflict)
		service := NewEmployeeService(mockRepo, mockWV)
		result, err := service.CreateEmployee(mockEmployeeAttr)

		assert.Equal(t, internal.Employee{}, result)
		assert.Equal(t, utils.ErrConflict, err)
	})

	t.Run("Create - Empty Fields", func(t *testing.T) {
		mockRepo := new(mockEmployeeRepository)
		mockWV := new(mockWarehouseValidation)
		service := NewEmployeeService(mockRepo, mockWV)

		result, err := service.CreateEmployee(mockErr1Employee.Attributes)

		assert.Equal(t, internal.Employee{}, result)
		assert.Equal(t, utils.ErrEmptyArguments, err)
	})

	t.Run("Create - Err on FindAll", func(t *testing.T) {
		mockRepo := new(mockEmployeeRepository)
		mockWV := new(mockWarehouseValidation)
		mockWV.On("GetByID", 1).Return(mockWarehouse, nil)
		mockRepo.On("FindAll").Return(map[int]internal.Employee{}, assert.AnError)
		mockRepo.On("CreateEmployee", mockEmployeeAttr).Return(internal.Employee{}, assert.AnError)
		service := NewEmployeeService(mockRepo, mockWV)
		result, err := service.CreateEmployee(mockEmployeeAttr)

		assert.Equal(t, internal.Employee{}, result)
		assert.NotNil(t, err)
	})

	t.Run("Create - Another Error on Warehouse", func(t *testing.T) {
		mockRepo := new(mockEmployeeRepository)
		mockWV := new(mockWarehouseValidation)
		mockWV.On("GetByID", 1).Return(mockWarehouse, assert.AnError)
		mockRepo.On("FindAll").Return(map[int]internal.Employee{}, nil)
		mockRepo.On("CreateEmployee", mockEmployeeAttr).Return(internal.Employee{}, assert.AnError)
		service := NewEmployeeService(mockRepo, mockWV)
		result, err := service.CreateEmployee(mockEmployeeAttr)

		assert.Equal(t, internal.Employee{}, result)
		assert.NotNil(t, err)
	})
}
