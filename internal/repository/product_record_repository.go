package repository

import (
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type ProductRecordDB struct {
	db *sql.DB
}

func NewProductRecordDB(db *sql.DB) *ProductRecordDB {
	return &ProductRecordDB{db: db}
}

func (p *ProductRecordDB) Read(productID int) ([]internal.ProductReport, error) {
	var rows *sql.Rows

	var err error

	if productID > 0 {
		rows, err = p.db.Query(`
			SELECT 
				p.id AS product_id,
				p.description,
				COUNT(pr.id) AS records_count
			FROM 
				products p
			INNER JOIN 
				product_records pr ON p.id = pr.product_id
			WHERE 
				p.id = ?
			GROUP BY 
				p.id, p.description
		`, productID)
	} else {
		rows, err = p.db.Query(`
			SELECT 
				p.id AS product_id,
				p.description,
				COUNT(pr.id) AS records_count
			FROM 
				products p
			INNER JOIN 
				product_records pr ON p.id = pr.product_id
			GROUP BY 
				p.id, p.description
		`)
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var listProducts []internal.ProductReport

	for rows.Next() {
		var productRecord internal.ProductReport

		err := rows.Scan(&productRecord.ProductID, &productRecord.Description, &productRecord.RecordsCount)
		if err != nil {
			return nil, err
		}

		listProducts = append(listProducts, productRecord)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return listProducts, nil
}

func (p *ProductRecordDB) Create(newProductRecord internal.ProductRecords) (internal.ProductRecords, error) {
	statement, err := p.db.Prepare("INSERT INTO product_records (last_update_date, purchase_price, sale_price, product_id) VALUES(?, ?, ?, ?)")
	if err != nil {
		return internal.ProductRecords{}, err
	}
	defer statement.Close()

	result, err := statement.Exec(newProductRecord.LastUpdateDate, newProductRecord.PurchasePrice, newProductRecord.SalePrice, newProductRecord.ProductID)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			switch mysqlErr.Number {
			case 1062:
				err = utils.ErrConflict
				fallthrough
			default:
				return internal.ProductRecords{}, err
			}
		}

		return internal.ProductRecords{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return internal.ProductRecords{}, err
	}

	newProductRecord.ID = int(id)

	return newProductRecord, nil
}

func (p *ProductRecordDB) FindByID(productRecordID int) (internal.ProductRecords, error) {
	var row *sql.Row

	var err error

	query := "select `id` from product_records where id = ?"

	row = p.db.QueryRow(query, productRecordID)

	var pr internal.ProductRecords

	err = row.Scan(&pr.ID)
	if err != nil {
		return internal.ProductRecords{}, err
	}

	err = row.Err()
	if err != nil {
		return internal.ProductRecords{}, err
	}

	return pr, nil
}
