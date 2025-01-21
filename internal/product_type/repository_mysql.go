package product_type

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type ProductTypeDB struct {
	db *sql.DB
}

func NewProductTypeDB(db *sql.DB) *ProductTypeDB {
	return &ProductTypeDB{db: db}
}

// GetAll returns all product types
func (p *ProductTypeDB) GetAll() (listProductTypes []internal.ProductType, err error) {
	rows, err := p.db.Query("SELECT id, description FROM product_types")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var productType internal.ProductType

		err := rows.Scan(&productType.ID, &productType.Description)
		if err != nil {
			return nil, err
		}

		listProductTypes = append(listProductTypes, productType)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return listProductTypes, nil
}

// GetByID returns a product type by id
func (p *ProductTypeDB) GetByID(id int) (productType internal.ProductType, err error) {
	row := p.db.QueryRow("SELECT id, description FROM product_types WHERE id = ?", id)
	if err := row.Err(); err != nil {
		return internal.ProductType{}, err
	}

	err = row.Scan(&productType.ID, &productType.Description)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return internal.ProductType{}, utils.ErrNotFound
		}

		return internal.ProductType{}, err
	}

	return productType, nil
}

// Create a product type
func (p *ProductTypeDB) Create(newProductType internal.ProductType) (productType internal.ProductType, err error) {
	statement, err := p.db.Prepare("INSERT INTO product_types (description) VALUES(?)")
	if err != nil {
		return internal.ProductType{}, err
	}
	defer statement.Close()

	result, err := statement.Exec(newProductType.Description)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			switch mysqlErr.Number {
			case 1062:
				err = utils.ErrConflict
				fallthrough
			default:
				return internal.ProductType{}, err
			}
		}

		return internal.ProductType{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return internal.ProductType{}, err
	}

	newProductType.ID = int(id)

	return newProductType, nil
}

// Update a product type
func (p *ProductTypeDB) Update(inputProductType internal.ProductType) (productType internal.ProductType, err error) {
	_, err = p.GetByID(inputProductType.ID)
	if err != nil {
		return internal.ProductType{}, err
	}

	statement, err := p.db.Prepare(
		"UPDATE product_types SET description=? WHERE id=?",
	)
	if err != nil {
		return internal.ProductType{}, err
	}
	defer statement.Close()

	_, err = statement.Exec(inputProductType.Description, inputProductType.ID)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			switch mysqlErr.Number {
			case 1062:
				err = utils.ErrConflict
				fallthrough
			default:
				return internal.ProductType{}, err
			}
		}

		return internal.ProductType{}, err
	}

	return inputProductType, nil
}

// Delete a product type
func (p *ProductTypeDB) Delete(id int) error {
	_, err := p.GetByID(id)
	if err != nil {
		return err
	}

	statement, err := p.db.Prepare("DELETE FROM product_types WHERE id = ?")
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
