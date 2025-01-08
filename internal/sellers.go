package internal

type SellerService interface {
	GetAll() (map[int]Seller, error)
	GetById(id int) (Seller, error)
	Create(SellerRequest) (Seller, error)
	Update(int, SellerRequestPointer) (Seller, error)
	Delete(int) (bool, error)
}

type SellerRepository interface {
	GetAll() (s map[int]Seller, err error)
	GetById(id int) (Seller, error)
	GetByCid(cid int) (Seller, error)
	Create(SellerRequest) (Seller, error)
	Update(Seller) (Seller, error)
	Delete(int) (bool, error)
}

type Seller struct {
	ID          int    `json:"id"`
	Cid         int    `json:"cid"`
	CompanyName string `json:"company_name"`
	Address     string `json:"address"`
	Telephone   string `json:"telephone"`
}

type SellerRequest struct {
	Cid         int    `json:"cid"`
	CompanyName string `json:"company_name"`
	Address     string `json:"address"`
	Telephone   string `json:"telephone"`
}

type SellerRequestPointer struct {
	Cid         *int    `json:"cid"`
	CompanyName *string `json:"company_name"`
	Address     *string `json:"address"`
	Telephone   *string `json:"telephone"`
}
