package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockCarryService struct {
	mock.Mock
}

func (m *mockCarryService) GetByID(id int) (internal.Carry, error) {
	args := m.Called(id)
	return args.Get(0).(internal.Carry), args.Error(1)
}

func (m *mockCarryService) GetAll() ([]internal.Carry, error) {
	args := m.Called()
	return args.Get(0).([]internal.Carry), args.Error(1)
}

func (m *mockCarryService) Save(carry *internal.Carry) error {
	args := m.Called(carry)
	return args.Error(0)
}

func (m *mockCarryService) Update(carry *internal.Carry) error {
	args := m.Called(carry)
	return args.Error(0)
}

func (m *mockCarryService) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestUnitCarryHandler_SaveCarry(t *testing.T) {
	cases := []struct {
		TestName           string
		ErrorToReturn      error
		Body               string
		ExpectedBody       string
		ExpectedStatusCode int
	}{
		{
			TestName:      "SaveCarry",
			ErrorToReturn: nil,
			Body: `{
					"cid": 6,
					"company_name": "New Alkemy",
					"address": "Monroe 860",
					"telephone": "47470000",
					"locality_id": 2
					}`,
			ExpectedBody: `{
									"data": {
										"id": 0,
										"cid": 6,
										"company_name": "New Alkemy",
										"address": "Monroe 860",
										"telephone": "47470000",
										"locality_id": 2
									}
								}`,
			ExpectedStatusCode: http.StatusCreated,
		},
		{
			TestName:           "SaveCarryError_ErrorUnprocessableEntity",
			ErrorToReturn:      utils.EZeroValue("Carry"),
			Body:               `{"cid": 6,"company_name": "New Alkemy","address": "Monroe 860","telephone": "47470000","locality_id": 2}`,
			ExpectedBody:       `{"status":"Unprocessable Entity","message":"invalid arguments: Carry cannot be empty/null"}`,
			ExpectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			TestName:      "SaveCarryError_ErrorConflict",
			ErrorToReturn: utils.EConflict("Carry", "CID"),
			Body:          `{"cid": 6,"company_name": "New Alkemy","address": "Monroe 860","telephone": "47470000","locality_id": 2}`,
			ExpectedBody: `{
								"status": "Conflict",
								"message": "entity already exists: Carry with attribute 'CID' already exists"
							}`,
			ExpectedStatusCode: http.StatusConflict,
		},
	}
	for _, c := range cases {
		t.Run(c.TestName, func(t *testing.T) {
			service := new(mockCarryService)
			service.On("Save", mock.Anything).Return(c.ErrorToReturn)
			handler := CarryHandler{service: service}
			req, _ := http.NewRequest("POST", "http://localhost:8080/api/v1/carries/", strings.NewReader(c.Body))
			res := httptest.NewRecorder()
			funcHandler := handler.SaveCarry()
			funcHandler(res, req)
			require.Equal(t, c.ExpectedStatusCode, res.Result().StatusCode)
			require.Equal(t, "application/json", res.Header().Get("Content-Type"))
			require.JSONEq(t, c.ExpectedBody, res.Body.String())

		})
	}
}

func TestUnitCarryHandler_GetAllCarries(t *testing.T) {
	cases := []struct {
		TestName           string
		ErrorToReturn      error
		ExpectedBody       string
		ExpectedStatusCode int
	}{
		{
			TestName:      "GetAllCarries",
			ErrorToReturn: nil,
			ExpectedBody: `{
							"data": []
						}`,
			ExpectedStatusCode: http.StatusOK,
		},
		{
			TestName:           "GetAllCarriesError_ErrorNotFound",
			ErrorToReturn:      utils.ENotFound("Carry"),
			ExpectedBody:       `{"status":"Not Found","message":"entity not found: Carry doesn't exist"}`,
			ExpectedStatusCode: http.StatusNotFound,
		},
	}
	for _, c := range cases {
		t.Run(c.TestName, func(t *testing.T) {
			service := new(mockCarryService)
			service.On("GetAll").Return([]internal.Carry{}, c.ErrorToReturn)
			handler := CarryHandler{service: service}
			req, _ := http.NewRequest("GET", "http://localhost:8080/api/v1/carries/", nil)
			res := httptest.NewRecorder()
			funcHandler := handler.GetAllCarries()
			funcHandler(res, req)
			require.Equal(t, c.ExpectedStatusCode, res.Result().StatusCode)
			require.Equal(t, "application/json", res.Header().Get("Content-Type"))
			require.JSONEq(t, c.ExpectedBody, res.Body.String())

		})
	}
}

