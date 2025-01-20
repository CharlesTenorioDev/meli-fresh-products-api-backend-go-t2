package repository

import (
	"database/sql"
	"errors"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type MysqlProvinceRepository struct {
	db *sql.DB
}

func NewMysqlProvinceRepository(db *sql.DB) internal.ProvinceRepository {
	return &MysqlProvinceRepository{db: db}
}

func (r *MysqlProvinceRepository) GetByName(name string) (internal.Province, error) {
	stmt, err := r.db.Prepare("SELECT id, province_name, country_id FROM provinces WHERE province_name=?;")
	if err != nil {
		return internal.Province{}, err
	}

	row := stmt.QueryRow(name)

	var province internal.Province

	err = row.Scan(&province.ID, &province.ProvinceName, &province.CountryID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return internal.Province{}, utils.ErrNotFound
		}

		return internal.Province{}, err
	}

	return province, nil
}

func (r *MysqlProvinceRepository) Save(province *internal.Province) error {
	stmt, err := r.db.Prepare("INSERT INTO provinces(province_name, country_id) VALUES(?, ?);")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(province.ProvinceName, province.CountryID)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	(*province).ID = int(id)

	return nil
}
