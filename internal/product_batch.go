package internal

type ProductBatch struct {
	ID int `json:"id"`
	ProductBatchRequest
}

type ProductBatchRequest struct {
	BatchNumber        int     `json:"batch_number"`
	CurrentQuantity    int     `json:"current_quantity"`
	CurrentTemperature float64 `json:"current_temperature"`
	DueDate            string  `json:"due_date"`
	InitialQuantity    int     `json:"initial_quantity"`
	ManufacturingDate  string  `json:"manufacturing_date"`
	ManufacturingHour  int     `json:"manufacturing_hour"`
	MinimumTemperature float64 `json:"minimum_temperature"`
	ProductId          int     `json:"product_id"`
	SectionId          int     `json:"section_id"`
}
type (
	ProductBatchRepository interface {
		Save(*ProductBatchRequest) (ProductBatch, error)
		GetBatchNumber(int) (int, error)
	}
	ProductBatchService interface {
		Save(*ProductBatchRequest) (ProductBatch, error)
	}
)
