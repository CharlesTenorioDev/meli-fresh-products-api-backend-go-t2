package repository

import (
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"testing"

	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
	"github.com/stretchr/testify/require"
)

var mockProduct = internal.Product{
	ID: 1,
	ProductAttributes: internal.ProductAttributes{
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

func TestProductDB_GetAll(t *testing.T) {
	repo := NewProductDB(map[int]internal.Product{
		1: mockProduct,
	})
	products, _ := repo.GetAll()
	require.Equal(t, 1, len(products))
	require.Equal(t, mockProduct, products[0])
}

func TestProductDB_GetByID_WhenExists(t *testing.T) {
	repo := NewProductDB(map[int]internal.Product{
		1: mockProduct,
	})
	product, _ := repo.GetByID(1)
	require.Equal(t, mockProduct, product)
}

func TestProductDB_GetByID_WhenNotExists(t *testing.T) {
	repo := NewProductDB(nil)
	product, _ := repo.GetByID(1)
	require.Empty(t, product)
}

func TestProductDB_Create(t *testing.T) {
	repo := NewProductDB(nil)
	product, _ := repo.Create(mockProduct.ProductAttributes)
	require.Equal(t, mockProduct.ProductAttributes, product.ProductAttributes)
}

func TestProductDB_Update(t *testing.T) {
	repo := NewProductDB(map[int]internal.Product{
		1: mockProduct,
	})
	mockProduct.ProductAttributes.Description = "updated"
	product, _ := repo.Update(mockProduct)
	require.Equal(t, mockProduct, product)
}

func TestProductDB_Delete_WhenExist(t *testing.T) {
	repo := NewProductDB(map[int]internal.Product{
		1: mockProduct,
	})
	err := repo.Delete(1)
	require.Nil(t, err)
}

func TestProductDB_Delete_WhenNotExist(t *testing.T) {
	repo := NewProductDB(nil)
	err := repo.Delete(1)
	require.ErrorIs(t, err, utils.ErrNotFound)
}
