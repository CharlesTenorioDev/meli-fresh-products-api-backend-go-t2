package repository

import (
	"database/sql"
	"errors"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type BuyerRepo struct {
	db *sql.DB
}

func NewBuyerDb(db *sql.DB) *BuyerRepo {
	return &BuyerRepo{db}
}

func (repo *BuyerRepo) GetAll() ([]internal.Buyer, error) {
	query := "SELECT id, id_card_number, first_name, last_name FROM buyers"
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var buyers []internal.Buyer
	for rows.Next() {
		var b internal.Buyer
		err := rows.Scan(&b.ID, &b.CardNumberID, &b.FirstName, &b.LastName)
		if err != nil {
			return nil, err
		}
		buyers = append(buyers, b)
	}
	return buyers, rows.Err()
}

func (repo *BuyerRepo) GetOne(id int) (*internal.Buyer, error) {
	query := "SELECT id, id_card_number, first_name, last_name FROM buyers WHERE id = ?"
	row := repo.db.QueryRow(query, id)

	var buyer internal.Buyer
	err := row.Scan(&buyer.ID, &buyer.CardNumberID, &buyer.FirstName, &buyer.LastName)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, utils.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &buyer, nil
}

func (repo *BuyerRepo) CreateBuyer(newBuyer internal.Buyer) (*internal.Buyer, error) {
	query := "INSERT INTO buyers (id_card_number, first_name, last_name) VALUES (?, ?, ?)"
	result, err := repo.db.Exec(query, newBuyer.CardNumberID, newBuyer.FirstName, newBuyer.LastName)
	if err != nil {
		return nil, err
	}

	insertedID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	newBuyer.ID = int64(insertedID)
	return &newBuyer, nil
}

func (repo *BuyerRepo) UpdateBuyer(updatedBuyer *internal.Buyer) (*internal.Buyer, error) {
	query := "UPDATE buyers SET id_card_number = ?, first_name = ?, last_name = ? WHERE id = ?"
	_, err := repo.db.Exec(query, updatedBuyer.CardNumberID, updatedBuyer.FirstName, updatedBuyer.LastName, updatedBuyer.ID)
	if err != nil {
		return nil, err
	}
	return updatedBuyer, nil
}

func (repo *BuyerRepo) DeleteBuyer(id int) error {
	query := "DELETE FROM buyers WHERE id = ?"
	result, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return utils.ErrNotFound
	}
	return nil
}
