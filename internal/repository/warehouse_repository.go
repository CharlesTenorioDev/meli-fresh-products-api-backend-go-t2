package repository

import (
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type WarehouseDB struct {
	db     *sql.DB
	data   map[int]internal.Warehouse
	nextID int
}

func NewWarehouseDB(db *sql.DB) *WarehouseDB {
	return &WarehouseDB{db: db}
}

// GetAll retrieves all warehouses from the database.
// It returns a slice of Warehouse structs and an error if any occurs during the query execution or row scanning.
// The function queries the database for the following fields: id, address, telephone, warehouse_code, and locality_id.
func (w *WarehouseDB) GetAll() ([]internal.Warehouse, error) {
	var warehouseList []internal.Warehouse
	// query the database
	rows, err := w.db.Query("SELECT `id`, `address`, `telephone`, `warehouse_code`, `locality_id` FROM warehouses")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var warehouse internal.Warehouse
		// scan the row into the warehouse struct
		err := rows.Scan(&warehouse.ID, &warehouse.Address, &warehouse.Telephone, &warehouse.WarehouseCode, &warehouse.LocalityID)
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

	return warehouseList, nil
}

func (r *WarehouseDB) GetById(id int) (internal.Warehouse, error) {
	row := r.db.QueryRow("SELECT `id`, `address`, `telephone`, `warehouse_code`, `locality_id` FROM warehouses WHERE id = ?", id)
	if err := row.Err(); err != nil {
		return internal.Warehouse{}, err
	}

	var warehouse internal.Warehouse
	err := row.Scan(&warehouse.ID, &warehouse.Address, &warehouse.Telephone, &warehouse.WarehouseCode, &warehouse.LocalityID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return internal.Warehouse{}, utils.ErrNotFound
		}
		return internal.Warehouse{}, err
	}

	return warehouse, nil

}

func (r *WarehouseDB) Save(newWarehouse internal.Warehouse) (internal.Warehouse, error) {
	// prepare the query
	statement, err := r.db.Prepare("INSERT INTO warehouses (address, telephone, warehouse_code, locality_id) VALUES (?, ?, ?, ?)")
	if err != nil {
		return internal.Warehouse{}, err
	}
	defer statement.Close()

	// execute the query
	result, err := statement.Exec(newWarehouse.Address, newWarehouse.Telephone, newWarehouse.WarehouseCode, newWarehouse.LocalityID)
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

func (r *WarehouseDB) Update(updatedWarehouse internal.Warehouse) (internal.Warehouse, error) {
	_, err := r.GetById(updatedWarehouse.ID)
	if err != nil {
		return internal.Warehouse{}, err
	}
	// prepare the query
	statement, err := r.db.Prepare(
		"UPDATE `warehouses` AS `w` SET `address` = ?, `telephone` = ?, `warehouse_code` = ?, `locality_id` = ? WHERE `id` = ?",
	)
	if err != nil {
		return internal.Warehouse{}, err
	}
	defer statement.Close()

	// execute the query
	_, err = statement.Exec(updatedWarehouse.Address, updatedWarehouse.Telephone, updatedWarehouse.WarehouseCode, updatedWarehouse.LocalityID, updatedWarehouse.ID)
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

func (r *WarehouseDB) Delete(id int) error {
	_, err := r.GetById(id)
	if err != nil {
		return err
	}
	statement, err := r.db.Prepare("DELETE FROM warehouses WHERE id = ?")
	if err != nil {
		return err
	}
	defer statement.Close()
	_, err = statement.Exec(id)
	if err != nil {
		return err
	}
	return nil
}
