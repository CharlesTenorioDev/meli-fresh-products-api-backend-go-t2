package warehouse

import (
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type MySQLWarehouseRepository struct {
	db *sql.DB
}

func NewWarehouseDB(db *sql.DB) *MySQLWarehouseRepository {
	return &MySQLWarehouseRepository{db: db}
}

// GetAll retrieves all warehouses from the database.
// It returns a slice of Warehouse structs and an error if any occurs during the query execution or row scanning.
// The function queries the database for the following fields: id, address, telephone, warehouse_code, and locality_id.
func (w *MySQLWarehouseRepository) GetAll() ([]internal.Warehouse, error) {
	var warehouseList []internal.Warehouse
	// query the database
	rows, err := w.db.Query("SELECT `id`, `address`, `telephone`, `warehouse_code`, `locality_id`, `minimum_capacity`, `minimum_temperature` FROM warehouses")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var warehouse internal.Warehouse
		// scan the row into the warehouse struct
		err := rows.Scan(&warehouse.ID, &warehouse.Address, &warehouse.Telephone, &warehouse.WarehouseCode, &warehouse.LocalityID, &warehouse.MinimumCapacity, &warehouse.MinimumTemperature)
		if err != nil {
			return nil, err
		}
		// append the warehouse to the slice
		warehouseList = append(warehouseList, warehouse)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return warehouseList, nil
}

// GetByID retrieves a warehouse by its ID from the database.
// It returns the warehouse details if found, otherwise it returns an error.
// If the warehouse is not found, it returns a utils.ErrNotFound error.
//
// Parameters:
//   - id: the ID of the warehouse to retrieve.
//
// Returns:
//   - internal.Warehouse: the warehouse details.
//   - error: an error if the warehouse is not found or if there is a database issue.
func (w *MySQLWarehouseRepository) GetByID(id int) (internal.Warehouse, error) {
	row := w.db.QueryRow("SELECT `id`, `address`, `telephone`, `warehouse_code`, `locality_id`, `minimum_capacity`, `minimum_temperature` FROM warehouses WHERE id = ?", id)

	if err := row.Err(); err != nil {
		return internal.Warehouse{}, err
	}

	var warehouse internal.Warehouse
	err := row.Scan(&warehouse.ID, &warehouse.Address, &warehouse.Telephone, &warehouse.WarehouseCode, &warehouse.LocalityID, &warehouse.MinimumCapacity, &warehouse.MinimumTemperature)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return internal.Warehouse{}, utils.ErrNotFound
		}

		return internal.Warehouse{}, err
	}

	return warehouse, nil
}

// Save inserts a new warehouse record into the database and returns the saved warehouse with its ID.
// If there is a conflict (e.g., duplicate entry), it returns an error indicating the conflict.
//
// Parameters:
//   - newWarehouse: The warehouse object to be saved.
//
// Returns:
//   - internal.Warehouse: The saved warehouse object with its ID populated.
//   - error: An error object if there was an issue during the save operation.
func (w *MySQLWarehouseRepository) Save(newWarehouse internal.Warehouse) (internal.Warehouse, error) {
	// prepare the query
	statement, err := w.db.Prepare("INSERT INTO warehouses (address, telephone, warehouse_code, locality_id, minimum_capacity, minimum_temperature) VALUES (?, ?, ?, ?,?,?)")

	if err != nil {
		return internal.Warehouse{}, err
	}

	defer statement.Close()

	// execute the query
	result, err := statement.Exec(newWarehouse.Address, newWarehouse.Telephone, newWarehouse.WarehouseCode, newWarehouse.LocalityID, newWarehouse.MinimumCapacity, newWarehouse.MinimumTemperature)

	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			switch mysqlErr.Number {
			case 1062:
				err = utils.ErrConflict
				fallthrough
			default:
				return internal.Warehouse{}, err
			}
		}

		return internal.Warehouse{}, err
	}

	// get the last inserted id
	id, err := result.LastInsertId()

	if err != nil {
		return internal.Warehouse{}, err
	}

	// set the id of the warehouse
	newWarehouse.ID = int(id)

	return newWarehouse, nil
}

// Update updates an existing warehouse in the database with the provided updatedWarehouse data.
// It first checks if the warehouse exists by its ID. If it does not exist, it returns an error.
// If the warehouse exists, it prepares and executes an SQL update statement to update the warehouse details.
// If the update is successful, it returns the updated warehouse data.
// If there is a MySQL error, it checks for specific error codes and returns appropriate errors.
// Parameters:
// - updatedWarehouse: The warehouse data to be updated.
// Returns:
// - The updated warehouse data if the update is successful.
// - An error if the warehouse does not exist or if there is an issue with the update operation.
func (w *MySQLWarehouseRepository) Update(updatedWarehouse internal.Warehouse) (internal.Warehouse, error) {
	_, err := w.GetByID(updatedWarehouse.ID)

	if err != nil {
		return internal.Warehouse{}, err
	}
	// prepare the query
	statement, err := w.db.Prepare(
		"UPDATE `warehouses` AS `w` SET `address` = ?, `telephone` = ?, `warehouse_code` = ?, `locality_id` = ?, `minimum_capacity`= ?, `minimum_temperature`= ? WHERE `id` = ?",
	)

	if err != nil {
		return internal.Warehouse{}, err
	}

	defer statement.Close()

	// execute the query
	_, err = statement.Exec(updatedWarehouse.Address, updatedWarehouse.Telephone, updatedWarehouse.WarehouseCode, updatedWarehouse.LocalityID, updatedWarehouse.MinimumCapacity, updatedWarehouse.MinimumTemperature, updatedWarehouse.ID)

	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			switch mysqlErr.Number {
			case 1062:
				err = utils.ErrConflict
				fallthrough
			default:
				return internal.Warehouse{}, err
			}
		}

		return internal.Warehouse{}, err
	}

	return updatedWarehouse, nil
}

// Delete removes a warehouse record from the database by its ID.
// It first checks if the warehouse exists by calling GetByID.
// If the warehouse exists, it prepares and executes a DELETE SQL statement.
// If any error occurs during these operations, it returns the error.
// Parameters:
//   - id: the ID of the warehouse to be deleted.
//
// Returns:
//   - error: an error object if any error occurs, otherwise nil.
func (w *MySQLWarehouseRepository) Delete(warehouseID int) error {
	_, err := w.GetByID(warehouseID)

	if err != nil {
		return err
	}

	statement, err := w.db.Prepare("DELETE FROM warehouses WHERE id = ?")

	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(warehouseID)

	if err != nil {
		return err
	}

	return nil
}
