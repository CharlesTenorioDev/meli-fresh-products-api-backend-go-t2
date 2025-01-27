package internal

type Product struct {
	ID int `json:"id"`
	ProductAttributes
}

type ProductAttributes struct {
	ProductCode                    string  `json:"product_code"`
	Description                    string  `json:"description"`
	Width                          float64 `json:"width"`
	Height                         float64 `json:"height"`
	Length                         float64 `json:"length"`
	NetWeight                      float64 `json:"net_weight"`
	ExpirationRate                 float64 `json:"expiration_rate"`
	RecommendedFreezingTemperature float64 `json:"recommended_freezing_temperature"`
	FreezingRate                   float64 `json:"freezing_rate"`
	ProductType                    int     `json:"product_type"`
	SellerID                       int     `json:"seller_id"`
}

type ProductService interface {
	GetProducts() (listProducts []Product, err error)
	GetProductByID(id int) (product Product, err error)
	CreateProduct(newProduct ProductAttributes) (product Product, err error)
	UpdateProduct(inputProduct Product) (product Product, err error)
	DeleteProduct(id int) (err error)
}

type ProductRepository interface {
	GetAll() (listProducts []Product, err error)
	GetByID(id int) (product Product, err error)
	Create(newproduct ProductAttributes) (product Product, err error)
	Update(inputProduct Product) (product Product, err error)
	Delete(id int) (err error)
}

type ProductTypeValidation interface {
	GetProductTypeByID(id int) (productType ProductType, err error)
}
type SellerValidation interface {
	GetByID(id int) (Seller, error)
}
