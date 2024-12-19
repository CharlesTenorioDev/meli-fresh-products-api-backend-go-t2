package service

import (
	"testing"

	"github.com/meli-fresh-products-api-backend-go-t2/internal/pkg"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/repository"
	"github.com/stretchr/testify/require"
)

var mockProduct = pkg.Product{
	ID: 1,
	ProductAttributes: pkg.ProductAttributes{
		ProductCode:                    "123",
		Description:                    "test",
		Width:                          1,
		Height:                         1,
		Length:                         1,
		NetWeight:                      1,
		ExpirationRate:                 1,
		RecommendedFreezingTemperature: 1,
		FreezingRate:                   1,
		ProductType:                    1,
		SellerID:                       1,
	},
}

var mockRepo = repository.NewProductDB(map[int]pkg.Product{
	1: mockProduct,
})

var mockProductTypeService = NewProductTypeService(repository.NewProductTypeDB(map[int]pkg.ProductType{
	1: {ID: 1, Description: "test"},
}))

func TestProductService_GetProducts(t *testing.T) {
	s := NewProductService(mockRepo, mockProductTypeService)
	products, _ := s.GetProducts()
	require.Equal(t, 1, len(products))
	require.Equal(t, mockProduct, products[0])
}

func TestProductService_GetProductByID_WhenExists(t *testing.T) {
	s := NewProductService(mockRepo, mockProductTypeService)
	product, _ := s.GetProductByID(1)
	require.Equal(t, mockProduct, product)
}

func TestProductService_GetProductByID_WhenNotExists(t *testing.T) {
	s := NewProductService(mockRepo, mockProductTypeService)
	product, _ := s.GetProductByID(99)
	require.Empty(t, product)
}

func TestProductService_UpdateProduct(t *testing.T) {
	s := NewProductService(mockRepo, mockProductTypeService)
	mockProduct.ProductAttributes.Description = "updated"
	product, _ := s.UpdateProduct(mockProduct)
	require.Equal(t, mockProduct, product)
}

func TestProductService_DeleteProduct(t *testing.T) {
	s := NewProductService(mockRepo, mockProductTypeService)
	err := s.DeleteProduct(1)
	require.Nil(t, err)
}

func TestProductService_DeleteProduct_WhenNotExists(t *testing.T) {
	s := NewProductService(mockRepo, mockProductTypeService)
	err := s.DeleteProduct(99)
	require.NotNil(t, err)
}

func TestProductService_CreateProduct_WhenEmptyFields(t *testing.T) {
	s := NewProductService(mockRepo, mockProductTypeService)
	product, err := s.CreateProduct(pkg.ProductAttributes{})
	require.Empty(t, product)
	require.NotNil(t, err)
}

func TestProductService_CreateProduct_WhenDuplicated(t *testing.T) {
	s := NewProductService(mockRepo, mockProductTypeService)
	product, err := s.CreateProduct(mockProduct.ProductAttributes)
	require.Empty(t, product)
	require.NotNil(t, err)
}

func TestProductService_CreateProduct_WhenEmptyFieldsAndDuplicated(t *testing.T) {
	s := NewProductService(mockRepo, mockProductTypeService)
	product, err := s.CreateProduct(pkg.ProductAttributes{})
	require.Empty(t, product)
	require.NotNil(t, err)
}

func TestProductService_UpdateProduct_WhenNotExists(t *testing.T) {
	s := NewProductService(mockRepo, mockProductTypeService)
	product, err := s.UpdateProduct(pkg.Product{})
	require.Empty(t, product)
	require.NotNil(t, err)
}

func TestProductService_UpdateProduct_WhenEmptyFields(t *testing.T) {
	s := NewProductService(mockRepo, mockProductTypeService)
	product, err := s.UpdateProduct(pkg.Product{ID: 1, ProductAttributes: pkg.ProductAttributes{Description: "updated"}})
	mockProduct.Description = "updated"
	require.Equal(t, mockProduct, product)
	require.Nil(t, err)
}

func Test_prepareProductUpdate(t *testing.T) {
	internalProduct, _ := mockRepo.GetByID(1)
	preparedProduc := prepareProductUpdate(pkg.Product{ID: 1}, internalProduct)
	require.Equal(t, internalProduct, preparedProduc)
}

func Test_validateEmptyFields(t *testing.T) {
	s := NewProductService(mockRepo, mockProductTypeService)

	err := s.validateEmptyFields(pkg.ProductAttributes{})
	require.NotNil(t, err)

	err = s.validateEmptyFields(pkg.ProductAttributes{
		ProductCode: "123",
	})
	require.NotNil(t, err)

	err = s.validateEmptyFields(pkg.ProductAttributes{
		ProductCode: "123",
		Description: "test",
	})
	require.NotNil(t, err)

	err = s.validateEmptyFields(pkg.ProductAttributes{
		ProductCode: "123",
		Description: "test",
		Width:       1,
	})
	require.NotNil(t, err)

	err = s.validateEmptyFields(pkg.ProductAttributes{
		ProductCode: "123",
		Description: "test",
		Width:       1,
		Height:      1,
	})
	require.NotNil(t, err)

	err = s.validateEmptyFields(pkg.ProductAttributes{
		ProductCode: "123",
		Description: "test",
		Width:       1,
		Height:      1,
		Length:      1,
	})
	require.NotNil(t, err)

	err = s.validateEmptyFields(pkg.ProductAttributes{
		ProductCode: "123",
		Description: "test",
		Width:       1,
		Height:      1,
		Length:      1,
		NetWeight:   1,
	})
	require.NotNil(t, err)

	err = s.validateEmptyFields(pkg.ProductAttributes{
		ProductCode:    "123",
		Description:    "test",
		Width:          1,
		Height:         1,
		Length:         1,
		NetWeight:      1,
		ExpirationRate: 1,
	})
	require.NotNil(t, err)

	err = s.validateEmptyFields(pkg.ProductAttributes{
		ProductCode:                    "123",
		Description:                    "test",
		Width:                          1,
		Height:                         1,
		Length:                         1,
		NetWeight:                      1,
		ExpirationRate:                 1,
		RecommendedFreezingTemperature: 1,
	})
	require.NotNil(t, err)
}
