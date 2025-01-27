package handler

import (
	"bytes"
	"encoding/json"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

type mockProductRecordsService struct {
	mock.Mock
}

func (m *mockProductRecordsService) GetProductRecords(productID int) ([]internal.ProductReport, error) {
	args := m.Called(productID)
	return args.Get(0).([]internal.ProductReport), args.Error(1)
}

func (m *mockProductRecordsService) CreateProductRecord(newProductRecord internal.ProductRecords) (internal.ProductRecords, error) {
	args := m.Called(newProductRecord)
	return args.Get(0).(internal.ProductRecords), args.Error(1)
}

func TestProductRecordsHandler_GetProductRecords(t *testing.T) {
	cases := []struct {
		TestName           string
		IDQueryParam       string
		ServiceResponse    []internal.ProductReport
		ServiceError       error
		ExpectedStatusCode int
		ExpectedBody       string
	}{
		{
			TestName:     "GetProductRecords_OK",
			IDQueryParam: "1",
			ServiceResponse: []internal.ProductReport{
				{ProductID: 1, Description: "Product A", RecordsCount: 10},
			},
			ServiceError:       nil,
			ExpectedStatusCode: http.StatusOK,
			ExpectedBody:       `{"data":[{"product_id":1,"description":"Product A","records_count":10}]}`,
		},
		{
			TestName:           "GetProductRecords_InvalidID",
			IDQueryParam:       "invalid",
			ServiceResponse:    nil,
			ServiceError:       nil,
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedBody:       `{"message":"Invalid 'id' format", "status":"Bad Request"}`,
		},
	}

	for _, c := range cases {
		t.Run(c.TestName, func(t *testing.T) {
			service := new(mockProductRecordsService)
			if c.ServiceError == nil {
				id, _ := strconv.Atoi(c.IDQueryParam)
				service.On("GetProductRecords", id).Return(c.ServiceResponse, nil)
			} else {
				service.On("GetProductRecords", mock.Anything).Return(nil, c.ServiceError)
			}

			h := NewProductRecordsHandler(service)
			req := httptest.NewRequest(http.MethodGet, "/product-records?id="+c.IDQueryParam, nil)
			res := httptest.NewRecorder()

			h.GetProductRecords(res, req)
			require.Equal(t, c.ExpectedStatusCode, res.Result().StatusCode)
			require.JSONEq(t, c.ExpectedBody, res.Body.String())
		})
	}
}

func TestProductRecordsHandler_CreateProductRecord(t *testing.T) {
	cases := []struct {
		TestName           string
		RequestBody        internal.ProductRecords
		ServiceResponse    internal.ProductRecords
		ServiceError       error
		ExpectedStatusCode int
		ExpectedBody       string
	}{
		{
			TestName: "CreateProductRecord_OK",
			RequestBody: internal.ProductRecords{
				LastUpdateDate: "2025-01-27",
				PurchasePrice:  100.50,
				SalePrice:      150.00,
				ProductID:      1,
			},
			ServiceResponse: internal.ProductRecords{
				ID:             1,
				LastUpdateDate: "2025-01-27",
				PurchasePrice:  100.50,
				SalePrice:      150.00,
				ProductID:      1,
			},
			ServiceError:       nil,
			ExpectedStatusCode: http.StatusCreated,
			ExpectedBody:       `{"data":{"id":1,"last_update_date":"2025-01-27","purchase_price":100.50,"sale_price":150.00,"product_id":1}}`,
		},
		{
			TestName: "CreateProductRecord_Conflict",
			RequestBody: internal.ProductRecords{
				LastUpdateDate: "2025-01-27",
				PurchasePrice:  100.50,
				SalePrice:      150.00,
				ProductID:      1,
			},
			ServiceResponse:    internal.ProductRecords{},
			ServiceError:       utils.ErrConflict,
			ExpectedStatusCode: http.StatusConflict,
			ExpectedBody:       `{"message":"entity already exists", "status":"Conflict"}`,
		},
		{
			TestName: "CreateProductRecord_InvalidFormat",
			RequestBody: internal.ProductRecords{
				LastUpdateDate: "",
				PurchasePrice:  0,
				SalePrice:      0,
				ProductID:      0,
			},
			ServiceResponse:    internal.ProductRecords{},
			ServiceError:       utils.ErrInvalidArguments,
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			ExpectedBody:       `{"message":"invalid arguments", "status":"Unprocessable Entity"}`,
		},
	}

	for _, c := range cases {
		t.Run(c.TestName, func(t *testing.T) {
			service := new(mockProductRecordsService)
			service.On("CreateProductRecord", c.RequestBody).Return(c.ServiceResponse, c.ServiceError)

			h := NewProductRecordsHandler(service)

			body, _ := json.Marshal(c.RequestBody)
			req := httptest.NewRequest(http.MethodPost, "/product-records", bytes.NewBuffer(body))
			res := httptest.NewRecorder()

			h.CreateProductRecord(res, req)
			require.Equal(t, c.ExpectedStatusCode, res.Result().StatusCode)
			require.JSONEq(t, c.ExpectedBody, res.Body.String())
		})
	}
}
