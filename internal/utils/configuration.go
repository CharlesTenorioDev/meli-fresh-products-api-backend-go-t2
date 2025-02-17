package utils

import (
	"fmt"
	"github.com/melisource/fury_go-toolkit-config/v2/pkg/config"
	"github.com/melisource/fury_go-toolkit-secrets/pkg/secrets"
	"gopkg.in/yaml.v3"
	"os"
)

type Properties struct {
	GoEnv string `yaml:"goenv"`
	DB    struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		EndPoint string `yaml:"endpoint"`
		Name     string `yaml:"name"`
	} `yaml:"db"`
}

func LoadProperties() *Properties {
	var cfgProps Properties
	configBytes, err := config.ReadFromArgs("g2-gow5", []string{"", "--config-dir", "config"})
	if err != nil {
		fmt.Println("Could not ReadFromArgs ")
		configBytes, err = config.Read("g2-gow5")
		if err != nil {
			fmt.Println("Could not Read ")
			panic(err)
		}
	}
	fmt.Println(string(configBytes))
	err = yaml.Unmarshal(configBytes, &cfgProps)
	if err != nil {
		panic(err)
	}

	if os.Getenv("GO_ENVIRONMENT") == "production" || cfgProps.GoEnv == "production" {
		client, err := secrets.NewClient()
		if err != nil {
			panic(err)
		}
		var ok = true

		cfgProps.DB.User, ok = client.GetSecret(cfgProps.DB.User)
		if !ok {
			panic("could not find db user secret")
		}
		cfgProps.DB.Password, ok = client.GetSecret(cfgProps.DB.Password)
		if !ok {
			panic("could not find db password secret")
		}
		fmt.Printf("LEN %d, %d", len(cfgProps.DB.User), len(cfgProps.DB.Password))
		fmt.Println(os.Getenv(cfgProps.DB.EndPoint))
		cfgProps.DB.EndPoint = os.Getenv(cfgProps.DB.EndPoint)
	}
	return &cfgProps
}
