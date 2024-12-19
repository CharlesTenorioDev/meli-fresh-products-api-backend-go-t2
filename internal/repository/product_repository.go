package repository

import (
	"github.com/meli-fresh-products-api-backend-go-t2/internal/pkg"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type ProductDB struct {
	db map[int]pkg.Product
}

func NewProductDB(db map[int]pkg.Product) *ProductDB {
	// default db
	defaultDb := make(map[int]pkg.Product)
	if db != nil {
		defaultDb = db
	}
	return &ProductDB{db: defaultDb}
}

// GetAll returns all products
func (p *ProductDB) GetAll() (listProducts []pkg.Product, err error) {
	for _, product := range p.db {
		listProducts = append(listProducts, product)
	}
	return listProducts, nil
}

// GetByID returns a product by id
func (p *ProductDB) GetByID(id int) (product pkg.Product, err error) {
	product, ok := p.db[id]
	if !ok {
		return pkg.Product{}, utils.ErrNotFound
	}
	return product, nil
}

// Create a product
func (p *ProductDB) Create(newproduct pkg.ProductAttributes) (product pkg.Product, err error) {
	newID := utils.GetBiggestId(p.db) + 1
	product = pkg.Product{
		ID:                newID,
		ProductAttributes: newproduct,
	}
	p.db[product.ID] = product
	return product, nil
}

// Update a product
func (p *ProductDB) Update(inputProduct pkg.Product) (product pkg.Product, err error) {
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
