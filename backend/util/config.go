package util

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
	Schema   string
}

func LoadDatabaseConfig() (databaseConfig DatabaseConfig, err error) {
	godotenv.Load("../../.env")

	database := os.Getenv("DB_DATABASE")
	if len(database) == 0 {
		err = errors.New("Database is not specified")
		return
	}

	password := os.Getenv("DB_PASSWORD")
	if len(password) == 0 {
		err = errors.New("Database password is not specified")
		return
	}

	username := os.Getenv("DB_USERNAME")
	if len(username) == 0 {
		err = errors.New("Database username is not specified")
		return
	}

	port := os.Getenv("DB_PORT")
	if len(port) == 0 {
		err = errors.New("Database port is not specified")
		return
	}

	host := os.Getenv("DB_HOST")
	if len(host) == 0 {
		err = errors.New("Database host is not specified")
		return
	}

	schema := os.Getenv("DB_SCHEMA")
	if len(schema) == 0 {
		err = errors.New("Database schema is not specified")
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
