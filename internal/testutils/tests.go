package testutils

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"path/filepath"

	"github.com/DATA-DOG/go-txdb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

type DBTest struct {
	DB        *sql.DB
	Container *mysql.MySQLContainer
}

func (t *DBTest) End() {
	t.DB.Close()
	if err := testcontainers.TerminateContainer(t.Container); err != nil {
		log.Printf("failed to terminate container: %s", err)
	}
}

func GetTestDBConn() *DBTest {
	User := "foo"
	Password := "bar"
	Database := "fresh_products"
	container, err := mysql.Run(context.Background(),
		"mysql:8.0.36",
		mysql.WithDatabase(Database),
		mysql.WithUsername(User),
		mysql.WithPassword(Password),
		mysql.WithScripts(filepath.Join("../", "../", "docs", "db", "mysql", "create.sql")),
	)
	if err != nil {
		panic(err)
	}
	NatPort, err := container.MappedPort(context.Background(), "3306/tcp")
	if err != nil {
		fmt.Println(err.Error())
	}
	Port2 := NatPort.Port()

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
	return &DBTest{db, container}
}
