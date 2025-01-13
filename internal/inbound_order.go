package internal

type InboundOrder struct {
	ID         int                    `json:"id"`
	Attributes InboundOrderAttributes `json:"attributes"`
}

type InboundOrderAttributes struct {
	OrderDate      string `json: "order_date"`
	OrderNumber    string `json: "order_number"`
	EmployeeID     int    `json: "employee_id"`
	ProductBatchID int    `json: "product_batch_id"`
	WarehouseID    int    `json: "warehouse_id"`
}

type InboundOrderService interface {
	CreateOrder(newOrder InboundOrderAttributes) (InboundOrder, error)
	GetInboundOrdersReport(employeeID int) ([]EmployeeInboundOrdersReport, error)
}

type InboundOrderRepository interface {
	Create(newOrder InboundOrderAttributes) (InboundOrder, error)
	FindByID(id int) (InboundOrder, error)
	FindByOrderNumber(orderNumber string) (InboundOrder, error)
	GenerateReport(employeeID int) ([]EmployeeInboundOrdersReport, error)
}

type EmployeeInboundOrdersReport struct {
	ID                 int    `json: "id"`
	CardNumberID       string `json: "id_card_number"`
	FirstName          string `json: "first_name"`
	LastName           string `json: "last_name"`
	WarehouseID        int    `json: "warehouse_id"`
	InboundOrdersCount int    `json: "inbound_orders_count"`
}
