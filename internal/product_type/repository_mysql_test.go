package product_type

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/testutils"
	"github.com/stretchr/testify/require"
)

func TestUnitProductTypeDB(t *testing.T) {
	t.Run("GetAll", func(t *testing.T) {
		t.Run("Given an query.error THEN return an error", func(t *testing.T) {
			db, mock, err := sqlmock.New()
			defer db.Close()
			require.NoError(t, err)

			repo := NewProductTypeDB(db)

			mock.ExpectQuery("SELECT id, description FROM product_types").WillReturnError(errors.New("error"))
			listProductTypes, err := repo.GetAll()
			require.Error(t, err)
			require.Zero(t, listProductTypes)
		})
		t.Run("Given an OK query THEN return an error at rows.Scan", func(t *testing.T) {
			db, mock, err := sqlmock.New()
			defer db.Close()
			require.NoError(t, err)

			repo := NewProductTypeDB(db)
			rows := sqlmock.NewRows([]string{"id", "description"}).AddRow(1, "Fruits")

			mock.ExpectQuery("SELECT id, description FROM product_types").WillReturnRows(rows)
			rows.RowError(0, errors.New("error"))
			listProductTypes, err := repo.GetAll()
			require.Error(t, err)
			require.Zero(t, listProductTypes)
		})
	})
	t.Run("GetByID", func(t *testing.T) {
		t.Run("Given an query.error THEN return an error", func(t *testing.T) {
			db, mock, err := sqlmock.New()
			defer db.Close()
			require.NoError(t, err)

			repo := NewProductTypeDB(db)

			mock.ExpectQuery("SELECT id, description FROM product_types WHERE id = ?").WillReturnError(errors.New("error"))
			productType, err := repo.GetByID(1)
			require.Error(t, err)
			require.Zero(t, productType)
		})
		t.Run("Given an OK query THEN return an error at rows.Scan", func(t *testing.T) {
			db, mock, err := sqlmock.New()
			defer db.Close()
			require.NoError(t, err)

			repo := NewProductTypeDB(db)
			rows := sqlmock.NewRows([]string{"id", "description"}).AddRow(1, "Fruits")

			mock.ExpectQuery("SELECT id, description FROM product_types WHERE id = ?").WillReturnRows(rows)
			rows.RowError(0, errors.New("error"))
			productType, err := repo.GetByID(1)
			require.Error(t, err)
			require.Zero(t, productType)
		})
	})
	t.Run("Create", func(t *testing.T) {
		t.Run("Given an query.error THEN return an error", func(t *testing.T) {
			db, mock, err := sqlmock.New()
			defer db.Close()
			require.NoError(t, err)

			repo := NewProductTypeDB(db)

			mock.ExpectExec("INSERT INTO product_types").WillReturnError(errors.New("error"))
			productType, err := repo.Create(internal.ProductType{Description: "Fruits"})
			require.Error(t, err)
			require.Zero(t, productType)
		})
		t.Run("Given an OK query THEN return an error at rows.Scan", func(t *testing.T) {
			db, mock, err := sqlmock.New()
			defer db.Close()
			require.NoError(t, err)

			repo := NewProductTypeDB(db)
			stmt := mock.ExpectPrepare("INSERT INTO product_types").ExpectExec().WithArgs("Fruits").WillReturnResult(sqlmock.NewResult(1, 1))
			stmt.WillReturnError(&mysql.MySQLError{Number: 1062})
			productType, err := repo.Create(internal.ProductType{Description: "Fruits"})
			require.Error(t, err)
			require.Zero(t, productType)
		})
		t.Run("Given an OK query THEN return an error at statement.Exec", func(t *testing.T) {
			db, mock, err := sqlmock.New()
			defer db.Close()
			require.NoError(t, err)

			repo := NewProductTypeDB(db)
			stmt := mock.ExpectPrepare("INSERT INTO product_types").ExpectExec().WithArgs("Fruits").WillReturnResult(sqlmock.NewResult(1, 1))
			stmt.WillReturnError(errors.New("error"))
			productType, err := repo.Create(internal.ProductType{Description: "Fruits"})
			require.Error(t, err)
			require.Zero(t, productType)
		})
	})
	t.Run("Update", func(t *testing.T) {
		t.Run("Given an query.error THEN return an error", func(t *testing.T) {
			db, mock, err := sqlmock.New()
			defer db.Close()
			require.NoError(t, err)

			repo := NewProductTypeDB(db)

			mock.ExpectExec("UPDATE product_types").WillReturnError(errors.New("error"))
			productType, err := repo.Update(internal.ProductType{ID: 1, Description: "Fruits"})
			require.Error(t, err)
			require.Zero(t, productType)
		})
		t.Run("Given an OK query THEN return an error at rows.Scan", func(t *testing.T) {
			db, mock, err := sqlmock.New()
			defer db.Close()
			require.NoError(t, err)

			repo := NewProductTypeDB(db)

			statement := mock.ExpectPrepare("UPDATE product_types SET description=? WHERE id=?")
			result := statement.ExpectExec().WithArgs("Fruits", 1)
			result.WillReturnError(errors.New("error"))
			_, err = repo.Update(internal.ProductType{ID: 1, Description: "Fruits"})
			require.Error(t, err)
		})
	})
	t.Run("Delete", func(t *testing.T) {
		t.Run("Given an query.error THEN return an error", func(t *testing.T) {
			db, mock, err := sqlmock.New()
			defer db.Close()
			require.NoError(t, err)

			repo := NewProductTypeDB(db)

			mock.ExpectExec("DELETE FROM product_types").WillReturnError(errors.New("error"))
			err = repo.Delete(1)
			require.Error(t, err)
		})
		t.Run("Given an OK query THEN return an error at rows.Scan", func(t *testing.T) {
			db, mock, err := sqlmock.New()
			defer db.Close()
			require.NoError(t, err)

			repo := NewProductTypeDB(db)

			result := mock.ExpectPrepare("DELETE FROM product_types WHERE id = ?").ExpectExec().WithArgs(1)
			result.WillReturnError(errors.New("error"))
			err = repo.Delete(1)
			require.Error(t, err)
		})
	})

}

