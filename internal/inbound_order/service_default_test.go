package inbound_order_test

import (
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/inbound_order"
	"github.com/stretchr/testify/mock"
)

type MockInboundOrderRepository struct {
	mock.Mock
}

func (m *MockInboundOrderRepository) CreateInboundOrder(newOrder internal.InboundOrderAttributes) (internal.InboundOrder, error) {
	args := m.Called(newOrder)
	return args.Get(0).(internal.InboundOrder), args.Error(1)
}
func (m *MockInboundOrderRepository) GenerateInboundOrdersReport() ([]internal.EmployeeInboundOrdersReport, error) {
	args := m.Called()
	return args.Get(0).([]internal.EmployeeInboundOrdersReport), args.Error(1)
}
func (m *MockInboundOrderRepository) GenerateByIDInboundOrdersReport(employeeID int) (internal.EmployeeInboundOrdersReport, error) {
	args := m.Called(employeeID)
	return args.Get(0).(internal.EmployeeInboundOrdersReport), args.Error(1)
}

func (m *MockInboundOrderRepository) FindByID(id int) (internal.InboundOrder, error) {
	args := m.Called(id)
	return args.Get(0).(internal.InboundOrder), args.Error(1)
}

func (m *MockInboundOrderRepository) FindByOrderNumber(orderNumber string) (internal.InboundOrder, error) {
	args := m.Called(orderNumber)
	return args.Get(0).(internal.InboundOrder), args.Error(1)
}

func TestUnitInboundOrder_CreateInboundOrder(t *testing.T) {
	type testCase struct {
		name          string
		input         internal.InboundOrderAttributes
		mockSetup     func(repository *MockInboundOrderRepository)
		expectedError error
		expectedOrder internal.InboundOrder
	}

	cases := []testCase{
		{
			name: "201 Created - Successfully create a new inbound order",
			input: internal.InboundOrderAttributes{
				OrderDate:      "2021-04-04",
				OrderNumber:    "order#2742",
				EmployeeID:     1,
				ProductBatchID: 1,
				WarehouseID:    1,
			},
			mockSetup: func(repository *MockInboundOrderRepository) {
				repository.On("FindByID", mock.Anything).Return(internal.InboundOrder{}, nil)
				repository.On("FindByOrderNumber", "order#2742").Return(internal.InboundOrder{}, utils.ErrNotFound)
				repository.On("CreateInboundOrder", mock.Anything).Return(internal.InboundOrder{
					ID: 21,
					Attributes: internal.InboundOrderAttributes{
						OrderDate:      "2021-04-04",
						OrderNumber:    "order#2742",
						EmployeeID:     1,
						ProductBatchID: 1,
						WarehouseID:    1,
					},
				}, nil)
			},
			expectedError: nil,
			expectedOrder: internal.InboundOrder{
				ID: 21,
				Attributes: internal.InboundOrderAttributes{
					OrderDate:      "2021-04-04",
					OrderNumber:    "order#2742",
					EmployeeID:     1,
					ProductBatchID: 1,
					WarehouseID:    1,
				},
			},
		},
		{
			name: "409 Conflict - Order number already exists",
			input: internal.InboundOrderAttributes{
				OrderDate:      "2021-04-05",
				OrderNumber:    "order#1",
				EmployeeID:     2,
				ProductBatchID: 2,
				WarehouseID:    2,
			},
			mockSetup: func(repository *MockInboundOrderRepository) {
				repository.On("FindByID", mock.Anything).Return(internal.InboundOrder{}, utils.ErrConflict)
				// No further configuration is needed, the error should occur after calling FindByID
			},
			expectedError: utils.ErrConflict,
		},
		{
			name: "422 Unprocessable Entity - Missing required field (OrderDate empty)",
			input: internal.InboundOrderAttributes{
				OrderDate:      "",
				OrderNumber:    "order#3",
				EmployeeID:     3,
				ProductBatchID: 3,
				WarehouseID:    3,
			},
			mockSetup: func(repository *MockInboundOrderRepository) {
				// No further configuration is needed, the error should occur before calling the repository
			},
			expectedError: utils.ErrInvalidArguments,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repository := new(MockInboundOrderRepository)
			service := inbound_order.NewInboundOrderService(repository)

			if tc.mockSetup != nil {
				tc.mockSetup(repository)
			}

			newOrder, err := service.CreateInboundOrder(tc.input)

			if tc.expectedError != nil {
				require.ErrorIs(t, err, tc.expectedError)
				require.Empty(t, newOrder)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedOrder.Attributes.OrderDate, newOrder.Attributes.OrderDate)
				require.Equal(t, tc.expectedOrder.Attributes.OrderNumber, newOrder.Attributes.OrderNumber)
				require.Equal(t, tc.expectedOrder.Attributes.EmployeeID, newOrder.Attributes.EmployeeID)
				require.Equal(t, tc.expectedOrder.Attributes.ProductBatchID, newOrder.Attributes.ProductBatchID)
				require.Equal(t, tc.expectedOrder.Attributes.WarehouseID, newOrder.Attributes.WarehouseID)
			}

			repository.AssertExpectations(t)
		})
	}
}
