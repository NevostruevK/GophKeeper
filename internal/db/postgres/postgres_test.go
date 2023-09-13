package postgres_test

import (
	"context"
	"testing"

	"github.com/NevostruevK/GophKeeper/internal/db/postgres"
	"github.com/stretchr/testify/assert"
)

func TestPostgres(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	t.Run("New postgres client ok", func(t *testing.T) {
		conn, err := postgres.NewClient(ctx, "user=postgres sslmode=disable")
		assert.NoError(t, err)
		assert.NotNil(t, conn)
		conn.Close()
	})
	t.Run("New postgres client err", func(t *testing.T) {
		conn, err := postgres.NewClient(ctx, "wrong DSN")
		assert.Error(t, err)
		assert.Nil(t, conn)
	})
}
