package util

import (
	"errors"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load("../../.env")
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
	Schema   string
}

type JWTConfig struct {
	SecretKey           string
	AccessTokenDuration time.Duration
}

func LoadSeverEnv() string {
	serverEnv := os.Getenv("SERVER_ENV")
	if serverEnv != "prod" {
		return "dev"
	}

	return serverEnv
}

func LoadDatabaseConfig() (databaseConfig DatabaseConfig, err error) {
	serverEnv := LoadSeverEnv()

	var password string
	if serverEnv == "dev" {
		password = os.Getenv("DB_PASSWORD")
	} else {
		var secrets Secrets
		secrets, err = getSecrets()
		if err != nil {
			return
		}
		password = secrets.DBPassword
	}
	if len(password) == 0 {
		err = errors.New("database password is not specified")
		return
	}

	database := os.Getenv("DB_DATABASE")
	if len(database) == 0 {
		err = errors.New("database is not specified")
		return
	}

	username := os.Getenv("DB_USERNAME")
	if len(username) == 0 {
		err = errors.New("database username is not specified")
		return
	}

	port := os.Getenv("DB_PORT")
	if len(port) == 0 {
		err = errors.New("database port is not specified")
		return
	}

	host := os.Getenv("DB_HOST")
	if len(host) == 0 {
		err = errors.New("database host is not specified")
		return
	}

	schema := os.Getenv("DB_SCHEMA")
	if len(schema) == 0 {
		err = errors.New("database schema is not specified")
		return
	}

	databaseConfig.Host = host
	databaseConfig.Port = port
	databaseConfig.Database = database
	databaseConfig.Username = username
	databaseConfig.Password = password
	databaseConfig.Schema = schema
	return
}

func LoadJWTConfig() (jwtConfig JWTConfig, err error) {
	serverEnv := os.Getenv("SERVER_ENV")
	if serverEnv != "prod" {
		serverEnv = "dev"
	}

	var secretKey string
	if serverEnv == "dev" {
		secretKey = os.Getenv("JWT_SECRET_KEY")
	} else {
		var secrets Secrets
		secrets, err = getSecrets()
		if err != nil {
			return
		}
		secretKey = secrets.JWTSecretKey
	}
	if len(secretKey) == 0 {
		err = errors.New("JWT secret key is not specified")
		return
	}

	accessTokenDurationEnv := os.Getenv("JWT_ACCESS_TOKEN_DURATION")
	if len(accessTokenDurationEnv) == 0 {
		err = errors.New("access token duration is not specified")
		return
	}
	accessTokenDuration, err := time.ParseDuration(accessTokenDurationEnv)
	if err != nil {
		return
	}

	jwtConfig.SecretKey = secretKey
	jwtConfig.AccessTokenDuration = accessTokenDuration
	return
}
