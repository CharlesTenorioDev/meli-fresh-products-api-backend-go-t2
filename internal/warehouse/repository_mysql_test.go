package warehouse_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/warehouse"
	"github.com/stretchr/testify/assert"
)

func TestSave_PrepareError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := warehouse.NewWarehouseDB(db)

	mock.ExpectPrepare("INSERT INTO warehouses").WillReturnError(errors.New("prepare error"))

	wh := internal.Warehouse{Address: "Address 1", Telephone: "123456", WarehouseCode: "WH001", LocalityID: 10, MinimumCapacity: 100, MinimumTemperature: -5}
	_, err := repo.Save(wh)

	assert.Error(t, err)
	assert.Equal(t, "prepare error", err.Error())
}

func TestSave_ExecError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := warehouse.NewWarehouseDB(db)

	mock.ExpectPrepare("INSERT INTO warehouses").ExpectExec().
		WithArgs("Address 1", "123456", "WH001", 10, 100, -5).
		WillReturnError(errors.New("exec error"))

	wh := internal.Warehouse{Address: "Address 1", Telephone: "123456", WarehouseCode: "WH001", LocalityID: 10, MinimumCapacity: 100, MinimumTemperature: -5}
	_, err := repo.Save(wh)

	assert.Error(t, err)
	assert.Equal(t, "exec error", err.Error())
}

func TestSave_LastInsertIdError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := warehouse.NewWarehouseDB(db)

	mock.ExpectPrepare("INSERT INTO warehouses").ExpectExec().
		WithArgs("Address 1", "123456", "WH001", 10, 100, -5).
		WillReturnResult(sqlmock.NewErrorResult(errors.New("last insert id error")))

	wh := internal.Warehouse{Address: "Address 1", Telephone: "123456", WarehouseCode: "WH001", LocalityID: 10, MinimumCapacity: 100, MinimumTemperature: -5}
	_, err := repo.Save(wh)

	assert.Error(t, err)
	assert.Equal(t, "last insert id error", err.Error())
}

func TestSave_MySQLError1062(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := warehouse.NewWarehouseDB(db)

	mysqlErr := &mysql.MySQLError{Number: 1062, Message: "Duplicate entry"}
	mock.ExpectPrepare("INSERT INTO warehouses").ExpectExec().
		WithArgs("Address 1", "123456", "WH001", 10, 100, -5).
		WillReturnError(mysqlErr)

	wh := internal.Warehouse{Address: "Address 1", Telephone: "123456", WarehouseCode: "WH001", LocalityID: 10, MinimumCapacity: 100, MinimumTemperature: -5}
	_, err := repo.Save(wh)

	assert.ErrorIs(t, err, utils.ErrConflict)
}

func TestSave_MySQLErrorGeneric(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := warehouse.NewWarehouseDB(db)

	mysqlErr := &mysql.MySQLError{Number: 1234, Message: "Some MySQL error"}
	mock.ExpectPrepare("INSERT INTO warehouses").ExpectExec().
		WithArgs("Address 1", "123456", "WH001", 10, 100, -5).
		WillReturnError(mysqlErr)

	wh := internal.Warehouse{Address: "Address 1", Telephone: "123456", WarehouseCode: "WH001", LocalityID: 10, MinimumCapacity: 100, MinimumTemperature: -5}
	_, err := repo.Save(wh)

	assert.Error(t, err)
	assert.Equal(t, mysqlErr, err)
}

func TestSave_AssignID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := warehouse.NewWarehouseDB(db)

	expectedID := int64(42)
	mock.ExpectPrepare("INSERT INTO warehouses").ExpectExec().
		WithArgs("Address 1", "123456", "WH001", 10, 100, -5).
		WillReturnResult(sqlmock.NewResult(expectedID, 1))

	wh := internal.Warehouse{Address: "Address 1", Telephone: "123456", WarehouseCode: "WH001", LocalityID: 10, MinimumCapacity: 100, MinimumTemperature: -5}
	savedWh, err := repo.Save(wh)

	assert.NoError(t, err)
	assert.Equal(t, int(expectedID), savedWh.ID)
}

