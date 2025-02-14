package country_test

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/testutils"
	"testing"

	"github.com/meli-fresh-products-api-backend-go-t2/internal/country"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
	"github.com/stretchr/testify/require"
)

func TestUnitCountryRepository(t *testing.T) {
	t.Run("GetByName", func(t *testing.T) {
		t.Run("GIVEN an OK query WHEN r.db.Prepare THEN return an error", func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()
			repo := country.NewMysqlCountryRepository(db)

			// conditions
			wErr := errors.New("some internal db error")
			stmt := mock.ExpectPrepare("SELECT id, country_name")
			stmt.WillReturnError(wErr)

			res, err := repo.GetByName("Ostania")

			require.ErrorIs(t, err, wErr)
			require.Zero(t, res)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})

		t.Run("GIVEN an OK query WHEN row.Scan returns sql.ErrNoRows THEN return utils.ErrNotFound", func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()
			repo := country.NewMysqlCountryRepository(db)

			// conditions
			stmt := mock.ExpectPrepare("SELECT id, country_name")
			stmt.WillBeClosed()
			qr := stmt.ExpectQuery()
			qr.WillReturnError(sql.ErrNoRows)

			res, err := repo.GetByName("Ostania")

			require.ErrorIs(t, err, utils.ErrNotFound)
			require.Zero(t, res)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})

		t.Run("GIVEN an OK query WHEN row.Scan returns an unexpected error THEN return the same error", func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()
			repo := country.NewMysqlCountryRepository(db)

			// conditions
			wErr := errors.New("some internal db error")
			stmt := mock.ExpectPrepare("SELECT id, country_name")
			stmt.WillBeClosed()
			qr := stmt.ExpectQuery()
			qr.WillReturnError(wErr)

			res, err := repo.GetByName("Ostania")

			require.ErrorIs(t, err, wErr)
			require.Zero(t, res)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	})

	t.Run("Save", func(t *testing.T) {
		t.Run("GIVEN a valid country struct WHEN r.db.Prepare returns an error THEN return the same error", func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()
			repo := country.NewMysqlCountryRepository(db)

			// conditions
			wErr := errors.New("some internal db error")
			stmt := mock.ExpectPrepare("INSERT INTO")
			stmt.WillReturnError(wErr)

			err := repo.Save(&internal.Country{CountryName: "Aincrad"})

			require.ErrorIs(t, err, wErr)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})

		t.Run("GIVEN a valid country struct WHEN stmt.Exec returns an error THEN return the same error", func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()
			repo := country.NewMysqlCountryRepository(db)

			// conditions
			wErr := errors.New("some internal db error")
			stmt := mock.ExpectPrepare("INSERT INTO")
			stmt.WillBeClosed()
			exc := stmt.ExpectExec()
			exc.WillReturnError(wErr)

			err := repo.Save(&internal.Country{CountryName: "Aincrad"})

			require.ErrorIs(t, err, wErr)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})

		t.Run("GIVEN a valid country struct WHEN res.LastInsertId returns an error THEN return the same error", func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()
			repo := country.NewMysqlCountryRepository(db)

			// conditions
			wErr := errors.New("some internal db error")
			stmt := mock.ExpectPrepare("INSERT INTO")
			stmt.WillBeClosed()
			exc := stmt.ExpectExec()
			exc.WillReturnResult(sqlmock.NewErrorResult(wErr))

			err := repo.Save(&internal.Country{CountryName: "Aincrad"})

			require.ErrorIs(t, err, wErr)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})

	})
}

func TestIntegrationCountryRepository(t *testing.T) {
	ts := testutils.GetTestDBConn()
	defer ts.End()

	repo := country.NewMysqlCountryRepository(ts.DB)

	t.Run("GetByName", func(t *testing.T) {
		t.Run("GIVEN an existing country USA WHEN searching by USA THEN return the country", func(t *testing.T) {
			c, err := repo.GetByName("USA")
			require.NoError(t, err)
			require.Equal(t, "USA", c.CountryName)
			require.NotZero(t, c.ID)
		})
		t.Run("GIVEN a not existing country Ostania WHEN searching by Ostania THEN return empty country and utils.ErrNotFound", func(t *testing.T) {
			c, err := repo.GetByName("Ostania")
			require.ErrorIs(t, err, utils.ErrNotFound)
			require.Empty(t, c)
		})
	})

	t.Run("Save", func(t *testing.T) {
		t.Run("GIVEN a non existing country with valid fields WHEN saving THEN save the country", func(t *testing.T) {
			newCountry := internal.Country{
				CountryName: "Westails",
			}
			err := repo.Save(&newCountry)
			require.NoError(t, err)
			require.Equal(t, "Westails", newCountry.CountryName)
			require.NotZero(t, newCountry.ID)
		})
	})
}
