package locality_test

import (
	"database/sql"
	"errors"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/testutils"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/locality"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
	"github.com/stretchr/testify/require"
)

func TestUnitLocalityRepository(t *testing.T) {
	t.Run("GetByID", func(t *testing.T) {
		t.Run("GIVEN an OK query WHEN r.db.Prepare THEN return an error", func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()
			repo := locality.NewMysqlLocalityRepository(db)

			// conditions
			wErr := errors.New("some internal db error")
			stmt := mock.ExpectPrepare("SELECT id, locality_name, province_id")
			stmt.WillReturnError(wErr)

			res, err := repo.GetByID(1)

			require.ErrorIs(t, err, wErr)
			require.Zero(t, res)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})

		t.Run("GIVEN an OK query WHEN row.Scan returns sql.ErrNoRows THEN return utils.ErrNotFound", func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()
			repo := locality.NewMysqlLocalityRepository(db)

			// conditions
			stmt := mock.ExpectPrepare("SELECT id, locality_name, province_id")
			stmt.WillBeClosed()
			qr := stmt.ExpectQuery()
			qr.WillReturnError(sql.ErrNoRows)

			res, err := repo.GetByID(1)

			require.ErrorIs(t, err, utils.ErrNotFound)
			require.Zero(t, res)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})

		t.Run("GIVEN an OK query WHEN row.Scan returns an unexpected error THEN return the same error", func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()
			repo := locality.NewMysqlLocalityRepository(db)

			// conditions
			wErr := errors.New("some internal db error")
			stmt := mock.ExpectPrepare("SELECT id, locality_name, province_id")
			stmt.WillBeClosed()
			qr := stmt.ExpectQuery()
			qr.WillReturnError(wErr)

			res, err := repo.GetByID(1)

			require.ErrorIs(t, err, wErr)
			require.Zero(t, res)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	})

	t.Run("GetCarriesByLocalityID", func(t *testing.T) {
		t.Run("GIVEN an 0 locality ID WHEN calling r.db.Query RETURN an error", func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()
			repo := locality.NewMysqlLocalityRepository(db)

			// conditions
			wErr := errors.New("some internal db error")
			stmt := mock.ExpectQuery("SELECT l.id, l.locality_name, COUNT\\(c.id\\) AS 'carries_count'")
			stmt.WillReturnError(wErr)

			res, err := repo.GetCarriesByLocalityID(0)

			require.ErrorIs(t, err, wErr)
			require.Empty(t, res)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})

		t.Run("GIVEN an 0 locality ID WHEN calling rows.Scan RETURN an error", func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()
			repo := locality.NewMysqlLocalityRepository(db)

			// conditions
			rows := mock.NewRows([]string{"id", "locality_name", "carries_count"}).AddRow(nil, "Westails", "1")
			stmt := mock.ExpectPrepare("SELECT l.id, l.locality_name, COUNT\\(c.id\\) AS 'carries_count'")
			qr := stmt.ExpectQuery()
			qr.WillReturnRows(rows)

			res, err := repo.GetCarriesByLocalityID(1)

			require.Contains(t, err.Error(), "converting NULL to int is unsupported")
			require.Empty(t, res)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})

		t.Run("GIVEN an 1 locality ID WHEN r.db.Prepare RETURN an error", func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()
			repo := locality.NewMysqlLocalityRepository(db)

			// conditions
			wErr := errors.New("some internal db error")
			stmt := mock.ExpectPrepare("SELECT l.id, l.locality_name, COUNT\\(c.id\\) AS 'carries_count'")
			stmt.WillReturnError(wErr)

			res, err := repo.GetCarriesByLocalityID(1)

			require.ErrorIs(t, err, wErr)
			require.Empty(t, res)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})

		t.Run("GIVEN an 1 locality ID WHEN stmt.Query RETURN an error", func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()
			repo := locality.NewMysqlLocalityRepository(db)

			// conditions
			wErr := errors.New("some internal db error")
			stmt := mock.ExpectPrepare("SELECT l.id, l.locality_name, COUNT\\(c.id\\) AS 'carries_count'")
			stmt.WillBeClosed()
			qr := stmt.ExpectQuery()
			qr.WillReturnError(wErr)

			res, err := repo.GetCarriesByLocalityID(1)

			require.ErrorIs(t, err, wErr)
			require.Empty(t, res)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	})

	t.Run("GetSellersByLocalityID", func(t *testing.T) {
		t.Run("GIVEN an 0 locality ID WHEN calling r.db.Query RETURN an error", func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()
			repo := locality.NewMysqlLocalityRepository(db)

			// conditions
			wErr := errors.New("some internal db error")
			stmt := mock.ExpectQuery("SELECT l.id, l.locality_name, COUNT\\(s.id\\) AS 'sellers_count'")
			stmt.WillReturnError(wErr)

			res, err := repo.GetSellersByLocalityID(0)

			require.ErrorIs(t, err, wErr)
			require.Empty(t, res)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})

		t.Run("GIVEN an 1 locality ID WHEN calling r.db.Prepare RETURN an error", func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()
			repo := locality.NewMysqlLocalityRepository(db)

			// conditions
			wErr := errors.New("some internal db error")
			stmt := mock.ExpectPrepare("SELECT l.id, l.locality_name, COUNT\\(s.id\\) AS 'sellers_count'")
			stmt.WillReturnError(wErr)

			res, err := repo.GetSellersByLocalityID(1)

			require.ErrorIs(t, err, wErr)
			require.Empty(t, res)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})

		t.Run("GIVEN an 1 locality ID WHEN calling r.db.Prepare RETURN an error", func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()
			repo := locality.NewMysqlLocalityRepository(db)

			// conditions
			wErr := errors.New("some internal db error")
			stmt := mock.ExpectPrepare("SELECT l.id, l.locality_name, COUNT\\(s.id\\) AS 'sellers_count'")
			stmt.WillBeClosed()
			stmt.ExpectQuery().WillReturnError(wErr)

			res, err := repo.GetSellersByLocalityID(1)

			require.ErrorIs(t, err, wErr)
			require.Empty(t, res)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})

		t.Run("GIVEN an 0 locality ID WHEN calling rows.Scan RETURN an error", func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()
			repo := locality.NewMysqlLocalityRepository(db)

			// conditions
			rows := mock.NewRows([]string{"id", "locality_name", "carries_count"}).AddRow(nil, "Westails", "1")
			stmt := mock.ExpectPrepare("SELECT l.id, l.locality_name, COUNT\\(s.id\\) AS 'sellers_count'")
			stmt.WillBeClosed()
			stmt.ExpectQuery().WillReturnRows(rows)

			res, err := repo.GetSellersByLocalityID(1)

			require.Contains(t, err.Error(), "converting NULL to int is unsupported")
			require.Empty(t, res)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})

	})

	t.Run("Save", func(t *testing.T) {
		t.Run("GIVEN a valid locality struct WHEN r.db.Prepare THEN return an error", func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()
			repo := locality.NewMysqlLocalityRepository(db)

			// conditions
			wErr := errors.New("some internal db error")
			mock.ExpectPrepare("INSERT INTO").WillReturnError(wErr)

			err := repo.Save(&internal.Locality{
				LocalityName: "Westails",
				ProvinceID:   2,
			})

			require.ErrorIs(t, err, wErr)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})

		t.Run("GIVEN a valid locality struct WHEN stmt.Exec returns an error THEN return the same error", func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()
			repo := locality.NewMysqlLocalityRepository(db)

			// conditions
			wErr := errors.New("some internal db error")
			stmt := mock.ExpectPrepare("INSERT INTO")
			stmt.WillBeClosed()
			exc := stmt.ExpectExec()
			exc.WillReturnError(wErr)

			err := repo.Save(&internal.Locality{
				LocalityName: "Westails",
				ProvinceID:   2,
			})

			require.ErrorIs(t, err, wErr)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})

	})

}