func TestGetAll_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := warehouse.NewWarehouseDB(db)

	mock.ExpectQuery("SELECT `id`, `address`, `telephone`, `warehouse_code`, `locality_id`, `minimum_capacity`, `minimum_temperature` FROM warehouses").
		WillReturnRows(sqlmock.NewRows([]string{"id", "address", "telephone", "warehouse_code", "locality_id", "minimum_capacity", "minimum_temperature"}).
			AddRow(1, "Address 1", "123456", "WH001", 10, 100, -5).
			AddRow(2, "Address 2", "654321", "WH002", 20, 200, -10))

	warehouses, err := repo.GetAll()

	assert.NoError(t, err)
	assert.Len(t, warehouses, 2)
	assert.Equal(t, 1, warehouses[0].ID)
	assert.Equal(t, "Address 1", warehouses[0].Address)
	assert.Equal(t, 2, warehouses[1].ID)
	assert.Equal(t, "Address 2", warehouses[1].Address)
}

func TestGetAll_QueryError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := warehouse.NewWarehouseDB(db)

	mock.ExpectQuery("SELECT `id`, `address`, `telephone`, `warehouse_code`, `locality_id`, `minimum_capacity`, `minimum_temperature` FROM warehouses").
		WillReturnError(errors.New("query error"))

	warehouses, err := repo.GetAll()

	assert.Error(t, err)
	assert.Nil(t, warehouses)
	assert.Equal(t, "query error", err.Error())
}

func TestGetAll_NoRows(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := warehouse.NewWarehouseDB(db)

	mock.ExpectQuery("SELECT `id`, `address`, `telephone`, `warehouse_code`, `locality_id`, `minimum_capacity`, `minimum_temperature` FROM warehouses").
		WillReturnRows(sqlmock.NewRows([]string{"id", "address", "telephone", "warehouse_code", "locality_id", "minimum_capacity", "minimum_temperature"}))

	warehouses, err := repo.GetAll()

	assert.NoError(t, err)
	assert.Len(t, warehouses, 0)
}

func TestGetAll_RowsCloseError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := warehouse.NewWarehouseDB(db)

	mock.ExpectQuery("SELECT `id`, `address`, `telephone`, `warehouse_code`, `locality_id`, `minimum_capacity`, `minimum_temperature` FROM warehouses").
		WillReturnRows(sqlmock.NewRows([]string{"id", "address", "telephone", "warehouse_code", "locality_id", "minimum_capacity", "minimum_temperature"}).
			AddRow(1, "Address 1", "123456", "WH001", 10, 100, -5).
			CloseError(errors.New("close error")))

	warehouses, err := repo.GetAll()

	assert.Error(t, err)
	assert.Nil(t, warehouses)
	assert.Equal(t, "close error", err.Error())
}

func TestGetAll_ScanError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := warehouse.NewWarehouseDB(db)

	mock.ExpectQuery("SELECT `id`, `address`, `telephone`, `warehouse_code`, `locality_id`, `minimum_capacity`, `minimum_temperature` FROM warehouses").
		WillReturnRows(sqlmock.NewRows([]string{"id", "address", "telephone", "warehouse_code", "locality_id", "minimum_capacity", "minimum_temperature"}).
			AddRow("invalid-id", "Address 1", "123456", "WH001", 10, 100, -5)) // "invalid-id" é um valor inválido para o campo ID (esperado int)

	warehouses, err := repo.GetAll()

	assert.Error(t, err)
	assert.Nil(t, warehouses)
	assert.Contains(t, err.Error(), "converting") // Verifica se o erro é relacionado à conversão de tipos
}

