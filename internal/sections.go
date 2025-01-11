package internal

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

type SectionProductsReport struct {
	SectionId     int `json:"section_id"`
	SectionNumber int `json:"section_number"`
	ProductsCount int `json:"products_count"`
}

type (
	SectionRepository interface {
		GetAll() ([]Section, error)
		Save(*Section) (Section, error)
		Update(Section) (Section, error)
		GetById(int) (Section, error)
		GetBySectionNumber(int) (Section, error)
		Delete(int) error
		GetSectionProductsReport() ([]SectionProductsReport, error)
		GetSectionProductsReportById(int) ([]SectionProductsReport, error)
	}
	SectionService interface {
		GetAll() ([]Section, error)
		Save(Section) (Section, error)
		Update(int, SectionPointers) (Section, error)
		GetById(int) (Section, error)
		Delete(int) error
		GetSectionProductsReport(int) ([]SectionProductsReport, error)
	}
	SectionWarehouseValidation interface {
		GetById(int) (Warehouse, error)
	}
	SectionProductTypeValidation interface {
		GetProductTypeByID(int) (ProductType, error)
	}
)
