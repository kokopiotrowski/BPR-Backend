package conf

import (
	"fmt"

	"github.com/spf13/viper"
)

type Configurations struct {
	Server   ServerConfigurations   `yaml:"server"`
	Email    EmailConfigurations    `yaml:"email"`
	StockAPI StockAPIConfigurations `yaml:"stockapi"`
}

type ServerConfigurations struct {
	ProdPort       string `yaml:"prodPort"`
	DevPort        string `yaml:"devPort"`
	SSLKey         string `yaml:"sslkey"`
	SSLCertificate string `yaml:"sslcert"`
	VerKey         string `yaml:"verkey"`
}

type EmailConfigurations struct {
	EmailAddress string `yaml:"emailAddress"`
	Password     string `yaml:"password"`
}

type StockAPIConfigurations struct {
	Key string `yaml:"key"`
}

var (
	Conf *Configurations
)

func ReadConfig() (*Configurations, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetConfigType("yaml")

	var configuration Configurations

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// Set undefined variables
	viper.SetDefault("database.dbname", "test_db")

	err := viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}

	// Reading variables using the model
	fmt.Println("Configuration read.")

	Conf = &configuration

	return &configuration, nil
}
