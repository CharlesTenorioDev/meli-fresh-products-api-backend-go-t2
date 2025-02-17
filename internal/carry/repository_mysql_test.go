package carry

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
	"github.com/stretchr/testify/require"
)

func TestMySQLCarryRepository_Save(t *testing.T) {
	t.Run("GIVEN a valid carry struct WHEN r.db.Prepare THEN return an error", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()
		repo := NewMySQLCarryRepository(db)

		// conditions
		wErr := errors.New("some internal db error")
		mock.ExpectPrepare("INSERT INTO carriers").WillReturnError(wErr)

		err := repo.Save(&internal.Carry{
			CID:         123,
			CompanyName: "Test Company",
			Address:     "123 Test St",
			Telephone:   "1234567890",
			LocalityID:  1,
		})

		require.ErrorIs(t, err, wErr)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("GIVEN a valid carry struct WHEN stmt.Exec returns a conflict error THEN return utils.EConflict", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()
		repo := NewMySQLCarryRepository(db)

		// conditions
		stmt := mock.ExpectPrepare("INSERT INTO carriers")
		stmt.WillBeClosed()
		exc := stmt.ExpectExec()
		exc.WillReturnError(&mysql.MySQLError{Number: 1062, Message: "Duplicate entry"})

		err := repo.Save(&internal.Carry{
			CID:         123,
			CompanyName: "Test Company",
			Address:     "123 Test St",
			Telephone:   "1234567890",
			LocalityID:  1,
		})

		require.Error(t, err)
	})

	t.Run("GIVEN a valid carry struct WHEN stmt.Exec returns an error THEN return the same error", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()
		repo := NewMySQLCarryRepository(db)

		// conditions
		wErr := errors.New("some internal db error")
		stmt := mock.ExpectPrepare("INSERT INTO carriers")
		stmt.WillBeClosed()
		exc := stmt.ExpectExec()
		exc.WillReturnError(wErr)

		err := repo.Save(&internal.Carry{
			CID:         123,
			CompanyName: "Test Company",
			Address:     "123 Test St",
			Telephone:   "1234567890",
			LocalityID:  1,
		})

		require.ErrorIs(t, err, wErr)

	})

	t.Run("GIVEN a valid carry struct WHEN result.LastInsertId returns an error THEN return the same error", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()
		repo := NewMySQLCarryRepository(db)

		// conditions
		stmt := mock.ExpectPrepare("INSERT INTO carriers")
		stmt.WillBeClosed()
		exc := stmt.ExpectExec()
		exc.WillReturnResult(sqlmock.NewErrorResult(errors.New("some internal db error")))

		err := repo.Save(&internal.Carry{
			CID:         123,
			CompanyName: "Test Company",
			Address:     "123 Test St",
			Telephone:   "1234567890",
			LocalityID:  1,
		})

		require.Error(t, err)

	})

	t.Run("GIVEN a valid carry struct WHEN everything is OK THEN save the carry", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()
		repo := NewMySQLCarryRepository(db)

		// conditions
		stmt := mock.ExpectPrepare("INSERT INTO carriers")
		stmt.WillBeClosed()
		exc := stmt.ExpectExec()
		exc.WillReturnResult(sqlmock.NewResult(1, 1))

		c := &internal.Carry{
			CID:         123,
			CompanyName: "Test Company",
			Address:     "123 Test St",
			Telephone:   "1234567890",
			LocalityID:  1,
		}

		err := repo.Save(c)

		require.NoError(t, err)
		require.Equal(t, 1, c.ID)

	})
}
func TestMySQLCarryRepository_GetAll(t *testing.T) {
	t.Run("GIVEN a database error WHEN r.db.Query THEN return utils.ENotFound", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()
		repo := NewMySQLCarryRepository(db)

		// conditions
		wErr := errors.New("some internal db error")
		mock.ExpectQuery("SELECT id, cid, company_name, address, telephone, locality_id FROM carriers;").WillReturnError(wErr)

		_, err := repo.GetAll()

		require.Error(t, err)
		require.Equal(t, utils.ENotFound("Carry"), err)
	})

	t.Run("GIVEN a valid query WHEN everything is OK THEN return all carries", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()
		repo := NewMySQLCarryRepository(db)

		// conditions
		rows := sqlmock.NewRows([]string{"id", "cid", "company_name", "address", "telephone", "locality_id"}).
			AddRow(1, 123, "Test Company", "123 Test St", "1234567890", 1).
			AddRow(2, 124, "Another Company", "456 Another St", "0987654321", 2)
		mock.ExpectQuery("SELECT id, cid, company_name, address, telephone, locality_id FROM carriers;").WillReturnRows(rows)

		carries, err := repo.GetAll()

		require.NoError(t, err)
		require.Len(t, carries, 2)

	})
}
func TestMySQLCarryRepository_GetByID(t *testing.T) {
	t.Run("GIVEN a database error WHEN r.db.Prepare THEN return the same error", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()
		repo := NewMySQLCarryRepository(db)

		// conditions
		wErr := errors.New("some internal db error")
		mock.ExpectPrepare("SELECT id, cid, company_name, address, telephone, locality_id FROM carriers WHERE id=?;").WillReturnError(wErr)

		_, err := repo.GetByID(1)
		require.Error(t, err)
	})

	t.Run("GIVEN a valid id WHEN row.Scan returns sql.ErrNoRows THEN return utils.ENotFound", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()
		repo := NewMySQLCarryRepository(db)

		// Define the expected query without the semicolon at the end
		expectedQuery := "SELECT id, cid, company_name, address, telephone, locality_id FROM carriers WHERE id=?"

		// Expect the prepare statement and the query execution
		mock.ExpectPrepare(expectedQuery).
			ExpectQuery().
			WithArgs(1).
			WillReturnError(sql.ErrNoRows) // Simulate no rows found

		_, err := repo.GetByID(1)

		require.Error(t, err)
		require.Equal(t, utils.ENotFound("Carry"), err)
	})

	t.Run("GIVEN a valid id WHEN row.Scan returns an error THEN return the same error", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()
		repo := NewMySQLCarryRepository(db)

		// conditions
		// Define the expected query without the semicolon at the end
		expectedQuery := "SELECT id, cid, company_name, address, telephone, locality_id FROM carriers WHERE id=?"

		// Expect the prepare statement and the query execution
		mock.ExpectPrepare(expectedQuery).
			ExpectQuery().
			WithArgs(1).
			WillReturnError(errors.ErrUnsupported)

		_, err := repo.GetByID(1)

		require.Error(t, err)
		require.ErrorIs(t, err, errors.ErrUnsupported)
	})

	t.Run("GIVEN a valid id WHEN everything is OK THEN return the carry", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()
		repo := NewMySQLCarryRepository(db)

		// Define the expected query without the semicolon at the end
		expectedQuery := "SELECT id, cid, company_name, address, telephone, locality_id FROM carriers WHERE id=?"

		// Expect the prepare statement and the query execution
		rows := sqlmock.NewRows([]string{"id", "cid", "company_name", "address", "telephone", "locality_id"}).
			AddRow(1, 123, "Test Company", "123 Test St", "1234567890", 1)
		mock.ExpectPrepare(expectedQuery).
			ExpectQuery().
			WithArgs(1).
			WillReturnRows(rows)

		carry, err := repo.GetByID(1)

		expectedCarry := internal.Carry{
			ID:          1,
			CID:         123,
			CompanyName: "Test Company",
			Address:     "123 Test St",
			Telephone:   "1234567890",
			LocalityID:  1,
		}

		require.NoError(t, err)
		require.Equal(t, expectedCarry, carry)
	})
}

