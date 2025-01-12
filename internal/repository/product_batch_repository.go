package repository

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type ProductBatchRepository struct {
	db *sql.DB
}

func NewProductBatchRepository(db *sql.DB) *ProductBatchRepository {
	return &ProductBatchRepository{db: db}
}

// Get all the sections and return in asc order
func (r *ProductBatchRepository) Save(newBatch *internal.ProductBatchRequest) (internal.ProductBatch, error) {
	result, err := r.db.Exec("INSERT INTO product_batches (batch_number, current_quantity, current_temperature, due_date, initial_quantity, manufacturing_date, manufacturing_hour, minimum_temperature, product_id, section_id) VALUES (?,?,?,?,?,?,?,?,?,?)",
		(*newBatch).BatchNumber, (*newBatch).CurrentQuantity, (*newBatch).CurrentTemperature, (*newBatch).DueDate, (*newBatch).InitialQuantity, (*newBatch).ManufacturingDate, (*newBatch).ManufacturingHour, (*newBatch).MinimumTemperature, (*newBatch).ProductId, (*newBatch).SectionId,
	)
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 {
				err = utils.ErrConflict
			}
			return internal.ProductBatch{}, err
		}
		return internal.ProductBatch{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return internal.ProductBatch{}, err
	}
	createdBatch := internal.ProductBatch{
		ID:                  int(id),
		ProductBatchRequest: *newBatch,
	}

	return createdBatch, nil
}

func (r *SectionMysqlRepository) GetBatchNumber(batchNumber int) error {
	var exists bool

	row := r.db.QueryRow("SELECT batch_number FROM product_batches WHERE batch_number=?", batchNumber)

	err := row.Scan(&exists)

	if err != nil {
		if err == sql.ErrNoRows {
			err = utils.ErrNotFound
			return nil
		}
		return err
	}

	return utils.ErrConflict

}
