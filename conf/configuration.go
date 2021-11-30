package conf

import (
	"fmt"

	"github.com/spf13/viper"
)

type Configurations struct {
	Server ServerConfigurations
	Email  EmailConfigurations
}

type ServerConfigurations struct {
	ProdPort string
	DevPort  string
}

type EmailConfigurations struct {
	EmailAddress string
	Password     string
}

var (
	Conf Configurations
)

func ReadConfig() (*Configurations, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetConfigType("yml")

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
	fmt.Println("Reading variables using the model..")
	fmt.Println("Production port is\t", configuration.Server.ProdPort)
	fmt.Println("Development port is\t\t", configuration.Server.DevPort)
	fmt.Println("Email address is\t", configuration.Email.EmailAddress)

	Conf = configuration

	return &configuration, nil
}
