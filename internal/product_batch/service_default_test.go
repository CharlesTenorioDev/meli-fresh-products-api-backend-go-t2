package product_batch

import (
	"errors"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"

	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockProductBatchRepository struct {
	mock.Mock
}

func (mpb *MockProductBatchRepository) Save(newBatch *internal.ProductBatchRequest) (internal.ProductBatch, error) {
	args := mpb.Called(newBatch)
	return args.Get(0).(internal.ProductBatch), args.Error(1)
}

func (mpb *MockProductBatchRepository) GetBatchNumber(batchNumber int) (int, error) {
	args := mpb.Called(batchNumber)
	return args.Int(0), args.Error(1)
}

// Mock of ProductRepository
type MockProductRepository struct {
	mock.Mock
}

func (mp *MockProductRepository) GetAll() (listProducts []internal.Product, err error) {
	args := mp.Called()
	return args.Get(0).([]internal.Product), args.Error(1)
}

func (mp *MockProductRepository) Create(newproduct internal.ProductAttributes) (product internal.Product, err error) {
	args := mp.Called(newproduct)
	return args.Get(0).(internal.Product), args.Error(1)
}

func (mp *MockProductRepository) Update(inputProduct internal.Product) (product internal.Product, err error) {
	args := mp.Called(inputProduct)
	return args.Get(0).(internal.Product), args.Error(1)
}

func (mp *MockProductRepository) Delete(id int) (err error) {
	args := mp.Called(id)
	return args.Error(0)
}

func (mp *MockProductRepository) GetByID(id int) (internal.Product, error) {
	args := mp.Called(id)
	return args.Get(0).(internal.Product), args.Error(1)
}

// Mock of SectionRepository
type MockSectionRepository struct {
	mock.Mock
}

func (ms *MockSectionRepository) GetAll() ([]internal.Section, error) {
	args := ms.Called()
	return args.Get(0).([]internal.Section), args.Error(1)
}

func (ms *MockSectionRepository) Save(section *internal.Section) error {
	args := ms.Called(section)
	return args.Error(0)
}

func (ms *MockSectionRepository) Update(section *internal.Section) error {
	args := ms.Called(section)
	return args.Error(0)
}

func (ms *MockSectionRepository) GetBySectionNumber(id int) (internal.Section, error) {
	args := ms.Called(id)
	return args.Get(0).(internal.Section), args.Error(1)
}

func (ms *MockSectionRepository) Delete(id int) error {
	args := ms.Called(id)
	return args.Error(0)
}

func (ms *MockSectionRepository) GetSectionProductsReport() ([]internal.SectionProductsReport, error) {
	args := ms.Called()
	return args.Get(0).([]internal.SectionProductsReport), args.Error(1)
}

func (ms *MockSectionRepository) GetSectionProductsReportByID(id int) ([]internal.SectionProductsReport, error) {
	args := ms.Called(id)
	return args.Get(0).([]internal.SectionProductsReport), args.Error(1)
}

func (ms *MockSectionRepository) GetByID(id int) (internal.Section, error) {
	args := ms.Called(id)
	return args.Get(0).(internal.Section), args.Error(1)
}

func TestUnitProductBatch_Save_Success(t *testing.T) {
	newBatch := internal.ProductBatchRequest{
		BatchNumber:        100,
		CurrentQuantity:    50,
		CurrentTemperature: 22.4,
		DueDate:            "2022-01-01",
		InitialQuantity:    10,
		ManufacturingDate:  "2022-01-01",
		ManufacturingHour:  18,
		MinimumTemperature: -3,
		ProductID:          1,
		SectionID:          1,
	}

	batchCreated := internal.ProductBatch{
		ID:                  1,
		ProductBatchRequest: newBatch,
	}

	batchRepo := new(MockProductBatchRepository)
	productRepo := new(MockProductRepository)
	sectionRepo := new(MockSectionRepository)

	batchRepo.On("GetBatchNumber", mock.Anything).Return(0, nil)
	sectionRepo.On("GetByID", newBatch.SectionID).Return(internal.Section{ID: 1}, nil)
	productRepo.On("GetByID", newBatch.ProductID).Return(internal.Product{ID: 1}, nil)
	batchRepo.On("Save", mock.Anything).Return(batchCreated, nil)

	service := NewProductBatchService(batchRepo, productRepo, sectionRepo)

	expectedResult := batchCreated
	result, err := service.Save(&newBatch)

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)

}

