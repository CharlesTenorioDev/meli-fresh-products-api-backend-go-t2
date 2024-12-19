package service

import "github.com/meli-fresh-products-api-backend-go-t2/internal/pkg"

type ProductTypeSvc struct {
	repo pkg.ProductTypeRepository
}

func NewProductTypeService(repo pkg.ProductTypeRepository) *ProductTypeSvc {
	return &ProductTypeSvc{repo: repo}
}

func (s *ProductTypeSvc) GetProductTypes() (listProductTypes []pkg.ProductType, err error) {
	return s.repo.GetAll()
}

func (s *ProductTypeSvc) GetProductTypeByID(id int) (productType pkg.ProductType, err error) {
	return s.repo.GetByID(id)
}

func (s *ProductTypeSvc) CreateProductType(newProductType pkg.ProductType) (productType pkg.ProductType, err error) {
	return s.repo.Create(newProductType)
}

func (s *ProductTypeSvc) UpdateProductType(inputProductType pkg.ProductType) (productType pkg.ProductType, err error) {
	_, err = s.repo.GetByID(inputProductType.ID)
	if err != nil {
		return pkg.ProductType{}, err
	}
	return s.repo.Update(inputProductType)
}

func (s *ProductTypeSvc) DeleteProductType(id int) (err error) {
	return s.repo.Delete(id)
}
