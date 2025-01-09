package internal

type Locality struct {
	ID           int
	LocalityName string
	ProvinceID   int
}

type SellersByLocality struct {
	LocalityId   int
	LocalityName string
	SellersCount int
}

type LocalityRepository interface {
	Save(*Locality) error
	GetById(id int) (Locality, error)
	GetSellersByLocalityId(localityId int) ([]SellersByLocality, error)
}

type LocalityService interface {
	Save(*Locality, *Province, *Country) error
	GetSellersByLocalityId(localityId int) ([]SellersByLocality, error)
}
