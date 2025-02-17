package main

import (
	"fmt"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/application"
)

// @title			Meli Fresh Products API
// @version		0.4.0
// @description	This is a HTTP REST API api.
// @termsOfService	http://swagger.io/terms/
func main() {
	cfgProps := utils.LoadProperties()
	var cfg *application.ConfigApplicationDefault
	if os.Getenv("GO_ENVIRONMENT") == "production" || cfgProps.GoEnv == "production" {
		cfg = &application.ConfigApplicationDefault{
			DB: &mysql.Config{
				User:                 cfgProps.DB.User,
				Passwd:               cfgProps.DB.Password,
				Net:                  "tcp",
				Addr:                 cfgProps.DB.EndPoint,
				DBName:               cfgProps.DB.Name,
				Timeout:              100 * time.Millisecond,
				ReadTimeout:          100 * time.Millisecond,
				WriteTimeout:         100 * time.Millisecond,
				ParseTime:            true,
				AllowNativePasswords: true,
			},
			Addr: ":8080",
		}
	} else {
		cfg = &application.ConfigApplicationDefault{
			DB: &mysql.Config{
				User:         cfgProps.DB.User,
				Passwd:       cfgProps.DB.Password,
				Net:          "tcp",
				Addr:         cfgProps.DB.EndPoint,
				DBName:       cfgProps.DB.Name,
				Timeout:      100 * time.Millisecond,
				ReadTimeout:  100 * time.Millisecond,
				WriteTimeout: 100 * time.Millisecond,
				ParseTime:    true,
			},
			Addr: "127.0.0.1:8080",
		}
	}
	app := application.NewApplicationDefault(cfg)
	// - set up
	if err := app.SetUp(); err != nil {
		fmt.Println(err)
		return
	}
	// - run
	if err := app.Run(); err != nil {
		fmt.Println(err)
		return
	}
}
