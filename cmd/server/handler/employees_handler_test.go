package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-go-t2/cmd/server/handler"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockEmployeesService struct {
	mock.Mock
}

func (m *MockEmployeesService) FindAll() (map[int]internal.Employee, error) {
	args := m.Called()
	return args.Get(0).(map[int]internal.Employee), args.Error(1)
}

func (m *MockEmployeesService) FindByID(id int) (internal.Employee, error) {
	args := m.Called(id)
	return args.Get(0).(internal.Employee), args.Error(1)
}

func (m *MockEmployeesService) CreateEmployee(newEmployee internal.EmployeeAttributes) (internal.Employee, error) {
	args := m.Called(newEmployee)
	return args.Get(0).(internal.Employee), args.Error(1)
}

func (m *MockEmployeesService) UpdateEmployee(inputEmployee internal.Employee) (internal.Employee, error) {
	args := m.Called(inputEmployee)
	return args.Get(0).(internal.Employee), args.Error(1)
}

func (m *MockEmployeesService) DeleteEmployee(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestUnitEmployees_GetAllEmployees(t *testing.T) {
	type testCase struct {
		name               string
		mockEmployees      map[int]internal.Employee
		mockError          error
		expectedStatusCode int
		expectedResponse   map[string]any
	}

	tests := []testCase{
		{
			name: "OK",
			mockEmployees: map[int]internal.Employee{
				1: {
					ID: 1,
					Attributes: internal.EmployeeAttributes{
						CardNumberID: "E001",
						FirstName:    "Alice",
						LastName:     "Johnson",
						WarehouseID:  1,
					},
				},
			},
			mockError:          nil,
			expectedStatusCode: http.StatusOK,
			expectedResponse: map[string]any{
				"message": "success",
				"data": map[string]any{
					"1": map[string]any{
						"id": float64(1),
						"attributes": map[string]any{
							"card_number_id": "E001",
							"first_name":     "Alice",
							"last_name":      "Johnson",
							"warehouse_id":   float64(1),
						},
					},
				},
			},
		},
		{
			name:               "NOT_FOUND",
			mockEmployees:      nil,
			mockError:          utils.ErrNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedResponse: map[string]interface{}{
				"status":  http.StatusText(http.StatusNotFound),
				"message": "No employees found",
			},
		},
		{
			name:               "INTERNAL_SERVER_ERROR",
			mockEmployees:      nil,
			mockError:          errors.New("some internal error"),
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: map[string]interface{}{
				"status":  http.StatusText(http.StatusInternalServerError),
				"message": "An error occurred while retrieving employees",
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			service := new(MockEmployeesService)
			handler := handler.NewEmployeeHandler(service)

			service.On("FindAll").Return(tc.mockEmployees, tc.mockError)

			request := httptest.NewRequest("GET", "/employees", nil)
			request.Header.Set("Content-Type", "application/json")

			response := httptest.NewRecorder()

			handler.GetAllEmployees()(response, request)

			require.Equal(t, tc.expectedStatusCode, response.Code)

			var actualResponse map[string]any
			err := json.Unmarshal(response.Body.Bytes(), &actualResponse)
			require.NoError(t, err)

			require.Equal(t, tc.expectedResponse, actualResponse)

			service.AssertExpectations(t)
		})
	}
}

func TestUnitEmployees_GetEmployeesByID(t *testing.T) {
	type testCase struct {
		name               string
		paramID            string
		mockEmployee       internal.Employee
		mockError          error
		expectedStatusCode int
		expectedResponse   map[string]any
	}

	employeeOK := internal.Employee{
		ID: 1,
		Attributes: internal.EmployeeAttributes{
			CardNumberID: "E001",
			FirstName:    "Alice",
			LastName:     "Johnson",
			WarehouseID:  1,
		},
	}

	tests := []testCase{
		{
			name:               "OK",
			paramID:            "1",
			mockEmployee:       employeeOK,
			mockError:          nil,
			expectedStatusCode: http.StatusOK,
			expectedResponse: map[string]any{
				"message": "success",
				"data": map[string]any{
					"id": float64(employeeOK.ID),
					"attributes": map[string]any{
						"card_number_id": employeeOK.Attributes.CardNumberID,
						"first_name":     employeeOK.Attributes.FirstName,
						"last_name":      employeeOK.Attributes.LastName,
						"warehouse_id":   float64(employeeOK.Attributes.WarehouseID),
					},
				},
			},
		},
		{
			name:               "INVALID ID",
			paramID:            "abc",
			mockEmployee:       internal.Employee{},
			mockError:          nil,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: map[string]any{
				"status":  http.StatusText(http.StatusBadRequest),
				"message": "invalid format",
			},
		},
		{
			name:               "NOT_FOUND",
			paramID:            "2",
			mockEmployee:       internal.Employee{},
			mockError:          utils.ErrNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedResponse: map[string]any{
				"status":  http.StatusText(http.StatusNotFound),
				"message": "entity not found",
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			service := new(MockEmployeesService)
			handler := handler.NewEmployeeHandler(service)

			if tc.name != "INVALID ID" {
				idInt, _ := strconv.Atoi(tc.paramID)
				service.On("FindByID", idInt).Return(tc.mockEmployee, tc.mockError)
			}

			request := httptest.NewRequest("GET", "/employees/"+tc.paramID, nil)
			response := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tc.paramID)

			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

			handler.GetEmployeesByID()(response, request)

			require.Equal(t, tc.expectedStatusCode, response.Code)

			var actualResponse map[string]any
			err := json.Unmarshal(response.Body.Bytes(), &actualResponse)
			require.NoError(t, err)

			require.Equal(t, tc.expectedResponse, actualResponse)

			service.AssertExpectations(t)
		})
	}
}

