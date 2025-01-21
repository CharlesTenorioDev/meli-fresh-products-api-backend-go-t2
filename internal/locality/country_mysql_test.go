package locality_test

import (
	"database/sql"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/locality"
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

func TestIntegrationCountry_GetByName(t *testing.T) {
	db, err := sql.Open("txdb", "fantasy_products")
	require.NoError(t, err)
	defer db.Close()

	repo := locality.NewMysqlCountryRepository(db)

	t.Run("Given an existing name, return the country", func(t *testing.T) {
		country, err := repo.GetByName("USA")
		require.NoError(t, err)
		require.Equal(t, "USA", country.CountryName)
		require.NotZero(t, country.ID)
	})

	t.Run("Given a not existing name, return empty country and utils.ErrNotFound", func(t *testing.T) {
		country, err := repo.GetByName("Ostania")
		require.ErrorIs(t, err, utils.ErrNotFound)
		require.Empty(t, country)
	})
}

func TestIntegrationCountry_Save(t *testing.T) {
	db, err := sql.Open("txdb", "fantasy_products")
	require.NoError(t, err)
	defer db.Close()

	repo := locality.NewMysqlCountryRepository(db)

	t.Run("Given a country with valid fields, save the country", func(t *testing.T) {
		newCountry := internal.Country{
			CountryName: "Westails",
		}
		err := repo.Save(&newCountry)
		require.NoError(t, err)
		require.Equal(t, "Westails", newCountry.CountryName)
		require.NotZero(t, newCountry.ID)
	})
}
