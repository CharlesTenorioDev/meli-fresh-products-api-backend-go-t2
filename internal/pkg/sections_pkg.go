package pkg

type Section struct {
	ID                 int
	SectionNumber      int
	CurrentCapacity    int
	MaximumCapacity    int
	MinimumCapacity    int
	CurrentTemperature float64
	MinimumTemperature float64
	ProductTypeID      int
	WarehouseID        int
}

type SectionPointers struct {
	ID                 *int
	SectionNumber      *int
	CurrentCapacity    *int
	MaximumCapacity    *int
	MinimumCapacity    *int
	CurrentTemperature *float64
	MinimumTemperature *float64
	ProductTypeID      *int
	WarehouseID        *int
}

type (
	SectionRepository interface {
		GetAll() ([]Section, error)
		// Save(Section) (Section, error)
		// Update(Section) (Section, error)
		// GetById(int) (Section, error)
		// GetBySectionNumber(int) (Section, error)
		// Delete(int) error
	}
	SectionService interface {
		GetAll() ([]Section, error)
		// Save(Section) (Section, error)
		// Update(int, Section) (Section, error)
		// UpdateByFields(int, SectionPointers) (Section, error)
		// GetById(int) (Section, error)
		// Delete(int) error
	}
)
