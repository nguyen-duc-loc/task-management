package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"github.com/nguyen-duc-loc/task-management/backend/util"
)

var (
	database string
	password string
	username string
	port     string
	host     string
	schema   string
)

func NewConnPool() *pgxpool.Pool {
	databaseConfig, err := util.LoadDatabaseConfig()
	if err != nil {
		log.Fatal(err)
	}
	host = databaseConfig.Host
	port = databaseConfig.Port
	database = databaseConfig.Database
	username = databaseConfig.Username
	password = databaseConfig.Password
	schema = databaseConfig.Schema

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", username, password, host, port, database, schema)
	connPool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatal(err)
	}

	return connPool
}
