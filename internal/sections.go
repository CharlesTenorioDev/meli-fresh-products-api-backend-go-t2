package internal

import "time"

type Section struct {
	ID                 int     `json:"id"`
	SectionNumber      int     `json:"section_number"`
	CurrentCapacity    int     `json:"current_capacity"`
	MaximumCapacity    int     `json:"maximum_capacity"`
	MinimumCapacity    int     `json:"minimum_capacity"`
	CurrentTemperature float64 `json:"current_temperature"`
	MinimumTemperature float64 `json:"minimum_temperature"`
	ProductTypeID      int     `json:"warehouse_id"`
	WarehouseID        int     `json:"product_type_id"`
}

type SectionPointers struct {
	SectionNumber      *int     `json:"section_number"`
	CurrentCapacity    *int     `json:"current_capacity"`
	MaximumCapacity    *int     `json:"maximum_capacity"`
	MinimumCapacity    *int     `json:"minimum_capacity"`
	CurrentTemperature *float64 `json:"current_temperature"`
	MinimumTemperature *float64 `json:"minimum_temperature"`
	ProductTypeID      *int     `json:"product_type_id"`
	WarehouseID        *int     `json:"warehouse_id"`
}

type ProductBatchRequest struct {
	BatchNumber        int       `json:"batch_number"`
	CurrentQuantity    int       `json:"current_quantity"`
	CurrentTemperature float64   `json:"current_temperature"`
	DueDate            time.Time `json:"due_date"`
	InitialQuantity    int       `json:"initial_quantity"`
	ManufacturingDate  time.Time `json:"manufacturing_date"`
	ManufacturingHour  int       `json:"manufacturing_hour"`
	MininumTemperature float64   `json:"mininum_temperature"`
	ProductId          int       `json:"product_id"`
	SectionId          int       `json:"section_id"`
}

type (
	SectionRepository interface {
		GetAll() ([]Section, error)
		Save(*Section) (Section, error)
		Update(Section) (Section, error)
		GetById(int) (Section, error)
		GetBySectionNumber(int) (Section, error)
		Delete(int) error
	}
	SectionService interface {
		GetAll() ([]Section, error)
		Save(Section) (Section, error)
		Update(int, SectionPointers) (Section, error)
		GetById(int) (Section, error)
		Delete(int) error
	}
	SectionWarehouseValidation interface {
		GetById(int) (Warehouse, error)
	}
	SectionProductTypeValidation interface {
		GetProductTypeByID(int) (ProductType, error)
	}
)
