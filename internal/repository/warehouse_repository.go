package repository

import (
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type WarehouseDB struct {
	data   map[int]internal.Warehouse
	nextID int
}

func NewWarehouseDB(load map[int]internal.Warehouse) *WarehouseDB {
	db := make(map[int]internal.Warehouse)
	nextID := 1

	if len(load) > 0 {
		db = load
		nextID = utils.GetBiggestId(load) + 1
	}

	return &WarehouseDB{
		data:   db,
		nextID: nextID,
	}
}

func (r *WarehouseDB) GetAll() ([]internal.Warehouse, error) {
	if len(r.data) == 0 {
		return nil, nil
	}

	warehouses := make([]internal.Warehouse, 0, len(r.data))
	for _, warehouse := range r.data {
		warehouses = append(warehouses, warehouse)
	}

	return warehouses, nil
}

func (r *WarehouseDB) GetById(id int) (internal.Warehouse, error) {
	warehouse, exists := r.data[id]
	if !exists {
		return internal.Warehouse{}, utils.ErrNotFound
	}
	return warehouse, nil
}

func (r *WarehouseDB) Save(newWarehouse internal.Warehouse) (internal.Warehouse, error) {
	newWarehouse.ID = r.nextID
	r.nextID++

	r.data[newWarehouse.ID] = newWarehouse
	return newWarehouse, nil
}

func (r *WarehouseDB) Update(updatedWarehouse internal.Warehouse) (internal.Warehouse, error) {
	if _, exists := r.data[updatedWarehouse.ID]; !exists {
		return internal.Warehouse{}, utils.ErrNotFound
	}

	r.data[updatedWarehouse.ID] = updatedWarehouse
	return updatedWarehouse, nil
}

func (r *WarehouseDB) Delete(id int) error {
	if _, exists := r.data[id]; !exists {
		return utils.ErrNotFound
	}

	delete(r.data, id)
	return nil
}
