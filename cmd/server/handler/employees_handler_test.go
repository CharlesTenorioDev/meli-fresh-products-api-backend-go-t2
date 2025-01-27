package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockEmployeeService struct {
	mock.Mock
}

func (m *mockEmployeeService) FindAll() (map[int]internal.Employee, error) {
	args := m.Called()
	return args.Get(0).(map[int]internal.Employee), args.Error(1)
}

func (m *mockEmployeeService) FindByID(id int) (internal.Employee, error) {
	args := m.Called(id)
	return args.Get(0).(internal.Employee), args.Error(1)
}

func (m *mockEmployeeService) CreateEmployee(inputEmployee internal.EmployeeAttributes) (employee internal.Employee, err error) {
	args := m.Called(inputEmployee)
	return args.Get(0).(internal.Employee), args.Error(1)
}

func (m *mockEmployeeService) UpdateEmployee(newEmployee internal.Employee) (internal.Employee, error) {
	args := m.Called(newEmployee)
	return args.Get(0).(internal.Employee), args.Error(1)
}

func (m *mockEmployeeService) DeleteEmployee(id int) (err error) {
	args := m.Called(id)
	return args.Error(0)
}

type mockWarehouseValidation struct {
	mock.Mock
}

func (m *mockWarehouseValidation) GetByID(id int) (internal.Warehouse, error) {
	args := m.Called(id)
	return args.Get(0).(internal.Warehouse), args.Error(1)
}

var (
	mockEmployee = internal.Employee{
		ID: 1,
		Attributes: internal.EmployeeAttributes{
			CardNumberID: "12345",
			FirstName:    "Aelin",
			LastName:     "Galanthynius",
			WarehouseID:  1,
		},
	}
	mockEmployee2 = internal.Employee{
		ID: 2,
		Attributes: internal.EmployeeAttributes{
			CardNumberID: "67890",
			FirstName:    "Rowan",
			LastName:     "Withethorn",
			WarehouseID:  1,
		},
	}
	mockEmployeeAttr = internal.EmployeeAttributes{
		CardNumberID: "67890",
		FirstName:    "Rowan",
		LastName:     "Withethorn",
		WarehouseID:  1,
	}
	mockInputEmployee = internal.Employee{
		ID: 1,
		Attributes: internal.EmployeeAttributes{
			CardNumberID: "12345",
			FirstName:    "Celaena",
			LastName:     "Sardothien",
			WarehouseID:  1,
		},
	}
	mockInputEmployeeInvalidID = internal.Employee{
		ID: 99,
		Attributes: internal.EmployeeAttributes{
			CardNumberID: "12345",
			FirstName:    "Celaena",
			LastName:     "Sardothien",
			WarehouseID:  1,
		},
	}
	mockJsonEmployee = `{
		"attributes": {
		  "card_number_id": "123",
		  "first_name": "Celaena",
		  "last_name": "Sardothien",
		  "warehouse_id": 1
		}
	  }`
	mockUpdatedEmployee = internal.Employee{
		ID: 1,
		Attributes: internal.EmployeeAttributes{
			CardNumberID: "123",
			FirstName:    "Celaena",
			LastName:     "Sardothien",
			WarehouseID:  1,
		},
	}
	mockJsonNewEmployee = `{
		  "card_number_id": "123",
		  "first_name": "Celaena",
		  "last_name": "Sardothien",
		  "warehouse_id": 1
	  }`
	mockNewEmployee = internal.EmployeeAttributes{
		CardNumberID: "123",
		FirstName:    "Celaena",
		LastName:     "Sardothien",
		WarehouseID:  1,
	}
	mockWarehouse = internal.Warehouse{
		ID:                 1,
		Address:            "Terrasen 452",
		Telephone:          "0123456789",
		WarehouseCode:      "XYZ",
		MinimumCapacity:    10,
		MinimumTemperature: 10,
	}
)

func TestEmployeeHandler_FindAll(t *testing.T) {
	mockService := new(mockEmployeeService)
	handler := NewEmployeeHandler(mockService)
	t.Run("FindAll - Success", func(t *testing.T) {
		mockService.On("FindAll").Return(map[int]internal.Employee{}, nil)

		req := httptest.NewRequest("GET", "/employees", nil)
		res := httptest.NewRecorder()
		handler.GetAllEmployees()(res, req)

		assert.Equal(t, http.StatusOK, res.Result().StatusCode)
		assert.Equal(t, "application/json", res.Header().Get("Content-Type"))
	})

	t.Run("FindAll - Internal Error", func(t *testing.T) {
		mockService.On("FindAll").Return(map[int]internal.Employee{}, assert.AnError)

		req := httptest.NewRequest("GET", "/employees", nil)
		res := httptest.NewRecorder()
		handler.GetAllEmployees()(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Result().StatusCode)
		assert.Equal(t, "application/json", res.Header().Get("Content-Type"))
	})
}

