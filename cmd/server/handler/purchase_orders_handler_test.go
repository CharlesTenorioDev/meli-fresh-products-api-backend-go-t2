package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockPurchaseOrderService struct {
	mock.Mock
}

func (m *mockPurchaseOrderService) FindAll() (map[int]internal.PurchaseOrder, error) {
	args := m.Called()
	return args.Get(0).(map[int]internal.PurchaseOrder), args.Error(1)
}

func (m *mockPurchaseOrderService) FindAllByBuyerID(buyerID int) (PurchaseOrders []internal.PurchaseOrderSummary, err error) {
	args := m.Called(buyerID)
	return args.Get(0).([]internal.PurchaseOrderSummary), args.Error(1)
}

func (m *mockPurchaseOrderService) CreatePurchaseOrder(inputPurchaseOrder internal.PurchaseOrderAttributes) (PurchaseOrder internal.PurchaseOrder, err error) {
	args := m.Called(inputPurchaseOrder)
	return args.Get(0).(internal.PurchaseOrder), args.Error(1)
}

var (
	mockJsonPurchaseOrder = `{
		"order_number": "order#1",
		"order_date": "2021-04-04",
		"tracking_code": "abscf123",
		"buyer_id": 1,
		"product_record_id": 1
	}`
	mockNewPurchaseOrder = internal.PurchaseOrderAttributes{
		OrderNumber:     "order#1",
		OrderDate:       "2021-04-04",
		TrackingCode:    "abscf123",
		BuyerID:         1,
		ProductRecordID: 1,
	}
	mockPurchaseOrder = internal.PurchaseOrder{
		ID: 1,
		Attributes: internal.PurchaseOrderAttributes{
			OrderNumber:     "order#1",
			OrderDate:       "2021-04-04",
			TrackingCode:    "abscf123",
			BuyerID:         1,
			ProductRecordID: 1,
		}}
	mockPurchaseOrderSummary = internal.PurchaseOrderSummary{
		OrderCodes:  "order#1",
		TotalOrders: 1,
		BuyerID:     1,
	}
)

func TestPurchaseOrdersHandler_FindAll(t *testing.T) {
	mockService := new(mockPurchaseOrderService)
	handler := NewPurchaseOrdersHandler(mockService)
	t.Run("FindAllByBuyerID - Valid ID", func(t *testing.T) {
		mockService.On("FindAllByBuyerID", 1).Return([]internal.PurchaseOrderSummary{mockPurchaseOrderSummary}, nil)

		req := httptest.NewRequest("GET", "/buyers/reportPurchaseOrders?id=1", nil)
		res := httptest.NewRecorder()
		handler.GetAllPurchaseOrders()(res, req)

		assert.Equal(t, http.StatusOK, res.Result().StatusCode)
		assert.Equal(t, "application/json", res.Header().Get("Content-Type"))
	})

	t.Run("FindAllByBuyerID - Invalid ID", func(t *testing.T) {
		mockService.On("FindAllByBuyerID", 1).Return(map[int]internal.PurchaseOrderSummary{1: mockPurchaseOrderSummary}, nil)

		req := httptest.NewRequest("GET", "/buyers/reportPurchaseOrders?id=x", nil)
		res := httptest.NewRecorder()
		handler.GetAllPurchaseOrders()(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
		assert.Equal(t, "application/json", res.Header().Get("Content-Type"))
	})

	t.Run("FindAll - Success", func(t *testing.T) {
		mockService.On("FindAllByBuyerID", mock.Anything).Return([]internal.PurchaseOrderSummary{mockPurchaseOrderSummary}, nil)
		mockService.On("FindAll").Return(map[int]internal.PurchaseOrder{1: mockPurchaseOrder}, nil)

		req := httptest.NewRequest("GET", "/buyers/reportPurchaseOrders", nil)
		res := httptest.NewRecorder()
		handler.GetAllPurchaseOrders()(res, req)

		assert.Equal(t, http.StatusOK, res.Result().StatusCode)
		assert.Equal(t, "application/json", res.Header().Get("Content-Type"))
	})
}

func TestPurchaseOrdersHandler_Create(t *testing.T) {
	t.Run("Create - Success", func(t *testing.T) {
		mockService := new(mockPurchaseOrderService)
		handler := NewPurchaseOrdersHandler(mockService)
		mockService.On("CreatePurchaseOrder", mockNewPurchaseOrder).Return(mockPurchaseOrder, nil)

		req := httptest.NewRequest("POST", "/purchaseOrders", bytes.NewBufferString(mockJsonPurchaseOrder))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		handler.PostPurchaseOrders()(res, req)

		assert.Equal(t, http.StatusCreated, res.Result().StatusCode)
		assert.Equal(t, "application/json", res.Header().Get("Content-Type"))
	})

	t.Run("Create - Conflict", func(t *testing.T) {
		mockService := new(mockPurchaseOrderService)
		handler := NewPurchaseOrdersHandler(mockService)
		mockService.On("CreatePurchaseOrder", mockNewPurchaseOrder).Return(internal.PurchaseOrder{}, utils.ErrConflict)

		req := httptest.NewRequest("POST", "/purchaseOrders", bytes.NewBufferString(mockJsonPurchaseOrder))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		handler.PostPurchaseOrders()(res, req)

		assert.Equal(t, http.StatusConflict, res.Result().StatusCode)
		assert.Equal(t, "application/json", res.Header().Get("Content-Type"))
	})

	t.Run("Create - Empty Arguments", func(t *testing.T) {
		mockService := new(mockPurchaseOrderService)
		handler := NewPurchaseOrdersHandler(mockService)
		mockService.On("CreatePurchaseOrder", mockNewPurchaseOrder).Return(internal.PurchaseOrder{}, utils.ErrEmptyArguments)

		req := httptest.NewRequest("POST", "/purchaseOrders", bytes.NewBufferString(mockJsonPurchaseOrder))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		handler.PostPurchaseOrders()(res, req)

		assert.Equal(t, http.StatusUnprocessableEntity, res.Result().StatusCode)
		assert.Equal(t, "application/json", res.Header().Get("Content-Type"))
	})
}
