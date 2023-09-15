// For creating postgres connection.
package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewClient returns postgres connection.
func NewClient(ctx context.Context, dsn string) (conn *pgxpool.Pool, err error) {
	conn, err = pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database, %w", err)
	}
	if err := createTables(ctx, conn); err != nil {
		return nil, err
	}
	return conn, nil
}

func createTables(ctx context.Context, conn *pgxpool.Pool) error {
	for i, query := range []string{createUsersSQL, createRecordsSQL} {
		_, err := conn.Exec(ctx, query)
		if err != nil {
			return fmt.Errorf("failed to create table %d: %w", i, err)
		}
	}
	return nil
}
