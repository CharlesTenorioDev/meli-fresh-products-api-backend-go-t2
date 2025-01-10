package repository

import (
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type MySQLCarryRepository struct {
	db *sql.DB
}

func NewMySQLCarryRepository(db *sql.DB) *MySQLCarryRepository {
	return &MySQLCarryRepository{db: db}
}

func (r *MySQLCarryRepository) Save(carry *internal.Carry) error {
	// Prepare statement for inserting data
	stmt, err := r.db.Prepare("INSERT INTO carriers(cid, company_name, address, telephone, locality_id) VALUES(?, ?, ?, ?, ?);")
	if err != nil {
		return err
	}
	// execute the statement
	result, err := stmt.Exec(carry.CID, carry.CompanyName, carry.Address, carry.Telephone, carry.LocalityID)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			if mysqlErr.Number == 1062 {
				return utils.ErrConflict
			}
		}
		return err
	}
	// get the last inserted id
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	carry.ID = int(id)
	return nil
}

func (r *MySQLCarryRepository) GetAll() ([]internal.Carry, error) {
	rows, err := r.db.Query("SELECT id, cid, company_name, address, telephone, locality_id FROM carriers;")
	if err != nil {
		return []internal.Carry{}, err
	}
	defer rows.Close()

	carries := []internal.Carry{}
	for rows.Next() {
		var carry internal.Carry
		err := rows.Scan(&carry.ID, &carry.CID, &carry.CompanyName, &carry.Address, &carry.Telephone, &carry.LocalityID)
		if err != nil {
			return []internal.Carry{}, err
		}
		carries = append(carries, carry)
	}
	return carries, nil
}

func (r *MySQLCarryRepository) GetById(id int) (internal.Carry, error) {
	stmt, err := r.db.Prepare("SELECT id, cid, company_name, address, telephone, locality_id FROM carriers WHERE id=?;")
	if err != nil {
		return internal.Carry{}, err
	}
	row := stmt.QueryRow(id)
	var carry internal.Carry
	err = row.Scan(&carry.ID, &carry.CID, &carry.CompanyName, &carry.Address, &carry.Telephone, &carry.LocalityID)
	if err !=
		nil {
		if errors.Is(err, sql.ErrNoRows) {
			return internal.Carry{}, utils.ErrNotFound
		}
		return internal.Carry{}, err
	}
	return carry, nil
}

func (r *MySQLCarryRepository) Update(carry *internal.Carry) error {
	stmt, err := r.db.Prepare("UPDATE carriers SET cid=?, company_name=?, address=?, telephone=?, locality_id=? WHERE id=?;")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(carry.CID, carry.CompanyName, carry.Address, carry.Telephone, carry.LocalityID, carry.ID)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			if mysqlErr.Number == 1062 {
				return utils.ErrConflict
			}
		}
		return err
	}
	return nil
}

func (r *MySQLCarryRepository) Delete(id int) error {
	_, err := r.GetById(id)
	if err != nil {
		return err
	}
	stmt, err := r.db.Prepare("DELETE FROM carriers WHERE id=?;")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}
