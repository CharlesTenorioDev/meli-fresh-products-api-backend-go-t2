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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type BuyerServiceMock struct {
	mock.Mock
}

func (m *BuyerServiceMock) GetAll() ([]internal.Buyer, error) {
	args := m.Called()
	return args.Get(0).([]internal.Buyer), args.Error(1)
}

func (m *BuyerServiceMock) GetOne(id int) (*internal.Buyer, error) {
	args := m.Called(id)
	return args.Get(0).(*internal.Buyer), args.Error(1)
}

func (m *BuyerServiceMock) CreateBuyer(buyer internal.BuyerAttributes) (*internal.Buyer, error) {
	args := m.Called(buyer)
	return args.Get(0).(*internal.Buyer), args.Error(1)
}

func (m *BuyerServiceMock) UpdateBuyer(buyer *internal.Buyer) (*internal.Buyer, error) {
	args := m.Called(buyer)
	return args.Get(0).(*internal.Buyer), args.Error(1)
}

func (m *BuyerServiceMock) DeleteBuyer(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestUnitBuyer_GetAllBuyers(t *testing.T) {
	type testCase struct {
		name            string
		mockBuyers      []internal.Buyer
		mockError       error
		expectedStatus  int
		expectedContent string
	}

	testCases := []testCase{
		{
			name:            "OK",
			mockBuyers:      []internal.Buyer{},
			mockError:       nil,
			expectedStatus:  http.StatusOK,
			expectedContent: "[]",
		},
		{
			name: "OK - one buyer",
			mockBuyers: []internal.Buyer{
				{
					ID: 1,
					BuyerAttributes: internal.BuyerAttributes{
						CardNumberID: "CARD123",
						FirstName:    "John",
						LastName:     "Doe",
					},
				},
			},
			mockError:       nil,
			expectedStatus:  http.StatusOK,
			expectedContent: "John",
		},
		{
			name:            "INTERNAL_SERVER_ERROR",
			mockBuyers:      nil,
			mockError:       errors.New("some service error"),
			expectedStatus:  http.StatusInternalServerError,
			expectedContent: "500 Erro Internal server error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service := new(BuyerServiceMock)
			service.
				On("GetAll").
				Return(tc.mockBuyers, tc.mockError)

			h := handler.NewBuyerHandler(service)

			req := httptest.NewRequest(http.MethodGet, "/buyers", nil)
			rr := httptest.NewRecorder()

			h.GetAll()(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)

			if tc.expectedStatus == http.StatusOK {
				if len(tc.mockBuyers) == 0 {
					assert.JSONEq(t, tc.expectedContent, rr.Body.String())
				} else {
					assert.Contains(t, rr.Body.String(), tc.expectedContent)
				}
			} else {
				assert.Contains(t, rr.Body.String(), tc.expectedContent)
			}
			service.AssertExpectations(t)
		})
	}
}

