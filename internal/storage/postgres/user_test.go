package postgres_test

import (
	"context"
	"errors"
	"testing"

	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/NevostruevK/GophKeeper/internal/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorage_AddUser(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("Add user ok", func(t *testing.T) {
		id1, err := addUser(ctx)
		require.NoError(t, err)
		id2, err := addUser(ctx)
		require.NoError(t, err)
		assert.NotEqual(t, id1, id2)
	})
	t.Run("Add the same login error", func(t *testing.T) {
		user := models.NewUser(newLogPass())
		_, err := testStorage.AddUser(ctx, *user)
		require.NoError(t, err)

		id, err := testStorage.AddUser(ctx, models.User{Login: user.Login, Password: "any password"})
		require.Error(t, err)
		assert.True(t, errors.Is(err, storage.ErrDuplicateLogin))
		assert.Equal(t, uuid.Nil, id)
	})
	t.Run("Add short login error", func(t *testing.T) {
		id, err := testStorage.AddUser(ctx, models.User{Login: "short", Password: "any password"})
		require.Error(t, err)
		assert.False(t, errors.Is(err, storage.ErrDuplicateLogin))
		assert.Equal(t, uuid.Nil, id)
	})
}

func TestPostgres_GetUser(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ids := make([]uuid.UUID, 0, 4)
	t.Run("Get user ok", func(t *testing.T) {
		user := models.NewUser("test_login", "test_password")
		idAdd, err := testStorage.AddUser(ctx, *user)
		require.NoError(t, err)
		ids = append(ids, idAdd)
		idGet, err := testStorage.GetUser(ctx, *user)
		require.NoError(t, err)
		assert.Equal(t, idAdd, idGet)
	})
}