func TestGetByID_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := warehouse.NewWarehouseDB(db)

	mock.ExpectQuery("SELECT `id`, `address`, `telephone`, `warehouse_code`, `locality_id`, `minimum_capacity`, `minimum_temperature` FROM warehouses WHERE id = ?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "address", "telephone", "warehouse_code", "locality_id", "minimum_capacity", "minimum_temperature"}).
			AddRow(1, "Address 1", "123456", "WH001", 10, 100, -5))

	warehouse, err := repo.GetByID(1)

	assert.NoError(t, err)
	assert.Equal(t, 1, warehouse.ID)
	assert.Equal(t, "Address 1", warehouse.Address)
	assert.Equal(t, "123456", warehouse.Telephone)
	assert.Equal(t, "WH001", warehouse.WarehouseCode)
	assert.Equal(t, 10, warehouse.LocalityID)
	assert.Equal(t, 100, warehouse.MinimumCapacity)
	assert.Equal(t, -5, warehouse.MinimumTemperature)
}

func TestGetByID_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := warehouse.NewWarehouseDB(db)

	mock.ExpectQuery("SELECT `id`, `address`, `telephone`, `warehouse_code`, `locality_id`, `minimum_capacity`, `minimum_temperature` FROM warehouses WHERE id = ?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{}).RowError(0, errors.New("sql: no rows in result set")))

	warehouse, err := repo.GetByID(1)

	assert.Error(t, err)
	assert.Equal(t, internal.Warehouse{}, warehouse)
	assert.Equal(t, "entity not found", err.Error())

	// Verifica se todas as expectativas foram atendidas
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByID_QueryError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := warehouse.NewWarehouseDB(db)

	mock.ExpectQuery("SELECT `id`, `address`, `telephone`, `warehouse_code`, `locality_id`, `minimum_capacity`, `minimum_temperature` FROM warehouses WHERE id = ?").
		WithArgs(1).
		WillReturnError(errors.New("query error"))

	warehouse, err := repo.GetByID(1)

	assert.Error(t, err)
	assert.Equal(t, internal.Warehouse{}, warehouse)
	assert.Equal(t, "query error", err.Error())
}

func TestGetByID_ScanError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := warehouse.NewWarehouseDB(db)

	mock.ExpectQuery("SELECT `id`, `address`, `telephone`, `warehouse_code`, `locality_id`, `minimum_capacity`, `minimum_temperature` FROM warehouses WHERE id = ?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "address", "telephone", "warehouse_code", "locality_id", "minimum_capacity", "minimum_temperature"}).
			AddRow("invalid-id", "Address 1", "123456", "WH001", 10, 100, -5))

	warehouse, err := repo.GetByID(1)

	assert.Error(t, err)
	assert.Equal(t, internal.Warehouse{}, warehouse)
	assert.Contains(t, err.Error(), "converting") // Verifica se o erro é relacionado à conversão de tipos
}

