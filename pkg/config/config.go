package compserv

import (
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
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

	// nolint:gocritic // gocritic suggests to rewrite this as switch-case..nope..
	if len(viper.GetString("database.host")) == 0 {
		log.Fatal("Database host not provided.")
	} else if len(viper.GetString("database.secret_arn")) == 0 {
		log.Fatal("Database password not provided as an ARN.")
	} else if len(viper.GetString("database.username")) == 0 {
		log.Fatal("Database username not provided.")
	} else if len(viper.GetString("database.secret_region")) == 0 {
		log.Fatal("Database secret region not provided.")
	}

	log.Printf("Loaded configuration file: %s", viper.ConfigFileUsed())
	return viper.GetViper()
}

func GetDatabaseConnectionString(v *viper.Viper) string {
	secret := getSecret(v.GetString("database.secret_arn"), v.GetString("database.secret_region"))
	connectionString := fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%s",
		v.GetString("database.host"), v.GetString("database.username"),
		v.GetString("database.name"), secret, v.GetString("database.port"))
	return connectionString
}

// This code comes from AWS Secret Manager
func getSecret(secretName string, region string) string {
	var secretString, decodedBinarySecret string

	// Create a Secrets Manager client
	sess, err := session.NewSession()
	if err != nil {
		// Handle session creation error
		fmt.Println(err.Error())
		return secretString
	}
	svc := secretsmanager.New(sess,
		aws.NewConfig().WithRegion(region))
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	// In this sample we only handle the specific exceptions for the 'GetSecretValue' API.
	// See https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html

	result, err := svc.GetSecretValue(input)
	if err != nil {
		//nolint:errorlint // let's keep borrowed code the way it is
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeDecryptionFailure:
				// Secrets Manager can't decrypt the protected secret text using the provided KMS key.
				fmt.Println(secretsmanager.ErrCodeDecryptionFailure, aerr.Error())

			case secretsmanager.ErrCodeInternalServiceError:
				// An error occurred on the server side.
				fmt.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())

			case secretsmanager.ErrCodeInvalidParameterException:
				// You provided an invalid value for a parameter.
				fmt.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())

			case secretsmanager.ErrCodeInvalidRequestException:
				// You provided a parameter value that is not valid for the current state of the resource.
				fmt.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())

			case secretsmanager.ErrCodeResourceNotFoundException:
				// We can't find the resource that you asked for.
				fmt.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return secretString
	}

	// Decrypts secret using the associated KMS key.
	// Depending on whether the secret is a string or binary, one of these fields will be populated.
	// nolint:golint,revive // let's keep borrowed code the way it is
	if result.SecretString != nil {
		secretString = *result.SecretString
		return secretString
	} else {
		decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(result.SecretBinary)))
		secretLen, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, result.SecretBinary)
		if err != nil {
			fmt.Println("Base64 Decode Error:", err)
			return secretString
		}
		decodedBinarySecret = string(decodedBinarySecretBytes[:secretLen])
		return decodedBinarySecret
	}
}
