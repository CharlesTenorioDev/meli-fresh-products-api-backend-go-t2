package repository

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type ProductDB struct {
	db *sql.DB
}

func NewProductDB(db *sql.DB) *ProductDB {

	return &ProductDB{db: db}
}

// GetAll returns all products
func (p *ProductDB) GetAll() (listProducts []internal.Product, err error) {
	var productList []internal.Product

	rows, err := p.db.Query("SELECT id, description, expiration_rate, freezing_rate, height, length, net_weight, product_code, recommended_freezing_temperature, width, product_type_id, seller_id FROM fresh_products.products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var product internal.Product

		err := rows.Scan(&product.ID, &product.Description, &product.ExpirationRate, &product.FreezingRate, &product.Height, &product.Length, &product.NetWeight, &product.ProductCode, &product.RecommendedFreezingTemperature, &product.Width, &product.ProductType, &product.SellerID)
		if err != nil {
			return nil, err
		}

		productList = append(productList, product)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return productList, nil
}

// GetByID returns a product by id
func (p *ProductDB) GetByID(id int) (product internal.Product, err error) {
	row := p.db.QueryRow("SELECT id, description, expiration_rate, freezing_rate, height, length, net_weight, product_code, recommended_freezing_temperature, width, product_type_id, seller_id FROM products WHERE id = ?", id)
	if err := row.Err(); err != nil {
		return internal.Product{}, err
	}

	err = row.Scan(&product.ID, &product.Description, &product.ExpirationRate, &product.FreezingRate, &product.Height, &product.Length, &product.NetWeight, &product.ProductCode, &product.RecommendedFreezingTemperature, &product.Width, &product.ProductType, &product.SellerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return internal.Product{}, utils.ErrNotFound
		}
		return internal.Product{}, err
	}

	return product, nil

}

// Create a product
func (p *ProductDB) Create(newproduct internal.ProductAttributes) (product internal.Product, err error) {
	statement, err := p.db.Prepare("INSERT INTO products (description, expiration_rate, freezing_rate, height, `length`, net_weight, product_code, recommended_freezing_temperature, width, product_type_id, seller_id) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return internal.Product{}, err
	}
	defer statement.Close()

	result, err := statement.Exec(newproduct.Description, newproduct.ExpirationRate, newproduct.FreezingRate, newproduct.Height, newproduct.Length, newproduct.NetWeight, newproduct.ProductCode, newproduct.RecommendedFreezingTemperature, newproduct.Width, newproduct.ProductType, newproduct.SellerID)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			switch mysqlErr.Number {
			case 1062:
				err = utils.ErrConflict
				fallthrough
			default:
				return internal.Product{}, err
			}
		}
		return internal.Product{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return internal.Product{}, err
	}

	product = internal.Product{
		ID:                int(id),
		ProductAttributes: newproduct,
	}

	return product, nil
}

// Update a product
func (p *ProductDB) Update(inputProduct internal.Product) (product internal.Product, err error) {
	_, err = p.GetByID(inputProduct.ID)
	if err != nil {
		return internal.Product{}, err
	}

	statement, err := p.db.Prepare(
		"UPDATE products SET description=?, expiration_rate=?, freezing_rate=?, height=?, `length`=?, net_weight=?, product_code=?, recommended_freezing_temperature=?, width=?, product_type_id=?, seller_id=? WHERE id=?",
	)
	if err != nil {
		return internal.Product{}, err
	}
	defer statement.Close()

	_, err = statement.Exec(inputProduct.Description, inputProduct.ExpirationRate, inputProduct.FreezingRate, inputProduct.Height, inputProduct.Length, inputProduct.NetWeight, inputProduct.ProductCode, inputProduct.RecommendedFreezingTemperature, inputProduct.Width, inputProduct.ProductType, inputProduct.SellerID, inputProduct.ID)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			switch mysqlErr.Number {
			case 1062:
				err = utils.ErrConflict
				fallthrough
			default:
				return internal.Product{}, err
			}
		}
		return internal.Product{}, err
	}

	return inputProduct, nil
}

// Delete a product
func (p *ProductDB) Delete(id int) error {
	_, err := p.GetByID(id)

	if err != nil {
		return err
	}

	statement, err := p.db.Prepare("DELETE FROM products WHERE id = ?")

	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(id)

	if err != nil {
		return err
	}

	return nil
}
