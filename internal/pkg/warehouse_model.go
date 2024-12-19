package pkg

type Warehouse struct {
	ID                 int    `json:"id"`
	Address            string `json:"address"`
	Telephone          string `json:"telephone"`
	WarehouseCode      string `json:"warehouse_code"`
	MinimumCapacity    int    `json:"minimum_capacity"`
	MinimumTemperature int    `json:"minimum_temperature"`
}

type WarehousePointers struct {
	Address            *string `json:"address"`
	Telephone          *string `json:"telephone"`
	WarehouseCode      *string `json:"warehouse_code"`
	MinimumCapacity    *int    `json:"minimum_capacity"`
	MinimumTemperature *int    `json:"minimum_temperature"`
}

type WarehouseService interface {
	GetAll() ([]Warehouse, error)
	Save(Warehouse) (Warehouse, error)
	Update(int, WarehousePointers) (Warehouse, error)
	GetById(int) (Warehouse, error)
	Delete(int) error
}

type WarehouseRepository interface {
	GetAll() ([]Warehouse, error)
	Save(Warehouse) (Warehouse, error)
	Update(Warehouse) (Warehouse, error)
	GetById(int) (Warehouse, error)
	Delete(int) error
}