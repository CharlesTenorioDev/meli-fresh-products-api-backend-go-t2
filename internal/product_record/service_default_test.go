package product_record_test

import (
	"errors"
	"testing"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/product_record"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockProductRecordsRepository struct {
	mock.Mock
}

func (m *mockProductRecordsRepository) Read(productID int) ([]internal.ProductReport, error) {
	args := m.Called(productID)
	return args.Get(0).([]internal.ProductReport), args.Error(1)
}

func (m *mockProductRecordsRepository) Create(newProductRecord internal.ProductRecords) (internal.ProductRecords, error) {
	args := m.Called(newProductRecord)
	return args.Get(0).(internal.ProductRecords), args.Error(1)
}

type mockProductValidation struct {
	mock.Mock
}

func (m *mockProductValidation) GetProductByID(id int) (internal.Product, error) {
	args := m.Called(id)
	return args.Get(0).(internal.Product), args.Error(1)
}

func TestProductRecordsService_GetProductRecords(t *testing.T) {
	cases := []struct {
		TestName        string
		ProductID       int
		RepoResponse    []internal.ProductReport
		RepoError       error
		ValidationError error
		ExpectedError   error
	}{
		{
			TestName:        "GetProductRecords_RepoError",
			ProductID:       1,
			RepoResponse:    nil,
			RepoError:       errors.New("repository error"),
			ValidationError: nil,
			ExpectedError:   errors.New("repository error"),
		},
		{
			TestName:  "GetProductRecords_Success",
			ProductID: 1,
			RepoResponse: []internal.ProductReport{
				{ProductID: 1, Description: "Product A", RecordsCount: 10},
			},
			RepoError:       nil,
			ValidationError: nil,
			ExpectedError:   nil,
		},
		{
			TestName:        "GetProductRecords_ProductNotFound",
			ProductID:       2,
			RepoResponse:    nil,
			RepoError:       nil,
			ValidationError: utils.ErrNotFound,
			ExpectedError:   utils.ErrNotFound,
		},
		{
			TestName:        "GetProductRecords_NoReportsFound",
			ProductID:       1,
			RepoResponse:    []internal.ProductReport{},
			RepoError:       nil,
			ValidationError: nil,
			ExpectedError:   utils.ErrNotFound,
		},
	}

	for _, c := range cases {
		t.Run(c.TestName, func(t *testing.T) {
			repo := new(mockProductRecordsRepository)
			validation := new(mockProductValidation)

			if c.ValidationError != nil {
				validation.On("GetProductByID", c.ProductID).Return(internal.Product{}, c.ValidationError)
			} else {
				validation.On("GetProductByID", c.ProductID).Return(internal.Product{ID: c.ProductID}, nil)
			}

			repo.On("Read", c.ProductID).Return(c.RepoResponse, c.RepoError)

			service := product_record.NewProductRecordService(repo, validation)
			result, err := service.GetProductRecords(c.ProductID)

			if c.ExpectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, c.RepoResponse, result)
			}
		})
	}
}

func TestProductRecordsService_CreateProductRecord(t *testing.T) {
	cases := []struct {
		TestName        string
		NewProduct      internal.ProductRecords
		RepoResponse    internal.ProductRecords
		RepoError       error
		ValidationError error
		ExpectedError   error
	}{
		{
			TestName: "CreateProductRecord_Success",
			NewProduct: internal.ProductRecords{
				LastUpdateDate: "2025-01-27",
				PurchasePrice:  100.50,
				SalePrice:      150.00,
				ProductID:      1,
			},
			RepoResponse: internal.ProductRecords{
				ID:             1,
				LastUpdateDate: "2025-01-27",
				PurchasePrice:  100.50,
				SalePrice:      150.00,
				ProductID:      1,
			},
			RepoError:       nil,
			ValidationError: nil,
			ExpectedError:   nil,
		},
		{
			TestName: "CreateProductRecord_InvalidArguments",
			NewProduct: internal.ProductRecords{
				LastUpdateDate: "",
				PurchasePrice:  0,
				SalePrice:      0,
				ProductID:      0,
			},
			RepoResponse:    internal.ProductRecords{},
			RepoError:       nil,
			ValidationError: nil,
			ExpectedError:   utils.ErrInvalidArguments,
		},
		{
			TestName: "CreateProductRecord_InvalidArguments",
			NewProduct: internal.ProductRecords{
				LastUpdateDate: "2025-01-27",
				PurchasePrice:  0,
				SalePrice:      0,
				ProductID:      0,
			},
			RepoResponse:    internal.ProductRecords{},
			RepoError:       nil,
			ValidationError: nil,
			ExpectedError:   utils.ErrInvalidArguments,
		},
		{
			TestName: "CreateProductRecord_InvalidArguments",
			NewProduct: internal.ProductRecords{
				LastUpdateDate: "2025-01-27",
				PurchasePrice:  150.00,
				SalePrice:      0,
				ProductID:      0,
			},
			RepoResponse:    internal.ProductRecords{},
			RepoError:       nil,
			ValidationError: nil,
			ExpectedError:   utils.ErrInvalidArguments,
		},
		{
			TestName: "CreateProductRecord_InvalidArguments",
			NewProduct: internal.ProductRecords{
				LastUpdateDate: "2025-01-27",
				PurchasePrice:  100.50,
				SalePrice:      120.50,
				ProductID:      0,
			},
			RepoResponse:    internal.ProductRecords{},
			RepoError:       nil,
			ValidationError: nil,
			ExpectedError:   utils.ErrInvalidArguments,
		},
		{
			TestName: "CreateProductRecord_Conflict",
			NewProduct: internal.ProductRecords{
				LastUpdateDate: "2025-01-27",
				PurchasePrice:  100.50,
				SalePrice:      150.00,
				ProductID:      1,
			},
			RepoResponse:    internal.ProductRecords{},
			RepoError:       nil,
			ValidationError: utils.ErrConflict,
			ExpectedError:   utils.ErrConflict,
		},
	}

	for _, c := range cases {
		t.Run(c.TestName, func(t *testing.T) {
			repo := new(mockProductRecordsRepository)
			validation := new(mockProductValidation)

			validation.On("GetProductByID", c.NewProduct.ProductID).Return(internal.Product{}, c.ValidationError)
			repo.On("Create", c.NewProduct).Return(c.RepoResponse, c.RepoError)

			service := product_record.NewProductRecordService(repo, validation)
			result, err := service.CreateProductRecord(c.NewProduct)

			if c.ExpectedError != nil {
				assert.ErrorIs(t, err, c.ExpectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, c.RepoResponse, result)
			}
		})
	}
}
