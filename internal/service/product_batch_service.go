package service

import (
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type ProductBatchService struct {
	batchRepo   internal.ProductBatchRepository
	productRepo internal.ProductRepository
	sectionRepo internal.SectionRepository
}

func NewProductBatchesService(batch internal.ProductBatchRepository,
	product internal.ProductRepository, section internal.SectionRepository) *ProductBatchService {
	return &ProductBatchService{
		batchRepo:   batch,
		productRepo: product,
		sectionRepo: section,
	}
}

func (s *ProductBatchService) Save(newBatch internal.ProductBatchRequest) (internal.ProductBatch, error) {
	batchValidation := s.verify(newBatch)

	if batchValidation != nil {
		return internal.ProductBatch{}, batchValidation
	}

	createdBatch, err := s.batchRepo.Save(&newBatch)
	if err != nil {
		return internal.ProductBatch{}, err
	}

	return createdBatch, nil

}

func (s *ProductBatchService) verify(newBatch internal.ProductBatchRequest) error {

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
	if newBatch.ProductId <= 0 {
		return utils.ErrInvalidArguments
	}
	if newBatch.SectionId <= 0 {
		return utils.ErrInvalidArguments
	}

	batchExists := s.batchRepo.GetBatchNumber(newBatch.BatchNumber)
	if batchExists != nil {
		return batchExists
	}

	sectionExists, err := s.sectionRepo.GetById(newBatch.SectionId)

	if sectionExists == (internal.Section{}) {
		return utils.ErrConflict
	}
	if err != nil {
		return err
	}

	productExists, err := s.productRepo.GetByID(newBatch.ProductId)
	if productExists == (internal.Product{}) {
		return utils.ErrConflict
	}
	if err != nil {
		return err
	}

	return nil
}
