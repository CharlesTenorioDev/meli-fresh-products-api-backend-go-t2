package repository

import (
	"database/sql"
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
func (r *ProductBatchRepository) Save(p internal.ProductBatchRequest) (internal.ProductBatch, error) {

	return internal.ProductBatch{}, nil
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
