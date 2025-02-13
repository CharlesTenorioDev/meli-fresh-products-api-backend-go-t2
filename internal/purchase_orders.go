package internal

import "time"

// PurchaseOrder represents an PurchaseOrder entity with its unique ID and attributes
type PurchaseOrder struct {
	ID         int `json:"id"`
	Attributes PurchaseOrderAttributes
}

// PurchaseOrderAttributes defines the details associated with an PurchaseOrder
type PurchaseOrderAttributes struct {
	OrderNumber     string `json:"order_number"`
	OrderDate       string `json:"order_date"`
	TrackingCode    string `json:"tracking_code"`
	BuyerID         int    `json:"buyer_id"`
	ProductRecordID int    `json:"product_record_id"`
}

// PurchaseOrderJSON defines the structure of the PurchaseOrder data as it appears in a json file
type PurchaseOrderJSON struct {
	ID              int       `json:"id"`
	OrderNumber     string    `json:"order_number"`
	OrderDate       time.Time `json:"order_date"`
	TrackingCode    string    `json:"tracking_code"`
	BuyerID         int       `json:"buyer_id"`
	ProductRecordID int       `json:"product_record_id"`
}

// PurchaseOrderRepository defines the interface for PurchaseOrder data persistence
// it specifies methods for fetching and creating PurchaseOrder data
type PurchaseOrderRepository interface {
	FindAll() ([]PurchaseOrder, error)
	FindAllByBuyerID(buyerID int) (PurchaseOrders []PurchaseOrderSummary, err error)
	CreatePurchaseOrder(newPurchaseOrder PurchaseOrderAttributes) (PurchaseOrder PurchaseOrder, err error)
}

// PurchaseOrderService defines the interface for PurchaseOrder-related business logic
// it includes methods for fetching and creating PurchaseOrders
type PurchaseOrderService interface {
	FindAllByBuyerID(buyerID int) (PurchaseOrders []PurchaseOrderSummary, err error)
	CreatePurchaseOrder(newPurchaseOrder PurchaseOrderAttributes) (PurchaseOrder PurchaseOrder, err error)
}

type PurchaseOrdersBuyerValidation interface {
	GetOne(int) (*Buyer, error)
}
type PurchaseOrdersProductRecordValidation interface {
	FindByID(productRecordID int) (ProductRecords, error)
}

type PurchaseOrderSummary struct {
	BuyerID     int    `json:"buyer_id"`
	TotalOrders int    `json:"total_orders"`
	OrderCodes  string `json:"order_codes"`
}
