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
}
