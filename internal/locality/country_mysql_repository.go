package locality

import (
	"database/sql"
	"errors"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type MysqlContryRepository struct {
	db *sql.DB
}

func NewMysqlCountryRepository(db *sql.DB) internal.CountryRepository {
	return &MysqlContryRepository{db: db}
}

func (r *MysqlContryRepository) Save(country *internal.Country) error {
	stmt, err := r.db.Prepare("INSERT INTO countries(country_name) VALUES(?);")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(country.CountryName)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	(*country).ID = int(id)

	return nil
}
func (r *MysqlContryRepository) GetByName(name string) (internal.Country, error) {
	stmt, err := r.db.Prepare("SELECT id, country_name FROM countries WHERE country_name=?;")
	if err != nil {
		return internal.Country{}, err
	}

	row := stmt.QueryRow(name)

	var country internal.Country

	err = row.Scan(&country.ID, &country.CountryName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return internal.Country{}, utils.ErrNotFound
		}

		return internal.Country{}, err
	}

	return country, nil
}
