package handler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-go-t2/cmd/server/handler"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockWarehouseService struct {
	mock.Mock
}

func (m *mockWarehouseService) GetAll() ([]internal.Warehouse, error) {
	args := m.Called()
	return args.Get(0).([]internal.Warehouse), args.Error(1)
}

func (m *mockWarehouseService) GetByID(id int) (internal.Warehouse, error) {
	args := m.Called(id)
	return args.Get(0).(internal.Warehouse), args.Error(1)
}

func (m *mockWarehouseService) Save(warehouse internal.Warehouse) (internal.Warehouse, error) {
	args := m.Called(warehouse)
	return args.Get(0).(internal.Warehouse), args.Error(1)
}

func (m *mockWarehouseService) Update(id int, warehouse internal.WarehousePointers) (internal.Warehouse, error) {
	args := m.Called(id, warehouse)
	return args.Get(0).(internal.Warehouse), args.Error(1)
}

func (m *mockWarehouseService) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestWarehouseHandler_GetAll(t *testing.T) {
	cases := []struct {
		TestName           string
		ErrorToReturn      error
		ExpectedBody       string
		ExpectedStatusCode int
	}{
		{
			TestName:           "GetAll_OK",
			ErrorToReturn:      nil,
			ExpectedBody:       `{"data":[{"id":1,"address":"1234 Cold Storage St, LA","telephone":"555-3456","warehouse_code":"WH001","locality_id":1,"minimum_capacity":30,"minimum_temperature":20},{"id":2,"address":"5678 Cool Goods Ave, Toronto","telephone":"555-7890","warehouse_code":"WH002","locality_id":3,"minimum_capacity":30,"minimum_temperature":15}]}`,
			ExpectedStatusCode: http.StatusOK,
		},
		{
			TestName:           "GetAll_NotFound",
			ErrorToReturn:      utils.ENotFound("Warehouse"),
			ExpectedBody:       `{"message":"entity not found: Warehouse doesn't exist", "status":"Not Found"}`,
			ExpectedStatusCode: http.StatusNotFound,
		},
	}

	for _, c := range cases {
		t.Run(c.TestName, func(t *testing.T) {
			service := new(mockWarehouseService)
			if c.ErrorToReturn == nil {
				service.On("GetAll").Return([]internal.Warehouse{
					{
						ID: 1, Address: "1234 Cold Storage St, LA", Telephone: "555-3456", WarehouseCode: "WH001", LocalityID: 1, MinimumCapacity: 30, MinimumTemperature: 20,
					},
					{
						ID: 2, Address: "5678 Cool Goods Ave, Toronto", Telephone: "555-7890", WarehouseCode: "WH002", LocalityID: 3, MinimumCapacity: 30, MinimumTemperature: 15,
					},
				}, nil)
			} else {
				service.On("GetAll").Return([]internal.Warehouse{}, c.ErrorToReturn)
			}

			h := handler.NewWarehouseHandler(service)
			req := httptest.NewRequest(http.MethodGet, "/warehouses", nil)
			res := httptest.NewRecorder()

			h.GetAll()(res, req)
			require.Equal(t, c.ExpectedStatusCode, res.Result().StatusCode)
			require.JSONEq(t, c.ExpectedBody, res.Body.String())
		})
	}
}