func TestUnitCarryHandler_GetCarryByID(t *testing.T) {
	cases := []struct {
		TestName           string
		ErrorToReturn      error
		ExpectedBody       string
		ExpectedStatusCode int
	}{
		{
			TestName:      "GetCarryByID",
			ErrorToReturn: nil,
			ExpectedBody: `{
							"data": {
								"id": 0,
								"cid": 6,
								"company_name": "New Alkemy",
								"address": "Monroe 860",
								"telephone": "47470000",
								"locality_id": 2
							}
						}`,
			ExpectedStatusCode: http.StatusOK,
		},
		{
			TestName:           "GetCarryByIDError_ErrorNotFound",
			ErrorToReturn:      utils.ENotFound("Carry"),
			ExpectedBody:       `{"status":"Not Found","message":"entity not found: Carry doesn't exist"}`,
			ExpectedStatusCode: http.StatusNotFound,
		},
		{
			TestName:           "GetCarryByIDError_ErrorBadRequest",
			ErrorToReturn:      utils.EBadRequest("Invalid ID"),
			ExpectedBody:       `{"status":"Bad Request","message":"invalid format: Invalid ID with invalid format"}`,
			ExpectedStatusCode: http.StatusBadRequest,
		},
	}
	for _, c := range cases {
		t.Run(c.TestName, func(t *testing.T) {
			service := new(mockCarryService)
			service.On("GetByID", 1).Return(internal.Carry{
				ID:          0,
				CID:         6,
				CompanyName: "New Alkemy",
				Address:     "Monroe 860",
				Telephone:   "47470000",
				LocalityID:  2,
			}, c.ErrorToReturn)
			handler := CarryHandler{service: service}
			req, _ := http.NewRequest("GET", "http://localhost:8080/api/v1/carries/", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", "1")
			if c.TestName == "GetCarryByIDError_ErrorBadRequest" {
				rctx.URLParams.Add("id", "a")
			}
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			res := httptest.NewRecorder()
			funcHandler := handler.GetCarryByID()
			funcHandler(res, req)
			require.Equal(t, c.ExpectedStatusCode, res.Result().StatusCode)
			require.Equal(t, "application/json", res.Header().Get("Content-Type"))
			require.JSONEq(t, c.ExpectedBody, res.Body.String())

		})
	}
}

func TestUnitCarryHandler_UpdateCarry(t *testing.T) {
	cases := []struct {
		TestName           string
		ErrorToReturn      error
		Body               string
		ExpectedBody       string
		ExpectedStatusCode int
	}{
		{
			TestName:      "UpdateCarry",
			ErrorToReturn: nil,
			Body: `{
					"cid": 6,
					"company_name": "New Alkemy",
					"address": "Monroe 860",
					"telephone": "47470000",
					"locality_id": 2
					}`,
			ExpectedBody: `{
									"data": {
										"id": 1,
										"cid": 6,
										"company_name": "New Alkemy",
										"address": "Monroe 860",
										"telephone": "47470000",
										"locality_id": 2
									}
								}`,
			ExpectedStatusCode: http.StatusOK,
		},
		{
			TestName:           "UpdateCarryError_ErrorUnprocessableEntity",
			ErrorToReturn:      utils.EZeroValue("Carry"),
			Body:               `{"cid": 6,"company_name": "New Alkemy","address": "Monroe 860","telephone": "47470000","locality_id": 2}`,
			ExpectedBody:       `{"status":"Unprocessable Entity","message":"invalid arguments: Carry cannot be empty/null"}`,
			ExpectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			TestName:      "UpdateCarryError_ErrorNotFound",
			ErrorToReturn: utils.ENotFound("Carry"),
			Body: `{
					"cid": 6,
					"company_name": "New Alkemy",
					"address": "Monroe 860",
					"telephone": "47470000",
					"locality_id": 2
					}`,
			ExpectedBody: `{
								"status": "Not Found",
								"message": "entity not found: Carry doesn't exist"
							}`,
			ExpectedStatusCode: http.StatusNotFound,
		},
		{
			TestName:      "UpdateCarryError_ErrorBadRequest",
			ErrorToReturn: utils.EBadRequest("Invalid ID"),
			Body: `{"cid": 6,"company
			_name": "New Alkemy","address": "Monroe 860","telephone": "47470000","locality_id": 2}`,
			ExpectedBody:       `{"status":"Bad Request","message":"invalid format: Invalid ID with invalid format"}`,
			ExpectedStatusCode: http.StatusBadRequest,
		},
	}
	for _, c := range cases {
		t.Run(c.TestName, func(t *testing.T) {
			service := new(mockCarryService)
			service.On("Update", mock.Anything).Return(c.ErrorToReturn)
			handler := CarryHandler{service: service}
			req, _ := http.NewRequest("PUT", "http://localhost:8080/api/v1/carries/", strings.NewReader(c.Body))
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", "1")
			if c.TestName == "UpdateCarryError_ErrorBadRequest" {
				rctx.URLParams.Add("id", "a")
			}
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			res := httptest.NewRecorder()
			funcHandler := handler.UpdateCarry()
			funcHandler(res, req)
			require.Equal(t, c.ExpectedStatusCode, res.Result().StatusCode)
			require.Equal(t, "application/json", res.Header().Get("Content-Type"))
			require.JSONEq(t, c.ExpectedBody, res.Body.String())

		})
	}
}

func TestUnitCarryHandler_DeleteCarry(t *testing.T) {
	cases := []struct {
		TestName           string
		ErrorToReturn      error
		ExpectedBody       string
		ExpectedStatusCode int
	}{
		{
			TestName:      "DeleteCarry",
			ErrorToReturn: nil,
			ExpectedBody: `{
							"data": "Carry deleted successfully"
						}`,
			ExpectedStatusCode: http.StatusNoContent,
		},
		{
			TestName:           "DeleteCarryError_ErrorNotFound",
			ErrorToReturn:      utils.ENotFound("Carry"),
			ExpectedBody:       `{"status":"Not Found","message":"entity not found: Carry doesn't exist"}`,
			ExpectedStatusCode: http.StatusNotFound,
		},
		{
			TestName:           "DeleteCarryError_ErrorBadRequest",
			ErrorToReturn:      utils.EBadRequest("Invalid ID"),
			ExpectedBody:       `{"status":"Bad Request","message":"invalid format: Invalid ID with invalid format"}`,
			ExpectedStatusCode: http.StatusBadRequest,
		},
	}
	for _, c := range cases {
		t.Run(c.TestName, func(t *testing.T) {
			service := new(mockCarryService)
			service.On("Delete", 1).Return(c.ErrorToReturn)
			handler := CarryHandler{service: service}
			req, _ := http.NewRequest("DELETE", "http://localhost:8080/api/v1/carries/", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", "1")
			if c.TestName == "DeleteCarryError_ErrorBadRequest" {
				rctx.URLParams.Add("id", "a")
			}
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			res := httptest.NewRecorder()
			funcHandler := handler.DeleteCarry()
			funcHandler(res, req)
			require.Equal(t, c.ExpectedStatusCode, res.Result().StatusCode)
			require.Equal(t, "application/json", res.Header().Get("Content-Type"))
			require.JSONEq(t, c.ExpectedBody, res.Body.String())

		})
	}
}
