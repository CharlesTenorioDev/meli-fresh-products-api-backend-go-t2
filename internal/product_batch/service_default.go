package product_batch

import (
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type DefaultProductBatchService struct {
	batchRepo   internal.ProductBatchRepository
	productRepo internal.ProductRepository
	sectionRepo internal.SectionRepository
}

func NewProductBatchesService(batch internal.ProductBatchRepository,
	product internal.ProductRepository, section internal.SectionRepository) internal.ProductBatchService {
	return &DefaultProductBatchService{
		batchRepo:   batch,
		productRepo: product,
		sectionRepo: section,
	}
}

func (s *DefaultProductBatchService) Save(newBatch *internal.ProductBatchRequest) (internal.ProductBatch, error) {
	batchValidation := s.verify(newBatch)

	if batchValidation != nil {
		return internal.ProductBatch{}, batchValidation
	}

	createdBatch, err := s.batchRepo.Save(newBatch)
	if err != nil {
		return internal.ProductBatch{}, err
	}

	return createdBatch, nil
}

func (s *DefaultProductBatchService) verify(newBatch *internal.ProductBatchRequest) error {
	if newBatch.BatchNumber <= 0 {
		return utils.ErrInvalidArguments
	}

	if newBatch.CurrentQuantity < 0 {
		return utils.ErrInvalidArguments
	}

	if newBatch.CurrentTemperature <= 0.0 {
		return utils.ErrInvalidArguments
	}

	if len(newBatch.DueDate) == 0 {
		return utils.ErrInvalidArguments
	}

	if newBatch.InitialQuantity < 0 {
		return utils.ErrInvalidArguments
	}

	if len(newBatch.ManufacturingDate) == 0 {
		return utils.ErrInvalidArguments
	}

	if newBatch.ManufacturingHour < 0 {
		return utils.ErrInvalidArguments
	}
	//MinimumTemperature não validada porque pode ser positiva, negativa ou zero
	if newBatch.ProductID <= 0 {
		return utils.ErrInvalidArguments
	}

	if newBatch.SectionID <= 0 {
		return utils.ErrInvalidArguments
	}

	batchExists, err := s.batchRepo.GetBatchNumber(newBatch.BatchNumber)
	if err != nil {
		return err
	}

	if batchExists != 0 {
		return utils.ErrConflict
	}

	sectionExists, err := s.sectionRepo.GetByID(newBatch.SectionID)

	if sectionExists == (internal.Section{}) {
		return utils.ErrConflict
	}

	if err != nil {
		return err
	}

	productExists, err := s.productRepo.GetByID(newBatch.ProductID)
	if productExists == (internal.Product{}) {
		return utils.ErrConflict
	}

	if err != nil {
		return err
	}

	return nil
}