func TestUnitBuyer_GetOne(t *testing.T) {
	type testCase struct {
		name           string
		paramID        string
		mockID         int
		mockBuyer      *internal.Buyer
		mockError      error
		expectedStatus int
		expectedBody   string
	}

	testCases := []testCase{
		{
			name:    "OK",
			paramID: "1",
			mockID:  1,
			mockBuyer: &internal.Buyer{
				ID: 1,
				BuyerAttributes: internal.BuyerAttributes{
					CardNumberID: "CARD123",
					FirstName:    "John",
					LastName:     "Doe",
				},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   "John",
		},
		{
			name:           "BAD_REQUEST",
			paramID:        "abc",
			mockID:         0,
			mockBuyer:      nil,
			mockError:      errors.New("some user error"),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid syntax",
		},
		{
			name:           "INTERNAL_SERVER_ERROR",
			paramID:        "2",
			mockID:         2,
			mockBuyer:      nil,
			mockError:      errors.New("some service error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "some service error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service := new(BuyerServiceMock)

			service.On("GetOne", tc.mockID).Return(tc.mockBuyer, tc.mockError)

			h := handler.NewBuyerHandler(service)

			req := httptest.NewRequest(http.MethodGet, "/buyers/"+tc.paramID, nil)
			rr := httptest.NewRecorder()

			ctx := chi.NewRouteContext()
			ctx.URLParams.Add("id", tc.paramID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

			h.GetOne()(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)

			if tc.expectedBody != "" {
				assert.Contains(t, rr.Body.String(), tc.expectedBody)
			}

			service.AssertExpectations(t)
		})
	}
}

func TestBuyerHandler_CreateBuyer(t *testing.T) {
	cases := []struct {
		TestName           string
		ErrorToReturn      error
		Body               string
		ExpectedBody       string
		ExpectedStatusCode int
	}{
		{
			TestName:      "CreateBuyer - Success",
			ErrorToReturn: nil,
			Body: `{
					"card_number_id": "402323",
					"first_name": "Jhon",
					"last_name": "Doe"
					}`,
			ExpectedBody:       `{"data": {"id":1,"card_number_id":"402323","first_name":"Jhon","last_name":"Doe"}}`,
			ExpectedStatusCode: 201,
		},
		{
			TestName:      "CreateBuyer - Error Conflict",
			ErrorToReturn: utils.ErrConflict,
			Body: `{
					"card_number_id": "402323",
					"first_name": "Jhon",
					"last_name": "Doe"
					}`,
			ExpectedBody:       `{"message":"entity already exists", "status": "Conflict"}`,
			ExpectedStatusCode: 409,
		},
	}
	for _, tc := range cases {
		t.Run(tc.TestName, func(t *testing.T) {
			service := new(BuyerServiceMock)
			service.On("CreateBuyer", mock.Anything).Return(&internal.Buyer{
				ID: 1,
				BuyerAttributes: internal.BuyerAttributes{
					CardNumberID: "402323",
					FirstName:    "Jhon",
					LastName:     "Doe",
				},
			}, tc.ErrorToReturn)
			handler := handler.NewBuyerHandler(service)

			req, _ := http.NewRequest("POST", "http://localhost:8080/api/v1/buyers/", strings.NewReader(tc.Body))
			res := httptest.NewRecorder()
			handler.CreateBuyer()(res, req)
			require.Equal(t, tc.ExpectedStatusCode, res.Result().StatusCode)
			require.Equal(t, "application/json", res.Header().Get("Content-Type"))
			require.JSONEq(t, tc.ExpectedBody, res.Body.String())
		})
	}
}

func TestUnitBuyerHandler_UpdateBuyer(t *testing.T) {
	cases := []struct {
		TestName           string
		ErrorToReturn      error
		Body               string
		ExpectedBody       string
		ExpectedStatusCode int
	}{
		{
			TestName:      "UpdateBuyer - Success",
			ErrorToReturn: nil,
			Body: `{
					"card_number_id": "402323",
					"first_name": "Jhon",
					"last_name": "Doe"
					}`,
			ExpectedBody:       `{"data": {"id":1,"card_number_id":"402323","first_name":"Jhon","last_name":"Doe"}}`,
			ExpectedStatusCode: 201,
		},
		{
			TestName:      "UpdateBuyer - Error Conflict",
			ErrorToReturn: utils.ErrConflict,
			Body: `{
					"card_number_id": "402323",
					"first_name": "Jhon",
					"last_name": "Doe"
					}`,
			ExpectedBody:       `{"message":"entity already exists", "status": "Conflict"}`,
			ExpectedStatusCode: 409,
		},
	}
	for _, tc := range cases {
		t.Run(tc.TestName, func(t *testing.T) {
			service := new(BuyerServiceMock)
			service.On("UpdateBuyer", mock.Anything).Return(&internal.Buyer{
				ID: 1,
				BuyerAttributes: internal.BuyerAttributes{
					CardNumberID: "402323",
					FirstName:    "Jhon",
					LastName:     "Doe",
				},
			}, tc.ErrorToReturn)
			handler := handler.NewBuyerHandler(service)

			req, _ := http.NewRequest("PUT", "http://localhost:8080/api/v1/buyers/1", strings.NewReader(tc.Body))
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", "1")
			if tc.TestName == "UpdateProduct_ErrorBadRequest" {
				rctx.URLParams.Add("id", "a")
			}
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			res := httptest.NewRecorder()
			handler.UpdateBuyer()(res, req)
			require.Equal(t, tc.ExpectedStatusCode, res.Result().StatusCode)
			require.Equal(t, "application/json", res.Header().Get("Content-Type"))
			require.JSONEq(t, tc.ExpectedBody, res.Body.String())
		})
	}
}

func TestUnitBuyerHandler_DeleteBuyer(t *testing.T) {
	cases := []struct {
		TestName           string
		ErrorToReturn      error
		ExpectedBody       string
		ExpectedStatusCode int
	}{
		{
			TestName:           "DeleteBuyer - Success",
			ErrorToReturn:      nil,
			ExpectedBody:       `{"message":"Buyer deleted successfully"}`,
			ExpectedStatusCode: 204,
		},
		{
			TestName:           "DeleteBuyer - Error",
			ErrorToReturn:      utils.ErrNotFound,
			ExpectedBody:       `{"message":"entity not found", "status": "Internal Server Error"}`,
			ExpectedStatusCode: 500,
		},
	}
	for _, tc := range cases {
		t.Run(tc.TestName, func(t *testing.T) {
			service := new(BuyerServiceMock)
			service.On("DeleteBuyer", 1).Return(tc.ErrorToReturn)
			handler := handler.NewBuyerHandler(service)

			req, _ := http.NewRequest("DELETE", "http://localhost:8080/api/v1/buyers/1", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", "1")
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			res := httptest.NewRecorder()
			handler.DeleteBuyer()(res, req)
			require.Equal(t, tc.ExpectedStatusCode, res.Result().StatusCode)
			if tc.TestName == "DeleteBuyer - Error" {
				require.JSONEq(t, tc.ExpectedBody, res.Body.String())
			}
		})
	}
}
