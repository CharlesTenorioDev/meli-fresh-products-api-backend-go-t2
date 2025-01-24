package service

import (
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

// PurchaseOrderDefault is the default implementation of the PurchaseOrder service
// it handles business logic and delegates data operations to the repository
type PurchaseOrderDefault struct {
	rp                   internal.PurchaseOrderRepository
	buyerService         internal.PurchaseOrdersBuyerValidation
	productRecordService internal.PurchaseOrdersProductRecordValidation
}

// NewPurchaseOrderService creates a new instance of PurchaseOrderDefault
// takes an PurchaseOrderRepository as a parameter to handle data operations
func NewPurchaseOrderService(rp internal.PurchaseOrderRepository, buyerService internal.PurchaseOrdersBuyerValidation, productRecordService internal.PurchaseOrdersProductRecordValidation) *PurchaseOrderDefault {
	return &PurchaseOrderDefault{rp: rp, buyerService: buyerService, productRecordService: productRecordService}
}

// FindAll retrieves all PurchaseOrders from the repository
func (s *PurchaseOrderDefault) FindAllByBuyerId(buyerId int) ([]internal.PurchaseOrderSummary, error) {
	purchaseOrdersSummary, err := s.rp.FindAllByBuyerId(buyerId)
	if err != nil {
		if err == utils.ErrBuyerDoesNotExists {
			return nil, utils.ErrBuyerDoesNotExists
		}
		return nil, err
	}
	return purchaseOrdersSummary, nil
}

// CreatePurchaseOrder adds a new purchaseOrder to the repository
func (s *PurchaseOrderDefault) CreatePurchaseOrder(newPurchaseOrder internal.PurchaseOrderAttributes) (purchaseOrder internal.PurchaseOrder, err error) {
	// validate required fields
	err = s.validateFields(newPurchaseOrder)
	if err != nil {
		return
	}

	// check for duplicates
	purchaseOrders, _ := s.rp.FindAll()
	err = s.validateDuplicates(purchaseOrders, newPurchaseOrder)
	if err != nil {
		return
	}

	// verify if buyer_id exists
	err = s.buyerExistsById(newPurchaseOrder.BuyerId)
	if err != nil {
		return
	}

	// verify if product_record_id exists
	err = s.productRecordExistsById(newPurchaseOrder.ProductRecordId)
	if err != nil {
		return
	}

	// attempt to create the new purchaseOrder
	return s.rp.CreatePurchaseOrder(newPurchaseOrder)
}

// validateFields checks if the required fields of a new purchaseOrder are not empty
func (s *PurchaseOrderDefault) validateFields(newPurchaseOrder internal.PurchaseOrderAttributes) (err error) {
	if newPurchaseOrder.OrderNumber == "" || newPurchaseOrder.OrderDate == "" || newPurchaseOrder.TrackingCode == "" || newPurchaseOrder.BuyerId == 0 || newPurchaseOrder.ProductRecordId == 0 {
		return utils.ErrEmptyArguments
	}
	return
}

// validateDuplicates ensures that no existing purchaseOrder has the same CardNumberId as the new purchaseOrder
func (s *PurchaseOrderDefault) validateDuplicates(purchaseOrders []internal.PurchaseOrder, newPurchaseOrder internal.PurchaseOrderAttributes) error {
	for _, purchaseOrder := range purchaseOrders {
		if purchaseOrder.Attributes.OrderNumber == newPurchaseOrder.OrderNumber {
			return utils.ErrConflict
		}
	}
	return nil
}

func (s *PurchaseOrderDefault) buyerExistsById(id int) error {
	possibleBuyer, err := s.buyerService.GetOne(id)
	// When internal server error
	if err != nil && err != utils.ErrNotFound {
		return err
	}
	if possibleBuyer == nil {
		return utils.ErrBuyerDoesNotExists
	}
	return nil
}

func (s *PurchaseOrderDefault) productRecordExistsById(id int) error {
	product, err := s.productRecordService.FindById(id)

	if err != nil {

		if err == utils.ErrNotFound {
			return utils.ErrProductDoesNotExists
		}
		return err
	}

	if product.ID == 0 {
		return utils.ErrProductDoesNotExists
	}

	return nil
}
