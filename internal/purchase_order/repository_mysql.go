package purchase_order

import (
	"database/sql"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type PurchaseOrderRepository struct {
	db *sql.DB
}

func NewPurchaseOrderDB(db *sql.DB) *PurchaseOrderRepository {
	return &PurchaseOrderRepository{db}
}

// FindAll retrieves all purchase orders
func (repo *PurchaseOrderRepository) FindAll() ([]internal.PurchaseOrder, error) {
	query := `
		SELECT po.id, po.order_number, po.tracking_code, po.order_date, po.buyer_id
		FROM purchase_orders po
		INNER JOIN buyers b ON po.buyer_id = b.id`

	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var purchaseOrders []internal.PurchaseOrder

	for rows.Next() {
		var po internal.PurchaseOrder

		err := rows.Scan(&po.ID, &po.Attributes.OrderNumber, &po.Attributes.TrackingCode, &po.Attributes.OrderDate, &po.Attributes.BuyerID)
		if err != nil {
			return nil, err
		}

		purchaseOrders = append(purchaseOrders, po)
	}

	return purchaseOrders, rows.Err()
}

// FindAllByBuyerID retrieves all purchase orders by buyer id
func (repo *PurchaseOrderRepository) FindAllByBuyerID(buyerID int) ([]internal.PurchaseOrderSummary, error) {
	var query string

	var rows *sql.Rows

	var err error

	if buyerID != 0 {
		query = `
			SELECT po.buyer_id, COUNT(po.id) AS total_orders, GROUP_CONCAT(po.order_number ORDER BY po.order_date) AS order_codes
			FROM purchase_orders po
			INNER JOIN buyers b ON po.buyer_id = b.id
			WHERE po.buyer_id = ?
			GROUP BY po.buyer_id`
		rows, err = repo.db.Query(query, buyerID)
	} else {
		query = `
			SELECT po.buyer_id, COUNT(po.id) AS total_orders, GROUP_CONCAT(po.order_number ORDER BY po.order_date) AS order_codes
			FROM purchase_orders po
			INNER JOIN buyers b ON po.buyer_id = b.id
			GROUP BY po.buyer_id`
		rows, err = repo.db.Query(query)
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var purchaseOrders []internal.PurchaseOrderSummary

	for rows.Next() {
		var summary internal.PurchaseOrderSummary

		err := rows.Scan(&summary.BuyerID, &summary.TotalOrders, &summary.OrderCodes)
		if err != nil {
			return nil, err
		}

		purchaseOrders = append(purchaseOrders, summary)
	}

	if len(purchaseOrders) == 0 {
		return nil, utils.ErrNotFound
	}

	return purchaseOrders, rows.Err()
}

// CreatePurchaseOrder adds a new purchase order
func (repo *PurchaseOrderRepository) CreatePurchaseOrder(newOrder internal.PurchaseOrderAttributes) (internal.PurchaseOrder, error) {
	query := "INSERT INTO purchase_orders (order_number, order_date, tracking_code, buyer_id, product_record_id) VALUES (?, ?, ?, ?, ?)"

	result, err := repo.db.Exec(query, newOrder.OrderNumber, newOrder.OrderDate, newOrder.TrackingCode, newOrder.BuyerID, newOrder.ProductRecordID)
	if err != nil {
		return internal.PurchaseOrder{}, err
	}

	insertedID, err := result.LastInsertId()
	if err != nil {
		return internal.PurchaseOrder{}, err
	}

	purchaseOrder := internal.PurchaseOrder{
		ID:         int(insertedID),
		Attributes: newOrder,
	}

	return purchaseOrder, nil
}
