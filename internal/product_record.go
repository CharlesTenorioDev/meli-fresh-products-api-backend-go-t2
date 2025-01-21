package internal

// ProductRecords represents a single product record with its metadata.
type ProductRecords struct {
	ID             int     `json:"id"`
	LastUpdateDate string  `json:"last_update_date"`
	PurchasePrice  float64 `json:"purchase_price"`
	SalePrice      float64 `json:"sale_price"`
	ProductID      int     `json:"product_id"`
}

type ProductReport struct {
	ProductID    int    `json:"product_id"`
	Description  string `json:"description"`
	RecordsCount int    `json:"records_count"`
}

type ProductRecordsRepository interface {
	Read(productID int) ([]ProductReport, error)
	Create(newProductRecord ProductRecords) (ProductRecords, error)
}

type ProductRecordsService interface {
	GetProductRecords(productID int) ([]ProductReport, error)
	CreateProductRecord(newProductRecord ProductRecords) (ProductRecords, error)
}

type ProductValidation interface {
	GetProductByID(id int) (product Product, err error)
}
