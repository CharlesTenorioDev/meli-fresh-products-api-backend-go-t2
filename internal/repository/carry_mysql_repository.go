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

// NewMySQLCarryRepository creates a new instance of MySQLCarryRepository with the given database connection.
// It takes a pointer to an sql.DB object as a parameter and returns a pointer to a MySQLCarryRepository.
//
// Parameters:
//   - db: A pointer to an sql.DB object representing the database connection.
//
// Returns:
//   - A pointer to a MySQLCarryRepository.
func NewMySQLCarryRepository(db *sql.DB) *MySQLCarryRepository {
	return &MySQLCarryRepository{db: db}
}

// Save inserts a new carry record into the carriers table in the MySQL database.
// It takes a pointer to a Carry struct as input and returns an error if any occurs during the process.
// If the insertion is successful, it sets the ID of the carry struct to the last inserted ID.
//
// Parameters:
//   - carry: A pointer to a Carry struct containing the data to be inserted.
//
// Returns:
//   - error: An error if any occurs during the insertion process. If a conflict error (duplicate entry) occurs,
//     it returns a predefined conflict error.
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

// GetAll retrieves all carriers from the database.
// It returns a slice of Carry objects and an error if any occurs during the query execution or row scanning.
// If no carriers are found, it returns an empty slice.
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

// GetByID retrieves a carrier record from the database by its ID.
// It returns an internal.Carry object and an error if any occurs during the process.
// If the carrier with the specified ID is not found, it returns a utils.ErrNotFound error.
//
// Parameters:
//   - id: The ID of the carrier to retrieve.
//
// Returns:
//   - internal.Carry: The carrier object retrieved from the database.
//   - error: An error object if any error occurs, including sql.ErrNoRows if the carrier is not found.
func (r *MySQLCarryRepository) GetByID(id int) (internal.Carry, error) {
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

// Update updates an existing carrier record in the database.
// It takes a pointer to a Carry struct as input and returns an error if the update fails.
// If the update fails due to a duplicate entry (MySQL error 1062), it returns a conflict error.
// Parameters:
//   - carry: a pointer to the Carry struct containing the updated carrier information.
//
// Returns:
//   - error: an error if the update fails, otherwise nil.
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

// Delete removes a carrier record from the database by its ID.
// It first checks if the carrier exists by calling GetByID.
// If the carrier does not exist or an error occurs during the check, it returns an error.
// If the carrier exists, it prepares a DELETE SQL statement and executes it.
// If any error occurs during the preparation or execution of the statement, it returns an error.
// Otherwise, it returns nil indicating the deletion was successful.
//
// Parameters:
//   - id: The ID of the carrier to be deleted.
//
// Returns:
//   - error: An error if the carrier does not exist or if there is an issue with the database operation, otherwise nil.
func (r *MySQLCarryRepository) Delete(id int) error {
	_, err := r.GetByID(id)
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