func TestUnitEmployees_CreateEmployee(t *testing.T) {
	type testCase struct {
		name               string
		requestBody        string
		mockInput          internal.EmployeeAttributes
		mockOutput         internal.Employee
		mockError          error
		expectedStatusCode int
		expectedResponse   map[string]any
		callService        bool
	}

	employeeOKInput := internal.EmployeeAttributes{
		CardNumberID: "E001",
		FirstName:    "Alice",
		LastName:     "Johnson",
		WarehouseID:  1,
	}

	employeeOKOutput := internal.Employee{
		ID: 1,
		Attributes: internal.EmployeeAttributes{
			CardNumberID: "E001",
			FirstName:    "Alice",
			LastName:     "Johnson",
			WarehouseID:  1,
		},
	}

	tests := []testCase{
		{
			name:               "OK",
			requestBody:        `{"card_number_id":"E001","first_name":"Alice","last_name":"Johnson","warehouse_id":1}`,
			mockInput:          employeeOKInput,
			mockOutput:         employeeOKOutput,
			mockError:          nil,
			expectedStatusCode: http.StatusCreated,
			expectedResponse: map[string]any{
				"message": "success",
				"data": map[string]any{
					"id": float64(employeeOKOutput.ID),
					"attributes": map[string]any{
						"card_number_id": employeeOKOutput.Attributes.CardNumberID,
						"first_name":     employeeOKOutput.Attributes.FirstName,
						"last_name":      employeeOKOutput.Attributes.LastName,
						"warehouse_id":   float64(employeeOKOutput.Attributes.WarehouseID),
					},
				},
			},
			callService: true,
		},
		{
			name:               "INVALID BODY",
			requestBody:        `{"card_number_id":"E001"`,
			mockInput:          internal.EmployeeAttributes{},
			mockOutput:         internal.Employee{},
			mockError:          nil,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: map[string]any{
				"status":  http.StatusText(http.StatusBadRequest),
				"message": "invalid format: employee with invalid format",
			},
			callService: false,
		},
		{
			name:               "CONFLICT",
			requestBody:        `{"card_number_id":"E001","first_name":"Alice","last_name":"Johnson","warehouse_id":1}`,
			mockInput:          employeeOKInput,
			mockOutput:         internal.Employee{},
			mockError:          utils.ErrConflict,
			expectedStatusCode: http.StatusConflict,
			expectedResponse: map[string]any{
				"status":  http.StatusText(http.StatusConflict),
				"message": "entity already exists",
			},
			callService: true,
		},
		{
			name:               "INTERNAL_SERVER_ERROR",
			requestBody:        `{"card_number_id":"E001","first_name":"Alice","last_name":"Johnson","warehouse_id":1}`,
			mockInput:          employeeOKInput,
			mockOutput:         internal.Employee{},
			mockError:          errors.New("some internal error"),
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: map[string]any{
				"status":  http.StatusText(http.StatusInternalServerError),
				"message": "internal server error",
			},
			callService: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			service := new(MockEmployeesService)
			empHandler := handler.NewEmployeeHandler(service)

			if tc.callService {
				service.On("CreateEmployee", tc.mockInput).Return(tc.mockOutput, tc.mockError)
			}

			request := httptest.NewRequest("POST", "/employees",

				strings.NewReader(tc.requestBody),
			)
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()

			empHandler.PostEmployees()(response, request)

			require.Equal(t, tc.expectedStatusCode, response.Code)

			var actualResponse map[string]any
			err := json.Unmarshal(response.Body.Bytes(), &actualResponse)
			require.NoError(t, err)

			require.Equal(t, tc.expectedResponse, actualResponse)

			service.AssertExpectations(t)
		})
	}
}

