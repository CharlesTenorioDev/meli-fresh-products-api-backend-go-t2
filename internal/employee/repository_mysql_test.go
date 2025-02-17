package employee_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/employee"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
	"github.com/stretchr/testify/require"
)

func TestEmployeeRepository_FindAll(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := employee.NewEmployeeRepository(db)

	t.Run("GIVEN a database error WHEN FindAll is called THEN return an error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, id_card_number, first_name, last_name, warehouse_id FROM employees").
			WillReturnError(errors.New("some internal db error"))

		res, err := repo.FindAll()

		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("GIVEN a valid query WHEN executing FindAll THEN return employees", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "id_card_number", "first_name", "last_name", "warehouse_id"}).
			AddRow(1, "123ABC", "Aelin", "Galanthinius", 10).
			AddRow(2, "456DEF", "Rowan", "Whitethorn", 20)

		mock.ExpectQuery("SELECT id, id_card_number, first_name, last_name, warehouse_id FROM employees").
			WillReturnRows(rows)

		res, err := repo.FindAll()

		require.NoError(t, err)
		require.Len(t, res, 2)
		require.Equal(t, "Aelin", res[1].Attributes.FirstName)
		require.Equal(t, "Rowan", res[2].Attributes.FirstName)
	})

	t.Run("GIVEN a database error WHEN executing FindAll THEN return an error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, id_card_number, first_name, last_name, warehouse_id FROM employees").
			WillReturnError(errors.New("db error"))

		res, err := repo.FindAll()

		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("GIVEN an invalid row scan WHEN executing FindAll THEN return an error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "id_card_number", "first_name", "last_name", "warehouse_id"}).
			AddRow(1, nil, "Aelin", "Galanthinius", 10)

		mock.ExpectQuery("SELECT id, id_card_number, first_name, last_name, warehouse_id FROM employees").
			WillReturnRows(rows)

		res, err := repo.FindAll()

		require.Error(t, err)
		require.Nil(t, res)
	})
}

func TestEmployeeRepository_FindByID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := employee.NewEmployeeRepository(db)

	t.Run("GIVEN a valid ID WHEN FindByID is called THEN return the employee", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, id_card_number, first_name, last_name, warehouse_id FROM employees WHERE id = ?").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "id_card_number", "first_name", "last_name", "warehouse_id"}).
				AddRow(1, "1234", "Aelin", "Galanthinius", 10))

		emp, err := repo.FindByID(1)

		require.NoError(t, err)
		require.Equal(t, 1, emp.ID)
	})

	t.Run("GIVEN a non-existent ID WHEN FindByID is called THEN return ErrNotFound", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, id_card_number, first_name, last_name, warehouse_id FROM employees WHERE id = ?").
			WithArgs(99).
			WillReturnError(sql.ErrNoRows)

		emp, err := repo.FindByID(99)

		require.ErrorIs(t, err, utils.ErrNotFound)
		require.Zero(t, emp)
	})

	t.Run("GIVEN a database error WHEN executing FindByID THEN return an error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, card_number_id, first_name, last_name, warehouse_id FROM employees WHERE id = ?").
			WithArgs(1).
			WillReturnError(errors.New("db error"))

		res, err := repo.FindByID(1)

		require.Error(t, err)
		require.Equal(t, internal.Employee{}, res)
	})
}

func TestEmployeeRepository_CreateEmployee(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := employee.NewEmployeeRepository(db)

	newEmp := internal.EmployeeAttributes{
		CardNumberID: "5678",
		FirstName:    "Celaena",
		LastName:     "Sardothien",
		WarehouseID:  20,
	}

	t.Run("GIVEN a valid employee WHEN CreateEmployee is called THEN insert and return the employee", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO employees").
			WithArgs(newEmp.CardNumberID, newEmp.FirstName, newEmp.LastName, newEmp.WarehouseID).
			WillReturnResult(sqlmock.NewResult(2, 1))

		mock.ExpectQuery("SELECT id, id_card_number, first_name, last_name, warehouse_id FROM employees WHERE id = ?").
			WithArgs(2).
			WillReturnRows(sqlmock.NewRows([]string{"id", "id_card_number", "first_name", "last_name", "warehouse_id"}).
				AddRow(2, "5678", "Celaena", "Sardothien", 20))

		emp, err := repo.CreateEmployee(newEmp)

		require.NoError(t, err)
		require.Equal(t, 2, emp.ID)
	})

	t.Run("GIVEN a database error WHEN executing CreateEmployee THEN return an error", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO employees").
			WithArgs("12345", "Aelin", "Galanthinius", 1).
			WillReturnError(errors.New("db error"))

		newEmp := internal.EmployeeAttributes{
			CardNumberID: "12345",
			FirstName:    "Aelin",
			LastName:     "Galanthinius",
			WarehouseID:  1,
		}

		res, err := repo.CreateEmployee(newEmp)

		require.Error(t, err)
		require.Equal(t, internal.Employee{}, res)
	})

	t.Run("GIVEN a database error WHEN retrieving LastInsertId THEN return an error", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO employees").
			WithArgs("12345", "John", "Doe", 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectQuery("SELECT LAST_INSERT_ID()").WillReturnError(errors.New("last insert id error"))

		newEmp := internal.EmployeeAttributes{
			CardNumberID: "12345",
			FirstName:    "Aelin",
			LastName:     "Galanthinius",
			WarehouseID:  1,
		}

		res, err := repo.CreateEmployee(newEmp)

		require.Error(t, err)
		require.Equal(t, internal.Employee{}, res)
	})

}

func TestEmployeeRepository_UpdateEmployee(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := employee.NewEmployeeRepository(db)

	updatedEmp := internal.Employee{
		ID: 1,
		Attributes: internal.EmployeeAttributes{
			CardNumberID: "9999",
			FirstName:    "Updated",
			LastName:     "Name",
			WarehouseID:  30,
		},
	}

	t.Run("GIVEN an existing employee WHEN UpdateEmployee is called THEN update and return the employee", func(t *testing.T) {
		mock.ExpectExec("UPDATE employees").
			WithArgs(updatedEmp.Attributes.CardNumberID, updatedEmp.Attributes.FirstName, updatedEmp.Attributes.LastName, updatedEmp.Attributes.WarehouseID, updatedEmp.ID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectQuery("SELECT id, id_card_number, first_name, last_name, warehouse_id FROM employees WHERE id = ?").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "id_card_number", "first_name", "last_name", "warehouse_id"}).
				AddRow(1, "9999", "Updated", "Name", 30))

		emp, err := repo.UpdateEmployee(updatedEmp)

		require.NoError(t, err)
		require.Equal(t, "9999", emp.Attributes.CardNumberID)
	})
}

func TestEmployeeRepository_DeleteEmployee(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := employee.NewEmployeeRepository(db)

	t.Run("GIVEN an existing employee WHEN DeleteEmployee is called THEN remove the employee", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM employees WHERE id = ?").
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.DeleteEmployee(1)

		require.NoError(t, err)
	})
}