func TestUpdate_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := warehouse.NewWarehouseDB(db)

	mock.ExpectQuery("SELECT `id`, `address`, `telephone`, `warehouse_code`, `locality_id`, `minimum_capacity`, `minimum_temperature` FROM warehouses WHERE id = ?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "address", "telephone", "warehouse_code", "locality_id", "minimum_capacity", "minimum_temperature"}).
			AddRow(1, "Old Address", "123456", "WH001", 10, 100, -5))

	mock.ExpectPrepare("UPDATE `warehouses` AS `w` SET `address` = \\?, `telephone` = \\?, `warehouse_code` = \\?, `locality_id` = \\?, `minimum_capacity`= \\?, `minimum_temperature`= \\? WHERE `id` = \\?")

	mock.ExpectExec("UPDATE `warehouses` AS `w` SET `address` = \\?, `telephone` = \\?, `warehouse_code` = \\?, `locality_id` = \\?, `minimum_capacity`= \\?, `minimum_temperature`= \\? WHERE `id` = \\?").
		WithArgs("New Address", "654321", "WH001", 10, 200, -10, 1).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 1 linha afetada

	// Dados atualizados do warehouse
	updatedWarehouse := internal.Warehouse{
		ID:                 1,
		Address:            "New Address",
		Telephone:          "654321",
		WarehouseCode:      "WH001",
		LocalityID:         10,
		MinimumCapacity:    200,
		MinimumTemperature: -10,
	}

	// Executa a função Update
	result, err := repo.Update(updatedWarehouse)

	// Verifica se não houve erro
	assert.NoError(t, err)
	assert.Equal(t, updatedWarehouse, result)

	// Verifica se todas as expectativas foram atendidas
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate_Prepare_Error(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := warehouse.NewWarehouseDB(db)

	mock.ExpectQuery("SELECT `id`, `address`, `telephone`, `warehouse_code`, `locality_id`, `minimum_capacity`, `minimum_temperature` FROM warehouses WHERE id = ?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "address", "telephone", "warehouse_code", "locality_id", "minimum_capacity", "minimum_temperature"}).
			AddRow(1, "Old Address", "123456", "WH001", 10, 100, -5))

	mock.ExpectPrepare("UPDATE `warehouses` AS `w` SET `address` = \\?, `telephone` = \\?, `warehouse_code` = \\?, `locality_id` = \\?, `minimum_capacity`= \\?, `minimum_temperature`= \\? WHERE `id` = \\?").
		WillReturnError(errors.New("prepare error"))

	// Dados do warehouse
	updatedWarehouse := internal.Warehouse{
		ID: 1,
	}

	// Executa a função Update
	result, err := repo.Update(updatedWarehouse)

	// Verifica se o erro é retornado corretamente
	assert.Error(t, err)
	assert.Equal(t, internal.Warehouse{}, result)
	assert.Equal(t, "prepare error", err.Error())

	// Verifica se todas as expectativas foram atendidas
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate_Exec_Error(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := warehouse.NewWarehouseDB(db)

	mock.ExpectQuery("SELECT `id`, `address`, `telephone`, `warehouse_code`, `locality_id`, `minimum_capacity`, `minimum_temperature` FROM warehouses WHERE id = ?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "address", "telephone", "warehouse_code", "locality_id", "minimum_capacity", "minimum_temperature"}).
			AddRow(1, "Old Address", "123456", "WH001", 10, 100, -5))

	mock.ExpectPrepare("UPDATE `warehouses` AS `w` SET `address` = \\?, `telephone` = \\?, `warehouse_code` = \\?, `locality_id` = \\?, `minimum_capacity`= \\?, `minimum_temperature`= \\? WHERE `id` = \\?")

	mock.ExpectExec("UPDATE `warehouses` AS `w` SET `address` = \\?, `telephone` = \\?, `warehouse_code` = \\?, `locality_id` = \\?, `minimum_capacity`= \\?, `minimum_temperature`= \\? WHERE `id` = \\?").
		WithArgs("New Address", "654321", "WH001", 10, 200, -10, 1).
		WillReturnError(errors.New("exec error"))

	updatedWarehouse := internal.Warehouse{
		ID:                 1,
		Address:            "New Address",
		Telephone:          "654321",
		WarehouseCode:      "WH001",
		LocalityID:         10,
		MinimumCapacity:    200,
		MinimumTemperature: -10,
	}

	result, err := repo.Update(updatedWarehouse)

	assert.Error(t, err)
	assert.Equal(t, internal.Warehouse{}, result)
	assert.Equal(t, "exec error", err.Error())

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate_Conflict_Error(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := warehouse.NewWarehouseDB(db)

	mock.ExpectQuery("SELECT `id`, `address`, `telephone`, `warehouse_code`, `locality_id`, `minimum_capacity`, `minimum_temperature` FROM warehouses WHERE id = ?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "address", "telephone", "warehouse_code", "locality_id", "minimum_capacity", "minimum_temperature"}).
			AddRow(1, "Old Address", "123456", "WH001", 10, 100, -5))

	mock.ExpectPrepare("UPDATE `warehouses` AS `w` SET `address` = \\?, `telephone` = \\?, `warehouse_code` = \\?, `locality_id` = \\?, `minimum_capacity`= \\?, `minimum_temperature`= \\? WHERE `id` = \\?")

	mock.ExpectExec("UPDATE `warehouses` AS `w` SET `address` = \\?, `telephone` = \\?, `warehouse_code` = \\?, `locality_id` = \\?, `minimum_capacity`= \\?, `minimum_temperature`= \\? WHERE `id` = \\?").
		WithArgs("New Address", "654321", "WH001", 10, 200, -10, 1).
		WillReturnError(&mysql.MySQLError{Number: 1062, Message: "Duplicate entry"})

	updatedWarehouse := internal.Warehouse{
		ID:                 1,
		Address:            "New Address",
		Telephone:          "654321",
		WarehouseCode:      "WH001",
		LocalityID:         10,
		MinimumCapacity:    200,
		MinimumTemperature: -10,
	}

	result, err := repo.Update(updatedWarehouse)

	assert.Error(t, err)
	assert.Equal(t, internal.Warehouse{}, result)
	assert.True(t, errors.Is(err, utils.ErrConflict))

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate_GetByID_GenericError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := warehouse.NewWarehouseDB(db)

	mock.ExpectQuery("SELECT `id`, `address`, `telephone`, `warehouse_code`, `locality_id`, `minimum_capacity`, `minimum_temperature` FROM warehouses WHERE id = ?").
		WithArgs(1).
		WillReturnError(errors.New("generic error"))

	updatedWarehouse := internal.Warehouse{
		ID: 1,
	}

	result, err := repo.Update(updatedWarehouse)

	assert.Error(t, err)
	assert.Equal(t, internal.Warehouse{}, result)
	assert.Equal(t, "generic error", err.Error())

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDelete_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := warehouse.NewWarehouseDB(db)

	mock.ExpectQuery("SELECT `id`, `address`, `telephone`, `warehouse_code`, `locality_id`, `minimum_capacity`, `minimum_temperature` FROM warehouses WHERE id = ?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "address", "telephone", "warehouse_code", "locality_id", "minimum_capacity", "minimum_temperature"}).
			AddRow(1, "Address 1", "123456", "WH001", 10, 100, -5))

	mock.ExpectPrepare("DELETE FROM warehouses WHERE id = \\?")

	mock.ExpectExec("DELETE FROM warehouses WHERE id = \\?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 1 linha afetada

	err := repo.Delete(1)

	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDelete_GetByID_Error(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := warehouse.NewWarehouseDB(db)

	mock.ExpectQuery("SELECT `id`, `address`, `telephone`, `warehouse_code`, `locality_id`, `minimum_capacity`, `minimum_temperature` FROM warehouses WHERE id = ?").
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)

	err := repo.Delete(1)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDelete_Prepare_Error(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := warehouse.NewWarehouseDB(db)

	mock.ExpectQuery("SELECT `id`, `address`, `telephone`, `warehouse_code`, `locality_id`, `minimum_capacity`, `minimum_temperature` FROM warehouses WHERE id = ?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "address", "telephone", "warehouse_code", "locality_id", "minimum_capacity", "minimum_temperature"}).
			AddRow(1, "Address 1", "123456", "WH001", 10, 100, -5))

	mock.ExpectPrepare("DELETE FROM warehouses WHERE id = \\?").
		WillReturnError(errors.New("prepare error"))

	err := repo.Delete(1)

	assert.Error(t, err)
	assert.Equal(t, "prepare error", err.Error())

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDelete_Exec_Error(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := warehouse.NewWarehouseDB(db)

	mock.ExpectQuery("SELECT `id`, `address`, `telephone`, `warehouse_code`, `locality_id`, `minimum_capacity`, `minimum_temperature` FROM warehouses WHERE id = ?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "address", "telephone", "warehouse_code", "locality_id", "minimum_capacity", "minimum_temperature"}).
			AddRow(1, "Address 1", "123456", "WH001", 10, 100, -5))

	mock.ExpectPrepare("DELETE FROM warehouses WHERE id = \\?")

	mock.ExpectExec("DELETE FROM warehouses WHERE id = \\?").
		WithArgs(1).
		WillReturnError(errors.New("exec error"))

	err := repo.Delete(1)

	assert.Error(t, err)
	assert.Equal(t, "exec error", err.Error())

	assert.NoError(t, mock.ExpectationsWereMet())
}
