package pkg


type SellerService interface {
	GetAll() (map[int]Seller, error)
	GetById(id int) (Seller, error)
	Create(SellerRequest)(Seller, error)

}

type SellerRepository interface {
	GetAll() (s map[int]Seller, err error)
	GetById(id int) (Seller, error)
	GetByCid(cid int) (Seller, error)
	Create(SellerRequest)(Seller, error)

}

type Seller struct {
	
	ID            int     `json:"id"`
	Cid           int     `json:"cid"`
	CompanyName   string  `json:"company_name"`
	Address       string  `json:"adress"`
	Telephone     string  `json:"telephone"`
}

type SellerRequest struct {
	Cid           int     `json:"cid"`
	CompanyName   string  `json:"company_name"`
	Address       string  `json:"adress"`
	Telephone     string  `json:"telephone"`
}

