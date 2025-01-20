package internal

type Locality struct {
	ID           int    `json:"id"`
	LocalityName string `json:"locality_name"`
	ProvinceID   int    `json:"province_id"`
}

type SellersByLocality struct {
	LocalityID   int    `json:"locality_id"`
	LocalityName string `json:"locality_name"`
	SellersCount int    `json:"seller_count"`
}

type CarriesByLocality struct {
	LocalityID   int    `json:"locality_id"`
	LocalityName string `json:"locality_name"`
	CarriesCount int    `json:"carries_count"`
}

type LocalityRepository interface {
	Save(*Locality) error
	GetByID(id int) (Locality, error)
	GetSellersByLocalityID(localityID int) ([]SellersByLocality, error)
	GetCarriesByLocalityID(localityID int) ([]CarriesByLocality, error)
}

type LocalityService interface {
	Save(*Locality, *Province, *Country) error
	GetSellersByLocalityID(localityID int) ([]SellersByLocality, error)
	GetCarriesByLocalityID(localityID int) ([]CarriesByLocality, error)
}
