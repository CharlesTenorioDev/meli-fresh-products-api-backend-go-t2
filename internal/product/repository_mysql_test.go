package product

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-txdb"
	"github.com/go-sql-driver/mysql"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
	"github.com/stretchr/testify/require"
)

func init() {
	cfg := mysql.Config{
		User:   "root",
		Passwd: "example",
		Net:    "tcp",
		Addr:   "localhost:3307",
		DBName: "fresh_products",
	}
	txdb.Register("txdb", "mysql", cfg.FormatDSN())
}

func TestIntegrationProduct_GetAll(t *testing.T) {
	db, err := sql.Open("txdb", "fantasy_products")
	require.NoError(t, err)
	defer db.Close()

	repo := NewProductDB(db)

	t.Run("Given existing products, return all products", func(t *testing.T) {
		products, err := repo.GetAll()
		require.NoError(t, err)
		require.NotEmpty(t, products)
	})

	t.Run("Given no products, return empty list", func(t *testing.T) {
		// Assuming there's a way to clear the products table for this test
		_, err := db.Exec("SET FOREIGN_KEY_CHECKS=0")
		require.NoError(t, err)

		_, err = db.Exec("DELETE FROM fresh_products.products")
		require.NoError(t, err)

		_, err = db.Exec("SET FOREIGN_KEY_CHECKS=1")
		require.NoError(t, err)

		products, err := repo.GetAll()
		require.NoError(t, err)
		require.Empty(t, products)
	})
}

func TestIntegrationProduct_GetByID(t *testing.T) {
	db, err := sql.Open("txdb", "fantasy_products")
	require.NoError(t, err)
	defer db.Close()

	repo := NewProductDB(db)

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

}

func TestIntegrationProduct_Create(t *testing.T) {
	db, err := sql.Open("txdb", "fantasy_products")
	require.NoError(t, err)
	defer db.Close()

	repo := NewProductDB(db)

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
}

func TestIntegrationProduct_Update(t *testing.T) {
	db, err := sql.Open("txdb", "fantasy_products")
	require.NoError(t, err)
	defer db.Close()

	repo := NewProductDB(db)

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
}

func TestIntegrationProduct_Delete(t *testing.T) {
	db, err := sql.Open("txdb", "fantasy_products")
	require.NoError(t, err)
	defer db.Close()

	repo := NewProductDB(db)

	t.Run("Given an existing ID, delete the product", func(t *testing.T) {
		product, err := repo.GetByID(1)
		require.NoError(t, err)

		_, err = db.Exec("SET FOREIGN_KEY_CHECKS=0")
		require.NoError(t, err)

		err = repo.Delete(product.ID)
		require.NoError(t, err)

		_, err = db.Exec("SET FOREIGN_KEY_CHECKS=1")
		require.NoError(t, err)

		deletedProduct, err := repo.GetByID(1)
		require.ErrorIs(t, err, utils.ErrNotFound)
		require.Empty(t, deletedProduct)
	})

	t.Run("Given a not existing ID, return error", func(t *testing.T) {
		_, err := db.Exec("SET FOREIGN_KEY_CHECKS=0")
		require.NoError(t, err)

		err = repo.Delete(9999)
		require.EqualError(t, err, utils.ErrNotFound.Error())

		_, err = db.Exec("SET FOREIGN_KEY_CHECKS=1")
		require.NoError(t, err)

	})
}
