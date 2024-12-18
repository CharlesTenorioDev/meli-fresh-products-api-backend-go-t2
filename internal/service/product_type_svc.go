package service

import "github.com/meli-fresh-products-api-backend-go-t2/internal/pkg"

type ProductTypeSvcDefault struct {
	repo pkg.ProductTypeRepository
}

func NewProductTypeServiceDefault(repo pkg.ProductTypeRepository) *ProductTypeSvcDefault {
	return &ProductTypeSvcDefault{repo: repo}
}

func (s *ProductTypeSvcDefault) GetProductTypes() (listProductTypes []pkg.ProductType, err error) {
	return s.repo.GetAll()
}

func (s *ProductTypeSvcDefault) GetProductTypeByID(id int) (productType pkg.ProductType, err error) {
	return s.repo.GetByID(id)
}

func (s *ProductTypeSvcDefault) CreateProductType(newProductType pkg.ProductType) (productType pkg.ProductType, err error) {
	return s.repo.Create(newProductType)
}

func (s *ProductTypeSvcDefault) UpdateProductType(inputProductType pkg.ProductType) (productType pkg.ProductType, err error) {
	internalProductType, err := s.repo.GetByID(inputProductType.ID)
	if err != nil {
		return pkg.ProductType{}, err
	}
	return s.repo.Update(internalProductType)
}

func (s *ProductTypeSvcDefault) DeleteProductType(id int) (err error) {
	return s.repo.Delete(id)
}
