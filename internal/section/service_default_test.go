package section

import (
	"errors"
	"testing"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"

	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockSectionRepository struct {
	mock.Mock
}

func (m *MockSectionRepository) GetAll() ([]internal.Section, error) {
	args := m.Called()
	return args.Get(0).([]internal.Section), args.Error(1)
}

func (m *MockSectionRepository) Save(section *internal.Section) error {
	args := m.Called(section)
	return args.Error(0)
}

func (m *MockSectionRepository) Update(section *internal.Section) error {
	args := m.Called(section)
	return args.Error(0)
}

func (m *MockSectionRepository) GetByID(id int) (internal.Section, error) {
	args := m.Called(id)
	return args.Get(0).(internal.Section), args.Error(1)
}

func (m *MockSectionRepository) GetBySectionNumber(sectionNumber int) (internal.Section, error) {
	args := m.Called(sectionNumber)
	return args.Get(0).(internal.Section), args.Error(1)
}

func (m *MockSectionRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockSectionRepository) GetSectionProductsReport() ([]internal.SectionProductsReport, error) {
	args := m.Called()
	return args.Get(0).([]internal.SectionProductsReport), args.Error(1)
}

func (m *MockSectionRepository) GetSectionProductsReportByID(id int) ([]internal.SectionProductsReport, error) {
	args := m.Called(id)
	return args.Get(0).([]internal.SectionProductsReport), args.Error(1)
}

type MockSectionWarehouseService struct {
	mock.Mock
}

func (m *MockSectionWarehouseService) GetByID(id int) (internal.Warehouse, error) {
	args := m.Called(id)
	return args.Get(0).(internal.Warehouse), args.Error(1)
}

type MockSectionProductTypeService struct {
	mock.Mock
}

func (m *MockSectionProductTypeService) GetProductTypeByID(id int) (internal.ProductType, error) {
	args := m.Called(id)
	return args.Get(0).(internal.ProductType), args.Error(1)
}

var (
	mockSection = internal.Section{
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
	mockSection2 = internal.Section{
		ID:                 2,
		SectionNumber:      2,
		CurrentTemperature: 2,
		MinimumTemperature: 2,
		CurrentCapacity:    2,
		MinimumCapacity:    2,
		MaximumCapacity:    2,
		WarehouseID:        2,
		ProductTypeID:      2,
	}
	zero          = 0
	one           = 1
	two           = 2
	mockWarehouse = internal.Warehouse{
		ID:                 1,
		Address:            "Monroe 860",
		Telephone:          "00900009999",
		WarehouseCode:      "DHM",
		MinimumCapacity:    10,
		MinimumTemperature: 10,
	}
	mockProductType = internal.ProductType{
		ID:          1,
		Description: "Foo Bar",
	}
	mockSectionProductsReport = []internal.SectionProductsReport{
		{
			SectionID:     1,
			SectionNumber: 923,
			ProductsCount: 20,
		},
	}
)

func TestUnitSection_GetAll(t *testing.T) {
	t.Run("WHEN repository returns no error, RETURN successfully", func(t *testing.T) {
		repo := new(MockSectionRepository)
		repo.On("GetAll").Return([]internal.Section{mockSection}, nil)
		service := NewBasicSectionService(repo, nil, nil)
		sections, _ := service.GetAll()
		require.Equal(t, 1, len(sections))
	})

	t.Run("WHEN repository returns some error, RETURN the error", func(t *testing.T) {
		repo := new(MockSectionRepository)
		repo.On("GetAll").Return([]internal.Section{}, errors.New("some error"))
		service := NewBasicSectionService(repo, nil, nil)
		sections, _ := service.GetAll()
		require.Equal(t, 0, len(sections))
	})

}

func TestUnitSection_GetById(t *testing.T) {
	t.Run("GIVEN a valid id, WHEN section exists, RETURN successfully", func(s *testing.T) {
		repo := new(MockSectionRepository)
		repo.On("GetByID", 1).Return(mockSection, nil)
		service := NewBasicSectionService(repo, nil, nil)
		section, err := service.GetByID(1)
		require.Equal(t, mockSection, section)
		require.Nil(t, err)
	})

	t.Run("GIVEN a valid id, WHEN section does not exists, RETURN successfully", func(s *testing.T) {
		repo := new(MockSectionRepository)
		repo.On("GetByID", 1).Return(internal.Section{}, nil)
		service := NewBasicSectionService(repo, nil, nil)
		section, err := service.GetByID(1)
		require.Empty(t, section)
		require.ErrorIs(t, err, utils.ErrNotFound)
	})

	t.Run("GIVEN a valid id, WHEN calling GetByID, RETURN internal error", func(s *testing.T) {
		repo := new(MockSectionRepository)
		repo.On("GetByID", 1).Return(internal.Section{}, errors.New("internal error"))
		service := NewBasicSectionService(repo, nil, nil)
		section, err := service.GetByID(1)
		require.Empty(t, section)
		require.Equal(t, err.Error(), "internal error")
	})
}

func TestUnitSection_Delete(t *testing.T) {
	t.Run("GIVEN a valid id, WHEN section exists, DELETE successfully", func(s *testing.T) {
		repo := new(MockSectionRepository)
		repo.On("GetByID", 1).Return(mockSection, nil)
		repo.On("Delete", 1).Return(nil)
		service := NewBasicSectionService(repo, nil, nil)
		err := service.Delete(1)
		require.Nil(t, err)
	})
	t.Run("GIVEN a valid id, WHEN section does not exists, RETURN utils.ErrNotFound", func(s *testing.T) {
		repo := new(MockSectionRepository)
		repo.On("GetByID", 1).Return(internal.Section{}, nil)
		service := NewBasicSectionService(repo, nil, nil)
		err := service.Delete(1)
		require.ErrorIs(t, err, utils.ErrNotFound)
	})
	t.Run("GIVEN a valid id, WHEN calling GetByID, RETURN internal error", func(s *testing.T) {
		repo := new(MockSectionRepository)
		repo.On("GetByID", 1).Return(internal.Section{}, errors.New("internal error"))
		service := NewBasicSectionService(repo, nil, nil)
		err := service.Delete(1)
		require.Equal(t, err.Error(), "internal error")
	})
	t.Run("GIVEN a valid id, WHEN calling GetByID, RETURN internal error", func(s *testing.T) {
		repo := new(MockSectionRepository)
		repo.On("GetByID", 1).Return(mockSection, nil)
		repo.On("Delete", 1).Return(errors.New("internal error"))
		service := NewBasicSectionService(repo, nil, nil)
		err := service.Delete(1)
		require.Equal(t, err.Error(), "internal error")
	})
}

func TestUnitSection_Save(t *testing.T) {
	internalError := errors.New("internal error")
	cases := []struct {
		Name string
		Data internal.Section
		Mock func(*MockSectionRepository, *MockSectionWarehouseService, *MockSectionProductTypeService) (internal.Section, error)
	}{
		{
			Name: "GIVEN a non valid section, WHEN zero value section_number, RETURNS error utils.ErrInvalidArguments",
			Data: internal.Section{},
			Mock: func(repo *MockSectionRepository, warehouseService *MockSectionWarehouseService, productTypeService *MockSectionProductTypeService) (internal.Section, error) {
				return internal.Section{}, utils.EZeroValue("section_number")
			},
		},
		{
			Name: "GIVEN a non valid section, WHEN zero value warehouse_id, RETURNS error utils.ErrInvalidArguments",
			Data: internal.Section{SectionNumber: 1},
			Mock: func(repo *MockSectionRepository, warehouseService *MockSectionWarehouseService, productTypeService *MockSectionProductTypeService) (internal.Section, error) {
				return internal.Section{}, utils.EZeroValue("warehouse_id")
			},
		},
		{
			Name: "GIVEN a non valid section, WHEN zero value product_type_id, RETURNS error utils.ErrInvalidArguments",
			Data: internal.Section{SectionNumber: 1, WarehouseID: 1},
			Mock: func(repo *MockSectionRepository, warehouseService *MockSectionWarehouseService, productTypeService *MockSectionProductTypeService) (internal.Section, error) {
				return internal.Section{}, utils.EZeroValue("product_type_id")
			},
		},
		{
			Name: "GIVEN a non valid section, WHEN warehouse does not exist, RETURNS error utils.ErrInvalidArguments",
			Data: internal.Section{SectionNumber: 1, WarehouseID: 1, ProductTypeID: 1},
			Mock: func(repo *MockSectionRepository, warehouseService *MockSectionWarehouseService, productTypeService *MockSectionProductTypeService) (internal.Section, error) {
				warehouseService.On("GetByID", mock.Anything).Return(internal.Warehouse{}, utils.ErrNotFound)
				return internal.Section{}, utils.EDependencyNotFound("warehouse", "id: 1")
			},
		},

		{
			Name: "GIVEN a non valid section, WHEN fail to call warehouseService.GetByID, RETURNS internal error",
			Data: internal.Section{SectionNumber: 1, WarehouseID: 1, ProductTypeID: 1},
			Mock: func(repo *MockSectionRepository, warehouseService *MockSectionWarehouseService, productTypeService *MockSectionProductTypeService) (internal.Section, error) {
				warehouseService.On("GetByID", mock.Anything).Return(internal.Warehouse{}, internalError)
				return internal.Section{}, internalError
			},
		},

		{
			Name: "GIVEN a non valid section, WHEN product_type does not exist, RETURNS error utils.ErrInvalidArguments",
			Data: internal.Section{SectionNumber: 1, WarehouseID: 1, ProductTypeID: 1},
			Mock: func(repo *MockSectionRepository, warehouseService *MockSectionWarehouseService, productTypeService *MockSectionProductTypeService) (internal.Section, error) {
				warehouseService.On("GetByID", mock.Anything).Return(mockWarehouse, nil)
				productTypeService.On("GetProductTypeByID", mock.Anything).Return(internal.ProductType{}, utils.ErrNotFound)
				return internal.Section{}, utils.EDependencyNotFound("product_type", "id: 1")
			},
		},

		{
			Name: "GIVEN a non valid section, WHEN fail to call productTypeService.GetProductTypeByID, RETURNS internal error",
			Data: internal.Section{SectionNumber: 1, WarehouseID: 1, ProductTypeID: 1},
			Mock: func(repo *MockSectionRepository, warehouseService *MockSectionWarehouseService, productTypeService *MockSectionProductTypeService) (internal.Section, error) {
				warehouseService.On("GetByID", mock.Anything).Return(mockWarehouse, nil)
				productTypeService.On("GetProductTypeByID", mock.Anything).Return(internal.ProductType{}, internalError)
				return internal.Section{}, internalError
			},
		},

		{
			Name: "GIVEN a non valid section, WHEN minimum_capacity is greater than maximum_capacity, RETURNS error utils.ErrInvalidArguments",
			Data: internal.Section{SectionNumber: 1, WarehouseID: 1, ProductTypeID: 1, MinimumCapacity: 5, MaximumCapacity: 4},
			Mock: func(repo *MockSectionRepository, warehouseService *MockSectionWarehouseService, productTypeService *MockSectionProductTypeService) (internal.Section, error) {
				warehouseService.On("GetByID", mock.Anything).Return(mockWarehouse, nil)
				productTypeService.On("GetProductTypeByID", mock.Anything).Return(mockProductType, nil)
				return internal.Section{}, utils.EBR("minimum_capacity cannot be greater than maximum_capacity")
			},
		},
		{
			Name: "GIVEN a non valid section, WHEN minimum_temperature is less than -273.15 Celsius, RETURNS error utils.ErrInvalidArguments",
			Data: internal.Section{SectionNumber: 1, WarehouseID: 1, ProductTypeID: 1, MinimumCapacity: 3, MaximumCapacity: 4, MinimumTemperature: -300},
			Mock: func(repo *MockSectionRepository, warehouseService *MockSectionWarehouseService, productTypeService *MockSectionProductTypeService) (internal.Section, error) {
				warehouseService.On("GetByID", mock.Anything).Return(mockWarehouse, nil)
				productTypeService.On("GetProductTypeByID", mock.Anything).Return(mockProductType, nil)
				return internal.Section{}, utils.EBR("minimum_temperature cannot be less than -273.15 Celsius")
			},
		},
		{
			Name: "GIVEN a non valid section, WHEN current_temperature is less than -273.15 Celsius, RETURNS error utils.ErrInvalidArguments",
			Data: internal.Section{SectionNumber: 1, WarehouseID: 1, ProductTypeID: 1, MinimumCapacity: 3, MaximumCapacity: 4, MinimumTemperature: 0, CurrentTemperature: -300},
			Mock: func(repo *MockSectionRepository, warehouseService *MockSectionWarehouseService, productTypeService *MockSectionProductTypeService) (internal.Section, error) {
				warehouseService.On("GetByID", mock.Anything).Return(mockWarehouse, nil)
				productTypeService.On("GetProductTypeByID", mock.Anything).Return(mockProductType, nil)
				return internal.Section{}, utils.EBR("current_temperature cannot be less than -273.15 Celsius")
			},
		},
		{
			Name: "GIVEN a non valid section, WHEN section_number already exists, RETURNS error utils.ErrConflict",
			Data: mockSection,
			Mock: func(repo *MockSectionRepository, warehouseService *MockSectionWarehouseService, productTypeService *MockSectionProductTypeService) (internal.Section, error) {
				repo.On("GetBySectionNumber", mock.Anything).Return(internal.Section{ID: 2}, nil)
				warehouseService.On("GetByID", mock.Anything).Return(mockWarehouse, nil)
				productTypeService.On("GetProductTypeByID", mock.Anything).Return(mockProductType, nil)
				return internal.Section{}, utils.EConflict("section", "id: 1")
			},
		},
		{
			Name: "GIVEN a non valid section, WHEN fail to call repo.GetBySectionNumber, RETURNS internal error",
			Data: mockSection,
			Mock: func(repo *MockSectionRepository, warehouseService *MockSectionWarehouseService, productTypeService *MockSectionProductTypeService) (internal.Section, error) {
				repo.On("GetBySectionNumber", mock.Anything).Return(internal.Section{}, internalError)
				warehouseService.On("GetByID", mock.Anything).Return(mockWarehouse, nil)
				productTypeService.On("GetProductTypeByID", mock.Anything).Return(mockProductType, nil)
				return internal.Section{}, internalError
			},
		},
		{
			Name: "GIVEN a valid section, WHEN no error occurs, SAVE successfully",
			Data: mockSection,
			Mock: func(repo *MockSectionRepository, warehouseService *MockSectionWarehouseService, productTypeService *MockSectionProductTypeService) (internal.Section, error) {
				repo.On("Save", mock.Anything).Return(nil)
				repo.On("GetBySectionNumber", mock.Anything).Return(internal.Section{}, nil)
				warehouseService.On("GetByID", mock.Anything).Return(mockWarehouse, nil)
				productTypeService.On("GetProductTypeByID", mock.Anything).Return(mockProductType, nil)
				return mockSection, nil
			},
		},
		{
			Name: "GIVEN a valid section, WHEN calling repo.Save, RETURNS internal error",
			Data: mockSection,
			Mock: func(repo *MockSectionRepository, warehouseService *MockSectionWarehouseService, productTypeService *MockSectionProductTypeService) (internal.Section, error) {
				repo.On("Save", mock.Anything).Return(internalError)
				repo.On("GetBySectionNumber", mock.Anything).Return(internal.Section{}, nil)
				warehouseService.On("GetByID", mock.Anything).Return(mockWarehouse, nil)
				productTypeService.On("GetProductTypeByID", mock.Anything).Return(mockProductType, nil)
				return internal.Section{}, internalError
			},
		},
	}
	for _, scenario := range cases {
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

func TestUnitSection_Update(t *testing.T) {
	internalError := errors.New("internal error")
	cases := []struct {
		Name   string
		DataID int
		Data   internal.SectionPointers
		Mock   func(*MockSectionRepository, *MockSectionWarehouseService, *MockSectionProductTypeService) (internal.Section, error)
	}{
		{
			Name:   "GIVEN a non valid section, WHEN calling GetByID, RETURNS internal error",
			DataID: 1,
			Data:   internal.SectionPointers{},
			Mock: func(repo *MockSectionRepository, warehouseService *MockSectionWarehouseService, productTypeService *MockSectionProductTypeService) (internal.Section, error) {
				repo.On("GetByID", mock.Anything).Return(internal.Section{}, internalError)
				return internal.Section{}, internalError
			},
		},
		{
			Name:   "GIVEN a non valid section, WHEN calling GetByID, RETURNS utils.ErrNotFound",
			DataID: 1,
			Data:   internal.SectionPointers{},
			Mock: func(repo *MockSectionRepository, warehouseService *MockSectionWarehouseService, productTypeService *MockSectionProductTypeService) (internal.Section, error) {
				repo.On("GetByID", mock.Anything).Return(internal.Section{}, utils.ENotFound("section"))
				return internal.Section{}, utils.ENotFound("section")
			},
		},
		{
			Name:   "GIVEN a non valid section, WHEN SectionNumber <= 0, RETURNS utils.ErrInvalidArguments",
			DataID: 1,
			Data:   internal.SectionPointers{SectionNumber: &zero},
			Mock: func(repo *MockSectionRepository, warehouseService *MockSectionWarehouseService, productTypeService *MockSectionProductTypeService) (internal.Section, error) {
				repo.On("GetByID", mock.Anything).Return(mockSection, nil)
				return internal.Section{}, utils.EZeroValue("section_number")
			},
		},
		{
			Name:   "GIVEN a non valid section, WHEN SectionNumber already exists, RETURNS utils.ErrConflict",
			DataID: 1,
			Data:   internal.SectionPointers{SectionNumber: &two},
			Mock: func(repo *MockSectionRepository, warehouseService *MockSectionWarehouseService, productTypeService *MockSectionProductTypeService) (internal.Section, error) {
				repo.On("GetByID", mock.Anything).Return(mockSection, nil)
				repo.On("GetBySectionNumber", mock.Anything).Return(mockSection2, nil)
				return internal.Section{}, utils.EConflict("section", "id: 2")
			},
		},
		{
			Name:   "GIVEN a non valid section, WHEN ProductTypeID <= 0, RETURNS utils.ErrInvalidArguments",
			DataID: 1,
			Data:   internal.SectionPointers{ProductTypeID: &zero},
			Mock: func(repo *MockSectionRepository, warehouseService *MockSectionWarehouseService, productTypeService *MockSectionProductTypeService) (internal.Section, error) {
				repo.On("GetByID", mock.Anything).Return(mockSection, nil)
				return internal.Section{}, utils.EZeroValue("product_type_id")
			},
		},
		{
			Name:   "GIVEN a non valid section, WHEN ProductTypeID already exists, RETURNS utils.ErrConflict",
			DataID: 1,
			Data:   internal.SectionPointers{ProductTypeID: &two},
			Mock: func(repo *MockSectionRepository, warehouseService *MockSectionWarehouseService, productTypeService *MockSectionProductTypeService) (internal.Section, error) {
				repo.On("GetByID", mock.Anything).Return(mockSection, nil)
				repo.On("GetBySectionNumber", mock.Anything).Return(mockSection2, nil)
				return internal.Section{}, utils.EConflict("section", "id: 2")
			},
		},
		//{
		//	Name:   "GIVEN a non valid section, WHEN calling GetByID, RETURNS utils.ErrNotFound",
		//	DataID: 1,
		//	Data: internal.SectionPointers{
		//		SectionNumber: &sectionNumber2,
		//	},
		//	Mock: func(repo *MockSectionRepository, warehouseService *MockSectionWarehouseService, productTypeService *MockSectionProductTypeService) (internal.Section, error) {
		//		repo.On("GetByID", mock.Anything).Return(internal.Section{}, utils.ENotFound("section"))
		//		return internal.Section{}, utils.ENotFound("section")
		//	},
		//},
		//{
		//	Name: "GIVEN a non valid section, WHEN zero value product_type_id, RETURNS error utils.ErrInvalidArguments",
		//	Data: internal.Section{SectionNumber: 1, WarehouseID: 1},
		//	Mock: func(repo *MockSectionRepository, warehouseService *MockSectionWarehouseService, productTypeService *MockSectionProductTypeService) (internal.Section, error) {
		//		return internal.Section{}, utils.EZeroValue("product_type_id")
		//	},
		//},
		//{
		//	Name: "GIVEN a non valid section, WHEN warehouse does not exist, RETURNS error utils.ErrInvalidArguments",
		//	Data: internal.Section{SectionNumber: 1, WarehouseID: 1, ProductTypeID: 1},
		//	Mock: func(repo *MockSectionRepository, warehouseService *MockSectionWarehouseService, productTypeService *MockSectionProductTypeService) (internal.Section, error) {
		//		warehouseService.On("GetByID", mock.Anything).Return(internal.Warehouse{}, utils.ErrNotFound)
		//		return internal.Section{}, utils.EDependencyNotFound("warehouse", "id: 1")
		//	},
		//},
		//{
		//	Name: "GIVEN a non valid section, WHEN fail to call warehouseService.GetByID, RETURNS internal error",
		//	Data: internal.Section{SectionNumber: 1, WarehouseID: 1, ProductTypeID: 1},
		//	Mock: func(repo *MockSectionRepository, warehouseService *MockSectionWarehouseService, productTypeService *MockSectionProductTypeService) (internal.Section, error) {
		//		warehouseService.On("GetByID", mock.Anything).Return(internal.Warehouse{}, internalError)
		//		return internal.Section{}, internalError
		//	},
		//},
		//{
		//	Name: "GIVEN a non valid section, WHEN product_type does not exist, RETURNS error utils.ErrInvalidArguments",
		//	Data: internal.Section{SectionNumber: 1, WarehouseID: 1, ProductTypeID: 1},
		//	Mock: func(repo *MockSectionRepository, warehouseService *MockSectionWarehouseService, productTypeService *MockSectionProductTypeService) (internal.Section, error) {
		//		warehouseService.On("GetByID", mock.Anything).Return(mockWarehouse, nil)
		//		productTypeService.On("GetProductTypeByID", mock.Anything).Return(internal.ProductType{}, utils.ErrNotFound)
		//		return internal.Section{}, utils.EDependencyNotFound("product_type", "id: 1")
		//	},
		//},
		//{
		//	Name: "GIVEN a non valid section, WHEN fail to call productTypeService.GetProductTypeByID, RETURNS internal error",
		//	Data: internal.Section{SectionNumber: 1, WarehouseID: 1, ProductTypeID: 1},
		//	Mock: func(repo *MockSectionRepository, warehouseService *MockSectionWarehouseService, productTypeService *MockSectionProductTypeService) (internal.Section, error) {
		//		warehouseService.On("GetByID", mock.Anything).Return(mockWarehouse, nil)
		//		productTypeService.On("GetProductTypeByID", mock.Anything).Return(internal.ProductType{}, internalError)
		//		return internal.Section{}, internalError
		//	},
		//},
		//{
		//	Name: "GIVEN a non valid section, WHEN minimum_capacity is greater than maximum_capacity, RETURNS error utils.ErrInvalidArguments",
		//	Data: internal.Section{SectionNumber: 1, WarehouseID: 1, ProductTypeID: 1, MinimumCapacity: 5, MaximumCapacity: 4},
		//	Mock: func(repo *MockSectionRepository, warehouseService *MockSectionWarehouseService, productTypeService *MockSectionProductTypeService) (internal.Section, error) {
		//		warehouseService.On("GetByID", mock.Anything).Return(mockWarehouse, nil)
		//		productTypeService.On("GetProductTypeByID", mock.Anything).Return(mockProductType, nil)
		//		return internal.Section{}, utils.EBR("minimum_capacity cannot be greater than maximum_capacity")
		//	},
		//},
		//{
		//	Name: "GIVEN a non valid section, WHEN minimum_temperature is less than -273.15 Celsius, RETURNS error utils.ErrInvalidArguments",
		//	Data: internal.Section{SectionNumber: 1, WarehouseID: 1, ProductTypeID: 1, MinimumCapacity: 3, MaximumCapacity: 4, MinimumTemperature: -300},
		//	Mock: func(repo *MockSectionRepository, warehouseService *MockSectionWarehouseService, productTypeService *MockSectionProductTypeService) (internal.Section, error) {
		//		warehouseService.On("GetByID", mock.Anything).Return(mockWarehouse, nil)
		//		productTypeService.On("GetProductTypeByID", mock.Anything).Return(mockProductType, nil)
		//		return internal.Section{}, utils.EBR("minimum_temperature cannot be less than -273.15 Celsius")
		//	},
		//},
		//{
		//	Name: "GIVEN a non valid section, WHEN current_temperature is less than -273.15 Celsius, RETURNS error utils.ErrInvalidArguments",
		//	Data: internal.Section{SectionNumber: 1, WarehouseID: 1, ProductTypeID: 1, MinimumCapacity: 3, MaximumCapacity: 4, MinimumTemperature: 0, CurrentTemperature: -300},
		//	Mock: func(repo *MockSectionRepository, warehouseService *MockSectionWarehouseService, productTypeService *MockSectionProductTypeService) (internal.Section, error) {
		//		warehouseService.On("GetByID", mock.Anything).Return(mockWarehouse, nil)
		//		productTypeService.On("GetProductTypeByID", mock.Anything).Return(mockProductType, nil)
		//		return internal.Section{}, utils.EBR("current_temperature cannot be less than -273.15 Celsius")
		//	},
		//},
		//{
		//	Name: "GIVEN a non valid section, WHEN section_number already exists, RETURNS error utils.ErrConflict",
		//	Data: mockSection,
		//	Mock: func(repo *MockSectionRepository, warehouseService *MockSectionWarehouseService, productTypeService *MockSectionProductTypeService) (internal.Section, error) {
		//		repo.On("GetBySectionNumber", mock.Anything).Return(internal.Section{ID: 2}, nil)
		//		warehouseService.On("GetByID", mock.Anything).Return(mockWarehouse, nil)
		//		productTypeService.On("GetProductTypeByID", mock.Anything).Return(mockProductType, nil)
		//		return internal.Section{}, utils.EConflict("section", "id: 1")
		//	},
		//},
		//{
		//	Name: "GIVEN a non valid section, WHEN fail to call repo.GetBySectionNumber, RETURNS internal error",
		//	Data: mockSection,
		//	Mock: func(repo *MockSectionRepository, warehouseService *MockSectionWarehouseService, productTypeService *MockSectionProductTypeService) (internal.Section, error) {
		//		repo.On("GetBySectionNumber", mock.Anything).Return(internal.Section{}, internalError)
		//		warehouseService.On("GetByID", mock.Anything).Return(mockWarehouse, nil)
		//		productTypeService.On("GetProductTypeByID", mock.Anything).Return(mockProductType, nil)
		//		return internal.Section{}, internalError
		//	},
		//},
		//{
		//	Name: "GIVEN a valid section, WHEN no error occurs, SAVE successfully",
		//	Data: mockSection,
		//	Mock: func(repo *MockSectionRepository, warehouseService *MockSectionWarehouseService, productTypeService *MockSectionProductTypeService) (internal.Section, error) {
		//		repo.On("Save", mock.Anything).Return(nil)
		//		repo.On("GetBySectionNumber", mock.Anything).Return(internal.Section{}, nil)
		//		warehouseService.On("GetByID", mock.Anything).Return(mockWarehouse, nil)
		//		productTypeService.On("GetProductTypeByID", mock.Anything).Return(mockProductType, nil)
		//		return mockSection, nil
		//	},
		//},
		//{
		//	Name: "GIVEN a valid section, WHEN calling repo.Save, RETURNS internal error",
		//	Data: mockSection,
		//	Mock: func(repo *MockSectionRepository, warehouseService *MockSectionWarehouseService, productTypeService *MockSectionProductTypeService) (internal.Section, error) {
		//		repo.On("Save", mock.Anything).Return(internalError)
		//		repo.On("GetBySectionNumber", mock.Anything).Return(internal.Section{}, nil)
		//		warehouseService.On("GetByID", mock.Anything).Return(mockWarehouse, nil)
		//		productTypeService.On("GetProductTypeByID", mock.Anything).Return(mockProductType, nil)
		//		return internal.Section{}, internalError
		//	},
		//},
	}
	for _, scenario := range cases {
		t.Run(scenario.Name, func(m *testing.T) {
			// Create the mocks
			repo := new(MockSectionRepository)
			warehouseService := new(MockSectionWarehouseService)
			productTypeService := new(MockSectionProductTypeService)

			// Mock and get the expected
			expectedData, expectedError := scenario.Mock(repo, warehouseService, productTypeService)

			service := NewBasicSectionService(repo, warehouseService, productTypeService)
			savedSection, err := service.Update(scenario.DataID, scenario.Data)
			require.Equal(m, expectedData, savedSection)
			if expectedError == nil {
				require.Nil(m, err)
			} else {
				require.Equal(m, err.Error(), expectedError.Error())
			}
		})
	}
}

func TestUnitSection_GetSectionProductsReport(t *testing.T) {
	internalError := errors.New("internal error")

	t.Run("GIVEN a id == 0, WHEN no errors, RETURN successfully", func(s *testing.T) {
		repo := new(MockSectionRepository)
		repo.On("GetSectionProductsReport").Return(mockSectionProductsReport, nil)
		service := NewBasicSectionService(repo, nil, nil)
		report, err := service.GetSectionProductsReport(0)
		require.Equal(t, mockSectionProductsReport, report)
		require.Nil(t, err)
	})
	t.Run("GIVEN a id == 0, WHEN calling repo.GetSectionProductsReport(), RETURN internal error", func(s *testing.T) {
		repo := new(MockSectionRepository)
		repo.On("GetSectionProductsReport").Return([]internal.SectionProductsReport{}, internalError)
		service := NewBasicSectionService(repo, nil, nil)
		report, err := service.GetSectionProductsReport(0)
		require.ErrorIs(t, err, internalError)
		require.Nil(t, report)
	})

	t.Run("GIVEN a id != 0, WHEN calling repo.GetByID(), RETURN internal error", func(s *testing.T) {
		repo := new(MockSectionRepository)
		repo.On("GetByID", mock.Anything).Return(internal.Section{}, internalError)
		service := NewBasicSectionService(repo, nil, nil)
		report, err := service.GetSectionProductsReport(1)
		require.ErrorIs(t, err, internalError)
		require.Nil(t, report)
	})

	t.Run("GIVEN a id != 0, WHEN calling repo.GetByID(), RETURN utils.ErrNotFound", func(s *testing.T) {
		repo := new(MockSectionRepository)
		repo.On("GetByID", mock.Anything).Return(internal.Section{}, utils.ENotFound("section"))
		service := NewBasicSectionService(repo, nil, nil)
		report, err := service.GetSectionProductsReport(1)
		require.ErrorIs(t, err, utils.ErrNotFound)
		require.Nil(t, report)
	})

	t.Run("GIVEN a id != 0, WHEN calling repo.GetSectionProductsReportByID(), RETURN internal error", func(s *testing.T) {
		repo := new(MockSectionRepository)
		repo.On("GetByID", mock.Anything).Return(mockSection, nil)
		repo.On("GetSectionProductsReportByID", mock.Anything).Return([]internal.SectionProductsReport{}, internalError)
		service := NewBasicSectionService(repo, nil, nil)
		report, err := service.GetSectionProductsReport(1)
		require.ErrorIs(t, err, internalError)
		require.Nil(s, report)
	})

	t.Run("GIVEN a id != 0, WHEN calling repo.GetByID(), RETURN utils.ErrNotFound", func(s *testing.T) {
		repo := new(MockSectionRepository)
		repo.On("GetByID", mock.Anything).Return(mockSection, nil)
		repo.On("GetSectionProductsReportByID", mock.Anything).Return(mockSectionProductsReport, nil)
		service := NewBasicSectionService(repo, nil, nil)
		report, err := service.GetSectionProductsReport(1)
		require.NotNil(s, report)
		require.NoError(s, err)
	})

}
