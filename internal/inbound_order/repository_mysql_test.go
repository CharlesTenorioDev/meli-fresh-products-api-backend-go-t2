package inbound_order

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/testutils"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUnitInboundOrderRepository(t *testing.T) {
	t.Run("GenerateInboundOrdersReport", func(t *testing.T) {
		t.Run("FAIL Generate Inbound Orders Report", func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			mock.ExpectQuery("SELECT e.id, e.id_card_number, e.first_name, e.last_name, e.warehouse_id, COUNT(o.id) as inbound_orders_count").
				WillReturnError(sql.ErrConnDone)

			repo := NewMySqlInboundOrderRepository(db)

			result, err := repo.GenerateInboundOrdersReport()
			require.Error(t, err)
			require.Empty(t, result)
		})

	})

	t.Run("GenerateByIDInboundOrdersReport", func(t *testing.T) {

		t.Run("FAIL Generate Inbound Orders Report by ID", func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			employeeID := 1

			mock.ExpectQuery("").
				WillReturnError(sql.ErrNoRows)

			repo := NewMySqlInboundOrderRepository(db)

			result, err := repo.GenerateByIDInboundOrdersReport(employeeID)

			if errors.Is(err, sql.ErrNoRows) {
				err = utils.ErrNotFound
			}
			require.Error(t, err)
			require.Empty(t, result)
			require.Equal(t, utils.ErrNotFound, err)
		})
	})

	t.Run("FindByID", func(t *testing.T) {
		t.Run("FAIL Find Inbound Order by ID (Error Not Found)", func(t *testing.T) {
			var (
				expectedOrder = internal.InboundOrder{}
				id            = 1
			)
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			mock.ExpectQuery("").
				WillReturnError(sql.ErrNoRows)

			repo := NewMySqlInboundOrderRepository(db)

			result, err := repo.FindByID(id)
			if errors.Is(err, sql.ErrNoRows) {
				err = utils.ErrNotFound
			}

			require.Error(t, err)
			require.Empty(t, result)
			require.Equal(t, utils.ErrNotFound, err)
			require.Equal(t, expectedOrder, result)
		})

		t.Run("FAIL Find Inbound Order by ID (Any other error)", func(t *testing.T) {
			var (
				expectedOrder = internal.InboundOrder{}
				id            = 1
			)
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			mock.ExpectQuery("").
				WillReturnError(sql.ErrConnDone)

			repo := NewMySqlInboundOrderRepository(db)

			result, err := repo.FindByID(id)
			if errors.Is(err, sql.ErrNoRows) {
				err = utils.ErrNotFound
			}

			require.Error(t, err)
			require.Empty(t, result)
			require.Equal(t, sql.ErrConnDone, err)
			require.Equal(t, expectedOrder, result)
		})
	})

	t.Run("FindByOrderNumber", func(t *testing.T) {
		t.Run("FAIL Find Inbound Order by Order Number", func(t *testing.T) {
			var (
				wantInboundOrder = internal.InboundOrder{}
				orderNumber      = "order#1"
			)
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			mock.ExpectQuery("").
				WillReturnError(sql.ErrNoRows)

			repo := NewMySqlInboundOrderRepository(db)

			result, err := repo.FindByOrderNumber(orderNumber)
			if errors.Is(err, sql.ErrNoRows) {
				err = utils.ErrNotFound
			}

			require.Error(t, err)
			require.Empty(t, result)
			require.Equal(t, utils.ErrNotFound, err)
			require.Equal(t, wantInboundOrder, result)
		})
	})
}

func TestIntegrationInboundOrderRepository(t *testing.T) {
	ts := testutils.GetTestDBConn()
	defer ts.End()
	repo := NewMySqlInboundOrderRepository(ts.DB)

	t.Run("NewMySqlInboundOrderRepository", func(t *testing.T) {
		require.NotNil(t, repo)
	})

	t.Run("Create", func(t *testing.T) {
		t.Run("SUCCESS Create Inbound Order", func(t *testing.T) {
			var (
				newOrder = internal.InboundOrderAttributes{
					OrderDate:      "2023-10-01",
					OrderNumber:    "12345",
					EmployeeID:     1,
					ProductBatchID: 1,
					WarehouseID:    1,
				}
			)
			order, err := repo.CreateInboundOrder(newOrder)
			require.NoError(t, err)
			require.NotZero(t, order.ID)
			require.Equal(t, newOrder.OrderNumber, order.Attributes.OrderNumber)
		})

		t.Run("FAIL Create Inbound Order", func(t *testing.T) {
			var (
				newOrder       = internal.InboundOrderAttributes{}
				wantErr  error = nil
			)
			order, err := repo.CreateInboundOrder(newOrder)
			require.Error(t, err)
			require.Empty(t, order.ID)
			require.NotEqual(t, nil, order)
			require.NotEqual(t, wantErr, err)
		})
	})

	t.Run("GenerateInboundOrdersReport", func(t *testing.T) {
		t.Run("SUCCESS Generate Inbound Orders Report", func(t *testing.T) {
			var (
				expectedReport = []internal.EmployeeInboundOrdersReport{
					{ID: 1, CardNumberID: "E001", FirstName: "Alice", LastName: "Johnson", WarehouseID: 1, InboundOrdersCount: 2},
					{ID: 2, CardNumberID: "E002", FirstName: "Bob", LastName: "Smith", WarehouseID: 2, InboundOrdersCount: 1}}
			)

			result, err := repo.GenerateInboundOrdersReport()
			require.NoError(t, err)
			require.NotEmpty(t, result)
			require.Equal(t, expectedReport, result)
		})
	})

	t.Run("GenerateByIDInboundOrdersReport", func(t *testing.T) {
		t.Run("SUCCESS Generate Inbound Orders Report by ID", func(t *testing.T) {
			var (
				expectedReport = internal.EmployeeInboundOrdersReport{ID: 1, CardNumberID: "E001", FirstName: "Alice", LastName: "Johnson", WarehouseID: 1, InboundOrdersCount: 2}
				employeeId     = 1
			)
			result, err := repo.GenerateByIDInboundOrdersReport(employeeId)
			require.NoError(t, err)
			require.NotEmpty(t, result)
			require.Equal(t, expectedReport, result)
		})
	})

	t.Run("FindByID", func(t *testing.T) {
		t.Run("SUCCESS Find Inbound Order by ID", func(t *testing.T) {
			var (
				expectedOrder = internal.InboundOrder{
					ID:         1,
					Attributes: internal.InboundOrderAttributes{OrderDate: "2025-01-05 12:00:00.000000", OrderNumber: "IN001", EmployeeID: 1, ProductBatchID: 1, WarehouseID: 1},
				}
				id = 1
			)
			result, err := repo.FindByID(id)

			require.NoError(t, err)
			require.NotEmpty(t, result)
			require.Equal(t, expectedOrder, result)
		})
	})

	t.Run("FindByOrderNumber", func(t *testing.T) {
		t.Run("SUCCESS Find Inbound Order by Order Number", func(t *testing.T) {
			var (
				expectedOrder = internal.InboundOrder{ID: 1, Attributes: internal.InboundOrderAttributes{OrderDate: "2025-01-05 12:00:00.000000", OrderNumber: "IN001", EmployeeID: 1, ProductBatchID: 1, WarehouseID: 1}}
				orderNumber   = "IN001"
			)
			result, err := repo.FindByOrderNumber(orderNumber)

			require.NoError(t, err)
			require.NotEmpty(t, result)
			require.Equal(t, expectedOrder, result)
		})
	})
}
