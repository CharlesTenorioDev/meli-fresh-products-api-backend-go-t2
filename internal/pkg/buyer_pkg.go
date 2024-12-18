package pkg

type Buyer struct {
	ID           int64  `json:"id"`
	CardNumberID string `json:"card_number_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
}

type BuyerService interface {
	GetAll() ([]Buyer, error)
}

type BuyerRepository interface {
	LoadBuyers() (map[int]Buyer, error)
	GetAll() ([]Buyer, error)
}
