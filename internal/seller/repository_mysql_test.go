package seller

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
)

func TestUnitSellerRepository_GetAll_Success(t *testing.T) {
	sellers := []internal.Seller{
		{1, 55, "Company", "Address", "1199999999", 1}}

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repository := NewSellerRepository(db)

	rows := sqlmock.NewRows([]string{"id", "cid", "company_name", "address", "telephone", "locality_id"}).
		AddRow(1, 55, "Company", "Address", "1199999999", 1)

	mock.ExpectQuery("SELECT `id`, `cid`, `company_name`, `address`, `telephone`, `locality_id` FROM `sellers`").
		WillReturnRows(rows)

	result, err := repository.GetAll()
	assert.NoError(t, err)
	assert.Equal(t, sellers, result)

}

func TestUnitSellerRepository_GetAll_QueryError(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repository := NewSellerRepository(db)

	queryError := errors.New("query error")

	mock.ExpectQuery("SELECT `id`, `cid`, `company_name`, `address`, `telephone`, `locality_id` FROM `sellers`").
		WillReturnError(queryError)

	_, err = repository.GetAll()
	require.ErrorIs(t, err, queryError)

}

func TestUnitSellerRepository_GetAll_RowError(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repository := NewSellerRepository(db)

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(1)

	mock.ExpectQuery("SELECT `id`, `cid`, `company_name`, `address`, `telephone`, `locality_id` FROM `sellers`").
		WillReturnRows(rows)

	_, err = repository.GetAll()
	require.Error(t, err)

}

func TestUnitSellerRepository_GetById_ErrNoRows(t *testing.T) {
	seller := internal.Seller{1, 55, "Company", "Address", "1199999999", 1}

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repository := NewSellerRepository(db)

	mock.ExpectQuery("SELECT `id`, `cid`, `company_name`, `address`, `telephone`, `locality_id` FROM `sellers` WHERE `id` = ?").
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)

	_, err = repository.GetByID(seller.ID)

	require.ErrorIs(t, err, utils.ErrNotFound)
}

func TestUnitSellerRepository_GetById_InternalError(t *testing.T) {
	seller := internal.Seller{1, 55, "Company", "Address", "1199999999", 1}

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repository := NewSellerRepository(db)
	internalError := errors.New("internal error")

	mock.ExpectQuery("SELECT `id`, `cid`, `company_name`, `address`, `telephone`, `locality_id` FROM `sellers` WHERE `id` = ?").
		WithArgs(1).
		WillReturnError(internalError)

	_, err = repository.GetByID(seller.ID)

	require.ErrorIs(t, err, internalError)
}

func TestUnitSellerRepository_GetById_Success(t *testing.T) {
	seller := internal.Seller{1, 55, "Company", "Address", "1199999999", 1}

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repository := NewSellerRepository(db)

	rows := sqlmock.NewRows([]string{"id", "cid", "company_name", "address", "telephone", "locality_id"}).
		AddRow(1, 55, "Company", "Address", "1199999999", 1)

	mock.ExpectQuery("SELECT `id`, `cid`, `company_name`, `address`, `telephone`, `locality_id` FROM `sellers` WHERE `id` = ?").
		WithArgs(1).
		WillReturnRows(rows)

	result, err := repository.GetByID(seller.ID)
	assert.NoError(t, err)
	assert.Equal(t, seller, result)

}

func TestUnitSellerRepository_GetByCid_ErrNoRows(t *testing.T) {
	seller := internal.Seller{1, 55, "Company", "Address", "1199999999", 1}

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repository := NewSellerRepository(db)

	mock.ExpectQuery("SELECT `id`, `cid`, `company_name`, `address`, `telephone`, `locality_id` FROM `sellers` WHERE `cid` = ?").
		WithArgs(55).
		WillReturnError(sql.ErrNoRows)

	_, err = repository.GetByCid(seller.Cid)
	require.ErrorIs(t, err, utils.ErrNotFound)

}

func TestUnitSellerRepository_GetByCid_InternalError(t *testing.T) {
	seller := internal.Seller{1, 55, "Company", "Address", "1199999999", 1}

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repository := NewSellerRepository(db)
	internalError := errors.New("internal error")

	mock.ExpectQuery("SELECT `id`, `cid`, `company_name`, `address`, `telephone`, `locality_id` FROM `sellers` WHERE `cid` = ?").
		WithArgs(55).
		WillReturnError(internalError)

	_, err = repository.GetByCid(seller.Cid)
	require.ErrorIs(t, err, internalError)

}

func TestUnitSellerRepository_GetByCid_Success(t *testing.T) {
	seller := internal.Seller{1, 55, "Company", "Address", "1199999999", 1}

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repository := NewSellerRepository(db)

	rows := sqlmock.NewRows([]string{"id", "cid", "company_name", "address", "telephone", "locality_id"}).
		AddRow(1, 55, "Company", "Address", "1199999999", 1)

	mock.ExpectQuery("SELECT `id`, `cid`, `company_name`, `address`, `telephone`, `locality_id` FROM `sellers` WHERE `cid` = ?").
		WithArgs(55).
		WillReturnRows(rows)

	result, err := repository.GetByCid(seller.Cid)
	assert.NoError(t, err)
	assert.Equal(t, seller, result)

}

