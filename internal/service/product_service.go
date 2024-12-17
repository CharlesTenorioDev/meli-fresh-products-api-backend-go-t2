package service

import (
	"github.com/meli-fresh-products-api-backend-go-t2/internal/pkg"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type ProductServiceDefault struct {
	repo pkg.ProductRepository
}

func NewProductServiceDefault(repo pkg.ProductRepository) *ProductServiceDefault {
	return &ProductServiceDefault{repo: repo}
}

func (s *ProductServiceDefault) GetProducts() (listProducts []pkg.Product, err error) {
	return s.repo.GetAll()
}

func (s *ProductServiceDefault) GetProductByID(id int) (product pkg.Product, err error) {
	return s.repo.GetByID(id)
}

func (s *ProductServiceDefault) CreateProduct(newProduct pkg.ProductAttributes) (product pkg.Product, err error) {
	listProducts, _ := s.repo.GetAll()
	err = validateDuplicates(listProducts, newProduct)
	if err != nil {
		return pkg.Product{}, err
	}
	return s.repo.Create(newProduct)
}

func (s *ProductServiceDefault) UpdateProduct(inputProduct pkg.Product) (product pkg.Product, err error) {
	internalProduct, err := s.repo.GetByID(inputProduct.ID)
	if err != nil {
		return pkg.Product{}, utils.ErrNotFound
	}
	preparedProduct := prepareProductUpdate(inputProduct, internalProduct)
	return s.repo.Update(preparedProduct)
}

func (s *ProductServiceDefault) DeleteProduct(id int) (err error) {
	return s.repo.Delete(id)
}

func validateDuplicates(listProducts []pkg.Product, newProduct pkg.ProductAttributes) error {
	for _, product := range listProducts {
		if product.ProductCode == newProduct.ProductCode {
			return utils.ErrConflict
		}
	}
	return nil
}

func prepareProductUpdate(inputProduct, internalProduct pkg.Product) (preparedProduct pkg.Product) {
	preparedProduct.ID = internalProduct.ID
	if inputProduct.ProductCode != "" {
		preparedProduct.ProductCode = inputProduct.ProductCode
	} else {
		preparedProduct.ProductCode = internalProduct.ProductCode
	}
	if inputProduct.Description != "" {
		preparedProduct.Description = inputProduct.Description
	} else {
		preparedProduct.Description = internalProduct.Description
	}
	if inputProduct.Width != 0 {
		preparedProduct.Width = inputProduct.Width
	} else {
		preparedProduct.Width = internalProduct.Width
	}
	if inputProduct.Height != 0 {
		preparedProduct.Height = inputProduct.Height
	} else {
		preparedProduct.Height = internalProduct.Height
	}
	if inputProduct.Length != 0 {
		preparedProduct.Length = inputProduct.Length
	} else {
		preparedProduct.Length = internalProduct.Length
	}
	if inputProduct.NetWeight != 0 {
		preparedProduct.NetWeight = inputProduct.NetWeight
	} else {
		preparedProduct.NetWeight = internalProduct.NetWeight
	}
	if inputProduct.ExpirationRate != 0 {
		preparedProduct.ExpirationRate = inputProduct.ExpirationRate
	} else {
		preparedProduct.ExpirationRate = internalProduct.ExpirationRate
	}
	if inputProduct.RecommendedFreezingTemperature != 0 {
		preparedProduct.RecommendedFreezingTemperature = inputProduct.RecommendedFreezingTemperature
	} else {
		preparedProduct.RecommendedFreezingTemperature = internalProduct.RecommendedFreezingTemperature
	}
	if inputProduct.FreezingRate != 0 {
		preparedProduct.FreezingRate = inputProduct.FreezingRate
	} else {
		preparedProduct.FreezingRate = internalProduct.FreezingRate
	}
	if inputProduct.ProuctType != "" {
		preparedProduct.ProuctType = inputProduct.ProuctType
	} else {
		preparedProduct.ProuctType = internalProduct.ProuctType
	}
	if inputProduct.SellerID != 0 {
		preparedProduct.SellerID = inputProduct.SellerID
	} else {
		preparedProduct.SellerID = internalProduct.SellerID
	}
	return preparedProduct
}
