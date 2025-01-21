package main

import (
	"fmt"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/application"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

func main() {
	err := utils.LoadProperties("./.env")
	if err != nil {
		panic(err)
	}

	// - config
	cfg := &application.ConfigApplicationDefault{
		DB: &mysql.Config{
			User:   os.Getenv("DB.USERNAME"),
			Passwd: os.Getenv("DB.PASSWORD"),
			Net:    "tcp",
			Addr:   "localhost" + os.Getenv("DB.ADDRESS"),
			DBName: os.Getenv("DB.NAME"),
		},
		Addr: "127.0.0.1" + os.Getenv("SERVER.PORT"),
	}
	app := application.NewApplicationDefault(cfg)
	// - set up
	err = app.SetUp()
	if err != nil {
		fmt.Println(err)
		return
	}
	// - run
	err = app.Run()
	if err != nil {
		fmt.Println(err)
		return
	}

}
