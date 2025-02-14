package province_test

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/testutils"
	"testing"

	"github.com/meli-fresh-products-api-backend-go-t2/internal/province"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
	"github.com/stretchr/testify/require"
)

func TestUnitProvinceRepository(t *testing.T) {
	t.Run("GetByName", func(t *testing.T) {
		t.Run("GIVEN an OK query WHEN r.db.Prepare THEN return an error", func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()
			repo := province.NewMysqlProvinceRepository(db)

			// conditions
			wErr := errors.New("some internal db error")
			stmt := mock.ExpectPrepare("SELECT id, province_name, country_id")
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
			repo := province.NewMysqlProvinceRepository(db)

			// conditions
			stmt := mock.ExpectPrepare("SELECT id, province_name, country_id")
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
			repo := province.NewMysqlProvinceRepository(db)

			// conditions
			wErr := errors.New("some internal db error")
			stmt := mock.ExpectPrepare("SELECT id, province_name, country_id")
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
		t.Run("GIVEN a valid province struct WHEN r.db.Prepare returns an error THEN return the same error", func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()
			repo := province.NewMysqlProvinceRepository(db)

			// conditions
			wErr := errors.New("some internal db error")
			stmt := mock.ExpectPrepare("INSERT INTO")
			stmt.WillReturnError(wErr)

			err := repo.Save(&internal.Province{
				ProvinceName: "Westails",
				CountryID:    3,
			})

			require.ErrorIs(t, err, wErr)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})

		t.Run("GIVEN a valid province struct WHEN stmt.Exec returns an error THEN return the same error", func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()
			repo := province.NewMysqlProvinceRepository(db)

			// conditions
			wErr := errors.New("some internal db error")
			stmt := mock.ExpectPrepare("INSERT INTO")
			stmt.WillBeClosed()
			exc := stmt.ExpectExec()
			exc.WillReturnError(wErr)

			err := repo.Save(&internal.Province{
				ProvinceName: "Westails",
				CountryID:    3,
			})

			require.ErrorIs(t, err, wErr)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})

		t.Run("GIVEN a valid province struct WHEN res.LastInsertId returns an error THEN return the same error", func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()
			repo := province.NewMysqlProvinceRepository(db)

			// conditions
			wErr := errors.New("some internal db error")
			stmt := mock.ExpectPrepare("INSERT INTO")
			stmt.WillBeClosed()
			exc := stmt.ExpectExec()
			exc.WillReturnResult(sqlmock.NewErrorResult(wErr))

			err := repo.Save(&internal.Province{
				ProvinceName: "Westails",
				CountryID:    3,
			})

			require.ErrorIs(t, err, wErr)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})

	})
}

func TestIntegrationProvinceProvince(t *testing.T) {
	ts := testutils.GetTestDBConn()
	defer ts.End()

	repo := province.NewMysqlProvinceRepository(ts.DB)

	t.Run("GetByName", func(t *testing.T) {
		t.Run("GIVEN an existing province Ontario WHEN searching by Ontario THEN return the province", func(t *testing.T) {
			p, err := repo.GetByName("Ontario")
			require.NoError(t, err)
			require.Equal(t, "Ontario", p.ProvinceName)
			require.NotZero(t, p.CountryID)
			require.NotZero(t, p.ID)
		})
		t.Run("GIVEN a not existing province Berlint WHEN searching by Berlint THEN return empty province and utils.ErrNotFound", func(t *testing.T) {
			p, err := repo.GetByName("Berlint")
			require.ErrorIs(t, err, utils.ErrNotFound)
			require.Empty(t, p)
		})
	})

	t.Run("Save", func(t *testing.T) {
		t.Run("GIVEN a non existing province with valid fields WHEN saving THEN save the province", func(t *testing.T) {
			newProvince := internal.Province{
				ProvinceName: "Westails",
				CountryID:    2,
			}
			err := repo.Save(&newProvince)
			require.NoError(t, err)
			require.Equal(t, "Westails", newProvince.ProvinceName)
			require.Equal(t, 2, newProvince.CountryID)
			require.NotZero(t, newProvince.ID)
		})
	})
}
