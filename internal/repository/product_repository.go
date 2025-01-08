package repository

import (
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type ProductDB struct {
	db map[int]internal.Product
}

func NewProductDB(db map[int]internal.Product) *ProductDB {
	// default db
	defaultDb := make(map[int]internal.Product)
	if db != nil {
		defaultDb = db
	}
	return &ProductDB{db: defaultDb}
}

// GetAll returns all products
func (p *ProductDB) GetAll() (listProducts []internal.Product, err error) {
	for _, product := range p.db {
		listProducts = append(listProducts, product)
	}
	return listProducts, nil
}

// GetByID returns a product by id
func (p *ProductDB) GetByID(id int) (product internal.Product, err error) {
	product, ok := p.db[id]
	if !ok {
		return internal.Product{}, utils.ErrNotFound
	}
	return product, nil
}

// Create a product
func (p *ProductDB) Create(newproduct internal.ProductAttributes) (product internal.Product, err error) {
	newID := utils.GetBiggestId(p.db) + 1
	product = internal.Product{
		ID:                newID,
		ProductAttributes: newproduct,
	}
	p.db[product.ID] = product
	return product, nil
}

// Update a product
func (p *ProductDB) Update(inputProduct internal.Product) (product internal.Product, err error) {
	p.db[inputProduct.ID] = inputProduct
	return inputProduct, nil
}

// Delete a product
func (p *ProductDB) Delete(id int) error {
	_, ok := p.db[id]
	if !ok {
		return utils.ErrNotFound
	}
	delete(p.db, id)
	return nil
}
