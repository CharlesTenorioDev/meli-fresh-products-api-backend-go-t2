package internal

type BuyerAttributes struct {
	CardNumberID string `json:"card_number_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
}

type Buyer struct {
	ID int64 `json:"id"`
	BuyerAttributes
}

type BuyerService interface {
	GetAll() ([]Buyer, error)
	GetOne(id int) (*Buyer, error)
	CreateBuyer(BuyerAttributes) (*Buyer, error)
	UpdateBuyer(*Buyer) (*Buyer, error)
	DeleteBuyer(id int) error
}

type BuyerRepository interface {
	LoadBuyers() (map[int]Buyer, error)
	GetAll() ([]Buyer, error)
	GetOne(id int) (*Buyer, error)
	CreateBuyer(Buyer) (*Buyer, error)
	UpdateBuyer(*Buyer) (*Buyer, error)
	DeleteBuyer(id int) error
}
