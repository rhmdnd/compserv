package compserv

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

func ParseConfig(configDir string) map[string]string {
	viper.SetDefault("database.port", "5432")
	viper.SetDefault("database.name", "compliance")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error reading config file: %w \n", err))
	}
	db_host := viper.GetString("database.host")
	db_port := viper.GetString("database.port")
	db_username := viper.GetString("database.username")
	db_name := viper.GetString("database.name")
	db_secret_arn := viper.GetString("database.secret_arn")
	db_secret_region := viper.GetString("database.secret_region")

	if len(db_host) == 0 {
		log.Fatal("Database host not provided.")
		os.Exit(1)
	} else if len(db_secret_arn) == 0 {
		log.Fatal("Database password not provided as an ARN.")
		os.Exit(1)
	} else if len(db_username) == 0 {
		log.Fatal("Database username not provided.")
		os.Exit(1)
	} else if len(db_secret_region) == 0 {
		log.Fatal("Database secret region not provided.")
		os.Exit(1)
	}

	m := make(map[string]string)
	m["db_host"] = db_host
	m["db_port"] = db_port
	m["db_secret_arn"] = db_secret_arn
	m["db_username"] = db_username
	m["db_name"] = db_name
	m["secret_region"] = db_secret_region

	log.Printf("Loaded configuration file: %s", viper.ConfigFileUsed())
	return m
}
