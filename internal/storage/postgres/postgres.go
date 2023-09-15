// Postgres storage implementaton.
package postgres

import (
	"context"

	db "github.com/NevostruevK/GophKeeper/internal/db/postgres"
	"github.com/NevostruevK/GophKeeper/internal/tools/crypto"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Storage postgres storage.
type Storage struct {
	*pgxpool.Pool
	crypto *crypto.Crypto
}

// NewStorage returns postgres storage.
func NewStorage(ctx context.Context, dsn string, crypto *crypto.Crypto) (*Storage, error) {
	conn, err := db.NewClient(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return &Storage{conn, crypto}, nil
}
