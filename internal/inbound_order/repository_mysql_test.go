package inbound_order_test

import (
	"database/sql"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/inbound_order"
	"testing"

	"github.com/DATA-DOG/go-txdb"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	// "github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
	"github.com/stretchr/testify/require"
)

func init() {
	txdb.Register("txdb", "mysql", "root:password@tcp(localhost:3306)/dbname")
}

func TestInboundOrderRepository_Create(t *testing.T) {
	db, err := sql.Open("txdb", "fantasy_order")
	require.NoError(t, err)
	defer db.Close()

	repo := inbound_order.NewMySqlInboundOrderRepository(db)

	t.Run("Creating a new inbound order", func(t *testing.T) {
		newInboundOrder := internal.InboundOrderAttributes{
			OrderDate:      "2021-04-04",
			OrderNumber:    "order#666",
			EmployeeID:     1,
			ProductBatchID: 1,
			WarehouseID:    1,
		}
		inboundOrder, err := repo.CreateOrder(newInboundOrder)
		require.NoError(t, err)
		require.Equal(t, "2021-04-04", inboundOrder.Attributes.OrderDate)
		require.Equal(t, "order#666", inboundOrder.Attributes.OrderNumber)
		require.NotZero(t, inboundOrder.Attributes.EmployeeID)
		require.NotZero(t, inboundOrder.Attributes.ProductBatchID)
		require.NotZero(t, inboundOrder.Attributes.WarehouseID)
	})
}
