package store

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage interface {
	Querier
}

type SQLStorage struct {
	*Queries
}

func NewStorage(connPool *pgxpool.Pool) Storage {
	return &SQLStorage{
		Queries: New(connPool),
	}
}
