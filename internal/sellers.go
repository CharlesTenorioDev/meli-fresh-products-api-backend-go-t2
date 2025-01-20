package internal

type SellerService interface {
	GetAll() ([]Seller, error)
	GetByID(id int) (Seller, error)
	Create(*Seller) error
	Update(int, SellerRequestPointer) (Seller, error)
	Delete(int) error
}

type SellerRepository interface {
	GetAll() (s []Seller, err error)
	GetById(id int) (Seller, error)
	GetByCid(cid int) (Seller, error)
	Create(*Seller) error
	Update(*Seller) error
	Delete(int) error
}

type SellerLocalityValidation interface {
	GetByID(int) (Locality, error)
}

type Seller struct {
	ID          int    `json:"id"`
	Cid         int    `json:"cid"`
	CompanyName string `json:"company_name"`
	Address     string `json:"address"`
	Telephone   string `json:"telephone"`
	LocalityId  int    `json:"locality_id"`
}

type SellerRequest struct {
	Cid         int    `json:"cid"`
	CompanyName string `json:"company_name"`
	Address     string `json:"address"`
	Telephone   string `json:"telephone"`
	LocalityId  int    `json:"locality_id"`
}

type SellerRequestPointer struct {
	Cid         *int    `json:"cid"`
	CompanyName *string `json:"company_name"`
	Address     *string `json:"address"`
	Telephone   *string `json:"telephone"`
}
