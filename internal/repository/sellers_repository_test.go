package repository

import (
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/stretchr/testify/require"
	"testing"
)

var mockSellerDb = internal.Seller{
	ID:          1,
	Cid:         100,
	CompanyName: "company",
	Address:     "address",
	Telephone:   "11900000000",
}

var mockNewSeller = internal.SellerRequest{
	Cid:         200,
	CompanyName: "company2",
	Address:     "address2",
	Telephone:   "11922222222",
}

var mockCreatedSeller = internal.Seller{
	ID:          2,
	Cid:         200,
	CompanyName: "company2",
	Address:     "address2",
	Telephone:   "11922222222",
}

var mockUpdatedSeller = internal.Seller{
	ID:          1,
	Cid:         888,
	CompanyName: "change company",
	Address:     "change address",
	Telephone:   "11955555555",
}

func Test_GetAll_Sellers_Success(t *testing.T) {
	repo := NewSellerDbRepository(map[int]internal.Seller{
		1: mockSellerDb,
	})

	sellers, _ := repo.GetAll()
	require.Equal(t, 1, len(sellers))
	require.Equal(t, mockSellerDb, sellers[1])
}

func Test_GetAll_Sellers_Not_Found(t *testing.T) {
	repo := NewSellerDbRepository(map[int]internal.Seller{})

	sellers, _ := repo.GetAll()
	require.Equal(t, 0, len(sellers))
}

func Test_GetById_Sellers_Success(t *testing.T) {
	repo := NewSellerDbRepository(map[int]internal.Seller{
		1: mockSellerDb,
	})

	sellers, _ := repo.GetById(1)
	require.Equal(t, mockSellerDb, sellers)
}

func Test_GetById_Sellers_Not_Found(t *testing.T) {
	repo := NewSellerDbRepository(nil)

	sellers, _ := repo.GetById(1)
	require.Empty(t, sellers)
}

func Test_GetByCid_Sellers_Success(t *testing.T) {
	repo := NewSellerDbRepository(map[int]internal.Seller{
		1: mockSellerDb,
	})

	sellers, _ := repo.GetByCid(100)
	require.Equal(t, mockSellerDb, sellers)
}

func Test_GetByCid_Sellers_Not_Found(t *testing.T) {
	repo := NewSellerDbRepository(map[int]internal.Seller{})

	sellers, _ := repo.GetByCid(100)
	require.Empty(t, sellers)
}

func Test_Create_Sellers_Success(t *testing.T) {
	repo := NewSellerDbRepository(map[int]internal.Seller{
		1: mockSellerDb,
	})

	createdSeller, _ := repo.Create(mockNewSeller)
	require.Equal(t, mockCreatedSeller, createdSeller)
}
func Test_Update_Sellers_Success(t *testing.T) {
	repo := NewSellerDbRepository(map[int]internal.Seller{
		1: mockSellerDb,
	})

	updatedSeller, _ := repo.Update(mockUpdatedSeller)
	require.Equal(t, mockUpdatedSeller, updatedSeller)
}

func Test_Delete_Sellers_Success(t *testing.T) {
	repo := NewSellerDbRepository(map[int]internal.Seller{
		1: mockSellerDb,
	})

	deletedSeller, _ := repo.Delete(1)
	require.Equal(t, true, deletedSeller)
}
