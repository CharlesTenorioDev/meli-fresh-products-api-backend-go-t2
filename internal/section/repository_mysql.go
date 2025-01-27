package section

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

// NewSectionMysql Initialize a MemorySectionRepository and set the next id to 1
func NewSectionMysql(db *sql.DB) internal.SectionRepository {
	return &SectionMysqlRepository{db}
}

// GetAll the sections and return in asc order
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

func (r *SectionMysqlRepository) GetByID(id int) (internal.Section, error) {
	var section internal.Section

	row := r.db.QueryRow("SELECT s.id, s.section_number, s.current_temperature, s.minimum_temperature, "+
		"s.current_capacity, s.minimum_capacity, s.maximum_capacity, s.warehouse_id, s.product_type_id FROM sections AS s WHERE id=?", id)

	err := row.Scan(&section.ID, &section.SectionNumber, &section.CurrentTemperature,
		&section.MinimumTemperature, &section.CurrentCapacity, &section.MinimumCapacity,
		&section.MaximumCapacity, &section.WarehouseID, &section.ProductTypeID)

	if err != nil {
		if err == sql.ErrNoRows {
			err = utils.ErrNotFound
			return internal.Section{}, err
		}

		return internal.Section{}, err
	}

	return section, nil
}

// GetBySectionNumber Finds the section by its sectionNumber
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
}

// Save Generate a new ID and save the entity
// All validatinos should be made on service layer
func (r *SectionMysqlRepository) Save(newSection *internal.Section) error {
	result, err := r.db.Exec("INSERT INTO sections (section_number, current_temperature, minimum_temperature, current_capacity, minimum_capacity, maximum_capacity, warehouse_id, product_type_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		(*newSection).SectionNumber, (*newSection).CurrentTemperature, (*newSection).MinimumTemperature, (*newSection).CurrentCapacity, (*newSection).MinimumCapacity, (*newSection).MaximumCapacity, (*newSection).WarehouseID, (*newSection).ProductTypeID)

	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			switch mysqlErr.Number {
			case 1062:
				err = utils.ErrConflict
			}

			return err
		}

		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	(*newSection).ID = int(id)

	return err
}

func (r *SectionMysqlRepository) Update(newSection *internal.Section) error {
	_, err := r.db.Exec(
		"UPDATE sections SET section_number=?, current_temperature=?, minimum_temperature=?, current_capacity=?, minimum_capacity=?, maximum_capacity=?, warehouse_id=?, product_type_id=? WHERE id=?",
		(*newSection).SectionNumber, (*newSection).CurrentTemperature, (*newSection).MinimumTemperature,
		(*newSection).CurrentCapacity, (*newSection).MinimumCapacity, (*newSection).MaximumCapacity, (*newSection).WarehouseID,
		(*newSection).ProductTypeID, (*newSection).ID,
	)

	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			switch mysqlErr.Number {
			case 1062:
				err = utils.ErrConflict
			}

			return err
		}

		return err
	}

	return nil
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

func (r *SectionMysqlRepository) GetSectionProductsReport() ([]internal.SectionProductsReport, error) {
	var reports []internal.SectionProductsReport

	rows, err := r.db.Query("SELECT s.id, s.section_number, ifnull(sum(p.current_quantity), 0) as products_count FROM sections s left join product_batches p on s.id = p.section_id group by s.id, s.section_number")

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var report internal.SectionProductsReport

		err = rows.Scan(&report.SectionID, &report.SectionNumber, &report.ProductsCount)
		if err != nil {
			return nil, err
		}

		reports = append(reports, report)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return reports, nil
}
func (r *SectionMysqlRepository) GetSectionProductsReportByID(id int) ([]internal.SectionProductsReport, error) {
	var report internal.SectionProductsReport

	var reports []internal.SectionProductsReport

	row := r.db.QueryRow("SELECT "+
		"s.id, "+
		"s.section_number, "+
		"sum(p.current_quantity) as products_count "+
		"FROM sections s "+
		"left join product_batches p "+
		"on s.id = p.section_id "+
		"where s.id=? group by s.id", id)

	err := row.Scan(&report.SectionID, &report.SectionNumber, &report.ProductsCount)
	if err != nil && err == sql.ErrNoRows {
		err = utils.ErrNotFound
		return nil, err
	}

	reports = append(reports, report)

	return reports, nil
}