func TestUnitProductBatch_Save_BatchNumberAlreadyExists(t *testing.T) {
	newBatch := internal.ProductBatchRequest{
		BatchNumber:        100,
		CurrentQuantity:    50,
		CurrentTemperature: 22.4,
		DueDate:            "2022-01-01",
		InitialQuantity:    10,
		ManufacturingDate:  "2022-01-01",
		ManufacturingHour:  18,
		MinimumTemperature: -3,
		ProductID:          1,
		SectionID:          1,
	}

	batchRepo := new(MockProductBatchRepository)
	productRepo := new(MockProductRepository)
	sectionRepo := new(MockSectionRepository)

	batchRepo.On("GetBatchNumber", mock.Anything).Return(1, nil)

	service := NewProductBatchService(batchRepo, productRepo, sectionRepo)

	_, err := service.Save(&newBatch)

	require.Equal(t, utils.EConflict("batch number", "Product batch"), err)

}

func TestUnitProductBatch_Save_SectionIdDoesNotExist(t *testing.T) {
	newBatch := internal.ProductBatchRequest{
		BatchNumber:        100,
		CurrentQuantity:    50,
		CurrentTemperature: 22.4,
		DueDate:            "2022-01-01",
		InitialQuantity:    10,
		ManufacturingDate:  "2022-01-01",
		ManufacturingHour:  18,
		MinimumTemperature: -3,
		ProductID:          1,
		SectionID:          1,
	}

	batchRepo := new(MockProductBatchRepository)
	productRepo := new(MockProductRepository)
	sectionRepo := new(MockSectionRepository)

	batchRepo.On("GetBatchNumber", mock.Anything).Return(0, nil)
	sectionRepo.On("GetByID", newBatch.SectionID).Return(internal.Section{}, utils.ENotFound("Section ID"))

	service := NewProductBatchService(batchRepo, productRepo, sectionRepo)

	_, err := service.Save(&newBatch)

	require.Equal(t, utils.ENotFound("Section ID"), err)

}

func TestUnitProductBatch_Save_ProductIdDoesNotExist(t *testing.T) {
	newBatch := internal.ProductBatchRequest{
		BatchNumber:        100,
		CurrentQuantity:    50,
		CurrentTemperature: 22.4,
		DueDate:            "2022-01-01",
		InitialQuantity:    10,
		ManufacturingDate:  "2022-01-01",
		ManufacturingHour:  18,
		MinimumTemperature: -3,
		ProductID:          1,
		SectionID:          1,
	}

	batchRepo := new(MockProductBatchRepository)
	productRepo := new(MockProductRepository)
	sectionRepo := new(MockSectionRepository)

	batchRepo.On("GetBatchNumber", mock.Anything).Return(0, nil)
	sectionRepo.On("GetByID", newBatch.SectionID).Return(internal.Section{ID: 1}, nil)
	productRepo.On("GetByID", newBatch.ProductID).Return(internal.Product{}, utils.ENotFound("Product ID"))

	service := NewProductBatchService(batchRepo, productRepo, sectionRepo)

	_, err := service.Save(&newBatch)

	require.Equal(t, utils.ENotFound("Product ID"), err)

}

func TestUnitProductBatch_Save_InvalidOrEmptyCompanyName(t *testing.T) {
	newBatch := internal.ProductBatchRequest{
		BatchNumber:        100,
		CurrentQuantity:    50,
		CurrentTemperature: 22.4,
		DueDate:            "", //empty argument
		InitialQuantity:    10,
		ManufacturingDate:  "2022-01-01",
		ManufacturingHour:  18,
		MinimumTemperature: -3,
		ProductID:          0,
		SectionID:          1,
	}

	batchRepo := new(MockProductBatchRepository)
	productRepo := new(MockProductRepository)
	sectionRepo := new(MockSectionRepository)

	service := NewProductBatchService(batchRepo, productRepo, sectionRepo)

	_, err := service.Save(&newBatch)

	require.Equal(t, utils.EZeroValue("Due date"), err)
}

