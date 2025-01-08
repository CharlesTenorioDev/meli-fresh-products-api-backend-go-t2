package service

import (
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
)

type ProductTypeSvc struct {
	repo internal.ProductTypeRepository
}

func NewProductTypeService(repo internal.ProductTypeRepository) *ProductTypeSvc {
	return &ProductTypeSvc{repo: repo}
}

func (s *ProductTypeSvc) GetProductTypes() (listProductTypes []internal.ProductType, err error) {
	return s.repo.GetAll()
}

func (s *ProductTypeSvc) GetProductTypeByID(id int) (productType internal.ProductType, err error) {
	return s.repo.GetByID(id)
}

func (s *ProductTypeSvc) CreateProductType(newProductType internal.ProductType) (productType internal.ProductType, err error) {
	return s.repo.Create(newProductType)
}

func (s *ProductTypeSvc) UpdateProductType(inputProductType internal.ProductType) (productType internal.ProductType, err error) {
	_, err = s.repo.GetByID(inputProductType.ID)
	if err != nil {
		return internal.ProductType{}, err
	}
	return s.repo.Update(inputProductType)
}

func (s *ProductTypeSvc) DeleteProductType(id int) (err error) {
	return s.repo.Delete(id)
}
