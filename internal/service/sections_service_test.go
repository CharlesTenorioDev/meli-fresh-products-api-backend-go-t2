package service

import (
	"errors"
	"fmt"
	"testing"

	"github.com/meli-fresh-products-api-backend-go-t2/internal/pkg"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockSectionRepository struct {
	mock.Mock
}

func (m *MockSectionRepository) GetAll() ([]pkg.Section, error) {
	args := m.Called()
	return args.Get(0).([]pkg.Section), args.Error(1)
}

func (m *MockSectionRepository) Save(section pkg.Section) (pkg.Section, error) {
	args := m.Called(section)
	return args.Get(0).(pkg.Section), args.Error(1)
}

func (m *MockSectionRepository) Update(section pkg.Section) (pkg.Section, error) {
	args := m.Called(section)
	return args.Get(0).(pkg.Section), args.Error(1)
}

func (m *MockSectionRepository) GetById(id int) (pkg.Section, error) {
	args := m.Called(id)
	return args.Get(0).(pkg.Section), args.Error(1)
}

func (m *MockSectionRepository) GetBySectionNumber(sectionNumber int) (pkg.Section, error) {
	args := m.Called(sectionNumber)
	return args.Get(0).(pkg.Section), args.Error(1)
}

func (m *MockSectionRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

type MockSectionWarehouseService struct {
	mock.Mock
}

func (m *MockSectionWarehouseService) GetById(id int) (pkg.Warehouse, error) {
	args := m.Called(id)
	return args.Get(0).(pkg.Warehouse), args.Error(1)
}

type MockSectionProductTypeService struct {
	mock.Mock
}

func (m *MockSectionProductTypeService) GetProductTypeByID(id int) (pkg.ProductType, error) {
	args := m.Called(id)
	return args.Get(0).(pkg.ProductType), args.Error(1)
}

var (
	simpleSection = pkg.Section{
		ID:                 1,
		SectionNumber:      1,
		CurrentTemperature: 1,
		MinimumTemperature: 1,
		CurrentCapacity:    1,
		MinimumCapacity:    1,
		MaximumCapacity:    1,
		WarehouseID:        1,
		ProductTypeID:      1,
	}
	simpleSectionPointers = pkg.SectionPointers{
		SectionNumber:      &simpleSection.SectionNumber,
		CurrentTemperature: &simpleSection.CurrentTemperature,
		MinimumTemperature: &simpleSection.MinimumTemperature,
		CurrentCapacity:    &simpleSection.CurrentCapacity,
		MinimumCapacity:    &simpleSection.MinimumCapacity,
		MaximumCapacity:    &simpleSection.MaximumCapacity,
		WarehouseID:        &simpleSection.WarehouseID,
		ProductTypeID:      &simpleSection.ProductTypeID,
	}
	simpleWarehouse = pkg.Warehouse{
		ID:                 1,
		Address:            "Monroe 860",
		Telephone:          "00900009999",
		WarehouseCode:      "DHM",
		MinimumCapacity:    10,
		MinimumTemperature: 10,
	}
	simpleProductType = pkg.ProductType{
		ID:          1,
		Description: "Foo Bar",
	}
)

func Test_GetAll(t *testing.T) {
	repo := new(MockSectionRepository)
	repo.On("GetAll").Return([]pkg.Section{simpleSection}, nil)
	service := NewBasicSectionService(repo, nil, nil)
	sections, _ := service.GetAll()
	require.Equal(t, 1, len(sections))
}

func Test_GetById(t *testing.T) {
	t.Run("when exist", func(s *testing.T) {
		repo := new(MockSectionRepository)
		repo.On("GetById", 1).Return(simpleSection, nil)
		service := NewBasicSectionService(repo, nil, nil)
		section, err := service.GetById(1)
		require.Equal(t, simpleSection, section)
		require.Nil(t, err)
	})
	t.Run("when does not exist", func(s *testing.T) {
		repo := new(MockSectionRepository)
		repo.On("GetById", 1).Return(pkg.Section{}, nil)
		service := NewBasicSectionService(repo, nil, nil)
		section, err := service.GetById(1)
		require.Empty(t, section)
		require.ErrorIs(t, err, utils.ErrNotFound)
	})
	t.Run("when internal error occurs", func(s *testing.T) {
		repo := new(MockSectionRepository)
		repo.On("GetById", 1).Return(pkg.Section{}, errors.New("internal error"))
		service := NewBasicSectionService(repo, nil, nil)
		section, err := service.GetById(1)
		require.Empty(t, section)
		require.Equal(t, err.Error(), "internal error")
	})
}

func Test_Delete(t *testing.T) {
	t.Run("when exist", func(s *testing.T) {
		repo := new(MockSectionRepository)
		repo.On("GetById", 1).Return(simpleSection, nil)
		repo.On("Delete", 1).Return(nil)
		service := NewBasicSectionService(repo, nil, nil)
		err := service.Delete(1)
		require.Nil(t, err)
	})
	t.Run("when does not exist", func(s *testing.T) {
		repo := new(MockSectionRepository)
		repo.On("GetById", 1).Return(pkg.Section{}, nil)
		service := NewBasicSectionService(repo, nil, nil)
		err := service.Delete(1)
		require.ErrorIs(t, err, utils.ErrNotFound)
	})
	t.Run("when internal error occurs at getById", func(s *testing.T) {
		repo := new(MockSectionRepository)
		repo.On("GetById", 1).Return(pkg.Section{}, errors.New("internal error"))
		service := NewBasicSectionService(repo, nil, nil)
		err := service.Delete(1)
		require.Equal(t, err.Error(), "internal error")
	})
	t.Run("when internal error occurs at delete", func(s *testing.T) {
		repo := new(MockSectionRepository)
		repo.On("GetById", 1).Return(simpleSection, nil)
		repo.On("Delete", 1).Return(errors.New("internal error"))
		service := NewBasicSectionService(repo, nil, nil)
		err := service.Delete(1)
		require.Equal(t, err.Error(), "internal error")
	})
}

func Test_Save(t *testing.T) {
	scenarios := []struct {
		Name string
		Data pkg.Section
		Mock func(*MockSectionRepository, *MockSectionWarehouseService, *MockSectionProductTypeService) (pkg.Section, error)
	}{
		{
			Name: "when ok",
			Data: simpleSection,
			Mock: func(repo *MockSectionRepository, warehouseService *MockSectionWarehouseService, productTypeService *MockSectionProductTypeService) (pkg.Section, error) {
				repo.On("Save", mock.Anything).Return(simpleSection, nil)
				repo.On("GetBySectionNumber", mock.Anything).Return(pkg.Section{}, nil)
				warehouseService.On("GetById", mock.Anything).Return(simpleWarehouse, nil)
				productTypeService.On("GetProductTypeByID", mock.Anything).Return(simpleProductType, nil)
				return simpleSection, nil
			},
		},
		{
			Name: "when zero value section_number",
			Data: pkg.Section{},
			Mock: func(repo *MockSectionRepository, warehouseService *MockSectionWarehouseService, productTypeService *MockSectionProductTypeService) (pkg.Section, error) {
				repo.On("Save", mock.Anything).Return(simpleSection, nil)
				repo.On("GetBySectionNumber", mock.Anything).Return(pkg.Section{}, nil)
				warehouseService.On("GetById", mock.Anything).Return(simpleWarehouse, nil)
				productTypeService.On("GetProductTypeByID", mock.Anything).Return(simpleProductType, nil)
				return pkg.Section{}, errors.Join(utils.ErrInvalidArguments, errors.New("section_number cannot be empty/null"))
			},
		},
		{
			Name: "when zero value warehouse_id",
			Data: pkg.Section{SectionNumber: 1},
			Mock: func(repo *MockSectionRepository, warehouseService *MockSectionWarehouseService, productTypeService *MockSectionProductTypeService) (pkg.Section, error) {
				return pkg.Section{}, errors.Join(utils.ErrInvalidArguments, errors.New("warehouse_id cannot be empty/null"))
			},
		},
		{
			Name: "when zero value product_type_id",
			Data: pkg.Section{SectionNumber: 1, WarehouseID: 1},
			Mock: func(repo *MockSectionRepository, warehouseService *MockSectionWarehouseService, productTypeService *MockSectionProductTypeService) (pkg.Section, error) {
				return pkg.Section{}, errors.Join(utils.ErrInvalidArguments, errors.New("product_type_id cannot be empty/null"))
			},
		},
		{
			Name: "when warehouse does not exist",
			Data: pkg.Section{SectionNumber: 1, WarehouseID: 1, ProductTypeID: 1},
			Mock: func(repo *MockSectionRepository, warehouseService *MockSectionWarehouseService, productTypeService *MockSectionProductTypeService) (pkg.Section, error) {
				warehouseService.On("GetById", mock.Anything).Return(pkg.Warehouse{}, utils.ErrNotFound)
				return pkg.Section{}, errors.Join(utils.ErrInvalidArguments, fmt.Errorf("warehouse not found for id %d", 1))
			},
		},
		{
			Name: "when product_type does not exist",
			Data: pkg.Section{SectionNumber: 1, WarehouseID: 1, ProductTypeID: 1},
			Mock: func(repo *MockSectionRepository, warehouseService *MockSectionWarehouseService, productTypeService *MockSectionProductTypeService) (pkg.Section, error) {
				warehouseService.On("GetById", mock.Anything).Return(simpleWarehouse, nil)
				productTypeService.On("GetProductTypeByID", mock.Anything).Return(pkg.ProductType{}, utils.ErrNotFound)
				return pkg.Section{}, errors.Join(utils.ErrInvalidArguments, fmt.Errorf("product_type not found for id %d", 1))
			},
		},
		{
			Name: "when minimum_capacity is greater than maximum_capacity",
			Data: pkg.Section{SectionNumber: 1, WarehouseID: 1, ProductTypeID: 1, MinimumCapacity: 5, MaximumCapacity: 4},
			Mock: func(repo *MockSectionRepository, warehouseService *MockSectionWarehouseService, productTypeService *MockSectionProductTypeService) (pkg.Section, error) {
				warehouseService.On("GetById", mock.Anything).Return(simpleWarehouse, nil)
				productTypeService.On("GetProductTypeByID", mock.Anything).Return(simpleProductType, nil)
				return pkg.Section{}, errors.Join(utils.ErrInvalidArguments, errors.New("minimum_capacity cannot be greater than maximum_capacity"))
			},
		},
		{
			Name: "when minimum_temperature is less than -273.15 Celsius",
			Data: pkg.Section{SectionNumber: 1, WarehouseID: 1, ProductTypeID: 1, MinimumCapacity: 3, MaximumCapacity: 4, MinimumTemperature: -300},
			Mock: func(repo *MockSectionRepository, warehouseService *MockSectionWarehouseService, productTypeService *MockSectionProductTypeService) (pkg.Section, error) {
				warehouseService.On("GetById", mock.Anything).Return(simpleWarehouse, nil)
				productTypeService.On("GetProductTypeByID", mock.Anything).Return(simpleProductType, nil)
				return pkg.Section{}, errors.Join(utils.ErrInvalidArguments, errors.New("minimum_temperature cannot be less than -273.15 Celsius"))
			},
		},
		{
			Name: "when current_temperature is less than -273.15 Celsius",
			Data: pkg.Section{SectionNumber: 1, WarehouseID: 1, ProductTypeID: 1, MinimumCapacity: 3, MaximumCapacity: 4, MinimumTemperature: 0, CurrentTemperature: -300},
			Mock: func(repo *MockSectionRepository, warehouseService *MockSectionWarehouseService, productTypeService *MockSectionProductTypeService) (pkg.Section, error) {
				warehouseService.On("GetById", mock.Anything).Return(simpleWarehouse, nil)
				productTypeService.On("GetProductTypeByID", mock.Anything).Return(simpleProductType, nil)
				return pkg.Section{}, errors.Join(utils.ErrInvalidArguments, errors.New("current_temperature cannot be less than -273.15 Celsius"))
			},
		},
		{
			Name: "when section_number already exists",
			Data: simpleSection,
			Mock: func(repo *MockSectionRepository, warehouseService *MockSectionWarehouseService, productTypeService *MockSectionProductTypeService) (pkg.Section, error) {
				repo.On("GetBySectionNumber", mock.Anything).Return(pkg.Section{ID: 2}, nil)
				warehouseService.On("GetById", mock.Anything).Return(simpleWarehouse, nil)
				productTypeService.On("GetProductTypeByID", mock.Anything).Return(simpleProductType, nil)
				return pkg.Section{}, utils.ErrConflict
			},
		},
	}
	for _, scenario := range scenarios {
		t.Run(scenario.Name, func(m *testing.T) {
			// Create the mocks
			repo := new(MockSectionRepository)
			warehouseService := new(MockSectionWarehouseService)
			productTypeService := new(MockSectionProductTypeService)

			// Mock and get the expected
			expectedData, expectedError := scenario.Mock(repo, warehouseService, productTypeService)

			service := NewBasicSectionService(repo, warehouseService, productTypeService)
			savedSection, err := service.Save(scenario.Data)
			require.Equal(m, expectedData, savedSection)
			if expectedError == nil {
				require.Nil(m, err)
			} else {
				require.Equal(m, err.Error(), expectedError.Error())
			}
		})
	}
}

func Test_Update(t *testing.T) {
	scenarios := []struct {
		Name   string
		Data   pkg.SectionPointers
		DataId int
		Mock   func(*MockSectionRepository, *MockSectionWarehouseService, *MockSectionProductTypeService) (pkg.Section, error)
	}{
		{
			Name:   "when ok",
			Data:   simpleSectionPointers,
			DataId: 1,
			Mock: func(repo *MockSectionRepository, warehouseService *MockSectionWarehouseService, productTypeService *MockSectionProductTypeService) (pkg.Section, error) {
				repo.On("Update", mock.Anything, mock.Anything).Return(simpleSection, nil)
				repo.On("GetById", mock.Anything).Return(simpleSection, nil)
				repo.On("GetBySectionNumber", mock.Anything).Return(pkg.Section{}, nil)
				warehouseService.On("GetById", mock.Anything).Return(simpleWarehouse, nil)
				productTypeService.On("GetProductTypeByID", mock.Anything).Return(simpleProductType, nil)
				return simpleSection, nil
			},
		},
	}
	for _, scenario := range scenarios {
		t.Run(scenario.Name, func(m *testing.T) {
			// Create the mocks
			repo := new(MockSectionRepository)
			warehouseService := new(MockSectionWarehouseService)
			productTypeService := new(MockSectionProductTypeService)

			// Mock and get the expected
			expectedData, expectedError := scenario.Mock(repo, warehouseService, productTypeService)

			service := NewBasicSectionService(repo, warehouseService, productTypeService)
			savedSection, err := service.Update(scenario.DataId, scenario.Data)
			require.Equal(m, expectedData, savedSection)
			if expectedError == nil {
				require.Nil(m, err)
			} else {
				require.Equal(m, err.Error(), expectedError.Error())
			}
		})
	}
}
