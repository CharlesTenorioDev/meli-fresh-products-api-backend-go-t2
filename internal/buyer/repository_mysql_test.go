package buyer_test

import (
	"database/sql"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/buyer"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/testutils"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
	"github.com/stretchr/testify/require"
)

func TestUnitBuyerRepo(t *testing.T) {
	t.Run("GetAll", func(t *testing.T) {
		t.Run("FAIL Get All Buyers", func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			mock.ExpectQuery("").
				WillReturnError(sql.ErrNoRows)

			repo := buyer.NewBuyerDB(db)

			buyers, err := repo.GetAll()
			require.Error(t, err)
			require.Empty(t, buyers)
		})
	})
	t.Run("GetOne", func(t *testing.T) {
		t.Run("FAIL Get One Buyer", func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			mock.ExpectQuery("").
				WillReturnError(sql.ErrNoRows)

			repo := buyer.NewBuyerDB(db)

			buyer, err := repo.GetOne(1)
			require.Error(t, err)
			require.Empty(t, buyer)
		})
	})

	t.Run("CreateBuyer", func(t *testing.T) {
		t.Run("FAIL Create Buyer", func(t *testing.T) {
			var (
				newBuyer = internal.Buyer{}
			)
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			mock.ExpectQuery("").
				WillReturnError(sql.ErrNoRows)

			repo := buyer.NewBuyerDB(db)

			buyer, err := repo.CreateBuyer(newBuyer)
			require.Error(t, err)
			require.Empty(t, buyer)
		})
	})

	t.Run("UpdateBuyer", func(t *testing.T) {
		t.Run("FAIL Update Buyer", func(t *testing.T) {
			var (
				updateBuyer = &internal.Buyer{}
			)
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			mock.ExpectQuery("").
				WillReturnError(sql.ErrNoRows)

			repo := buyer.NewBuyerDB(db)

			buyer, err := repo.UpdateBuyer(updateBuyer)
			require.Error(t, err)
			require.Empty(t, buyer)
		})
	})

	t.Run("DeleteBuyer", func(t *testing.T) {
		t.Run("FAIL Delete Buyer", func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			mock.ExpectQuery("").
				WillReturnError(sql.ErrNoRows)

			repo := buyer.NewBuyerDB(db)

			err = repo.DeleteBuyer(1)
			require.Error(t, err)
		})

		t.Run("FAIL Delete Buyer - No Rows Affected", func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			mock.ExpectExec("DELETE FROM buyers WHERE id = ?").
				WithArgs(1).
				WillReturnResult(sqlmock.NewResult(0, 0))

			repo := buyer.NewBuyerDB(db)

			err = repo.DeleteBuyer(1)
			require.Error(t, err)
			require.Equal(t, utils.ErrNotFound, err)
		})
	})
}

func TestIntegrationBuyerRepo(t *testing.T) {
	ts := testutils.GetTestDBConn()
	defer ts.End()

	repo := buyer.NewBuyerDB(ts.DB)
	t.Run("NewBuyerDB", func(t *testing.T) {
		require.NotNil(t, repo)
	})

	t.Run("GetAll", func(t *testing.T) {
		t.Run("SUCCESS Get All Buyers", func(t *testing.T) {
			buyers, err := repo.GetAll()
			require.NoError(t, err)
			require.NotEmpty(t, buyers)
		})
	})

	t.Run("GetOne", func(t *testing.T) {
		t.Run("SUCCESS Get One Buyer", func(t *testing.T) {
			var (
				wantBuyer = &internal.Buyer{
					ID: 1,
					BuyerAttributes: internal.BuyerAttributes{
						CardNumberID: "B001",
						FirstName:    "Charlie",
						LastName:     "Brown",
					},
				}
			)
			buyer, err := repo.GetOne(1)
			require.NoError(t, err)
			require.NotEmpty(t, buyer)
			require.Equal(t, wantBuyer, buyer)
		})
	})

	t.Run("CreateBuyer", func(t *testing.T) {
		t.Run("SUCCESS Create Buyer", func(t *testing.T) {
			var (
				newBuyer = internal.Buyer{
					BuyerAttributes: internal.BuyerAttributes{
						CardNumberID: "B0012",
						FirstName:    "Charlie",
						LastName:     "Brown",
					},
				}
			)
			buyer, err := repo.CreateBuyer(newBuyer)
			require.NoError(t, err)
			require.NotEmpty(t, buyer)
		})
	})

	t.Run("UpdateBuyer", func(t *testing.T) {
		t.Run("SUCCESS Update Buyer", func(t *testing.T) {
			var (
				updateBuyer = &internal.Buyer{
					BuyerAttributes: internal.BuyerAttributes{
						CardNumberID: "B0012",
						FirstName:    "Charlie",
						LastName:     "Brown",
					},
				}
			)
			buyer, err := repo.UpdateBuyer(updateBuyer)
			require.NoError(t, err)
			require.NotEmpty(t, buyer)
		})
	})

	t.Run("DeleteBuyer", func(t *testing.T) {
		t.Run("SUCCESS Delete Buyer", func(t *testing.T) {
			err := repo.DeleteBuyer(3)
			require.NoError(t, err)
		})
	})

}
