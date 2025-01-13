package repository

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type PurchaseOrderRepository struct {
	purchaseOrderTable map[int]internal.PurchaseOrder
}

func NewPurchaseOrderDb(orderTab map[int]internal.PurchaseOrder) *PurchaseOrderRepository {
	orderDb := make(map[int]internal.PurchaseOrder)
	if orderTab != nil {
		orderDb = orderTab
	}
	return &PurchaseOrderRepository{purchaseOrderTable: orderDb}
}

var ordersFile = "./docs/db/purchase_orders.json"

// LoadPurchaseOrders reads purchase orders from a JSON file
func (repo *PurchaseOrderRepository) LoadPurchaseOrders() (map[int]internal.PurchaseOrder, error) {
	file, err := os.ReadFile(ordersFile)
	if err != nil {
		log.Println("Error reading file:", err)
		return nil, err
	}

	var orders []internal.PurchaseOrderJson
	if err := json.Unmarshal(file, &orders); err != nil {
		log.Println("Error unmarshaling JSON:", err)
		return nil, err
	}

	repo.purchaseOrderTable = make(map[int]internal.PurchaseOrder)
	for _, order := range orders {
		repo.purchaseOrderTable[order.ID] = internal.PurchaseOrder{
			ID: order.ID,
			Attributes: internal.PurchaseOrderAttributes{
				OrderNumber:     order.OrderNumber,
				OrderDate:       order.OrderDate,
				TrackingCode:    order.TrackingCode,
				BuyerId:         order.BuyerId,
				ProductRecordId: order.ProductRecordId,
			},
		}
	}

	return repo.purchaseOrderTable, nil
}

// FindAll retrieves all purchase orders
func (repo *PurchaseOrderRepository) FindAll() (map[int]internal.PurchaseOrder, error) {
	ordersMap, err := repo.LoadPurchaseOrders()
	if err != nil {
		return nil, err
	}

	return ordersMap, nil
}

// FindById retrieves a specific purchase order by ID
func (repo *PurchaseOrderRepository) FindById(id int) (internal.PurchaseOrder, error) {
	ordersMap, err := repo.LoadPurchaseOrders()
	if err != nil {
		log.Println("Error loading orders:", err)
		return internal.PurchaseOrder{}, err
	}

	if order, exists := ordersMap[id]; exists {
		return order, nil
	}
	return internal.PurchaseOrder{}, utils.ErrNotFound
}

// CreatePurchaseOrder adds a new purchase order
func (repo *PurchaseOrderRepository) CreatePurchaseOrder(newOrder internal.PurchaseOrderAttributes) (internal.PurchaseOrder, error) {
	ordersMap, err := repo.FindAll()
	if err != nil {
		log.Println("Error loading orders:", err)
		return internal.PurchaseOrder{}, err
	}

	for _, o := range ordersMap {
		if o.Attributes.OrderNumber == newOrder.OrderNumber {
			log.Println("Order with this number already exists")
			return internal.PurchaseOrder{}, utils.ErrConflict
		}
	}

	newID := int(time.Now().UnixNano())
	order := internal.PurchaseOrder{
		ID:         newID,
		Attributes: newOrder,
	}

	ordersMap[newID] = order

	file, err := os.Create(ordersFile)
	if err != nil {
		log.Println("Error reopening file:", err)
		return internal.PurchaseOrder{}, err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(ordersMap); err != nil {
		log.Println("Error encoding JSON:", err)
		return internal.PurchaseOrder{}, err
	}

	log.Println("Order saved successfully!")
	return order, nil
}

// UpdatePurchaseOrder updates an existing purchase order
func (repo *PurchaseOrderRepository) UpdatePurchaseOrder(updatedOrder internal.PurchaseOrder) (internal.PurchaseOrder, error) {
	orders, err := repo.FindAll()
	if err != nil {
		log.Println("Error loading orders:", err)
		return internal.PurchaseOrder{}, err
	}

	updated := false
	for i, order := range orders {
		if order.ID == updatedOrder.ID {
			orders[i] = updatedOrder
			updated = true
			break
		}
	}

	if !updated {
		log.Println("Order with the given ID not found")
		return internal.PurchaseOrder{}, utils.ErrNotFound
	}

	file, err := os.Create(ordersFile)
	if err != nil {
		log.Println("Error opening file:", err)
		return internal.PurchaseOrder{}, err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(orders); err != nil {
		log.Println("Error encoding JSON:", err)
		return internal.PurchaseOrder{}, err
	}
	log.Println("Order updated successfully!")

	return updatedOrder, nil
}

// DeletePurchaseOrder deletes a purchase order by ID
func (repo *PurchaseOrderRepository) DeletePurchaseOrder(id int) error {
	ordersMap, err := repo.FindAll()
	if err != nil {
		log.Println("Error loading orders:", err)
		return err
	}

	if _, exists := ordersMap[id]; !exists {
		log.Println("Order with the given ID not found")
		return utils.ErrNotFound
	}

	delete(ordersMap, id)

	file, err := os.Create(ordersFile)
	if err != nil {
		log.Println("Error opening file for writing:", err)
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(ordersMap); err != nil {
		log.Println("Error encoding JSON:", err)
		return err
	}

	log.Println("Order deleted successfully!")
	return nil
}
