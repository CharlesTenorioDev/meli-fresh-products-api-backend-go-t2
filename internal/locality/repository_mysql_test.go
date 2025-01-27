package locality_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-txdb"
	"github.com/go-sql-driver/mysql"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/locality"

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

func TestIntegrationLocality_GetById(t *testing.T) {
	db, err := sql.Open("txdb", "fantasy_products")
	require.NoError(t, err)
	defer db.Close()

	repo := locality.NewMysqlLocalityRepository(db)

	t.Run("Given an existing ID, return the ", func(t *testing.T) {
		locality, err := repo.GetByID(1)
		require.NoError(t, err)
		require.Equal(t, "Los Angeles", locality.LocalityName)
		require.NotZero(t, locality.ProvinceID)
		require.NotZero(t, locality.ID)
	})

	t.Run("Given a not existing ID, return empty locality and utils.ErrNotFound", func(t *testing.T) {
		locality, err := repo.GetByID(9999)
		require.ErrorIs(t, err, utils.ErrNotFound)
		require.Empty(t, locality)
	})
}

func TestIntegrationLocality_Save(t *testing.T) {
	db, err := sql.Open("txdb", "fantasy_products")
	require.NoError(t, err)
	defer db.Close()

	repo := locality.NewMysqlLocalityRepository(db)

	t.Run("Given a locality with valid fields, save the locality", func(t *testing.T) {
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
}

func TestIntegrationLocality_GetSellersByLocalityId(t *testing.T) {
	db, err := sql.Open("txdb", "fantasy_products")
	require.NoError(t, err)
	defer db.Close()

	repo := locality.NewMysqlLocalityRepository(db)
	t.Run("Given an existing locality ID, return the sellers count for that ID", func(t *testing.T) {
		report, err := repo.GetSellersByLocalityID(1)
		require.NoError(t, err)
		require.Len(t, report, 1)
	})

	t.Run("Given an 0 as locality ID, return the sellers count for all locations", func(t *testing.T) {
		report, err := repo.GetSellersByLocalityID(0)
		require.NoError(t, err)
		require.Len(t, report, 2)
	})
}
