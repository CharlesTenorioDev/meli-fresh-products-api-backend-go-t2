package repository

import (
	"database/sql"
	"errors"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type MysqlLocalityRepository struct {
	db *sql.DB
}

func NewMysqlLocalityRepository(db *sql.DB) internal.LocalityRepository {
	return &MysqlLocalityRepository{db: db}
}

func (r *MysqlLocalityRepository) Save(locality *internal.Locality) error {
	stmt, err := r.db.Prepare("INSERT INTO localities(id, locality_name, province_id) VALUES(?, ?, ?);")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(locality.ID, locality.LocalityName, locality.ProvinceID)
	if err != nil {
		return err
	}
	return nil
}

func (r *MysqlLocalityRepository) GetById(id int) (internal.Locality, error) {
	stmt, err := r.db.Prepare("SELECT id, locality_name, province_id FROM localities WHERE id=?;")
	if err != nil {
		return internal.Locality{}, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(id)
	var locality internal.Locality
	err = row.Scan(&locality.ID, &locality.LocalityName, &locality.ProvinceID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return internal.Locality{}, utils.ErrNotFound
		}
		return internal.Locality{}, err
	}
	return locality, nil
}

func (r *MysqlLocalityRepository) GetSellersByLocalityId(localityId int) ([]internal.SellersByLocality, error) {
	report := []internal.SellersByLocality{}
	var rows *sql.Rows
	if localityId == 0 {
		var err error
		rows, err = r.db.Query(`SELECT l.id, l.locality_name, COUNT(s.id) AS 'sellers_count' 
			FROM localities l 
			INNER JOIN sellers s ON s.locality_id=l.id 
			GROUP BY l.id;`)
		if err != nil {
			return []internal.SellersByLocality{}, err
		}
		defer rows.Close()
	} else {
		stmt, err := r.db.Prepare(`SELECT l.id, l.locality_name, COUNT(s.id) AS 'sellers_count' 
			FROM localities l 
			INNER JOIN sellers s ON s.locality_id=l.id 
			WHERE l.id = ?
			GROUP BY l.id;`)
		if err != nil {
			return []internal.SellersByLocality{}, err
		}
		defer stmt.Close()
		rows, err = stmt.Query(localityId)
		if err != nil {
			return []internal.SellersByLocality{}, err
		}
		defer rows.Close()
	}
	for rows.Next() {
		var row internal.SellersByLocality
		err := rows.Scan(&row.LocalityId, &row.LocalityName, &row.SellersCount)
		if err != nil {
			return []internal.SellersByLocality{}, err
		}
		report = append(report, row)
	}
	return report, nil
}

func (r *MysqlLocalityRepository) GetCarriesByLocalityId(localityId int) ([]internal.CarriesByLocality, error) {
	report := []internal.CarriesByLocality{}
	var rows *sql.Rows
	if localityId == 0 {
		var err error
		rows, err = r.db.Query(`SELECT l.id, l.locality_name, COUNT(c.id) AS 'carries_count' 
			FROM localities l 
			INNER JOIN carriers c ON c.locality_id=l.id 
			GROUP BY l.id;`)
		if err != nil {
			return []internal.CarriesByLocality{}, err
		}
	} else {
		stmt, err := r.db.Prepare(`SELECT l.id, l.locality_name, COUNT(c.id) AS 'carries_count' 
			FROM localities l 
			INNER JOIN carriers c ON c.locality_id=l.id 
			WHERE l.id = ?
			GROUP BY l.id;`)
		if err != nil {
			return []internal.CarriesByLocality{}, err
		}
		rows, err = stmt.Query(localityId)
		if err != nil {
			return []internal.CarriesByLocality{}, err
		}
	}
	for rows.Next() {
		var row internal.CarriesByLocality
		err := rows.Scan(&row.LocalityId, &row.LocalityName, &row.CarriesCount)
		if err != nil {
			return []internal.CarriesByLocality{}, err
		}
		report = append(report, row)
	}
	return report, nil
}
