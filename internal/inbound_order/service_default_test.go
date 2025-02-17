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
		repo          func(repository *MockInboundOrderRepository)
		expectedError error
		expectedOrder internal.InboundOrder
	}

	cases := []testCase{
		{
			name: "Created - Successfully create a new inbound order",
			input: internal.InboundOrderAttributes{
				OrderDate:      "2021-04-04",
				OrderNumber:    "order#2742",
				EmployeeID:     1,
				ProductBatchID: 1,
				WarehouseID:    1,
			},
			repo: func(repository *MockInboundOrderRepository) {
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
			name: "Conflict - Order number already exists",
			input: internal.InboundOrderAttributes{
				OrderDate:      "2021-04-05",
				OrderNumber:    "order#1",
				EmployeeID:     2,
				ProductBatchID: 2,
				WarehouseID:    2,
			},
			repo: func(repository *MockInboundOrderRepository) {
				repository.On("FindByID", mock.Anything).Return(internal.InboundOrder{}, utils.ErrConflict)
				// No further configuration is needed, the error should occur after calling FindByID
			},
			expectedError: utils.ErrConflict,
		},
		{
			name: "Unprocessable Entity - Missing required field (OrderDate empty)",
			input: internal.InboundOrderAttributes{
				OrderDate:      "",
				OrderNumber:    "order#3",
				EmployeeID:     3,
				ProductBatchID: 3,
				WarehouseID:    3,
			},
			repo: func(repository *MockInboundOrderRepository) {
				// No further configuration is needed, the error should occur before calling the repository
			},
			expectedError: utils.ErrInvalidArguments,
		},
		{
			name: "Conflict - Order number cause error in FindByID",
			input: internal.InboundOrderAttributes{
				OrderDate:      "2021-04-05",
				OrderNumber:    "order#1",
				EmployeeID:     2,
				ProductBatchID: 2,
				WarehouseID:    2,
			},
			repo: func(repository *MockInboundOrderRepository) {
				repository.On("FindByID", mock.Anything).Return(internal.InboundOrder{}, utils.ErrNotFound)
				// No further configuration is needed, the error should occur after calling FindByID
			},
			expectedError: utils.ErrNotFound,
		},
		{
			name: "Conflict - Order number cause error in FindByOrderNumber",
			input: internal.InboundOrderAttributes{
				OrderDate:      "2021-04-05",
				OrderNumber:    "order#1",
				EmployeeID:     2,
				ProductBatchID: 2,
				WarehouseID:    2,
			},
			repo: func(repository *MockInboundOrderRepository) {
				repository.On("FindByID", mock.Anything).Return(internal.InboundOrder{
					ID: 1,
					Attributes: internal.InboundOrderAttributes{
						OrderDate:      "2021-04-05",
						OrderNumber:    "order#1",
						EmployeeID:     2,
						ProductBatchID: 2,
						WarehouseID:    2,
					},
				}, nil)
				repository.On("FindByOrderNumber", mock.Anything).Return(internal.InboundOrder{
					ID: 1,
					Attributes: internal.InboundOrderAttributes{
						OrderDate:      "2021-04-05",
						OrderNumber:    "order#1",
						EmployeeID:     2,
						ProductBatchID: 2,
						WarehouseID:    2,
					},
				}, utils.ErrConflict)
				// No further configuration is needed, the error should occur after calling FindByID
			},
			expectedError: utils.ErrConflict,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repository := new(MockInboundOrderRepository)
			service := inbound_order.NewInboundOrderService(repository)

			if tc.repo != nil {
				tc.repo(repository)
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

func TestUnitInboundOrder_GenerateInboundOrdersReport(t *testing.T) {
	tests := []struct {
		name           string
		ids            []int
		repo           func(repository *MockInboundOrderRepository)
		expectedReport []internal.EmployeeInboundOrdersReport
		expectedError  error
	}{
		{
			name: "Success - Generate report for all employees without ids",
			ids:  []int{},
			repo: func(repository *MockInboundOrderRepository) {
				repository.On("GenerateInboundOrdersReport").Return([]internal.EmployeeInboundOrdersReport{{ID: 2}}, nil)
			},
			expectedReport: []internal.EmployeeInboundOrdersReport{{ID: 2}},
			expectedError:  nil,
		},
		{
			name: "Error - Generate report for all employees withou ids",
			ids:  []int{},
			repo: func(repository *MockInboundOrderRepository) {
				repository.On("GenerateInboundOrdersReport").Return([]internal.EmployeeInboundOrdersReport{}, utils.ErrNotFound)
			},
			expectedReport: []internal.EmployeeInboundOrdersReport{},
			expectedError:  utils.ErrNotFound,
		},
		{
			name: "Error - Generate report for all employees with ids - Error Not Found",
			ids:  []int{2},
			repo: func(repository *MockInboundOrderRepository) {
				repository.On("FindByID", mock.Anything).Return(internal.InboundOrder{}, utils.ErrNotFound)
			},
			expectedReport: nil,
			expectedError:  utils.ErrNotFound,
		},
		{
			name: "Error - Generate report for all employees with ids - GenerateByIDInboundOrdersReport",
			ids:  []int{2},
			repo: func(repository *MockInboundOrderRepository) {
				repository.On("FindByID", mock.Anything).Return(internal.InboundOrder{}, nil)
				repository.On("GenerateByIDInboundOrdersReport", mock.Anything).Return(internal.EmployeeInboundOrdersReport{}, utils.ErrNotFound)
			},
			expectedReport: []internal.EmployeeInboundOrdersReport{},
			expectedError:  utils.ErrNotFound,
		},
		{
			name: "Error - Generate report for all employees with ids - GenerateByIDInboundOrdersReport (error different)",
			ids:  []int{2},
			repo: func(repository *MockInboundOrderRepository) {
				repository.On("FindByID", mock.Anything).Return(internal.InboundOrder{}, nil)
				repository.On("GenerateByIDInboundOrdersReport", mock.Anything).Return(internal.EmployeeInboundOrdersReport{}, utils.ErrConflict)
			},
			expectedReport: nil,
			expectedError:  utils.ErrConflict,
		},
		{
			name: "Success - Generate report for all employees with ids",
			ids:  []int{2},
			repo: func(repository *MockInboundOrderRepository) {
				repository.On("FindByID", mock.Anything).Return(internal.InboundOrder{}, nil)
				repository.On("GenerateByIDInboundOrdersReport", mock.Anything).Return(internal.EmployeeInboundOrdersReport{ID: 2}, nil)
			},
			expectedReport: []internal.EmployeeInboundOrdersReport{{ID: 2}},
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := new(MockInboundOrderRepository)
			service := inbound_order.NewInboundOrderService(repository)

			if tt.repo != nil {
				tt.repo(repository)
			}
			defer repository.AssertExpectations(t)

			report, err := service.GenerateInboundOrdersReport(tt.ids)

			require.Equal(t, tt.expectedReport, report)
			require.ErrorIs(t, err, tt.expectedError)

		})
	}
}
