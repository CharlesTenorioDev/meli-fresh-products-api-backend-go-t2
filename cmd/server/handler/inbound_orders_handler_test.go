package handler_test

import (
	"errors"
	"github.com/meli-fresh-products-api-backend-go-t2/cmd/server/handler"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockInBoundService struct {
	mock.Mock
}

func (m *MockInBoundService) CreateInboundOrder(inboundOrder internal.InboundOrderAttributes) (internal.InboundOrder, error) {
	args := m.Called(inboundOrder)
	return args.Get(0).(internal.InboundOrder), args.Error(1)
}

func (m *MockInBoundService) GenerateInboundOrdersReport(ids []int) ([]internal.EmployeeInboundOrdersReport, error) {
	args := m.Called(ids)
	return args.Get(0).([]internal.EmployeeInboundOrdersReport), args.Error(1)
}

func (m *MockInBoundService) FindByID(id int) (internal.InboundOrder, error) {
	args := m.Called(id)
	return args.Get(0).(internal.InboundOrder), args.Error(1)
}

func (m *MockInBoundService) FindByOrderNumber(orderNumber string) (internal.InboundOrder, error) {
	args := m.Called(orderNumber)
	return args.Get(0).(internal.InboundOrder), args.Error(1)
}

func (m *MockInBoundService) GenerateReportForEmployee(employeeID int) (internal.EmployeeInboundOrdersReport, error) {
	args := m.Called(employeeID)
	return args.Get(0).(internal.EmployeeInboundOrdersReport), args.Error(1)
}

func TestUnitInboundOrder_CreateInboundOrder(t *testing.T) {
	cases := []struct {
		TestName           string
		ErrorToReturn      error
		DataToReturn       internal.InboundOrder
		Body               string
		ExpectedStatusCode int
	}{
		{
			TestName:           "CREATED",
			Body:               `{"data":{"order_date": "2021-04-04", "order_number": "order#1999", "employee_id": 1, "product_batch_id": 1, "warehouse_id": 1}}`,
			DataToReturn:       internal.InboundOrder{},
			ExpectedStatusCode: 201,
			ErrorToReturn:      nil,
		},
		{
			TestName:           "BAD_REQUEST",
			Body:               `"data":{"order_date": "", "order_number": "order#19751", "employee_id": 1, "product_batch_id": 1, "warehouse_id": 1}}`,
			DataToReturn:       internal.InboundOrder{},
			ExpectedStatusCode: 400,
			ErrorToReturn:      nil,
		},
		{
			TestName:           "CONFLICT",
			Body:               `{"data":{"order_date": "2021-04-04", "order_number": "order#19751", "employee_id": 1, "product_batch_id": 1, "warehouse_id": 1}}`,
			DataToReturn:       internal.InboundOrder{},
			ExpectedStatusCode: 409,
			ErrorToReturn:      utils.ErrConflict,
		},
		{
			TestName:           "NOT_FOUND",
			Body:               `{"data":{"order_date": "2021-04-04", "order_number": "order#6762", "employee_id": 999, "product_batch_id": 1, "warehouse_id": 1}}`,
			DataToReturn:       internal.InboundOrder{},
			ExpectedStatusCode: 409,
			ErrorToReturn:      utils.ErrConflict,
		},
		{
			TestName:           "UNPROCESSABLE_ENTITY",
			Body:               `{"data":{"order_date": "", "order_number": "order#1999", "employee_id": 1, "product_batch_id": 1, "warehouse_id": 1}}`,
			DataToReturn:       internal.InboundOrder{},
			ExpectedStatusCode: 422,
			ErrorToReturn:      utils.ErrInvalidArguments,
		},
		{
			TestName:           "INTERNAL_SERVER_ERROR",
			Body:               `{"data":{"order_date": "2021-04-04", "order_number": "order#1975", "employee_id": 1, "product_batch_id": 1, "warehouse_id": 1}}`,
			DataToReturn:       internal.InboundOrder{},
			ExpectedStatusCode: 500,
			ErrorToReturn:      errors.New("Internal Server Error"),
		},
	}
	for _, c := range cases {
		t.Run(c.TestName, func(t *testing.T) {
			service := new(MockInBoundService)
			service.On("CreateInboundOrder", mock.Anything).Return(c.DataToReturn, c.ErrorToReturn)
			handler := handler.NewInboundOrderHandler(service)

			request := &http.Request{
				Body:   io.NopCloser(strings.NewReader(c.Body)),
				Header: http.Header{"Content-Type": []string{"application/json"}},
			}
			response := httptest.NewRecorder()
			handler.CreateInboundOrder()(response, request)

			require.Equal(t, c.ExpectedStatusCode, response.Result().StatusCode)
			require.Equal(t, "application/json", response.Header().Get("Content-Type"))
		})
	}
}

func TestUnitInboundOrder_GenerateInboundOrdersReport(t *testing.T) {
	cases := []struct {
		TestName           string
		DataToReturn       []internal.EmployeeInboundOrdersReport
		ErrorToReturn      error
		ExpectedStatusCode int
		Query              string
	}{
		{
			TestName:           "OK",
			DataToReturn:       []internal.EmployeeInboundOrdersReport{},
			ExpectedStatusCode: 200,
			ErrorToReturn:      nil,
			Query:              "",
		},
		{
			TestName:           "BAD_REQUEST",
			DataToReturn:       []internal.EmployeeInboundOrdersReport{},
			ExpectedStatusCode: 400,
			ErrorToReturn:      nil,
			Query:              "id=notID",
		},
		{
			TestName:           "NOT_FOUND",
			DataToReturn:       []internal.EmployeeInboundOrdersReport{},
			ExpectedStatusCode: 404,
			ErrorToReturn:      utils.ErrNotFound,
			Query:              "id=1109",
		},
		{
			TestName:           "INTERNAL_SERVER_ERROR",
			DataToReturn:       []internal.EmployeeInboundOrdersReport{},
			ExpectedStatusCode: 500,
			ErrorToReturn:      errors.New("Internal Server Error"),
			Query:              "id=1",
		},
	}
	for _, c := range cases {
		t.Run(c.TestName, func(t *testing.T) {
			service := new(MockInBoundService)
			service.On("GenerateInboundOrdersReport", mock.Anything).Return(c.DataToReturn, c.ErrorToReturn)
			handler := handler.NewInboundOrderHandler(service)
			request := &http.Request{
				URL:    &url.URL{RawQuery: c.Query},
				Header: http.Header{"Content-Type": []string{"application/json"}},
			}
			response := httptest.NewRecorder()
			handler.GenerateInboundOrdersReport()(response, request)

			require.Equal(t, c.ExpectedStatusCode, response.Result().StatusCode)
		})
	}

}
