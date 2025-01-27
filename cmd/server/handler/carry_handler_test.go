package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
