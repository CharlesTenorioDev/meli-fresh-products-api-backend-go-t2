package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockProductBatchService struct {
	mock.Mock
}

func (m *MockProductBatchService) Save(newBatch *internal.ProductBatchRequest) (internal.ProductBatch, error) {
	args := m.Called(newBatch)
	return args.Get(0).(internal.ProductBatch), args.Error(1)
}

func TestUnitSeller_Create_Success(t *testing.T) {
	productBatch := internal.ProductBatch{
		ID: 1,
		ProductBatchRequest: internal.ProductBatchRequest{
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
		},
	}
	service := new(MockProductBatchService)
	service.On("Save", mock.Anything).Return(productBatch, nil)

	bodyByte, _ := json.Marshal(productBatch)
	bodyReader := bytes.NewReader(bodyByte)
	request := httptest.NewRequest(http.MethodPost, "/api/v1/productBatches", bodyReader)
	writer := httptest.NewRecorder()
	handler := NewProductBatchHandler(service)

	handler.Create()(writer, request)

	response := writer.Result()
	require.Equal(t, http.StatusCreated, response.StatusCode)

	responseBody, _ := io.ReadAll(response.Body)
	expectedResponseBody := `{"data":{"id":1,"batch_number":100,"current_quantity":50,"current_temperature":22.4,"due_date":"2022-01-01","initial_quantity":10,"manufacturing_date":"2022-01-01","manufacturing_hour":18,"minimum_temperature":-3,"product_id":1,"section_id":1}}`
	require.Equal(t, expectedResponseBody, string(responseBody))

}

func TestUnitSeller_Create_BadRequest(t *testing.T) {

	service := new(MockProductBatchService)

	bodyByte, _ := json.Marshal(`{"a":1}`)
	bodyReader := bytes.NewReader(bodyByte)
	request := httptest.NewRequest(http.MethodPost, "/api/v1/productBatches", bodyReader)
	writer := httptest.NewRecorder()
	handler := NewProductBatchHandler(service)

	handler.Create()(writer, request)

	response := writer.Result()
	require.Equal(t, http.StatusBadRequest, response.StatusCode)

	responseBody, _ := io.ReadAll(response.Body)
	expectedResponseBody := `{"status":"Bad Request","message":"invalid format"}`
	require.Equal(t, expectedResponseBody, string(responseBody))

}

func TestUnitSeller_Create_Conflict(t *testing.T) {
	productBatch := internal.ProductBatch{
		ID: 1,
		ProductBatchRequest: internal.ProductBatchRequest{
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
		},
	}

	service := new(MockProductBatchService)
	service.On("Save", mock.Anything).Return(internal.ProductBatch{}, utils.ErrConflict)

	bodyByte, _ := json.Marshal(productBatch)
	bodyReader := bytes.NewReader(bodyByte)
	request := httptest.NewRequest(http.MethodPost, "/api/v1/productBatches", bodyReader)
	writer := httptest.NewRecorder()
	handler := NewProductBatchHandler(service)

	handler.Create()(writer, request)

	response := writer.Result()
	require.Equal(t, http.StatusConflict, response.StatusCode)

	responseBody, _ := io.ReadAll(response.Body)
	expectedResponseBody := `{"status":"Conflict","message":"entity already exists"}`
	require.Equal(t, expectedResponseBody, string(responseBody))

}

func TestUnitSeller_Create_InvalidOrEmptyAguments(t *testing.T) {
	productBatch := internal.ProductBatch{
		ID: 1,
		ProductBatchRequest: internal.ProductBatchRequest{
			BatchNumber:        100,
			CurrentQuantity:    50,
			CurrentTemperature: 22.4,
			DueDate:            "",
			InitialQuantity:    10,
			ManufacturingDate:  "2022-01-01",
			ManufacturingHour:  18,
			MinimumTemperature: -3,
			ProductID:          1,
			SectionID:          1,
		},
	}

	service := new(MockProductBatchService)
	service.On("Save", mock.Anything).Return(internal.ProductBatch{}, utils.ErrInvalidArguments)

	bodyByte, _ := json.Marshal(productBatch)
	bodyReader := bytes.NewReader(bodyByte)
	request := httptest.NewRequest(http.MethodPost, "/api/v1/productBatches", bodyReader)
	writer := httptest.NewRecorder()
	handler := NewProductBatchHandler(service)

	handler.Create()(writer, request)

	response := writer.Result()
	require.Equal(t, http.StatusUnprocessableEntity, response.StatusCode)

	responseBody, _ := io.ReadAll(response.Body)
	expectedResponseBody := `{"status":"Unprocessable Entity","message":"invalid arguments"}`
	require.Equal(t, expectedResponseBody, string(responseBody))

}

func TestUnitSeller_Create_InternalServerError(t *testing.T) {
	productBatch := internal.ProductBatch{
		ID: 1,
		ProductBatchRequest: internal.ProductBatchRequest{
			BatchNumber:        100,
			CurrentQuantity:    50,
			CurrentTemperature: 22.4,
			DueDate:            "",
			InitialQuantity:    10,
			ManufacturingDate:  "2022-01-01",
			ManufacturingHour:  18,
			MinimumTemperature: -3,
			ProductID:          1,
			SectionID:          1,
		},
	}

	service := new(MockProductBatchService)
	service.On("Save", mock.Anything).Return(internal.ProductBatch{}, errors.New("some error"))

	bodyByte, _ := json.Marshal(productBatch)
	bodyReader := bytes.NewReader(bodyByte)
	request := httptest.NewRequest(http.MethodPost, "/api/v1/productBatches", bodyReader)
	writer := httptest.NewRecorder()
	handler := NewProductBatchHandler(service)

	handler.Create()(writer, request)

	response := writer.Result()
	require.Equal(t, http.StatusInternalServerError, response.StatusCode)

	responseBody, _ := io.ReadAll(response.Body)
	expectedResponseBody := `{"status":"Internal Server Error","message":"Some error occurs"}`
	require.Equal(t, expectedResponseBody, string(responseBody))

}
