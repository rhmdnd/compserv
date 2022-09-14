package compserv

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

func ParseConfig(configDir, configFile string) *viper.Viper {
	viper.SetDefault("app.host", "localhost")
	viper.SetDefault("app.port", "50051")
	viper.SetDefault("database.port", "5432")
	viper.SetDefault("database.name", "compliance")
	configType := "yaml"
	parts := strings.Split(configFile, ".")
	if ln := len(parts); ln > 1 {
		configType = parts[ln-1]
	}
	viper.SetConfigName(configFile)
	viper.SetConfigType(configType)
	viper.AddConfigPath(configDir)
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error reading config file: %w", err))
	}

	v := viper.GetViper()
	validateConfig(v)
	log.Printf("Loaded configuration file: %s", v.ConfigFileUsed())
	return v
}

func GetDatabaseConnectionString(v *viper.Viper) string {
	p := v.GetString("database.password.provider")
	var s string
	switch p {
	case "aws":
		s = getSecretFromAws(v)
	case "kubernetes":
		s = getSecretFromKubernetes(v)
	}

	connectionString := fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%s",
		v.GetString("database.host"), v.GetString("database.username"),
		v.GetString("database.name"), s, v.GetString("database.port"))
	return connectionString
}

// Ensure we can use the provided configuration
func validateConfig(v *viper.Viper) {
	if v.GetString("database.host") == "" {
		log.Fatal("Database host not provided (database.host)")
	}
	if v.GetString("database.username") == "" {
		log.Fatal("Database username not provided (database.username)")
	}

	p := v.GetString("database.password.provider")
	switch p {
	case "aws":
		if v.GetString("database.password.secret_arn") == "" {
			log.Fatal("Database password not provided as a secret ARN (database.password.secret_arn)")
		}
		if v.GetString("database.password.secret_region") == "" {
			log.Fatal("Missing database secret region (database.password.secret_region)")
		}
	case "kubernetes":
		if v.GetString("database.password.secret_name") == "" {
			log.Fatal("Database password not provided as a Kubernetes secret (database.password.secret_name)")
		}
		if v.GetString("database.password.secret_namespace") == "" {
			log.Fatal("Missing database secret namespace (database.password.secret_namespace)")
		}
	default:
		log.Fatalf("Invalid password provider: %s", p)
	}
}
