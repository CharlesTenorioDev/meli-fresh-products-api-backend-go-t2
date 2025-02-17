package product_batch

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/testutils"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
)

func TestUnitProductBatchRepository(t *testing.T) {

	t.Run("Save", func(t *testing.T) {
		t.Run("should save the product batch successfully", func(t *testing.T) {
			newBatch := internal.ProductBatchRequest{
				BatchNumber:        100,
				CurrentQuantity:    50,
				CurrentTemperature: 22.4,
				DueDate:            "2022-01-01",
				InitialQuantity:    10,
				ManufacturingDate:  "2022-01-01",
				ManufacturingHour:  18,
				MinimumTemperature: -3,
				ProductID:          1,
				SectionID:          1,
			}

			createdBatch := internal.ProductBatch{
				ID:                  1,
				ProductBatchRequest: newBatch,
			}

			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			repository := NewProductBatchRepository(db)

			mock.ExpectExec(regexp.QuoteMeta("INSERT INTO product_batches (batch_number, current_quantity, current_temperature, due_date, initial_quantity, manufacturing_date, manufacturing_hour, minimum_temperature, product_id, section_id) VALUES (?,?,?,?,?,?,?,?,?,?)")).
				WithArgs(newBatch.BatchNumber, newBatch.CurrentQuantity, newBatch.CurrentTemperature, newBatch.DueDate, newBatch.InitialQuantity, newBatch.ManufacturingDate, newBatch.ManufacturingHour, newBatch.MinimumTemperature, newBatch.ProductID, newBatch.SectionID).
				WillReturnResult(sqlmock.NewResult(1, 1))

			result, err := repository.Save(&newBatch)
			require.NoError(t, err)
			require.NoError(t, mock.ExpectationsWereMet())
			require.Equal(t, result, createdBatch)

		})

		t.Run("should return a mySQL conflict error", func(t *testing.T) {
			newBatch := internal.ProductBatchRequest{
				BatchNumber:        100,
				CurrentQuantity:    50,
				CurrentTemperature: 22.4,
				DueDate:            "2022-01-01",
				InitialQuantity:    10,
				ManufacturingDate:  "2022-01-01",
				ManufacturingHour:  18,
				MinimumTemperature: -3,
				ProductID:          1,
				SectionID:          1,
			}

			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			repository := NewProductBatchRepository(db)
			mysqlErr := mysql.MySQLError{1062, [5]byte{0, 1, 2, 3, 4}, ""}

			mock.ExpectExec(regexp.QuoteMeta("INSERT INTO product_batches (batch_number, current_quantity, current_temperature, due_date, initial_quantity, manufacturing_date, manufacturing_hour, minimum_temperature, product_id, section_id) VALUES (?,?,?,?,?,?,?,?,?,?)")).
				WithArgs(newBatch.BatchNumber, newBatch.CurrentQuantity, newBatch.CurrentTemperature, newBatch.DueDate, newBatch.InitialQuantity, newBatch.ManufacturingDate, newBatch.ManufacturingHour, newBatch.MinimumTemperature, newBatch.ProductID, newBatch.SectionID).
				WillReturnError(&mysqlErr)

			result, err := repository.Save(&newBatch)
			require.ErrorIs(t, err, utils.ErrConflict)
			require.Equal(t, result, internal.ProductBatch{})

		})

		t.Run("should return any mySQL error except conflict error", func(t *testing.T) {

			newBatch := internal.ProductBatchRequest{
				BatchNumber:        100,
				CurrentQuantity:    50,
				CurrentTemperature: 22.4,
				DueDate:            "2022-01-01",
				InitialQuantity:    10,
				ManufacturingDate:  "2022-01-01",
				ManufacturingHour:  18,
				MinimumTemperature: -3,
				ProductID:          1,
				SectionID:          1,
			}

			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			repository := NewProductBatchRepository(db)
			mysqlErr := mysql.MySQLError{1099, [5]byte{0, 1, 2, 3, 4}, ""}

			mock.ExpectExec(regexp.QuoteMeta("INSERT INTO product_batches (batch_number, current_quantity, current_temperature, due_date, initial_quantity, manufacturing_date, manufacturing_hour, minimum_temperature, product_id, section_id) VALUES (?,?,?,?,?,?,?,?,?,?)")).
				WithArgs(newBatch.BatchNumber, newBatch.CurrentQuantity, newBatch.CurrentTemperature, newBatch.DueDate, newBatch.InitialQuantity, newBatch.ManufacturingDate, newBatch.ManufacturingHour, newBatch.MinimumTemperature, newBatch.ProductID, newBatch.SectionID).
				WillReturnError(&mysqlErr)

			result, err := repository.Save(&newBatch)
			require.ErrorIs(t, err, &mysqlErr)
			require.Equal(t, result, internal.ProductBatch{})

		})

		t.Run("should return an internal error", func(t *testing.T) {
			newBatch := internal.ProductBatchRequest{
				BatchNumber:        100,
				CurrentQuantity:    50,
				CurrentTemperature: 22.4,
				DueDate:            "2022-01-01",
				InitialQuantity:    10,
				ManufacturingDate:  "2022-01-01",
				ManufacturingHour:  18,
				MinimumTemperature: -3,
				ProductID:          1,
				SectionID:          1,
			}

			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			repository := NewProductBatchRepository(db)
			internalError := errors.New("internal error")

			mock.ExpectExec(regexp.QuoteMeta("INSERT INTO product_batches (batch_number, current_quantity, current_temperature, due_date, initial_quantity, manufacturing_date, manufacturing_hour, minimum_temperature, product_id, section_id) VALUES (?,?,?,?,?,?,?,?,?,?)")).
				WithArgs(newBatch.BatchNumber, newBatch.CurrentQuantity, newBatch.CurrentTemperature, newBatch.DueDate, newBatch.InitialQuantity, newBatch.ManufacturingDate, newBatch.ManufacturingHour, newBatch.MinimumTemperature, newBatch.ProductID, newBatch.SectionID).
				WillReturnError(internalError)

			result, err := repository.Save(&newBatch)
			require.ErrorIs(t, err, internalError)
			require.Equal(t, result, internal.ProductBatch{})
		})

	})

	t.Run("GetBatchNumber", func(t *testing.T) {
		t.Run("should return the batch number successfully", func(t *testing.T) {
			batchNumber := 100

			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			repository := NewProductBatchRepository(db)

			rows := sqlmock.NewRows([]string{"batch_number"}).
				AddRow(100)

			mock.ExpectQuery(regexp.QuoteMeta("SELECT batch_number FROM product_batches WHERE batch_number=?")).
				WithArgs(100).
				WillReturnRows(rows)

			result, err := repository.GetBatchNumber(batchNumber)
			require.NoError(t, err)
			require.Equal(t, batchNumber, result)

		})

		t.Run("should return not found error while getting the batch number", func(t *testing.T) {
			batchNumber := 100

			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			repository := NewProductBatchRepository(db)

			mock.ExpectQuery(regexp.QuoteMeta("SELECT batch_number FROM product_batches WHERE batch_number=?")).
				WithArgs(100).
				WillReturnError(sql.ErrNoRows)

			result, err := repository.GetBatchNumber(batchNumber)
			require.ErrorIs(t, err, nil)
			require.Equal(t, 0, result)

		})

		t.Run("should return an internal error while getting the batch number", func(t *testing.T) {
			batchNumber := 100

			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			repository := NewProductBatchRepository(db)
			internalError := errors.New("internal error")

			mock.ExpectQuery(regexp.QuoteMeta("SELECT batch_number FROM product_batches WHERE batch_number=?")).
				WithArgs(100).
				WillReturnError(internalError)

			result, err := repository.GetBatchNumber(batchNumber)
			require.ErrorIs(t, err, internalError)
			require.Equal(t, 0, result)

		})

	})

}

func TestIntegrationProductBatchRepository(t *testing.T) {
	ts := testutils.GetTestDBConn()
	defer ts.End()

	repository := NewProductBatchRepository(ts.DB)

	t.Run("Save", func(t *testing.T) {
		t.Run("should save a product batch successfully", func(t *testing.T) {
			newBatch := internal.ProductBatchRequest{
				BatchNumber:        100,
				CurrentQuantity:    50,
				CurrentTemperature: 22.4,
				DueDate:            "2022-01-01",
				InitialQuantity:    10,
				ManufacturingDate:  "2022-01-01",
				ManufacturingHour:  18,
				MinimumTemperature: -3,
				ProductID:          1,
				SectionID:          1,
			}

			createdBatch := internal.ProductBatch{
				ID:                  1,
				ProductBatchRequest: newBatch,
			}

			batch, err := repository.Save(&newBatch)
			require.NoError(t, err)
			require.Equal(t, createdBatch, batch)
		})
	})

}
