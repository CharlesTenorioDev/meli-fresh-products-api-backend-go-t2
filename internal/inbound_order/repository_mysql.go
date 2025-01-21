package inbound_order

import (
	"database/sql"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type InboundOrderRepository struct {
	db *sql.DB
}

func NewInboundOrderRepository(db *sql.DB) *InboundOrderRepository {
	return &InboundOrderRepository{db: db}
}

func (r *InboundOrderRepository) FindByID(id int) (internal.InboundOrder, error) {
	var order internal.InboundOrder
	order.Attributes = internal.InboundOrderAttributes{}

	err := r.db.QueryRow(`
		SELECT id, order_date, order_number, employee_id, product_batch_id, warehouse_id
		FROM inbound_orders
		WHERE id = ?`, id).Scan(&order.ID, &order.Attributes.OrderDate, &order.Attributes.OrderNumber, &order.Attributes.EmployeeID, &order.Attributes.ProductBatchID, &order.Attributes.WarehouseID)
	if err == sql.ErrNoRows {
		return internal.InboundOrder{}, utils.ErrNotFound
	}

	if err != nil {
		return internal.InboundOrder{}, err
	}

	return order, nil
}

func (r *InboundOrderRepository) Create(newOrder internal.InboundOrderAttributes) (internal.InboundOrder, error) {
	result, err := r.db.Exec("INSERT INTO inbound_orders (order_date, order_number, employee_id, product_batch_id, warehouse_id) VALUES (?, ?, ?, ?, ?)", newOrder.OrderDate, newOrder.OrderNumber, newOrder.EmployeeID, newOrder.ProductBatchID, newOrder.WarehouseID)
	if err != nil {
		return internal.InboundOrder{}, err
	}

	id, _ := result.LastInsertId()

	return r.FindByID(int(id))
}

func (r *InboundOrderRepository) FindByOrderNumber(orderNumber string) (internal.InboundOrder, error) {
	var order internal.InboundOrder
	order.Attributes = internal.InboundOrderAttributes{}

	err := r.db.QueryRow("SELECT id, order_date, order_number, employee_id, product_batch_id, warehouse_id FROM inbound_orders WHERE order_number = ?", orderNumber).Scan(&order.ID, &order.Attributes.OrderDate, &order.Attributes.OrderNumber, &order.Attributes.EmployeeID, &order.Attributes.ProductBatchID, &order.Attributes.WarehouseID)
	if err == sql.ErrNoRows {
		return internal.InboundOrder{}, utils.ErrNotFound
	}

	return order, err
}

func (r *InboundOrderRepository) GenerateReportForEmployee(employeeID int) (internal.EmployeeInboundOrdersReport, error) {
	var report internal.EmployeeInboundOrdersReport
	err := r.db.QueryRow(`
		SELECT e.id, e.id_card_number, e.first_name, e.last_name, e.warehouse_id, COUNT(o.id) as inbound_orders_count
		FROM employees e
		LEFT JOIN inbound_orders o ON e.id = o.employee_id
		WHERE e.id = ?
		GROUP BY e.id
	`, employeeID).Scan(&report.ID, &report.CardNumberID, &report.FirstName, &report.LastName, &report.WarehouseID, &report.InboundOrdersCount)

	if err == sql.ErrNoRows {
		return report, utils.ErrNotFound
	}

	if err != nil {
		return report, err
	}

	return report, nil
}

func (r *InboundOrderRepository) GenerateReport() ([]internal.EmployeeInboundOrdersReport, error) {
	report := []internal.EmployeeInboundOrdersReport{}
	rows, err := r.db.Query(`
		SELECT e.id, e.id_card_number, e.first_name, e.last_name, e.warehouse_id, COUNT(o.id) as inbound_orders_count
		FROM employees e
		LEFT JOIN inbound_orders o ON e.id = o.employee_id
		GROUP BY e.id
	`)

	if err != nil {
		return report, err
	}

	for rows.Next() {
		var row internal.EmployeeInboundOrdersReport
		if err := rows.Scan(&row.ID, &row.CardNumberID, &row.FirstName, &row.LastName, &row.WarehouseID, &row.InboundOrdersCount); err != nil {
			return []internal.EmployeeInboundOrdersReport{}, err
		}

		report = append(report, row)
	}

	return report, nil
}
