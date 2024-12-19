package repository

import (
	"github.com/meli-fresh-products-api-backend-go-t2/internal/pkg"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type WarehouseDB struct {
	data   map[int]pkg.Warehouse
	nextID int
}

func NewWarehouseDB(load map[int]pkg.Warehouse) *WarehouseDB {
	db := make(map[int]pkg.Warehouse)
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

func (r *WarehouseDB) GetAll() ([]pkg.Warehouse, error) {
	if len(r.data) == 0 {
		return nil, nil
	}

	warehouses := make([]pkg.Warehouse, 0, len(r.data))
	for _, warehouse := range r.data {
		warehouses = append(warehouses, warehouse)
	}

	return warehouses, nil
}

func (r *WarehouseDB) GetById(id int) (pkg.Warehouse, error) {
	warehouse, exists := r.data[id]
	if !exists {
		return pkg.Warehouse{}, utils.ErrNotFound
	}
	return warehouse, nil
}

func (r *WarehouseDB) Save(newWarehouse pkg.Warehouse) (pkg.Warehouse, error) {
	newWarehouse.ID = r.nextID
	r.nextID++

	r.data[newWarehouse.ID] = newWarehouse
	return newWarehouse, nil
}

func (r *WarehouseDB) Update(updatedWarehouse pkg.Warehouse) (pkg.Warehouse, error) {
	if _, exists := r.data[updatedWarehouse.ID]; !exists {
		return pkg.Warehouse{}, utils.ErrNotFound
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