func TestMySQLCarryRepository_Update(t *testing.T) {
	t.Run("GIVEN a valid carry struct WHEN r.db.Prepare THEN return an error", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()
		repo := NewMySQLCarryRepository(db)

		// conditions
		wErr := errors.New("some internal db error")
		mock.ExpectPrepare("UPDATE carriers SET").WillReturnError(wErr)

		err := repo.Update(&internal.Carry{
			ID:          1,
			CID:         123,
			CompanyName: "Test Company",
			Address:     "123 Test St",
			Telephone:   "1234567890",
			LocalityID:  1,
		})

		require.ErrorIs(t, err, wErr)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("GIVEN a valid carry struct WHEN stmt.Exec returns a conflict error THEN return utils.EConflict", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()
		repo := NewMySQLCarryRepository(db)

		// conditions
		stmt := mock.ExpectPrepare("UPDATE carriers SET")
		stmt.WillBeClosed()
		exc := stmt.ExpectExec()
		exc.WillReturnError(&mysql.MySQLError{Number: 1062, Message: "Duplicate entry"})

		err := repo.Update(&internal.Carry{
			ID:          1,
			CID:         123,
			CompanyName: "Test Company",
			Address:     "123 Test St",
			Telephone:   "1234567890",
			LocalityID:  1,
		})

		require.Error(t, err)
	})

	t.Run("GIVEN a valid carry struct WHEN stmt.Exec returns an error THEN return the same error", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()
		repo := NewMySQLCarryRepository(db)

		// conditions
		wErr := errors.New("some internal db error")
		stmt := mock.ExpectPrepare("UPDATE carriers SET")
		stmt.WillBeClosed()
		exc := stmt.ExpectExec()
		exc.WillReturnError(wErr)

		err := repo.Update(&internal.Carry{
			ID:          1,
			CID:         123,
			CompanyName: "Test Company",
			Address:     "123 Test St",
			Telephone:   "1234567890",
			LocalityID:  1,
		})

		require.ErrorIs(t, err, wErr)
	})

	t.Run("GIVEN a valid carry struct return a nil error", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()
		repo := NewMySQLCarryRepository(db)

		// conditions
		stmt := mock.ExpectPrepare("UPDATE carriers SET")
		stmt.WillBeClosed()
		exc := stmt.ExpectExec()
		exc.WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Update(&internal.Carry{
			ID:          1,
			CID:         123,
			CompanyName: "Test Company",
			Address:     "123 Test St",
			Telephone:   "1234567890",
			LocalityID:  1,
		})

		require.NoError(t, err)
	})
}
func TestMySQLCarryRepository_Delete(t *testing.T) {
	t.Run("GIVEN a valid id WHEN r.GetByID returns an error THEN return the same error", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()
		repo := NewMySQLCarryRepository(db)

		// conditions
		wErr := errors.New("some internal db error")
		mock.ExpectPrepare("SELECT id, cid, company_name, address, telephone, locality_id FROM carriers WHERE id=?").
			ExpectQuery().
			WithArgs(1).
			WillReturnError(wErr)

		err := repo.Delete(1)
		require.ErrorIs(t, err, wErr)
	})

	t.Run("GIVEN a valid id WHEN r.db.Prepare returns an error THEN return the same error", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()
		repo := NewMySQLCarryRepository(db)

		// conditions
		mock.ExpectPrepare("SELECT id, cid, company_name, address, telephone, locality_id FROM carriers WHERE id=?").
			ExpectQuery().
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "cid", "company_name", "address", "telephone", "locality_id"}).
				AddRow(1, 123, "Test Company", "123 Test St", "1234567890", 1))

		wErr := errors.New("some internal db error")
		mock.ExpectPrepare("DELETE FROM carriers WHERE id=?").WillReturnError(wErr)

		err := repo.Delete(1)
		require.ErrorIs(t, err, wErr)
	})

	t.Run("GIVEN a valid id WHEN stmt.Exec returns an error THEN return the same error", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()
		repo := NewMySQLCarryRepository(db)

		// conditions
		mock.ExpectPrepare("SELECT id, cid, company_name, address, telephone, locality_id FROM carriers WHERE id=?").
			ExpectQuery().
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "cid", "company_name", "address", "telephone", "locality_id"}).
				AddRow(1, 123, "Test Company", "123 Test St", "1234567890", 1))

		mock.ExpectPrepare("DELETE FROM carriers WHERE id=?").
			ExpectExec().
			WithArgs(1).
			WillReturnError(errors.New("some internal db error"))

		err := repo.Delete(1)
		require.Error(t, err)
	})

	t.Run("GIVEN a valid id WHEN everything is OK THEN delete the carry", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()
		repo := NewMySQLCarryRepository(db)

		// conditions
		mock.ExpectPrepare("SELECT id, cid, company_name, address, telephone, locality_id FROM carriers WHERE id=?").
			ExpectQuery().
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "cid", "company_name", "address", "telephone", "locality_id"}).
				AddRow(1, 123, "Test Company", "123 Test St", "1234567890", 1))

		mock.ExpectPrepare("DELETE FROM carriers WHERE id=?").
			ExpectExec().
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Delete(1)
		require.NoError(t, err)
	})
}
