package service

import (
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

// PurchaseOrderDefault is the default implementation of the PurchaseOrder service
// it handles business logic and delegates data operations to the repository
type PurchaseOrderDefault struct {
	rp           internal.PurchaseOrderRepository
	buyerService internal.PurchaseOrdersBuyerValidation
}

// NewPurchaseOrderService creates a new instance of PurchaseOrderDefault
// takes an PurchaseOrderRepository as a parameter to handle data operations
func NewPurchaseOrderService(rp internal.PurchaseOrderRepository, buyerService internal.PurchaseOrdersBuyerValidation) *PurchaseOrderDefault {
	return &PurchaseOrderDefault{rp: rp, buyerService: buyerService}
}

// FindAll retrieves all PurchaseOrders from the repository
func (s *PurchaseOrderDefault) FindAll() (purchaseOrders map[int]internal.PurchaseOrder, err error) {
	purchaseOrders, err = s.rp.FindAll()
	return
}

// FindById retrieves an purchaseOrder by ID from the repository
func (s *PurchaseOrderDefault) FindById(id int) (purchaseOrder internal.PurchaseOrder, err error) {
	purchaseOrder, err = s.rp.FindById(id)
	return
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

	// attempt to create the new purchaseOrder
	return s.rp.CreatePurchaseOrder(newPurchaseOrder)
}

// UpdatepPurchaseOrder updates an purchaseOrder in the repository
func (s *PurchaseOrderDefault) UpdatePurchaseOrder(inputPurchaseOrder internal.PurchaseOrder) (purchaseOrder internal.PurchaseOrder, err error) {
	// find the existing purchaseOrder
	internalPurchaseOrder, err := s.rp.FindById(inputPurchaseOrder.ID)
	if err != nil {
		err = utils.ErrNotFound
		return
	}

	// verify if buyer_id exists
	err = s.buyerExistsById(inputPurchaseOrder.Attributes.BuyerId)
	if err != nil {
		return
	}

	// merge input fields with the existing purchaseOrder
	updatedPurchaseOrder := mergePurchaseOrderFields(inputPurchaseOrder, internalPurchaseOrder)

	// update the purchaseOrder in the repository
	purchaseOrder, err = s.rp.UpdatePurchaseOrder(updatedPurchaseOrder)
	return
}

// DeletePurchaseOrder deletes an purchaseOrder from the repository based on the provided ID
func (s *PurchaseOrderDefault) DeletePurchaseOrder(id int) (err error) {
	// find the purchaseOrder to ensure it exists
	purchaseOrder, err := s.rp.FindById(id)
	if err != nil {
		return utils.ErrNotFound
	}

	// delete the purchaseOrder by passing only the ID to the repository
	err = s.rp.DeletePurchaseOrder(purchaseOrder.ID)
	if err != nil {
		return utils.ErrInvalidArguments
	}

	return nil
}

// validateFields checks if the required fields of a new purchaseOrder are not empty
func (s *PurchaseOrderDefault) validateFields(newPurchaseOrder internal.PurchaseOrderAttributes) (err error) {
	if newPurchaseOrder.OrderNumber == "" || newPurchaseOrder.TrackingCode == "" || newPurchaseOrder.BuyerId == 0 || newPurchaseOrder.ProductRecordId == 0 {
		return utils.ErrEmptyArguments
	}
	return
}

// validateDuplicates ensures that no existing purchaseOrder has the same CardNumberId as the new purchaseOrder
func (s *PurchaseOrderDefault) validateDuplicates(purchaseOrders map[int]internal.PurchaseOrder, newPurchaseOrder internal.PurchaseOrderAttributes) error {
	for _, purchaseOrder := range purchaseOrders {
		if purchaseOrder.Attributes.OrderNumber == newPurchaseOrder.OrderNumber {
			return utils.ErrConflict
		}
	}
	return nil
}

// mergePurchaseOrderFields merges the fields of the input purchaseOrder with the internal purchaseOrder
func mergePurchaseOrderFields(inputPurchaseOrder, internalPurchaseOrder internal.PurchaseOrder) (updatedPurchaseOrder internal.PurchaseOrder) {
	updatedPurchaseOrder.ID = internalPurchaseOrder.ID

	if inputPurchaseOrder.Attributes.OrderNumber != "" {
		updatedPurchaseOrder.Attributes.OrderNumber = inputPurchaseOrder.Attributes.OrderNumber
	} else {
		updatedPurchaseOrder.Attributes.OrderNumber = internalPurchaseOrder.Attributes.OrderNumber
	}

	if inputPurchaseOrder.Attributes.TrackingCode != "" {
		updatedPurchaseOrder.Attributes.TrackingCode = inputPurchaseOrder.Attributes.TrackingCode
	} else {
		updatedPurchaseOrder.Attributes.TrackingCode = internalPurchaseOrder.Attributes.TrackingCode
	}

	if inputPurchaseOrder.Attributes.BuyerId != 0 {
		updatedPurchaseOrder.Attributes.BuyerId = inputPurchaseOrder.Attributes.BuyerId
	} else {
		updatedPurchaseOrder.Attributes.BuyerId = internalPurchaseOrder.Attributes.BuyerId
	}

	if inputPurchaseOrder.Attributes.ProductRecordId != 0 {
		updatedPurchaseOrder.Attributes.ProductRecordId = inputPurchaseOrder.Attributes.ProductRecordId
	} else {
		updatedPurchaseOrder.Attributes.ProductRecordId = internalPurchaseOrder.Attributes.ProductRecordId
	}

	return updatedPurchaseOrder
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
