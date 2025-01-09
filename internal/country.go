package internal

type Country struct {
	ID          int
	CountryName string
}

type CountryRepository interface {
	Save(*Country) error
	GetByName(string) (Country, error)
}
