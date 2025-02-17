package purchase_order_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/purchase_order"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

func TestPurchaseOrderRepository_FindAll(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := purchase_order.NewPurchaseOrderDB(db)

	t.Run("GIVEN a valid query WHEN executing FindAll THEN return purchase orders", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "order_number", "tracking_code", "order_date", "buyer_id"}).
			AddRow(1, "order#101", "abc123", "2024-02-16", 10)

		mock.ExpectQuery("SELECT po.id, po.order_number, po.tracking_code, po.order_date, po.buyer_id").
			WillReturnRows(rows)

		res, err := repo.FindAll()

		require.NoError(t, err)
		require.Len(t, res, 1)
		require.Equal(t, 1, res[0].ID)
	})

	t.Run("GIVEN no specific buyer ID WHEN executing FindAllByBuyerID THEN return all purchase orders summary", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"buyer_id", "total_orders", "order_codes"}).
			AddRow(10, 3, "order#101,67890,11223").
			AddRow(20, 2, "54321,98765")

		mock.ExpectQuery("SELECT po.buyer_id, COUNT.*GROUP BY po.buyer_id").
			WillReturnRows(rows)

		res, err := repo.FindAllByBuyerID(0)

		require.NoError(t, err)
		require.Len(t, res, 2)
		require.Equal(t, 10, res[0].BuyerID)
		require.Equal(t, 3, res[0].TotalOrders)
		require.Equal(t, "order#101,67890,11223", res[0].OrderCodes)
	})

	t.Run("GIVEN an invalid row WHEN executing FindAll THEN return a scan error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "order_number", "tracking_code", "order_date", "buyer_id"}).
			AddRow("invalid_id", "order#101", "abc123", "2024-02-16", 10)

		mock.ExpectQuery("SELECT po.id, po.order_number, po.tracking_code, po.order_date, po.buyer_id").
			WillReturnRows(rows)

		res, err := repo.FindAll()

		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("GIVEN a database error WHEN executing FindAll THEN return an error", func(t *testing.T) {
		mock.ExpectQuery("SELECT po.id, po.order_number, po.tracking_code, po.order_date, po.buyer_id").
			WillReturnError(errors.New("db error"))

		res, err := repo.FindAll()

		require.Error(t, err)
		require.Nil(t, res)
	})
}

func TestPurchaseOrderRepository_FindAllByBuyerID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := purchase_order.NewPurchaseOrderDB(db)

	t.Run("GIVEN a valid buyer ID WHEN executing FindAllByBuyerID THEN return purchase orders summary", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"buyer_id", "total_orders", "order_codes"}).
			AddRow(10, 2, "order#101,67890")

		mock.ExpectQuery("SELECT po.buyer_id, COUNT.*WHERE po.buyer_id = ?").
			WithArgs(10).
			WillReturnRows(rows)

		res, err := repo.FindAllByBuyerID(10)

		require.NoError(t, err)
		require.Len(t, res, 1)
		require.Equal(t, 10, res[0].BuyerID)
	})

	t.Run("GIVEN a database error WHEN executing FindAllByBuyerID THEN return an error", func(t *testing.T) {
		mock.ExpectQuery("SELECT po.buyer_id, COUNT.*").
			WillReturnError(errors.New("db error"))

		res, err := repo.FindAllByBuyerID(0)

		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("GIVEN a scan error WHEN executing FindAllByBuyerID THEN return an error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"buyer_id", "total_orders"}).
			AddRow(10, 2)

		mock.ExpectQuery("SELECT po.buyer_id, COUNT.*").
			WillReturnRows(rows)

		res, err := repo.FindAllByBuyerID(0)

		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("GIVEN no purchase orders WHEN executing FindAllByBuyerID THEN return ErrNotFound", func(t *testing.T) {
		emptyRows := sqlmock.NewRows([]string{"buyer_id", "total_orders", "order_codes"})

		mock.ExpectQuery("SELECT po.buyer_id, COUNT.*").
			WillReturnRows(emptyRows)

		res, err := repo.FindAllByBuyerID(0)

		require.Error(t, err)
		require.Nil(t, res)
		require.Equal(t, utils.ErrNotFound, err)
	})

}

func TestPurchaseOrderRepository_CreatePurchaseOrder(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := purchase_order.NewPurchaseOrderDB(db)

	t.Run("GIVEN a valid purchase order WHEN executing CreatePurchaseOrder THEN return inserted purchase order", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO purchase_orders").
			WithArgs("order#101", "2024-02-16", "abc123", 10, 5).
			WillReturnResult(sqlmock.NewResult(1, 1))

		order := internal.PurchaseOrderAttributes{
			OrderNumber:     "order#101",
			OrderDate:       "2024-02-16",
			TrackingCode:    "abc123",
			BuyerID:         10,
			ProductRecordID: 5,
		}

		res, err := repo.CreatePurchaseOrder(order)

		require.NoError(t, err)
		require.Equal(t, 1, res.ID)
	})

	t.Run("GIVEN a database error WHEN executing CreatePurchaseOrder THEN return an error", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO purchase_orders").
			WithArgs("order#101", "2024-02-16", "abc123", 10, 5).
			WillReturnError(errors.New("db insert error"))

		order := internal.PurchaseOrderAttributes{
			OrderNumber:     "order#101",
			OrderDate:       "2024-02-16",
			TrackingCode:    "abc123",
			BuyerID:         10,
			ProductRecordID: 5,
		}

		res, err := repo.CreatePurchaseOrder(order)

		require.Error(t, err)
		require.Equal(t, internal.PurchaseOrder{}, res)
		require.Equal(t, "db insert error", err.Error())
	})

	t.Run("GIVEN a database error WHEN retrieving LastInsertId THEN return an error", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO purchase_orders").
			WithArgs("order#101", "2024-02-16", "abc123", 10, 5).
			WillReturnResult(sqlmock.NewErrorResult(errors.New("last insert id error")))

		order := internal.PurchaseOrderAttributes{
			OrderNumber:     "order#101",
			OrderDate:       "2024-02-16",
			TrackingCode:    "abc123",
			BuyerID:         10,
			ProductRecordID: 5,
		}

		res, err := repo.CreatePurchaseOrder(order)

		require.Error(t, err)
		require.Equal(t, internal.PurchaseOrder{}, res)
		require.Equal(t, "last insert id error", err.Error())
	})

}
