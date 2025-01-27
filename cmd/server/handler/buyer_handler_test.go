package handler_test

import (
	"errors"
	"github.com/meli-fresh-products-api-backend-go-t2/cmd/server/handler"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockBuyerService struct {
	mock.Mock
}

func (m *MockBuyerService) GetAll() ([]internal.Buyer, error) {
	args := m.Called()
	return args.Get(0).([]internal.Buyer), args.Error(1)
}

func (m *MockBuyerService) GetOne(id int) (*internal.Buyer, error) {
	args := m.Called(id)
	return args.Get(0).(*internal.Buyer), args.Error(1)
}

func (m *MockBuyerService) CreateBuyer(internal.BuyerAttributes) (*internal.Buyer, error) {
	args := m.Called()
	return args.Get(0).(*internal.Buyer), args.Error(1)
}

func (m *MockBuyerService) UpdateBuyer(*internal.Buyer) (*internal.Buyer, error) {
	args := m.Called()
	return args.Get(0).(*internal.Buyer), args.Error(1)
}

func (m *MockBuyerService) DeleteBuyer(id int) error {
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
			service := new(MockBuyerService)
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
