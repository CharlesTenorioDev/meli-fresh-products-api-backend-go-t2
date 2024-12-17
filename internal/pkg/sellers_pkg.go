package pkg


type SellerService interface {
	GetAll() (map[int]Seller, error)
	GetById(id int) (Seller, error)


}

type SellerRepository interface {
	GetAll() (s map[int]Seller, err error)
	GetById(id int) (Seller, error)

}

type Seller struct {
	
	ID           int
	Cid 		 int
	CompanyName  string
	Address 	 string
	Telephone 	 string
}
