package internal

type Carry struct {
	ID          int    `json:"id"`
	CID         int    `json:"cid"`
	CompanyName string `json:"company_name"`
	Address     string `json:"address"`
	Telephone   string `json:"telephone"`
	LocalityID  int    `json:"locality_id"`
}

type CarryRepository interface {
	Save(*Carry) error
	GetAll() ([]Carry, error)
	GetByID(id int) (Carry, error)
	Update(*Carry) error
	Delete(id int) error
}

type CarryService interface {
	Save(*Carry) error
	GetAll() ([]Carry, error)
	GetByID(id int) (Carry, error)
	Update(*Carry) error
	Delete(id int) error
}

type LocalityValidation interface {
	GetByID(id int) (Locality, error)
}
