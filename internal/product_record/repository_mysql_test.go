package product_record_test

import (
	"database/sql"
	"errors"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/product_record"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestRead(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := product_record.NewProductRecordDB(db)

	rows := sqlmock.NewRows([]string{"product_id", "description", "records_count"}).
		AddRow(1, "Product 1", 10)

	mock.ExpectQuery("SELECT .* FROM products").WillReturnRows(rows)

	products, err := repo.Read(0)
	assert.NoError(t, err)
	assert.Len(t, products, 1)
	assert.Equal(t, 1, products[0].ProductID)
	assert.Equal(t, "Product 1", products[0].Description)
	assert.Equal(t, 10, products[0].RecordsCount)
}

func TestRead_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := product_record.NewProductRecordDB(db)
	mock.ExpectQuery("SELECT .* FROM products").WillReturnError(errors.New("query error"))

	_, err = repo.Read(0)
	assert.Error(t, err)
}

func TestRead_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := product_record.NewProductRecordDB(db)

	rows := sqlmock.NewRows([]string{"product_id", "description"}).AddRow(1, "Product 1")
	mock.ExpectQuery("SELECT .* FROM products").WillReturnRows(rows)

	_, err = repo.Read(0)
	assert.Error(t, err)
}

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := product_record.NewProductRecordDB(db)

	mock.ExpectPrepare("INSERT INTO product_records").ExpectExec().
		WithArgs("2025-02-17", 100.0, 150.0, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	productRecord := internal.ProductRecords{
		LastUpdateDate: "2025-02-17",
		PurchasePrice:  100.0,
		SalePrice:      150.0,
		ProductID:      1,
	}

	created, err := repo.Create(productRecord)
	assert.NoError(t, err)
	assert.Equal(t, 1, created.ID)
}

func TestCreate_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := product_record.NewProductRecordDB(db)
	mock.ExpectPrepare("INSERT INTO product_records").WillReturnError(errors.New("prepare error"))

	_, err = repo.Create(internal.ProductRecords{})
	assert.Error(t, err)
}

func TestCreate_ExecError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := product_record.NewProductRecordDB(db)

	mock.ExpectPrepare("INSERT INTO product_records").ExpectExec().WillReturnError(errors.New("exec error"))

	_, err = repo.Create(internal.ProductRecords{})
	assert.Error(t, err)
}

func TestCreate_ConflictError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := product_record.NewProductRecordDB(db)

	mock.ExpectPrepare("INSERT INTO product_records").ExpectExec().
		WillReturnError(&mysql.MySQLError{Number: 1062, Message: "Duplicate entry"})

	_, err = repo.Create(internal.ProductRecords{})
	assert.Equal(t, utils.ErrConflict, err)
}

func TestFindByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := product_record.NewProductRecordDB(db)
	row := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery("SELECT `id` FROM product_records WHERE id = ?").
		WithArgs(1).WillReturnRows(row)

	productRecord, err := repo.FindByID(1)
	assert.NoError(t, err)
	assert.Equal(t, 1, productRecord.ID)
}

func TestFindByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := product_record.NewProductRecordDB(db)
	mock.ExpectQuery("SELECT `id` FROM product_records WHERE id = ?").
		WithArgs(1).WillReturnError(sql.ErrNoRows)

	_, err = repo.FindByID(1)
	assert.Error(t, err)
}

func TestFindByID_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := product_record.NewProductRecordDB(db)

	rows := sqlmock.NewRows([]string{"invalid_column"}).AddRow("invalid_data")
	mock.ExpectQuery("SELECT `id` FROM product_records WHERE id = ?").WithArgs(1).WillReturnRows(rows)

	_, err = repo.FindByID(1)
	assert.Error(t, err)
}

func TestRead_WithProductID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := product_record.NewProductRecordDB(db)

	rows := sqlmock.NewRows([]string{"product_id", "description", "records_count"}).
		AddRow(1, "Product 1", 10)

	mock.ExpectQuery("SELECT .* WHERE p.id = ?").WithArgs(1).WillReturnRows(rows)

	products, err := repo.Read(1)
	assert.NoError(t, err)
	assert.Len(t, products, 1)
	assert.Equal(t, 1, products[0].ProductID)
}

func TestRead_RowsError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := product_record.NewProductRecordDB(db)

	rows := sqlmock.NewRows([]string{"product_id", "description", "records_count"}).
		AddRow(1, "Product 1", 10).
		RowError(0, errors.New("rows error"))

	mock.ExpectQuery("SELECT .* FROM products").WillReturnRows(rows)

	_, err = repo.Read(0)
	assert.Error(t, err)
}

func TestCreate_LastInsertIdError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := product_record.NewProductRecordDB(db)

	mock.ExpectPrepare("INSERT INTO product_records").ExpectExec().
		WillReturnResult(sqlmock.NewErrorResult(errors.New("last insert id error")))

	_, err = repo.Create(internal.ProductRecords{})
	assert.Error(t, err)
}
