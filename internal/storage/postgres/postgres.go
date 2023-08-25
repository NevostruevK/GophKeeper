package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	*pgxpool.Pool
}

func NewStorage(client *pgxpool.Pool) *Storage {
	return &Storage{
		client,
	}
}