func TestUnitEmployees_UpdateEmployee(t *testing.T) {
	type testCase struct {
		name               string
		paramID            string
		requestBody        string
		mockInputEmployee  internal.Employee
		mockOutputEmployee internal.Employee
		mockError          error
		expectedStatusCode int
		expectedResponse   map[string]any
		callService        bool
	}

	validBody := `{
		"attributes": {
			"card_number_id":"E001",
			"first_name":"AliceUpdated",
			"last_name":"Johnson",
			"warehouse_id":2
		}
	}`

	inputEmployee := internal.Employee{
		ID: 0,
		Attributes: internal.EmployeeAttributes{
			CardNumberID: "E001",
			FirstName:    "AliceUpdated",
			LastName:     "Johnson",
			WarehouseID:  2,
		},
	}

	outputEmployee := internal.Employee{
		ID: 1,
		Attributes: internal.EmployeeAttributes{
			CardNumberID: "E001",
			FirstName:    "AliceUpdated",
			LastName:     "Johnson",
			WarehouseID:  2,
		},
	}

	tests := []testCase{
		{
			name:               "OK",
			paramID:            "1",
			requestBody:        validBody,
			mockInputEmployee:  inputEmployee,
			mockOutputEmployee: outputEmployee,
			mockError:          nil,
			expectedStatusCode: http.StatusOK,
			expectedResponse: map[string]any{
				"message": "success",
				"data": map[string]any{
					"id": float64(outputEmployee.ID),
					"attributes": map[string]any{
						"card_number_id": outputEmployee.Attributes.CardNumberID,
						"first_name":     outputEmployee.Attributes.FirstName,
						"last_name":      outputEmployee.Attributes.LastName,
						"warehouse_id":   float64(outputEmployee.Attributes.WarehouseID),
					},
				},
			},
			callService: true,
		},
		{
			name:               "INVALID ID",
			paramID:            "abc",
			requestBody:        validBody,
			mockInputEmployee:  internal.Employee{},
			mockOutputEmployee: internal.Employee{},
			mockError:          nil,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: map[string]any{
				"status":  http.StatusText(http.StatusBadRequest),
				"message": "invalid format",
			},
			callService: false,
		},
		{
			name:               "INVALID BODY",
			paramID:            "1",
			requestBody:        `{"attributes": { "card_number_id": "E001"`,
			mockInputEmployee:  internal.Employee{},
			mockOutputEmployee: internal.Employee{},
			mockError:          nil,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: map[string]any{
				"status":  http.StatusText(http.StatusBadRequest),
				"message": "invalid format",
			},
			callService: false,
		},
		{
			name:               "NOT_FOUND",
			paramID:            "2",
			requestBody:        validBody,
			mockInputEmployee:  inputEmployee,
			mockOutputEmployee: internal.Employee{},
			mockError:          utils.ErrNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedResponse: map[string]any{
				"status":  http.StatusText(http.StatusNotFound),
				"message": "entity not found",
			},
			callService: true,
		},
		{
			name:               "CONFLICT",
			paramID:            "2",
			requestBody:        validBody,
			mockInputEmployee:  inputEmployee,
			mockOutputEmployee: internal.Employee{},
			mockError:          utils.ErrConflict,
			expectedStatusCode: http.StatusConflict,
			expectedResponse: map[string]any{
				"status":  http.StatusText(http.StatusConflict),
				"message": "entity already exists",
			},
			callService: true,
		},
		{
			name:               "INTERNAL_SERVER_ERROR",
			paramID:            "2",
			requestBody:        validBody,
			mockInputEmployee:  inputEmployee,
			mockOutputEmployee: internal.Employee{},
			mockError:          errors.New("some internal error"),
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: map[string]any{
				"status":  http.StatusText(http.StatusInternalServerError),
				"message": "internal server error",
			},
			callService: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			service := new(MockEmployeesService)
			empHandler := handler.NewEmployeeHandler(service)

			if tc.callService {
				idInt, _ := strconv.Atoi(tc.paramID)
				tc.mockInputEmployee.ID = idInt

				service.On("UpdateEmployee", tc.mockInputEmployee).Return(tc.mockOutputEmployee, tc.mockError)
			}

			request := httptest.NewRequest("PATCH", "/employees/"+tc.paramID, strings.NewReader(tc.requestBody))
			request.Header.Set("Content-Type", "application/json")

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tc.paramID)
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

			response := httptest.NewRecorder()

			empHandler.PatchEmployees()(response, request)

			require.Equal(t, tc.expectedStatusCode, response.Code)

			var actualResponse map[string]any
			err := json.Unmarshal(response.Body.Bytes(), &actualResponse)
			require.NoError(t, err)

			require.Equal(t, tc.expectedResponse, actualResponse)

			service.AssertExpectations(t)
		})
	}
}

