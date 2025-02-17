package section_test

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/testutils"
	"github.com/stretchr/testify/require"

	"github.com/meli-fresh-products-api-backend-go-t2/internal/section"
)

func TestUnitSectionRepository(t *testing.T) {
	t.Run("GetAll", func(t *testing.T) {
		t.Run("WHEN calling r.db.Query returns an error RETURN the same error", func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()
			repo := section.NewSectionMysql(db)

			wErr := errors.New("some internal db error")
			mock.ExpectQuery("SELECT s.id, s.section_number, s.current_temperature, s.minimum_temperature, s.current_capacity, s.minimum_capacity, s.maximum_capacity, s.warehouse_id, s.product_type_id").WillReturnError(wErr)

			_, err := repo.GetAll()

			require.ErrorIs(t, wErr, err)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
		t.Run("GIVEN an result with invalid columns WHEN rows.Scan RETURN a conversion type error", func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()
			repo := section.NewSectionMysql(db)

			rows := sqlmock.NewRows([]string{"id", "section_number", "current_temperature", "minimum_temperature", "current_capacity", "minimum_capacity", "maximum_capacity", "warehouse_id", "product_type_id"}).AddRow(
				"ab", 1, 1, 1, 1, 1, 1, 1, 1,
			)
			mock.ExpectQuery("SELECT s.id, s.section_number, s.current_temperature, s.minimum_temperature, s.current_capacity, s.minimum_capacity, s.maximum_capacity, s.warehouse_id, s.product_type_id").WillReturnRows(rows).RowsWillBeClosed()
			_, err := repo.GetAll()

			require.Contains(t, err.Error(), "converting driver.Value type string")

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	})

	t.Run("GetByID", func(t *testing.T) {
		t.Run("WHEN calling r.db.QueryRow returns sql.ErrNoRows RETURN utils.ErrNotFound", func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()
			repo := section.NewSectionMysql(db)

			mock.ExpectQuery("SELECT s.id, s.section_number, s.current_temperature, s.minimum_temperature, s.current_capacity, s.minimum_capacity, s.maximum_capacity, s.warehouse_id, s.product_type_id").WillReturnError(sql.ErrNoRows)

			_, err := repo.GetByID(0)

			require.ErrorIs(t, utils.ErrNotFound, err)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
		t.Run("WHEN calling r.db.QueryRow returns an error RETURN the same error", func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()
			repo := section.NewSectionMysql(db)

			wErr := errors.New("some internal db error")
			mock.ExpectQuery("SELECT s.id, s.section_number, s.current_temperature, s.minimum_temperature, s.current_capacity, s.minimum_capacity, s.maximum_capacity, s.warehouse_id, s.product_type_id").WillReturnError(wErr)

			_, err := repo.GetByID(0)

			require.ErrorIs(t, wErr, err)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	})

	t.Run("GetBySectionNumber", func(t *testing.T) {
		t.Run("WHEN calling r.db.QueryRow returns sql.ErrNoRows RETURN utils.ErrNotFound", func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()
			repo := section.NewSectionMysql(db)

			mock.ExpectQuery("SELECT s.id, s.section_number, s.current_temperature, s.minimum_temperature, s.current_capacity, s.minimum_capacity, s.maximum_capacity, s.warehouse_id, s.product_type_id").WillReturnError(sql.ErrNoRows)

			_, err := repo.GetBySectionNumber(0)

			require.ErrorIs(t, utils.ErrNotFound, err)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
		t.Run("WHEN calling r.db.QueryRow returns an error RETURN the same error", func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()
			repo := section.NewSectionMysql(db)

			wErr := errors.New("some internal db error")
			mock.ExpectQuery("SELECT s.id, s.section_number, s.current_temperature, s.minimum_temperature, s.current_capacity, s.minimum_capacity, s.maximum_capacity, s.warehouse_id, s.product_type_id").WillReturnError(wErr)

			_, err := repo.GetBySectionNumber(0)

			require.ErrorIs(t, wErr, err)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	})

	t.Run("GetSectionProductsReport", func(t *testing.T) {
		t.Run("WHEN calling r.db.Query returns an error THEN return the same error", func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()
			repo := section.NewSectionMysql(db)

			wErr := errors.New("some internal db error")
			mock.ExpectQuery("SELECT s.id, s.section_number, ifnull\\(sum\\(p.current_quantity\\), 0\\) as products_count").WillReturnError(wErr)

			_, err := repo.GetSectionProductsReport()

			require.ErrorIs(t, wErr, err)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})

		t.Run("WHEN calling r.db.Query returns an error THEN return the same error", func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()
			repo := section.NewSectionMysql(db)

			rows := sqlmock.NewRows([]string{"id", "section_number", "products_count"}).AddRow(
				"ab", 1, 2,
			)
			mock.ExpectQuery("SELECT s.id, s.section_number, ifnull\\(sum\\(p.current_quantity\\), 0\\) as products_count").WillReturnRows(rows)

			_, err := repo.GetSectionProductsReport()

			require.Contains(t, err.Error(), "converting driver.Value type string")

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	})

	t.Run("Save", func(t *testing.T) {
		t.Run("WHEN calling r.db.Exec returns an error THEN return the same error", func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()
			repo := section.NewSectionMysql(db)

			wErr := errors.New("some internal db error")
			mock.ExpectExec("INSERT INTO").WillReturnError(wErr)

			err := repo.Save(&internal.Section{
				SectionNumber:      6,
				CurrentCapacity:    10,
				MaximumCapacity:    20,
				MinimumCapacity:    5,
				CurrentTemperature: 10,
				MinimumTemperature: 5,
				ProductTypeID:      1,
				WarehouseID:        1,
			})

			require.ErrorIs(t, wErr, err)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})

		t.Run("WHEN calling result.LastInsertId returns an error THEN return the same error", func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()
			repo := section.NewSectionMysql(db)

			wErr := errors.New("some internal db error")
			mock.ExpectExec("INSERT INTO").WillReturnResult(sqlmock.NewErrorResult(wErr))

			err := repo.Save(&internal.Section{
				SectionNumber:      6,
				CurrentCapacity:    10,
				MaximumCapacity:    20,
				MinimumCapacity:    5,
				CurrentTemperature: 10,
				MinimumTemperature: 5,
				ProductTypeID:      1,
				WarehouseID:        1,
			})

			require.ErrorIs(t, wErr, err)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	})

	t.Run("Update", func(t *testing.T) {
		t.Run("WHEN calling r.db.Query returns an error THEN return the same error", func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()
			repo := section.NewSectionMysql(db)

			wErr := errors.New("some internal db error")
			mock.ExpectExec("UPDATE sections SET").WillReturnError(wErr)

			err := repo.Update(&internal.Section{
				SectionNumber:      6,
				CurrentCapacity:    10,
				MaximumCapacity:    20,
				MinimumCapacity:    5,
				CurrentTemperature: 10,
				MinimumTemperature: 5,
				ProductTypeID:      1,
				WarehouseID:        1,
			})

			require.ErrorIs(t, wErr, err)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	})
}

func TestIntegrationSection(t *testing.T) {
	ts := testutils.GetTestDBConn()
	defer ts.End()

	repo := section.NewSectionMysql(ts.DB)

	t.Run("GetAll", func(t *testing.T) {
		t.Run("GIVEN 4 existing sections WHEN calling repo.GetAll THEN return the 4 rows", func(t *testing.T) {
			res, err := repo.GetAll()
			require.NoError(t, err)
			require.Len(t, res, 4)
		})
	})

	t.Run("GetByID", func(t *testing.T) {
		t.Run("GIVEN an existing section with ID 1 WHEN calling repo.GetByID(1) THEN return the section", func(t *testing.T) {
			res, err := repo.GetByID(1)
			require.NoError(t, err)
			require.Equal(t, res.ID, 1)
		})
		t.Run("GIVEN a non existing section WHEN calling repo.GetByID THEN return utils.ErrNotFound", func(t *testing.T) {
			res, err := repo.GetByID(9999)
			require.ErrorIs(t, utils.ErrNotFound, err)
			require.Empty(t, res)
		})

	})

	t.Run("GetBySectionNumber", func(t *testing.T) {
		t.Run("GIVEN an existing section with SectionNumber 1 WHEN calling repo.GetBySectionNumber(1) THEN return the section", func(t *testing.T) {
			res, err := repo.GetBySectionNumber(1)
			require.NoError(t, err)
			require.Equal(t, res.SectionNumber, 1)
		})
		t.Run("GIVEN a non existing section WHEN calling repo.GetBySectionNumber THEN return utils.ErrNotFound", func(t *testing.T) {
			res, err := repo.GetBySectionNumber(9999)
			require.ErrorIs(t, utils.ErrNotFound, err)
			require.Empty(t, res)
		})
	})

	t.Run("GetSectionProductsReport", func(t *testing.T) {
		t.Run("GIVEN WHEN repo.GetSectionProductsReport THEN return report", func(t *testing.T) {
			res, err := repo.GetSectionProductsReport()
			require.NoError(t, err)
			require.Len(t, res, 4)
		})
	})

	t.Run("GetSectionProductsReportByID", func(t *testing.T) {
		t.Run("GIVEN WHEN repo.GetSectionProductsReportByID THEN return report", func(t *testing.T) {
			res, err := repo.GetSectionProductsReportByID(1)
			require.NoError(t, err)
			require.Len(t, res, 1)
		})
		t.Run("GIVEN a non existing section WHEN calling repo.GetSectionProductsReportByID THEN return utils.ErrNotFound", func(t *testing.T) {
			res, err := repo.GetSectionProductsReportByID(9999)
			require.ErrorIs(t, utils.ErrNotFound, err)
			require.Empty(t, res)
		})
	})

	t.Run("Save", func(t *testing.T) {
		t.Run("GIVEN a non existing section to create WHEN repo.Save RETURN no error", func(t *testing.T) {
			s := &internal.Section{
				SectionNumber:      6,
				CurrentCapacity:    10,
				MaximumCapacity:    20,
				MinimumCapacity:    5,
				CurrentTemperature: 10,
				MinimumTemperature: 5,
				ProductTypeID:      1,
				WarehouseID:        1,
			}
			err := repo.Save(s)
			require.NoError(t, err)
		})
	})

	t.Run("Update", func(t *testing.T) {
		t.Run("GIVEN an existing section to update WHEN repo.Update RETURN no error", func(t *testing.T) {
			s := &internal.Section{
				SectionNumber:      6,
				CurrentCapacity:    10,
				MaximumCapacity:    20,
				MinimumCapacity:    5,
				CurrentTemperature: 10,
				MinimumTemperature: 5,
				ProductTypeID:      1,
				WarehouseID:        1,
			}
			err := repo.Update(s)
			require.NoError(t, err)
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("GIVEN an existing section with ID 5 WHEN repo.Delete(5) RETURN no error", func(t *testing.T) {
			ts.DB.Exec("INSERT INTO sections (section_number, current_capacity, maximum_capacity, minimum_capacity, current_temperature, minimum_temperature, product_type_id, warehouse_id) VALUES (5, 50, 100, 20, 5.0, -2.0, 1, 1)")
			err := repo.Delete(5)
			require.NoError(t, err)
		})
		t.Run("GIVEN an existing section with ID 1 and associated product_batch WHEN repo.Delete(1) RETURN constraint error", func(t *testing.T) {
			err := repo.Delete(1)
			fmt.Println(err.Error())
			require.Contains(t, err.Error(), "foreign key constraint fails")
		})
	})
}
