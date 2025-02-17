package inbound_order

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/DATA-DOG/go-txdb"
	"github.com/go-sql-driver/mysql"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
	"github.com/stretchr/testify/require"
	"testing"
)

func init() {
	cfg := mysql.Config{
		User:   "root",
		Passwd: "example",
		Net:    "tcp",
		Addr:   "localhost:3307",
		DBName: "fresh_products",
	}
	txdb.Register("txdb", "mysql", cfg.FormatDSN())
}

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("txdb", "identifier")
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	return db
}

func TestNewMySqlInboundOrderRepository(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewMySqlInboundOrderRepository(db)
	require.NotNil(t, repo)
}

func TestCreateInboundOrder(t *testing.T) {
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
		db := setupTestDB(t)
		defer db.Close()

		repo := NewMySqlInboundOrderRepository(db)

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
		db := setupTestDB(t)
		defer db.Close()

		repo := NewMySqlInboundOrderRepository(db)
		order, err := repo.CreateInboundOrder(newOrder)
		require.Error(t, err)
		require.Empty(t, order.ID)
		require.NotEqual(t, nil, order)
		require.NotEqual(t, wantErr, err)
	})
}

func TestMysqlInboundOrderRepository_GenerateInboundOrdersReport(t *testing.T) {
	t.Run("SUCCESS Generate Inbound Orders Report", func(t *testing.T) {
		var (
			expectedReport = []internal.EmployeeInboundOrdersReport{
				{
					ID:                 1,
					CardNumberID:       "123456",
					FirstName:          "John",
					LastName:           "Doe",
					WarehouseID:        1,
					InboundOrdersCount: 0,
				},
				{
					ID:                 2,
					CardNumberID:       "654321",
					FirstName:          "Jane",
					LastName:           "Smith",
					WarehouseID:        2,
					InboundOrdersCount: 1,
				},
			}
		)

		db := setupTestDB(t)
		defer db.Close()

		repo := NewMySqlInboundOrderRepository(db)

		result, err := repo.GenerateInboundOrdersReport()
		require.NoError(t, err)
		require.NotEmpty(t, result)
		require.Equal(t, expectedReport, result)
	})

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

}

func TestMysqlInboundOrderRepository_GenerateByIDInboundOrdersReport(t *testing.T) {
	t.Run("SUCCESS Generate Inbound Orders Report by ID", func(t *testing.T) {
		var (
			expectedReport = internal.EmployeeInboundOrdersReport{
				ID:                 1,
				CardNumberID:       "123456",
				FirstName:          "John",
				LastName:           "Doe",
				WarehouseID:        1,
				InboundOrdersCount: 0,
			}
			employeeId = 1
		)

		db := setupTestDB(t)
		defer db.Close()

		repo := NewMySqlInboundOrderRepository(db)

		result, err := repo.GenerateByIDInboundOrdersReport(employeeId)
		require.NoError(t, err)
		require.NotEmpty(t, result)
		require.Equal(t, expectedReport, result)
	})

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
}

func TestMysqlInboundOrderRepository_FindByID(t *testing.T) {
	t.Run("SUCCESS Find Inbound Order by ID", func(t *testing.T) {
		var (
			expectedOrder = internal.InboundOrder{
				ID: 1,
				Attributes: internal.InboundOrderAttributes{
					OrderDate:      "2021-04-04 00:00:00.000000",
					OrderNumber:    "order#1",
					EmployeeID:     2,
					ProductBatchID: 1,
					WarehouseID:    1,
				},
			}
			id = 1
		)

		db := setupTestDB(t)
		defer db.Close()

		repo := NewMySqlInboundOrderRepository(db)

		result, err := repo.FindByID(id)

		require.NoError(t, err)
		require.NotEmpty(t, result)
		require.Equal(t, expectedOrder, result)
	})

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
}

func TestMysqlInboundOrderRepository_FindByOrderNumber(t *testing.T) {
	t.Run("SUCCESS Find Inbound Order by Order Number", func(t *testing.T) {
		var (
			expectedOrder = internal.InboundOrder{
				ID: 1,
				Attributes: internal.InboundOrderAttributes{
					OrderDate:      "2021-04-04 00:00:00.000000",
					OrderNumber:    "order#1",
					EmployeeID:     2,
					ProductBatchID: 1,
					WarehouseID:    1,
				},
			}
			orderNumber = "order#1"
		)

		db := setupTestDB(t)
		defer db.Close()

		repo := NewMySqlInboundOrderRepository(db)

		result, err := repo.FindByOrderNumber(orderNumber)

		require.NoError(t, err)
		require.NotEmpty(t, result)
		require.Equal(t, expectedOrder, result)
	})

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
}
