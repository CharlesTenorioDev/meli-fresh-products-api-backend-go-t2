package repository

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type SectionMysqlRepository struct {
	db *sql.DB
}

// Initialize a MemorySectionRepository and set the next id to 1
func NewSectionMysql(db *sql.DB) *SectionMysqlRepository {
	return &SectionMysqlRepository{db}
}

// Get all the sections and return in asc order
func (r SectionMysqlRepository) GetAll() ([]internal.Section, error) {
	var sections []internal.Section

	rows, err := r.db.Query("SELECT s.id, s.section_number, s.current_temperature, s.minimum_temperature, " +
		"s.current_capacity, s.minimum_capacity, s.maximum_capacity, s.warehouse_id, s.product_type_id FROM sections AS s")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var section internal.Section
		err = rows.Scan(&section.ID, &section.SectionNumber, &section.CurrentTemperature,
			&section.MinimumTemperature, &section.CurrentCapacity, &section.MinimumCapacity,
			&section.MaximumCapacity, &section.WarehouseID, &section.ProductTypeID)

		if err != nil {
			return nil, err
		}

		sections = append(sections, section)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return sections, nil
}

func (r *SectionMysqlRepository) GetById(id int) (internal.Section, error) {
	var section internal.Section

	row := r.db.QueryRow("SELECT s.id, s.section_number, s.current_temperature, s.minimum_temperature, "+
		"s.current_capacity, s.minimum_capacity, s.maximum_capacity, s.warehouse_id, s.product_type_id FROM sections AS s WHERE id=?", id)

	err := row.Scan(&section.ID, &section.SectionNumber, &section.CurrentTemperature,
		&section.MinimumTemperature, &section.CurrentCapacity, &section.MinimumCapacity,
		&section.MaximumCapacity, &section.WarehouseID, &section.ProductTypeID)

	if err != nil {
		if err == sql.ErrNoRows {
			err = utils.ErrNotFound
		}
		return internal.Section{}, err
	}

	return section, nil

}

// Finds the section by its sectionNumber
// If not found, pkg.Section{} is returned
func (r *SectionMysqlRepository) GetBySectionNumber(sectionNumber int) (internal.Section, error) {

	var section internal.Section

	row := r.db.QueryRow("SELECT s.id, s.section_number, s.current_temperature, s.minimum_temperature, "+
		"s.current_capacity, s.minimum_capacity, s.maximum_capacity, s.warehouse_id, s.product_type_id FROM sections AS s WHERE section_number=?", sectionNumber)

	err := row.Scan(&section.ID, &section.SectionNumber, &section.CurrentTemperature,
		&section.MinimumTemperature, &section.CurrentCapacity, &section.MinimumCapacity,
		&section.MaximumCapacity, &section.WarehouseID, &section.ProductTypeID)

	if err != nil {
		if err == sql.ErrNoRows {
			err = utils.ErrNotFound
		}
		return internal.Section{}, err
	}

	return section, nil

	//for _, section := range r.db {
	//	if section.SectionNumber == sectionNumber {
	//		return section, nil
	//	}
	//}
	//return internal.Section{}, nil
}

// Generate a new ID and save the entity
// All validatinos should be made on service layer
func (r *SectionMysqlRepository) Save(newSection *internal.Section) (internal.Section, error) {

	result, err := r.db.Exec(
		"INSERT INTO sections (section_number, current_temperature, minimum_temperature, current_capacity, minimum_capacity, maximum_capacity, warehouse_id, product_type_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		(*newSection).SectionNumber, (*newSection).SectionNumber, (*newSection).CurrentTemperature, (*newSection).MinimumTemperature, (*newSection).CurrentCapacity, (*newSection).MinimumCapacity, (*newSection).MaximumCapacity, (*newSection).WarehouseID, (*newSection).ProductTypeID,
	)

	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			switch mysqlErr.Number {
			case 1062:
				err = utils.ErrConflict
			}
			return internal.Section{}, err
		}
		return internal.Section{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return internal.Section{}, err
	}

	(*newSection).ID = int(id)

	return *newSection, err

}

// Update the map
// If a section does not exist for the id, a new one will be created
func (r *SectionMysqlRepository) Update(section internal.Section) (internal.Section, error) {
	//r.db[section.ID] = section
	//return r.db[section.ID], nil

	return internal.Section{}, nil
}

// Delete the section by its id
// If no sections exists by id attribute, utils.ErrNotFound is returned
func (r *SectionMysqlRepository) Delete(id int) error {

	_, err := r.db.Exec("DELETE FROM sections WHERE id=?", id)

	if err != nil {
		return err
	}
	return nil
}
