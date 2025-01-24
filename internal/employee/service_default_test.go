package employee

import (
	"log"
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
			WarehouseID:  2,
		},
	}
	mockEmployeeAttr = internal.EmployeeAttributes{
		CardNumberID: "12345",
		FirstName:    "Aelin",
		LastName:     "Galanthynius",
		WarehouseID:  1,
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
	// t.Run("FindAll - Success", func(t *testing.T) {
	// 	mockRepo := new(mockEmployeeRepository)
	// 	mockWV := new(mockWarehouseValidation)

	// 	employees := map[int]internal.Employee{
	// 		1: {
	// 			ID: 1,
	// 			Attributes: internal.EmployeeAttributes{
	// 				CardNumberID: "12345",
	// 				FirstName:    "Aelin",
	// 				LastName:     "Galanthynius",
	// 				WarehouseID:  1,
	// 			},
	// 		},
	// 		2: {
	// 			ID: 2,
	// 			Attributes: internal.EmployeeAttributes{
	// 				CardNumberID: "67890",
	// 				FirstName:    "Rowan",
	// 				LastName:     "Withethorn",
	// 				WarehouseID:  2,
	// 			},
	// 		},
	// 	}

	// 	mockRepo.On("FindAll").Return(employees, nil)

	// 	service := NewEmployeeService(mockRepo, mockWV)

	// 	result, err := service.FindAll()

	// 	assert.NoError(t, err)
	// 	assert.Equal(t, employees, result)

	// 	mockRepo.AssertExpectations(t)
	// })

	// t.Run("FindAll - Error", func(t *testing.T) {
	// 	mockRepo := new(mockEmployeeRepository)
	// 	mockWV := new(mockWarehouseValidation)

	// 	mockRepo.On("FindAll").Return(nil, assert.AnError)

	// 	service := NewEmployeeService(mockRepo, mockWV)

	// 	result, err := service.FindAll()

	// 	assert.Error(t, err)
	// 	assert.Nil(t, result)

	// 	mockRepo.AssertExpectations(t)
	// })

	t.Run("FindAll - Success", func(t *testing.T) {
		mockRepo := new(mockEmployeeRepository)
		mockRepo.On("FindAll").Return([]internal.Employee{mockEmployee, mockEmployee2}, nil)
		service := NewEmployeeService(mockRepo, nil)
		result, err := service.FindAll()

		assert.Equal(t, []internal.Employee{mockEmployee, mockEmployee2}, result)
		assert.Nil(t, err)
	})

	t.Run("FindAll - Error", func(t *testing.T) {
		mockRepo := new(mockEmployeeRepository)
		mockRepo.On("FindAll").Return(map[int]internal.Employee{}, assert.AnError)
		service := NewEmployeeService(mockRepo, nil)
		result, err := service.FindAll()

		log.Println("ERRO: ", err)
		log.Println("RESULT:", result)

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

	t.Run("FindByID - Another Error", func(t *testing.T) {
		mockRepo := new(mockEmployeeRepository)
		mockRepo.On("FindByID", 99).Return(internal.Employee{}, assert.AnError)
		service := NewEmployeeService(mockRepo, nil)
		result, err := service.FindByID(99)

		assert.Equal(t, internal.Employee{}, result)
		assert.NotNil(t, err)
	})
}