func TestWarehouseHandler_GetByID(t *testing.T) {
	cases := []struct {
		TestName           string
		ID                 string
		ErrorToReturn      error
		ExpectedBody       string
		ExpectedStatusCode int
	}{
		{
			TestName:           "GetByID_OK",
			ID:                 "1",
			ErrorToReturn:      nil,
			ExpectedBody:       `{"data": {"id":1,"address":"1234 Cold Storage St, LA","telephone":"555-3456","warehouse_code":"WH001","locality_id":1,"minimum_capacity":30,"minimum_temperature":20}}`,
			ExpectedStatusCode: http.StatusOK,
		},
		{
			TestName:           "GetByID_NotFound",
			ID:                 "2",
			ErrorToReturn:      utils.ENotFound("Warehouse"),
			ExpectedBody:       `{"status":"Not Found","message":"entity not found: Warehouse doesn't exist"}`,
			ExpectedStatusCode: http.StatusNotFound,
		},
		{
			TestName:           "GetByID_BadRequest",
			ID:                 "abc",
			ErrorToReturn:      utils.EBadRequest("Invalid ID"),
			ExpectedBody:       `{"status":"Bad Request","message":"invalid format: Invalid ID with invalid format"}`,
			ExpectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, c := range cases {
		t.Run(c.TestName, func(t *testing.T) {
			service := new(mockWarehouseService)
			if c.ErrorToReturn == nil {
				service.On("GetByID", 1).Return(internal.Warehouse{
					ID: 1, Address: "1234 Cold Storage St, LA", Telephone: "555-3456", WarehouseCode: "WH001", LocalityID: 1, MinimumCapacity: 30, MinimumTemperature: 20,
				}, nil)
			} else {
				service.On("GetByID", mock.Anything).Return(internal.Warehouse{}, c.ErrorToReturn)
			}

			h := handler.NewWarehouseHandler(service)
			req := httptest.NewRequest(http.MethodGet, "/warehouses/"+c.ID, nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", c.ID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			res := httptest.NewRecorder()

			h.GetByID()(res, req)
			require.Equal(t, c.ExpectedStatusCode, res.Result().StatusCode)
			require.JSONEq(t, c.ExpectedBody, res.Body.String())
		})
	}
}

func TestWarehouseHandler_Post(t *testing.T) {
	cases := []struct {
		TestName           string
		RequestBody        string
		ErrorToReturn      error
		ExpectedBody       string
		ExpectedStatusCode int
	}{
		{
			TestName:           "Post_OK",
			RequestBody:        `{"warehouse_code":"WH001","address":"1234 Cold Storage St, LA","telephone":"555-3456","locality_id":1,"minimum_capacity":30,"minimum_temperature":20}`,
			ErrorToReturn:      nil,
			ExpectedBody:       `{"data": {"id":1,"address":"1234 Cold Storage St, LA","telephone":"555-3456","warehouse_code":"WH001","locality_id":1,"minimum_capacity":30,"minimum_temperature":20}}`,
			ExpectedStatusCode: http.StatusCreated,
		},
		{
			TestName:           "Post_UnprocessableEntity",
			RequestBody:        `{"warehouse_code":"","address":"1234 Warehouse St","telephone":"555-1234","locality_id":1,"minimum_capacity":50,"minimum_temperature":10}`,
			ErrorToReturn:      utils.EZeroValue("Warehouse Code"),
			ExpectedBody:       `{"status":"Unprocessable Entity","message":"invalid arguments: Warehouse Code cannot be empty/null"}`,
			ExpectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			TestName:           "Post_InvalidJson",
			RequestBody:        `{JSON_INVALID}`,
			ErrorToReturn:      errors.New("invalid Json"),
			ExpectedBody:       `{"message":"internal server error", "status":"Internal Server Error"}`,
			ExpectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, c := range cases {
		t.Run(c.TestName, func(t *testing.T) {
			service := new(mockWarehouseService)
			if c.ErrorToReturn == nil {
				service.On("Save", mock.Anything).Return(internal.Warehouse{
					ID: 1, WarehouseCode: "WH001",
					Address:            "1234 Cold Storage St, LA",
					Telephone:          "555-3456",
					LocalityID:         1,
					MinimumCapacity:    30,
					MinimumTemperature: 20,
				}, nil)
			} else {
				service.On("Save", mock.Anything).Return(internal.Warehouse{}, c.ErrorToReturn)
			}

			h := handler.NewWarehouseHandler(service)
			req := httptest.NewRequest(http.MethodPost, "/warehouses", strings.NewReader(c.RequestBody))
			req.Header.Set("Content-Type", "application/json")
			res := httptest.NewRecorder()

			h.Post()(res, req)
			require.Equal(t, c.ExpectedStatusCode, res.Result().StatusCode)
			require.JSONEq(t, c.ExpectedBody, res.Body.String())
		})
	}
}

func TestWarehouseHandler_Update(t *testing.T) {
	cases := []struct {
		TestName           string
		ID                 string
		RequestBody        string
		ErrorToReturn      error
		ExpectedBody       string
		ExpectedStatusCode int
	}{
		{
			TestName:           "Update_OK",
			ID:                 "1",
			RequestBody:        `{"warehouse_code":"WH001","address":"1234 Cold Storage St, LA","telephone":"555-3456","locality_id":1,"minimum_capacity":30,"minimum_temperature":20}`,
			ErrorToReturn:      nil,
			ExpectedBody:       `{"data": {"id":1,"address":"1234 Cold Storage St, LA","telephone":"555-3456","warehouse_code":"WH001","locality_id":1,"minimum_capacity":30,"minimum_temperature":20}}`,
			ExpectedStatusCode: http.StatusOK,
		},
		{
			TestName:           "Update_NotFound",
			ID:                 "2",
			RequestBody:        `{"address":"Updated Address"}`,
			ErrorToReturn:      utils.ENotFound("Warehouse"),
			ExpectedBody:       `{"message":"entity not found: Warehouse doesn't exist", "status":"Not Found"}`,
			ExpectedStatusCode: http.StatusNotFound,
		},
		{
			TestName:           "Update_BadRequest",
			ID:                 "abc",
			RequestBody:        `{"address":"Updated Address"}`,
			ErrorToReturn:      utils.EBadRequest("Invalid ID"),
			ExpectedBody:       `{"message":"invalid format: Invalid ID with invalid format", "status":"Bad Request"}`,
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			TestName:           "Update_InvalidJson",
			ID:                 "1",
			RequestBody:        `{INVALID_JSON}`,
			ErrorToReturn:      errors.New("invalid Json"),
			ExpectedBody:       `{"message":"internal server error", "status":"Internal Server Error"}`,
			ExpectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, c := range cases {
		t.Run(c.TestName, func(t *testing.T) {
			service := new(mockWarehouseService)
			service.On("Update", mock.Anything, mock.Anything).Return(internal.Warehouse{
				ID:                 1,
				WarehouseCode:      "WH001",
				Address:            "1234 Cold Storage St, LA",
				Telephone:          "555-3456",
				LocalityID:         1,
				MinimumCapacity:    30,
				MinimumTemperature: 20,
			}, c.ErrorToReturn)

			h := handler.NewWarehouseHandler(service)
			req := httptest.NewRequest(http.MethodPut, "/warehouses/"+c.ID, strings.NewReader(c.RequestBody))
			req.Header.Set("Content-Type", "application/json")
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", c.ID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			res := httptest.NewRecorder()

			h.Update()(res, req)
			require.Equal(t, c.ExpectedStatusCode, res.Result().StatusCode)
			require.JSONEq(t, c.ExpectedBody, res.Body.String())
		})
	}
}

func TestWarehouseHandler_Delete(t *testing.T) {
	cases := []struct {
		TestName           string
		ID                 string
		ErrorToReturn      error
		ExpectedStatusCode int
	}{
		{
			TestName:           "Delete_OK",
			ID:                 "1",
			ErrorToReturn:      nil,
			ExpectedStatusCode: http.StatusNoContent,
		},
		{
			TestName:           "Delete_NotFound",
			ID:                 "2",
			ErrorToReturn:      utils.ENotFound("Warehouse"),
			ExpectedStatusCode: http.StatusNotFound,
		},
		{
			TestName:           "Delete_BadRequest",
			ID:                 "abc",
			ErrorToReturn:      utils.EBadRequest("Invalid ID"),
			ExpectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, c := range cases {
		t.Run(c.TestName, func(t *testing.T) {
			service := new(mockWarehouseService)
			service.On("Delete", mock.Anything).Return(c.ErrorToReturn)

			h := handler.NewWarehouseHandler(service)
			req := httptest.NewRequest(http.MethodDelete, "/warehouses/"+c.ID, nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", c.ID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			res := httptest.NewRecorder()

			h.Delete()(res, req)
			require.Equal(t, c.ExpectedStatusCode, res.Result().StatusCode)
		})
	}
}