func TestUnitProductBatch_Save_InvalidOrEmptyBatchNumber(t *testing.T) {
	newBatch := internal.ProductBatchRequest{
		BatchNumber:        0,
		CurrentQuantity:    50,
		CurrentTemperature: 22.4,
		DueDate:            "Company",
		InitialQuantity:    10,
		ManufacturingDate:  "2022-01-01",
		ManufacturingHour:  18,
		MinimumTemperature: -3,
		ProductID:          0,
		SectionID:          1,
	}

	batchRepo := new(MockProductBatchRepository)
	productRepo := new(MockProductRepository)
	sectionRepo := new(MockSectionRepository)

	service := NewProductBatchService(batchRepo, productRepo, sectionRepo)

	_, err := service.Save(&newBatch)

	require.Equal(t, utils.EZeroValue("Batch number"), err)
}

func TestUnitProductBatch_Save_InvalidOrEmptyProductID(t *testing.T) {
	newBatch := internal.ProductBatchRequest{
		BatchNumber:        1,
		CurrentQuantity:    50,
		CurrentTemperature: 22.4,
		DueDate:            "Company",
		InitialQuantity:    10,
		ManufacturingDate:  "2022-01-01",
		ManufacturingHour:  18,
		MinimumTemperature: -3,
		ProductID:          0,
		SectionID:          1,
	}

	batchRepo := new(MockProductBatchRepository)
	productRepo := new(MockProductRepository)
	sectionRepo := new(MockSectionRepository)

	service := NewProductBatchService(batchRepo, productRepo, sectionRepo)

	_, err := service.Save(&newBatch)

	require.Equal(t, utils.EZeroValue("Product ID"), err)
}

func TestUnitProductBatch_Save_InvalidOrEmptySectionID(t *testing.T) {
	newBatch := internal.ProductBatchRequest{
		BatchNumber:        1,
		CurrentQuantity:    50,
		CurrentTemperature: 22.4,
		DueDate:            "Company",
		InitialQuantity:    10,
		ManufacturingDate:  "2022-01-01",
		ManufacturingHour:  18,
		MinimumTemperature: -3,
		ProductID:          1,
		SectionID:          0,
	}

	batchRepo := new(MockProductBatchRepository)
	productRepo := new(MockProductRepository)
	sectionRepo := new(MockSectionRepository)

	service := NewProductBatchService(batchRepo, productRepo, sectionRepo)

	_, err := service.Save(&newBatch)

	require.Equal(t, utils.EZeroValue("Section ID"), err)
}

func TestUnitProductBatch_InternalServerError(t *testing.T) {
	newBatch := internal.ProductBatchRequest{
		BatchNumber:        100,
		CurrentQuantity:    50,
		CurrentTemperature: 22.4,
		DueDate:            "2022-01-01",
		InitialQuantity:    10,
		ManufacturingDate:  "2022-01-01",
		ManufacturingHour:  18,
		MinimumTemperature: -3,
		ProductID:          1,
		SectionID:          1,
	}

	batchCreated := internal.ProductBatch{
		ID:                  1,
		ProductBatchRequest: newBatch,
	}

	batchRepo := new(MockProductBatchRepository)
	productRepo := new(MockProductRepository)
	sectionRepo := new(MockSectionRepository)

	internalErr := errors.New("internal server error")

	batchRepo.On("GetBatchNumber", mock.Anything).Return(0, nil)
	sectionRepo.On("GetByID", newBatch.SectionID).Return(internal.Section{ID: 1}, nil)
	productRepo.On("GetByID", newBatch.ProductID).Return(internal.Product{ID: 1}, nil)
	batchRepo.On("Save", mock.Anything).Return(batchCreated, internalErr)

	service := NewProductBatchService(batchRepo, productRepo, sectionRepo)

	_, err := service.Save(&newBatch)

	require.ErrorIs(t, err, internalErr)

}
