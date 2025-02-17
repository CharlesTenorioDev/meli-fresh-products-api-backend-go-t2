package product

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/testutils"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
	"github.com/stretchr/testify/require"
)

func TestUnitProductRepository(t *testing.T) {
	t.Run("GetByID", func(t *testing.T) {
		t.Run("Given a r.db.Query error, return nil list of products", func(t *testing.T) {
			db, mock, err := sqlmock.New()
			defer db.Close()
			require.NoError(t, err)

			repo := NewProductDB(db)

			mock.ExpectQuery("SELECT id, description, expiration_rate, freezing_rate, height, length, net_weight, product_code, recommended_freezing_temperature, width, product_type_id, seller_id FROM products WHERE id = ?").WillReturnError(errors.New("error"))

			product, err := repo.GetByID(1)
			require.Error(t, err)
			require.Empty(t, product)
		})
		t.Run("Given a r.db.Query OK, When Row.Scan, return nil list of products", func(t *testing.T) {
			db, mock, err := sqlmock.New()
			defer db.Close()
			require.NoError(t, err)

			repo := NewProductDB(db)

			rows := sqlmock.NewRows([]string{"id", "description", "expiration_rate", "freezing_rate", "height", "length", "net_weight", "product_code", "recommended_freezing_temperature", "width", "product_type_id", "seller_id"}).
				AddRow(1, "description", 0.1, 0.2, 1.2, 1.3, 1.4, "product_code", -2.5, 5, 2, 1)

			mock.ExpectQuery("SELECT id, description, expiration_rate, freezing_rate, height, length, net_weight, product_code, recommended_freezing_temperature, width, product_type_id, seller_id FROM products WHERE id = ?").WillReturnRows(rows)

			rows.RowError(0, errors.New("error"))
			product, err := repo.GetByID(1)
			require.Error(t, err)
			require.Empty(t, product)
		})
	})
	t.Run("GetAll", func(t *testing.T) {
		t.Run("Given a r.db.Query error, return nil list of products", func(t *testing.T) {
			db, mock, err := sqlmock.New()
			defer db.Close()
			require.NoError(t, err)

			repo := NewProductDB(db)

			mock.ExpectQuery("SELECT id, description, expiration_rate, freezing_rate, height, length, net_weight, product_code, recommended_freezing_temperature, width, product_type_id, seller_id FROM fresh_products.products").WillReturnError(errors.New("error"))

			products, err := repo.GetAll()
			require.Error(t, err)
			require.Empty(t, products)
		})
		t.Run("Given a r.db.Query OK, When Row.Scan, return nil list of products", func(t *testing.T) {
			db, mock, err := sqlmock.New()
			defer db.Close()
			require.NoError(t, err)

			repo := NewProductDB(db)

			rows := sqlmock.NewRows([]string{"id", "description", "expiration_rate", "freezing_rate", "height", "length", "net_weight", "product_code", "recommended_freezing_temperature", "width", "product_type_id", "seller_id"}).
				AddRow(1, "description", 0.1, 0.2, 1.2, 1.3, 1.4, "product_code", -2.5, 5, 2, 1)

			mock.ExpectQuery("SELECT id, description, expiration_rate, freezing_rate, height, length, net_weight, product_code, recommended_freezing_temperature, width, product_type_id, seller_id FROM fresh_products.products").WillReturnRows(rows)

			rows.RowError(0, errors.New("error"))
			products, err := repo.GetAll()
			require.Error(t, err)
			require.Empty(t, products)
		})
	})
	t.Run("Create", func(t *testing.T) {
		t.Run("Given a r.db.Query error, return nil list of products", func(t *testing.T) {
			db, mock, err := sqlmock.New()
			defer db.Close()
			require.NoError(t, err)

			repo := NewProductDB(db)

			mock.ExpectExec("INSERT INTO products").WillReturnError(errors.New("error"))

			product, err := repo.Create(internal.ProductAttributes{})
			require.Error(t, err)
			require.Empty(t, product)
		})
		t.Run("Given a r.db.Query OK, When Row.Scan, return nil list of products", func(t *testing.T) {
			db, mock, err := sqlmock.New()
			defer db.Close()
			require.NoError(t, err)

			repo := NewProductDB(db)

			statement := mock.ExpectPrepare("INSERT INTO products description, expiration_rate, freezing_rate, height, `length`, net_weight, product_code, recommended_freezing_temperature, width, product_type_id, seller_id")
			statement.ExpectExec().WillReturnError(&mysql.MySQLError{Number: 1062})

			product, err := repo.Create(internal.ProductAttributes{})
			require.Error(t, err)
			require.Empty(t, product)
		})

	})
	t.Run("Update", func(t *testing.T) {
		t.Run("Given a r.db.Query error, return nil list of products", func(t *testing.T) {
			db, mock, err := sqlmock.New()
			defer db.Close()
			require.NoError(t, err)

			repo := NewProductDB(db)

			mock.ExpectExec("UPDATE products").WillReturnError(errors.New("error"))

			_, err = repo.Update(internal.Product{})
			require.Error(t, err)
		})
		t.Run("Given a r.db.Query OK, When Row.Scan, return nil list of products", func(t *testing.T) {
			db, mock, err := sqlmock.New()
			defer db.Close()
			require.NoError(t, err)

			repo := NewProductDB(db)

			statement := mock.ExpectPrepare("UPDATE products SET description").WillReturnError(errors.New("error"))
			statement.WillReturnError(errors.New("error"))
			_, err = repo.Update(internal.Product{})
			require.Error(t, err)
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("Given a r.db.Query error, return nil list of products", func(t *testing.T) {
			db, mock, err := sqlmock.New()
			defer db.Close()
			require.NoError(t, err)

			repo := NewProductDB(db)

			mock.ExpectExec("DELETE FROM products").WillReturnError(errors.New("error"))

			err = repo.Delete(1)
			require.Error(t, err)
		})
		t.Run("Given a r.db.Query OK, When Row.Scan, return nil list of products", func(t *testing.T) {
			db, mock, err := sqlmock.New()
			defer db.Close()
			require.NoError(t, err)

			repo := NewProductDB(db)

			mock.ExpectExec("DELETE FROM products").WillReturnError(errors.New("error"))

			err = repo.Delete(1)
			require.Error(t, err)
		})
	})

}

func TestIntegrationProduct(t *testing.T) {
	ts := testutils.GetTestDBConn()
	defer ts.End()

	repo := NewProductDB(ts.DB)

	t.Run("GetByID", func(t *testing.T) {
		t.Run("Given an existing ID, return the product", func(t *testing.T) {

			product, err := repo.GetByID(1)
			expectedProduct := internal.Product{
				ID: 1,
				ProductAttributes: internal.ProductAttributes{
					ProductCode:                    "PA001",
					Description:                    "Fresh Apples",
					Width:                          5,
					Height:                         4.5,
					Length:                         7.5,
					NetWeight:                      1.2,
					ExpirationRate:                 0.1,
					RecommendedFreezingTemperature: -2.5,
					FreezingRate:                   0.2,
					ProductType:                    2,
					SellerID:                       1,
				}}

			require.NoError(t, err)
			require.Equal(t, expectedProduct, product)

		})

		t.Run("Given a not existing ID, return empty product and utils.ErrNotFound", func(t *testing.T) {
			product, err := repo.GetByID(9999)
			require.ErrorIs(t, err, utils.ErrNotFound)
			require.Empty(t, product)
		})
	})

	t.Run("Create", func(t *testing.T) {
		t.Run("Given a product with valid fields, save the product", func(t *testing.T) {
			newProduct := internal.ProductAttributes{
				ProductCode:                    "PA005",
				Description:                    "Fresh Bananas",
				Width:                          5,
				Height:                         4.5,
				Length:                         7.5,
				NetWeight:                      1.2,
				ExpirationRate:                 0.1,
				RecommendedFreezingTemperature: -2.5,
				FreezingRate:                   0.2,
				ProductType:                    2,
				SellerID:                       1,
			}
			product, err := repo.Create(newProduct)
			require.NoError(t, err)
			require.Equal(t, newProduct, product.ProductAttributes)
			require.NotZero(t, product.ID)
		})
		t.Run("Given a product with invalid fields, return error", func(t *testing.T) {
			newProduct := internal.ProductAttributes{
				ProductCode:                    "PA005",
				Description:                    "Fresh Bananas",
				Width:                          5,
				Height:                         4.5,
				Length:                         7.5,
				NetWeight:                      1.2,
				ExpirationRate:                 0.1,
				RecommendedFreezingTemperature: -2.5,
				FreezingRate:                   0.2,
				ProductType:                    2,
				SellerID:                       0,
			}
			_, err := repo.Create(newProduct)
			require.Error(t, err)
			require.EqualError(t, err, "Error 1452 (23000): Cannot add or update a child row: a foreign key constraint fails (`fresh_products`.`products`, CONSTRAINT `products_ibfk_2` FOREIGN KEY (`seller_id`) REFERENCES `sellers` (`id`))")
		})
	})

	t.Run("Update", func(t *testing.T) {
		t.Run("Given a product with valid fields, update the product", func(t *testing.T) {
			product, err := repo.GetByID(1)
			require.NoError(t, err)

			product.Description = "Updated Description"
			_, err = repo.Update(product)
			require.NoError(t, err)

			updatedProduct, err := repo.GetByID(1)
			require.NoError(t, err)
			require.Equal(t, product, updatedProduct)
		})

		t.Run("Given a product with invalid fields, return error", func(t *testing.T) {
			product, err := repo.GetByID(1)
			require.NoError(t, err)

			product.SellerID = 0
			_, err = repo.Update(product)
			require.Error(t, err)
			require.EqualError(t, err, "Error 1452 (23000): Cannot add or update a child row: a foreign key constraint fails (`fresh_products`.`products`, CONSTRAINT `products_ibfk_2` FOREIGN KEY (`seller_id`) REFERENCES `sellers` (`id`))")
		})

	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("Given an existing ID, delete the product", func(t *testing.T) {
			product, err := repo.GetByID(1)
			require.NoError(t, err)

			_, err = repo.db.Exec("SET FOREIGN_KEY_CHECKS=0")
			require.NoError(t, err)

			err = repo.Delete(product.ID)
			require.NoError(t, err)

			_, err = repo.db.Exec("SET FOREIGN_KEY_CHECKS=1")
			require.NoError(t, err)

			deletedProduct, err := repo.GetByID(1)
			require.ErrorIs(t, err, utils.ErrNotFound)
			require.Empty(t, deletedProduct)
		})

		t.Run("Given a not existing ID, return error", func(t *testing.T) {
			_, err := repo.db.Exec("SET FOREIGN_KEY_CHECKS=0")
			require.NoError(t, err)

			err = repo.Delete(9999)
			require.EqualError(t, err, utils.ErrNotFound.Error())

			_, err = repo.db.Exec("SET FOREIGN_KEY_CHECKS=1")
			require.NoError(t, err)

		})

	})

	t.Run("GetAll", func(t *testing.T) {

		t.Run("Given existing products, return all products", func(t *testing.T) {
			products, err := repo.GetAll()
			require.NoError(t, err)
			require.NotEmpty(t, products)
		})

		t.Run("Given no products, return empty list", func(t *testing.T) {
			// Assuming there's a way to clear the products table for this test
			_, err := repo.db.Exec("SET FOREIGN_KEY_CHECKS=0")
			require.NoError(t, err)

			_, err = repo.db.Exec("DELETE FROM fresh_products.products")
			require.NoError(t, err)

			_, err = repo.db.Exec("SET FOREIGN_KEY_CHECKS=1")
			require.NoError(t, err)

			products, err := repo.GetAll()
			require.NoError(t, err)
			require.Empty(t, products)
		})

	})
}
