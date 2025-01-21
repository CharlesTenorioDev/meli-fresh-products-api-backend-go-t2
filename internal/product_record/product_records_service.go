package product_record

import (
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type ProductRecordsService struct {
	repo              internal.ProductRecordsRepository
	validationProduct internal.ProductValidation
}

func NewProductRecordService(repo internal.ProductRecordsRepository, validationProduct internal.ProductValidation) *ProductRecordsService {
	return &ProductRecordsService{
		repo:              repo,
		validationProduct: validationProduct,
	}
}

func (s *ProductRecordsService) GetProductRecords(productID int) ([]internal.ProductReport, error) {
	if productID > 0 {
		if _, err := s.validationProduct.GetProductByID(productID); err != nil {
			return nil, utils.ErrNotFound
		}
	}

	productReports, err := s.repo.Read(productID)
	if err != nil {
		return nil, err
	}

	if len(productReports) == 0 {
		return nil, utils.ErrNotFound
	}

	return productReports, nil
}

func (s *ProductRecordsService) CreateProductRecord(newProductRecord internal.ProductRecords) (internal.ProductRecords, error) {
	err := s.validateEmptyFields(newProductRecord)
	if err != nil {
		return internal.ProductRecords{}, err
	}

	return s.repo.Create(newProductRecord)
}

func (s *ProductRecordsService) validateEmptyFields(newProduct internal.ProductRecords) error {
	if newProduct.LastUpdateDate == "" {
		return utils.ErrInvalidArguments
	}

	if newProduct.PurchasePrice == 0 {
		return utils.ErrInvalidArguments
	}

	if newProduct.SalePrice == 0 {
		return utils.ErrInvalidArguments
	}

	if newProduct.ProductID == 0 {
		return utils.ErrInvalidArguments
	}

	if _, err := s.validationProduct.GetProductByID(newProduct.ProductID); err != nil {
		return utils.ErrConflict
	}

	return nil
}
