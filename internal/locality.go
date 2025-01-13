package internal

type Locality struct {
	ID           int    `json:"id"`
	LocalityName string `json:"locality_name"`
	ProvinceID   int    `json:"province_id"`
}

type SellersByLocality struct {
	LocalityId   int    `json:"locality_id"`
	LocalityName string `json:"locality_name"`
	SellersCount int    `json:"seller_count"`
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
