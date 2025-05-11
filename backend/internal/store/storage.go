package store

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage interface {
	Querier
	Health() map[string]string
}

type SQLStorage struct {
	connPool *pgxpool.Pool
	*Queries
}

func NewStorage(connPool *pgxpool.Pool) Storage {
	return &SQLStorage{
		connPool: connPool,
		Queries:  New(connPool),
	}
}
