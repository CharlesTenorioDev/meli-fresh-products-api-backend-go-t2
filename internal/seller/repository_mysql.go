package seller

import (
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

// MySQLSellerRepository is the mysql implementation of the seller repository
type MySQLSellerRepository struct {
	// db is the database connection to mysql
	db *sql.DB
}

// NewSellerRepository creates a new instance of the seller repository
func NewSellerRepository(db *sql.DB) internal.SellerRepository {
	return &MySQLSellerRepository{db}
}

// GetAll returns all sellers from the database
func (r *MySQLSellerRepository) GetAll() (sellers []internal.Seller, err error) {
	// execute the query
	rows, err := r.db.Query("SELECT `id`, `cid`, `company_name`, `address`, `telephone`, `locality_id` FROM `sellers`")
	if err != nil {
		return
	}

	// iterate over the rows
	for rows.Next() {
		// create a new seller
		var seller internal.Seller

		err = rows.Scan(&seller.ID, &seller.Cid, &seller.CompanyName, &seller.Address, &seller.Telephone, &seller.LocalityID)
		if err != nil {
			return
		}

		// append the seller to the slice
		sellers = append(sellers, seller)
	}

	// check for errors
	err = rows.Err()
	if err != nil {
		return
	}

	return
}

// GetByID returns a seller from the database by its id
func (r *MySQLSellerRepository) GetByID(id int) (seller internal.Seller, err error) {
	// execute the query
	row := r.db.QueryRow("SELECT `id`, `cid`, `company_name`, `address`, `telephone`, `locality_id` FROM `sellers` WHERE `id` = ?", id)

	// scan the row into the seller
	err = row.Scan(&seller.ID, &seller.Cid, &seller.CompanyName, &seller.Address, &seller.Telephone, &seller.LocalityID)
	if err != nil {
		if err == sql.ErrNoRows {
			err = utils.ErrNotFound
			return
		}

		return
	}

	return
}
func (r *MySQLSellerRepository) GetByCid(cid int) (seller internal.Seller, err error) {
	// execute the query
	row := r.db.QueryRow("SELECT `id`, `cid`, `company_name`, `address`, `telephone`, `locality_id` FROM `sellers` WHERE `cid` = ?", cid)

	// scan the row into the seller
	err = row.Scan(&seller.ID, &seller.Cid, &seller.CompanyName, &seller.Address, &seller.Telephone, &seller.LocalityID)
	if err != nil {
		if err == sql.ErrNoRows {
			err = utils.ErrNotFound
			return
		}

		return
	}

	return
}

// Create Save saves a seller into the database
func (r *MySQLSellerRepository) Create(seller *internal.Seller) (err error) {
	// execute the query
	result, err := r.db.Exec(
		"INSERT INTO `sellers` (`cid`, `company_name`, `address`, `telephone`, `locality_id`) VALUES (?, ?, ?, ?, ?)",
		(*seller).Cid, (*seller).CompanyName, (*seller).Address, (*seller).Telephone, (*seller).LocalityID,
	)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			switch mysqlErr.Number {
			case 1062:
				err = utils.ErrConflict
			default:
			}

			return err
		}

		return err
	}

	// get the last inserted id
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// set the id of the seller
	(*seller).ID = int(id)

	return err
}

// Update updates a seller in the database
func (r *MySQLSellerRepository) Update(seller *internal.Seller) (err error) {
	// execute the query
	_, err = r.db.Exec(
		"UPDATE `sellers` SET `cid` = ?, `company_name` = ?, `address` = ?, `telephone` = ?, `locality_id` = ? WHERE `id` = ?",
		(*seller).Cid, (*seller).CompanyName, (*seller).Address, (*seller).Telephone, (*seller).LocalityID, (*seller).ID,
	)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			switch mysqlErr.Number {
			case 1062:
				err = utils.ErrConflict
			default:
			}

			return
		}

		return
	}

	return
}

// Delete deletes a seller from the database
func (r *MySQLSellerRepository) Delete(id int) error {
	// execute the query
	_, err := r.db.Exec("DELETE FROM `sellers` WHERE `id` = ?", id)
	if err != nil {
		return err
	}

	return err
}
