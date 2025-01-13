package repository

import (
	"database/sql"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type WarehouseRepository struct {
	db *sql.DB
}

// NewWarehouseRepository cria uma nova instância de WarehouseRepository
func NewWarehouseRepository(db *sql.DB) *WarehouseRepository {
	return &WarehouseRepository{db: db}
}

// GetAll return all of warehouses in database
func (r *WarehouseRepository) GetAll() ([]internal.Warehouse, error) {
	rows, err := r.db.Query("SELECT id, warehouse_code, address, telephone, minimum_capacity, minimum_temperature FROM warehouses")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var warehouses []internal.Warehouse
	for rows.Next() {
		var warehouse internal.Warehouse
		err := rows.Scan(
			&warehouse.ID,
			&warehouse.WarehouseCode,
			&warehouse.Address,
			&warehouse.Telephone,
			&warehouse.MinimumCapacity,
			&warehouse.MinimumTemperature,
		)
		if err != nil {
			return nil, err
		}
		warehouses = append(warehouses, warehouse)
	}

	if len(warehouses) == 0 {
		return nil, utils.ErrNotFound
	}

	return warehouses, nil
}

// GetById is a function that return one warehouse by id
func (r *WarehouseRepository) GetById(id int) (internal.Warehouse, error) {
	var warehouse internal.Warehouse
	err := r.db.QueryRow(
		"SELECT id, warehouse_code, address, telephone, minimum_capacity, minimum_temperature FROM warehouses WHERE id = ?",
		id,
	).Scan(
		&warehouse.ID,
		&warehouse.WarehouseCode,
		&warehouse.Address,
		&warehouse.Telephone,
		&warehouse.MinimumCapacity,
		&warehouse.MinimumTemperature,
	)
	if err == sql.ErrNoRows {
		return internal.Warehouse{}, utils.ErrNotFound
	}
	if err != nil {
		return internal.Warehouse{}, err
	}

	return warehouse, nil
}

// Save is a function that save in database one warehouse
func (r *WarehouseRepository) Save(newWarehouse internal.Warehouse) (internal.Warehouse, error) {
	result, err := r.db.Exec(
		"INSERT INTO warehouses (warehouse_code, address, telephone, minimum_capacity, minimum_temperature) VALUES (?, ?, ?, ?, ?)",
		newWarehouse.WarehouseCode,
		newWarehouse.Address,
		newWarehouse.Telephone,
		newWarehouse.MinimumCapacity,
		newWarehouse.MinimumTemperature,
	)
	if err != nil {
		return internal.Warehouse{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return internal.Warehouse{}, err
	}

	newWarehouse.ID = int(id)
	return newWarehouse, nil
}

// Update is a function that update warehouse attributes
func (r *WarehouseRepository) Update(updatedWarehouse internal.Warehouse) (internal.Warehouse, error) {
	_, err := r.db.Exec(
		"UPDATE warehouses SET warehouse_code = ?, address = ?, telephone = ?, minimum_capacity = ?, minimum_temperature = ? WHERE id = ?",
		updatedWarehouse.WarehouseCode,
		updatedWarehouse.Address,
		updatedWarehouse.Telephone,
		updatedWarehouse.MinimumCapacity,
		updatedWarehouse.MinimumTemperature,
		updatedWarehouse.ID,
	)
	if err != nil {
		return internal.Warehouse{}, err
	}

	return r.GetById(updatedWarehouse.ID)
}

// Delete is a function that remove one warehouse
func (r *WarehouseRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM warehouses WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}
