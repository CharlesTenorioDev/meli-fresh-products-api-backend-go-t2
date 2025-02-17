package handler

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockProductTypeService struct {
	mock.Mock
}

func (m *MockProductTypeService) GetProductTypes() ([]internal.ProductType, error) {
	args := m.Called()
	return args.Get(0).([]internal.ProductType), args.Error(1)
}

func (m *MockProductTypeService) GetProductTypeByID(id int) (internal.ProductType, error) {
	args := m.Called(id)
	return args.Get(0).(internal.ProductType), args.Error(1)
}

func (m *MockProductTypeService) CreateProductType(newProductType internal.ProductType) (internal.ProductType, error) {
	args := m.Called(newProductType)
	return args.Get(0).(internal.ProductType), args.Error(1)
}

func (m *MockProductTypeService) UpdateProductType(inputProductType internal.ProductType) (internal.ProductType, error) {
	args := m.Called(inputProductType)
	return args.Get(0).(internal.ProductType), args.Error(1)
}

func (m *MockProductTypeService) DeleteProductType(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

var vegetableProductType = internal.ProductType{
	ID:          1,
	Description: "Vegetables",
}

func TestProductTypeHandler_GetProductTypes(t *testing.T) {
	t.Run("WHEN service returns no error, RETURN successfully", func(t *testing.T) {
		service := new(MockProductTypeService)
		service.On("GetProductTypes").Return([]internal.ProductType{vegetableProductType}, nil)

		handler := NewProductTypeHandler(service)

		request := httptest.NewRequest("GET", "/api/v1/product_types/", nil)
		response := httptest.NewRecorder()

		handler.GetProductTypes(response, request)

		var actualResponse struct {
			Data []internal.ProductType `json:"data"`
		}

		err := json.NewDecoder(response.Body).Decode(&actualResponse)
		expectedBody := []internal.ProductType{vegetableProductType}

		require.NoError(t, err)
		require.Equal(t, 200, response.Code)
		require.Equal(t, expectedBody, actualResponse.Data)

	})

	t.Run("WHEN service returns an error, RETURN the error", func(t *testing.T) {
		service := new(MockProductTypeService)
		service.On("GetProductTypes").Return([]internal.ProductType{}, utils.ErrNotFound)

		handler := NewProductTypeHandler(service)

		request := httptest.NewRequest("GET", "/api/v1/product_types/", nil)
		response := httptest.NewRecorder()

		handler.GetProductTypes(response, request)

		var actualResponse struct {
			Status string `json:"status"`
		}

		err := json.NewDecoder(response.Body).Decode(&actualResponse)

		require.NoError(t, err)
		require.Equal(t, 404, response.Code)
		require.Equal(t, "Not Found", actualResponse.Status)
	})

}
