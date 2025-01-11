package internal

// ProductRecord represents a single product record with its metadata.
type ProductRecord struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}

type ProductRecordsRepository interface {
	Read(productID int) ([]ProductRecord, error)
	Create(newProductRecord ProductRecord) (ProductRecord, error)
}

type ProductRecordsService interface {
	GetProductRecords(productID int) ([]ProductRecord, error)
	CreateProductRecord(newProductRecord ProductRecord) (ProductRecord, error)
}
