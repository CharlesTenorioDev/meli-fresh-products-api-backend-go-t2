package pkg


type SellerService interface {
	GetAll() (s map[int]Seller)
}

type SellerRepository interface {
	GetAll() (s map[int]Seller, err error)
}

type Seller struct {
	
	Id int
	Cid int
	CompanyName string
	Adress string
	Telephone string
}
