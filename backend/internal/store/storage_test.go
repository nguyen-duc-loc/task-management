package store

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nguyen-duc-loc/task-management/backend/util"
)

var testStore Storage

func TestMain(m *testing.M) {
	databaseConfig, err := util.LoadDatabaseConfig()
	if err != nil {
		log.Fatal(err)
	}

	host := databaseConfig.Host
	port := databaseConfig.Port
	database := databaseConfig.Database
	username := databaseConfig.Username
	password := databaseConfig.Password
	schema := databaseConfig.Schema

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", username, password, host, port, database, schema)
	connPool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatal(err)
	}

	testStore = NewStorage(connPool)
	os.Exit(m.Run())
}
