package repository

import (
	"database/sql"
	"errors"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type MysqlLocalityRepository struct {
	db *sql.DB
}

// NewMysqlLocalityRepository creates a new instance of MysqlLocalityRepository with the given database connection.
// It returns an implementation of the LocalityRepository interface.
//
// Parameters:
//   - db: A pointer to an sql.DB instance representing the database connection.
//
// Returns:
//   - An implementation of the LocalityRepository interface.
func NewMysqlLocalityRepository(db *sql.DB) internal.LocalityRepository {
	return &MysqlLocalityRepository{db: db}
}

// Save inserts a new locality record into the localities table in the MySQL database.
// It prepares an SQL statement for inserting the locality's ID, name, and province ID.
// If the preparation or execution of the statement fails, it returns an error.
//
// Parameters:
//
//	locality (*internal.Locality): A pointer to the Locality struct containing the data to be saved.
//
// Returns:
//
//	error: An error object if there is an issue with preparing or executing the SQL statement, otherwise nil.
func (r *MysqlLocalityRepository) Save(locality *internal.Locality) error {
	stmt, err := r.db.Prepare("INSERT INTO localities(id, locality_name, province_id) VALUES(?, ?, ?);")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(locality.ID, locality.LocalityName, locality.ProvinceID)
	if err != nil {
		return err
	}

	return nil
}

// GetByID retrieves a locality by its ID from the MySQL database.
// It prepares a SQL statement to select the locality with the given ID,
// executes the query, and scans the result into an internal.Locality struct.
// If the locality is not found, it returns an ErrNotFound error.
// If any other error occurs during the process, it returns that error.
//
// Parameters:
//   - id: The ID of the locality to retrieve.
//
// Returns:
//   - internal.Locality: The locality with the specified ID.
//   - error: An error if the locality is not found or if any other error occurs.
func (r *MysqlLocalityRepository) GetByID(id int) (internal.Locality, error) {
	stmt, err := r.db.Prepare("SELECT id, locality_name, province_id FROM localities WHERE id=?;")
	if err != nil {
		return internal.Locality{}, err
	}

	defer stmt.Close()
	row := stmt.QueryRow(id)

	var locality internal.Locality

	err = row.Scan(&locality.ID, &locality.LocalityName, &locality.ProvinceID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return internal.Locality{}, utils.ErrNotFound
		}

		return internal.Locality{}, err
	}

	return locality, nil
}

// GetSellersByLocalityID retrieves a list of sellers by locality ID.
// If the localityId is 0, it retrieves the count of sellers for all localities.
// Otherwise, it retrieves the count of sellers for the specified locality ID.
//
// Parameters:
//   - localityId: the ID of the locality to filter sellers by.
//
// Returns:
//   - []internal.SellersByLocality: a slice of SellersByLocality containing the locality ID, locality name, and sellers count.
//   - error: an error if the query fails or if there is an issue scanning the rows.
func (r *MysqlLocalityRepository) GetSellersByLocalityID(localityID int) ([]internal.SellersByLocality, error) {
	report := []internal.SellersByLocality{}

	var rows *sql.Rows

	if localityID == 0 {
		var err error

		rows, err = r.db.Query(`SELECT l.id, l.locality_name, COUNT(s.id) AS 'sellers_count' 
			FROM localities l 
			INNER JOIN sellers s ON s.locality_id=l.id 
			GROUP BY l.id;`)
		if err != nil {
			return []internal.SellersByLocality{}, err
		}

		defer rows.Close()
	} else {
		stmt, err := r.db.Prepare(`SELECT l.id, l.locality_name, COUNT(s.id) AS 'sellers_count' 
			FROM localities l 
			INNER JOIN sellers s ON s.locality_id=l.id 
			WHERE l.id = ?
			GROUP BY l.id;`)
		if err != nil {
			return []internal.SellersByLocality{}, err
		}

		defer stmt.Close()

		rows, err = stmt.Query(localityID)
		if err != nil {
			return []internal.SellersByLocality{}, err
		}

		defer rows.Close()
	}

	for rows.Next() {
		var row internal.SellersByLocality

		err := rows.Scan(&row.LocalityID, &row.LocalityName, &row.SellersCount)
		if err != nil {
			return []internal.SellersByLocality{}, err
		}

		report = append(report, row)
	}

	return report, nil
}

// GetCarriesByLocalityID retrieves the number of carriers associated with a given locality ID.
// If the locality ID is 0, it retrieves the carrier count for all localities.
//
// Parameters:
//   - localityId: The ID of the locality to filter by. If 0, retrieves data for all localities.
//
// Returns:
//   - []internal.CarriesByLocality: A slice of CarriesByLocality structs containing the locality ID, locality name, and carrier count.
//   - error: An error object if an error occurred during the query execution.
func (r *MysqlLocalityRepository) GetCarriesByLocalityID(localityID int) ([]internal.CarriesByLocality, error) {
	report := []internal.CarriesByLocality{}

	var rows *sql.Rows

	if localityID == 0 {
		var err error

		rows, err = r.db.Query(`SELECT l.id, l.locality_name, COUNT(c.id) AS 'carries_count' 
			FROM localities l 
			INNER JOIN carriers c ON c.locality_id=l.id 
			GROUP BY l.id;`)
		if err != nil {
			return []internal.CarriesByLocality{}, err
		}
	} else {
		stmt, err := r.db.Prepare(`SELECT l.id, l.locality_name, COUNT(c.id) AS 'carries_count' 
			FROM localities l 
			INNER JOIN carriers c ON c.locality_id=l.id 
			WHERE l.id = ?
			GROUP BY l.id;`)
		if err != nil {
			return []internal.CarriesByLocality{}, err
		}

		rows, err = stmt.Query(localityID)
		if err != nil {
			return []internal.CarriesByLocality{}, err
		}
	}

	for rows.Next() {
		var row internal.CarriesByLocality

		err := rows.Scan(&row.LocalityID, &row.LocalityName, &row.CarriesCount)
		if err != nil {
			return []internal.CarriesByLocality{}, err
		}

		report = append(report, row)
	}

	return report, nil
}