func TestUnitEmployees_DeleteEmployee(t *testing.T) {
	type testCase struct {
		name               string
		paramID            string
		mockError          error
		expectedStatusCode int
		expectedResponse   map[string]any
		callService        bool
	}

	tests := []testCase{
		{
			name:               "INVALID ID",
			paramID:            "abc",
			mockError:          nil,
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: map[string]any{
				"status":  http.StatusText(http.StatusBadRequest),
				"message": "invalid format",
			},
			callService: false,
		},
		{
			name:               "NOT FOUND",
			paramID:            "2",
			mockError:          utils.ErrNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedResponse: map[string]any{
				"status":  http.StatusText(http.StatusNotFound),
				"message": "entity not found",
			},
			callService: true,
		},
		{
			name:               "INVALID ARGUMENTS",
			paramID:            "3",
			mockError:          errors.New("some random error"),
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedResponse: map[string]any{
				"status":  http.StatusText(http.StatusUnprocessableEntity),
				"message": "invalid arguments",
			},
			callService: true,
		},
		{
			name:               "OK",
			paramID:            "1",
			mockError:          nil,
			expectedStatusCode: http.StatusNoContent,
			expectedResponse: map[string]any{
				"message": "employee deleted successfully",
			},
			callService: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			service := new(MockEmployeesService)
			empHandler := handler.NewEmployeeHandler(service)

			if tc.callService {
				idInt, _ := strconv.Atoi(tc.paramID)
				service.On("DeleteEmployee", idInt).Return(tc.mockError)
			}

			request := httptest.NewRequest("DELETE", "/employees/"+tc.paramID, nil)
			response := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tc.paramID)
			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

			empHandler.DeleteEmployees()(response, request)

			require.Equal(t, tc.expectedStatusCode, response.Code)

			var actualResponse map[string]any
			err := json.Unmarshal(response.Body.Bytes(), &actualResponse)

			require.NoError(t, err)

			require.Equal(t, tc.expectedResponse, actualResponse)

			service.AssertExpectations(t)
		})
	}
}
