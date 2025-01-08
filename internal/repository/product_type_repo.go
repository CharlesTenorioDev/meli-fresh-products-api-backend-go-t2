package repository

import (
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type ProductTypeDB struct {
	db map[int]internal.ProductType
}

func NewProductTypeDB(db map[int]internal.ProductType) *ProductTypeDB {
	// default db
	defaultDb := make(map[int]internal.ProductType)
	if db != nil {
		defaultDb = db
	}
	return &ProductTypeDB{db: defaultDb}
}

// GetAll returns all product types
func (p *ProductTypeDB) GetAll() (listProductTypes []internal.ProductType, err error) {
	for _, productType := range p.db {
		listProductTypes = append(listProductTypes, productType)
	}
	return listProductTypes, nil
}

// GetByID returns a product type by id
func (p *ProductTypeDB) GetByID(id int) (productType internal.ProductType, err error) {
	productType, ok := p.db[id]
	if !ok {
		return internal.ProductType{}, utils.ErrNotFound
	}
	return productType, nil
}

// Create a product type
func (p *ProductTypeDB) Create(newProductType internal.ProductType) (productType internal.ProductType, err error) {
	newID := utils.GetBiggestId(p.db) + 1
	productType = internal.ProductType{
		ID:          newID,
		Description: newProductType.Description,
	}
	p.db[productType.ID] = productType
	return productType, nil
}

// Update a product type
func (p *ProductTypeDB) Update(inputProductType internal.ProductType) (productType internal.ProductType, err error) {
	p.db[inputProductType.ID] = inputProductType
	return inputProductType, nil
}

// Delete a product type
func (p *ProductTypeDB) Delete(id int) error {
	_, ok := p.db[id]
	if !ok {
		return utils.ErrNotFound
	}
	delete(p.db, id)
	return nil
}