func TestIntegrationLocalityRepository(t *testing.T) {
	ts := testutils.GetTestDBConn()
	defer ts.End()

	repo := locality.NewMysqlLocalityRepository(ts.DB)

	t.Run("GetById", func(t *testing.T) {
		t.Run("GIVEN an existing ID THEN return the ", func(t *testing.T) {
			locality, err := repo.GetByID(1)
			require.NoError(t, err)
			require.Equal(t, "Los Angeles", locality.LocalityName)
			require.NotZero(t, locality.ProvinceID)
			require.NotZero(t, locality.ID)
		})
		t.Run("GIVEN a not existing ID THEN return empty locality and utils.ErrNotFound", func(t *testing.T) {
			locality, err := repo.GetByID(9999)
			require.ErrorIs(t, err, utils.ErrNotFound)
			require.Empty(t, locality)
		})
	})

	t.Run("Save", func(t *testing.T) {
		t.Run("GIVEN a locality with valid fields THEN save the locality", func(t *testing.T) {
			newLocality := internal.Locality{
				ID:           5,
				LocalityName: "Stella Castle",
				ProvinceID:   1,
			}
			err := repo.Save(&newLocality)
			require.NoError(t, err)
			require.Equal(t, "Stella Castle", newLocality.LocalityName)
			require.NotZero(t, newLocality.ProvinceID)
			require.NotZero(t, newLocality.ID)
		})
	})

	t.Run("GetSellersByLocalityID", func(t *testing.T) {
		t.Run("GIVEN an existing locality ID THEN return the sellers count for that ID", func(t *testing.T) {
			report, err := repo.GetSellersByLocalityID(1)
			require.NoError(t, err)
			require.Len(t, report, 1)
		})
		t.Run("GIVEN an 0 as locality ID THEN return the sellers count for all locations", func(t *testing.T) {
			report, err := repo.GetSellersByLocalityID(0)
			require.NoError(t, err)
			require.Len(t, report, 2)
		})
	})

	t.Run("GetCarriesByLocalityID", func(t *testing.T) {
		t.Run("GIVEN an existing locality ID THEN return the carriers count for that ID", func(t *testing.T) {
			report, err := repo.GetCarriesByLocalityID(1)
			require.NoError(t, err)
			require.Len(t, report, 1)
		})
		t.Run("GIVEN an 0 as locality ID THEN return the carriers count for all locations", func(t *testing.T) {
			report, err := repo.GetCarriesByLocalityID(0)
			require.NoError(t, err)
			require.Len(t, report, 2)
		})
	})

}
