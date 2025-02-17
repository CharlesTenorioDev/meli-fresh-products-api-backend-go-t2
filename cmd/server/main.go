package main

import (
	"fmt"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/application"
	"github.com/melisource/fury_go-toolkit-config/v2/pkg/config"
	"github.com/melisource/fury_go-toolkit-secrets/pkg/secrets"
	"gopkg.in/yaml.v3"
)

type properties struct {
	DB struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		EndPoint string `yaml:"endpoint"`
		Name     string `yaml:"name"`
	} `yaml:"db"`
}

// @title			Meli Fresh Products API
// @version		0.0.2.1
// @description	This is a HTTP REST API server.
// @termsOfService	http://swagger.io/terms/
func main() {
	// err := utils.LoadProperties("./.env")
	// if err != nil {
	// 	panic(err)
	// }
	fmt.Println(os.Getenv("GO_ENVIRONMENT"))

	var cfgProps properties
	configBytes, err := config.ReadFromArgs("g2-gow5", []string{"", "--config-dir", "config"})
	if err != nil {
		panic(err)
	}
	fmt.Println(string(configBytes))
	err = yaml.Unmarshal(configBytes, &cfgProps)
	if err != nil {
		panic(err)
	}

	// if os.Getenv("GO_ENVIRONMENT") == "local" {

	// } else
	if os.Getenv("GO_ENVIRONMENT") == "production" {
		client, err := secrets.NewClient()
		if err != nil {
			panic(err)
		}
		var ok = true

		cfgProps.DB.User, ok = client.GetSecret(cfgProps.DB.User)
		if !ok {
			panic("could not find db user secret")
			// secret not found
		}
		cfgProps.DB.Password, ok = client.GetSecret(cfgProps.DB.Password)
		if !ok {
			panic("could not find db password secret")
			// secret not found
		}
		// cfgProps.DB.EndPoint, ok = client.GetSecret(cfgProps.DB.EndPoint)
		// if !ok {
		// 	panic("could not find db endpoint secret")
		// 	// secret not found
		// }
		// cfgProps.DB.Name, ok = client.GetSecret(cfgProps.DB.Name)
		// if !ok {
		// 	panic("could not find db name secret")
		// 	// secret not found
		// }
	}

	// //Create the secret client

	// os.Getenv("DB.USERNAME")
	// os.Getenv("DB.PASSWORD")
	// os.Getenv("DB.DNS")
	// passphrase, ok := client.GetSecret("DB_MYSQL_DESAENV10_BOOTCAMPLEARNINGY_BOOTCAMPLEARNINGY_RPROD")
	// if !ok {
	// 	// secret not found
	// }

	//-config
	// cfg := &application.ConfigApplicationDefault{
	// 	DB: &mysql.Config{
	// 		User:   os.Getenv("DB.USERNAME"),
	// 		Passwd: os.Getenv("DB.PASSWORD"),
	// 		Net:    "tcp",
	// 		Addr:   "localhost" + os.Getenv("DB.ADDRESS"),
	// 		DBName: os.Getenv("DB.NAME"),
	// 	},
	// 	Addr: "127.0.0.1:8080",
	// }
	cfg := &application.ConfigApplicationDefault{
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
