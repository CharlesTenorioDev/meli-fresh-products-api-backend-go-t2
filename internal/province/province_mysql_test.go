package province_test

import (
	"database/sql"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/province"
	"testing"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
	"github.com/stretchr/testify/require"
)

func TestIntegrationProvince_GetByName(t *testing.T) {
	db, err := sql.Open("txdb", "fantasy_products")
	require.NoError(t, err)
	defer db.Close()

	repo := province.NewMysqlProvinceRepository(db)

	t.Run("Given an existing name, return the province", func(t *testing.T) {
		province, err := repo.GetByName("Ontario")
		require.NoError(t, err)
		require.Equal(t, "Ontario", province.ProvinceName)
		require.NotZero(t, province.CountryID)
		require.NotZero(t, province.ID)
	})

	t.Run("Given a not existing name, return empty province and utils.ErrNotFound", func(t *testing.T) {
		province, err := repo.GetByName("Berlint")
		require.ErrorIs(t, err, utils.ErrNotFound)
		require.Empty(t, province)
	})
}

func TestIntegrationProvince_Save(t *testing.T) {
	db, err := sql.Open("txdb", "fantasy_products")
	require.NoError(t, err)
	defer db.Close()

	repo := province.NewMysqlProvinceRepository(db)

	t.Run("Given a province with valid fields, save the province", func(t *testing.T) {
		newProvince := internal.Province{
			ProvinceName: "Westails",
			CountryID:    2,
		}
		err := repo.Save(&newProvince)
		require.NoError(t, err)
		require.Equal(t, "Westails", newProvince.ProvinceName)
		require.NotZero(t, newProvince.CountryID)
		require.NotZero(t, newProvince.ID)
	})
}
