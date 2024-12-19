package pkg

// ProductType represents a product type
type ProductType struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}

type ProductTypeRepository interface {
	GetAll() (listProductTypes []ProductType, err error)
	GetByID(id int) (productType ProductType, err error)
	Create(newProductType ProductType) (productType ProductType, err error)
	Update(inputProductType ProductType) (productType ProductType, err error)
	Delete(id int) (err error)
}

type ProductTypeService interface {
	GetProductTypes() (listProductTypes []ProductType, err error)
	GetProductTypeByID(id int) (productType ProductType, err error)
	CreateProductType(newProductType ProductType) (productType ProductType, err error)
	UpdateProductType(inputProductType ProductType) (productType ProductType, err error)
	DeleteProductType(id int) (err error)
}
