package internal

import "time"

// PurchaseOrder represents an PurchaseOrder entity with its unique ID and attributes
type PurchaseOrder struct {
	ID         int `json:"id"`
	Attributes PurchaseOrderAttributes
}

// PurchaseOrderAttributes defines the details associated with an PurchaseOrder
type PurchaseOrderAttributes struct {
	OrderNumber     string    `json:"order_number"`
	OrderDate       time.Time `json:"order_date"`
	TrackingCode    string    `json:"tracking_code"`
	BuyerId         int       `json:"buyer_id"`
	ProductRecordId int       `json:"product_record_id"`
}

// PurchaseOrderJson defines the structure of the PurchaseOrder data as it appears in a json file
type PurchaseOrderJson struct {
	ID              int       `json:"id"`
	OrderNumber     string    `json:"order_number"`
	OrderDate       time.Time `json:"order_date"`
	TrackingCode    string    `json:"tracking_code"`
	BuyerId         int       `json:"buyer_id"`
	ProductRecordId int       `json:"product_record_id"`
}

// PurchaseOrderRepository defines the interface for PurchaseOrder data persistence
// it specifies methods for fetching and creating PurchaseOrder data
type PurchaseOrderRepository interface {
	FindAll() (PurchaseOrders map[int]PurchaseOrder, err error)
	FindById(id int) (PurchaseOrder PurchaseOrder, err error)
	CreatePurchaseOrder(newPurchaseOrder PurchaseOrderAttributes) (PurchaseOrder PurchaseOrder, err error)
	UpdatePurchaseOrder(inputPurchaseOrder PurchaseOrder) (PurchaseOrder PurchaseOrder, err error)
	DeletePurchaseOrder(id int) (err error)
}

// PurchaseOrderService defines the interface for PurchaseOrder-related business logic
// it includes methods for fetching and creating PurchaseOrders
type PurchaseOrderService interface {
	FindAll() (PurchaseOrders map[int]PurchaseOrder, err error)
	FindById(id int) (PurchaseOrder PurchaseOrder, err error)
	CreatePurchaseOrder(newPurchaseOrder PurchaseOrderAttributes) (PurchaseOrder PurchaseOrder, err error)
	UpdatePurchaseOrder(inputPurchaseOrder PurchaseOrder) (PurchaseOrder PurchaseOrder, err error)
	DeletePurchaseOrder(id int) (err error)
}

type PurchaseOrdersBuyerValidation interface {
	GetOne(int) (*Buyer, error)
}