func TestIntegrationProductTypeDB(t *testing.T) {
	ts := testutils.GetTestDBConn()
	defer ts.End()

	repo := NewProductTypeDB(ts.DB)

	t.Run("GetAll", func(t *testing.T) {
		t.Run("Given existing product types, return all product types", func(t *testing.T) {
			productTypes, err := repo.GetAll()
			require.NoError(t, err)
			require.NotEmpty(t, productTypes)
		})

		t.Run("Given no product types, return empty list", func(t *testing.T) {
			// Assuming there's a way to clear the product types table for this test
			_, err := repo.db.Exec("SET FOREIGN_KEY_CHECKS=0")
			require.NoError(t, err)

			_, err = repo.db.Exec("DELETE FROM fresh_products.product_types")
			require.NoError(t, err)

			_, err = repo.db.Exec("SET FOREIGN_KEY_CHECKS=1")
			require.NoError(t, err)

			products, err := repo.GetAll()
			require.NoError(t, err)
			require.Empty(t, products)
		})
	})

	t.Run("GetByID", func(t *testing.T) {

		t.Run("Given existing product type, return product type", func(t *testing.T) {
			productType, err := repo.GetByID(1)
			require.NoError(t, err)
			require.NotEmpty(t, productType)
		})

		t.Run("Given no product type, return empty product type", func(t *testing.T) {
			productType, err := repo.GetByID(0)
			require.Error(t, err)
			require.Empty(t, productType)
		})
	})

	t.Run("Create", func(t *testing.T) {
		t.Run("Given a new product type, create the product type", func(t *testing.T) {
			productType := internal.ProductType{Description: "Fruits"}
			newProductType, err := repo.Create(productType)
			require.NoError(t, err)
			require.NotEmpty(t, newProductType)
		})
	})

	t.Run("Update", func(t *testing.T) {
		t.Run("Given an existing product type, update the product type", func(t *testing.T) {
			productType := internal.ProductType{ID: 1, Description: "Fruits"}
			newProductType, err := repo.Update(productType)
			require.NoError(t, err)
			require.NotEmpty(t, newProductType)
		})

		t.Run("Given a non-existing product type, return error", func(t *testing.T) {
			productType := internal.ProductType{ID: 0, Description: "Fruits"}
			newProductType, err := repo.Update(productType)
			require.Error(t, err)
			require.Empty(t, newProductType)
		})
	})

	t.Run("Delete", func(t *testing.T) {

		t.Run("Given an existing product type, delete the product type", func(t *testing.T) {

			_, err := repo.db.Exec("SET FOREIGN_KEY_CHECKS=0")
			require.NoError(t, err)

			err = repo.Delete(1)

			_, err = repo.db.Exec("SET FOREIGN_KEY_CHECKS=1")
			require.NoError(t, err)
		})

		t.Run("Given a non-existing product type, return error", func(t *testing.T) {
			err := repo.Delete(0)
			require.Error(t, err)
		})
	})
}
