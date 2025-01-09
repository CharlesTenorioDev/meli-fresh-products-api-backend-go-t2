package internal

type Province struct {
	ID           int
	ProvinceName string
	CountryID    int
}

type ProvinceRepository interface {
	Save(*Province) error
	GetByName(string) (Province, error)
}
