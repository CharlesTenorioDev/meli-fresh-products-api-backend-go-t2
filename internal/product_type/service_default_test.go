package product_type

import (
	"errors"
	"testing"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockProductTypeRepository struct {
	mock.Mock
}

func (m *MockProductTypeRepository) GetAll() ([]internal.ProductType, error) {
	args := m.Called()
	return args.Get(0).([]internal.ProductType), args.Error(1)
}

func (m *MockProductTypeRepository) GetByID(id int) (internal.ProductType, error) {
	args := m.Called(id)
	return args.Get(0).(internal.ProductType), args.Error(1)
}

func (m *MockProductTypeRepository) Create(productType internal.ProductType) (internal.ProductType, error) {
	args := m.Called(productType)
	return args.Get(0).(internal.ProductType), args.Error(1)
}

func (m *MockProductTypeRepository) Update(productType internal.ProductType) (internal.ProductType, error) {
	args := m.Called(productType)
	return args.Get(0).(internal.ProductType), args.Error(1)
}

func (m *MockProductTypeRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestUnitProductType_GetProductTypes(t *testing.T) {
	t.Run("WHEN repository returns no error, RETURN successfully", func(t *testing.T) {
		repo := new(MockProductTypeRepository)
		repo.On("GetAll").Return([]internal.ProductType{internal.ProductType{ID: 1, Description: "Vegetables"}}, nil)
		service := NewProductTypeService(repo)
		productTypes, err := service.GetProductTypes()
		require.NoError(t, err)
		require.Equal(t, 1, len(productTypes))
		require.Equal(t, internal.ProductType{ID: 1, Description: "Vegetables"}, productTypes[0])
	})

	t.Run("WHEN repository returns some error, RETURN the error", func(t *testing.T) {
		repo := new(MockProductTypeRepository)
		repo.On("GetAll").Return([]internal.ProductType{}, errors.New("some error"))
		service := NewProductTypeService(repo)
		productTypes, err := service.GetProductTypes()
		require.Error(t, err)
		require.Equal(t, 0, len(productTypes))
	})
}

func TestUnitProductType_GetProductTypeByID(t *testing.T) {
	t.Run("WHEN repository returns no error, RETURN successfully", func(t *testing.T) {
		repo := new(MockProductTypeRepository)
		repo.On("GetByID", 1).Return(internal.ProductType{ID: 1, Description: "Vegetables"}, nil)
		service := NewProductTypeService(repo)
		productType, err := service.GetProductTypeByID(1)
		require.NoError(t, err)
		require.Equal(t, internal.ProductType{ID: 1, Description: "Vegetables"}, productType)
	})

	t.Run("WHEN repository returns some error, RETURN the error", func(t *testing.T) {
		repo := new(MockProductTypeRepository)
		repo.On("GetByID", 1).Return(internal.ProductType{}, errors.New("some error"))
		service := NewProductTypeService(repo)
		productType, err := service.GetProductTypeByID(1)
		require.Error(t, err)
		require.Equal(t, internal.ProductType{}, productType)
	})
}

func TestUnitProductType_CreateProductType(t *testing.T) {
	t.Run("WHEN repository returns no error, RETURN successfully", func(t *testing.T) {
		repo := new(MockProductTypeRepository)
		repo.On("Create", internal.ProductType{Description: "Vegetables"}).Return(internal.ProductType{ID: 1, Description: "Vegetables"}, nil)
		service := NewProductTypeService(repo)
		productType, err := service.CreateProductType(internal.ProductType{Description: "Vegetables"})
		require.NoError(t, err)
		require.Equal(t, internal.ProductType{ID: 1, Description: "Vegetables"}, productType)
	})

	t.Run("WHEN repository returns some error, RETURN the error", func(t *testing.T) {
		repo := new(MockProductTypeRepository)
		repo.On("Create", internal.ProductType{Description: "Vegetables"}).Return(internal.ProductType{}, errors.New("some error"))
		service := NewProductTypeService(repo)
		productType, err := service.CreateProductType(internal.ProductType{Description: "Vegetables"})
		require.Error(t, err)
		require.Equal(t, internal.ProductType{}, productType)
	})
}

func TestUnitProductType_UpdateProductType(t *testing.T) {
	t.Run("WHEN repository returns no error, RETURN successfully", func(t *testing.T) {
		repo := new(MockProductTypeRepository)
		repo.On("GetByID", 1).Return(internal.ProductType{ID: 1, Description: "Vegetables"}, nil)
		repo.On("Update", internal.ProductType{ID: 1, Description: "Fruits"}).Return(internal.ProductType{ID: 1, Description: "Fruits"}, nil)
		service := NewProductTypeService(repo)
		productType, err := service.UpdateProductType(internal.ProductType{ID: 1, Description: "Fruits"})
		require.NoError(t, err)
		require.Equal(t, internal.ProductType{ID: 1, Description: "Fruits"}, productType)
	})

	t.Run("WHEN repository returns some error, RETURN the error", func(t *testing.T) {
		repo := new(MockProductTypeRepository)
		repo.On("GetByID", 1).Return(internal.ProductType{}, errors.New("some error"))
		service := NewProductTypeService(repo)
		productType, err := service.UpdateProductType(internal.ProductType{ID: 1, Description: "Fruits"})
		require.Error(t, err)
		require.Equal(t, internal.ProductType{}, productType)
	})
}

func TestUnitProductType_DeleteProductType(t *testing.T) {
	t.Run("WHEN repository returns no error, RETURN successfully", func(t *testing.T) {
		repo := new(MockProductTypeRepository)
		repo.On("Delete", 1).Return(nil)
		service := NewProductTypeService(repo)
		err := service.DeleteProductType(1)
		require.NoError(t, err)
	})

	t.Run("WHEN repository returns some error, RETURN the error", func(t *testing.T) {
		repo := new(MockProductTypeRepository)
		repo.On("Delete", 1).Return(errors.New("some error"))
		service := NewProductTypeService(repo)
		err := service.DeleteProductType(1)
		require.Error(t, err)
	})
}