func TestUnitSellerRepository_Create_Success(t *testing.T) {
	seller := internal.Seller{1, 55, "Company", "Address", "1199999999", 1}

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repository := NewSellerRepository(db)

	mock.ExpectExec("INSERT INTO `sellers`").
		WithArgs(seller.Cid, seller.CompanyName, seller.Address, seller.Telephone, seller.LocalityID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repository.Create(&seller)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

}

func TestUnitSellerRepository_Create_ErrorConflict(t *testing.T) {
	seller := internal.Seller{1, 55, "Company", "Address", "1199999999", 1}

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repository := NewSellerRepository(db)
	mysqlErr := mysql.MySQLError{1062, [5]byte{0, 1, 2, 3, 4}, ""}

	mock.ExpectExec("INSERT INTO `sellers`").
		WithArgs(seller.Cid, seller.CompanyName, seller.Address, seller.Telephone, seller.LocalityID).
		WillReturnError(&mysqlErr)

	err = repository.Create(&seller)
	require.ErrorIs(t, err, utils.ErrConflict)

}

func TestUnitSellerRepository_Create_InternalError(t *testing.T) {
	seller := internal.Seller{1, 55, "Company", "Address", "1199999999", 1}

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repository := NewSellerRepository(db)
	internalError := errors.New("internal error")

	mock.ExpectExec("INSERT INTO `sellers`").
		WithArgs(seller.Cid, seller.CompanyName, seller.Address, seller.Telephone, seller.LocalityID).
		WillReturnError(internalError)

	err = repository.Create(&seller)
	require.ErrorIs(t, err, internalError)

}

func TestUnitSellerRepository_Update_Success(t *testing.T) {
	seller := internal.Seller{1, 55, "Company", "Address", "1199999999", 1}

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repository := NewSellerRepository(db)

	mock.ExpectExec(regexp.QuoteMeta("UPDATE `sellers` SET `cid` = ?, `company_name` = ?, `address` = ?, `telephone` = ?, `locality_id` = ? WHERE `id` = ?")).
		WithArgs(seller.Cid, seller.CompanyName, seller.Address, seller.Telephone, seller.LocalityID, seller.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repository.Update(&seller)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

}

func TestUnitSellerRepository_Update_ErrorConflict(t *testing.T) {
	seller := internal.Seller{1, 55, "Company", "Address", "1199999999", 1}

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repository := NewSellerRepository(db)
	mysqlErr := mysql.MySQLError{1062, [5]byte{0, 1, 2, 3, 4}, ""}

	mock.ExpectExec(regexp.QuoteMeta("UPDATE `sellers` SET `cid` = ?, `company_name` = ?, `address` = ?, `telephone` = ?, `locality_id` = ? WHERE `id` = ?")).
		WithArgs(seller.Cid, seller.CompanyName, seller.Address, seller.Telephone, seller.LocalityID, seller.ID).
		WillReturnError(&mysqlErr)

	err = repository.Update(&seller)
	require.ErrorIs(t, err, utils.ErrConflict)

}

func TestUnitSellerRepository_Update_InternalError(t *testing.T) {
	seller := internal.Seller{1, 55, "Company", "Address", "1199999999", 1}

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repository := NewSellerRepository(db)
	internalError := errors.New("internal error")

	mock.ExpectExec(regexp.QuoteMeta("UPDATE `sellers` SET `cid` = ?, `company_name` = ?, `address` = ?, `telephone` = ?, `locality_id` = ? WHERE `id` = ?")).
		WithArgs(seller.Cid, seller.CompanyName, seller.Address, seller.Telephone, seller.LocalityID, seller.ID).
		WillReturnError(internalError)

	err = repository.Update(&seller)
	require.ErrorIs(t, err, internalError)

}

func TestUnitSellerRepository_Delete_Success(t *testing.T) {
	seller := internal.Seller{1, 55, "Company", "Address", "1199999999", 1}

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repository := NewSellerRepository(db)

	mock.ExpectExec("DELETE FROM `sellers` WHERE `id` = ?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repository.Delete(seller.ID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUnitSellerRepository_Delete_InternalError(t *testing.T) {
	seller := internal.Seller{1, 55, "Company", "Address", "1199999999", 1}

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repository := NewSellerRepository(db)
	internalError := errors.New("internal error")

	mock.ExpectExec("DELETE FROM `sellers` WHERE `id` = ?").
		WithArgs(1).
		WillReturnError(internalError)

	err = repository.Delete(seller.ID)
	require.ErrorIs(t, err, internalError)

}