func TestEmployeeHandler_FindByID(t *testing.T) {
	mockService := new(mockEmployeeService)
	handler := NewEmployeeHandler(mockService)

	t.Run("FindByID - Valid ID", func(t *testing.T) {
		mockService.On("FindByID", 1).Return(mockEmployee, nil)

		req := httptest.NewRequest("GET", "/employees/1", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		res := httptest.NewRecorder()
		handler.GetEmployeesByID()(res, req)

		assert.Equal(t, http.StatusOK, res.Result().StatusCode)
		assert.Equal(t, "application/json", res.Header().Get("Content-Type"))
	})

	t.Run("FindByID - Invalid ID", func(t *testing.T) {
		mockService.On("FindByID", 99).Return(internal.Employee{}, utils.ErrNotFound)

		req := httptest.NewRequest("GET", "/employees/99", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "99")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		res := httptest.NewRecorder()
		handler.GetEmployeesByID()(res, req)

		assert.Equal(t, http.StatusNotFound, res.Result().StatusCode)
		assert.Equal(t, "application/json", res.Header().Get("Content-Type"))
	})

	t.Run("FindByID - Invalid Param Format", func(t *testing.T) {
		mockService.On("FindByID").Return(internal.Employee{}, utils.ErrInvalidFormat)

		req := httptest.NewRequest("GET", "/employees/x", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "x")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		res := httptest.NewRecorder()
		handler.GetEmployeesByID()(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
		assert.Equal(t, "application/json", res.Header().Get("Content-Type"))
	})
}

func TestEmployeeHandler_Delete(t *testing.T) {
	mockService := new(mockEmployeeService)
	handler := NewEmployeeHandler(mockService)

	t.Run("Delete - Valid ID", func(t *testing.T) {
		mockService.On("DeleteEmployee", 1).Return(nil)

		req := httptest.NewRequest("DELETE", "/employees/1", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		res := httptest.NewRecorder()
		handler.DeleteEmployees()(res, req)

		assert.Equal(t, http.StatusNoContent, res.Result().StatusCode)
		assert.Equal(t, "application/json", res.Header().Get("Content-Type"))
	})

	t.Run("Delete - Invalid ID", func(t *testing.T) {
		mockService.On("DeleteEmployee", 99).Return(utils.ErrNotFound)

		req := httptest.NewRequest("DELETE", "/employees/99", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "99")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		res := httptest.NewRecorder()
		handler.DeleteEmployees()(res, req)

		assert.Equal(t, http.StatusNotFound, res.Result().StatusCode)
		assert.Equal(t, "application/json", res.Header().Get("Content-Type"))
	})

	t.Run("Delete - Invalid Param Format", func(t *testing.T) {
		mockService.On("DeleteEmployee").Return(utils.ErrInvalidFormat)

		req := httptest.NewRequest("DELETE", "/employees/x", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "x")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		res := httptest.NewRecorder()
		handler.DeleteEmployees()(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
		assert.Equal(t, "application/json", res.Header().Get("Content-Type"))
	})
}

func TestEmployeeHandler_Update(t *testing.T) {
	mockService := new(mockEmployeeService)
	handler := NewEmployeeHandler(mockService)

	t.Run("Update - Valid ID", func(t *testing.T) {
		mockService.On("UpdateEmployee", mockUpdatedEmployee).Return(mockUpdatedEmployee, nil)

		req := httptest.NewRequest("PATCH", "/employees/1", bytes.NewBufferString(mockJsonEmployee))
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		res := httptest.NewRecorder()
		handler.PatchEmployees()(res, req)

		assert.Equal(t, http.StatusOK, res.Result().StatusCode)
		assert.Equal(t, "application/json", res.Header().Get("Content-Type"))
	})

	t.Run("Update - Invalid ID", func(t *testing.T) {
		mockService.On("UpdateEmployee", mock.Anything).Return(internal.Employee{}, utils.ErrNotFound)

		req := httptest.NewRequest("PATCH", "/employees/99", bytes.NewBufferString(mockJsonEmployee))
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "99")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		res := httptest.NewRecorder()
		handler.PatchEmployees()(res, req)

		assert.Equal(t, http.StatusNotFound, res.Result().StatusCode)
		assert.Equal(t, "application/json", res.Header().Get("Content-Type"))
	})
}

func TestEmployeeHandler_Create(t *testing.T) {
	mockService := new(mockEmployeeService)
	handler := NewEmployeeHandler(mockService)

	t.Run("Create - Success", func(t *testing.T) {
		mockService.On("CreateEmployee", mockNewEmployee).Return(mockUpdatedEmployee, nil)

		req := httptest.NewRequest("POST", "/employees", bytes.NewBufferString(mockJsonNewEmployee))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		handler.PostEmployees()(res, req)

		assert.Equal(t, http.StatusCreated, res.Result().StatusCode)
		assert.Equal(t, "application/json", res.Header().Get("Content-Type"))
	})

	t.Run("Create - Conflict", func(t *testing.T) {
		mockService.On("CreateEmployee", mockNewEmployee).Return(internal.Employee{}, utils.ErrConflict)

		req := httptest.NewRequest("POST", "/employees", bytes.NewBufferString(mockJsonNewEmployee))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		handler.PostEmployees()(res, req)

		assert.Equal(t, http.StatusConflict, res.Result().StatusCode)
		assert.Equal(t, "application/json", res.Header().Get("Content-Type"))
	})

	t.Run("Create - Empty Arguments", func(t *testing.T) {
		mockService.On("CreateEmployee", mockNewEmployee).Return(internal.Employee{}, utils.ErrEmptyArguments)

		req := httptest.NewRequest("POST", "/employees", bytes.NewBufferString(mockJsonNewEmployee))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		handler.PostEmployees()(res, req)

		assert.Equal(t, http.StatusUnprocessableEntity, res.Result().StatusCode)
		assert.Equal(t, "application/json", res.Header().Get("Content-Type"))
	})
}
