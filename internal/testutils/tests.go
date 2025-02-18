package testutils

import (
	"database/sql"
	"fmt"

	"github.com/DATA-DOG/go-txdb"
	_ "github.com/go-sql-driver/mysql"
)

type DBTest struct {
	DB *sql.DB
}

func (t *DBTest) End() {
	t.DB.Close()
}

func GetTestDBConn() *DBTest {
	User := "root"
	Password := "test"
	Database := "fresh_products"
	Port2 := "3306"

	txdb.Register("txdb", "mysql", fmt.Sprintf("%s:%s@tcp(localhost:%s)/%s?allowNativePasswords=false&checkConnLiveness=false&maxAllowedPacket=0",
		User,
		Password,
		Port2,
		Database,
	))

	db, err := sql.Open("txdb", "fantasy_products")
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return &DBTest{db}
}
