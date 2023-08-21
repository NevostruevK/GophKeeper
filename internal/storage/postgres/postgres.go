package postgres

import (
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrDuplicateLogin = errors.New("duplicate login")

type Storage struct {
	*pgxpool.Pool
}

func NewStorage(client *pgxpool.Pool) *Storage {
	return &Storage{
		client,
	}
}
