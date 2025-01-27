package purchase_order

import (
	"testing"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockPurchaseOrderRepository struct {
	mock.Mock
}

func (m *mockPurchaseOrderRepository) FindAll() ([]internal.PurchaseOrder, error) {
	args := m.Called()
	return args.Get(0).([]internal.PurchaseOrder), args.Error(1)
}

func (m *mockPurchaseOrderRepository) FindAllByBuyerID(buyerID int) ([]internal.PurchaseOrderSummary, error) {
	args := m.Called(buyerID)
	return args.Get(0).([]internal.PurchaseOrderSummary), args.Error(1)
}

func (m *mockPurchaseOrderRepository) CreatePurchaseOrder(newPurchaseOrder internal.PurchaseOrderAttributes) (purchaseOrder internal.PurchaseOrder, err error) {
	args := m.Called(newPurchaseOrder)
	return args.Get(0).(internal.PurchaseOrder), args.Error(1)
}

type mockPurchaseOrderBuyerValidation struct {
	mock.Mock
}

func (m *mockPurchaseOrderBuyerValidation) GetOne(id int) (*internal.Buyer, error) {
	args := m.Called(id)
	return args.Get(0).(*internal.Buyer), args.Error(1)
}

type mockPurchaseOrderProductRecordValidation struct {
	mock.Mock
}

func (m *mockPurchaseOrderProductRecordValidation) FindByID(productRecordID int) (internal.ProductRecords, error) {
	args := m.Called(productRecordID)
	return args.Get(0).(internal.ProductRecords), args.Error(1)
}

var (
	mockJsonPurchaseOrder = `{
		"order_number": "order#101",
		"order_date": "2021-04-04",
		"tracking_code": "abscf1234",
		"buyer_id": 1,
		"product_record_id": 1
	}`
	mockNewPurchaseOrder = internal.PurchaseOrderAttributes{
		OrderNumber:     "order#101",
		OrderDate:       "2021-04-04",
		TrackingCode:    "abscf1234",
		BuyerID:         1,
		ProductRecordID: 1,
	}
	mockInvalidNewPurchaseOrder = internal.PurchaseOrderAttributes{
		OrderNumber:     "order#101",
		OrderDate:       "2021-04-04",
		TrackingCode:    "abscf1234",
		BuyerID:         1,
		ProductRecordID: 99,
	}
	mockPurchaseOrder = internal.PurchaseOrder{
		ID: 1,
		Attributes: internal.PurchaseOrderAttributes{
			OrderNumber:     "order#101",
			OrderDate:       "2021-04-04",
			TrackingCode:    "abscf1234",
			BuyerID:         1,
			ProductRecordID: 1,
		}}
	mockPurchaseOrder2 = internal.PurchaseOrder{
		ID: 2,
		Attributes: internal.PurchaseOrderAttributes{
			OrderNumber:     "order#10101",
			OrderDate:       "2021-04-04",
			TrackingCode:    "xyz123456",
			BuyerID:         1,
			ProductRecordID: 1,
		}}
	mockPurchaseOrderSummary = internal.PurchaseOrderSummary{
		OrderCodes:  "order#1",
		TotalOrders: 1,
		BuyerID:     1,
	}
	mockBuyer = internal.Buyer{
		ID: 1,
		BuyerAttributes: internal.BuyerAttributes{
			CardNumberID: "1",
			FirstName:    "Manon",
			LastName:     "Blackbeak",
		},
	}

	mockProductRecord = internal.ProductRecords{
		ID:             1,
		LastUpdateDate: "2025-01-01",
		PurchasePrice:  10.50,
		SalePrice:      15.00,
		ProductID:      1,
	}
)

func TestPurchaseOrdersService_FindAll(t *testing.T) {
	mockRepo := new(mockPurchaseOrderRepository)
	mockBV := new(mockPurchaseOrderBuyerValidation)
	mockPRV := new(mockPurchaseOrderProductRecordValidation)
	service := NewPurchaseOrderService(mockRepo, mockBV, mockPRV)
	t.Run("FindAllByBuyerID - Valid ID", func(t *testing.T) {
		mockRepo.On("FindAllByBuyerID", 1).Return([]internal.PurchaseOrderSummary{mockPurchaseOrderSummary}, nil)
		result, err := service.FindAllByBuyerID(1)

		assert.Equal(t, []internal.PurchaseOrderSummary{mockPurchaseOrderSummary}, result)
		assert.Nil(t, err)
	})

	t.Run("FindAllByBuyerID - Invalid ID", func(t *testing.T) {
		mockRepo.On("FindAllByBuyerID", 99).Return([]internal.PurchaseOrderSummary{}, utils.ErrNotFound)
		result, err := service.FindAllByBuyerID(99)

		assert.Equal(t, []internal.PurchaseOrderSummary(nil), result)
		assert.Equal(t, utils.ErrNotFound, err)
	})

	t.Run("FindAll - Success", func(t *testing.T) {
		mockRepo.On("FindAllByBuyerID", 0).Return([]internal.PurchaseOrderSummary{mockPurchaseOrderSummary}, nil)
		result, err := service.FindAllByBuyerID(0)

		assert.Equal(t, []internal.PurchaseOrderSummary{mockPurchaseOrderSummary}, result)
		assert.Nil(t, err)
	})
}

func TestPurchaseOrdersService_Create(t *testing.T) {
	t.Run("Create - Success", func(t *testing.T) {
		mockRepo := new(mockPurchaseOrderRepository)
		mockBV := new(mockPurchaseOrderBuyerValidation)
		mockPRV := new(mockPurchaseOrderProductRecordValidation)
		service := NewPurchaseOrderService(mockRepo, mockBV, mockPRV)

		mockBV.On("GetOne", 1).Return(&mockBuyer, nil)
		mockPRV.On("FindByID", 1).Return(mockProductRecord, nil)
		mockRepo.On("FindAll").Return([]internal.PurchaseOrder{mockPurchaseOrder2}, nil)
		mockRepo.On("CreatePurchaseOrder", mockNewPurchaseOrder).Return(mockPurchaseOrder, nil)

		result, err := service.CreatePurchaseOrder(mockNewPurchaseOrder)

		assert.Equal(t, mockPurchaseOrder, result)
		assert.Nil(t, err)
	})

	t.Run("Create - Conflict", func(t *testing.T) {
		mockRepo := new(mockPurchaseOrderRepository)
		mockBV := new(mockPurchaseOrderBuyerValidation)
		mockPRV := new(mockPurchaseOrderProductRecordValidation)
		service := NewPurchaseOrderService(mockRepo, mockBV, mockPRV)

		mockBV.On("GetOne", 1).Return(&mockBuyer, nil)
		mockPRV.On("FindByID", 1).Return(mockProductRecord, nil)
		mockRepo.On("FindAll").Return([]internal.PurchaseOrder{}, utils.ErrConflict)
		mockRepo.On("CreatePurchaseOrder", mockNewPurchaseOrder).Return(internal.PurchaseOrder{}, utils.ErrConflict)

		result, err := service.CreatePurchaseOrder(mockNewPurchaseOrder)

		assert.Equal(t, internal.PurchaseOrder{}, result)
		assert.Equal(t, utils.ErrConflict, err)
	})

	t.Run("Create - Empty Arguments", func(t *testing.T) {
		mockRepo := new(mockPurchaseOrderRepository)
		mockBV := new(mockPurchaseOrderBuyerValidation)
		mockPRV := new(mockPurchaseOrderProductRecordValidation)
		service := NewPurchaseOrderService(mockRepo, mockBV, mockPRV)

		mockBV.On("GetOne", 1).Return(&mockBuyer, nil)
		mockPRV.On("FindByID", 1).Return(mockProductRecord, nil)
		mockRepo.On("FindAll").Return([]internal.PurchaseOrder{mockPurchaseOrder2}, nil)
		mockRepo.On("CreatePurchaseOrder", internal.PurchaseOrderAttributes{}).Return(internal.PurchaseOrder{}, nil)

		result, err := service.CreatePurchaseOrder(internal.PurchaseOrderAttributes{})

		assert.Equal(t, internal.PurchaseOrder{}, result)
		assert.Equal(t, utils.ErrEmptyArguments, err)
	})

	t.Run("Create - Product Record Invalid", func(t *testing.T) {
		mockRepo := new(mockPurchaseOrderRepository)
		mockBV := new(mockPurchaseOrderBuyerValidation)
		mockPRV := new(mockPurchaseOrderProductRecordValidation)
		service := NewPurchaseOrderService(mockRepo, mockBV, mockPRV)

		mockBV.On("GetOne", 1).Return(&mockBuyer, nil)
		mockPRV.On("FindByID", 99).Return(internal.ProductRecords{}, utils.ErrNotFound)
		mockRepo.On("FindAll").Return([]internal.PurchaseOrder{mockPurchaseOrder2}, nil)
		mockRepo.On("CreatePurchaseOrder", mockInvalidNewPurchaseOrder).Return(internal.PurchaseOrder{}, utils.ErrNotFound)

		result, err := service.CreatePurchaseOrder(mockInvalidNewPurchaseOrder)

		assert.Equal(t, internal.PurchaseOrder{}, result)
		assert.Equal(t, utils.EDependencyNotFound("product", "id: "+"99"), err)
	})
}
