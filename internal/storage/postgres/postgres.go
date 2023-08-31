package postgres

import (
	"context"

	db "github.com/NevostruevK/GophKeeper/internal/db/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	*pgxpool.Pool
}

func NewStorage(ctx context.Context, dsn string) (*Storage, error) {
	conn, err := db.NewClient(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return &Storage{conn}, nil
}
