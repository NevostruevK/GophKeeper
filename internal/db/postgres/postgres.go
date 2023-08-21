package postgres

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewClient(ctx context.Context, dsn string) (conn *pgxpool.Pool, err error) {
	conn, err = pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("failed to connect to database, %v", err)
	}
	if err := createTables(ctx, conn); err != nil {
		log.Fatalf("error: %v", err)
	}
	return
}

func createTables(ctx context.Context, conn *pgxpool.Pool) error {
	for i, query := range []string{createUsersSQL, createMetasSQL, createDatasSQL, createDescriptionsSQL} {
		_, err := conn.Exec(ctx, query)
		if err != nil {
			return fmt.Errorf("failed to create table %d: %w", i, err)
		}
	}
	return nil
}
