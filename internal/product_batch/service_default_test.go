package product_batch

import (
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockProductBatchRepository struct {
	mock.Mock
}

//17h tem touchbase

func (m *MockProductBatchRepository) Save(newBatch *internal.ProductBatchRequest) (internal.ProductBatch, error) {
	args := m.Called(newBatch)
	return args.Get(0).(internal.ProductBatch), args.Error(1)
}

func (m *MockProductBatchRepository) GetBatchNumber(batchNumber int) (int, error) {
	args := m.Called(batchNumber)
	return args.Int(0), args.Error(1)
}

func TestUnitSeller_GetAll(t *testing.T) {

}
